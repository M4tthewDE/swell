package class

import (
	"bufio"
	"errors"
	"log"
)

type AttributeInfo interface{}

type Attribute struct {
	nameIndex     uint16
	attributeInfo AttributeInfo
}

func NewAttributes(reader *bufio.Reader, count uint16, cp *ConstantPool) ([]Attribute, error) {
	attributes := make([]Attribute, count)
	for i := range count {
		attribute, err := NewAttribute(reader, cp)
		if err != nil {
			return nil, err
		}

		attributes[i] = *attribute
	}

	return attributes, nil
}

func NewAttribute(reader *bufio.Reader, cp *ConstantPool) (*Attribute, error) {
	nameIndex, err := readUint16(reader)
	if err != nil {
		return nil, err
	}

	_, err = readUint32(reader)
	if err != nil {
		return nil, err
	}

	name, err := cp.GetUtf8(nameIndex)
	if err != nil {
		return nil, err
	}

	var info AttributeInfo
	log.Println("parsing attribute " + name)
	switch name {
	case "Code":
		info, err = NewCodeInfo(reader, cp)
	default:
		return nil, errors.New("unknown attribute: " + name)
	}

	if err != nil {
		return nil, err
	}

	return &Attribute{nameIndex: nameIndex, attributeInfo: info}, nil
}

type CodeAttributeInfo struct {
	maxStack   uint16
	maxLocals  uint16
	code       []byte
	attributes []Attribute
}

func NewCodeInfo(reader *bufio.Reader, cp *ConstantPool) (AttributeInfo, error) {
	maxStack, err := readUint16(reader)
	if err != nil {
		return nil, err
	}

	maxLocals, err := readUint16(reader)
	if err != nil {
		return nil, err
	}

	codeLength, err := readUint32(reader)
	if err != nil {
		return nil, err
	}

	code := make([]byte, codeLength)
	_, err = reader.Read(code)

	exceptionTableLength, err := readUint16(reader)
	if err != nil {
		return nil, err
	}

	if exceptionTableLength != 0 {
		return nil, errors.New("exceptions not supported")
	}

	attributesCount, err := readUint16(reader)
	if err != nil {
		return nil, err
	}

	attributes, err := NewAttributes(reader, attributesCount, cp)
	if err != nil {
		return nil, err
	}

	return CodeAttributeInfo{
		maxStack:   maxStack,
		maxLocals:  maxLocals,
		code:       code,
		attributes: attributes,
	}, nil
}
