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

type SoilFileData struct {
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

func NewSoilFileData(soilID string) SoilFileData {

	return SoilFileData{
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

func LoadSoil(withGroundwater bool, LOGID string, hPath *HFilePath, soilID string) (SoilFileData, error) {

	soildata := NewSoilFileData(soilID)
	_, scannerSoilFile, err := Open(&FileDescriptior{FilePath: hPath.bofile, FileDescription: "soil file", UseFilePool: true})
	if err != nil {
		return SoilFileData{}, err
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
				soildata.GRHI = int(ValAsInt(bodenLine[70:72], "none", bodenLine))
				soildata.GRLO = soildata.GRHI
				soildata.GRW = float64(soildata.GRLO+soildata.GRHI) / 2
				soildata.GW = float64(soildata.GRLO+soildata.GRHI) / 2
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
				(&soildata).BulkDensityClassToDensity(i)
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
				return SoilFileData{}, fmt.Errorf("%s total number of 10cm layers from last soil layer is %d, should be > 1 and < 20", LOGID, soildata.N)
			}
		}
	}
	if !bofind {
		return SoilFileData{}, fmt.Errorf("SoilID '%s' not found", soilID)
	}
	return soildata, nil
}

func LoadSoilCSV(withGroundwater bool, LOGID string, hPath *HFilePath, soilID string) (SoilFileData, error) {

	soildata := NewSoilFileData(soilID)
	_, scanner, err := Open(&FileDescriptior{FilePath: hPath.bofile, FileDescription: "soil file", UseFilePool: true})
	if err != nil {
		return SoilFileData{}, err
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
				soildata.GRHI = int(ValAsInt(tokens[header[groundwaterlevel]], "none", bodenLine))
				soildata.GRLO = soildata.GRHI
				soildata.GRW = float64(soildata.GRLO+soildata.GRHI) / 2
				soildata.GW = float64(soildata.GRLO+soildata.GRHI) / 2
				soildata.AMPL = 0
			}
			soildata.DRAIDEP = int(ValAsInt(tokens[header[drainagedepth]], "none", bodenLine))
			soildata.DRAIFAK = ValAsFloat(tokens[header[drainagepercetage]], "none", bodenLine)
			soildata.UKT[0] = 0
			for i := 0; i < soildata.AZHO; i++ {
				soildata.BART[i] = tokens[header[texture]]
				textLenght := len(soildata.BART[i])
				if textLenght > 3 || textLenght == 0 {
					return SoilFileData{}, fmt.Errorf("invalid texture '%s'", soildata.BART[i])
				} else if textLenght == 2 {
					soildata.BART[i] = soildata.BART[i] + " "
				} else if textLenght == 1 {
					soildata.BART[i] = soildata.BART[i] + "  "
				}
				soildata.UKT[i+1] = int(ValAsInt(tokens[header[layerdepth]], "none", bodenLine))
				soildata.LD[i] = int(ValAsInt(tokens[header[bulkdensityclass]], "none", bodenLine))
				// read buld density classes (LD = Lagerungsdichte) set bulk density values
				(&soildata).BulkDensityClassToDensity(i)
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
				return SoilFileData{}, fmt.Errorf("%s total number of 10cm layers from last soil layer is %d, should be > 1 and < 20", LOGID, soildata.N)
			}
		}
	}
	if !bofind {
		return SoilFileData{}, fmt.Errorf("SoilID '%s' not found", soilID)
	}
	return soildata, nil
}

