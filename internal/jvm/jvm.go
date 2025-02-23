package jvm

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/m4tthewde/swell/internal/class"
	"github.com/m4tthewde/swell/internal/loader"
	"github.com/m4tthewde/swell/internal/logger"
)

type Runner struct {
	classBeingInitialized string
	initializedClasses    map[string]struct{}
	pc                    int
	loader                loader.Loader
	stack                 Stack
	heap                  Heap
}

func NewRunner() Runner {
	return Runner{
		classBeingInitialized: "",
		initializedClasses:    make(map[string]struct{}),
		pc:                    0,
		loader:                loader.NewLoader(),
		stack:                 NewStack(),
		heap:                  NewHeap(),
	}

}

func (r *Runner) RunMain(ctx context.Context, className string) error {
	err := r.initializeClass(ctx, className)
	if err != nil {
		return err
	}

	c, err := r.loader.Load(ctx, className)
	if err != nil {
		return err
	}

	main, ok, err := c.GetMainMethod()
	if !ok {
		return errors.New("no main method found")
	}

	if err != nil {
		return err
	}

	code, err := main.CodeAttribute()
	if err != nil {
		return err
	}

	err = r.runMethod(ctx, code.Code, *c, *main, make([]Value, 0))
	if err != nil {
		return err
	}

	return nil
}

const LDC = 0x12
const ALOAD_0 = 0x2a
const RET = 0xb1
const GET_STATIC = 0xb2
const INVOKE_SPECIAL = 0xb7
const INVOKE_STATIC = 0xb8
const NEW = 0xbb
const DUP = 0x59

func (r *Runner) run(ctx context.Context, code []byte) error {
	log := logger.FromContext(ctx)

	for {
		instruction := code[r.pc]

		var err error
		switch instruction {
		case LDC:
			log.Info("ldc")
			err = ldc(r, ctx, code)
		case ALOAD_0:
			log.Info("aload_0")
			err = aload(r, 0)
		case GET_STATIC:
			log.Info("getstatic")
			err = getStatic(r, ctx, code)
		case INVOKE_STATIC:
			log.Info("invokestatic")
			err = invokeStatic(r, ctx, code)
		case NEW:
			log.Info("new")
			err = new(r, ctx, code)
		case DUP:
			log.Info("dup")
			err = dup(r)
		case INVOKE_SPECIAL:
			log.Info("invokespecial")
			err = invokeSpecial(r, ctx, code)
		case RET:
			log.Info("ret")
			return ret(r)
		default:
			err = errors.New(
				fmt.Sprintf("unknown instruction %x", instruction),
			)
		}

		if err != nil {
			return err
		}

		if r.pc == len(code) {
			return nil
		}
	}
}

func (r *Runner) initializeClass(ctx context.Context, className string) error {
	log := logger.FromContext(ctx)

	_, exists := r.initializedClasses[className]
	if exists {
		log.Infof("already initialized %s", className)
		return nil
	}

	if r.classBeingInitialized == className {
		return nil
	} else {
		r.classBeingInitialized = className
	}

	log.Infof("initializing %s", className)

	c, err := r.loader.Load(ctx, className)
	if err != nil {
		return err
	}

	clinit, ok, err := c.GetMethod("<clinit>")
	if !ok {
		r.initializedClasses[className] = struct{}{}
		log.Infof("initialized %s", className)
		return nil
	}

	if err != nil {
		return err
	}

	code, err := clinit.CodeAttribute()
	if err != nil {
		return err
	}

	err = r.runMethod(ctx, code.Code, *c, *clinit, make([]Value, 0))
	if err != nil {
		return err
	}

	r.initializedClasses[className] = struct{}{}
	log.Infof("initialized %s", className)
	return nil
}

func (r *Runner) runMethod(ctx context.Context, code []byte, c class.Class, method class.Method, parameters []Value) error {
	log := logger.FromContext(ctx)

	name, err := c.ConstantPool.GetUtf8(method.NameIndex)
	if err != nil {
		return err
	}

	log.Infof("running %s %s %s % x", c.Name, name, parameters, code)
	r.stack.Push(c.Name, method, c.ConstantPool, make([]Value, 0), parameters)

	returnPc := r.pc
	r.pc = 0

	err = r.run(ctx, code)
	if err != nil {
		err = fmt.Errorf("%v\n\t%s.%s()", err, strings.ReplaceAll(c.Name, "/", "."), name)
		return err
	}

	err = r.stack.Pop()
	if err != nil {
		return err
	}

	r.pc = returnPc
	return nil
}

func (r *Runner) runNative(ctx context.Context, c class.Class, method *class.Method, operands []Value) (Value, error) {
	descriptor, err := c.ConstantPool.GetUtf8(method.DescriptorIndex)
	if err != nil {
		return nil, err
	}

	methodDescriptor, err := class.NewMethodDescriptor(descriptor)
	if err != nil {
		return nil, err
	}

	methodName, err := c.ConstantPool.GetUtf8(method.NameIndex)
	if err != nil {
		return nil, err
	}

	if c.Name == "java/lang/System" &&
		methodName == "registerNatives" &&
		methodDescriptor.ReturnDescriptor == 'V' &&
		len(methodDescriptor.Parameters) == 0 {

		method, ok, err := c.GetMethod("initPhase1")
		if err != nil {
			return nil, err
		}

		if !ok {
			return nil, errors.New("method 'initPhase1' not found")
		}

		code, err := method.CodeAttribute()
		if err != nil {
			return nil, err
		}

		return nil, r.runMethod(ctx, code.Code, c, *method, make([]Value, 0))
	} else {
		return nil, errors.New("native method not implemented")
	}
}
