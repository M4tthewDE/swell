package class

import (
	"bufio"
)

type Attribute interface {
	Name() string
}

type UnknownAttribute struct{}

func (u UnknownAttribute) Name() string {
	return "Unknown"
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

	length, err := readUint32(reader)
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
	default:
		// skip unknown attributes
		_, err = reader.Discard(int(length))
		if err != nil {
			return nil, err
		}

		return UnknownAttribute{}, nil
	}
}

type Exception struct {
	StartPc   uint16
	EndPc     uint16
	HandlerPc uint16
	CatchType uint16
}

func NewException(reader *bufio.Reader) (*Exception, error) {
	startPc, err := readUint16(reader)
	if err != nil {
		return nil, err
	}

	endPc, err := readUint16(reader)
	if err != nil {
		return nil, err
	}

	handlerPc, err := readUint16(reader)
	if err != nil {
		return nil, err
	}

	catchType, err := readUint16(reader)
	if err != nil {
		return nil, err
	}

	return &Exception{
		StartPc:   startPc,
		EndPc:     endPc,
		HandlerPc: handlerPc,
		CatchType: catchType,
	}, nil
}

type CodeAttribute struct {
	MaxStack   uint16      `json:"max_stack"`
	MaxLocals  uint16      `json:"max_locals"`
	Code       []byte      `json:"code"`
	Exceptions []Exception `json:"exceptions"`
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

	exceptions := make([]Exception, exceptionTableLength)
	for i := range exceptionTableLength {
		exception, err := NewException(reader)
		if err != nil {
			return nil, err
		}

		exceptions[i] = *exception
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
		Exceptions: exceptions,
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