// BulkDensityClassToDensity set bulk density from class
func (soildata *SoilFileData) BulkDensityClassToDensity(i int) {
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

// BulkDensityToClass get bulk density class from bulk density
func (soildata *SoilFileData) BulkDensityToClass(bulkDensity float64) (bulkDensityClass int) {
	bulkDensityClass = 1
	bd := bulkDensity / 1000
	if bd < 1.3 {
		bulkDensityClass = 1
	} else if bd < 1.5 {
		bulkDensityClass = 2
	} else if bd < 1.7 {
		bulkDensityClass = 3
	} else if bd < 1.85 {
		bulkDensityClass = 4
	} else {
		bulkDensityClass = 5
	}
	return bulkDensityClass
}

// CalculatePoreVolume calculate pore volume from bulk density
func CalculatePoreVolume(bulkDensity float64) float64 {
	return 1 - ((bulkDensity / 1000) / 2.65)
}

// SandAndClayToHa5Texture get a rough KA5 soil texture class from given sand and soil content
func SandAndClayToHa5Texture(sand, clay float64) string {
	silt := 1.0 - sand - clay
	soil_texture := ""

	// SS silt 0-10% clay 0-5%
	if silt >= 0.0 && silt < 0.1 && clay >= 0.0 && clay < 0.05 {
		soil_texture = "SS "
	} else
	// ST2 silt 0-10% clay 5-17%
	if silt >= 0.0 && silt < 0.1 && clay >= 0.05 && clay < 0.17 {
		soil_texture = "ST2"
	} else
	// ST3 silt 0-15% clay 17-25%
	if silt >= 0.0 && silt < 0.15 && clay >= 0.17 && clay < 0.25 {
		soil_texture = "ST3"
	} else
	// SU2 silt 10-25% clay 0-5%
	if silt >= 0.1 && silt < 0.25 && clay >= 0.0 && clay < 0.05 {
		soil_texture = "SU2"
	} else
	// SU3 silt 25-40% clay 0-8%
	if silt >= 0.25 && silt < 0.4 && clay >= 0.0 && clay < 0.08 {
		soil_texture = "SU3"
	} else
	// SU4 silt 40-50% clay 0-8%
	if silt >= 0.4 && silt < 0.5 && clay >= 0.0 && clay < 0.08 {
		soil_texture = "SU4"
	} else
	// SL2 silt 10-25% clay 5-8%
	if silt >= 0.1 && silt < 0.25 && clay >= 0.05 && clay < 0.08 {
		soil_texture = "SL2"
	} else
	// SL3 silt 10-40% clay 8-12%
	if silt >= 0.1 && silt < 0.4 && clay >= 0.08 && clay < 0.12 {
		soil_texture = "SL3"
	} else
	// SL4 silt 10-40% clay 12-17%
	if silt >= 0.1 && silt < 0.4 && clay >= 0.12 && clay < 0.17 {
		soil_texture = "SL4"
	} else
	// SLU silt 40-50% clay 8-17%
	if silt >= 0.4 && silt < 0.5 && clay >= 0.08 && clay < 0.17 {
		soil_texture = "SLU"
	} else
	// LS2 silt 40-50% clay 17-25%
	if silt >= 0.4 && silt < 0.5 && clay >= 0.17 && clay < 0.25 {
		soil_texture = "LS2"
	} else
	// LS3 silt 30-40% clay 17-25%
	if silt >= 0.3 && silt < 0.4 && clay >= 0.17 && clay < 0.25 {
		soil_texture = "LS3"
	} else
	// LS4 silt 15-30% clay 17-25%
	if silt >= 0.15 && silt < 0.3 && clay >= 0.17 && clay < 0.25 {
		soil_texture = "LS4"
	} else
	// LT2 silt 30-50% clay 25-35%
	if silt >= 0.3 && silt < 0.5 && clay >= 0.25 && clay < 0.35 {
		soil_texture = "LT2"
	} else
	// LT3 silt 30-50% clay 35-45%
	if silt >= 0.3 && silt < 0.5 && clay >= 0.35 && clay < 0.45 {
		soil_texture = "LT3"
	} else
	// LTS silt 15-30% clay 25-45%
	if silt >= 0.15 && silt < 0.3 && clay >= 0.25 && clay < 0.45 {
		soil_texture = "LTS"
	} else
	// LU silt 50-65% clay 17-30%
	if silt >= 0.5 && silt < 0.65 && clay >= 0.17 && clay < 0.3 {
		soil_texture = "LU "
	} else
	// ULS silt 50-65% clay 8-17%
	if silt >= 0.5 && silt < 0.65 && clay >= 0.08 && clay < 0.17 {
		soil_texture = "ULS"
	} else
	// US silt 50-80% clay 0-8%
	if silt >= 0.5 && silt < 0.8 && clay >= 0.0 && clay < 0.08 {
		soil_texture = "US "
	} else
	// UU silt >80% clay 0-8%
	if silt >= 0.8 && clay >= 0.0 && clay < 0.08 {
		soil_texture = "UU "
	} else
	// UT2 silt >65% clay 8-12%
	if silt >= 0.65 && clay >= 0.08 && clay < 0.12 {
		soil_texture = "UT2"
	} else
	// UT3 silt >65% clay 12-17%
	if silt >= 0.65 && clay >= 0.12 && clay < 0.17 {
		soil_texture = "UT3"
	} else
	// UT4 silt >65% clay 17-25%
	if silt >= 0.65 && clay >= 0.17 && clay < 0.25 {
		soil_texture = "UT4"
	} else
	// TS2 silt 0-15% clay 45-65%
	if silt >= 0.0 && silt < 0.15 && clay >= 0.45 && clay < 0.65 {
		soil_texture = "TS2"
	} else
	// TS3 silt 0-15% clay 35-45%
	if silt >= 0.0 && silt < 0.15 && clay >= 0.35 && clay < 0.45 {
		soil_texture = "TS3"
	} else
	// TS4 silt 0-15% clay 25-35%
	if silt >= 0.0 && silt < 0.15 && clay >= 0.25 && clay < 0.35 {
		soil_texture = "TS4"
	} else
	// TL silt 15-30% clay 45-65%
	if silt >= 0.15 && silt < 0.3 && clay >= 0.45 && clay < 0.65 {
		soil_texture = "TL "
	} else
	// TU3 silt 50-65% clay 30-45%
	if silt >= 0.5 && silt < 0.65 && clay >= 0.3 && clay < 0.45 {
		soil_texture = "TU3"
	} else
	// TU2 silt > 30% clay 45-65%
	if silt >= 0.3 && clay >= 0.45 && clay < 0.65 {
		soil_texture = "TU2"
	} else
	// TU4 silt > 65% clay >25%
	if silt >= 0.65 && clay >= 0.25 {
		soil_texture = "TU4"
	} else
	// TT clay > 65
	if clay >= 0.65 {
		soil_texture = "TT "
	}

	return soil_texture
}

// SandAndClayToHa5Texture with percent sand and clay as integer
func SandAndClayToKa5Texture(sand, clay int) string {
	var soil_texture string
	silt := 100 - sand - clay

	// SS silt 0-10% clay 0-5%
	if silt >= 0 && silt < 10 && clay >= 0.0 && clay < 5 {
		soil_texture = "SS "
	} else
	// ST2 silt 0-10% clay 5-17%
	if silt >= 0 && silt < 10 && clay >= 5 && clay < 17 {
		soil_texture = "ST2"
	} else
	// ST3 silt 0-15% clay 17-25%
	if silt >= 0 && silt < 15 && clay >= 17 && clay < 25 {
		soil_texture = "ST3"
	} else
	// SU2 silt 10-25% clay 0-5%
	if silt >= 10 && silt < 25 && clay >= 0 && clay < 5 {
		soil_texture = "SU2"
	} else
	// SU3 silt 25-40% clay 0-8%
	if silt >= 25 && silt < 40 && clay >= 0 && clay < 8 {
		soil_texture = "SU3"
	} else
	// SU4 silt 40-50% clay 0-8%
	if silt >= 40 && silt < 50 && clay >= 0 && clay < 8 {
		soil_texture = "SU4"
	} else
	// SL2 silt 10-25% clay 5-8%
	if silt >= 10 && silt < 25 && clay >= 5 && clay < 8 {
		soil_texture = "SL2"
	} else
	// SL3 silt 10-40% clay 8-12%
	if silt >= 10 && silt < 40 && clay >= 8 && clay < 12 {
		soil_texture = "SL3"
	} else
	// SL4 silt 10-40% clay 12-17%
	if silt >= 10 && silt < 40 && clay >= 12 && clay < 17 {
		soil_texture = "SL4"
	} else
	// SLU silt 40-50% clay 8-17%
	if silt >= 40 && silt < 50 && clay >= 8 && clay < 17 {
		soil_texture = "SLU"
	} else
	// LS2 silt 40-50% clay 17-25%
	if silt >= 40 && silt < 50 && clay >= 17 && clay < 25 {
		soil_texture = "LS2"
	} else
	// LS3 silt 30-40% clay 17-25%
	if silt >= 30 && silt < 40 && clay >= 17 && clay < 25 {
		soil_texture = "LS3"
	} else
	// LS4 silt 15-30% clay 17-25%
	if silt >= 15 && silt < 30 && clay >= 17 && clay < 25 {
		soil_texture = "LS4"
	} else
	// LT2 silt 30-50% clay 25-35%
	if silt >= 30 && silt < 50 && clay >= 25 && clay < 35 {
		soil_texture = "LT2"
	} else
	// LT3 silt 30-50% clay 35-45%
	if silt >= 30 && silt < 50 && clay >= 35 && clay < 45 {
		soil_texture = "LT3"
	} else
	// LTS silt 15-30% clay 25-45%
	if silt >= 15 && silt < 30 && clay >= 25 && clay < 45 {
		soil_texture = "LTS"
	} else
	// LU silt 50-65% clay 17-30%
	if silt >= 50 && silt < 65 && clay >= 17 && clay < 30 {
		soil_texture = "LU "
	} else
	// ULS silt 50-65% clay 8-17%
	if silt >= 50 && silt < 65 && clay >= 8 && clay < 17 {
		soil_texture = "ULS"
	} else
	// US silt 50-80% clay 0-8%
	if silt >= 50 && silt < 80 && clay >= 0 && clay < 8 {
		soil_texture = "US "
	} else
	// UU silt >80% clay 0-8%
	if silt >= 80 && clay >= 0 && clay < 8 {
		soil_texture = "UU "
	} else
	// UT2 silt >65% clay 8-12%
	if silt >= 65 && clay >= 8 && clay < 12 {
		soil_texture = "UT2"
	} else
	// UT3 silt >65% clay 12-17%
	if silt >= 65 && clay >= 12 && clay < 17 {
		soil_texture = "UT3"
	} else
	// UT4 silt >65% clay 17-25%
	if silt >= 65 && clay >= 17 && clay < 25 {
		soil_texture = "UT4"
	} else
	// TS2 silt 0-15% clay 45-65%
	if silt >= 0 && silt < 15 && clay >= 45 && clay < 65 {
		soil_texture = "TS2"
	} else
	// TS3 silt 0-15% clay 35-45%
	if silt >= 0 && silt < 15 && clay >= 35 && clay < 45 {
		soil_texture = "TS3"
	} else
	// TS4 silt 0-15% clay 25-35%
	if silt >= 0 && silt < 15 && clay >= 25 && clay < 35 {
		soil_texture = "TS4"
	} else
	// TL silt 15-30% clay 45-65%
	if silt >= 15 && silt < 30 && clay >= 45 && clay < 65 {
		soil_texture = "TL "
	} else
	// TU3 silt 50-65% clay 30-45%
	if silt >= 50 && silt < 65 && clay >= 30 && clay < 45 {
		soil_texture = "TU3"
	} else
	// TU2 silt > 30% clay 45-65%
	if silt >= 30 && clay >= 45 && clay < 65 {
		soil_texture = "TU2"
	} else
	// TU4 silt > 65% clay >25%
	if silt >= 65 && clay >= 25 {
		soil_texture = "TU4"
	} else
	// TT clay > 65
	if clay >= 65 {
		soil_texture = "TT "
	}
	return soil_texture
}

func (soildata *SoilFileData) cNSetup(i int) {

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
