package main

import (
	"flag"
	"fmt"
	"math"
	"strings"

	"github.com/zalf-rpm/Hermes2Go/hermes"
)

func main() {
	// read hermes soil file
	// calculate FC, WP, and GPV(PS)

	// input file
	inputFile := flag.String("input", "", "input file")
	// output file
	outputFile := flag.String("output", "", "output file")

	calcTexture := flag.Bool("texture", false, "calculate texture")
	calcFC := flag.Bool("fc", false, "calculate field capacity")
	calcWP := flag.Bool("wp", false, "calculate wilting point")
	calcGPV := flag.Bool("gpv", false, "calculate gpv")
	ptf := flag.Int("ptf", 0, "calculate with ptf")
	calBulkDensity := flag.Bool("bulk", false, "calculate bulk density")

	flag.Parse()

	hpath := hermes.NewHermesFilePath("", "0", "", "", "")
	hpath.OverrideBofile(*inputFile)
	defer hermes.HermesFilePool.Close()

	listOfSoilIds := readSoilIds(*inputFile)
	soilData := make([]hermes.SoilFileData, 0)
	if strings.HasSuffix(*inputFile, ".csv") {

		// read csv file
		for _, soilId := range listOfSoilIds {
			data, err := hermes.LoadSoilCSV(true, "any", &hpath, soilId)
			if err != nil {
				panic(err)
			}
			soilData = append(soilData, data)
		}
	} else {
		// read txt file
		for _, soilId := range listOfSoilIds {
			data, err := hermes.LoadSoil(true, "any", &hpath, soilId)
			if err != nil {
				panic(err)
			}
			soilData = append(soilData, data)
		}
	}

	// calculate texture
	for i := range soilData {
		for layer := 0; layer < soilData[i].AZHO; layer++ {
			if *calcTexture {
				soilData[i].BART[layer] = hermes.SandAndClayToKa5Texture(int(soilData[i].SSAND[layer]), int(soilData[i].TON[layer]))
			}
			if *calBulkDensity {
				bulk := stdBulk(soilData[i].UKT[layer+1] - 1)
				soilData[i].LD[layer] = (&soilData[i]).BulkDensityToClass(bulk * 1000)

			}
			if *calcGPV {
				bulkDenssity := soilData[i].BULK[layer]
				if *calBulkDensity {
					bulkDenssity = stdBulk(soilData[i].UKT[layer+1] - 1)
				}
				soilData[i].GPV[layer] = hermes.CalculatePoreVolume(bulkDenssity*1000) * 100
			}
			if *ptf > 0 {
				var fc, wp float64
				if *ptf == 1 {
					fc, wp = hermes.PTF1(soilData[i].CGEHALT[layer], soilData[i].TON[layer], soilData[i].SLUF[layer])
				} else if *ptf == 2 {
					fc, wp = hermes.PTF2(soilData[i].CGEHALT[layer], soilData[i].TON[layer], soilData[i].SLUF[layer])
				} else if *ptf == 3 {
					fc, wp = hermes.PTF3(soilData[i].CGEHALT[layer], soilData[i].TON[layer], soilData[i].SLUF[layer])
				} else if *ptf == 4 {
					fc, wp = hermes.PTF4(soilData[i].CGEHALT[layer], soilData[i].TON[layer], soilData[i].SSAND[layer])
				}
				if *calcFC {
					soilData[i].FKA[layer] = fc * 100
				}
				if *calcWP {
					soilData[i].WP[layer] = wp * 100
				}
			}
		}
	}
	out := hermes.OpenResultFile(*outputFile, false)
	defer out.Close()

	// write output file
	if strings.HasSuffix(*outputFile, ".csv") {
		// write soil file header
		out.Write("SID,C_org,Texture,LayerDepth,BulkDensityClass,Stone,C/N,C/S,RootDepth,NumberHorizon,FieldCapacity,WiltingPoint,PoreVolume,Sand,Silt,Clay,DrainageDepth,Drainage%,GroundWaterLevel\n")
		// write csv file
		for _, soilData := range soilData {
			err := WriteSoilCSV(soilData, out)
			if err != nil {
				panic(err)
			}
		}
	} else {
		// write soil file header
		out.Write("SID Corg Tex lb B St C/N C/S Hy Rd NuHo FC WP PS S% Si C% lmd drdp drfGW\n")
		// write txt file
		for _, soilData := range soilData {
			err := WriteSoil(soilData, out)
			if err != nil {
				panic(err)
			}
		}
	}

}

