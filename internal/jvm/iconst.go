package jvm

import (
	"context"

	"github.com/m4tthewde/swell/internal/jvm/stack"
)

func iconst(ctx context.Context, r *Runner, n int32) error {
	r.pc += 1
	return r.stack.PushOperand(ctx, stack.IntValue{Value: n})
}
