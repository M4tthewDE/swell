package jvm

func goTo(r *Runner, code []byte) error {
	index := (uint16(code[r.pc+1])<<8 | uint16(code[r.pc+2]))
	r.pc += int(index)
	return nil
}
