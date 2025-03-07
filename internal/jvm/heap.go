package jvm

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/m4tthewde/swell/internal/class"
	"github.com/m4tthewde/swell/internal/jvm/stack"
	"github.com/m4tthewde/swell/internal/logger"
)

type Object struct {
	className string
	fields    map[string]stack.Value
}

func NewObject(name string, fields map[string]stack.Value) Object {
	return Object{className: name, fields: fields}
}

func (o *Object) GetFieldValue(name string) (stack.Value, error) {
	value, ok := o.fields[name]
	if !ok {
		return nil, fmt.Errorf("field with name %s not found on %s", name, o.className)
	}

	return value, nil
}

type Heap struct {
	objects map[uuid.UUID]Object
}

func NewHeap() Heap {
	return Heap{objects: make(map[uuid.UUID]Object, 0)}
}

func (h *Heap) Allocate(ctx context.Context, c *class.Class) (*uuid.UUID, error) {
	log := logger.FromContext(ctx)
	log.Infof("allocating %s object", c.Name)

	fields := make(map[string]stack.Value)

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

		value, err := stack.DefaultValue(fieldType)
		if err != nil {
			return nil, err
		}

		fields[name] = value
	}

	id := uuid.New()
	h.objects[id] = NewObject(c.Name, fields)

	return &id, nil
}

func (h *Heap) GetObject(id uuid.UUID) (*Object, error) {
	object, ok := h.objects[id]
	if !ok {
		return nil, fmt.Errorf("object with id %s not found", id)
	}

	return &object, nil
}
