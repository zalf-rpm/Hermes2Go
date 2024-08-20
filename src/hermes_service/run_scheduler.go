package main

// event loop for the scheduler

import (
	"fmt"

	"github.com/zalf-rpm/Hermes2Go/hermes"
)

// RunScheduler runs the scheduler
func runScheduler(newSessionChan, closedSession <-chan *Hermes_Session, hermesRun <-chan *Hermes_Run, maxConcurrent uint, writeLogoutput bool) {

	// list of all sessions
	sessions := make([]*Hermes_Session, 0)
	toDoRuns := make([]*Hermes_Run, 0)
	var activeRuns uint
	logOutputChan := make(chan string)
	resultChannel := make(chan string)
	for {
		// check if we can start a new run
		for activeRuns < maxConcurrent && len(toDoRuns) > 0 {
			run := toDoRuns[0]
			toDoRuns = toDoRuns[1:]
			if run.session.done {
				// drop left over runs if session is done
				continue
			}
			activeRuns++
			go hermes.Run(run.session.workingDir, run.args, run.runID, resultChannel, logOutputChan)
		}
		select {
		case result := <-resultChannel:
			activeRuns--
			if writeLogoutput {
				fmt.Println(result)
			}
		case log := <-logOutputChan:
			if writeLogoutput {
				fmt.Println(log)
			}
		case newSession := <-newSessionChan:
			sessions = append(sessions, newSession)
		case sessionClosed := <-closedSession:
			// TBI: stop session
			for i, session := range sessions {
				session.done = true
				// close all runs, do not send results
				if session == sessionClosed {
					sessions = append(sessions[:i], sessions[i+1:]...)
					break
				}
			}
		case run := <-hermesRun:
			toDoRuns = append(toDoRuns, run)
		}

	}

}

type Hermes_Run struct {
	session *Hermes_Session
	runID   string
	args    []string
}
