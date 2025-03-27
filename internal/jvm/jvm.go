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

	err = r.runMethod(ctx, code, *c, *main, make([]stack.Value, 0))
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
const PutStatic = 0xb3
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
const IReturn = 0xac
const IfNe = 0x9a
const GoTo = 0xa7
const LdcWide = 0x13
const PutField = 0xb5
const ArrayLength = 0xbe
const IfEq = 0x99
const IntShiftRight = 0x7a
const IStore2 = 0x3d
const ILoad0 = 0x1a
const ILoad1 = 0x1b
const ILoad2 = 0x1c
const ISub = 0x64
const BiPush = 0x10
const IfLt = 0x9b
const IfICmpLt = 0xa1
const NewArray = 0xbc

func (r *Runner) run(ctx context.Context, code []byte) error {
	log := logger.FromContext(ctx)

	for {
		instruction := code[r.pc]
		log.Debugw("instruction", "pc", r.pc)

		var err error
		switch instruction {
		case GetField:
			log.Debug("getfield")
			err = getField(r, ctx, code)
		case InvokeVirtual:
			log.Debug("invokevirtual")
			err = invokeVirtual(r, ctx, code)
		case LdcOp:
			log.Debug("ldc")
			err = ldcNormal(r, ctx, code)
		case LdcWide:
			log.Debug("ldc_w")
			err = ldcWide(r, ctx, code)
		case Aload0:
			log.Debug("aload_0")
			err = aload(ctx, r, 0)
		case Aload1:
			log.Debug("aload_1")
			err = aload(ctx, r, 1)
		case Aload2:
			log.Debug("aload_2")
			err = aload(ctx, r, 2)
		case Aload3:
			log.Debug("aload_3")
			err = aload(ctx, r, 3)
		case GetStaticOp:
			log.Debug("getstatic")
			err = getStatic(r, ctx, code)
		case InvokeStaticOp:
			log.Debug("invokestatic")
			err = invokeStatic(r, ctx, code)
		case NewOp:
			log.Debug("new")
			err = new(r, ctx, code)
		case DupOp:
			log.Debug("dup")
			err = dup(r)
		case InvokeSpecialOp:
			log.Debug("invokespecial")
			err = invokeSpecial(r, ctx, code)
		case RetOp:
			log.Debug("ret")
			return ret(r)
		case AReturn:
			log.Debug("areturn")
			return areturn(r)
		case Astore0:
			log.Debug("astore_0")
			err = astore(ctx, r, 0)
		case Astore1:
			log.Debug("astore_1")
			err = astore(ctx, r, 1)
		case Astore2:
			log.Debug("astore_2")
			err = astore(ctx, r, 2)
		case Astore3:
			log.Debug("astore_3")
			err = astore(ctx, r, 3)
		case IfNonNull:
			log.Debug("ifnonull")
			err = ifnonnull(r, code)
		case IConstM1:
			log.Debug("iconst_m1")
			err = iconst(r, -1)
		case IConst0:
			log.Debug("iconst_0")
			err = iconst(r, 0)
		case IConst1:
			log.Debug("iconst_1")
			err = iconst(r, 1)
		case IConst2:
			log.Debug("iconst_2")
			err = iconst(r, 2)
		case IConst3:
			log.Debug("iconst_3")
			err = iconst(r, 3)
		case IConst4:
			log.Debug("iconst_4")
			err = iconst(r, 4)
		case IConst5:
			log.Debug("iconst_5")
			err = iconst(r, 5)
		case ANewArray:
			log.Debug("anwarray")
			err = anewarray(r, ctx, code)
		case PutStatic:
			log.Debug("putstatic")
			err = putstatic(r, ctx, code)
		case Nop:
			log.Debug("nop")
			r.pc += 1
		case IReturn:
			log.Debug("ireturn")
			return ireturn(r)
		case IfNe:
			log.Debug("ifne")
			err = ifne(r, code)
		case GoTo:
			log.Debug("goto")
			err = goTo(r, code)
		case PutField:
			log.Debug("putfield")
			err = putField(r, ctx, code)
		case ArrayLength:
			log.Debug("arraylength")
			err = arrayLength(r)
		case IfEq:
			log.Debug("ifeq")
			err = ifeq(r, code)
		case IntShiftRight:
			log.Debug("ishr")
			err = intShiftRight(r)
		case IStore2:
			log.Debug("istore_2")
			err = istore(ctx, r, 2)
		case ILoad0:
			log.Debug("iload_0")
			err = iload(ctx, r, 0)
		case ILoad1:
			log.Debug("iload_1")
			err = iload(ctx, r, 1)
		case ILoad2:
			log.Debug("iload_2")
			err = iload(ctx, r, 2)
		case ISub:
			log.Debug("isub")
			err = isub(r)
		case BiPush:
			log.Debug("bipush")
			value := code[r.pc+1]
			r.pc += 2
			err = r.stack.PushOperand(stack.IntValue{Value: int32(value)})
		case IfLt:
			log.Debug("iflt")
			err = iflt(r, code)
		case IfICmpLt:
			log.Debug("if_icmplt")
			err = ifICmpLt(r, code)
		case NewArray:
			log.Debug("newarray")
			err = newArray(r, ctx, code)
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
		log.Debugw("already initialized", "className", className)
		return nil
	}

	if r.classBeingInitialized == className {
		return nil
	} else {
		r.classBeingInitialized = className
	}

	log.Infow("initializing", "className", className)

	c, err := r.loader.Load(ctx, className)
	if err != nil {
		return err
	}

	clinit, ok, err := c.GetMethod("<clinit>")
	if !ok {
		r.initializedClasses[className] = struct{}{}
		log.Infow("initialized without clinit", "className", className)
		return nil
	}

	if err != nil {
		return err
	}

	code, err := clinit.CodeAttribute()
	if err != nil {
		return err
	}

	err = r.runMethod(ctx, code, *c, *clinit, make([]stack.Value, 0))
	if err != nil {
		return err
	}

	r.initializedClasses[className] = struct{}{}
	log.Infow("initialized", "className", className)
	return nil
}

func (r *Runner) runMethod(ctx context.Context, code *class.CodeAttribute, c class.Class, method class.Method, parameters []stack.Value) error {
	log := logger.FromContext(ctx)

	name, err := c.ConstantPool.GetUtf8(method.NameIndex)
	if err != nil {
		return err
	}

	log.Infow("executing method",
		"class", c.Name,
		"name", name,
		"parameters", fmt.Sprintf("%s", parameters),
		"code", fmt.Sprintf("% x", code.Code), // Use hex formatting for binary data
	)
	r.stack.Push(c.Name, method, c.ConstantPool, parameters)

	returnPc := r.pc
	r.pc = 0

	err = r.run(ctx, code.Code)
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
