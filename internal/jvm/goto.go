package jvm

func goTo(r *Runner, code []byte) error {
	r.pc += 1
	index := (uint16(code[r.pc])<<8 | uint16(code[r.pc+1]))
	r.pc += int(index)
	return nil
}
