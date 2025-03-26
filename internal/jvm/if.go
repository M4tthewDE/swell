package jvm

import (
	"fmt"

	"github.com/m4tthewde/swell/internal/jvm/stack"
)

func value(r *Runner) (int, error) {
	operands, err := r.stack.PopOperands(1)
	if err != nil {
		return 0, err
	}

	if val, ok := operands[0].(stack.IntValue); ok {
		return int(val.Value), nil
	}

	if val, ok := operands[0].(stack.ByteValue); ok {
		return int(val.Value), nil
	}

	if val, ok := operands[0].(stack.BooleanValue); ok {
		if val.Value {
			return 1, nil
		} else {
			return 0, nil
		}
	}

	return 0, fmt.Errorf("operand has to be int, is %s", operands[0])

}

func jump(r *Runner, code []byte) {
	index := (uint16(code[r.pc+1])<<8 | uint16(code[r.pc+2]))
	r.pc += int(index)
}

func cont(r *Runner) {
	r.pc += 3
}

func ifne(r *Runner, code []byte) error {
	val, err := value(r)
	if err != nil {
		return err
	}

	if val != 0 {
		jump(r, code)
		return nil
	}

	cont(r)
	return nil
}

func ifeq(r *Runner, code []byte) error {
	val, err := value(r)
	if err != nil {
		return err
	}

	if val == 0 {
		jump(r, code)
		return nil
	}

	cont(r)
	return nil
}

func iflt(r *Runner, code []byte) error {
	val, err := value(r)
	if err != nil {
		return err
	}

	if val < 0 {
		jump(r, code)
		return nil
	}

	cont(r)
	return nil
}

func ifICmpLt(r *Runner, code []byte) error {
	val1, err := value(r)
	if err != nil {
		return err
	}

	val2, err := value(r)
	if err != nil {
		return err
	}

	if val1 < val2 {
		jump(r, code)
		return nil
	}

	cont(r)
	return nil
}
