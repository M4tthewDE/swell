package jvm

import (
	"context"
	"fmt"

	"github.com/m4tthewde/swell/internal/jvm/stack"
)

func astore(ctx context.Context, r *Runner, n int) error {
	r.pc += 1
	operands, err := r.stack.PopOperands(1)
	if err != nil {
		return err
	}

	if objectref, ok := operands[0].(stack.ReferenceValue); ok {
		r.stack.SetLocalVariable(ctx, n, objectref)
		return nil
	}

	// FIXME: can also be a return address
	return fmt.Errorf("operand has to be reference, is %s", operands[0])
}
