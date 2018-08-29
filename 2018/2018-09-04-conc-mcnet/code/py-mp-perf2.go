package main

import "github.com/sbinet/play"

func main() {
	play.RunPy(code)
}

const code = `
# START OMIT
from multiprocessing import Pool
import time

COUNT = 50000000
def countdown(n):
    while n>0:
        n -= 1

pool = Pool(processes=2)
start = time.time()
r1 = pool.apply_async(countdown, [COUNT//2])
r2 = pool.apply_async(countdown, [COUNT//2])
pool.close()
pool.join()
end = time.time()
print('Time taken in seconds -', end - start)
# STOP OMIT
`
