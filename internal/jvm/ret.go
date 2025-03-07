package jvm

import (
	"fmt"

	"github.com/m4tthewde/swell/internal/class"
)

func ret(r *Runner) error {
	pool, err := r.stack.CurrentConstantPool()
	if err != nil {
		return err
	}

	method, err := r.stack.CurrentMethod()
	if err != nil {
		return err
	}

	descriptor, err := pool.GetUtf8(method.DescriptorIndex)
	if err != nil {
		return err
	}

	methodDescriptor, err := class.NewMethodDescriptor(descriptor)
	if err != nil {
		return err
	}

	if methodDescriptor.ReturnDescriptor != 'V' {
		return fmt.Errorf("method has to be void, is %s", methodDescriptor.ReturnDescriptor)
	}

	return nil
}
