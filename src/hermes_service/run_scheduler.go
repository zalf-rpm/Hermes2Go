package main

// event loop for the scheduler

import (
	"container/list"
	"fmt"

	"github.com/zalf-rpm/Hermes2Go/hermes"
)

// RunScheduler runs the scheduler
func runScheduler(closedSession <-chan *Hermes_Session, hermesRun <-chan *Hermes_Run, maxConcurrent uint, writeLogoutput bool) {

	toDoRuns := list.New()
	var activeRuns uint
	logOutputChan := make(chan string)
	resultChannel := make(chan *hermes.RunReturn)
	for {
		// check if we can start a new run
		for activeRuns < maxConcurrent && toDoRuns.Len() > 0 {
			runEl := toDoRuns.Front()
			toDoRuns.Remove(runEl)
			run := runEl.Value.(*Hermes_Run)
			if run.session.done {
				// drop left over runs if session is done
				continue
			}
			activeRuns++
			go run.session.hermesSession.Run(run.session.workingDir, run.args, run.runID, resultChannel, logOutputChan)
		}
		select {
		case result := <-resultChannel:
			activeRuns--
			if !result.Success {
				// TODO send error to client

			}
			if writeLogoutput {
				fmt.Println(result)
			}
		case log := <-logOutputChan:
			if writeLogoutput {
				fmt.Println(log)
			}
		case session := <-closedSession:
			session.done = true
			session.callBack.Release()
			session.hermesSession.Close()

		case run := <-hermesRun:
			toDoRuns.PushBack(run)
		}

	}

}

type Hermes_Run struct {
	session *Hermes_Session
	runID   string
	args    []string
}
