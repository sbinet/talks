package main

func main() {
	code := Code{
		Prog: []Instruction{
			OpLoadValue, 0,
			OpStoreName, 0,
			OpLoadValue, 1,
			OpStoreName, 1,
			OpLoadName, 0,
			OpLoadName, 1,
			OpAdd,
			OpPrint,
		},
		Numbers: []int{1, 2},
		Names:   []string{"a", "b"},
	}

	interp := New()
	interp.Run(code)
}
