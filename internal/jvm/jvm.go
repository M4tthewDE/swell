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
	currentClass          *class.Class
	classBeingInitialized string
	pc                    int
	returnPc              int
	loader                loader.Loader
	stack                 Stack
	heap                  Heap
}

func NewRunner() Runner {
	return Runner{
		currentClass:          nil,
		classBeingInitialized: "",
		pc:                    0,
		returnPc:              0,
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

	r.currentClass = c

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

	err = r.runMethod(ctx, code.Code, "main", make([]Value, 0))
	if err != nil {
		r.printStacktrace(ctx)
	}

	return err
}

const ALOAD_0 = 0x2a
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
		case ALOAD_0:
			err = aload(r, 0)
		case GET_STATIC:
			log.Info("executing getstatic")
			err = getStatic(r, ctx, code)
		case INVOKE_STATIC:
			log.Info("executing invokestatic")
			err = invokeStatic(r, ctx, code)
		case NEW:
			log.Info("executing new")
			err = new(r, ctx, code)
		case DUP:
			log.Info("executing dup")
			err = dup(r)
		case INVOKE_SPECIAL:
			log.Info("executing invokespecial")
			err = invokeSpecial(r, ctx, code)
		default:
			err = errors.New(
				fmt.Sprintf("unknown instruction %x", instruction),
			)
		}

		if err != nil {
			return err
		}

		if r.pc >= len(code)-1 {
			return nil
		}

	}
}

func (r *Runner) initializeClass(ctx context.Context, className string) error {
	if r.classBeingInitialized == className {
		return nil
	} else {
		r.classBeingInitialized = className
	}

	c, err := r.loader.Load(ctx, className)
	if err != nil {
		return err
	}

	r.currentClass = c

	clinit, ok, err := r.currentClass.GetMethod("<clinit>")
	if !ok {
		return nil
	}

	if err != nil {
		return err
	}

	code, err := clinit.CodeAttribute()
	if err != nil {
		return err
	}

	return r.runMethod(ctx, code.Code, "clinit", make([]Value, 0))
}

func (r *Runner) printStacktrace(ctx context.Context) {
	log := logger.FromContext(ctx)

	stackTrace := "\n"
	for i := len(r.stack.frames) - 1; i >= 0; i-- {
		frame := r.stack.frames[i]
		className := strings.ReplaceAll(frame.className, "/", ".")
		stackTrace += fmt.Sprintf("\t%s.%s()\n", className, frame.methodName)
	}

	stackTrace = stackTrace[:len(stackTrace)-1]
	log.Info(stackTrace)
}

func (r *Runner) runMethod(ctx context.Context, code []byte, name string, parameters []Value) error {
	log := logger.FromContext(ctx)

	log.Infof("running method '%s'", name)
	log.Debugf("with code % x", code)
	log.Debugf("and %d parameters", len(parameters))
	r.stack.Push(r.currentClass.Name, name, parameters)

	oldClass := r.currentClass
	r.returnPc = r.pc
	r.pc = 0

	err := r.run(ctx, code)
	if err != nil {
		return err
	}

	r.currentClass = oldClass
	r.pc = r.returnPc
	return nil
}

func (r *Runner) runNative(ctx context.Context, method *class.Method, operands []Value) (Value, error) {
	descriptor, err := r.currentClass.ConstantPool.GetUtf8(method.DescriptorIndex)
	if err != nil {
		return nil, err
	}

	methodDescriptor, err := class.NewMethodDescriptor(descriptor)
	if err != nil {
		return nil, err
	}

	methodName, err := r.currentClass.ConstantPool.GetUtf8(method.NameIndex)
	if err != nil {
		return nil, err
	}

	if r.currentClass.Name == "java/lang/System" &&
		methodName == "registerNatives" &&
		methodDescriptor.ReturnDescriptor == 'V' &&
		len(methodDescriptor.Parameters) == 0 {

		method, ok, err := r.currentClass.GetMethod("initPhase1")
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

		return nil, r.runMethod(ctx, code.Code, "initPhase1", make([]Value, 0))
	} else {
		return nil, errors.New("native method not implemented")
	}
}
