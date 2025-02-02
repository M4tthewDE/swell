package class

import (
	"bufio"
	"errors"
	"fmt"
)

type ConstantPool struct {
	infos []CpInfo
}

func (cp *ConstantPool) GetUtf8(n uint16) (string, error) {
	if info, ok := cp.infos[n-1].(Utf8Info); ok {
		return string(info.bytes), nil
	}

	return "", nil
}

func NewConstantPool(reader *bufio.Reader, count uint16) (*ConstantPool, error) {
	infos := make([]CpInfo, count)

	for i := range count - 1 {
		cpInfo, err := NewCpInfo(reader)
		if err != nil {
			return nil, err
		}

		infos[i] = *cpInfo
	}

	return &ConstantPool{infos: infos}, nil
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
	classIndex       uint16
	nameAndTypeIndex uint16
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

	return RefInfo{classIndex: classIndex, nameAndTypeIndex: nameAndTypeIndex}, nil
}

type ClassInfo struct {
	nameIndex uint16
}

func NewClassInfo(reader *bufio.Reader) (CpInfo, error) {
	nameIndex, err := readUint16(reader)
	if err != nil {
		return nil, err
	}

	return ClassInfo{nameIndex: nameIndex}, nil
}

type NameAndTypeInfo struct {
	nameIndex       uint16
	descriptorIndex uint16
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

	return NameAndTypeInfo{nameIndex: nameIndex, descriptorIndex: descriptorIndex}, nil
}

type Utf8Info struct {
	bytes []byte
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

	return Utf8Info{bytes: bytes}, nil
}

type StringInfo struct {
	stringIndex uint16
}

func NewStringInfo(reader *bufio.Reader) (CpInfo, error) {
	stringIndex, err := readUint16(reader)
	if err != nil {
		return nil, err
	}

	return StringInfo{stringIndex: stringIndex}, nil
}
