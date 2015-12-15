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

	tmpdir, err := ioutil.TempDir("", "ecole-ci-")
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
		"docker", "build", "-t", "binet/web-base", ".",
	)
	cmd.Dir = filepath.Join(repodir, "docker-web-base")
	err = run(cmd)
	if err != nil {
		log.Fatalf("error build web-base image: %v\n", err)
	}

	cmd = exec.Command(
		"docker", "tag", "-f",
		"binet/web-base",
		"cc-ecole2015-docker.in2p3.fr:5000/binet/web-base",
	)
	err = run(cmd)
	if err != nil {
		log.Fatalf("error tagging web-base to cc-ecole/web-base: %v\n", err)
	}

	if *doPush {
		cmd = exec.Command(
			"docker", "push", "cc-ecole2015-docker.in2p3.fr:5000/binet/web-base",
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
	cmd = exec.Command(
		"docker", "run", "-d", "-p", "8080:8080",
		"--name=binet-web-app",
		"binet/web-app:v1",
	)
	err = run(cmd)
	if err != nil {
		log.Fatalf("error running web-app: %v\n", err)
	}

	cmd = exec.Command(
		"docker", "tag", "-f",
		"binet/web-app:v1",
		"cc-ecole2015-docker.in2p3.fr:5000/binet/web-app:v1",
	)
	err = run(cmd)
	if err != nil {
		log.Fatalf("error tagging web-app to cc-ecole/web-app: %v\n", err)
	}

	if *doPush {
		cmd = exec.Command(
			"docker", "push", "cc-ecole2015-docker.in2p3.fr:5000/binet/web-app:v1",
		)
		err = run(cmd)
		if err != nil {
			log.Fatalf("error pushing web-app: %v\n", err)
		}
	}

	// FIXME(sbinet): we rely on the docker-push to take some time
	// so the http.Get will see a container exposing a (running) web server...
	testWebServer("http://"+localhost+":8080/", "<h1>Bienvenue ")

	// now create a v2
	const v2 = `package fr.in2p3.informatique.ecole2015.web;

import org.eclipse.jetty.server.QuietServletException;

import java.io.IOException;

import javax.servlet.http.HttpServlet;
import javax.servlet.http.HttpServletRequest;
import javax.servlet.http.HttpServletResponse;

public class MyServlet extends HttpServlet
{
    private String greeting="Welcome to école informatique IN2P3 2015";
    public MyServlet(){}
    public MyServlet(String greeting)
    {
        this.greeting=greeting;
    }
    protected void doGet(HttpServletRequest request, HttpServletResponse response) throws QuietServletException, IOException
    {
        response.setContentType("text/html");
        response.setStatus(HttpServletResponse.SC_OK);
        response.getWriter().println("<h1>"+greeting+"</h1>");
        response.getWriter().println("<a href='/analyse'>Analyse de données</a>");
        response.getWriter().println("session=" + request.getSession(true).getId());
    }
}`

	myservlet, err := os.Create(filepath.Join(
		repodir,
		"src/main/java/fr/in2p3/informatique/ecole2015/web/MyServlet.java",
	))
	if err != nil {
		log.Fatalf("error opening file myservlet file: %v\n", err)
	}
	defer myservlet.Close()

	_, err = myservlet.WriteString(v2)
	if err != nil {
		log.Fatalf("error updating [%s]: %v\n", myservlet.Name(), err)
	}
	err = myservlet.Close()
	if err != nil {
		log.Fatalf("error closing [%s]: %v\n", myservlet.Name(), err)
	}

	// now create and run server-v2
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
	cmd = exec.Command(
		"docker", "run", "-d", "-p", "8082:8080",
		"--name=binet-web-app-v2",
		"binet/web-app:v2",
	)
	err = run(cmd)
	if err != nil {
		log.Fatalf("error running web-app: %v\n", err)
	}

	cmd = exec.Command(
		"docker", "tag", "-f",
		"binet/web-app:v2",
		"cc-ecole2015-docker.in2p3.fr:5000/binet/web-app:v2",
	)
	err = run(cmd)
	if err != nil {
		log.Fatalf("error tagging web-app to cc-ecole/web-app:v2: %v\n", err)
	}

	if *doPush {
		cmd = exec.Command(
			"docker", "push", "cc-ecole2015-docker.in2p3.fr:5000/binet/web-app:v2",
		)
		err = run(cmd)
		if err != nil {
			log.Fatalf("error pushing web-app:v2: %v\n", err)
		}
	}
	testWebServer("http://"+localhost+":8082/", "<h1>Welcome to ")
	testWebServer("http://"+localhost+":8080/", "<h1>Bienvenue ")

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
