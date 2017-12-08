// Copyright 2017 The SoLid read-out software Authors.  All rights reserved.

package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"golang.org/x/net/websocket"
)

var (
	errInvalidMethod = errors.New("invalid HTTP method")
)

func main() {
	srv := newServer(*addrFlag) // HL
	mux := http.NewServeMux()
	mux.Handle("/", srv)

	// data/cmd channels (as websockets)
	mux.Handle("/cmds", websocket.Handler(srv.handleCmds))
	mux.Handle("/msgs", websocket.Handler(srv.handleMsgs))

	// authentication
	//mux.HandleFunc("/login", srv.handleLogin)
	//mux.HandleFunc("/logout", srv.handleLogout)

	mux.HandleFunc("/start", srv.wrap(srv.handleStart))
	mux.HandleFunc("/stop", srv.wrap(srv.handleStop))
	mux.HandleFunc("/status", srv.wrap(srv.handleStatus))
	mux.HandleFunc("/configure", srv.wrap(srv.handleConfigure))

	log.Printf("starting SoLiD server on [%s]...", srv.addr)
	err := http.ListenAndServe(srv.addr, mux)
	if err != nil {
		log.Fatal(err)
	}
}

type server struct {
	addr string
	cmds chan cmdRequest
	msgs chan msgData

	cmdsReg registry // clients interested in sending/receiving RunCtl commands
	msgsReg registry // clients interested in receiving RunCtl messages

	rc RunControler // handle to the actual run control
}

func newServer(addr string) *server {
	if addr == "" {
		addr = getHostIP() + ":8080"
	}

	srv := &server{
		addr:    addr,
		cmds:    make(chan cmdRequest),
		msgs:    make(chan msgData),
		cmdsReg: newRegistry(),
		msgsReg: newRegistry(),
		rc:      newRunCtl(),
	}

	go srv.mon()
	go srv.run() // HL

	return srv
}

func (srv *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	srv.wrap(func(w http.ResponseWriter, r *http.Request) error {
		_, err := fmt.Fprintf(w, tmplMain)
		return err
	})(w, r)
}

