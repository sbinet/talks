package main

import (
	"io/ioutil"
	"log"
	"os"
	"os/exec"
)

func main() {
	src := `
// START OMIT
package main

import "fmt"

func main() {
	lang := "Go"
	fmt.Printf("Hello from %s\n", lang)
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

	err = ioutil.WriteFile("hello.go", []byte(src), 0644)
	if err != nil {
		log.Fatal(err)
	}

	cmd := exec.Command("go", "run", "./hello.go")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}
