package hermes

import (
	"fmt"
	"strings"
)

// ! ++++++++++++++++++++++ Einlesen der Bodenkartiereinheit (Boden-ID) +++++++++++++++++++++++++++++
// ! INPUTS:
// ! BOD$			= Boden-ID (zum Auffinden von BOF$)
// ! AZHO           = Anzahl Horizonte des BOD$ Bodenprofils
// ! WURZMAX        = effektive Wurzeltiefe des Profils
// ! DRAIDEP        = Tiefe der Drainung (dm)
// ! DRAIFAK        = Anteil des Drainwassers am Sickerwasseranfakk (fraction)
// ! I              = Horizontzähler
// ! BART$(I)       = Bodenart nach KA5 (Besonderheits der Schreibweise beachten, Manual)
// ! UKT(I)			= Unterkante von Horizont I
// ! LD(I)          = Lagerungsdichtestufe nach KA5 (1-5)
// ! CGEHALT(I)     = Corg-Gehalt in Horizont I (Gew.%)
// ! CNRATIO(I)     = CN-Verhaeltnis in Hor. I
// ! STEIN(I)       = Steingehalt (%)
// ! optional (überschreibt Defaultwerte aus KA5 Tabelle, wenn <> 0)):
// ! FKA(I)         = Wassergehalt bei Feldkapazität (Vol. %)
// ! PWP(I)         = Wassergehalt bei PWP (Vol. %)
// ! GPV(I)         = Gesamtporenvolumen (Vol%)
// ! Ableitungen
// ! BULK(I)		= Zuweisung mittlere Lagerungsdichte von LD(I) (g/cm^3)
// ! NGEHALT(I)     = Norg-Gehalt (Gew. %)
// ! HUMUS(I)       = Humusgehalt in Hor. I (Gew.%)
// ! ------------------------------------------------------------------------------------------------

type soilFileData struct {
	SoilID                     string
	N                          int
	AZHO                       int
	WURZMAX                    int
	useGroundwaterFromSoilfile bool
	GRHI                       int
	GRLO                       int
	GRW                        float64
	GW                         float64
	AMPL                       int // amplitude by layer
	DRAIDEP                    int
	DRAIFAK                    float64
	UKT                        [11]int
	BART                       [10]string  // soil texture (Bodenart)
	LD                         [10]int     // bulk density class (Lagerungsdichteklasse)
	BULK                       [10]float64 // bulk density
	CGEHALT                    [10]float64 // C-content soil class specific in %
	CNRATIO                    [10]float64 // C/N ratio
	CNRAT1                     float64     // C/N ratio in top layer
	NGEHALT                    [10]float64
	HUMUS                      [21]float64
	STEIN                      [10]float64 // stone content
	FKA                        [10]float64 // Field capacity
	WP                         [10]float64 // wilting point
	GPV                        [10]float64 // general pore volume
	SSAND                      [10]float64 // sand in %
	SLUF                       [10]float64 // silt in %
	TON                        [10]float64 // clay in %
}

func NewSoilFileData(soilID string) soilFileData {

	return soilFileData{
		SoilID:                     soilID,
		N:                          20,
		AZHO:                       0,
		WURZMAX:                    0,
		useGroundwaterFromSoilfile: false,
		GRHI:                       0,
		GRLO:                       0,
		GRW:                        0,
		GW:                         0,
		AMPL:                       0,
		DRAIDEP:                    0,
		DRAIFAK:                    0,
		UKT:                        [11]int{},
		BART:                       [10]string{},
		LD:                         [10]int{},
		BULK:                       [10]float64{},
		CGEHALT:                    [10]float64{},
		CNRATIO:                    [10]float64{},
		CNRAT1:                     0,
		NGEHALT:                    [10]float64{},
		HUMUS:                      [21]float64{},
		STEIN:                      [10]float64{},
		FKA:                        [10]float64{},
		WP:                         [10]float64{},
		GPV:                        [10]float64{},
		SSAND:                      [10]float64{},
		SLUF:                       [10]float64{},
		TON:                        [10]float64{},
	}
}

