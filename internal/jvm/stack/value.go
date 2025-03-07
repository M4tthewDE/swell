package stack

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/m4tthewde/swell/internal/class"
)

type Value interface {
	isValue()
	String() string
}

func DefaultValue(typ class.FieldType) (Value, error) {
	if _, ok := typ.(class.ObjectType); ok {
		return ReferenceValue{Value: nil}, nil
	}

	switch typ {
	case class.BaseType('I'):
		return IntValue{Value: 0}, nil
	case class.BaseType('J'):
		return LongValue{Value: 0}, nil
	default:
		// TODO: add DefaultValue function to FieldType interface so that switch becomes obsolete
		return nil, fmt.Errorf("unknown field type %s", typ)
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
	// reference to the Class object
	Value *uuid.UUID
	Class *class.Class
}

func (v ClassReferenceValue) String() string {
	return fmt.Sprintf("ClassReference=%s", v.Class.Name)
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
