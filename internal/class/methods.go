package class

import (
	"bufio"
	"errors"
)

const MAIN_DESCRIPTOR = "([Ljava/lang/String;)V"

const ACC_PUBLIC = 0x0001
const ACC_STATIC = 0x0008

type Method struct {
	AccessFlags     uint16      `json:"access_flags"`
	NameIndex       uint16      `json:"name_index"`
	DescriptorIndex uint16      `json:"descriptor_index"`
	Attributes      []Attribute `json:"attributes"`
}

func (m Method) IsMain(cp *ConstantPool) (bool, error) {
	name, err := cp.GetUtf8(m.NameIndex)
	if err != nil {
		return false, err
	}

	if name != "main" {
		return false, nil
	}

	if !m.isPublic() || !m.isStatic() {
		return false, nil
	}

	descriptor, err := cp.GetUtf8(m.DescriptorIndex)
	if err != nil {
		return false, err
	}

	return descriptor == MAIN_DESCRIPTOR, nil
}

func (m Method) isPublic() bool {
	return (m.AccessFlags & ACC_PUBLIC) != 0
}

func (m Method) isStatic() bool {
	return (m.AccessFlags & ACC_STATIC) != 0
}

func (m Method) CodeAttribute() (*CodeAttribute, error) {
	for _, attribute := range m.Attributes {
		if codeAttribute, ok := attribute.(CodeAttribute); ok {
			return &codeAttribute, nil
		}
	}

	return nil, errors.New("no code attribute found")
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
			AccessFlags:     accessFlags,
			NameIndex:       nameIndex,
			DescriptorIndex: descriptorIndex,
			Attributes:      attributes},
		nil
}
