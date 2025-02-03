package class

import "bufio"

type Field struct {
	AccessFlags     uint16      `json:"access_flags"`
	NameIndex       uint16      `json:"name_index"`
	DescriptorIndex uint16      `json:"descriptor_index"`
	Attributes      []Attribute `json:"attributes"`
}

func NewFields(reader *bufio.Reader, count uint16, cp *ConstantPool) ([]Field, error) {
	fields := make([]Field, count)
	for i := range count {
		field, err := NewField(reader, cp)
		if err != nil {
			return nil, err
		}

		fields[i] = *field
	}

	return fields, nil
}

func NewField(reader *bufio.Reader, cp *ConstantPool) (*Field, error) {
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

	return &Field{
			AccessFlags:     accessFlags,
			NameIndex:       nameIndex,
			DescriptorIndex: descriptorIndex,
			Attributes:      attributes},
		nil
}
