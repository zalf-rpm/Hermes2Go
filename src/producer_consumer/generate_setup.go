package main

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

type setup struct {
	file           *os.File
	scanner        *bufio.Scanner
	currentLineIdx uint
}

func setupFromBatch(batchPath string) (*setup, error) {
	f, err := os.Open(batchPath)
	if err != nil {
		return nil, err
	}
	return &setup{file: f, scanner: bufio.NewScanner(f)}, nil
}

type runList struct {
	id     string
	params []string
}

func (s *setup) nextRun() (*runList, bool) {
	// read the next line from the file
	// and return it as a runList
	if !s.scanner.Scan() {
		return nil, false
	}

	line := s.scanner.Text()
	// extract the id and params from the line
	tokens := strings.Fields(line)
	polygonID := ""
	plotNr := ""
	//plotNr=10002 poligonID=30413
	for _, t := range tokens {
		if strings.HasPrefix(t, "poligonID=") {
			polygonID = t[10:]

		} else if strings.HasPrefix(t, "plotNr=") {
			plotNr = t[7:]
		}
	}
	id := polygonID + plotNr
	if id == "" {
		id = strconv.Itoa(int(s.currentLineIdx))
	}

	s.currentLineIdx++
	return &runList{
		id:     id,
		params: tokens,
	}, true
}

func (s *setup) close() {
	s.file.Close()
}
