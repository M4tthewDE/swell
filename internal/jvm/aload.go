package jvm

import (
	"fmt"

	"github.com/m4tthewde/swell/internal/jvm/stack"
)

func aload(r *Runner, n int) error {
	r.pc += 1
	variable, err := r.stack.GetLocalVariable(n)
	if err != nil {
		return err
	}

	switch val := variable.(type) {
	case stack.ReferenceValue:
		return r.stack.PushOperand(val)
	case stack.ClassReferenceValue:
		return r.stack.PushOperand(val)
	case stack.StringReferenceValue:
		return r.stack.PushOperand(val)
	default:
		return fmt.Errorf("invalid variable type: %s", val)
	}

}