func WriteSoil(soilData hermes.SoilFileData, out *hermes.Fout) error {

	formatCGEHALT := func(val float64) string {
		if math.Round(val) >= 10 {
			return fmt.Sprintf("%4.1f", val)
		}
		return fmt.Sprintf("%4.2f", val)
	}
	// write soil file
	for layer := 0; layer < soilData.AZHO; layer++ {
		if layer == 0 {
			//001 2.09 LS3 03 1 00 10      00 10 04   0  0  0  35 44 21 00  20   00 55

			_, err := out.Write(fmt.Sprintf("%s %s %s %02d %d %02d %02d      %02d %02d %02d   %02d %02d %02d  %02d %02d %02d %02d  %02d   %02d %02d\n",
				soilData.SoilID,
				formatCGEHALT(soilData.CGEHALT[layer]),
				soilData.BART[layer],
				soilData.UKT[layer+1],
				soilData.LD[layer],
				int(soilData.STEIN[layer]*100),
				int(soilData.CNRATIO[layer]),
				0,
				int(soilData.WURZMAX),
				int(soilData.AZHO),
				int(soilData.FKA[layer]),
				int(soilData.WP[layer]),
				int(soilData.GPV[layer]),
				int(soilData.SSAND[layer]),
				int(soilData.SLUF[layer]),
				int(soilData.TON[layer]),
				0,
				int(soilData.DRAIDEP),
				int(soilData.DRAIFAK),
				int(soilData.GW)))
			if err != nil {
				return err
			}

		} else {

			_, err := out.Write(fmt.Sprintf("%s %s %s %02d %d %02d %02d      %02d         %02d %02d %02d  %02d %02d %02d %02d  %02d   %02d\n",
				soilData.SoilID,
				formatCGEHALT(soilData.CGEHALT[layer]),
				soilData.BART[layer],
				soilData.UKT[layer+1],
				soilData.LD[layer],
				int(soilData.STEIN[layer]*100),
				int(soilData.CNRATIO[layer]),
				0,
				int(soilData.FKA[layer]),
				int(soilData.WP[layer]),
				int(soilData.GPV[layer]),
				int(soilData.SSAND[layer]),
				int(soilData.SLUF[layer]),
				int(soilData.TON[layer]),
				0,
				int(soilData.DRAIDEP),
				int(soilData.DRAIFAK)))

			if err != nil {
				return err
			}
		}
	}

	return nil
}

func WriteSoilCSV(soilData hermes.SoilFileData, out *hermes.Fout) error {
	for layer := 0; layer < soilData.AZHO; layer++ {
		if layer == 0 {
			//001 2.09 LS3 03 1 00 10      00 10 04   0  0  0  35 44 21 00  20   00 55

			_, err := out.Write(fmt.Sprintf("%s,f%4.2f,%s,%02d,%d,%02d,%02d,,%02d,%02d,%02d,%d,%d,%d,%02d,%02d,%02d,%02d,%02d,%02d %02d\n",
				soilData.SoilID,
				soilData.CGEHALT[layer],
				soilData.BART[layer],
				soilData.UKT[layer+1],
				soilData.LD[layer],
				int(soilData.STEIN[layer]*100),
				int(soilData.CNRATIO[layer]),
				0,
				int(soilData.WURZMAX),
				int(soilData.AZHO),
				int(soilData.FKA[layer]),
				int(soilData.WP[layer]),
				int(soilData.GPV[layer]),
				int(soilData.SSAND[layer]),
				int(soilData.SLUF[layer]),
				int(soilData.TON[layer]),
				0,
				int(soilData.DRAIDEP),
				int(soilData.DRAIFAK),
				int(soilData.GW)))
			if err != nil {
				return err
			}

		} else {

			_, err := out.Write(fmt.Sprintf("%s,f%4.2f,%s,%02d,%d,%02d,%02d,,%02d,,,%d,%d,%d,%02d,%02d,%02d,%02d,%02d,%02d,\n",
				soilData.SoilID,
				soilData.CGEHALT[layer],
				soilData.BART[layer],
				soilData.UKT[layer+1],
				soilData.LD[layer],
				int(soilData.STEIN[layer]*100),
				int(soilData.CNRATIO[layer]),
				0,
				int(soilData.FKA[layer]),
				int(soilData.WP[layer]),
				int(soilData.GPV[layer]),
				int(soilData.SSAND[layer]),
				int(soilData.SLUF[layer]),
				int(soilData.TON[layer]),
				0,
				int(soilData.DRAIDEP),
				int(soilData.DRAIFAK)))

			if err != nil {
				return err
			}
		}
	}
	return nil
}

func readSoilIds(inFile string) []string {
	_, scanner, err := hermes.Open(&hermes.FileDescriptior{FilePath: inFile, FileDescription: "soil file", UseFilePool: true})
	if err != nil {
		panic(err)
	}
	// dump header
	hermes.LineInut(scanner)

	soilIds := make(map[string]bool)
	soilIdList := make([]string, 0, 20)
	for scanner.Scan() {
		bodenLine := scanner.Text()
		tokens := strings.FieldsFunc(bodenLine, func(r rune) bool {
			return r == ',' || r == ' '
		})

		soilId := tokens[0]
		if soilId == "end" {
			break
		}
		if _, ok := soilIds[soilId]; !ok {
			soilIdList = append(soilIdList, soilId)
		}
		soilIds[soilId] = true

	}
	return soilIdList
}
func stdBulk(layer int) float64 {
	if layer < 3 {
		return 1.4
	} else if layer < 6 {
		return 1.5
	} else if layer < 9 {
		return 1.6
	} else {
		return 1.7
	}
}
