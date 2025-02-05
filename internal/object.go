package internal

import (
	"github.com/google/uuid"
	"github.com/m4tthewde/swell/internal/class"
)

type Object struct {
	className string
	fields    map[string]Value
}

func NewObject(name string, fields map[string]Value) Object {
	return Object{className: name, fields: fields}
}

type Heap struct {
	objects map[uuid.UUID]Object
}

func NewHeap() Heap {
	return Heap{objects: make(map[uuid.UUID]Object, 0)}
}

func (h *Heap) Allocate(c *class.Class) (*uuid.UUID, error) {
	fields := make(map[string]Value)

	for _, field := range c.Fields {
		name, err := c.ConstantPool.GetUtf8(field.NameIndex)
		if err != nil {
			return nil, err
		}

		descriptor, err := c.ConstantPool.GetUtf8(field.DescriptorIndex)
		if err != nil {
			return nil, err
		}

		fieldType, err := class.NewFieldType(descriptor)
		if err != nil {
			return nil, err
		}

		value, err := DefaultValue(fieldType)
		if err != nil {
			return nil, err
		}

		fields[name] = value
	}

	id := uuid.New()
	h.objects[id] = NewObject(c.Name, fields)

	return &id, nil
}