func (srv *server) wrap(f func(w http.ResponseWriter, r *http.Request) error) http.HandlerFunc {
	handler := func(w http.ResponseWriter, r *http.Request) {
		err := f(w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	return handler
}

func (srv *server) handleStart(w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodPost {
		return errInvalidMethod
	}
	log.Printf("received req: /start\n")
	cmd := cmdRequest{Cmd: CmdStart, Reply: make(chan cmdReply)} // HL
	srv.cmds <- cmd                                              // HL
	reply := <-cmd.Reply                                         // HL
	if reply.Err != "" {
		return fmt.Errorf(reply.Err)
	}
	return nil
}

func (srv *server) handleStop(w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodPost {
		return errInvalidMethod
	}
	log.Printf("received req: /stop\n")
	cmd := cmdRequest{Cmd: CmdStop, Reply: make(chan cmdReply)}
	srv.cmds <- cmd
	reply := <-cmd.Reply
	if reply.Err != "" {
		return fmt.Errorf(reply.Err)
	}
	return nil
}

func (srv *server) handleStatus(w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodGet {
		return errInvalidMethod
	}
	log.Printf("received req: /status\n")
	cmd := cmdRequest{Cmd: CmdStatus, Reply: make(chan cmdReply)}
	srv.cmds <- cmd
	reply := <-cmd.Reply
	if reply.Err != "" {
		return fmt.Errorf(reply.Err)
	}
	return nil
}

func (srv *server) handleConfigure(w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodPost {
		return errInvalidMethod
	}
	log.Printf("received req: /configure\n")
	cmd := cmdRequest{Cmd: CmdConfigure, Reply: make(chan cmdReply)}
	srv.cmds <- cmd
	reply := <-cmd.Reply
	if reply.Err != "" {
		return fmt.Errorf(reply.Err)
	}
	return nil
}

func (srv *server) handleCmds(ws *websocket.Conn) {
	c := &client{
		srv:   srv,
		reg:   &srv.cmdsReg,
		datac: make(chan msgData, 256),
		ws:    ws,
	}
	c.reg.register <- c
	defer c.Release()

cmdLoop:
	for {
		log.Printf("waiting for commands from client %v...\n", c.ws.RemoteAddr())
		cmd := cmdRequest{Reply: make(chan cmdReply)}
		err := websocket.JSON.Receive(ws, &cmd)
		if err != nil {
			log.Printf("error rcv: %v\n", err)
			if err == io.EOF {
				break cmdLoop
			}
			continue
		}
		srv.cmds <- cmd
		reply := <-cmd.Reply
		err = websocket.JSON.Send(ws, reply)
		if err != nil {
			log.Printf("error sending reply: %v\n", err)
			if err == io.EOF {
				break cmdLoop
			}
			continue
		}
	}
}

func (srv *server) handleMsgs(ws *websocket.Conn) {
	c := &client{
		srv:   srv,
		reg:   &srv.msgsReg,
		datac: make(chan msgData, 256),
		ws:    ws,
	}
	c.reg.register <- c
	defer c.Release()

msgLoop:
	for {
		select {
		case msg := <-c.datac:
			log.Printf("received msg: %v\n", msg)
			err := websocket.JSON.Send(ws, msg)
			if err != nil {
				log.Printf("error sending msg: %v\n", err)
				if err == io.EOF {
					break msgLoop
				}
				continue
			}
		}
	}
}

func (srv *server) mon() {
	lines := make([]string, 0, 32)
	for {
		select {
		case c := <-srv.msgsReg.register:
			srv.msgsReg.clients[c] = true
		case c := <-srv.msgsReg.unregister:
			_, ok := srv.msgsReg.clients[c]
			if !ok {
				continue
			}
			delete(srv.msgsReg.clients, c)
			close(c.datac)
			log.Printf("client disconnected [%v]\n", c.ws.RemoteAddr())
		case msg := <-srv.msgs:
			if n := len(lines); n == cap(lines) {
				copy(lines[2:n-1], lines[3:n])
				lines = lines[:n-1]
				lines[1] = "[...]"
			}
			lines = append(lines, msg.Lines)

			if len(srv.msgsReg.clients) == 0 {
				// no client connected.
				continue
			}
			msg.Lines = strings.Join(lines, "<br>")
			for c := range srv.msgsReg.clients {
				select {
				case c.datac <- msg:
				default:
					close(c.datac)
					delete(srv.msgsReg.clients, c)
				}
			}
		}
	}
}

func (srv *server) run() {
	for {
		select {
		case cmd := <-srv.cmds: // HL
			log.Printf("received cmd %q\n", cmd.Cmd)
			srv.handleCmd(context.Background(), cmd) // HL
		case c := <-srv.cmdsReg.register:
			srv.cmdsReg.clients[c] = true
		case c := <-srv.cmdsReg.unregister:
			_, ok := srv.cmdsReg.clients[c]
			if !ok {
				continue
			}
			delete(srv.cmdsReg.clients, c)
			close(c.datac)
			log.Printf("client disconnected [%v]\n", c.ws.RemoteAddr())
		}
	}
}

func (srv *server) handleCmd(ctx context.Context, cmd cmdRequest) {
	var err error
	reply := cmdReply{Req: cmd.Cmd}
	srv.sendMsg(fmt.Sprintf("==> %q at: %v", cmd.Cmd, time.Now().Format(time.RFC3339)))

	switch cmd.Cmd {
	case CmdStart:
		err = srv.rc.start(ctx)
	case CmdStop:
		err = srv.rc.stop(ctx)
	case CmdConfigure:
		err = srv.rc.configure(ctx)
	case CmdStatus:
		err = srv.rc.status(ctx)
	default:
		err = fmt.Errorf("invalid command (%v)", cmd.Cmd)
	}
	if err != nil {
		reply.Err = err.Error()
		srv.sendMsg(fmt.Sprintf("=== err=%v", reply.Err))
	}
	srv.sendMsg(fmt.Sprintf("<== %q at: %v", cmd.Cmd, time.Now().Format(time.RFC3339)))
	cmd.Reply <- reply
}

func (srv *server) sendMsg(msgs ...string) {
	go func(msgs []string) {
		for _, msg := range msgs {
			srv.msgs <- msgData{Lines: msg}
		}
	}(msgs)
}

func getHostIP() string {
	host, err := os.Hostname()
	if err != nil {
		log.Fatalf("could not retrieve hostname: %v\n", err)
	}

	addrs, err := net.LookupIP(host)
	if err != nil {
		log.Fatalf("could not lookup hostname IP: %v\n", err)
	}

	for _, addr := range addrs {
		ipv4 := addr.To4()
		if ipv4 == nil {
			continue
		}
		return ipv4.String()
	}

	log.Fatalf("could not infer host IP")
	return ""
}

// START CMD OMIT
const (
	CmdStart     = "start"
	CmdStop      = "stop"
	CmdConfigure = "configure"
	CmdStatus    = "status"
)

type cmdRequest struct {
	Cmd   string `json:"cmd"`
	Reply chan cmdReply
}

// END CMD OMIT

type cmdReply struct {
	Req string `json:"req"`
	Err string `json:"err"`
}

type registry struct {
	clients    map[*client]bool
	register   chan *client
	unregister chan *client
}

func newRegistry() registry {
	return registry{
		clients:    make(map[*client]bool),
		register:   make(chan *client),
		unregister: make(chan *client),
	}
}

type client struct {
	srv   *server
	reg   *registry
	ws    *websocket.Conn
	datac chan msgData
	acl   byte // acl notes whether the client is authentified and has r/w access
}

func (c *client) Release() {
	c.reg.unregister <- c
	c.ws.Close()
	c.reg = nil
	c.srv = nil
}

func (c *client) run() {
	//c.ws.SetReadLimit(maxMessageSize)
	//c.ws.SetReadDeadline(time.Now().Add(pongWait))
	//c.ws.SetPongHandler(func(string) error { c.ws.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for data := range c.datac {
		err := websocket.JSON.Send(c.ws, data)
		if err != nil {
			log.Printf(
				"error sending data to [%v]: %v\n",
				c.ws.LocalAddr(),
				err,
			)
			break
		}
	}
}

type msgData struct {
	Lines string `json:"lines"`
}

const tmplMain = `
<html>
	<head>
		<title>SoLiD Run Control server</title>
		<meta name="viewport" content="width=device-width, initial-scale=1">
		<link rel="stylesheet" href="https://www.w3schools.com/w3css/3/w3.css">

		<script type="text/javascript">
		var sockCmds = null;
		var sockMsgs = null;

		function updateCmd(data) {
		};

		function updateMsg(data) {
			var node = document.getElementById("rc-log");
			node.innerHTML = data["lines"];
			node.scrollTop = node.scrollHeight;
		};

		window.onload = function() {
			sockCmds = new WebSocket("ws://"+location.host+"/cmds");

			sockCmds.onmessage = function(event) {
				var data = JSON.parse(event.data);
				console.log("CMD: received JSON: "+JSON.stringify(data));
				updateCmd(data);
			};

			sockMsgs = new WebSocket("ws://"+location.host+"/msgs");
			sockMsgs.onmessage = function(event) {
				var data = JSON.parse(event.data);
				console.log("MSG: received JSON: "+JSON.stringify(data));
				updateMsg(data);
			};

		};

		function rcStartRun() {
			console.log("sending start-run...");
			sockCmds.send(JSON.stringify({"cmd": "start"}));
		};

		function rcStopRun() {
			console.log("sending stop-run...");
			sockCmds.send(JSON.stringify({"cmd": "stop"}));
		};

		function rcConfigureRun() {
			console.log("sending configure-run...");
			sockCmds.send(JSON.stringify({"cmd": "configure"}));
		};

		function rcStatusRun() {
			console.log("sending status-run...");
			sockCmds.send(JSON.stringify({"cmd": "status"}));
		};

		</script>
	</head>

	<body>
		<header class="w3-container w3-black">
		    <h1>Run Control panel</h1>
			<div class="w3-container">
			<!--
				<button class="w3-button w3-blue" onclick="run();">Launch</button>
			-->
			</div>
		</header>
		<div id="rc-logo"><img src="/favicon.ico" style="float:right; vertical-align: bottom; height: 150px;"></img></div>

		<button class="w3-button w3-blue" onclick="rcStartRun();">Start</button>
		<button class="w3-button w3-blue" onclick="rcStopRun();">Stop</button>
		<button class="w3-button w3-blue" onclick="rcConfigureRun();">Configure</button>
		<button class="w3-button w3-blue" onclick="rcStatusRun();">Status</button>

		<div class="w3-cell-row">
		  <div class="w3-container w3-light-gray w3-cell w3-mobile">
			<div class="w3-cell-row">
				<div id="rc-log"></div>
				<div class="w3-cell-row w3-display-center w3-panel w3-container">
					<div class="w3-container w-cell w3-center w3-display-center">
						<div id="rc-fsm" class="w3-container w3-card-4 w3-center w3-white"></div>
					</div>
				</div>
			</div>
		  </div>
		</div>

	</body>
</html>
`
