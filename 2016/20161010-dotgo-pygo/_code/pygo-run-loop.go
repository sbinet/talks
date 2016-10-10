func (interp *Interpreter) Run(code Code) {
	prog := code.Prog
	for pc := 0; pc < len(prog); pc++ {
		op := prog[pc].(Opcode)
		switch op {
		case OpLoadValue:
			pc++
			val := code.Numbers[prog[pc].(int)]
			interp.stack.push(val)
		case OpAdd:
			lhs := interp.stack.pop()
			rhs := interp.stack.pop()
			sum := lhs + rhs
			interp.stack.push(sum)
		case OpPrint:
			val := interp.stack.pop()
			fmt.Println(val)
		case OpLoadName:
			pc++
			// ...
		}
	}
}
