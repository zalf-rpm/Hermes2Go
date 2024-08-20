package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"capnproto.org/go/capnp/v3"
	"capnproto.org/go/capnp/v3/rpc"
	"github.com/zalf-rpm/Hermes2Go/hermes"
	hermes_service_capnp "github.com/zalf-rpm/Hermes2Go/hermes_service/capnp/hermes_service_capnp"
)

var version = "undefined"
var concurrentOperations uint = 10

func main() {

	// cmd line arguments:
	// working directory for project setups
	workingDir := flag.String("workingdir", "", "working directory for project setups")
	// number of concurrent operations
	concurrentOperations = *flag.Uint("concurrent", 10, "number of concurrent operations")
	// print version
	printVersion := flag.Bool("version", false, "print version")
	writeLogoutput := *flag.Bool("log", false, "write log output")

	flag.Parse()

	if *printVersion {
		fmt.Println("Version: ", version)
		return
	}
	// listen on a socket
	l, err := net.Listen("tcp", "localhost:1234")
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()

	// catch signals to close the listener
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		l.Close()
		hermes.HermesFilePool.Close()
	}()
	sessionChan := make(chan *Hermes_Session)
	closedSession := make(chan *Hermes_Session)
	hermesRun := make(chan *Hermes_Run)
	go runScheduler(sessionChan, closedSession, hermesRun, concurrentOperations, writeLogoutput)

	for {
		// accept connections and serve
		c, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("server: accepted connection from", c.RemoteAddr())
		go func() {
			server := Hermes_SessionServer{
				workingDir:     *workingDir,
				writeLogoutput: writeLogoutput,
			}
			client := hermes_service_capnp.SessionServer_ServerToClient(&server)
			errorChan := make(chan error)
			msgChan := make(chan string)
			conn := rpc.NewConn(rpc.NewPackedStreamTransport(c), &rpc.Options{BootstrapClient: capnp.Client(client), Logger: &ConnError{Out: errorChan, Msg: msgChan}})

			defer conn.Close()
			fmt.Println("Bootstraping" + c.RemoteAddr().String())
			for {
				select {
				case <-conn.Done():
					fmt.Println("Connection closed")
					return
				case err := <-errorChan:
					fmt.Println("Error reported:", err)
					return
				case msg := <-msgChan:
					fmt.Println("Message reported:", msg)
				}
			}

		}()
	}
}

// SessionServer_Server interface
// 	NewSession(context.Context, Server_newSession) error

// Hermes_Server implements the interface for the capnp schema SessionServer_Server
type Hermes_SessionServer struct {
	workingDir     string
	writeLogoutput bool
	sessionChan    chan<- *Hermes_Session
}

func (a *Hermes_SessionServer) NewSession(ctx context.Context, call hermes_service_capnp.SessionServer_newSession) error {
	env, err := call.Args().Env()
	if err != nil {
		return err
	}
	if a.writeLogoutput {
		fmt.Println("server: NewSession Received", env)
	}
	// create a new session
	session := &Hermes_Session{
		workingDir:   a.workingDir,
		runParams:    make(map[string][]string),
		runCallbacks: make(map[string]hermes_service_capnp.Callback),
	}
	// send the session to the scheduler
	a.sessionChan <- session
	// return the session
	results, err := call.AllocResults()
	if err != nil {
		return err
	}
	err = results.SetSession(hermes_service_capnp.Session_ServerToClient(session))
	if err != nil {
		return err
	}

	return nil
}

// implement the interface for the capnp schema
// Session_Server
//	Send(context.Context, Session_send) error
//	Close(context.Context, Session_close) error

type Hermes_Session struct {
	workingDir    string
	hermesRun     chan<- *Hermes_Run
	runParams     map[string][]string
	runCallbacks  map[string]hermes_service_capnp.Callback
	done          bool
	closedSession chan<- *Hermes_Session
}

func (a *Hermes_Session) Send(ctx context.Context, call hermes_service_capnp.Session_send) error {
	runId, err := call.Args().RunId()
	if err != nil {
		return err
	}
	params, err := call.Args().Params()
	if err != nil {
		return err
	}
	if params.Len() == 0 {
		fmt.Println("server: Received empty params, run default")
		// run default
		a.runParams[runId] = []string{}
	} else {
		paramList := make([]string, params.Len())
		// print the params
		fmt.Println("server: Received Params:")
		for i := 0; i < params.Len(); i++ {
			param, err := params.At(i)
			if err != nil {
				return err
			}
			fmt.Println(param)
			paramList[i] = param
		}
		a.runParams[runId] = paramList
	}
	a.runCallbacks[runId] = call.Args().ResultCallback()
	// send the run to the scheduler
	a.hermesRun <- &Hermes_Run{
		session: a,
		runID:   runId,
		args:    a.runParams[runId],
	}

	return nil
}

func (a *Hermes_Session) Close(ctx context.Context, call hermes_service_capnp.Session_close) error {
	fmt.Println("server: Close Received")
	a.closedSession <- a
	// close all runs, do not send results
	return nil
}

// type Callback_Server interface {
// 	SendHeader(context.Context, Callback_sendHeader) error

// 	SendResult(context.Context, Callback_sendResult) error

// 	Done(context.Context, Callback_done) error
// }

type ConnError struct {
	Out chan<- error
	Msg chan<- string
}

// Logger interface
func (cerr *ConnError) Debug(message string, args ...any) {
	cerr.Msg <- message
}
func (cerr *ConnError) Info(message string, args ...any) {
	cerr.Msg <- message
}
func (cerr *ConnError) Warn(message string, args ...any) {
	cerr.Msg <- message
}

func (cerr *ConnError) Error(message string, args ...any) {
	cerr.Out <- fmt.Errorf(message)
}
