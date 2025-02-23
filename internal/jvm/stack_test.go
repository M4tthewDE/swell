package jvm

import (
	"testing"

	"github.com/m4tthewde/swell/internal/class"
	"github.com/stretchr/testify/assert"
)

func TestStackPushPop(t *testing.T) {
	stack := NewStack()
	stack.Push("Main", class.Method{}, class.ConstantPool{}, []Value{}, []Value{})
	stack.Push("Main2", class.Method{}, class.ConstantPool{}, []Value{}, []Value{})

	assert.Equal(t, 2, len(stack.frames))

	err := stack.Pop()
	assert.Nil(t, err)

	assert.Equal(t, 1, len(stack.frames))
}

func TestStackPushOperand(t *testing.T) {
	stack := NewStack()
	stack.Push("Main", class.Method{}, class.ConstantPool{}, []Value{}, []Value{})

	value := BooleanValue{value: false}
	err := stack.PushOperand(BooleanValue{value: false})
	assert.Nil(t, err)

	operands, err := stack.PopOperands(1)
	assert.Nil(t, err)
	assert.Equal(t, value, operands[0])
}

func TestStackPopOperand(t *testing.T) {
	stack := NewStack()
	stack.Push("Main", class.Method{}, class.ConstantPool{}, []Value{}, []Value{})

	value, err := DefaultValue(class.BOOLEAN)
	assert.NotNil(t, err)

	err = stack.PushOperand(value)
	assert.Nil(t, err)

	operands, err := stack.PopOperands(1)
	assert.Nil(t, err)
	assert.Equal(t, value, operands[0])

	assert.Panics(t, func() { stack.PopOperands(1) })
}
