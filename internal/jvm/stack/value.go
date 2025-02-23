package stack

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/m4tthewde/swell/internal/class"
)

type Value interface {
	isValue()
	String() string
}

func DefaultValue(typ class.FieldType) (Value, error) {
	switch typ.(type) {
	case BooleanValue:
		return BooleanValue{Value: false}, nil
	case ByteValue:
		return ByteValue{Value: 0}, nil
	case ShortValue:
		return ShortValue{Value: 0}, nil
	case IntValue:
		return IntValue{Value: 0}, nil
	case LongValue:
		return LongValue{Value: 0}, nil
	case CharValue:
		return CharValue{Value: 0}, nil
	case FloatValue:
		return FloatValue{Value: 0}, nil
	case DoubleValue:
		return DoubleValue{Value: 0}, nil
	case ReferenceValue:
		return ReferenceValue{Value: nil}, nil
	case ClassReferenceValue:
		return ClassReferenceValue{Value: nil}, nil
	default:
		return nil, errors.New("unknown field type")
	}
}

type BooleanValue struct {
	Value bool
}

func (v BooleanValue) String() string {
	return fmt.Sprintf("Boolean=%t", v.Value)
}

type ByteValue struct {
	Value uint8
}

func (v ByteValue) String() string {
	return fmt.Sprintf("Byte=%d", v.Value)
}

type ShortValue struct {
	Value uint16
}

func (v ShortValue) String() string {
	return fmt.Sprintf("Short=%d", v.Value)
}

type IntValue struct {
	Value uint32
}

func (v IntValue) String() string {
	return fmt.Sprintf("Int=%d", v.Value)
}

type LongValue struct {
	Value uint64
}

func (v LongValue) String() string {
	return fmt.Sprintf("Long=%d", v.Value)
}

type CharValue struct {
	Value rune
}

func (v CharValue) String() string {
	return fmt.Sprintf("Char=%c", v.Value)
}

type FloatValue struct {
	Value float32
}

func (v FloatValue) String() string {
	return fmt.Sprintf("Float=%f", v.Value)
}

type DoubleValue struct {
	Value float64
}

func (v DoubleValue) String() string {
	return fmt.Sprintf("Double=%f", v.Value)
}

type ReferenceValue struct {
	Value *uuid.UUID
}

func (v ReferenceValue) String() string {
	return fmt.Sprintf("Reference=%s", v.Value)
}

type ClassReferenceValue struct {
	Value *class.Class
}

func (v ClassReferenceValue) String() string {
	return fmt.Sprintf("ClassReference=%s", v.Value.Name)
}

func (v BooleanValue) isValue()        {}
func (v ByteValue) isValue()           {}
func (v ShortValue) isValue()          {}
func (v IntValue) isValue()            {}
func (v LongValue) isValue()           {}
func (v CharValue) isValue()           {}
func (v FloatValue) isValue()          {}
func (v DoubleValue) isValue()         {}
func (v ReferenceValue) isValue()      {}
func (v ClassReferenceValue) isValue() {}
