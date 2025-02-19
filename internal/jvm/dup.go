package jvm

func dup(r *Runner) error {
	r.pc += 1
	operand := r.stack.GetOperand()
	r.stack.PushOperand(operand)
	return nil
}