func loadSoil(withGroundwater bool, LOGID string, hPath *HFilePath, soilID string) (soilFileData, error) {

	soildata := NewSoilFileData(soilID)
	_, scannerSoilFile, err := Open(&FileDescriptior{FilePath: hPath.bofile, FileDescription: "soil file", UseFilePool: true})
	if err != nil {
		return soilFileData{}, err
	}
	LineInut(scannerSoilFile)
	bofind := false

	for scannerSoilFile.Scan() {
		bodenLine := scannerSoilFile.Text()
		boden := bodenLine[0:3] // SID - first 3 character
		if boden == soilID {

			bofind = true
			soildata.AZHO = int(ValAsInt(bodenLine[35:37], "none", bodenLine))
			soildata.WURZMAX = int(ValAsInt(bodenLine[32:34], "none", bodenLine))

			soildata.useGroundwaterFromSoilfile = withGroundwater
			if withGroundwater {
				gw := int(ValAsInt(bodenLine[70:72], "none", bodenLine))
				soildata.GRHI = gw
				soildata.GRLO = gw
				soildata.GRW = float64(gw)
				soildata.GW = float64(gw)
				soildata.AMPL = 0
			}
			soildata.DRAIDEP = int(ValAsInt(bodenLine[62:64], "none", bodenLine))
			soildata.DRAIFAK = ValAsFloat(bodenLine[67:70], "none", bodenLine)
			soildata.UKT[0] = 0
			for i := 0; i < soildata.AZHO; i++ {
				soildata.BART[i] = bodenLine[9:12]
				soildata.UKT[i+1] = int(ValAsInt(bodenLine[13:15], "none", bodenLine))
				soildata.LD[i] = int(ValAsInt(bodenLine[16:17], "none", bodenLine))
				// read buld density classes (LD = Lagerungsdichte) set bulk density values
				(&soildata).bulkDensityClassToDensity(i)
				// C-content soil class specific in %
				soildata.CGEHALT[i] = ValAsFloat(bodenLine[4:8], "none", bodenLine)
				// C/N ratio
				soildata.CNRATIO[i] = ValAsFloat(bodenLine[21:24], "none", bodenLine)
				(&soildata).cNSetup(i)
				soildata.STEIN[i] = ValAsFloat(bodenLine[18:20], "none", bodenLine) / 100
				// Field capacity

				value, err := TryValAsFloat(bodenLine[40:42])
				if err == nil {
					soildata.FKA[i] = value
				}
				// wilting point
				value, err = TryValAsFloat(bodenLine[43:45])
				if err == nil {
					soildata.WP[i] = value
				}
				// general pore volume
				value, err = TryValAsFloat(bodenLine[46:48])
				if err == nil {
					soildata.GPV[i] = value
				}
				// sand in %
				value, err = TryValAsFloat(bodenLine[49:51])
				if err == nil {
					soildata.SSAND[i] = value
				}
				// silt in %
				value, err = TryValAsFloat(bodenLine[52:54])
				if err == nil {
					soildata.SLUF[i] = value
				}
				// clay in %
				value, err = TryValAsFloat(bodenLine[55:57])
				if err == nil {
					soildata.TON[i] = value
				}
				if i+1 < soildata.AZHO {
					// scan next line in soil profile
					bodenLine = LineInut(scannerSoilFile)
				}
			}
			// get total number of 10cm layers from last soil layer
			soildata.N = soildata.UKT[soildata.AZHO]
			if soildata.N > 20 || soildata.N < 1 {
				return soilFileData{}, fmt.Errorf("%s total number of 10cm layers from last soil layer is %d, should be > 1 and < 20", LOGID, soildata.N)
			}
		}
	}
	if !bofind {
		return soilFileData{}, fmt.Errorf("SoilID '%s' not found", soilID)
	}
	return soildata, nil
}

