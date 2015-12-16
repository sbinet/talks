package main

import (
	"bytes"
	"encoding/json"
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
	doPush    = flag.Bool("do-push", false, "push to cc-ecole2015-docker")
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

	tmpdir, err := ioutil.TempDir("", "loops-tp-")
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
		"git", "clone", "git@github.com:sbinet/loops-20151217-tp.git",
	)
	err = run(cmd)
	if err != nil {
		log.Fatalf("could not retrieve git repo!")
	}

	repodir := filepath.Join(tmpdir, "loops-20151217-tp")
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

	if *doPush {
		cmd = exec.Command(
			"docker", "push", "binet/web-base",
		)
		err = run(cmd)
		if err != nil {
			log.Fatalf("error pushing web-base")
		}
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

	if *doPush {
		cmd = exec.Command(
			"docker", "push", "binet/web-app:v1",
		)
		err = run(cmd)
		if err != nil {
			log.Fatalf("error pushing web-app: %v\n", err)
		}
	}

	time.Sleep(2 * time.Second) // racing the web server...
	testWebServer("http://"+localhost+":8080/", "hello LoOPS 20151217!\n")

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
	fmt.Fprintf(w, "Welcome LoOPS 20151217!\n")
	fmt.Fprintf(w, "\n\n--- running external command...\n\n>>> pkg-config --cflags python2\n")
	cmd := exec.Command("pkg-config", "--cflags", "python2")
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

	if *doPush {
		cmd = exec.Command(
			"docker", "push", "binet/web-app:v2",
		)
		err = run(cmd)
		if err != nil {
			log.Fatalf("error pushing web-app:v2: %v\n", err)
		}
	}
	time.Sleep(2 * time.Second) // racing the web-servers...
	testWebServer("http://"+localhost+":8082/", "Welcome LoOPS ")
	testWebServer("http://"+localhost+":8080/", "hello LoOPS ")

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

func launchRegistry(tmpdir string) error {
	data := make([]struct {
		State struct {
			Running    bool
			Paused     bool
			Restarting bool
			OOMKilled  bool
			Dead       bool
			Pid        int
			ExitCode   int
			Error      string
			StartedAt  time.Time
			FinishedAt time.Time
		}
	}, 0)

	cmd := exec.Command(
		"docker", "inspect", "ecole-registry",
	)
	buf := new(bytes.Buffer)
	cmd.Stdout = buf
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	err := cmd.Run()
	if err != nil {
		log.Printf("starting a registry...\n")
		err = run(exec.Command("docker", "pull", "registry"))
		if err != nil {
			log.Printf("error retrieving registry image: %v\n", err)
			return err
		}

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
			// FIXME(sbinet): wait a bit for the server to get completely ready
			time.Sleep(5 * time.Second)
		}()
		return nil
	}

	err = json.Unmarshal(buf.Bytes(), &data)
	if err != nil {
		log.Printf("error unmarshaling JSON: %v\n", err)
		return err
	}

	if !data[0].State.Running {
		log.Printf("restarting registry...\n")
		go func() {
			cmd := exec.Command("docker", "restart", "ecole-registry")
			err = cmd.Run()
			if err != nil {
				log.Fatalf("could not restart ecole-registry: %v\n", err)
			}
		}()
		time.Sleep(5 * time.Second)
		return nil
	}
	return err
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
