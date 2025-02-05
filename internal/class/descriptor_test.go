package class

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMethodDescriptor(t *testing.T) {
	parameters := []FieldType{
		BaseType(INT),
		BaseType(DOUBLE),
		ObjectType("java/lang/Thread"),
	}

	methodDescriptor, err := NewMethodDescriptor("(IDLjava/lang/Thread;)Ljava/lang/Object;")
	assert.Nil(t, err)
	assert.NotNil(t, methodDescriptor)
	assert.Equal(t, parameters, methodDescriptor.Parameters)
	assert.Equal(t, ObjectType("java/lang/Object"), methodDescriptor.ReturnDescriptor)
}
