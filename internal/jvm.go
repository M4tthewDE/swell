package internal

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/m4tthewde/swell/internal/class"
	"github.com/m4tthewde/swell/internal/loader"
)

func Run(className string) error {
	classFilePath := fmt.Sprintf("%s.class", className)

	file, err := os.Open(classFilePath)
	if err != nil {
		return err
	}

	defer file.Close()

	reader := bufio.NewReader(file)

	class, err := class.NewClass(reader)
	if err != nil {
		return err
	}

	method, err := class.GetMainMethod()
	if err != nil {
		return err
	}

	log.Println("executing main method")

	runner := NewRunner(*class)
	return runner.runMethod(method)
}

type Runner struct {
	currentClass class.Class
	pc           int
	returnPc     int
	loader       loader.Loader
}

func NewRunner(class class.Class) Runner {
	return Runner{
		currentClass: class,
		pc:           0,
		returnPc:     0,
		loader:       loader.NewLoader(),
	}

}

const GET_STATIC = 0xb2

func (r *Runner) run(code []byte) error {
	for {
		instruction := code[r.pc]

		switch instruction {
		case GET_STATIC:
			err := r.getStatic(code)
			if err != nil {
				return err
			}
		default:
			return errors.New(
				fmt.Sprintf("unknown instruction %x", instruction),
			)
		}
	}
}

func (r *Runner) getStatic(code []byte) error {
	index := (uint16(code[r.pc+1]) | uint16(code[r.pc+2]))
	r.pc += 2

	refInfo, err := r.currentClass.ConstantPool.Ref(index)
	if err != nil {
		return err
	}

	err = r.initializeClass(refInfo.ClassIndex)
	if err != nil {
		return err
	}

	err = r.resolveField(refInfo)
	if err != nil {
		return err
	}

	return errors.New("not implemented: getstatic")
}

func (r *Runner) initializeClass(classIndex uint16) error {
	classInfo, err := r.currentClass.ConstantPool.Class(classIndex)
	if err != nil {
		return err
	}

	className, err := r.currentClass.ConstantPool.GetUtf8(classInfo.NameIndex)
	if err != nil {
		return err
	}

	class, err := r.loader.Load(className)
	if err != nil {
		return err
	}

	log.Printf("running <clinit> for %s", className)
	clinit, err := class.GetClinitMethod()
	if err != nil {
		return err
	}

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
