package stack

import (
	"errors"
	"fmt"

	"github.com/m4tthewde/swell/internal/class"
)

type Frame struct {
	className      string
	method         class.Method
	constantPool   class.ConstantPool
	operands       []Value
	localVariables []Value
}

func NewFrame(
	className string,
	method class.Method,
	constantPool class.ConstantPool,
	operands []Value,
	localVariables []Value,
) Frame {
	return Frame{className: className, method: method, constantPool: constantPool, operands: make([]Value, 0), localVariables: localVariables}
}

type Stack struct {
	frames []Frame
}

func NewStack() Stack {
	return Stack{frames: make([]Frame, 0)}
}
func (s *Stack) Push(
	className string,
	method class.Method,
	constantPool class.ConstantPool,
	localVariables []Value,
) {
	frame := NewFrame(className, method, constantPool, []Value{}, localVariables)
	s.frames = append(s.frames, frame)
}

func (s *Stack) Pop() error {
	if len(s.frames) == 0 {
		return errors.New("stack is empty")
	}

	s.frames = s.frames[:len(s.frames)-1]
	return nil
}

func (s *Stack) activeFrame() (*Frame, error) {
	if len(s.frames) == 0 {
		return nil, errors.New("stack is empty")
	}

	frame := s.frames[len(s.frames)-1]
	return &frame, nil
}

func (s *Stack) PopOperands(count int) ([]Value, error) {
	frame, err := s.activeFrame()
	if err != nil {
		return nil, err
	}

	operands := frame.operands[len(frame.operands)-count:]

	frame.operands = frame.operands[:len(frame.operands)-count]
	s.frames[len(s.frames)-1] = *frame
	return operands, nil
}

func (s *Stack) PushOperand(operand Value) error {
	frame, err := s.activeFrame()
	if err != nil {
		return err
	}

	frame.operands = append(frame.operands, operand)
	s.frames[len(s.frames)-1] = *frame

	return nil
}

func (s *Stack) GetOperand() (Value, error) {
	frame, err := s.activeFrame()
	if err != nil {
		return nil, err
	}
	return frame.operands[len(frame.operands)-1], nil
}

func (s *Stack) GetLocalVariable(n int) (Value, error) {
	frame, err := s.activeFrame()
	if err != nil {
		return nil, err
	}

	if n >= len(frame.localVariables) {
		return nil, fmt.Errorf("no localvariable at %d, len is %d", n, len(frame.localVariables))
	}

	return frame.localVariables[n], nil
}

func (s *Stack) CurrentConstantPool() (*class.ConstantPool, error) {
	frame, err := s.activeFrame()
	if err != nil {
		return nil, err
	}

	return &frame.constantPool, nil
}

func (s *Stack) CurrentMethod() (*class.Method, error) {
	frame, err := s.activeFrame()
	if err != nil {
		return nil, err
	}

	return &frame.method, nil
}