func loadSoilCSV(withGroundwater bool, LOGID string, hPath *HFilePath, soilID string) (soilFileData, error) {

	soildata := NewSoilFileData(soilID)
	_, scanner, err := Open(&FileDescriptior{FilePath: hPath.bofile, FileDescription: "soil file", UseFilePool: true})
	if err != nil {
		return soilFileData{}, err
	}
	headline := LineInut(scanner)
	header := readSoilHeader(headline)

	bofind := false

	for scanner.Scan() {
		bodenLine := scanner.Text()
		tokens := strings.Split(bodenLine, ",")
		boden := tokens[header[sid]]
		if boden == soilID {

			bofind = true
			soildata.AZHO = int(ValAsInt(tokens[header[numberhorizon]], "none", bodenLine))
			soildata.WURZMAX = int(ValAsInt(tokens[header[rootdepth]], "none", bodenLine))

			soildata.useGroundwaterFromSoilfile = withGroundwater
			if withGroundwater {
				gw := int(ValAsInt(tokens[header[groundwaterlevel]], "none", bodenLine))
				soildata.GRHI = gw
				soildata.GRLO = gw
				soildata.GRW = float64(gw)
				soildata.GW = float64(gw)
				soildata.AMPL = 0
			}
			soildata.DRAIDEP = int(ValAsInt(tokens[header[drainagedepth]], "none", bodenLine))
			soildata.DRAIFAK = ValAsFloat(tokens[header[drainagepercetage]], "none", bodenLine)
			soildata.UKT[0] = 0
			for i := 0; i < soildata.AZHO; i++ {
				soildata.BART[i] = tokens[header[texture]]
				textLenght := len(soildata.BART[i])
				if textLenght > 3 || textLenght == 0 {
					return soilFileData{}, fmt.Errorf("invalid texture '%s'", soildata.BART[i])
				} else if textLenght == 2 {
					soildata.BART[i] = soildata.BART[i] + " "
				} else if textLenght == 1 {
					soildata.BART[i] = soildata.BART[i] + "  "
				}
				soildata.UKT[i+1] = int(ValAsInt(tokens[header[layerdepth]], "none", bodenLine))
				soildata.LD[i] = int(ValAsInt(tokens[header[bulkdensityclass]], "none", bodenLine))
				// read buld density classes (LD = Lagerungsdichte) set bulk density values
				(&soildata).bulkDensityClassToDensity(i)
				// C-content soil class specific in %
				soildata.CGEHALT[i] = ValAsFloat(tokens[header[corg]], "none", bodenLine)
				// C/N ratio
				soildata.CNRATIO[i] = ValAsFloat(tokens[header[c_n]], "none", bodenLine)
				(&soildata).cNSetup(i)
				soildata.STEIN[i] = ValAsFloat(tokens[header[stone]], "none", bodenLine) / 100
				// Field capacity

				value, err := TryValAsFloat(tokens[header[fieldcapacity]])
				if err == nil {
					soildata.FKA[i] = value
				}
				// wilting point
				value, err = TryValAsFloat(tokens[header[wiltingpoint]])
				if err == nil {
					soildata.WP[i] = value
				}
				// general pore volume
				value, err = TryValAsFloat(tokens[header[porevolume]])
				if err == nil {
					soildata.GPV[i] = value
				}
				// sand in %
				value, err = TryValAsFloat(tokens[header[sand]])
				if err == nil {
					soildata.SSAND[i] = value
				}
				// silt in %
				value, err = TryValAsFloat(tokens[header[silt]])
				if err == nil {
					soildata.SLUF[i] = value
				}
				// clay in %
				value, err = TryValAsFloat(tokens[header[clay]])
				if err == nil {
					soildata.TON[i] = value
				}
				if i+1 < soildata.AZHO {
					// scan next line in soil profile
					bodenLine = LineInut(scanner)
					tokens = strings.Split(bodenLine, ",")
				}
			}
			// get total number of 10cm layers from last soil layer
			soildata.N = soildata.UKT[soildata.AZHO]
			if soildata.N > 20 || soildata.N < 1 {
				return soilFileData{}, fmt.Errorf("%s total number of 10cm layers from last soil layer is %d, should be > 1 and < 20", LOGID, soildata.N)
			}
		}
	}
	if !bofind {
		return soilFileData{}, fmt.Errorf("SoilID '%s' not found", soilID)
	}
	return soildata, nil
}

func (soildata *soilFileData) bulkDensityClassToDensity(i int) {
	// read buld density classes (LD = Lagerungsdichte) set bulk density values
	if soildata.LD[i] == 1 {
		soildata.BULK[i] = 1.1
	} else if soildata.LD[i] == 2 {
		soildata.BULK[i] = 1.3
	} else if soildata.LD[i] == 3 {
		soildata.BULK[i] = 1.5
	} else if soildata.LD[i] == 4 {
		soildata.BULK[i] = 1.7
	} else if soildata.LD[i] == 5 {
		soildata.BULK[i] = 1.85
	}
}

func (soildata *soilFileData) cNSetup(i int) {

	if soildata.CNRATIO[i] == 0 {
		soildata.CNRATIO[i] = 10
	}
	if i == 0 {
		soildata.CNRAT1 = soildata.CNRATIO[i]
	}
	soildata.NGEHALT[i] = soildata.CGEHALT[i] / soildata.CNRATIO[i]
	soildata.HUMUS[i] = soildata.CGEHALT[i] * 1.72 / 100
}

