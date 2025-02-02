package class

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"log"
	"os"
)

func readUint8(reader *bufio.Reader) (uint8, error) {
	byte, err := reader.ReadByte()
	if err != nil {
		return 0, err
	}

	return byte, nil
}

func readUint16(reader *bufio.Reader) (uint16, error) {
	data := make([]byte, 2)
	_, err := reader.Read(data)
	if err != nil {
		return 0, err
	}

	return binary.BigEndian.Uint16(data), nil
}

var MAGIC = []byte{0xCA, 0xFE, 0xBA, 0xBE}

type Class struct {
	constantPool ConstantPool
}

func NewClass(path string) (*Class, error) {
	log.Printf("parsing %s", path)

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	reader := bufio.NewReader(file)

	magic := make([]byte, 4)
	_, err = reader.Read(magic)
	if err != nil {
		return nil, err
	}

	if !bytes.Equal(magic, MAGIC) {
		return nil, errors.New("magic does not match")
	}

	// we have no use for minor and major version as of now
	_, err = reader.Discard(4)
	if err != nil {
		return nil, err
	}

	constantPoolCount, err := readUint16(reader)
	if err != nil {
		return nil, err
	}

	constantPool, err := NewConstantPool(reader, constantPoolCount)
	if err != nil {
		return nil, err
	}

	return &Class{constantPool: *constantPool}, nil
}

type ConstantPool struct {
	infos []CpInfo
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

type ConstantKind interface{}

type CpInfo struct {
	tag  uint8
	info ConstantKind
}

func NewCpInfo(reader *bufio.Reader) (*CpInfo, error) {
	tag, err := readUint8(reader)
	if err != nil {
		return nil, err
	}

	var info ConstantKind
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

	return &CpInfo{tag: tag, info: info}, nil
}

type RefInfo struct {
	classIndex       uint16
	nameAndTypeIndex uint16
}

func NewRefInfo(reader *bufio.Reader) (ConstantKind, error) {
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

func NewClassInfo(reader *bufio.Reader) (ConstantKind, error) {
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

func NewNameAndTypeInfo(reader *bufio.Reader) (ConstantKind, error) {
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

func NewUtf8Info(reader *bufio.Reader) (ConstantKind, error) {
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

func NewStringInfo(reader *bufio.Reader) (ConstantKind, error) {
	stringIndex, err := readUint16(reader)
	if err != nil {
		return nil, err
	}

	return StringInfo{stringIndex: stringIndex}, nil
}
