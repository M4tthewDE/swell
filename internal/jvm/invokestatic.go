package jvm

import (
	"context"
	"errors"
	"fmt"

	"github.com/m4tthewde/swell/internal/class"
	"github.com/m4tthewde/swell/internal/jvm/stack"
)

func invokeStatic(r *Runner, ctx context.Context, code []byte) error {
	index := (uint16(code[r.pc+1])<<8 | uint16(code[r.pc+2]))
	r.pc += 3

	pool, err := r.stack.CurrentConstantPool()
	if err != nil {
		return err
	}

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

	descriptor, err := pool.GetUtf8(nameAndType.DescriptorIndex)
	if err != nil {
		return err
	}

	method, ok, err := c.GetMethod(methodName, descriptor)
	if err != nil {
		return err
	}

	if !ok {
		return errors.New("method not found")
	}

	methodDescriptor, err := class.NewMethodDescriptor(descriptor)
	if err != nil {
		return err
	}

	operands, err := r.stack.PopOperands(len(methodDescriptor.Parameters))
	if err != nil {
		return err
	}

	if !method.IsNative() {
		code, err := method.CodeAttribute()
		if err != nil {
			return err
		}

		return r.runMethod(ctx, code, *c, *method, operands)
	} else {
		val, err := runNative(ctx, r, *c, method, operands)
		if err != nil {
			return err
		}

		if val != nil {
			return r.stack.PushOperand(ctx, val)
		}

		return nil
	}
}

func runNative(ctx context.Context, r *Runner, c class.Class, method *class.Method, operands []stack.Value) (stack.Value, error) {
	descriptor, err := c.ConstantPool.GetUtf8(method.DescriptorIndex)
	if err != nil {
		return nil, err
	}

	methodDescriptor, err := class.NewMethodDescriptor(descriptor)
	if err != nil {
		return nil, err
	}

	methodName, err := c.ConstantPool.GetUtf8(method.NameIndex)
	if err != nil {
		return nil, err
	}

	if c.Name == "java/lang/System" &&
		methodName == "registerNatives" &&
		methodDescriptor.ReturnDescriptor == 'V' &&
		len(methodDescriptor.Parameters) == 0 {

		method, ok, err := c.GetMethodByName("initPhase1")
		if err != nil {
			return nil, err
		}

		if !ok {
			return nil, errors.New("method 'initPhase1' not found")
		}

		code, err := method.CodeAttribute()
		if err != nil {
			return nil, err
		}

		return nil, r.runMethod(ctx, code, c, *method, operands)
	} else if c.Name == "java/lang/Class" && methodName == "registerNatives" {
		return nil, nil
	} else if c.Name == "java/lang/Class" && methodName == "desiredAssertionStatus0" {
		return stack.BooleanValue{Value: true}, nil
	} else if c.Name == "java/lang/StringUTF16" && methodName == "isBigEndian" {
		return stack.BooleanValue{Value: true}, nil
	} else {
		return nil, fmt.Errorf("native method %s in %s not implemented", methodName, c.Name)
	}
}
