package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
)

var version = "undefined"

func main() {
	argsWithoutProg := os.Args[1:]

	//calcHermesBatch -size
	//calcHermesBatch -list
	var numNodes uint64
	returnSize := false
	returnList := false
	var lines uint64
	for i, arg := range argsWithoutProg {
		if arg == "-v" {
			fmt.Println("Version: ", version)
			return
		}
		if arg == "-size" && i+1 < len(argsWithoutProg) {
			returnSize = true
			val, err := strconv.ParseUint(argsWithoutProg[i+1], 10, 64)
			if err != nil {
				log.Fatal(err)
			}
			numNodes = val
		} else if arg == "-list" {
			returnList = true
			val, err := strconv.ParseUint(argsWithoutProg[i+1], 10, 64)
			if err != nil {
				log.Fatal(err)
			}
			numNodes = val
		} else if arg == "-batch" && i+1 < len(argsWithoutProg) {
			batchFile := argsWithoutProg[i+1]
			lines = readProj(batchFile)
		}
	}
	if returnSize {
		if lines/numNodes == 0 {
			fmt.Print(lines)
		} else {
			fmt.Print(numNodes)
		}
		return
	}
	if returnList {
		//result := "("
		var result string
		if lines/numNodes == 0 {
			var i uint64 = 1
			for ; i < lines; i++ {
				if i == lines-1 {
					result += fmt.Sprintf("%d-%d", i, i)
				} else {
					result += fmt.Sprintf("%d-%d ", i, i)
				}
			}
		} else {
			sizePerSlice := lines / numNodes
			rest := lines % numNodes
			var i uint64 = 1
			var lastSlice uint64
			for ; i <= numNodes; i++ {
				strFormat := "%d-%d "
				if i == numNodes {
					strFormat = "%d-%d"
				}
				if i <= rest {
					result += fmt.Sprintf(strFormat, lastSlice+1, lastSlice+sizePerSlice+1)
					lastSlice = lastSlice + sizePerSlice + 1
				} else {
					result += fmt.Sprintf(strFormat, lastSlice+1, lastSlice+sizePerSlice)
					lastSlice = lastSlice + sizePerSlice
				}
			}
			fmt.Print(result)
		}
	}
}

func readProj(filename string) uint64 {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	count, err := lineCounter(file)
	if err != nil {
		log.Fatal(err)
	}
	return count
}

func lineCounter(r io.Reader) (uint64, error) {
	buf := make([]byte, 32*1024)
	var count uint64
	lineSep := []byte{'\n'}

	distanceCarryForward := 0
	prevNotCarageReturn := true
	for {
		c, err := r.Read(buf)
		//count += bytes.Count(buf[:c], lineSep)

		if c > 0 {
			index := -1
			startIndex := 0
			for ok := true; ok; ok = index != -1 && startIndex < c {
				distance := 0
				index = bytes.Index(buf[startIndex:c], lineSep)
				if index != -1 {
					distance = distanceCarryForward + index + 1
					distanceCarryForward = 0
					if index > 0 {
						prevNotCarageReturn = buf[index-1] != '\r'
					}
					if (distance > 1 && prevNotCarageReturn) || (distance > 2 && !prevNotCarageReturn) {
						count++
					}
				} else {
					distanceCarryForward = len(buf[startIndex:c])
					prevNotCarageReturn = buf[c-1] != '\r'
				}
				startIndex = startIndex + index + 1
			}
		}
		switch {
		case err == io.EOF:
			if distanceCarryForward > 0 {
				count++
			}
			return count, nil

		case err != nil:
			return count, err
		}
	}
}
