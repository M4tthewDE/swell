package stack

import (
	"testing"

	"github.com/m4tthewde/swell/internal/class"
	"github.com/m4tthewde/swell/internal/logger"
	"github.com/stretchr/testify/assert"
)

func TestStackPushPop(t *testing.T) {
	stack := NewStack()
	stack.Push("Main", class.Method{}, class.ConstantPool{}, []Value{})
	stack.Push("Main2", class.Method{}, class.ConstantPool{}, []Value{})

	assert.Equal(t, 2, len(stack.frames))

	err := stack.Pop()
	assert.Nil(t, err)

	assert.Equal(t, 1, len(stack.frames))
}

func TestStackPushOperand(t *testing.T) {
	log, err := logger.NewLogger()
	assert.Nil(t, err)

	ctx := logger.OnContext(t.Context(), log)

	stack := NewStack()
	stack.Push("Main", class.Method{
		NameIndex: 0,
	}, class.ConstantPool{
		Infos: []class.CpInfo{
			class.Utf8Info{Content: "testMethod"},
		},
	}, []Value{})

	value := BooleanValue{Value: false}
	err = stack.PushOperand(ctx, value)
	assert.Nil(t, err)

	operands, err := stack.PopOperands(1)
	assert.Nil(t, err)
	assert.Equal(t, value, operands[0])
}

func TestStackPopOperand(t *testing.T) {
	log, err := logger.NewLogger()
	assert.Nil(t, err)

	ctx := logger.OnContext(t.Context(), log)

	stack := NewStack()
	stack.Push("Main", class.Method{
		NameIndex: 0,
	}, class.ConstantPool{
		Infos: []class.CpInfo{
			class.Utf8Info{Content: "testMethod"},
		},
	}, []Value{})

	baseType, err := class.NewBaseType(class.BOOLEAN)
	assert.Nil(t, err)

	value, err := DefaultValue(baseType)
	assert.Nil(t, err)

	err = stack.PushOperand(ctx, value)
	assert.Nil(t, err)

	operands, err := stack.PopOperands(1)
	assert.Nil(t, err)
	assert.Equal(t, value, operands[0])

	assert.Panics(t, func() { stack.PopOperands(1) })
}

func TestStackPushOperandInvoker(t *testing.T) {
	log, err := logger.NewLogger()
	assert.Nil(t, err)

	ctx := logger.OnContext(t.Context(), log)

	stack := NewStack()
	stack.Push("Main1", class.Method{
		NameIndex: 0,
	}, class.ConstantPool{
		Infos: []class.CpInfo{
			class.Utf8Info{Content: "testMethod1"},
		},
	}, []Value{})

	stack.Push("Main2", class.Method{
		NameIndex: 0,
	}, class.ConstantPool{
		Infos: []class.CpInfo{
			class.Utf8Info{Content: "testMethod2"},
		},
	}, []Value{})

	value := BooleanValue{Value: false}
	err = stack.PushOperandInvoker(ctx, value)
	assert.Nil(t, err)

	err = stack.Pop()
	assert.Nil(t, err)

	operands, err := stack.PopOperands(1)
	assert.Nil(t, err)
	assert.Equal(t, value, operands[0])
}

func TestStackSetLocalVariable(t *testing.T) {
	stack := NewStack()
	stack.Push("Main", class.Method{}, class.ConstantPool{}, []Value{})

	value := BooleanValue{Value: false}

	log, err := logger.NewLogger()
	assert.Nil(t, err)

	ctx := logger.OnContext(t.Context(), log)

	err = stack.SetLocalVariable(ctx, 0, value)
	assert.Nil(t, err)

	variable, err := stack.GetLocalVariable(ctx, 0)
	assert.Nil(t, err)
	assert.Equal(t, value, variable)
}

/*
func TestStackPopMultipleOperands(t *testing.T) {
	log, err := logger.NewLogger()
	assert.Nil(t, err)

	ctx := logger.OnContext(t.Context(), log)

	stack := NewStack()
	stack.Push("Main", class.Method{}, class.ConstantPool{}, []Value{})

	err = stack.PushOperand(ctx, BooleanValue{Value: false})
	assert.Nil(t, err)

	err = stack.PushOperand(ctx, BooleanValue{Value: true})
	assert.Nil(t, err)

	operands, err := stack.PopOperands(2)
	assert.Nil(t, err)

	assert.Equal(t, []Value{BooleanValue{Value: true}, BooleanValue{Value: false}}, operands)
}
*/
