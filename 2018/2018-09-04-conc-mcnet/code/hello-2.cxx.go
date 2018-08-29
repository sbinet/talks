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
	// (1) // HL
	std::cout << "hello" << std::endl;
}

int main(int argc, char **argv) {
	std::thread t(func);

	t.join(); // (2) // HL
	return 0;
}
// STOP OMIT
`
