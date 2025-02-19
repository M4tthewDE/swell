package jvm

import "context"

func new(r *Runner, ctx context.Context, code []byte) error {
	index := (uint16(code[r.pc+1])<<8 | uint16(code[r.pc+2]))
	r.pc += 3

	class, err := r.currentClass.ConstantPool.Class(index)
	if err != nil {
		return err
	}

	className, err := r.currentClass.ConstantPool.GetUtf8(class.NameIndex)
	if err != nil {
		return err
	}

	oldClass := r.currentClass
	err = r.initializeClass(ctx, className)
	if err != nil {
		return err
	}
	r.currentClass = oldClass

	c, err := r.loader.Load(ctx, className)
	if err != nil {
		return err
	}

	id, err := r.heap.Allocate(c)
	if err != nil {
		return err
	}

	r.stack.PushOperand(ReferenceValue{value: id})

	return nil
}
