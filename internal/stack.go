package internal

import (
	"errors"

	"github.com/google/uuid"
	"github.com/m4tthewde/swell/internal/class"
)

type Value interface {
	isValue()
}

func DefaultValue(typ class.FieldType) (Value, error) {
	switch typ.(type) {
	case BooleanValue:
		return BooleanValue{value: false}, nil
	case ByteValue:
		return ByteValue{value: 0}, nil
	case ShortValue:
		return ShortValue{value: 0}, nil
	case IntValue:
		return IntValue{value: 0}, nil
	case LongValue:
		return LongValue{value: 0}, nil
	case CharValue:
		return CharValue{value: 0}, nil
	case FloatValue:
		return FloatValue{value: 0}, nil
	case DoubleValue:
		return DoubleValue{value: 0}, nil
	case ReferenceValue:
		return ReferenceValue{value: nil}, nil
	default:
		return nil, errors.New("unknown field type")
	}
}

type BooleanValue struct {
	value bool
}

type ByteValue struct {
	value uint8
}

type ShortValue struct {
	value uint16
}

type IntValue struct {
	value uint32
}

type LongValue struct {
	value uint64
}

type CharValue struct {
	value rune
}

type FloatValue struct {
	value float32
}

type DoubleValue struct {
	value float64
}

type ReferenceValue struct {
	value *uuid.UUID
}

func (b BooleanValue) isValue()   {}
func (b ByteValue) isValue()      {}
func (b ShortValue) isValue()     {}
func (b IntValue) isValue()       {}
func (b LongValue) isValue()      {}
func (b CharValue) isValue()      {}
func (b FloatValue) isValue()     {}
func (b DoubleValue) isValue()    {}
func (b ReferenceValue) isValue() {}

type Frame struct {
	className      string
	methodName     string
	operands       []Value
	localVariables []Value
}

func NewFrame(className string, methodName string, localVariables []Value) Frame {
	return Frame{className: className, methodName: methodName, operands: make([]Value, 0), localVariables: localVariables}
}

type Stack struct {
	frames []Frame
}

func NewStack() Stack {
	return Stack{frames: make([]Frame, 0)}
}

func (s *Stack) Push(className string, methodName string, localVariables []Value) {
	frame := NewFrame(className, methodName, localVariables)
	s.frames = append(s.frames, frame)
}

func (s *Stack) PopOperands(count int) []Value {
	frame := s.frames[len(s.frames)-1]
	operands := frame.operands[len(frame.operands)-count:]

	frame.operands = frame.operands[:len(frame.operands)-count]
	s.frames[len(s.frames)-1] = frame
	return operands
}

func (s *Stack) PushOperand(operand Value) {
	frame := s.frames[len(s.frames)-1]
	frame.operands = append(frame.operands, operand)
	s.frames[len(s.frames)-1] = frame
}

func (s *Stack) GetOperand() Value {
	frame := s.frames[len(s.frames)-1]
	return frame.operands[len(frame.operands)-1]
}

func (s *Stack) GetLocalVariable(n int) Value {
	frame := s.frames[len(s.frames)-1]
	return frame.localVariables[n]
}
