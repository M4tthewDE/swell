package jvm

import "errors"

func aload(r *Runner, n int) error {
	r.pc += 1
	variable, err := r.stack.GetLocalVariable(n)
	if err != nil {
		return err
	}

	if reference, ok := variable.(ReferenceValue); ok {
		return r.stack.PushOperand(reference)
	}

	return errors.New("invalid variable type")
}