// Header for csv weather files
type SoilHeader int

const (
	sid SoilHeader = iota
	corg
	texture
	layerdepth
	bulkdensityclass
	stone
	c_n
	c_s
	rootdepth
	numberhorizon
	fieldcapacity
	wiltingpoint
	porevolume
	sand
	silt
	clay
	drainagedepth
	drainagepercetage
	groundwaterlevel
)

var soilHeaderNames = map[string]SoilHeader{
	"SID":              sid,
	"C_org":            corg,
	"Texture":          texture,
	"LayerDepth":       layerdepth,
	"BulkDensityClass": bulkdensityclass,
	"Stone":            stone,
	"C/N":              c_n,
	"C/S":              c_s,
	"RootDepth":        rootdepth,
	"NumberHorizon":    numberhorizon,
	"FieldCapacity":    fieldcapacity,
	"WiltingPoint":     wiltingpoint,
	"PoreVolume":       porevolume,
	"Sand":             sand,
	"Silt":             silt,
	"Clay":             clay,
	"DrainageDepth":    drainagedepth,
	"Drainage%":        drainagepercetage,
	"GroundWaterLevel": groundwaterlevel,
}

func readSoilHeader(line string) map[SoilHeader]int {
	tokens := Explode(line, []rune{',', ';'})
	headers := make(map[SoilHeader]int)
	for kHeader, vHeader := range soilHeaderNames {
		for i, token := range tokens {
			if token == kHeader {
				headers[vHeader] = i
				break
			}
		}
	}
	return headers
}

// ReadGroundWaterTimeSeries read ground water time series
// from csv file, for a given soilID (sid)
func ReadGroundWaterTimeSeries(g *GlobalVarsMain, hPath *HFilePath, sid string) error {

	groundWaterHeader := map[string]int{
		"SID":   0,
		"Date":  1,
		"Level": 2,
	}
	// ground water time series
	// key: date as int
	// value: ground water level in layer
	g.GWTimeSeriesValues = make(map[int]float64)
	g.GWTimestamps = make([]int, 0)

	_, scanner, err := Open(&FileDescriptior{FilePath: hPath.gwtimeseries, FileDescription: "ground water time series file", UseFilePool: true})
	if err != nil {
		return err
	}
	scanner.Scan() // skip header

	// read data
	for scanner.Scan() {
		// check for soil id
		if HasPrefixWithSeperator(scanner.Text(), sid) {
			tokens := Explode(scanner.Text(), []rune{',', ';'})
			// read date
			_, date := g.Datum(tokens[groundWaterHeader["Date"]])
			// read ground water level
			level := ValAsFloat(tokens[groundWaterHeader["Level"]], hPath.gwtimeseries, scanner.Text())
			g.GWTimeSeriesValues[date] = level
			g.GWTimestamps = append(g.GWTimestamps, date)
		}
	}
	return nil
}

func HasPrefixWithSeperator(s, prefix string) bool {
	return len(s) >= len(prefix) && s[0:len(prefix)] == prefix && (s[len(prefix)] == ',' || s[len(prefix)] == ';' || s[len(prefix)] == ' ')
}

// GetGroundWaterLevel returns ground water level for a given date
// if no ground water level is found, a ground water level is calculated from the previous and next date
func GetGroundWaterLevel(g *GlobalVarsMain, date int) (float64, error) {

	// check if ground water level is given for date
	if level, ok := g.GWTimeSeriesValues[date]; ok {
		return level, nil
	}
	// no ground water level is given for date
	// calculate ground water level from previous and next date
	var prevDate, nextDate int
	for d := range g.GWTimestamps {
		if d < date {
			prevDate = d
		} else if d > date {
			nextDate = d
			break
		}
	}
	if prevDate == 0 && nextDate == 0 {
		return 0, fmt.Errorf("no ground water level found for date %d", date)
	} else if prevDate == 0 {
		return g.GWTimeSeriesValues[nextDate], nil
	} else if nextDate == 0 {
		return g.GWTimeSeriesValues[prevDate], nil
	}
	// calculate ground water level
	level := (g.GWTimeSeriesValues[nextDate]-g.GWTimeSeriesValues[prevDate])/float64(nextDate-prevDate)*float64(date-prevDate) + g.GWTimeSeriesValues[prevDate]
	return level, nil
}
