package class

import (
	"bufio"
	"context"
	"fmt"
	"io"

	"github.com/m4tthewde/swell/internal/logger"
)

type ConstantPool struct {
	Infos []CpInfo `json:"infos"`
}

func NewConstantPool(ctx context.Context, reader *bufio.Reader, count int) (*ConstantPool, error) {
	log := logger.FromContext(ctx)

	infos := make([]CpInfo, count)
	infos[0] = ReservedInfo{}

	for i := 1; i < count; i++ {
		cpInfo, err := NewCpInfo(reader)
		if err != nil {
			return nil, err
		}

		if _, ok := cpInfo.(LongInfo); ok {
			i += 1
		}

		log.Debugf("cp info %d: %s", i, cpInfo)

		infos[i] = cpInfo
	}

	return &ConstantPool{Infos: infos}, nil
}

func (cp *ConstantPool) Get(index int) (*CpInfo, error) {
	if index >= len(cp.Infos) {
		return nil, fmt.Errorf("invalid constant pool index: %d", index)
	}

	return &cp.Infos[index], nil
}

func (cp *ConstantPool) GetUtf8(n uint16) (string, error) {
	if info, ok := cp.Infos[n].(Utf8Info); ok {
		return info.Content, nil
	}

	return "", nil
}

func (cp *ConstantPool) Ref(n uint16) (*RefInfo, error) {
	if info, ok := cp.Infos[n].(RefInfo); ok {
		return &info, nil
	}

	return nil, fmt.Errorf("no ref info found at %d", n)
}

func (cp *ConstantPool) Class(n uint16) (*ClassInfo, error) {
	if info, ok := cp.Infos[n].(ClassInfo); ok {
		return &info, nil
	}

	return nil, fmt.Errorf("no class info found at %d", n)
}

func (cp *ConstantPool) NameAndType(n uint16) (*NameAndTypeInfo, error) {
	if info, ok := cp.Infos[n].(NameAndTypeInfo); ok {
		return &info, nil
	}

	return nil, fmt.Errorf("no NameAndType info found at %d", n)
}

const Utf8Tag = 1
const IntegerTag = 3
const LongTag = 5
const ClassTag = 7
const StringTag = 8
const FieldrefTag = 9
const MethodrefTag = 10
const InterfaceMethodrefTag = 11
const NameAndTypeTag = 12
const MethodHandleTag = 15
const MethodTypeTag = 16
const InvokeDynamicTag = 18

type CpInfo interface {
	String() string
}

type ReservedInfo struct{}

func (c ReservedInfo) String() string {
	return "ReservedInfo"
}

func NewCpInfo(reader *bufio.Reader) (CpInfo, error) {
	tag, err := readUint8(reader)
	if err != nil {
		return nil, err
	}

	switch tag {
	case Utf8Tag:
		return NewUtf8Info(reader)
	case IntegerTag:
		return NewIntegerInfo(reader)
	case LongTag:
		return NewLongInfo(reader)
	case ClassTag:
		return NewClassInfo(reader)
	case StringTag:
		return NewStringInfo(reader)
	case FieldrefTag:
		return NewRefInfo(reader)
	case MethodrefTag:
		return NewRefInfo(reader)
	case InterfaceMethodrefTag:
		return NewRefInfo(reader)
	case NameAndTypeTag:
		return NewNameAndTypeInfo(reader)
	case MethodHandleTag:
		return NewMethodHandleInfo(reader)
	case MethodTypeTag:
		return NewMethodTypeInfo(reader)
	case InvokeDynamicTag:
		return NewInvokeDynamicInfo(reader)
	default:
		return nil, fmt.Errorf("unknown tag: %d", tag)
	}
}

type RefInfo struct {
	ClassIndex       uint16 `json:"class_index"`
	NameAndTypeIndex uint16 `json:"name_and_type_index"`
}

func (c RefInfo) String() string {
	return fmt.Sprintf("RefInfo[ClassIndex: %d, NameAndTypeIndex: %d]", c.ClassIndex, c.NameAndTypeIndex)
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

func (c ClassInfo) String() string {
	return fmt.Sprintf("ClassInfo[NameIndex: %d]", c.NameIndex)
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

func (c NameAndTypeInfo) String() string {
	return fmt.Sprintf("NameAndTypeInfo[NameIndex: %d, DescriptorIndex: %d]", c.NameIndex, c.DescriptorIndex)
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

func (c Utf8Info) String() string {
	return fmt.Sprintf("Utf8Info['%s']", c.Content)
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

func (c StringInfo) String() string {
	return fmt.Sprintf("StringInfo[StringIndex: %d]", c.StringIndex)
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

func (c InvokeDynamicInfo) String() string {
	return fmt.Sprintf("InvokeDynamicInfo[BootstrapMethodAttributeIndex: %d, NameAndTypeIndex: %d]", c.BootstrapMethodAttributeIndex, c.NameAndTypeIndex)
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

type LongInfo struct {
	Value uint64 `json:"value"`
}

func (c LongInfo) String() string {
	return fmt.Sprintf("LongInfo[%d]", c.Value)
}

func NewLongInfo(reader *bufio.Reader) (CpInfo, error) {
	value, err := readUint64(reader)
	if err != nil {
		return nil, err
	}

	return LongInfo{Value: value}, nil
}

type IntegerInfo struct {
	Value uint32 `json:"value"`
}

func (c IntegerInfo) String() string {
	return fmt.Sprintf("IntegerInfo[%d]", c.Value)
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

func (c MethodHandleInfo) String() string {
	return fmt.Sprintf("MethodHandleInfo[ReferenceKind: %d, ReferenceIndex: %d]", c.ReferenceKind, c.ReferenceIndex)
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

func (c MethodTypeInfo) String() string {
	return fmt.Sprintf("MethodTypeInfo[DescriptorIndex: %d]", c.DescriptorIndex)
}

func NewMethodTypeInfo(reader *bufio.Reader) (CpInfo, error) {
	descriptorIndex, err := readUint16(reader)
	if err != nil {
		return nil, err
	}

	return MethodTypeInfo{DescriptorIndex: descriptorIndex}, nil
}
