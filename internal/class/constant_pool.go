package class

import (
	"bufio"
	"errors"
	"fmt"
	"io"
)

type ConstantPool struct {
	Infos []CpInfo `json:"infos"`
}

func NewConstantPool(reader *bufio.Reader, count uint16) (*ConstantPool, error) {
	infos := make([]CpInfo, count-1)

	for i := range count - 1 {
		cpInfo, err := NewCpInfo(reader)
		if err != nil {
			return nil, err
		}

		infos[i] = cpInfo
	}

	return &ConstantPool{Infos: infos}, nil
}

func (cp *ConstantPool) GetUtf8(n uint16) (string, error) {
	if info, ok := cp.Infos[n-1].(Utf8Info); ok {
		return info.Content, nil
	}

	return "", nil
}

func (cp *ConstantPool) Ref(n uint16) (*RefInfo, error) {
	if info, ok := cp.Infos[n-1].(RefInfo); ok {
		return &info, nil
	}

	return nil, errors.New(fmt.Sprintf("no ref info found at %d", n))
}

func (cp *ConstantPool) Class(n uint16) (*ClassInfo, error) {
	if info, ok := cp.Infos[n-1].(ClassInfo); ok {
		return &info, nil
	}

	return nil, errors.New(fmt.Sprintf("no class info found at %d", n))
}

const UTF8_TAG = 1
const INTEGER_TAG = 3
const CLASS_TAG = 7
const STRING_TAG = 8
const FIELDREF_TAG = 9
const METHODREF_TAG = 10
const INTERFACE_METHODREF_TAG = 11
const NAME_AND_TYPE_TAG = 12
const METHOD_HANDLE_TAG = 15
const METHOD_TYPE_TAG = 16
const INVOKE_DYNAMIC_TAG = 18

type CpInfo interface{}

func NewCpInfo(reader *bufio.Reader) (CpInfo, error) {
	tag, err := readUint8(reader)
	if err != nil {
		return nil, err
	}

	switch tag {
	case INTEGER_TAG:
		return NewIntegerInfo(reader)
	case UTF8_TAG:
		return NewUtf8Info(reader)
	case CLASS_TAG:
		return NewClassInfo(reader)
	case STRING_TAG:
		return NewStringInfo(reader)
	case FIELDREF_TAG:
		return NewRefInfo(reader)
	case METHODREF_TAG:
		return NewRefInfo(reader)
	case INTERFACE_METHODREF_TAG:
		return NewRefInfo(reader)
	case NAME_AND_TYPE_TAG:
		return NewNameAndTypeInfo(reader)
	case METHOD_HANDLE_TAG:
		return NewMethodHandleInfo(reader)
	case METHOD_TYPE_TAG:
		return NewMethodTypeInfo(reader)
	case INVOKE_DYNAMIC_TAG:
		return NewInvokeDynamicInfo(reader)
	default:
		return nil, errors.New(fmt.Sprintf("unknown tag: %d", tag))
	}
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
	_, err = io.ReadFull(reader, bytes)
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

type InvokeDynamicInfo struct {
	BootstrapMethodAttributeIndex uint16 `json:"bootstrap_method_attribute_index"`
	NameAndTypeIndex              uint16 `json:"name_and_type_index"`
}

func NewInvokeDynamicInfo(reader *bufio.Reader) (CpInfo, error) {
	bootstrapMethodAttributeIndex, err := readUint16(reader)
	if err != nil {
		return nil, err
	}

	nameAndTypeIndex, err := readUint16(reader)
	if err != nil {
		return nil, err
	}

	return InvokeDynamicInfo{
		BootstrapMethodAttributeIndex: bootstrapMethodAttributeIndex,
		NameAndTypeIndex:              nameAndTypeIndex,
	}, nil
}

type IntegerInfo struct {
	Value uint32 `json:"value"`
}

func NewIntegerInfo(reader *bufio.Reader) (CpInfo, error) {
	value, err := readUint32(reader)
	if err != nil {
		return nil, err
	}

	return IntegerInfo{Value: value}, nil
}

type MethodHandleInfo struct {
	ReferenceKind  uint8  `json:"reference_kind"`
	ReferenceIndex uint16 `json:"reference_index"`
}

func NewMethodHandleInfo(reader *bufio.Reader) (CpInfo, error) {
	referenceKind, err := readUint8(reader)
	if err != nil {
		return nil, err
	}

	referenceIndex, err := readUint16(reader)
	if err != nil {
		return nil, err
	}

	return MethodHandleInfo{
		ReferenceKind:  referenceKind,
		ReferenceIndex: referenceIndex,
	}, nil
}

type MethodTypeInfo struct {
	DescriptorIndex uint16 `json:"descriptor_index"`
}

func NewMethodTypeInfo(reader *bufio.Reader) (CpInfo, error) {
	descriptorIndex, err := readUint16(reader)
	if err != nil {
		return nil, err
	}

	return MethodTypeInfo{DescriptorIndex: descriptorIndex}, nil
}
