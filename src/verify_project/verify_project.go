package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/zalf-rpm/Hermes2Go/hermes"
	"gopkg.in/yaml.v2"
)

func main() {
	// cmdline args:
	// project folder
	projectDir := flag.String("projectDir", "./project", "project folder")
	project := flag.String("project", "ex2", "project")
	ext := flag.String("ext", "ex2", "ext")

	flag.Parse()

	if *projectDir == "" || *project == "" || *ext == "" {
		flag.PrintDefaults()
		return
	}

	// read config file
	configFile := filepath.Join(*projectDir, *project, "config.yml")
	hconfig := hermes.NewDefaultConfig()
	// if config files exists, read it into hconfig
	if _, err := os.Stat(configFile); err == nil {
		byteData := hermes.HermesFilePool.Get(&hermes.FileDescriptior{FilePath: configFile, ContinueOnError: true, UseFilePool: true})
		err := yaml.Unmarshal(byteData, &hconfig)
		if err != nil {
			log.Fatalf("error: %v", err)
		}
	} else {
		log.Fatalf("error: %v", err)
	}

	Datum := hermes.DateConverter(hconfig.DivideCentury, hconfig.Dateformat)
	Kalender := hermes.KalenderConverter(hconfig.Dateformat, ".")

	// verify crop rotation dates

	// open crop rotation file
	cropPath := filepath.Join(*projectDir, *project, "crop_"+*ext+".txt")
	cropFile, err := os.Open(cropPath)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer cropFile.Close()

	// read crop rotation file
	cropReader := bufio.NewScanner(cropFile)
	// skip header
	hermes.LineInut(cropReader)

	// previous line harvest
	var prevHarvest int
	// previous line rotation id
	var prevRotationID string
	lineCount := 1
	for cropReader.Scan() {
		lineCount++
		line := cropReader.Text()
		// split line
		tokens := strings.Fields(line)
		if len(tokens) > 4 {
			// Field_ID    crp  sowing harvst Rex yld autorg variety comment
			// SOYSM1    SM  05151980 09311980 080 050 0
			rotationID := tokens[0]
			_, dateSowing := Datum(tokens[2])
			_, dateHarvest := Datum(tokens[3])
			if dateSowing > dateHarvest {
				fmt.Printf("ERROR line %d: Sowing date after harvest date in crop rotation file (%s > %s)",
					lineCount, Kalender(dateSowing), Kalender(prevHarvest))
			}
			if prevRotationID != rotationID {
				prevHarvest = 0
			} else {
				if dateSowing < prevHarvest {
					fmt.Printf("ERROR line %d: Sowing date before previous harvest date in crop rotation file (%s < %s)",
						lineCount, Kalender(dateSowing), Kalender(prevHarvest))
				}
			}

		}
	}

}
