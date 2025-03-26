package jvm

import (
	"context"
	"fmt"

	"github.com/m4tthewde/swell/internal/jvm/stack"
)

func aload(ctx context.Context, r *Runner, n int) error {
	r.pc += 1
	variable, err := r.stack.GetLocalVariable(ctx, n)
	if err != nil {
		return err
	}

	switch val := variable.(type) {
	case stack.ReferenceValue:
		return r.stack.PushOperand(val)
	case stack.ClassReferenceValue:
		return r.stack.PushOperand(val)
	case Array:
		return r.stack.PushOperand(val)
	default:
		return fmt.Errorf("invalid variable type: %s", val)
	}

}
