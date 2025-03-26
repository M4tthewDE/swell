package class

import (
	"errors"
	"fmt"
	"strings"
)

const BYTE = 'B'
const CHAR = 'C'
const DOUBLE = 'D'
const FLOAT = 'F'
const INT = 'I'
const LONG = 'J'
const SHORT = 'S'
const BOOLEAN = 'Z'

type BaseType rune

func (b BaseType) isFieldType() {}
func (b BaseType) length() int {
	return 1
}

func (b BaseType) String() string {
	return fmt.Sprintf("BaseType[%s]", string(b))
}

func NewBaseType(r rune) (BaseType, error) {
	switch r {
	case BYTE:
		return BaseType(BYTE), nil
	case CHAR:
		return BaseType(CHAR), nil
	case DOUBLE:
		return BaseType(DOUBLE), nil
	case FLOAT:
		return BaseType(FLOAT), nil
	case INT:
		return BaseType(INT), nil
	case LONG:
		return BaseType(LONG), nil
	case SHORT:
		return BaseType(SHORT), nil
	case BOOLEAN:
		return BaseType(BOOLEAN), nil
	default:
		return 0, fmt.Errorf("invalid base type: %s", string(r))
	}
}

type ObjectType struct {
	ClassName string
}

func (o ObjectType) isFieldType() {}
func (o ObjectType) length() int {
	return len(o.ClassName) + 2
}

func (o ObjectType) String() string {
	return fmt.Sprintf("ObjectType[%s]", o.ClassName)
}

func NewObjectType(objectType string) (*ObjectType, error) {
	if objectType[0] != 'L' {
		return nil, fmt.Errorf("invalid object type: %s", string(objectType))
	}

	index := strings.Index(objectType, ";")
	if index == -1 {
		return nil, errors.New("invalid object type")
	}

	return &ObjectType{ClassName: objectType[1:index]}, nil
}

type ArrayType struct {
	FieldType FieldType
}

func (a ArrayType) isFieldType() {}

func (a ArrayType) length() int {
	return a.FieldType.length()
}

func NewArrayType(arrayType string, endIndex int) (*ArrayType, error) {
	if arrayType[0] != '[' {
		return nil, fmt.Errorf("invalid array type: %s", string(arrayType))
	}

	fieldType, err := NewFieldType(arrayType[1:])
	if err != nil {
		return nil, err
	}

	return &ArrayType{FieldType: fieldType}, nil
}

type FieldType interface {
	isFieldType()
	length() int
}

func NewFieldType(rawFieldType string) (FieldType, error) {
	baseType, err := NewBaseType(rune(rawFieldType[0]))
	if err == nil {
		return baseType, nil
	}

	objectType, err := NewObjectType(rawFieldType)
	if err == nil {
		return *objectType, nil
	}

	arrayType, err := NewArrayType(rawFieldType, 0)
	if err != nil {
		return nil, err
	}

	return *arrayType, nil
}

const VOID = 'V'

type MethodDescriptor struct {
	Parameters []FieldType
	// FIXME: this should not be interface{}
	ReturnDescriptor interface{}
}

func NewMethodDescriptor(methodDescriptor string) (*MethodDescriptor, error) {
	if methodDescriptor[0] != '(' {
		return nil, errors.New("invalid method descriptor")
	}

	parameters := make([]FieldType, 0)
	parameters, err := parseParameters(methodDescriptor[1:], parameters)
	if err != nil {
		return nil, err
	}

	index := strings.Index(methodDescriptor, ")")
	if index == -1 {
		return nil, errors.New("invalid method descriptor")
	}

	var returnDescriptor interface{}
	if methodDescriptor[index+1] == 'V' {
		returnDescriptor = VOID
	} else {
		returnDescriptor, err = NewFieldType(methodDescriptor[index+1:])
		if err != nil {
			return nil, err
		}
	}

	return &MethodDescriptor{Parameters: parameters, ReturnDescriptor: returnDescriptor}, nil
}

func parseParameters(raw string, parameters []FieldType) ([]FieldType, error) {
	if raw[0] == ')' {
		return parameters, nil
	}

	baseType, err := NewBaseType(rune(raw[0]))
	if err == nil {
		parameters = append(parameters, baseType)
		return parseParameters(raw[1:], parameters)
	}

	if raw[0] == 'L' {
		index := strings.Index(raw, ";")
		if index == -1 {
			return nil, errors.New("invalid object type")
		}

		objectType, err := NewObjectType(raw[:index+1])
		if err != nil {
			return nil, err
		}

		parameters = append(parameters, *objectType)
		return parseParameters(raw[index+1:], parameters)
	}

	if raw[0] == '[' {
		arrayType, err := NewArrayType(raw, 0)
		if err != nil {
			return nil, err
		}

		parameters = append(parameters, *arrayType)
		return parseParameters(raw[arrayType.length()+1:], parameters)
	}

	return nil, errors.New("invalid parameters")
}
