package jvm

import (
	"fmt"

	"github.com/m4tthewde/swell/internal/jvm/stack"
)

func ifnonnull(r *Runner, code []byte) error {
	operands, err := r.stack.PopOperands(1)
	if err != nil {
		return err
	}

	switch objectRef := operands[0].(type) {
	case stack.ReferenceValue:
		if objectRef.IsNull() {
			r.pc += 3
			return nil
		} else {
			index := (uint16(code[r.pc+1])<<8 | uint16(code[r.pc+2]))
			r.pc += int(index)
			return nil
		}
	default:
		return fmt.Errorf("invalid operand type: %s", operands[0])
	}
}
