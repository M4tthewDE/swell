package jvm

import (
	"context"
	"errors"
	"fmt"

	"github.com/m4tthewde/swell/internal/class"
	"github.com/m4tthewde/swell/internal/jvm/stack"
)

func putField(r *Runner, ctx context.Context, code []byte) error {
	index := (uint16(code[r.pc+1])<<8 | uint16(code[r.pc+2]))
	r.pc += 3

	pool, err := r.stack.CurrentConstantPool()
	if err != nil {
		return err
	}

	fieldRef, err := pool.Ref(index)
	if err != nil {
		return err
	}

	classInfo, err := pool.Class(fieldRef.ClassIndex)
	if err != nil {
		return err
	}

	className, err := pool.GetUtf8(classInfo.NameIndex)
	if err != nil {
		return err
	}

	_, err = r.loader.Load(ctx, className)
	if err != nil {
		return err
	}

	nameAndType, err := pool.NameAndType(fieldRef.NameAndTypeIndex)
	if err != nil {
		return err
	}

	descriptorString, err := pool.GetUtf8(nameAndType.DescriptorIndex)
	if err != nil {
		return err
	}

	fieldType, err := class.NewFieldType(descriptorString)
	if err != nil {
		return err
	}

	fieldName, err := pool.GetUtf8(nameAndType.NameIndex)
	if err != nil {
		return err
	}

	operands, err := r.stack.PopOperands(2)
	if err != nil {
		return err
	}

	value := operands[0]
	objectRef := operands[1]

	if !isCompatible(fieldType, value) {
		return fmt.Errorf("field type %v is incompatible with value %v", fieldType, value)
	}

	if objectRef, ok := objectRef.(stack.ReferenceValue); ok {
		return r.heap.SetField(*objectRef.Value, fieldName, value)
	}

	return errors.New("objectref has to be a reference")
}

func isCompatible(fieldType class.FieldType, value stack.Value) bool {
	switch fieldType.(type) {
	case class.ObjectType:
		_, ok := value.(stack.ReferenceValue)
		return ok
	default:
		return false
	}
}
