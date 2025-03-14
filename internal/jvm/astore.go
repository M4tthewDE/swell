package jvm

import (
	"fmt"

	"github.com/m4tthewde/swell/internal/jvm/stack"
)

func astore(r *Runner, n int) error {
	operands, err := r.stack.PopOperands(n)
	if err != nil {
		return err
	}

	if objectref, ok := operands[0].(stack.ReferenceValue); ok {
		r.stack.SetLocalVariable(n, objectref)
		r.pc += 1
		return nil
	}

	// FIXME: can also be a return address
	return fmt.Errorf("operand has to be reference, is %s", operands[0])
}
