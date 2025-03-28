package jvm

import (
	"context"
	"fmt"

	"github.com/m4tthewde/swell/internal/jvm/stack"
)

func arrayLength(ctx context.Context, r *Runner) error {
	r.pc += 1
	operands, err := r.stack.PopOperands(1)
	if err != nil {
		return err
	}

	if array, ok := operands[0].(Array); ok {
		return r.stack.PushOperand(ctx, stack.IntValue{Value: int32(len(array.items))})
	}

	return fmt.Errorf("has to be array, is %v", operands[0])
}
