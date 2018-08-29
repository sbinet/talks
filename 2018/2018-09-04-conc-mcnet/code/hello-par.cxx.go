package main

import "github.com/sbinet/play"

func main() {
	play.RunCxx(code)
}

const code = `
// START OMIT
#include <thread>
#include <iostream>

void func() {
	std::cout << "** inside thread "
			  << std::this_thread::get_id() // HL
			  << "!" << std::endl;
}

int main(int argc, char **argv) {
	std::thread t(func); // create and schedule thread to execute 'func' // HL

	t.join(); // wait for thread(s) to finish // HL
	return 0;
}
// STOP OMIT
`
