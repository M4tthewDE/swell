package class

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"os"
)

var MAGIC = []byte{0xCA, 0xFE, 0xBA, 0xBE}

type Class struct {
	ConstantPool ConstantPool `json:"constant_pool"`
	Methods      []Method     `json:"methods"`
	Attributes   []Attribute  `json:"attributes"`
}

func (c *Class) PrettyPrint() error {
	data, err := json.MarshalIndent(c, "", "\t")
	if err != nil {
		return err
	}

	log.Println(string(data))
	return nil
}

func (c *Class) GetMainMethod() (*Method, error) {
	for _, m := range c.Methods {
		isMain, err := m.IsMain(&c.ConstantPool)
		if err != nil {
			return nil, err
		}

		if isMain {
			return &m, nil
		}
	}

	return nil, errors.New("no main method found")
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
		ConstantPool: *constantPool,
		Methods:      methods,
		Attributes:   attributes,
	}, nil
}
