package main

import (
	"flag"
	"fmt"
)

func main() {

	inPtr := flag.String("in", "SEH-weather-new.csv", "input file name")
	outPtr := flag.String("out", "SEH-weather-new.w6d", "out file name")

	flag.Parse()

	err := ConvertFileMonicaToCZ(*inPtr, *outPtr)
	if err != nil {
		fmt.Println(err)
	}
}
