package stack

import (
	"testing"

	"github.com/m4tthewde/swell/internal/class"
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
	stack := NewStack()
	stack.Push("Main", class.Method{}, class.ConstantPool{}, []Value{})

	value := BooleanValue{Value: false}
	err := stack.PushOperand(value)
	assert.Nil(t, err)

	operands, err := stack.PopOperands(1)
	assert.Nil(t, err)
	assert.Equal(t, value, operands[0])
}

func TestStackPopOperand(t *testing.T) {
	stack := NewStack()
	stack.Push("Main", class.Method{}, class.ConstantPool{}, []Value{})

	baseType, err := class.NewBaseType(class.BOOLEAN)
	assert.Nil(t, err)

	value, err := DefaultValue(baseType)
	assert.NotNil(t, err)

	err = stack.PushOperand(value)
	assert.Nil(t, err)

	operands, err := stack.PopOperands(1)
	assert.Nil(t, err)
	assert.Equal(t, value, operands[0])

	assert.Panics(t, func() { stack.PopOperands(1) })
}

func TestStackPushOperandInvoker(t *testing.T) {
	stack := NewStack()
	stack.Push("Main", class.Method{}, class.ConstantPool{}, []Value{})
	stack.Push("Main2", class.Method{}, class.ConstantPool{}, []Value{})

	value := BooleanValue{Value: false}
	err := stack.PushOperandInvoker(value)
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
	err := stack.SetLocalVariable(0, value)
	assert.Nil(t, err)

	variable, err := stack.GetLocalVariable(0)
	assert.Nil(t, err)
	assert.Equal(t, value, variable)
}
