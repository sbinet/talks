package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"time"
)

func main() {
	src := `
// START OMIT
#include <iostream>
int main(int, char **) {
  std::cout << "Hello from C++" << std::endl;
  return 0;
}
// END OMIT
`

	dir, err := ioutil.TempDir("", "go-compile-")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(dir)
	err = os.Chdir(dir)
	if err != nil {
		log.Fatal(err)
	}

	err = ioutil.WriteFile("hello.cxx", []byte(src), 0644)
	if err != nil {
		log.Fatal(err)
	}

	beg := time.Now()
	defer func() {
		fmt.Printf("took: %v\n", time.Since(beg))
	}()

	err = exec.Command("c++", "-o", "hello", "hello.cxx").Run()
	if err != nil {
		log.Fatal(err)
	}

	cmd := exec.Command("./hello")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}
