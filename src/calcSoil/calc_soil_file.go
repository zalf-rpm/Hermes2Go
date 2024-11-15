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
	inputFile := flag.String("input", "", "input file (as .txt or .csv)")
	// output file
	outputFile := flag.String("output", "", "output file (as .txt or .csv)")

	calcTexture := flag.Bool("texture", false, "calculate texture (Ka5 Textur)")
	calcFC := flag.Bool("fc", false, "calculate field capacity %(Feldkapazität)")
	calcWP := flag.Bool("wp", false, "calculate wilting point % (Welke Punkt)")
	calcGPV := flag.Bool("gpv", false, "calculate total pore volume % (Gesamtporenvolumen)")
	ptf := flag.Int("ptf", 0, "calculate with ptf (1,2,3,4) 0=none (Pedotransferfunktion see Hermes2Go)")
	calBulkDensity := flag.Bool("stdbulk", false, "set default bulk density class (Lagerungsdichtenklasse)")
	withBulkDensity := flag.Bool("withBD", false, "add a BulkDensity column for measured values, set to defaults (Lagerungsdichte für gemessene Werte)")

	flag.Parse()

	hpath := hermes.NewHermesFilePath("", "0", "", "", "")
	hpath.OverrideBofile(*inputFile)
	session := hermes.NewHermesSession()
	defer session.Close()

	listOfSoilIds := readSoilIds(*inputFile, session)
	soilData := make([]hermes.SoilFileData, 0)
	if strings.HasSuffix(*inputFile, ".csv") {

		// read csv file
		for _, soilId := range listOfSoilIds {
			data, err := hermes.LoadSoilCSV(true, "any", &hpath, soilId, session)
			if err != nil {
				panic(err)
			}
			soilData = append(soilData, data)
		}
	} else {
		// read txt file
		for _, soilId := range listOfSoilIds {
			data, err := hermes.LoadSoil(true, "any", &hpath, soilId, session)
			if err != nil {
				panic(err)
			}
			soilData = append(soilData, data)
		}
	}

	for i := range soilData {
		for layer := 0; layer < soilData[i].AZHO; layer++ {
			// calculate texture
			if *calcTexture {
				soilData[i].BART[layer] = hermes.SandAndClayToKa5Texture(int(soilData[i].SSAND[layer]), int(soilData[i].TON[layer]))
			}
			if *calBulkDensity {
				bulk := stdBulk(soilData[i].UKT[layer+1] - 1)
				soilData[i].LD[layer] = (&soilData[i]).BulkDensityToClass(bulk)

			}
			if *calcGPV {
				// bulkDenssity := soilData[i].BULK[layer]
				// if *calBulkDensity {
				// 	bulkDenssity = stdBulk(soilData[i].UKT[layer+1] - 1)
				// }
				//soilData[i].GPV[layer] = hermes.CalculatePoreSpace(bulkDenssity*1000) * 100
				ts := 1.0
				if layer == 0 {
					ts = 0
				}

				soilData[i].GPV[layer] = hermes.CalculatePoreSpacePTF1(soilData[i].CGEHALT[layer], soilData[i].TON[layer], soilData[i].SLUF[layer], float64(soilData[i].LD[layer]), ts) * 100
				//soilData[i].GPV[layer] = hermes.CalculatePoreSpacePTF1(soilData[i].CGEHALT[layer], soilData[i].TON[layer], soilData[i].SLUF[layer], bulkDenssity, ts) * 100
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
			if *withBulkDensity {
				soilData[i].BulkDensityClassToDensity(layer)
			}
		}

	}
	out := session.OpenResultFile(*outputFile, false)
	defer out.Close()

	// write output file
	if strings.HasSuffix(*outputFile, ".csv") {
		// write soil file header
		if *withBulkDensity {
			out.Write("SID,C_org,Texture,LayerDepth,BulkDensityClass,BulkDensity,Stone,C/N,C/S,RootDepth,NumberHorizon,FieldCapacity,WiltingPoint,PoreVolume,Sand,Silt,Clay,DrainageDepth,Drainage%,GroundWaterLevel\n")
		} else {
			out.Write("SID,C_org,Texture,LayerDepth,BulkDensityClass,Stone,C/N,C/S,RootDepth,NumberHorizon,FieldCapacity,WiltingPoint,PoreVolume,Sand,Silt,Clay,DrainageDepth,Drainage%,GroundWaterLevel\n")
		}
		// write csv file
		for _, soilData := range soilData {
			err := WriteSoilCSV(soilData, out, *withBulkDensity)
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

func WriteSoil(soilData hermes.SoilFileData, out hermes.OutWriter) error {

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

			_, err := out.Write(fmt.Sprintf("%s %s %s %02d %d %02d %02d      %02d %02d %02d   %02d %02d %02d %02d %02d %02d %02d  %02d   %02d %02d\n",
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

			_, err := out.Write(fmt.Sprintf("%s %s %s %02d %d %02d %02d      %02d         %02d %02d %02d %02d %02d %02d %02d  %02d   %02d\n",
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

func WriteSoilCSV(soilData hermes.SoilFileData, out hermes.OutWriter, withBulkdensityColumn bool) error {
	for layer := 0; layer < soilData.AZHO; layer++ {

		if !withBulkdensityColumn {
			if layer == 0 {
				//001 2.09 LS3 03 1 00 10      00 10 04   0  0  0  35 44 21 00  20   00 55
				_, err := out.Write(fmt.Sprintf("%s,%4.2f,%s,%02d,%d,%02d,%02d,00,%02d,%02d,%d,%d,%d,%02d,%02d,%02d,%02d,%02d,%02d\n",
					soilData.SoilID,
					soilData.CGEHALT[layer],
					soilData.BART[layer],
					soilData.UKT[layer+1],
					soilData.LD[layer],
					int(soilData.STEIN[layer]*100),
					int(soilData.CNRATIO[layer]),
					int(soilData.WURZMAX),
					int(soilData.AZHO),
					int(soilData.FKA[layer]),
					int(soilData.WP[layer]),
					int(soilData.GPV[layer]),
					int(soilData.SSAND[layer]),
					int(soilData.SLUF[layer]),
					int(soilData.TON[layer]),
					int(soilData.DRAIDEP),
					int(soilData.DRAIFAK),
					int(soilData.GW)))
				if err != nil {
					return err
				}

			} else {

				_, err := out.Write(fmt.Sprintf("%s,%4.2f,%s,%02d,%d,%02d,%02d,00,,,%d,%d,%d,%02d,%02d,%02d,%02d,%02d,\n",
					soilData.SoilID,
					soilData.CGEHALT[layer],
					soilData.BART[layer],
					soilData.UKT[layer+1],
					soilData.LD[layer],
					int(soilData.STEIN[layer]*100),
					int(soilData.CNRATIO[layer]),
					int(soilData.FKA[layer]),
					int(soilData.WP[layer]),
					int(soilData.GPV[layer]),
					int(soilData.SSAND[layer]),
					int(soilData.SLUF[layer]),
					int(soilData.TON[layer]),
					int(soilData.DRAIDEP),
					int(soilData.DRAIFAK)))

				if err != nil {
					return err
				}
			}
		} else {
			if layer == 0 {
				//001 2.09 LS3 03 1 00 10      00 10 04   0  0  0  35 44 21 00  20   00 55
				_, err := out.Write(fmt.Sprintf("%s,%4.2f,%s,%02d,%d,%1.2f,%02d,%02d,00,%02d,%02d,%d,%d,%d,%02d,%02d,%02d,%02d,%02d,%02d\n",
					soilData.SoilID,
					soilData.CGEHALT[layer],
					soilData.BART[layer],
					soilData.UKT[layer+1],
					soilData.LD[layer],
					soilData.BULK[layer],
					int(soilData.STEIN[layer]*100),
					int(soilData.CNRATIO[layer]),
					int(soilData.WURZMAX),
					int(soilData.AZHO),
					int(soilData.FKA[layer]),
					int(soilData.WP[layer]),
					int(soilData.GPV[layer]),
					int(soilData.SSAND[layer]),
					int(soilData.SLUF[layer]),
					int(soilData.TON[layer]),
					int(soilData.DRAIDEP),
					int(soilData.DRAIFAK),
					int(soilData.GW)))
				if err != nil {
					return err
				}

			} else {

				_, err := out.Write(fmt.Sprintf("%s,%4.2f,%s,%02d,%d,%1.2f,%02d,%02d,00,,,%d,%d,%d,%02d,%02d,%02d,%02d,%02d,\n",
					soilData.SoilID,
					soilData.CGEHALT[layer],
					soilData.BART[layer],
					soilData.UKT[layer+1],
					soilData.LD[layer],
					soilData.BULK[layer],
					int(soilData.STEIN[layer]*100),
					int(soilData.CNRATIO[layer]),
					int(soilData.FKA[layer]),
					int(soilData.WP[layer]),
					int(soilData.GPV[layer]),
					int(soilData.SSAND[layer]),
					int(soilData.SLUF[layer]),
					int(soilData.TON[layer]),
					int(soilData.DRAIDEP),
					int(soilData.DRAIFAK)))

				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func readSoilIds(inFile string, session *hermes.HermesSession) []string {
	_, scanner, err := session.Open(&hermes.FileDescriptior{FilePath: inFile, FileDescription: "soil file", UseFilePool: true})
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
