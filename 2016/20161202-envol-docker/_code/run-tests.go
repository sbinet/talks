package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"
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
	start := time.Now()
	defer func() {
		log.Printf("time= %v\n", time.Since(start))
	}()

	flag.Parse()

	origdir, err := os.Getwd()
	if err != nil {
		log.Fatalf("could not fetch current working directory: %v\n", err)
	}

	tmpdir, err := ioutil.TempDir("", "ji-docker-")
	if err != nil {
		log.Fatalf("could not create tempdir: %v\n", err)
	}
	defer os.RemoveAll(tmpdir)

	err = os.Chdir(tmpdir)
	if err != nil {
		log.Fatalf("could not chdir to [%s]: %v\n", tmpdir, err)
	}

	// clone repository
	cmd := exec.Command(
		"git", "clone", "git@github.com:sbinet/ji-docker-2016.git",
	)
	err = run(cmd)
	if err != nil {
		log.Fatalf("could not retrieve git repo!")
	}

	repodir := filepath.Join(tmpdir, "ji-docker-2016")
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
		"docker", "build", "-t", "binet/web-base", ".",
	)
	cmd.Dir = filepath.Join(repodir, "docker-web-base")
	err = run(cmd)
	if err != nil {
		log.Fatalf("error build web-base image: %v\n", err)
	}

	err = cp(
		filepath.Join(repodir, "Dockerfile"),
		filepath.Join(origdir, "webapp-dockerfile"),
	)
	if err != nil {
		log.Fatalf("error creating web-app/Dockerfile: %v\n", err)
	}

	log.Printf("tagging binet/web-app:v1...\n")
	cmd = exec.Command(
		"docker", "build", "-t", "binet/web-app:v1", ".",
	)
	cmd.Dir = repodir
	err = run(cmd)
	if err != nil {
		log.Fatalf("error build web-app: %v\n", err)
	}

	_ = exec.Command("docker", "kill", "binet-web-app").Run()
	_ = exec.Command("docker", "rm", "binet-web-app").Run()
	log.Printf("running binet-web-app (v1)...\n")
	cmd = exec.Command(
		"docker", "run", "-d", "-p", "8080:8080",
		"--name=binet-web-app",
		"binet/web-app:v1",
	)
	err = run(cmd)
	if err != nil {
		log.Fatalf("error running web-app: %v\n", err)
	}

	time.Sleep(2 * time.Second) // racing the web server...
	testWebServer("http://"+localhost+":8080/", "hello JI-2016!\n")

	// now create a v2
	const v2 = `package main

import (
	"fmt"
	"log"
	"net/http"
	"os/exec"
)

func main() {
	http.HandleFunc("/", rootHandle)
	port := ":8080"
	log.Printf("listening on: http://localhost%s\n", port)
	err := http.ListenAndServe("0.0.0.0"+port, nil)
	if err != nil {
		log.Fatalf("error closing web server: %v\n", err)
	}
}

func rootHandle(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome JI-2016!\n")
	fmt.Fprintf(w, "\n\n--- running external command...\n\n>>> pkg-config --version systemd\n")
	cmd := exec.Command("pkg-config", "--version", "systemd")
	cmd.Stdout = w
	cmd.Stderr = w
	err := cmd.Run()
	if err != nil {
		fmt.Fprintf(w, "error: %v\n", err)
	}
}`

	web, err := os.Create(filepath.Join(
		repodir,
		"web-app/main.go",
	))
	if err != nil {
		log.Fatalf("error opening file: %v\n", err)
	}
	defer web.Close()

	_, err = web.WriteString(v2)
	if err != nil {
		log.Fatalf("error updating [%s]: %v\n", web.Name(), err)
	}
	err = web.Close()
	if err != nil {
		log.Fatalf("error closing [%s]: %v\n", web.Name(), err)
	}

	// now create and run server-v2
	log.Printf("building binet/web-app:v2...\n")
	cmd = exec.Command(
		"docker", "build", "-t", "binet/web-app:v2", ".",
	)
	cmd.Dir = repodir
	err = run(cmd)
	if err != nil {
		log.Fatalf("error build web-app: %v\n", err)
	}

	_ = exec.Command("docker", "kill", "binet-web-app-v2").Run()
	_ = exec.Command("docker", "rm", "binet-web-app-v2").Run()
	log.Printf("running binet-web-app (v2)...\n")
	cmd = exec.Command(
		"docker", "run", "-d", "-p", "8082:8080",
		"--name=binet-web-app-v2",
		"binet/web-app:v2",
	)
	err = run(cmd)
	if err != nil {
		log.Fatalf("error running web-app: %v\n", err)
	}

	time.Sleep(2 * time.Second) // racing the web-servers...
	testWebServer("http://"+localhost+":8082/", "Welcome JI-2016")
	testWebServer("http://"+localhost+":8080/", "hello JI-2016")

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

func testWebServer(url string, prefix string) {
	resp, err := http.Get(url)
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
	if !bytes.HasPrefix(hello.Bytes(), []byte(prefix)) {
		log.Fatalf("invalid homepage:\n%v\n", string(hello.Bytes()))
	}
}
