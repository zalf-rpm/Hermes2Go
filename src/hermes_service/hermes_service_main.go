package main

import (
	"bufio"
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

func main() {

	// cmd line arguments:
	// number of concurrent operations
	concurrentOperations := flag.Uint("concurrent", 10, "number of concurrent operations")
	// print version
	printVersion := flag.Bool("version", false, "print version")
	writeLogoutput := flag.Bool("log", false, "write log output")

	flag.Parse()

	if *printVersion {
		fmt.Println("Version: ", version)
		return
	}
	// listen on a socket
	l, err := net.Listen("tcp", "localhost:8841")
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

	}()
	closedSession := make(chan *Hermes_Session)
	hermesRun := make(chan *Hermes_Run)
	go runScheduler(closedSession, hermesRun, *concurrentOperations, *writeLogoutput)

	for {
		// accept connections and serve
		c, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("server: accepted connection from", c.RemoteAddr())
		go func() {
			server := Hermes_SessionServer{
				writeLogoutput: *writeLogoutput,
				sessions:       []*Hermes_Session{},
				runChan:        hermesRun,
				closeChan:      closedSession,
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
					for _, session := range server.sessions {
						session.closedSession <- session
					}
					return
				case err := <-errorChan:
					fmt.Println("Error reported:", err)
					for _, session := range server.sessions {
						session.closedSession <- session
					}
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
	writeLogoutput bool
	sessions       []*Hermes_Session
	runChan        chan<- *Hermes_Run
	closeChan      chan<- *Hermes_Session
}

func (a *Hermes_SessionServer) NewSession(ctx context.Context, call hermes_service_capnp.SessionServer_newSession) error {
	workdir, err := call.Args().Workdir()
	if err != nil {
		return err
	}
	if a.writeLogoutput {
		fmt.Println("server: NewSession Received for WORKDIR: ", workdir)
	}
	callback := call.Args().ResultCallback().AddRef()
	// create a new session
	session := &Hermes_Session{
		workingDir:    workdir,
		hermesRun:     a.runChan,
		done:          false,
		callBack:      callback,
		closedSession: a.closeChan,
		hermesSession: hermes.NewHermesSession(),
	}
	session.hermesSession.HermesOutWriter = NewOutWriterCallback(callback)
	// return the session
	results, err := call.AllocResults()
	if err != nil {
		return err
	}
	err = results.SetSession(hermes_service_capnp.Session_ServerToClient(session))
	if err != nil {
		return err
	}
	a.sessions = append(a.sessions, session)

	return nil
}

// implement the interface for the capnp schema
// Session_Server
//	Send(context.Context, Session_send) error
//	Close(context.Context, Session_close) error

type Hermes_Session struct {
	workingDir    string
	hermesRun     chan<- *Hermes_Run
	done          bool
	callBack      hermes_service_capnp.Callback
	closedSession chan<- *Hermes_Session
	hermesSession *hermes.HermesSession
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

	// send the run to the scheduler
	a.hermesRun <- &Hermes_Run{
		session: a,
		runID:   runId,
		args:    paramList,
	}

	return nil
}

func (a *Hermes_Session) Close(ctx context.Context, call hermes_service_capnp.Session_close) error {
	fmt.Println("session: Close Received")
	a.closedSession <- a
	// close all runs, do not send results
	return nil
}
func NewOutWriterCallback(c hermes_service_capnp.Callback) func(string, bool) (hermes.OutWriter, error) {
	return func(filename string, append bool) (hermes.OutWriter, error) {
		callbackWriter := &CallbackWriter{
			id:       filename,
			callback: c,
		}
		fwriter := bufio.NewWriter(callbackWriter)
		return &CallBackOutwriter{cWriter: callbackWriter, fwriter: fwriter}, nil
	}
}

// implements io.Writer
type CallbackWriter struct {
	id       string
	callback hermes_service_capnp.Callback
}

func (c *CallbackWriter) Write(data []byte) (n int, err error) {

	future, rel := c.callback.SendData(context.Background(), func(p hermes_service_capnp.Callback_sendData_Params) error {
		err := p.SetRunId(c.id)
		if err != nil {
			return err
		}
		err = p.SetOutData(string(data))
		if err != nil {
			return err
		}
		return nil
	})
	defer rel()
	_, err = future.Struct()
	if err != nil {
		return 0, err
	}

	return len(data), nil
}

// implement interface hermes.OutWriter
type CallBackOutwriter struct {
	cWriter *CallbackWriter
	fwriter *bufio.Writer
}

func (c *CallBackOutwriter) Write(s string) (int, error) {
	return c.fwriter.WriteString(s)
}
func (c *CallBackOutwriter) WriteBytes(b []byte) (int, error) {
	return c.fwriter.Write(b)
}
func (c *CallBackOutwriter) WriteRune(r rune) (int, error) {
	return c.fwriter.WriteRune(r)
}
func (c *CallBackOutwriter) WriteError(errOut error) (int, error) {
	_, _ = c.cWriter.callback.SendError(context.Background(), func(p hermes_service_capnp.Callback_sendError_Params) error {
		err := p.SetRunId(c.cWriter.id)
		if err != nil {
			return err
		}
		err = p.SetError(errOut.Error())
		if err != nil {
			return err
		}
		return nil
	})
	return 0, nil
}

func (c *CallBackOutwriter) Close() {
	err := c.fwriter.Flush()
	if err != nil {
		log.Fatalln(err)
	}
	// send done
	future, rel := c.cWriter.callback.Done(context.Background(), func(p hermes_service_capnp.Callback_done_Params) error {
		err := p.SetRunId(c.cWriter.id)
		if err != nil {
			return err
		}
		return nil
	})
	defer rel()
	_, err = future.Struct()
	if err != nil {
		log.Fatalln(err)
	}
}

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
	cerr.Out <- fmt.Errorf(message, args...)
}
