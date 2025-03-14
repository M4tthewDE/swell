package jvm

import "github.com/m4tthewde/swell/internal/jvm/stack"

func iconst(r *Runner, n int32) error {
	r.pc += 1
	return r.stack.PushOperand(stack.IntValue{Value: n})
}
