package main

import "github.com/sbinet/play"

func main() {
	play.RunPy(code)
}

const code = `
# START OMIT
import threading


def worker():
    """thread worker function"""
    print('Worker')


threads = []
for i in range(5):
    t = threading.Thread(target=worker)
    threads.append(t)
    t.start()
# STOP OMIT
`
