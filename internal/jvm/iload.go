package jvm

import (
	"context"
	"fmt"

	"github.com/m4tthewde/swell/internal/jvm/stack"
)

func iload(ctx context.Context, r *Runner, n int) error {
	r.pc += 1

	localVariable, err := r.stack.GetLocalVariable(ctx, n)
	if err != nil {
		return err
	}

	if value, ok := localVariable.(stack.IntValue); ok {
		return r.stack.PushOperand(value)
	}

	return fmt.Errorf("value has to be integer, is %v", localVariable)

}
