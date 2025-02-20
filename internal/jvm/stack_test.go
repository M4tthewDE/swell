package jvm

import (
	"testing"

	"github.com/m4tthewde/swell/internal/class"
	"github.com/stretchr/testify/assert"
)

func TestStackPushPop(t *testing.T) {
	stack := NewStack()
	stack.Push("Main", "main", class.ConstantPool{}, []Value{}, []Value{})
	stack.Push("Main2", "main2", class.ConstantPool{}, []Value{}, []Value{})

	assert.Equal(t, 2, len(stack.frames))

	stack.Pop()

	assert.Equal(t, 1, len(stack.frames))
}

func TestStackPushOperand(t *testing.T) {
	stack := NewStack()
	stack.Push("Main", "main", class.ConstantPool{}, []Value{}, []Value{})

	value, err := DefaultValue(class.BOOLEAN)
	assert.NotNil(t, err)

	stack.PushOperand(value)

	operands := stack.PopOperands(1)
	assert.Equal(t, value, operands[0])
}

func TestStackPopOperand(t *testing.T) {
	stack := NewStack()
	stack.Push("Main", "main", class.ConstantPool{}, []Value{}, []Value{})

	value, err := DefaultValue(class.BOOLEAN)
	assert.NotNil(t, err)
	stack.PushOperand(value)

	operands := stack.PopOperands(1)
	assert.Equal(t, value, operands[0])

	assert.Panics(t, func() { stack.PopOperands(1) })
}
