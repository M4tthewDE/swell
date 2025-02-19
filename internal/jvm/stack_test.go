package jvm

import (
	"testing"

	"github.com/m4tthewde/swell/internal/class"
	"github.com/stretchr/testify/assert"
)

func TestStackPush(t *testing.T) {
	stack := NewStack()
	stack.Push("Main", "main", []Value{})

	value, err := DefaultValue(class.BOOLEAN)
	assert.NotNil(t, err)

	stack.PushOperand(value)

	operands := stack.PopOperands(1)
	assert.Equal(t, value, operands[0])
}

func TestStackPop(t *testing.T) {
	stack := NewStack()
	stack.Push("Main", "main", []Value{})

	value, err := DefaultValue(class.BOOLEAN)
	assert.NotNil(t, err)
	stack.PushOperand(value)

	operands := stack.PopOperands(1)
	assert.Equal(t, value, operands[0])

	assert.Panics(t, func() { stack.PopOperands(1) })
}
