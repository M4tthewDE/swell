package internal

type Value interface {
	isValue()
}

type BooleanValue struct {
	value bool
}

type ByteValue struct {
	value uint8
}

type ShortValue struct {
	value uint16
}

type IntValue struct {
	value uint32
}

type LongValue struct {
	value uint64
}

type CharValue struct {
	value rune
}

type FloatValue struct {
	value float32
}

type DoubleValue struct {
	value float64
}

func (b BooleanValue) isValue() {}
func (b ByteValue) isValue()    {}
func (b ShortValue) isValue()   {}
func (b IntValue) isValue()     {}
func (b LongValue) isValue()    {}
func (b CharValue) isValue()    {}
func (b FloatValue) isValue()   {}
func (b DoubleValue) isValue()  {}

type Frame struct {
	methodName string
	operands   []Value
}

func NewFrame(methodName string, operands []Value) Frame {
	return Frame{methodName: methodName, operands: operands}
}

type Stack struct {
	frames []Frame
}

func NewStack() Stack {
	return Stack{frames: make([]Frame, 0)}
}

func (s *Stack) Push(methodName string, operands []Value) {
	frame := NewFrame(methodName, operands)
	s.frames = append(s.frames, frame)
}

func (s *Stack) PopOperands(count int) []Value {
	frame := s.frames[len(s.frames)-1]
	return frame.operands
}

func (s *Stack) PushOperand(operand Value) {
	frame := s.frames[len(s.frames)-1]
	frame.operands = append(frame.operands, operand)
}
