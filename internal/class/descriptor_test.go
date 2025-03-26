package class

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMethodDescriptor(t *testing.T) {
	parameters := []FieldType{
		BaseType(INT),
		BaseType(DOUBLE),
		ObjectType{ClassName: "java/lang/Thread"},
	}

	methodDescriptor, err := NewMethodDescriptor("(IDLjava/lang/Thread;)Ljava/lang/Object;")
	assert.Nil(t, err)
	assert.NotNil(t, methodDescriptor)
	assert.Equal(t, parameters, methodDescriptor.Parameters)
	assert.Equal(t, ObjectType{ClassName: "java/lang/Object"}, methodDescriptor.ReturnDescriptor)
}

func TestMethodDescriptorArray(t *testing.T) {
	parameters := []FieldType{
		ArrayType{FieldType: ObjectType{ClassName: "java/lang/String"}},
	}

	methodDescriptor, err := NewMethodDescriptor("([Ljava/lang/String;)V")
	assert.Nil(t, err)
	assert.NotNil(t, methodDescriptor)
	assert.Equal(t, parameters, methodDescriptor.Parameters)
	assert.Equal(t, 'V', methodDescriptor.ReturnDescriptor)
}

func TestMethodDescriptorArrayAndParameter(t *testing.T) {
	parameters := []FieldType{
		ArrayType{FieldType: BaseType(BYTE)},
		BaseType(INT),
	}

	methodDescriptor, err := NewMethodDescriptor("([BI)C")
	assert.Nil(t, err)
	assert.NotNil(t, methodDescriptor)
	assert.Equal(t, parameters, methodDescriptor.Parameters)
	assert.Equal(t, BaseType(CHAR), methodDescriptor.ReturnDescriptor)
}
