package jvm

import (
	"fmt"

	"github.com/m4tthewde/swell/internal/jvm/stack"
)

func ifne(r *Runner, code []byte) error {
	operands, err := r.stack.PopOperands(1)
	if err != nil {
		return err
	}

	if val, ok := operands[0].(stack.IntValue); ok {
		if val.Value != 0 {
			index := (uint16(code[r.pc+1])<<8 | uint16(code[r.pc+2]))
			r.pc += int(index)
			return nil
		}

		r.pc += 3
	}

	return fmt.Errorf("operand has to be int, is %s", operands[0])

}
