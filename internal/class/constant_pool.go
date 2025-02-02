package class

import (
	"bufio"
	"errors"
	"fmt"
)

type ConstantPool struct {
	Infos []CpInfo `json:"infos"`
}

func (cp *ConstantPool) GetUtf8(n uint16) (string, error) {
	if info, ok := cp.Infos[n-1].(Utf8Info); ok {
		return info.Content, nil
	}

	return "", nil
}

func NewConstantPool(reader *bufio.Reader, count uint16) (*ConstantPool, error) {
	infos := make([]CpInfo, count-1)

	for i := range count - 1 {
		cpInfo, err := NewCpInfo(reader)
		if err != nil {
			return nil, err
		}

		infos[i] = *cpInfo
	}

	return &ConstantPool{Infos: infos}, nil
}

type CpInfo interface{}

func NewCpInfo(reader *bufio.Reader) (*CpInfo, error) {
	tag, err := readUint8(reader)
	if err != nil {
		return nil, err
	}

	var info CpInfo
	switch tag {
	case 1:
		info, err = NewUtf8Info(reader)
	case 7:
		info, err = NewClassInfo(reader)
	case 8:
		info, err = NewStringInfo(reader)
	case 9:
		info, err = NewRefInfo(reader)
	case 10:
		info, err = NewRefInfo(reader)
	case 12:
		info, err = NewNameAndTypeInfo(reader)
	default:
		return nil, errors.New(fmt.Sprintf("unknown tag: %d", tag))
	}

	if err != nil {
		return nil, err
	}

	return &info, nil
}

type RefInfo struct {
	ClassIndex       uint16 `json:"class_index"`
	NameAndTypeIndex uint16 `json:"name_and_type_index"`
}

func NewRefInfo(reader *bufio.Reader) (CpInfo, error) {
	classIndex, err := readUint16(reader)
	if err != nil {
		return nil, err
	}

	nameAndTypeIndex, err := readUint16(reader)
	if err != nil {
		return nil, err
	}

	return RefInfo{ClassIndex: classIndex, NameAndTypeIndex: nameAndTypeIndex}, nil
}

type ClassInfo struct {
	NameIndex uint16 `json:"name_index"`
}

func NewClassInfo(reader *bufio.Reader) (CpInfo, error) {
	nameIndex, err := readUint16(reader)
	if err != nil {
		return nil, err
	}

	return ClassInfo{NameIndex: nameIndex}, nil
}

type NameAndTypeInfo struct {
	NameIndex       uint16 `json:"name_index"`
	DescriptorIndex uint16 `json:"descriptor_index"`
}

func NewNameAndTypeInfo(reader *bufio.Reader) (CpInfo, error) {
	nameIndex, err := readUint16(reader)
	if err != nil {
		return nil, err
	}

	descriptorIndex, err := readUint16(reader)
	if err != nil {
		return nil, err
	}

	return NameAndTypeInfo{NameIndex: nameIndex, DescriptorIndex: descriptorIndex}, nil
}

type Utf8Info struct {
	Content string `json:"content"`
}

func NewUtf8Info(reader *bufio.Reader) (CpInfo, error) {
	length, err := readUint16(reader)
	if err != nil {
		return nil, err
	}

	bytes := make([]byte, length)
	_, err = reader.Read(bytes)
	if err != nil {
		return nil, err
	}

	return Utf8Info{Content: string(bytes)}, nil
}

type StringInfo struct {
	StringIndex uint16 `json:"string_index"`
}

func NewStringInfo(reader *bufio.Reader) (CpInfo, error) {
	stringIndex, err := readUint16(reader)
	if err != nil {
		return nil, err
	}

	return StringInfo{StringIndex: stringIndex}, nil
}
