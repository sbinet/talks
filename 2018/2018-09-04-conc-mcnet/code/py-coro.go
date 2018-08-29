package main

import "github.com/sbinet/play"

func main() {
	play.RunPy(code)
}

const code = `
# START OMIT
def gen(n):
	i = 0
	while i < n:
		yield i
		i += 1

print(sum(gen(10)))
g = gen(10)
print(next(g),next(g),next(g))
print(next(g),next(g),next(g))
# STOP OMIT
`
