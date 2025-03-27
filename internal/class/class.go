package class

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
)

var MAGIC = []byte{0xCA, 0xFE, 0xBA, 0xBE}

type Class struct {
	Name         string       `json:"name"`
	ConstantPool ConstantPool `json:"constant_pool"`
	Methods      []Method     `json:"methods"`
	Interfaces   []uint16     `json:"interfaces"`
	Fields       []Field      `json:"fields"`
	Attributes   []Attribute  `json:"attributes"`
}

func (c *Class) GetMainMethod() (*Method, bool, error) {
	for _, m := range c.Methods {
		isMain, err := m.IsMain(&c.ConstantPool)
		if err != nil {
			return nil, false, err
		}

		if isMain {
			return &m, true, nil
		}
	}

	return nil, false, nil
}

func (c *Class) GetMethod(methodName string, descriptor string) (*Method, bool, error) {
	for _, m := range c.Methods {
		name, err := c.ConstantPool.GetUtf8(m.NameIndex)
		if err != nil {
			return nil, false, err
		}

		methodDescriptor, err := c.ConstantPool.GetUtf8(m.DescriptorIndex)
		if err != nil {
			return nil, false, err
		}

		if name == methodName && methodDescriptor == descriptor {
			return &m, true, nil
		}
	}

	return nil, false, nil
}

func (c *Class) GetMethodByName(methodName string) (*Method, bool, error) {
	for _, m := range c.Methods {
		name, err := c.ConstantPool.GetUtf8(m.NameIndex)
		if err != nil {
			return nil, false, err
		}

		if name == methodName {
			return &m, true, nil
		}
	}

	return nil, false, nil
}

// TODO: name is not enough to find the correct field
// will have to use descriptor in the future
func (c *Class) GetField(fieldName string) (*Field, bool, error) {
	for _, f := range c.Fields {
		name, err := c.ConstantPool.GetUtf8(f.NameIndex)
		if err != nil {
			return nil, false, err
		}

		if name == fieldName {
			return &f, true, nil
		}
	}

	return nil, false, nil
}

func NewClass(ctx context.Context, reader *bufio.Reader, name string) (*Class, error) {
	magic := make([]byte, 4)
	_, err := io.ReadFull(reader, magic)
	if err != nil {
		return nil, err
	}

	if !bytes.Equal(magic, MAGIC) {
		return nil, errors.New("magic does not match")
	}

	// skip minor and major version
	_, err = reader.Discard(4)
	if err != nil {
		return nil, err
	}

	constantPoolCount, err := readUint16(reader)
	if err != nil {
		return nil, err
	}

	constantPool, err := NewConstantPool(ctx, reader, int(constantPoolCount))
	if err != nil {
		return nil, fmt.Errorf("constant pool in %s: %v", name, err)
	}

	// skip to interfaces
	_, err = reader.Discard(6)
	if err != nil {
		return nil, err
	}

	interfacesCount, err := readUint16(reader)
	if err != nil {
		return nil, err
	}

	interfaces := make([]uint16, interfacesCount)
	for i := range interfacesCount {
		interfaceIndex, err := readUint16(reader)
		if err != nil {
			return nil, err
		}
		interfaces[i] = interfaceIndex
	}

	fieldsCount, err := readUint16(reader)
	if err != nil {
		return nil, err
	}

	fields, err := NewFields(reader, fieldsCount, constantPool)
	if err != nil {
		return nil, err
	}

	methodsCount, err := readUint16(reader)
	if err != nil {
		return nil, err
	}

	methods, err := NewMethods(reader, methodsCount, constantPool)
	if err != nil {
		return nil, err
	}

	attributesCount, err := readUint16(reader)
	if err != nil {
		return nil, err
	}

	attributes, err := NewAttributes(reader, attributesCount, constantPool)
	if err != nil {
		return nil, err
	}

	return &Class{
		Name:         name,
		ConstantPool: *constantPool,
		Methods:      methods,
		Fields:       fields,
		Interfaces:   interfaces,
		Attributes:   attributes,
	}, nil
}
