package main

import "github.com/sbinet/play"

func main() {
	play.RunPy(code)
}

const code = `
# START OMIT
import asyncio, random

async def produce(queue, n):
    for x in range(1, n + 1):
        print('producing {}/{}'.format(x, n))
        await asyncio.sleep(random.random())
        item = str(x)
        await queue.put(item)         // HL
    # indicate the producer is done
    await queue.put(None)             // HL

async def consume(queue):
    while True:
        item = await queue.get()      // HL
        if item is None:
            break
        print('consuming item {}...'.format(item))
        await asyncio.sleep(random.random())

loop = asyncio.get_event_loop()
queue = asyncio.Queue(loop=loop)
producer_coro = produce(queue, 10)
consumer_coro = consume(queue)
loop.run_until_complete(asyncio.gather(producer_coro, consumer_coro))
loop.close()
# STOP OMIT
`
