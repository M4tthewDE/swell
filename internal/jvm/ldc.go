package jvm

import (
	"context"
	"fmt"

	"github.com/m4tthewde/swell/internal/class"
	"github.com/m4tthewde/swell/internal/jvm/stack"
)

func ldcNormal(r *Runner, ctx context.Context, code []byte) error {
	index := code[r.pc+1]
	r.pc += 2

	return ldc(r, ctx, int(index))
}

func ldcWide(r *Runner, ctx context.Context, code []byte) error {
	index := (uint16(code[r.pc+1])<<8 | uint16(code[r.pc+2]))
	r.pc += 3

	return ldc(r, ctx, int(index))
}

func ldc(r *Runner, ctx context.Context, index int) error {
	pool, err := r.stack.CurrentConstantPool()
	if err != nil {
		return err
	}

	cpInfo, err := pool.Get(int(index))
	if err != nil {
		return err
	}

	if !isLoadable(cpInfo) {
		return fmt.Errorf("%s ist not loadable", cpInfo)
	}

	switch info := cpInfo.(type) {
	case class.ClassInfo:
		className, err := pool.GetUtf8(info.NameIndex)
		if err != nil {
			return err
		}

		c, err := r.loader.Load(ctx, className)
		if err != nil {
			return err
		}

		classClass, err := r.loader.Load(ctx, "java/lang/Class")
		if err != nil {
			return err
		}

		ref, err := r.heap.AllocateObject(ctx, classClass)
		if err != nil {
			return err
		}

		return r.stack.PushOperand(ctx, stack.ClassReferenceValue{Value: ref, Class: c})
	case class.IntegerInfo:
		return r.stack.PushOperand(ctx, stack.IntValue{Value: int32(info.Value)})
	case class.StringInfo:
		stringValue, err := pool.GetUtf8(info.StringIndex)
		if err != nil {
			return err
		}

		c, err := r.loader.Load(ctx, "java/lang/String")
		if err != nil {
			return err
		}

		strID, err := r.heap.AllocateObject(ctx, c)
		if err != nil {
			return err
		}

		byteArray := make([]stack.Value, 0)
		for _, b := range []byte(stringValue) {
			byteArray = append(byteArray, stack.ByteValue{Value: b})
		}

		var value stack.Value
		value = Array{items: byteArray}

		err = r.heap.SetField(*strID, "value", value)
		if err != nil {
			return err
		}

		err = r.heap.SetField(*strID, "coder", stack.ByteValue{Value: 1})
		if err != nil {
			return err
		}

		return r.stack.PushOperand(ctx, stack.ReferenceValue{Value: strID})
	default:
		return fmt.Errorf("ldc not implemented for %s", cpInfo)
	}

}

func isLoadable(cpInfo class.CpInfo) bool {
	switch cpInfo.(type) {
	case class.IntegerInfo:
		return true
	case class.ClassInfo:
		return true
	case class.StringInfo:
		return true
	case class.MethodHandleInfo:
		return true
	case class.MethodTypeInfo:
		return true
	default:
		return false
	}
}
