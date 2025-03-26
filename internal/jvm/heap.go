package jvm

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/m4tthewde/swell/internal/class"
	"github.com/m4tthewde/swell/internal/jvm/stack"
	"github.com/m4tthewde/swell/internal/logger"
)

type HeapItem interface {
	IsHeapItem()
}

type Object struct {
	className string
	fields    map[string]stack.Value
}

func (o Object) IsHeapItem() {}

func newObject(name string, fields map[string]stack.Value) Object {
	return Object{className: name, fields: fields}
}

func (o *Object) GetFieldValue(name string) (stack.Value, error) {
	value, ok := o.fields[name]
	if !ok {
		return nil, fmt.Errorf("field with name %s not found on %s", name, o.className)
	}

	return value, nil
}

type Array struct {
	items []stack.Value
}

func (a Array) IsHeapItem() {}
func (a Array) IsValue()    {}
func (a Array) String() string {
	return "Array[...]"
}

type Heap struct {
	items map[uuid.UUID]HeapItem
}

func NewHeap() Heap {
	return Heap{items: make(map[uuid.UUID]HeapItem, 0)}
}

func (h *Heap) AllocateObject(ctx context.Context, c *class.Class) (*uuid.UUID, error) {
	log := logger.FromContext(ctx)

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
	h.items[id] = newObject(c.Name, fields)

	log.Infow("allocated object", "id", id, "className", c.Name)

	return &id, nil
}

func (h *Heap) GetObject(id uuid.UUID) (*Object, error) {
	item, ok := h.items[id]
	if !ok {
		return nil, fmt.Errorf("object with id %s not found", id)
	}

	if object, ok := item.(Object); ok {
		return &object, nil
	}

	return nil, fmt.Errorf("object with id %s not found", id)
}

func (h *Heap) AllocateDefaultArray(ctx context.Context, size int, defaultValue stack.Value) (*uuid.UUID, error) {
	items := make([]stack.Value, 0)
	for range size {
		items = append(items, defaultValue)
	}

	id := uuid.New()
	h.items[id] = Array{items: items}

	return &id, nil
}

func (h *Heap) SetField(id uuid.UUID, fieldName string, value stack.Value) error {
	obj, err := h.GetObject(id)
	if err != nil {
		return err
	}

	obj.fields[fieldName] = value
	h.items[id] = *obj
	return nil
}
