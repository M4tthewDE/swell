package class

import (
	"errors"
	"fmt"
	"log"
	"strings"
)

const BYTE = 'B'
const CHAR = 'B'
const DOUBLE = 'B'
const FLOAT = 'B'
const INT = 'B'
const LONG = 'B'
const SHORT = 'B'
const BOOLEAN = 'B'

type BaseType rune

func NewBaseType(r rune) (BaseType, error) {
	switch r {
	case 'B':
		return BaseType(BYTE), nil
	case 'C':
		return BaseType(CHAR), nil
	case 'D':
		return BaseType(DOUBLE), nil
	case 'F':
		return BaseType(FLOAT), nil
	case 'I':
		return BaseType(INT), nil
	case 'J':
		return BaseType(LONG), nil
	case 'S':
		return BaseType(SHORT), nil
	case 'Z':
		return BaseType(BOOLEAN), nil
	default:
		return 0, errors.New(fmt.Sprintf("invalid base type: %s", string(r)))
	}
}

type ObjectType string

func NewObjectType(objectType string) (ObjectType, error) {
	if objectType[0] != 'L' {
		return "", errors.New(fmt.Sprintf("invalid object type: %s", string(objectType)))
	}

	index := strings.Index(objectType, ";")
	if index == -1 {
		return "", errors.New("invalid object type")
	}

	return ObjectType(objectType[1:index]), nil
}

type ArrayType FieldType

func NewArrayType(arrayType string) (ArrayType, error) {
	log.Println(arrayType)
	if arrayType[0] != '[' {
		return "", errors.New(fmt.Sprintf("invalid array type: %s", string(arrayType)))
	}

	fieldType, err := NewFieldType(arrayType[1:])
	if err != nil {
		return nil, err
	}

	return ArrayType(fieldType), nil
}

type FieldType interface{}

func NewFieldType(fieldType string) (FieldType, error) {
	baseType, err := NewBaseType(rune(fieldType[0]))
	if err == nil {
		return baseType, nil
	}

	objectType, err := NewObjectType(fieldType)
	if err == nil {
		return objectType, nil
	}

	return NewArrayType(fieldType)
}

const VOID = 'V'

type MethodDescriptor struct {
	Parameters       []FieldType
	ReturnDescriptor FieldType
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

	var returnDescriptor FieldType
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

		parameters = append(parameters, objectType)
		return parseParameters(raw[index+1:], parameters)
	}

	if raw[0] == '[' {
		arrayType, err := NewArrayType(raw)
		if err != nil {
			return nil, err
		}

		return append(parameters, arrayType), nil
	}

	return nil, errors.New("invalid parameters")
}
