package jvm

import (
	"context"
	"fmt"

	"github.com/m4tthewde/swell/internal/class"
	"github.com/m4tthewde/swell/internal/jvm/stack"
)

func anewarray(r *Runner, ctx context.Context, code []byte) error {
	index := (uint16(code[r.pc+1])<<8 | uint16(code[r.pc+2]))
	r.pc += 3

	operands, err := r.stack.PopOperands(1)
	if err != nil {
		return err
	}

	operand := operands[0]

	if o, ok := operand.(stack.IntValue); !ok {
		return fmt.Errorf("count has to be int, is  %s", o)
	}

	count := operand.(stack.IntValue)

	pool, err := r.stack.CurrentConstantPool()
	if err != nil {
		return err
	}

	cpInfo, err := pool.Get(int(index))
	if err != nil {
		return err
	}

	switch info := cpInfo.(type) {
	case class.ClassInfo:
		className, err := pool.GetUtf8(info.NameIndex)
		if err != nil {
			return err
		}

		_, err = r.loader.Load(ctx, className)
		if err != nil {
			return err
		}

		id, err := r.heap.AllocateDefaultArray(ctx, int(count.Value), stack.ReferenceValue{Value: nil})
		if err != nil {
			return err
		}

		return r.stack.PushOperand(ctx, stack.ReferenceValue{Value: id})
	default:
		return fmt.Errorf("anewarray not implemented for %s", cpInfo)
	}
}
