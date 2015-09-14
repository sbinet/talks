package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

var (
	localhost = "localhost"
)

func init() {
	switch runtime.GOOS {
	case "darwin", "windows":
		out, err := exec.Command("docker-machine", "ip", "default").Output()
		if err != nil {
			log.Fatalf("error retrieving ip-address of docker-VM: %v\n", err)
		}
		localhost = string(bytes.Trim(out, "\n"))
	}
}

func main() {
	origdir, err := os.Getwd()
	if err != nil {
		log.Fatalf("could not fetch current working directory: %v\n", err)
	}

	tmpdir, err := ioutil.TempDir("", "ecole-ci-")
	if err != nil {
		log.Fatalf("could not create tempdir: %v\n", err)
	}
	defer os.RemoveAll(tmpdir)

	err = os.Chdir(tmpdir)
	if err != nil {
		log.Fatalf("could not chdir to [%s]: %v\n", tmpdir, err)
	}

	err = os.MkdirAll(filepath.Join(tmpdir, "registry"), 0755)
	if err != nil {
		log.Fatalf("error creating 'registry' dir: %v\n", err)
	}

	// launch a registry if not there
	err = exec.Command("docker", "inspect", "ecole-registry").Run()
	if err != nil {
		go func() {
			cmd := exec.Command(
				"docker", "run",
				"-p", "5000:5000", "-v", tmpdir+"/registry:/var/lib/registry",
				"--name=ecole-registry", "registry",
			)
			err = cmd.Run()
			if err != nil {
				log.Fatalf("could not launch ecole-registry: %v\n", err)
			}
		}()
	}

	// clone repository
	cmd := exec.Command(
		"git", "clone", "git@gitlab.in2p3.fr:EcoleInfo2015/TP.git", "TP",
	)
	err = run(cmd)
	if err != nil {
		log.Fatalf("could not retrieve git repo!")
	}

	repodir := filepath.Join(tmpdir, "TP")
	err = os.Chdir(repodir)
	if err != nil {
		log.Fatalf("could not chdir to [%s]: %v\n", repodir, err)
	}

	err = os.MkdirAll(filepath.Join(repodir, "docker-web-base"), 0755)
	if err != nil {
		log.Fatalf("error creating [docker-web-base]: %v\n", err)
	}

	err = cp(
		filepath.Join(repodir, "docker-web-base", "Dockerfile"),
		filepath.Join(origdir, "base-dockerfile"),
	)
	if err != nil {
		log.Fatalf("error creating docker-web-base/Dockerfile: %v\n", err)
	}

	cmd = exec.Command(
		"docker", "build", "-t", localhost+":5000/binet/web-base", ".",
	)
	cmd.Dir = filepath.Join(repodir, "docker-web-base")
	err = run(cmd)
	if err != nil {
		log.Fatalf("error build web-base image: %v\n", err)
	}

	cmd = exec.Command(
		"docker", "push", localhost+":5000/binet/web-base",
	)
	err = run(cmd)
	if err != nil {
		log.Fatalf("error pushing web-base")
	}

	err = cp(
		filepath.Join(repodir, "Dockerfile"),
		filepath.Join(origdir, "webapp-dockerfile"),
	)
	if err != nil {
		log.Fatalf("error creating web-app/Dockerfile: %v\n", err)
	}

	cmd = exec.Command(
		"docker", "build", "-t", localhost+":5000/binet/web-app", ".",
	)
	cmd.Dir = repodir
	err = run(cmd)
	if err != nil {
		log.Fatalf("error build web-app: %v\n", err)
	}

	_ = run(exec.Command("docker", "kill", "binet-web-app"))
	_ = run(exec.Command("docker", "rm", "binet-web-app"))
	cmd = exec.Command(
		"docker", "run", "-d", "-p", "8080:8080",
		"--name=binet-web-app",
		localhost+":5000/binet/web-app",
	)
	err = run(cmd)
	if err != nil {
		log.Fatalf("error running web-app: %v\n", err)
	}

	cmd = exec.Command(
		"docker", "push", localhost+":5000/binet/web-app",
	)
	err = run(cmd)
	if err != nil {
		log.Fatalf("error pushing web-app: %v\n", err)
	}

	// FIXME(sbinet): we rely on the docker-push to take some time
	// so the http.Get will see a container exposing a (running) web server...
	resp, err := http.Get("http://" + localhost + ":8080/")
	if err != nil {
		log.Printf("response:\n%#v\n", resp)
		log.Fatalf("error GET-localhost:8080: %v\n", err)
	}
	defer resp.Body.Close()
	hello := new(bytes.Buffer)
	_, err = io.Copy(hello, resp.Body)
	if err != nil {
		log.Fatalf("error printing resp.body: %v\n", err)
	}

	fmt.Printf("===\n%v===\n", string(hello.Bytes()))
	if !bytes.HasPrefix(hello.Bytes(), []byte("<h1>Bienvenue ")) {
		log.Fatalf("invalid homepage:\n%v\n", string(hello.Bytes()))
	}

}

func run(cmd *exec.Cmd) error {
	if cmd.Stdin == nil {
		cmd.Stdin = os.Stdin
	}
	if cmd.Stdout == nil {
		cmd.Stdout = os.Stdout
	}
	if cmd.Stderr == nil {
		cmd.Stderr = os.Stderr
	}
	return cmd.Run()
}

func cp(dst, src string) error {
	fdst, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer fdst.Close()

	fsrc, err := os.Open(src)
	if err != nil {
		return err
	}
	defer fsrc.Close()

	_, err = io.Copy(fdst, fsrc)
	if err != nil {
		return err
	}

	return fdst.Close()
}
