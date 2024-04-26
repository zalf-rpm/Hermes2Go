package main

import (
	"flag"
	"log"

	"github.com/zalf-rpm/Hermes2Go/hermes"
)

// read a old hermes crop file and convert it to a new format
func main() {
	// flags:
	// -input file
	// -output file
	input := flag.String("input", "", "input file")
	output := flag.String("output", "", "output file")
	flag.Parse()

	if *input == "" {
		log.Fatal("input file not specified")
	}
	outName := *output
	if *output == "" {
		outName = *input + ".yml"
		log.Printf("output file not specified, using %s", outName)
	}

	// read the input file & convert the data
	cropParam, err := hermes.ConvertCropParamClassicToYml(*input)
	if err != nil {
		log.Fatal(err)
	}
	// write the output file
	err = hermes.WriteCropParam(outName, cropParam)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("wrote %s", outName)
}
