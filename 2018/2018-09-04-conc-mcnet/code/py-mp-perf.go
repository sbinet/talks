package main

import "github.com/sbinet/play"

func main() {
	play.RunPy(code)
}

const code = `
# START OMIT
import time

COUNT = 50000000
def countdown(n):
    while n>0:
        n -= 1

start = time.time()
countdown(COUNT)
end = time.time()
print('Time taken in seconds -', end - start)
# STOP OMIT
`
