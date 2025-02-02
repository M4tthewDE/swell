package class

import (
	"bufio"
)

type Method struct {
	accessFlags     uint16
	nameIndex       uint16
	descriptorIndex uint16
	attributes      []Attribute
}

func NewMethods(reader *bufio.Reader, count uint16, cp *ConstantPool) ([]Method, error) {
	methods := make([]Method, count)
	for i := range count {
		method, err := NewMethod(reader, cp)
		if err != nil {
			return nil, err
		}

		methods[i] = *method
	}

	return methods, nil
}

func NewMethod(reader *bufio.Reader, cp *ConstantPool) (*Method, error) {
	accessFlags, err := readUint16(reader)
	if err != nil {
		return nil, err
	}

	nameIndex, err := readUint16(reader)
	if err != nil {
		return nil, err
	}

	descriptorIndex, err := readUint16(reader)
	if err != nil {
		return nil, err
	}

	attributesCount, err := readUint16(reader)
	if err != nil {
		return nil, err
	}

	attributes, err := NewAttributes(reader, attributesCount, cp)
	if err != nil {
		return nil, err
	}

	return &Method{
			accessFlags:     accessFlags,
			nameIndex:       nameIndex,
			descriptorIndex: descriptorIndex,
			attributes:      attributes},
		nil
}
