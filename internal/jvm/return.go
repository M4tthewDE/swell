package jvm

import (
	"context"
	"fmt"

	"github.com/m4tthewde/swell/internal/class"
	"github.com/m4tthewde/swell/internal/jvm/stack"
)

func ret(r *Runner) error {
	pool, err := r.stack.CurrentConstantPool()
	if err != nil {
		return err
	}

	method, err := r.stack.CurrentMethod()
	if err != nil {
		return err
	}

	descriptor, err := pool.GetUtf8(method.DescriptorIndex)
	if err != nil {
		return err
	}

	methodDescriptor, err := class.NewMethodDescriptor(descriptor)
	if err != nil {
		return err
	}

	if methodDescriptor.ReturnDescriptor != 'V' {
		return fmt.Errorf("method has to be void, is %s", methodDescriptor.ReturnDescriptor)
	}

	return nil
}

func areturn(ctx context.Context, r *Runner) error {
	operands, err := r.stack.PopOperands(1)
	if err != nil {
		return err
	}

	if objectref, ok := operands[0].(stack.ReferenceValue); ok {
		// FIXME: check assignment compatibility according with JLS §5.2
		return r.stack.PushOperandInvoker(ctx, objectref)
	}

	return fmt.Errorf("operand has to be reference, is %s", operands[0])
}

func ireturn(ctx context.Context, r *Runner) error {
	operands, err := r.stack.PopOperands(1)
	if err != nil {
		return err
	}

	if bool, ok := operands[0].(stack.BooleanValue); ok {
		var val int32
		if bool.Value {
			val = 1
		} else {
			val = 0
		}

		return r.stack.PushOperandInvoker(ctx, stack.IntValue{Value: val})
	}

	if byteValue, ok := operands[0].(stack.ByteValue); ok {
		return r.stack.PushOperandInvoker(ctx, stack.IntValue{Value: int32(byteValue.Value)})
	}

	if intValue, ok := operands[0].(stack.IntValue); ok {
		return r.stack.PushOperandInvoker(ctx, intValue)
	}

	return fmt.Errorf("operand has to be int, is %s", operands[0])
}
