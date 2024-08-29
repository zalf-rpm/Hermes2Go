package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"time"

	"capnproto.org/go/capnp/v3/rpc"
	hermes_service_capnp "github.com/zalf-rpm/Hermes2Go/hermes_service/capnp/hermes_service_capnp"
)

// this is a simple producer-consumer example, to test the hermes_service package
// the producer part generates setups and sends it to the hermes_service
// the consumer recceives the results from the hermes_service and processes them

func main() {

	// command line flags
	// address of the hermes_service
	hServive := flag.String("hermes_service", "localhost:8841", "address of the hermes_service")
	// work directory
	workDir := flag.String("workdir", "", "working directory")
	// batch file
	batchFile := flag.String("batch", "", "batch file")

	flag.Parse()

	if *workDir == "" {
		log.Fatal("workdir not specified")
	}
	if *hServive == "" {
		log.Fatal("hermes_service not specified")
	}
	if *batchFile == "" {
		log.Fatal("batch file not specified")
	}
	setup, err := setupFromBatch(*batchFile)
	if err != nil {
		log.Fatal(err)
	}

	doneProducer := make(chan bool)
	doneConsumer := make(chan bool)
	// create a new ResultCallback
	cb := &ResultCallback{consumer: make(chan *resultData)}
	go runConsumer(cb.consumer, doneConsumer)
	go runProducer(*workDir, *hServive, cb, doneProducer, doneConsumer, setup)

	// wait for the producer and consumer to finish
	<-doneProducer
	fmt.Println("producer done")
}

// runProducer generates setups and sends them to the hermes_service
func runProducer(workDir, hService string, cb *ResultCallback, done chan<- bool, doneConsumer <-chan bool, setup *setup) {
	defer func() { done <- true }()

	conn, err := net.Dial("tcp", hService)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	// establish connection to registry
	connection := rpc.NewConn(rpc.NewPackedStreamTransport(conn), nil)
	// get ModelRunController Bootstrap
	client := connection.Bootstrap(context.Background())
	mrc := hermes_service_capnp.SessionServer(client)
	defer mrc.Release()

	// create a new session
	sessionFut, relSession := mrc.NewSession(context.Background(), func(p hermes_service_capnp.SessionServer_newSession_Params) error {
		err := p.SetWorkdir(workDir)
		if err != nil {
			return err
		}
		callback := hermes_service_capnp.Callback_ServerToClient(cb)
		err = p.SetResultCallback(callback)
		return err
	})
	defer relSession()

	if run, ok := setup.nextRun(); ok {

		_, relSend := sessionFut.Session().Send(context.Background(), func(p hermes_service_capnp.Session_send_Params) error {
			err := p.SetRunId(run.id)
			if err != nil {
				return err
			}
			paramList, err := p.NewParams(int32(len(run.params)))
			if err != nil {
				return err
			}
			for i, param := range run.params {
				err = paramList.Set(i, param)
				if err != nil {
					return err
				}
			}
			err = p.SetParams(paramList)
			return err
		})
		relSend()
	}
	<-doneConsumer // wait for consumer to finish

	// close the session
	doneFut, relDone := sessionFut.Session().Close(context.Background(), func(p hermes_service_capnp.Session_close_Params) error {
		return nil
	})
	doneFut.Struct()
	relDone()
}

// runConsumer receives the results from the hermes_service
// and processes them
func runConsumer(consumer <-chan *resultData, done chan<- bool) {

	timeout := false
	for !timeout {
		select {
		case r := <-consumer:
			if r.done {
				log.Println("run done", r.runId)
			} else {
				log.Println("run data", r.runId, r.data)
			}
		case <-time.After(60 * time.Second):
			log.Println("timeout")
			timeout = true
		}
	}
	done <- true
}

// data received from the hermes_service
type resultData struct {
	runId string
	data  string
	done  bool
}

// implement the ResultCallback interface
type ResultCallback struct {
	consumer chan *resultData
}

// SendData(context.Context, Callback_sendData) error
func (r *ResultCallback) SendData(ctx context.Context, call hermes_service_capnp.Callback_sendData) error {
	runId, err := call.Args().RunId()
	if err != nil {
		return err
	}
	data, err := call.Args().OutData()
	if err != nil {
		return err
	}
	r.consumer <- &resultData{runId: runId, data: data, done: false}

	return nil
}

// Done(context.Context, Callback_done) error
func (r *ResultCallback) Done(ctx context.Context, call hermes_service_capnp.Callback_done) error {
	runId, err := call.Args().RunId()
	if err != nil {
		return err
	}
	r.consumer <- &resultData{runId: runId, done: true}
	return nil
}
