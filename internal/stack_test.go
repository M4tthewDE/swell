package internal

import (
	"testing"

	"github.com/m4tthewde/swell/internal/class"
	"github.com/stretchr/testify/assert"
)

func TestStack(t *testing.T) {
	value, err := DefaultValue(class.BOOLEAN)
	assert.NotNil(t, err)

	stack := NewStack()
	stack.Push("Main", "main", []Value{value})

	operands := stack.PopOperands(1)
	assert.Equal(t, value, operands[0])

	stack.PushOperand(value)

	operands = stack.PopOperands(2)
	assert.Equal(t, value, operands[0])
	assert.Equal(t, value, operands[1])
}
