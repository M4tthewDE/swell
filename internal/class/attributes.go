package class

import (
	"bufio"
)

type Attribute struct {
	nameIndex uint16
}

func NewAttributes(reader *bufio.Reader, count uint16) ([]Attribute, error) {
	attributes := make([]Attribute, count)
	for i := range count {
		attribute, err := NewAttribute(reader)
		if err != nil {
			return nil, err
		}

		attributes[i] = *attribute
	}

	return attributes, nil
}

func NewAttribute(reader *bufio.Reader) (*Attribute, error) {
	nameIndex, err := readUint16(reader)
	if err != nil {
		return nil, err
	}

	attributeLength, err := readUint32(reader)
	if err != nil {
		return nil, err
	}

	// TODO: parse attribute
	_, err = reader.Discard(int(attributeLength))
	if err != nil {
		return nil, err
	}

	return &Attribute{nameIndex: nameIndex}, nil
}
