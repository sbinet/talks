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
	fortran := `
c     START OMIT
c     == hello.f ==
      program main
      implicit none
      write ( *, '(a)' ) 'Hello from FORTRAN'
      stop
      end
c     END OMIT
`

	dir, err := ioutil.TempDir("", "go-compile-fortran-")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(dir)
	err = os.Chdir(dir)
	if err != nil {
		log.Fatal(err)
	}

	err = ioutil.WriteFile("hello.f", []byte(fortran), 0644)
	if err != nil {
		log.Fatal(err)
	}

	beg := time.Now()
	defer func() {
		fmt.Printf("took: %v\n", time.Since(beg))
	}()

	err = exec.Command("gfortran", "-c", "hello.f").Run()
	if err != nil {
		log.Fatal(err)
	}

	err = exec.Command("gfortran", "-o", "hello", "hello.o").Run()
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
