package internal

import (
	"errors"
	"fmt"
	"log"

	"github.com/m4tthewde/swell/internal/class"
	"github.com/m4tthewde/swell/internal/loader"
)

func Run(className string) error {
	runner := NewRunner()

	return runner.runMain(className)
}

type Runner struct {
	currentClass          *class.Class
	classBeingInitialized string
	pc                    int
	returnPc              int
	loader                loader.Loader
	stack                 Stack
}

func NewRunner() Runner {
	return Runner{
		currentClass:          nil,
		classBeingInitialized: "",
		pc:                    0,
		returnPc:              0,
		loader:                loader.NewLoader(),
		stack:                 NewStack(),
	}

}

func (r *Runner) runMain(className string) error {
	err := r.initializeClass(className)
	if err != nil {
		return err
	}

	c, err := r.loader.Load(className)
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

	return r.runMethod(code.Code, "main", make([]Value, 0))
}

const GET_STATIC = 0xb2
const INVOKE_STATIC = 0xb8

func (r *Runner) run(code []byte) error {
	for {
		instruction := code[r.pc]

		var err error
		switch instruction {
		case GET_STATIC:
			log.Println("getstatic")
			err = r.getStatic(code)
		case INVOKE_STATIC:
			log.Println("invokestatic")
			err = r.invokeStatic(code)
		default:
			err = errors.New(
				fmt.Sprintf("unknown instruction %x", instruction),
			)
		}

		return err
	}
}

func (r *Runner) getStatic(code []byte) error {
	index := (uint16(code[r.pc+1])<<8 | uint16(code[r.pc+2]))
	r.pc += 2

	refInfo, err := r.currentClass.ConstantPool.Ref(index)
	if err != nil {
		return err
	}

	classInfo, err := r.currentClass.ConstantPool.Class(refInfo.ClassIndex)
	if err != nil {
		return err
	}

	className, err := r.currentClass.ConstantPool.GetUtf8(classInfo.NameIndex)
	if err != nil {
		return err
	}

	err = r.initializeClass(className)
	if err != nil {
		return err
	}

	nameAndType, err := r.currentClass.ConstantPool.NameAndType(refInfo.NameAndTypeIndex)
	if err != nil {
		return err
	}

	fieldName, err := r.currentClass.ConstantPool.GetUtf8(nameAndType.NameIndex)
	if err != nil {
		return err
	}

	field, ok, err := r.currentClass.GetFIeld(fieldName)
	if err != nil {
		return err
	}

	if !ok {
		return errors.New("static field not found")
	}

	log.Println(field)
	log.Println(fieldName)

	return errors.New("not implemented: getstatic")
}

func (r *Runner) invokeStatic(code []byte) error {
	index := (uint16(code[r.pc+1])<<8 | uint16(code[r.pc+2]))
	r.pc += 2

	ref, err := r.currentClass.ConstantPool.Ref(index)
	if err != nil {
		return err
	}

	c, err := r.currentClass.ConstantPool.Class(ref.ClassIndex)
	if err != nil {
		return err
	}

	className, err := r.currentClass.ConstantPool.GetUtf8(c.NameIndex)
	if err != nil {
		return err
	}

	if err = r.initializeClass(className); err != nil {
		return err
	}

	nameAndType, err := r.currentClass.ConstantPool.NameAndType(ref.NameAndTypeIndex)
	if err != nil {
		return err
	}

	methodName, err := r.currentClass.ConstantPool.GetUtf8(nameAndType.NameIndex)
	if err != nil {
		return err
	}

	log.Println(methodName)

	method, ok, err := r.currentClass.GetMethod(methodName)
	if err != nil {
		return err
	}

	if !ok {
		return errors.New("method not found")
	}

	descriptor, err := r.currentClass.ConstantPool.GetUtf8(nameAndType.DescriptorIndex)
	if err != nil {
		return err
	}

	methodDescriptor, err := class.NewMethodDescriptor(descriptor)
	if err != nil {
		return err
	}

	operands := r.popOperands(len(methodDescriptor.Parameters))

	if !method.IsNative() {
		code, err := method.CodeAttribute()
		if err != nil {
			return err
		}

		return r.runMethod(code.Code, methodName, operands)
	} else {
		val, err := r.RunNative(method, operands)
		if err != nil {
			return err
		}

		if val != nil {
			r.stack.PushOperand(val)
		}
	}

	return nil
}

func (r *Runner) popOperands(count int) []Value {
	return r.stack.PopOperands(count)
}

func (r *Runner) initializeClass(className string) error {
	if r.classBeingInitialized == className {
		return nil
	} else {
		r.classBeingInitialized = className
	}

	c, err := r.loader.Load(className)
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

	log.Printf("running <clinit> for %s", className)

	return r.runMethod(code.Code, "clinit", make([]Value, 0))
}

func (r *Runner) runMethod(code []byte, name string, parameters []Value) error {
	r.stack.Push(name, parameters)

	r.returnPc = r.pc
	r.pc = 0

	err := r.run(code)
	if err != nil {
		return err
	}

	r.pc = r.returnPc
	return nil
}

func (r *Runner) RunNative(method *class.Method, operands []Value) (Value, error) {
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

		return nil, r.runMethod(code.Code, "initPhase1", make([]Value, 0))
	} else {
		return nil, errors.New("native method not implemented")
	}
}
