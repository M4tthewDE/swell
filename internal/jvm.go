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

	codeAttribute, err := method.CodeAttribute()
	if err != nil {
		return err
	}

	runner := NewRunner(codeAttribute.Code, *class)
	return runner.run()
}

type Runner struct {
	currentClass class.Class
	code         []byte
	pc           int
	loader       loader.Loader
}

func NewRunner(code []byte, class class.Class) Runner {
	return Runner{code: code, pc: 0, currentClass: class}

}

const GET_STATIC = 0xb2

func (r *Runner) run() error {
	for {
		instruction := r.code[r.pc]

		switch instruction {
		case GET_STATIC:
			return r.getStatic()
		default:
			return errors.New(
				fmt.Sprintf("unknown instruction %x", instruction),
			)
		}
	}
}

func (r *Runner) getStatic() error {
	index := (uint16(r.code[r.pc+1]) | uint16(r.code[r.pc+2]))
	r.pc += 2

	refInfo, err := r.currentClass.ConstantPool.Ref(index)
	if err != nil {
		return err
	}

	err = r.resolve(refInfo)
	if err != nil {
		return err
	}

	return errors.New("not implemented: getstatic")
}

func (r *Runner) resolve(refInfo *class.RefInfo) error {
	classInfo, err := r.currentClass.ConstantPool.Class(refInfo.ClassIndex)
	if err != nil {
		return err
	}

	className, err := r.currentClass.ConstantPool.GetUtf8(classInfo.NameIndex)
	if err != nil {
		return err
	}

	err = r.loader.Load(className)
	if err != nil {
		return err
	}

	return errors.New("not implemented: field resolution")
}
