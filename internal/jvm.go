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
	currentClass *class.Class
	pc           int
	returnPc     int
	loader       loader.Loader
}

func NewRunner() Runner {
	return Runner{
		currentClass: nil,
		pc:           0,
		returnPc:     0,
		loader:       loader.NewLoader(),
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

	return r.runMethod(main)
}

const GET_STATIC = 0xb2
const INVOKE_STATIC = 0xb8

func (r *Runner) run(code []byte) error {
	for {
		instruction := code[r.pc]

		var err error
		switch instruction {
		case GET_STATIC:
			err = r.getStatic(code)
		case INVOKE_STATIC:
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
	index := (uint16(code[r.pc+1]) | uint16(code[r.pc+2]))
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

	err = r.resolveField(refInfo)
	if err != nil {
		return err
	}

	return errors.New("not implemented: getstatic")
}

func (r *Runner) invokeStatic(code []byte) error {
	index := (uint16(code[r.pc+1]) | uint16(code[r.pc+2]))
	r.pc += 2

	log.Println(index)
	log.Println(r.currentClass.Name)

	_, err := r.currentClass.ConstantPool.Ref(index)
	if err != nil {
		return err
	}

	return errors.New("not implemented: invokestatic")
}

func (r *Runner) initializeClass(className string) error {
	c, err := r.loader.Load(className)
	if err != nil {
		return err
	}

	r.currentClass = c

	clinit, ok, err := c.GetClinitMethod()
	if !ok {
		return nil
	}

	if err != nil {
		return err
	}

	log.Printf("running <clinit> for %s", className)

	return r.runMethod(clinit)
}

func (r *Runner) runMethod(method *class.Method) error {
	codeAttribute, err := method.CodeAttribute()
	if err != nil {
		return err
	}

	r.returnPc = r.pc
	r.pc = 0

	err = r.run(codeAttribute.Code)
	if err != nil {
		return err
	}

	r.pc = r.returnPc
	return nil
}

func (r *Runner) resolveField(refInfo *class.RefInfo) error {
	return errors.New("not implemented: field resolution")
}
