package class

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"log"
	"os"
)

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
}

func NewConstantPool(reader *bufio.Reader, count uint16) (*ConstantPool, error) {
	return nil, nil
}
