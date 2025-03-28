package jvm

import (
	"context"
	"errors"
)

func getStatic(r *Runner, ctx context.Context, code []byte) error {
	index := (uint16(code[r.pc+1])<<8 | uint16(code[r.pc+2]))
	r.pc += 3

	pool, err := r.stack.CurrentConstantPool()
	if err != nil {
		return err
	}

	refInfo, err := pool.Ref(index)
	if err != nil {
		return err
	}

	classInfo, err := pool.Class(refInfo.ClassIndex)
	if err != nil {
		return err
	}

	className, err := pool.GetUtf8(classInfo.NameIndex)
	if err != nil {
		return err
	}

	err = r.initializeClass(ctx, className)
	if err != nil {
		return err
	}

	c, err := r.loader.Load(ctx, className)
	if err != nil {
		return err
	}

	nameAndType, err := pool.NameAndType(refInfo.NameAndTypeIndex)
	if err != nil {
		return err
	}

	fieldName, err := pool.GetUtf8(nameAndType.NameIndex)
	if err != nil {
		return err
	}

	_, ok, err := c.GetField(fieldName)
	if err != nil {
		return err
	}

	if !ok {
		return errors.New("static field not found")
	}

	fieldValue, err := r.loader.GetField(className, fieldName)
	if err != nil {
		return err
	}

	return r.stack.PushOperand(ctx, fieldValue)
}
