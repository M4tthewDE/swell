package jvm

import (
	"context"

	"github.com/m4tthewde/swell/internal/class"
)

func invokeSpecial(r *Runner, ctx context.Context, code []byte) error {
	index := (uint16(code[r.pc+1])<<8 | uint16(code[r.pc+2]))
	r.pc += 3

	pool := r.stack.CurrentConstantPool()

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

	methodName, err := pool.GetUtf8(nameAndType.NameIndex)
	if err != nil {
		return err
	}

	method, ok, err := c.GetMethod(methodName)
	if err != nil {
		return err
	}

	if !ok {
		return err
	}

	descriptor, err := pool.GetUtf8(nameAndType.DescriptorIndex)
	if err != nil {
		return err
	}

	methodDescriptor, err := class.NewMethodDescriptor(descriptor)
	if err != nil {
		return err
	}

	codeAttribute, err := method.CodeAttribute()
	if err != nil {
		return err
	}

	operands := r.stack.PopOperands(1)
	operands = append(operands, r.stack.PopOperands(len(methodDescriptor.Parameters))...)

	return r.runMethod(ctx, codeAttribute.Code, *c, *method, operands)
}
