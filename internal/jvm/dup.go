package jvm

import "context"

func dup(ctx context.Context, r *Runner) error {
	r.pc += 1

	operand, err := r.stack.GetOperand()
	if err != nil {
		return err
	}

	return r.stack.PushOperand(ctx, operand)
}
