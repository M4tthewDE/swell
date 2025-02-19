package jvm

import "errors"

func aload(r *Runner, n int) error {
	r.pc += 1
	variable := r.stack.GetLocalVariable(n)

	if reference, ok := variable.(ReferenceValue); ok {
		r.stack.PushOperand(reference)
		return nil
	}

	return errors.New("invalid variable type")
}
