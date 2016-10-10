func run(instructions []instruction) {
	for _, instruction := range instructions {
		switch inst := instruction.(type) {
		case opADD:
			// perform a+b
		case opPRINT:
			// print values
			// ...
		}
	}
}

