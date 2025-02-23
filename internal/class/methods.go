package class

import (
	"bufio"
	"errors"
	"fmt"
)

const MainDescriptor = "([Ljava/lang/String;)V"

const AccPublic = 0x0001
const AccStatic = 0x0008
const AccVarargs = 0x0080
const AccNative = 0x0100

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

	return descriptor == MainDescriptor, nil
}

func (m Method) isPublic() bool {
	return (m.AccessFlags & AccPublic) != 0
}

func (m Method) isStatic() bool {
	return (m.AccessFlags & AccStatic) != 0
}

func (m Method) IsNative() bool {
	return (m.AccessFlags & AccNative) != 0
}

func (m Method) IsVarargs() bool {
	return (m.AccessFlags & AccVarargs) != 0
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
			return nil, fmt.Errorf("failed to parse method %d: %v", i, err)
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
		return nil, fmt.Errorf("failed to parse attributes, nameIndex=%d: %v", nameIndex, err)
	}

	return &Method{
			AccessFlags:     accessFlags,
			NameIndex:       nameIndex,
			DescriptorIndex: descriptorIndex,
			Attributes:      attributes},
		nil
}
