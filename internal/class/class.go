package class

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
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

func readUint32(reader *bufio.Reader) (uint32, error) {
	data := make([]byte, 4)
	_, err := reader.Read(data)
	if err != nil {
		return 0, err
	}

	return binary.BigEndian.Uint32(data), nil
}

var MAGIC = []byte{0xCA, 0xFE, 0xBA, 0xBE}

type Class struct {
	constantPool ConstantPool
	methods      []Method
	attributes   []Attribute
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

	// skip minor and major version
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

	// skip to methods_count (this relies on interfaces and fields being 0)
	_, err = reader.Discard(10)
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
		constantPool: *constantPool,
		methods:      methods,
		attributes:   attributes,
	}, nil
}
