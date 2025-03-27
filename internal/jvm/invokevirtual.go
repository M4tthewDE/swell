package jvm

import (
	"context"
	"errors"
	"fmt"

	"github.com/m4tthewde/swell/internal/class"
)

func invokeVirtual(r *Runner, ctx context.Context, code []byte) error {
	index := (uint16(code[r.pc+1])<<8 | uint16(code[r.pc+2]))
	r.pc += 3

	pool, err := r.stack.CurrentConstantPool()
	if err != nil {
		return err
	}

	methodRef, err := pool.Ref(index)
	if err != nil {
		return err
	}

	classInfo, err := pool.Class(methodRef.ClassIndex)
	if err != nil {
		return err
	}

	className, err := pool.GetUtf8(classInfo.NameIndex)
	if err != nil {
		return err
	}

	c, err := r.loader.Load(ctx, className)
	if err != nil {
		return err
	}

	nameAndType, err := pool.NameAndType(methodRef.NameAndTypeIndex)
	if err != nil {
		return err
	}

	methodName, err := pool.GetUtf8(nameAndType.NameIndex)
	if err != nil {
		return err
	}

	descriptorString, err := pool.GetUtf8(nameAndType.DescriptorIndex)
	if err != nil {
		return err
	}

	method, ok, err := c.GetMethod(methodName, descriptorString)
	if err != nil {
		return err
	}

	if !ok {
		return fmt.Errorf("method %s not found in %s", methodName, c.Name)
	}

	methodDescriptor, err := class.NewMethodDescriptor(descriptorString)
	if err != nil {
		return err
	}

	if isSignaturePolymorphic(c, method, methodDescriptor) {
		return errors.New("invokevirtual not implemented for signature polymorphic methods")
	} else {
		// +1 to include the objectref at position 0
		parameters, err := r.stack.PopOperands(len(methodDescriptor.Parameters) + 1)
		if err != nil {
			return err
		}

		codeAttribute, err := method.CodeAttribute()
		if err != nil {
			return err
		}

		return r.runMethod(ctx, codeAttribute, *c, *method, parameters)
	}
}

func isSignaturePolymorphic(c *class.Class, method *class.Method, methodDescriptor *class.MethodDescriptor) bool {
	return isMethodHandleOrVarHandle(c) &&
		hasSingleParamObjectArray(methodDescriptor) &&
		method.IsVarargs() &&
		method.IsNative()
}

func isMethodHandleOrVarHandle(c *class.Class) bool {
	return c.Name == "java/lang/invoke/MethodHandle" || c.Name == "java/lang/invoke/VarHandle"
}

func hasSingleParamObjectArray(methodDescriptor *class.MethodDescriptor) bool {
	if len(methodDescriptor.Parameters) != 1 {
		return false
	}

	param := methodDescriptor.Parameters[0]

	if arrayType, ok := param.(class.ArrayType); ok {
		if objectType, ok := arrayType.FieldType.(class.ObjectType); ok {
			return objectType.ClassName == "java/lang/Object"
		}
	}

	return false
}
