package main

import "github.com/sbinet/play"

func main() {
	play.RunPy(code)
}

const code = `
# START OMIT
import threading

def worker(num):
    """thread worker function"""
    print('Worker: {}'.format(num))
    return

threads = []
for i in range(5):
    t = threading.Thread(target=worker, args=(i,))
    threads.append(t)
    t.start()
# STOP OMIT
`
