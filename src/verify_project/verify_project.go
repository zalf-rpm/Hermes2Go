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
	ext := flag.String("ext", "txt", "extension")
	repair := flag.Bool("repair", false, "try to repair error")

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
	KalenderOut := hermes.KalenderConverter(hconfig.Dateformat, "")

	// verify crop rotation dates

	// open crop rotation file
	cropPath := filepath.Join(*projectDir, *project, "crop_"+*project+"."+*ext)
	cropFile, err := os.Open(cropPath)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer cropFile.Close()

	var cropFileOut *os.File
	if *repair {
		// open output file
		cropPathOut := filepath.Join(*projectDir, *project, "crop_"+*project+"_out."+*ext)
		cropFileOut, err = os.Create(cropPathOut)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer cropFileOut.Close()
	}

	// read crop rotation file
	cropReader := bufio.NewScanner(cropFile)
	// header
	header := hermes.LineInut(cropReader)
	if *repair {
		cropFileOut.WriteString(header + "\n")
	}

	// previous line harvest
	var prevHarvest int
	// previous line rotation id
	var prevRotationID string
	lineCount := 1
	for cropReader.Scan() {
		lineCount++
		line := cropReader.Text()
		outLine := line
		// split line
		tokens := strings.Fields(line)
		if len(tokens) > 4 {
			// Field_ID    crp  sowing harvst Rex yld autorg variety comment
			// SOYSM1    SM  05151980 09311980 080 050 0
			rotationID := tokens[0]
			_, dateSowing := Datum(tokens[2])
			_, dateHarvest := Datum(tokens[3])
			if dateSowing >= dateHarvest {
				fmt.Printf("ERROR line %d: Sowing date after harvest date in crop rotation file (%s > %s) \n",
					lineCount, Kalender(dateSowing), Kalender(prevHarvest))
			}
			if prevRotationID != rotationID {
				prevHarvest = 0
				prevRotationID = rotationID
			} else {
				if dateSowing <= prevHarvest {
					fmt.Printf("Correctable ERROR line %d: Sowing date before previous harvest date in crop rotation file (%s < %s) \n",
						lineCount, Kalender(dateSowing), Kalender(prevHarvest))
					if *repair {
						// correct sowing date
						outLine = strings.Replace(line, tokens[2], KalenderOut(prevHarvest+1), 1)
					}
				}
			}
			prevHarvest = dateHarvest
		}
		if *repair {
			cropFileOut.WriteString(outLine + "\n")
		}
	}

}
