package jvm

import (
	"context"
	"fmt"

	"github.com/m4tthewde/swell/internal/jvm/stack"
)

func isub(ctx context.Context, r *Runner) error {
	r.pc += 1
	operands, err := r.stack.PopOperands(2)
	if err != nil {
		return nil
	}

	value1, ok1 := operands[0].(stack.IntValue)
	value2, ok2 := operands[1].(stack.IntValue)

	if ok1 && ok2 {
		result := value1.Value - value2.Value
		return r.stack.PushOperand(ctx, stack.IntValue{Value: result})
	}

	return fmt.Errorf("values have to be int, are %v and %v", operands[0], operands[1])
}
