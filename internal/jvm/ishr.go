package jvm

import (
	"fmt"

	"github.com/m4tthewde/swell/internal/jvm/stack"
)

func intShiftRight(r *Runner) error {
	r.pc += 1
	operands, err := r.stack.PopOperands(2)
	if err != nil {
		return nil
	}

	value1, ok1 := operands[0].(stack.IntValue)
	value2, ok2 := operands[1].(stack.IntValue)

	if ok1 && ok2 {
		result := value1.Value >> (value2.Value & 31)
		return r.stack.PushOperand(stack.IntValue{Value: result})
	}

	return fmt.Errorf("values have to be int, are %v and %v", operands[0], operands[1])
}
