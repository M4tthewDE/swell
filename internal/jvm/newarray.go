package jvm

import (
	"context"
	"fmt"

	"github.com/m4tthewde/swell/internal/jvm/stack"
)

func newArray(r *Runner, ctx context.Context, code []byte) error {
	aType := code[r.pc+1]
	r.pc += 2

	var defaultValue stack.Value
	switch aType {
	case 8:
		defaultValue = stack.ByteValue{Value: 0}
	default:
		return fmt.Errorf("invalid atype: %v", aType)
	}

	operands, err := r.stack.PopOperands(1)
	if err != nil {
		return err
	}

	operand := operands[0]

	if o, ok := operand.(stack.IntValue); !ok {
		return fmt.Errorf("count has to be int, is  %s", o)
	}

	count := operand.(stack.IntValue)

	if count.Value < 0 {
		return fmt.Errorf("count has to be >= 0, is  %s", count)
	}

	id, err := r.heap.AllocateDefaultArray(ctx, int(count.Value), defaultValue)
	if err != nil {
		return err
	}

	return r.stack.PushOperand(ctx, stack.ReferenceValue{Value: id})
}
