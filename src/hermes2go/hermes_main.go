package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/zalf-rpm/Hermes2Go/hermes"
)

// allowed number of concurrent Operations (should be number of processor units)
var concurrentOperations uint16 = 10
var version = "undefined"

func main() {
	start := time.Now()
	var configLines []string   // lines from a batch file
	var workingDir string      // changed working directory
	var otherArgs []string     // direct arguments that are not specified in main
	var endLine = -1           // number of lines that will we executed from batch file
	var startLine = 0          // start index
	var writeLogoutput = false // write debug output
	module := "single"         // single or batch mode
	var locID string           // optional locID for single mode

	argsWithoutProg := os.Args[1:]

	// hermes session
	session := hermes.NewHermesSession()

	for i := 0; i < len(argsWithoutProg); i++ {
		arg := argsWithoutProg[i]
		// switch between single(simplace version) and batch (long version) e.g. "-module batch"
		if arg == "-module" && i+1 < len(argsWithoutProg) {
			module = argsWithoutProg[i+1]
			i++
			// batch file
		} else if arg == "-batch" && i+1 < len(argsWithoutProg) {
			// read batch file
			setupFilename := argsWithoutProg[i+1]
			if strings.HasPrefix(setupFilename, "~") {
				usr, _ := user.Current()
				dir := usr.HomeDir
				setupFilename = strings.TrimPrefix(setupFilename, "~")
				setupFilename = filepath.Join(dir, setupFilename)
			}
			absBatchFile, err := filepath.Abs(setupFilename)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("Batch File: %s \n", absBatchFile)
			file, err := os.Open(absBatchFile)
			if err != nil {
				log.Fatal(err)
			}
			defer file.Close()
			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				line := scanner.Text()
				if len(line) > 0 {
					configLines = append(configLines, line)
				}
			}
			if err := scanner.Err(); err != nil {
				log.Fatal(err)
			}
			// if no working dir was set
			if len(workingDir) == 0 {
				path, _ := filepath.Abs(file.Name())
				workingDir = filepath.Dir(path)
			}
			i++
		} else if arg == "-workingdir" && i+1 < len(argsWithoutProg) {
			dir := argsWithoutProg[i+1]
			if strings.HasPrefix(dir, "~") {
				usr, _ := user.Current()
				homedir := usr.HomeDir
				dir = strings.TrimPrefix(dir, "~")
				dir = filepath.Join(homedir, dir)
			}
			workingDir = dir
			i++
			// number of concurrent operations
		} else if arg == "-concurrent" && i+1 < len(argsWithoutProg) {
			cOps, err := strconv.ParseUint(argsWithoutProg[i+1], 10, 64)
			if err != nil {
				log.Fatal("ERROR: Failed to parse number of concurrent runs")
				return
			}
			concurrentOperations = uint16(cOps)
			i++
			// number of lines executed 11-21 (from index 10 until index 20) or 10 (first 10) or 11-end (from index 10 until the end)
		} else if arg == "-lines" && i+1 < len(argsWithoutProg) {
			splitstr := hermes.Explode(argsWithoutProg[i+1], []rune{'-'})
			if len(splitstr) == 2 {
				firstLine, err := strconv.ParseUint(splitstr[0], 10, 64)
				if err != nil {
					log.Fatal("ERROR: Failed to parse first number in -lines")
					return
				}
				if splitstr[1] != "end" {
					lastLine, err := strconv.ParseUint(splitstr[1], 10, 64)
					if err != nil {
						log.Fatal("ERROR: Failed to parse second number in -lines")
						return
					}
					if firstLine > lastLine {
						log.Fatal("ERROR: first number in -lines must be smaller equal then second number")
						return
					}
					endLine = int(lastLine)
				}

				startLine = int(firstLine) - 1
			} else {
				// evaluate the first n lines only (optional)
				numLines, err := strconv.ParseInt(argsWithoutProg[i+1], 10, 64)
				if err != nil {
					log.Fatal("ERROR: Failed to parse number of lines")
					return
				}
				endLine = int(numLines)
			}

			i++
			// write debug output
		} else if arg == "-logoutput" {
			writeLogoutput = true
		} else if arg == "-locid" && i+1 < len(argsWithoutProg) {
			locID = argsWithoutProg[i+1]
		} else if arg == "-v" || arg == "-version" {
			fmt.Println("Version: ", version)
			return
		} else if arg == "-rpc" && i+1 < len(argsWithoutProg) {
			address := argsWithoutProg[i+1]
			rpcService, err := hermes.NewRPCService(address)
			if err == nil {
				session.HermesRPCService = rpcService
			}
		} else {
			otherArgs = append(otherArgs, arg)
		}
	}

	if module == "single" {
		root := hermes.AskDirectory()
		file, scanner, err := session.Open(&hermes.FileDescriptior{FilePath: root + "/project/" + hermes.Modfil, FileDescription: "modinp"})
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
		for scanner.Scan() {
			textToken := strings.Fields(scanner.Text())
			if len(textToken) != 2 {
				continue
			}

			locid := textToken[0] // Name of location directory (character till space)
			snam := textToken[1]
			// check if location id is set, and run only for this location
			if locID == locid || locID == "" {
				singleArgs := []string{
					fmt.Sprintf("project=%s", locid),
					fmt.Sprintf("plotNr=%s", snam),
				}
				session.Run(workingDir, singleArgs, "1", nil, nil)
			}

		}
	} else if module == "batch" {
		if len(configLines) > 0 {
			doConcurrentBatchRun(session, workingDir, startLine, endLine, writeLogoutput, configLines)
		} else {
			session.Run(workingDir, otherArgs, "1", nil, nil)
		}
	} else {
		log.Fatalf("module type '%s' not recognized", module)
	}

	session.Close()
	end := time.Now()
	elapsed := end.Sub(start)
	fmt.Println("Execution time: ", elapsed)
}

func doConcurrentBatchRun(session *hermes.HermesSession, workingDir string, startLine, numberOfLines int, writeLogoutput bool, configLines []string) {
	fmt.Printf("Working Dir: %s \n", workingDir)
	fmt.Printf("Start Line: %d \n", startLine)
	fmt.Printf("End Line: %d \n", numberOfLines)

	logOutputChan := make(chan string)
	resultChannel := make(chan string)
	var activeRuns uint16
	errorSummary := checkResultForError()
	var errorSummaryResult []string
	for i, line := range configLines {
		if i < startLine {
			continue
		}
		if numberOfLines > 0 && i >= numberOfLines {
			// if number of lines is set and limit is reached
			break
		}
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
			go session.Run(workingDir, args, logID, resultChannel, logOutputChan)
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
