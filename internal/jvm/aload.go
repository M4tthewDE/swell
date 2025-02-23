package jvm

import (
	"errors"

	"github.com/m4tthewde/swell/internal/jvm/stack"
)

func aload(r *Runner, n int) error {
	r.pc += 1
	variable, err := r.stack.GetLocalVariable(n)
	if err != nil {
		return err
	}

	if reference, ok := variable.(stack.ReferenceValue); ok {
		return r.stack.PushOperand(reference)
	}

	return errors.New("invalid variable type")
}
