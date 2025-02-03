package class

import (
	"bufio"
	"errors"
)

type Attribute interface {
	Name() string
}

func NewAttributes(reader *bufio.Reader, count uint16, cp *ConstantPool) ([]Attribute, error) {
	attributes := make([]Attribute, count)
	for i := range count {
		attribute, err := NewAttribute(reader, cp)
		if err != nil {
			return nil, err
		}

		attributes[i] = attribute
	}

	return attributes, nil
}

func NewAttribute(reader *bufio.Reader, cp *ConstantPool) (Attribute, error) {
	nameIndex, err := readUint16(reader)
	if err != nil {
		return nil, err
	}

	_, err = readUint32(reader)
	if err != nil {
		return nil, err
	}

	name, err := cp.GetUtf8(nameIndex)
	if err != nil {
		return nil, err
	}

	switch name {
	case "Code":
		return NewCodeAttribute(reader, cp)
	case "LineNumberTable":
		return NewLineNumberTableAttribute(reader)
	case "SourceFile":
		return NewSourceFileAttribute(reader)
	case "ConstantValue":
		return NewConstantValueAttribute(reader)
	case "RuntimeVisibleAnnotations":
		return NewRuntimeVisibleAnnotationsAttribute(reader)
	default:
		return nil, errors.New("unknown attribute: " + name)
	}
}

type CodeAttribute struct {
	MaxStack   uint16      `json:"max_stack"`
	MaxLocals  uint16      `json:"max_locals"`
	Code       []byte      `json:"code"`
	Attributes []Attribute `json:"attributes"`
}

func (c CodeAttribute) Name() string {
	return "Code"
}

func NewCodeAttribute(reader *bufio.Reader, cp *ConstantPool) (Attribute, error) {
	maxStack, err := readUint16(reader)
	if err != nil {
		return nil, err
	}

	maxLocals, err := readUint16(reader)
	if err != nil {
		return nil, err
	}

	codeLength, err := readUint32(reader)
	if err != nil {
		return nil, err
	}

	code := make([]byte, codeLength)
	_, err = reader.Read(code)

	exceptionTableLength, err := readUint16(reader)
	if err != nil {
		return nil, err
	}

	if exceptionTableLength != 0 {
		return nil, errors.New("exceptions not supported")
	}

	attributesCount, err := readUint16(reader)
	if err != nil {
		return nil, err
	}

	attributes, err := NewAttributes(reader, attributesCount, cp)
	if err != nil {
		return nil, err
	}

	return CodeAttribute{
		MaxStack:   maxStack,
		MaxLocals:  maxLocals,
		Code:       code,
		Attributes: attributes,
	}, nil
}

type LineNumberTableAttribute struct {
	Table []LineNumberTableEntry `json:"table"`
}

func (l LineNumberTableAttribute) Name() string {
	return "LineNumberTable"
}

type LineNumberTableEntry struct {
	StartPc    uint16 `json:"start_pc"`
	LineNumber uint16 `json:"line_number"`
}

func NewLineNumberTableAttribute(reader *bufio.Reader) (Attribute, error) {
	length, err := readUint16(reader)
	if err != nil {
		return nil, err
	}

	table := make([]LineNumberTableEntry, length)

	for i := range length {
		startPc, err := readUint16(reader)
		if err != nil {
			return nil, err
		}

		lineNumber, err := readUint16(reader)
		if err != nil {
			return nil, err
		}

		table[i] = LineNumberTableEntry{
			StartPc:    startPc,
			LineNumber: lineNumber,
		}
	}

	return LineNumberTableAttribute{Table: table}, nil
}

type SourceFileAttribute struct {
	SourceFileIndex uint16 `json:"source_file_index"`
}

func (s SourceFileAttribute) Name() string {
	return "SourceFile"
}

func NewSourceFileAttribute(reader *bufio.Reader) (Attribute, error) {
	sourceFileIndex, err := readUint16(reader)
	if err != nil {
		return nil, err
	}

	return SourceFileAttribute{SourceFileIndex: sourceFileIndex}, nil
}

type ConstantValueAttribute struct {
	ConstantValueIndex uint16 `json:"constant_value_index"`
}

func (c ConstantValueAttribute) Name() string {
	return "ConstantValue"
}

func NewConstantValueAttribute(reader *bufio.Reader) (Attribute, error) {
	constantValueIndex, err := readUint16(reader)
	if err != nil {
		return nil, err
	}

	return ConstantValueAttribute{ConstantValueIndex: constantValueIndex}, nil
}

type Annotation struct {
	TypeIndex uint16 `json:"type_index"`
}

func NewAnnotation(reader *bufio.Reader) (*Annotation, error) {
	typeIndex, err := readUint16(reader)
	if err != nil {
		return nil, err
	}

	numElementValuePairs, err := readUint16(reader)
	if err != nil {
		return nil, err
	}

	if numElementValuePairs != 0 {
		return nil, errors.New("TODO: parse element value pairs")
	}

	return &Annotation{TypeIndex: typeIndex}, nil
}

type RuntimeVisibleAnnotationsAttribute struct {
	Annotations []Annotation `json:"annotations"`
}

func (r RuntimeVisibleAnnotationsAttribute) Name() string {
	return "RuntimeVisibleAnnotations"
}

func NewRuntimeVisibleAnnotationsAttribute(reader *bufio.Reader) (Attribute, error) {
	numAnnotations, err := readUint16(reader)
	if err != nil {
		return nil, err
	}

	annotations := make([]Annotation, numAnnotations)
	for i := range numAnnotations {
		annotation, err := NewAnnotation(reader)
		if err != nil {
			return nil, err
		}
		annotations[i] = *annotation
	}

	return RuntimeVisibleAnnotationsAttribute{Annotations: annotations}, nil
}
