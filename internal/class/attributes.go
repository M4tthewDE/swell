package class

import (
	"bufio"
	"fmt"
	"io"
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
			return nil, fmt.Errorf("failed to parse attribute %d/%d: %v", i, count, err)
		}

		attributes[i] = attribute
	}

	return attributes, nil
}

func NewAttribute(reader *bufio.Reader, cp *ConstantPool) (Attribute, error) {
	nameIndex, err := readUint16(reader)
	if err != nil {
		return nil, fmt.Errorf("nameIndex failed: %v", err)
	}

	length, err := readUint32(reader)
	if err != nil {
		return nil, fmt.Errorf("length failed for nameIndex=%d: %v", nameIndex, err)
	}

	name, err := cp.GetUtf8(nameIndex)
	if err != nil {
		return nil, fmt.Errorf("name failed for nameIndex=%d: %v", nameIndex, err)
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
			return nil, fmt.Errorf("skipping unknown attribute %s failed: %v", name, err)
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
		return nil, fmt.Errorf("startPc failed: %v", err)
	}

	endPc, err := readUint16(reader)
	if err != nil {
		return nil, fmt.Errorf("endPc failed, startPc=%d: %v", startPc, err)
	}

	handlerPc, err := readUint16(reader)
	if err != nil {
		return nil, fmt.Errorf("handlerPc failed, startPc=%d,endPc=%d: %v", startPc, endPc, err)
	}

	catchType, err := readUint16(reader)
	if err != nil {
		return nil, fmt.Errorf("catchType failed, startPc=%d,endPc=%d,handlerPc=%d: %v", startPc, endPc, handlerPc, err)
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
		return nil, fmt.Errorf("attribute Code failed to read maxStack: %v", err)
	}

	maxLocals, err := readUint16(reader)
	if err != nil {
		return nil, fmt.Errorf("attribute Code failed to read maxLocals: %v", err)
	}

	codeLength, err := readUint32(reader)
	if err != nil {
		return nil, fmt.Errorf("attribute Code failed to read codeLength: %v", err)
	}

	code := make([]byte, codeLength)
	_, err = io.ReadFull(reader, code)
	if err != nil {
		return nil, fmt.Errorf("attribute Code failed to read code: %v", err)
	}

	exceptionTableLength, err := readUint16(reader)
	if err != nil {
		return nil, fmt.Errorf("attribute Code failed to read exception table length: %v", err)
	}

	exceptions := make([]Exception, exceptionTableLength)
	for i := range exceptionTableLength {
		exception, err := NewException(reader)
		if err != nil {
			return nil, fmt.Errorf("attribute Code failed to read exception %d/%d: %v", i, exceptionTableLength, err)
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
