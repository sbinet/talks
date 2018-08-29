package main

import "github.com/sbinet/play"

func main() {
	play.RunPy(code)
}

const code = `
# START OMIT
import asyncio

async def outer():
    print('in outer')
    print('waiting for result1')
    result1 = await phase1()         # // HL
    print('waiting for result2')
    result2 = await phase2(result1)  # // HL
    return (result1, result2)

async def phase1():
    print('in phase1');
    return 'result1'

async def phase2(arg):
    print('in phase2')
    return 'result2 derived from {}'.format(arg)


event_loop = asyncio.get_event_loop()
try:
    return_value = event_loop.run_until_complete(outer())
    print('return value: {!r}'.format(return_value))
finally:
    event_loop.close()
# STOP OMIT
`
