package main

import "github.com/sbinet/play"

func main() {
	play.RunPy(code)
}

const code = `
# START OMIT
import sys
a = []
b = a
print("refcount(a) = {}".format(sys.getrefcount(a)))
# STOP OMIT
`
