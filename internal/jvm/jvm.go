package jvm

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/m4tthewde/swell/internal/class"
	"github.com/m4tthewde/swell/internal/jvm/stack"
	"github.com/m4tthewde/swell/internal/loader"
	"github.com/m4tthewde/swell/internal/logger"
)

type Runner struct {
	classBeingInitialized string
	initializedClasses    map[string]struct{}
	pc                    int
	loader                loader.Loader
	stack                 stack.Stack
	heap                  Heap
}

func NewRunner(classPath []string) Runner {
	return Runner{
		classBeingInitialized: "",
		// FIXME: is there a better data structure for this?
		initializedClasses: make(map[string]struct{}),
		pc:                 0,
		loader:             loader.NewLoader(classPath),
		stack:              stack.NewStack(),
		heap:               NewHeap(),
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

	err = r.runMethod(ctx, code.Code, *c, *main, make([]stack.Value, 0))
	if err != nil {
		return err
	}

	return nil
}

const Nop = 0x00
const IConstM1 = 0x2
const IConst0 = 0x3
const IConst1 = 0x4
const IConst2 = 0x5
const IConst3 = 0x6
const IConst4 = 0x7
const IConst5 = 0x8
const LdcOp = 0x12
const Aload0 = 0x2a
const Aload1 = 0x2b
const Aload2 = 0x2c
const Aload3 = 0x2d
const AReturn = 0xb0
const RetOp = 0xb1
const GetStaticOp = 0xb2
const GetField = 0xb4
const InvokeVirtual = 0xb6
const InvokeSpecialOp = 0xb7
const InvokeStaticOp = 0xb8
const NewOp = 0xbb
const ANewArray = 0xbd
const DupOp = 0x59
const Astore0 = 0x4b
const Astore1 = 0x4c
const Astore2 = 0x4d
const Astore3 = 0x4e
const IfNonNull = 0xc7

func (r *Runner) run(ctx context.Context, code []byte) error {
	log := logger.FromContext(ctx)

	for {
		instruction := code[r.pc]

		var err error
		switch instruction {
		case GetField:
			log.Info("getfield")
			err = getField(r, ctx, code)
		case InvokeVirtual:
			log.Info("invokevirtual")
			err = invokeVirtual(r, ctx, code)
		case LdcOp:
			log.Info("ldc")
			err = ldc(r, ctx, code)
		case Aload0:
			log.Info("aload_0")
			err = aload(r, 0)
		case Aload1:
			log.Info("aload_1")
			err = aload(r, 1)
		case Aload2:
			log.Info("aload_2")
			err = aload(r, 2)
		case Aload3:
			log.Info("aload_3")
			err = aload(r, 3)
		case GetStaticOp:
			log.Info("getstatic")
			err = getStatic(r, ctx, code)
		case InvokeStaticOp:
			log.Info("invokestatic")
			err = invokeStatic(r, ctx, code)
		case NewOp:
			log.Info("new")
			err = new(r, ctx, code)
		case DupOp:
			log.Info("dup")
			err = dup(r)
		case InvokeSpecialOp:
			log.Info("invokespecial")
			err = invokeSpecial(r, ctx, code)
		case RetOp:
			log.Info("ret")
			return ret(r)
		case AReturn:
			log.Info("areturn")
			return areturn(r)
		case Astore0:
			log.Info("astore_0")
			err = astore(r, 0)
		case Astore1:
			log.Info("astore_1")
			err = astore(r, 1)
		case Astore2:
			log.Info("astore_2")
			err = astore(r, 2)
		case Astore3:
			log.Info("astore_3")
			err = astore(r, 3)
		case IfNonNull:
			log.Info("ifnonull")
			err = ifnonnull(r)
		case IConstM1:
			log.Info("iconst_m1")
			err = iconst(r, -1)
		case IConst0:
			log.Info("iconst_0")
			err = iconst(r, 0)
		case IConst1:
			log.Info("iconst_1")
			err = iconst(r, 1)
		case IConst2:
			log.Info("iconst_2")
			err = iconst(r, 2)
		case IConst3:
			log.Info("iconst_3")
			err = iconst(r, 3)
		case IConst4:
			log.Info("iconst_4")
			err = iconst(r, 4)
		case IConst5:
			log.Info("iconst_5")
			err = iconst(r, 5)
		case ANewArray:
			log.Info("anwarray")
			err = anewarray(r, ctx, code)
		case Nop:
			log.Info("nop")
			r.pc += 1
		default:
			return fmt.Errorf("unknown instruction %x", instruction)

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

	err = r.runMethod(ctx, code.Code, *c, *clinit, make([]stack.Value, 0))
	if err != nil {
		return err
	}

	r.initializedClasses[className] = struct{}{}
	log.Infof("initialized %s", className)
	return nil
}

func (r *Runner) runMethod(ctx context.Context, code []byte, c class.Class, method class.Method, parameters []stack.Value) error {
	log := logger.FromContext(ctx)

	name, err := c.ConstantPool.GetUtf8(method.NameIndex)
	if err != nil {
		return err
	}

	log.Infof("running %s %s %s % x", c.Name, name, parameters, code)
	r.stack.Push(c.Name, method, c.ConstantPool, parameters)

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

func (r *Runner) runNative(ctx context.Context, c class.Class, method *class.Method, operands []stack.Value) (stack.Value, error) {
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

		return nil, r.runMethod(ctx, code.Code, c, *method, operands)
	} else if c.Name == "java/lang/Class" && methodName == "registerNatives" {
		return nil, nil
	} else {
		return nil, fmt.Errorf("native method %s in %s not implemented", methodName, c.Name)
	}
}
