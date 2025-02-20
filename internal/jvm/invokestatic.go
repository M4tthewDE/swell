package jvm

import (
	"context"
	"errors"

	"github.com/m4tthewde/swell/internal/class"
)

func invokeStatic(r *Runner, ctx context.Context, code []byte) error {
	index := (uint16(code[r.pc+1])<<8 | uint16(code[r.pc+2]))
	r.pc += 3

	pool := r.stack.CurrentConstantPool()

	ref, err := pool.Ref(index)
	if err != nil {
		return err
	}

	classInfo, err := pool.Class(ref.ClassIndex)
	if err != nil {
		return err
	}

	className, err := pool.GetUtf8(classInfo.NameIndex)
	if err != nil {
		return err
	}

	if err = r.initializeClass(ctx, className); err != nil {
		return err
	}

	nameAndType, err := pool.NameAndType(ref.NameAndTypeIndex)
	if err != nil {
		return err
	}

	methodName, err := pool.GetUtf8(nameAndType.NameIndex)
	if err != nil {
		return err
	}

	c, err := r.loader.Load(ctx, className)
	if err != nil {
		return err
	}

	method, ok, err := c.GetMethod(methodName)
	if err != nil {
		return err
	}

	if !ok {
		return errors.New("method not found")
	}

	descriptor, err := pool.GetUtf8(nameAndType.DescriptorIndex)
	if err != nil {
		return err
	}

	methodDescriptor, err := class.NewMethodDescriptor(descriptor)
	if err != nil {
		return err
	}

	operands := r.stack.PopOperands(len(methodDescriptor.Parameters))

	if !method.IsNative() {
		code, err := method.CodeAttribute()
		if err != nil {
			return err
		}

		return r.runMethod(ctx, code.Code, *c, methodName, operands)
	} else {
		val, err := r.runNative(ctx, *c, method, operands)
		if err != nil {
			return err
		}

		if val != nil {
			r.stack.PushOperand(val)
		}
	}

	return nil
}
