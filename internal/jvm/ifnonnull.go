package jvm

import (
	"errors"
	"fmt"

	"github.com/m4tthewde/swell/internal/jvm/stack"
)

func ifnonnull(r *Runner) error {
	operands, err := r.stack.PopOperands(1)
	if err != nil {
		return err
	}

	if objectref, ok := operands[0].(stack.ReferenceValue); ok {
		if objectref.IsNull() {
			r.pc += 3
			return nil
		} else {
			return errors.New("not implemented: ifnonnull jump")
		}

	}

	return fmt.Errorf("operand has to be reference, is %s", operands[0])
}
