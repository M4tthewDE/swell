package jvm

import (
	"fmt"

	"github.com/m4tthewde/swell/internal/jvm/stack"
)

func istore(r *Runner, n int) error {
	r.pc += 1

	operands, err := r.stack.PopOperands(1)
	if err != nil {
		return err
	}

	if value, ok := operands[0].(stack.IntValue); ok {
		return r.stack.SetLocalVariable(n, value)
	}

	return fmt.Errorf("value has to be integer, is %v", operands[0])
}
