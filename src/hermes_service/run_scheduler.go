package main

// event loop for the scheduler

import (
	"fmt"
	"strings"

	"github.com/zalf-rpm/Hermes2Go/hermes"
)

// RunScheduler runs the scheduler
func runScheduler(newSessionChan, closedSession <-chan *Hermes_Session, maxConcurrent uint) {

	// list of all sessions
	sessions := make([]*Hermes_Session, 0)

	for {
		select {
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
		}
	}

}

func doConcurrentBatchRun(workingDir string, writeLogoutput bool, configLines []string) {
	logOutputChan := make(chan string)
	resultChannel := make(chan string)
	var activeRuns uint
	errorSummary := checkResultForError()
	var errorSummaryResult []string
	for i, line := range configLines {
		for activeRuns == concurrentOperations {
			select {
			case result := <-resultChannel:
				activeRuns--
				errorSummaryResult = errorSummary(result)
			case log := <-logOutputChan:
				if writeLogoutput {
					fmt.Println(log)
				}
			}
		}

		if activeRuns < concurrentOperations {
			activeRuns++
			logID := fmt.Sprintf("[%v]", i)
			if writeLogoutput {
				fmt.Println(logID)
			}
			args := strings.Fields(line)
			go hermes.Run(workingDir, args, logID, resultChannel, logOutputChan)
		}
	}
	// fetch output of last runs
	for activeRuns > 0 {
		select {
		case result := <-resultChannel:
			activeRuns--
			errorSummaryResult = errorSummary(result)
		case log := <-logOutputChan:
			if writeLogoutput {
				fmt.Println(log)
			}
		}
	}
	var numErr int
	for _, line := range errorSummaryResult {
		fmt.Println(line)
		numErr++
	}

	fmt.Printf("Number of errors: %v \n", numErr-1)
}

// checkResultForError concurrent output for error/ success, and add it to a summary
func checkResultForError() func(string) []string {
	var errSummary = []string{"Error Summary:"}
	return func(result string) []string {
		if !strings.HasSuffix(result, "Success") {
			errSummary = append(errSummary, result)
		}
		return errSummary
	}
}
