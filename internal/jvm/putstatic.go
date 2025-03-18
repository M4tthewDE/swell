package jvm

import (
	"context"
)

func putstatic(r *Runner, ctx context.Context, code []byte) error {
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

	fieldName, err := pool.GetUtf8(nameAndType.NameIndex)
	if err != nil {
		return err
	}

	err = r.initializeClass(ctx, className)
	if err != nil {
		return err
	}

	operands, err := r.stack.PopOperands(1)
	if err != nil {
		return err
	}

	value := operands[0]

	return r.loader.SetField(className, fieldName, value)
}
