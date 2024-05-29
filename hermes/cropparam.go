package hermes

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	yaml "gopkg.in/yaml.v3"
)

// crop parameter IO

type CropParam struct {
	// CropParam is a struct to hold the crop parameters
	CropName string `yaml:"CropName" comment:"name of the crop"`
	ABBr     string `yaml:"CropAbbreviation" comment:"crop no./ abbreviation"` // Abbreviation of the crop
	Variety  string `yaml:"Variaty" comment:"variaty of the crop"`             // Variaty of the crop

	MAXAMAX           float64       `yaml:"MAXAMAX" comment:"AMAX Max. CO2 assimilation rate (kg CO2/ha leave/h)"`                               // AMAX Max. CO2 assimilation rate (kg CO2/ha leave/h)
	TempTyp           int           `yaml:"TempTyp" comment:"type of temperature dependency (C3 = 1/ C4 = 2)"`                                   // type of temperature dependency (C3 = 1/ C4 = 2)
	MINTMP            float64       `yaml:"MINTMP" comment:"minimum temperature crop growth (in C°)"`                                            // minimum temperature crop growth (in C°)
	WUMAXPF           float64       `yaml:"WUMAXPF" comment:"crop specific maximum effective rooting depth(dm)"`                                 // crop specific maximum effective rooting depth(dm)
	VELOC             float64       `yaml:"VELOC" comment:"root depth increase in mm/C°"`                                                        // root depth increase in mm/C°
	NGEFKT            int           `yaml:"NGEFKT" comment:"crop N-content function number for critical and max. N-contents"`                    // crop N-content function number for critical and max. N-contents
	RGA               float64       `yaml:"RGA,omitempty" comment:" RGA parameter for crop N-content function number 5"`                         // RGA parameter for crop N-content function number 5
	RGB               float64       `yaml:"RGB,omitempty" comment:"RGB parameter for crop N-content function number 5"`                          // RGB parameter for crop N-content function number 5
	SubOrgan          int           `yaml:"SubOrgan,omitempty" comment:"SubOrgan parameter for crop N-content function number"`                  // SubOrgan parameter for crop N-content function number 5
	AboveGroundOrgans []int         `yaml:"AboveGroundOrgans" comment:"list of above ground organs (numbers of compartiments increasing order)"` // SubOrgan parameter for crop N-content function number
	YORGAN            int           `yaml:"YORGAN" comment:"organ number for yield"`                                                             // organ number for yield
	YIFAK             float64       `yaml:"YIFAK" comment:"fraction of yield organ (90% = 0.90)"`                                                // fraction of yield organ (90% = 0.90)
	INITCONCNBIOM     float64       `yaml:"INITCONCNBIOM" comment:"start concentration N in above ground biomass (% i. d.m.)"`                   // start conzentration N in above ground biomass (% i. d.m.)
	INITCONCNROOT     float64       `yaml:"INITCONCNROOT" comment:"start concentration N in roots (% i. d.m.)"`                                  // start concentration N in roots (% i. d.m.)
	NRKOM             int           `yaml:"NRKOM" comment:"Number of crop compartiments"`                                                        // Number of crop compartiments
	CompartimentNames []string      `yaml:"CompartimentNames" comment:"list of compartiment names"`                                              // list of compartiment names
	DAUERKULT         FeatureSwitch `yaml:"DAUERKULT" comment:"Dauerkultur - Is Permaculture true/false 1/0"`                                    // Dauerkultur - Permaculture D / Non Permaculture 0
	LEGUM             FeatureSwitch `yaml:"LEGUM" comment:"Legume - Is Legume true/false 1/0"`                                                   // Legume L / Non Legume 0
	WORG              []float64     `yaml:"WORG" comment:"initial weight kg d.m./ha of organ I"`                                                 // initial weight kg d.m./ha of organ I
	MAIRT             []float64     `yaml:"MAIRT" comment:"maintainance rates of organ I (1/day)"`                                               // Maintainance rates of organ I
	KcIni             float64       `yaml:"KcIni" comment:"initial kc factor for evapotranspiration (uncovered soil)"`                           // initial kc factor for evapotranspiration (uncovered soil)
	// sulfonie

	SGEFKT       int     `yaml:"SGEFKT,omitempty" comment:"crop S-content function number for critical and max. S-contents"` // crop S-content function number for critical and max. S-contents
	SFunctExp    float64 `yaml:"SFunctExp,omitempty" comment:"exponent for S-content function"`                              // exponent for S-content function
	SCritContent float64 `yaml:"SCritContent,omitempty" comment:"base for S-content function"`                               // base for S-content function

	// crop development stages
	NRENTW                int                    `yaml:"NRENTW" comment:"number of development phases(max 10)"`   // number of development phases(max 10)
	CropDevelopmentStages []CropDevelopmentStage `yaml:"CropDevelopmentStages" comment:"development stage/phase"` // development stage/phase
}
type CropDevelopmentStage struct {
	DevelopmentStageName string    `yaml:"DevelopmentStageName" comment:"name of the development stage/phase"`       // name of the development stage/phase
	ENDBBCH              int       `yaml:"ENDBBCH,omitempty" comment:"end on BBCH-scale"`                            // end on BBCH-scale
	TSUM                 float64   `yaml:"TSUM" comment:"development phase temperature sum (°C days)"`               // development phase temperatur sum (°C days)
	BAS                  float64   `yaml:"BAS" comment:"base temperature in phase (°C)"`                             // base temperature in phase (°C)
	VSCHWELL             float64   `yaml:"VSCHWELL" comment:"vernalisation requirements (days)"`                     // vernalisation requirements (days)
	DAYL                 float64   `yaml:"DAYL" comment:"day length requirements (hours)"`                           // day length requirements (hours)
	DLBAS                float64   `yaml:"DLBAS" comment:"base day length in phase (hours)"`                         // base day length in phase (hours)
	DRYSWELL             float64   `yaml:"DRYSWELL" comment:"drought stress below ETA/ETP-quotient"`                 // drought stress below ETA/ETP-quotient
	LUKRIT               float64   `yaml:"LUKRIT" comment:"critical aircontent in topsoil (cm^3/cm^3)"`              // critical aircontent in topsoil (cm^3/cm^3)
	LAIFKT               float64   `yaml:"LAIFKT" comment:"specific leave area (LAI per mass) (ha/kg TM)"`           // specific leave area (area per mass) (m2/m2/kg TM)
	WGMAX                float64   `yaml:"WGMAX" comment:"N-content root at the end of phase (fraction)"`            // N-content root end at the of phase
	WGSMAX               float64   `yaml:"WGSMAX,omitempty" comment:"S-content root at the end of phase (fraction)"` // S-content root end at the of phase
	PRO                  []float64 `yaml:"PRO" comment:"Partitioning at end of phase (fraction, sum should be 1)"`   // Partitioning at end of phase (fraction)
	DEAD                 []float64 `yaml:"DEAD" comment:"death rate at end of phase (coefficient, 1/day)"`           // death rate at end of phase (coefficient)
	Kc                   float64   `yaml:"Kc" comment:"kc factor for evapotranspiration at end of phase"`            // kc factor for evapotranspiration at end of phase

}

// ReadCropParam reads the crop parameters from a yml file
func ReadCropParamFromFile(filename string) (CropParam, error) {

	cropParam := CropParam{}
	ymlData, err := os.ReadFile(filename)
	if err != nil {
		return cropParam, err
	}
	err = yaml.Unmarshal(ymlData, &cropParam)
	if err != nil {
		return cropParam, err
	}
	return cropParam, nil
}

// WriteCropParam writes the crop parameters to a yml file (with comments?)
func WriteCropParam(filename string, cropParam CropParam) error {

	ymlNode, err := toYamlNode(cropParam)
	if err != nil {
		return err
	}
	ymldata, err := yaml.Marshal(ymlNode)

	//ymldata, err := yaml.Marshal(cropParam)
	if err != nil {
		return err
	}
	err = os.WriteFile(filename, ymldata, 0644)
	if err != nil {
		return err
	}

	return nil
}

func ReadCropParamYml(PARANAM string, l *CropSharedVars, g *GlobalVarsMain) {
	ymlData, err := ReadFile(&FileDescriptior{FilePath: PARANAM, FileDescription: "crop file", UseFilePool: true})
	if err != nil {
		log.Fatalf("Error reading crop parameters from file %s: %s", PARANAM, err)
	}
	cropParam := CropParam{}
	err = yaml.Unmarshal(ymlData, &cropParam)
	if err != nil {
		log.Fatalf("Error reading crop parameters from file %s: %s", PARANAM, err)
	}

	g.MAXAMAX = cropParam.MAXAMAX
	l.temptyp = cropParam.TempTyp
	g.MINTMP = cropParam.MINTMP
	g.WUMAXPF = cropParam.WUMAXPF
	g.VELOC = cropParam.VELOC / 200
	g.NGEFKT = cropParam.NGEFKT
	g.RGA = cropParam.RGA
	g.RGB = cropParam.RGB
	g.SubOrgan = cropParam.SubOrgan
	g.YORGAN = cropParam.YORGAN
	g.YIFAK = cropParam.YIFAK
	g.CRITSGEHALT = cropParam.SCritContent
	g.SGEFKT = cropParam.SGEFKT
	g.CRITSEXP = cropParam.SFunctExp

	g.DAUERKULT = bool(cropParam.DAUERKULT)
	g.LEGUM = bool(cropParam.LEGUM)

	g.NRKOM = cropParam.NRKOM
	maxOrgans := 5
	if g.NRKOM > maxOrgans {
		log.Fatalf("Error: too many crop compartiments! File: %s \n", PARANAM)
	}
	l.AboveGroundOrgans = cropParam.AboveGroundOrgans
	// check that all organs are in the range of 1 to NRKOM
	for _, organ := range l.AboveGroundOrgans {
		if organ < 1 || organ > g.NRKOM {
			log.Fatalf("Error: reading above ground organs! File: %s \n", PARANAM)
		}
	}

	l.NRENTW = cropParam.NRENTW
	maxStages := 10
	if l.NRENTW > maxStages {
		log.Fatalf("Error: too many development stages! File: %s \n", PARANAM)
	}

	if !g.DAUERKULT {
		ResetStages(g)
		g.PHYLLO, g.VERNTAGE = 0, 0
		for i := 0; i < maxOrgans; i++ {
			for i2 := 0; i2 < maxStages; i2++ {
				g.SUM[i2] = 0
				g.DEV[i2] = 0
				g.PRO[i2][i] = 0
				g.DEAD[i2][i] = 0
				g.TROOTSUM = 0
			}
		}
	} else {
		for i := 0; i < maxOrgans; i++ {
			for i2 := 2; i2 < maxStages; i2++ {
				g.SUM[i2] = 0
				g.PRO[i2][i] = 0
				g.DEAD[i2][i] = 0
				g.TROOTSUM = 0
			}
		}
	}

	// check if NRKOM is the same as the length of the CompartimentNames, WORG and MAIRT
	if g.NRKOM != len(cropParam.CompartimentNames) {
		log.Fatalf("Error: reading crop compartiment names! File: %s \n", PARANAM)
	}
	if g.NRKOM != len(cropParam.WORG) {
		log.Fatalf("Error: reading crop WORG! File: %s \n", PARANAM)
	}
	if g.NRKOM != len(cropParam.MAIRT) {
		log.Fatalf("Error: reading crop MAIRT! File: %s \n", PARANAM)
	}
	for i := 0; i < g.NRKOM; i++ {
		if !(g.DAUERKULT && g.AKF.Num > 2 && g.FRUCHT[g.AKF.Index] == g.FRUCHT[g.AKF.Index-1]) {
			g.WORG[i] = cropParam.WORG[i]
		}
		g.MAIRT[i] = cropParam.MAIRT[i]
		g.WDORG[i] = 0
	}
	if !(g.DAUERKULT && g.AKF.Num > 2 && g.FRUCHT[g.AKF.Index] == g.FRUCHT[g.AKF.Index-1]) {
		g.GEHOB = cropParam.INITCONCNBIOM / 100
		g.WUGEH = cropParam.INITCONCNROOT / 100
	}

	l.kcini = cropParam.KcIni
	l.tendsum = 0
	l.useBBCH = false
	for i := 0; i < l.NRENTW; i++ {
		l.ENDBBCH[i] = float64(cropParam.CropDevelopmentStages[i].ENDBBCH)
		l.useBBCH = l.useBBCH || l.ENDBBCH[i] > 0
		g.TSUM[i] = cropParam.CropDevelopmentStages[i].TSUM
		g.BAS[i] = cropParam.CropDevelopmentStages[i].BAS
		g.VSCHWELL[i] = cropParam.CropDevelopmentStages[i].VSCHWELL
		g.DAYL[i] = cropParam.CropDevelopmentStages[i].DAYL
		g.DLBAS[i] = cropParam.CropDevelopmentStages[i].DLBAS
		g.DRYSWELL[i] = cropParam.CropDevelopmentStages[i].DRYSWELL
		g.LUKRIT[i] = cropParam.CropDevelopmentStages[i].LUKRIT
		g.LAIFKT[i] = cropParam.CropDevelopmentStages[i].LAIFKT
		g.WGMAX[i] = cropParam.CropDevelopmentStages[i].WGMAX
		g.WGSMAX[i] = cropParam.CropDevelopmentStages[i].WGSMAX
		for L := 0; L < g.NRKOM; L++ {
			g.PRO[i][L] = cropParam.CropDevelopmentStages[i].PRO[L]
			g.DEAD[i][L] = cropParam.CropDevelopmentStages[i].DEAD[L]
		}
		l.kc[i] = cropParam.CropDevelopmentStages[i].Kc
		l.tendsum = l.tendsum + g.TSUM[i]
	}
}

// ReadCropParamClassic reads the crop parameters from an hermes crop file (classic format)
func ReadCropParamClassic(PARANAM string, l *CropSharedVars, g *GlobalVarsMain) {

	_, scanner, _ := Open(&FileDescriptior{FilePath: PARANAM, FileDescription: "crop file", UseFilePool: true})
	LineInut(scanner)
	LineInut(scanner)
	LineInut(scanner)
	LINE0 := LineInut(scanner)
	// Amax (C-Assimilation bei Lichtsättigung) bei Optimaltemperatur (kg CO2/ha leave/h)
	g.MAXAMAX = ValAsFloat(LINE0[65:], PARANAM, LINE0)
	LINE0b := LineInut(scanner)
	// C-Typ (C3 = 1/ C4 = 2) für Temperaturfunktion
	l.temptyp = int(ValAsInt(LINE0b[65:], PARANAM, LINE0b))
	LINE01 := LineInut(scanner)
	// Minimumtemperatur für Wachstum (°C)
	g.MINTMP = ValAsFloat(LINE01[65:], PARANAM, LINE01)
	LINE02 := LineInut(scanner)
	// Pflanzenspezifische effektive Durchwurzelungstiefe (dm)
	g.WUMAXPF = ValAsFloat(LINE02[65:], PARANAM, LINE02)
	LINE03 := LineInut(scanner)

	// selection root distribution function over depth (actually only 1 available)
	// replace no. of root function by: root depth increase in mm/C°....
	//variable: RTVELOC -> Parameter VELOC= RTVELOC/200
	RTVELOC := ValAsFloat(LINE03[65:], PARANAM, LINE03)
	g.VELOC = RTVELOC / 200
	//TODO: Default value to all files

	// Auswahl Wurzeltiefenfunktion (nur 1 verfügbar)
	//g.WUFKT = int(ValAsInt(LINE03[65:], PARANAM, LINE03))

	LINE04 := LineInut(scanner)
	// crop N-content function no. (critical and max. N-contents)....
	g.NGEFKT = int(ValAsInt(LINE04[65:], PARANAM, LINE04))

	// crop N-content function no. (critical and max. N-contents)....
	// simplify functions and reading two parameters for function 5: RGA and RGB and
	// if (S0 = no additional below ground organ and which additional organ to
	// those included in OBMAS (e.g. beet = S4 =Worg[3]) should be included in the N function
	// N-Gehaltsfunktion Nr.  a=4.90 b=0.45 below gr. organ org=S4 ..

	if g.NGEFKT == 5 {
		// TODO: setup default
		line04Token := strings.Fields(LINE04)
		for _, token := range line04Token {
			if strings.HasPrefix(token, "a=") {
				subToken := strings.Split(token, "=")
				g.RGA = ValAsFloat(subToken[1], PARANAM, LINE04)
			}
			if strings.HasPrefix(token, "b=") {
				subToken := strings.Split(token, "=")
				g.RGB = ValAsFloat(subToken[1], PARANAM, LINE04)
			}
			if strings.HasPrefix(token, "org=") {
				subToken := strings.Split(token, "=")
				g.SubOrgan = int(ValAsInt(subToken[1][1:], PARANAM, LINE04))
				if g.SubOrgan > 5 {
					log.Fatalf("Error: parsing crop organ! File: %s \n   Line: %s \n", PARANAM, LINE04)
					return
				}
			}
		}
	}

	LINE05 := LineInut(scanner)
	//above ground organs (numbers of compartiments increasing order)
	progip1Trimed := strings.TrimSpace(LINE05[65:])
	nrkob := len(progip1Trimed)
	l.AboveGroundOrgans = make([]int, nrkob)
	for i := 0; i < nrkob; i++ {
		komp := int(ValAsInt(progip1Trimed[i:i+1], "none", LINE05))
		l.AboveGroundOrgans[i] = komp
	}

	Line05b := LineInut(scanner)
	// organ no. for yield and fraction of organ.(organ 4 80% =4.80).
	g.YORGAN = int(ValAsInt(Line05b[65:66], PARANAM, Line05b))
	// fraction of organ.(organ 4 80% =4.80)
	g.YIFAK = ValAsFloat(Line05b[66:], PARANAM, Line05b)
	LINE06 := LineInut(scanner)
	LINE06b := LineInut(scanner)
	LINE1 := LineInut(scanner)
	//Anzahl Pflanzenkompertimente
	g.NRKOM = int(ValAsInt(LINE1[65:], PARANAM, LINE1))
	LineInut(scanner)
	LINE1bVal := LineInut(scanner)
	LINE1b := []rune(LINE1bVal)
	// Kennzeichnung Dauerkultur = "D"
	g.DAUERKULT = LINE1b[32] == 'D'
	// Kennzeichnung Leguminose  = "L"
	g.LEGUM = LINE1b[40] == 'L'
	if !g.DAUERKULT {
		ResetStages(g)
		g.PHYLLO, g.VERNTAGE = 0, 0
		for i := 0; i < 5; i++ {
			for i2 := 0; i2 < 10; i2++ {
				g.SUM[i2] = 0
				g.DEV[i2] = 0
				g.PRO[i2][i] = 0
				g.DEAD[i2][i] = 0
				g.TROOTSUM = 0
			}
		}
	} else {
		for i := 0; i < 5; i++ {
			for i2 := 2; i2 < 10; i2++ {
				g.SUM[i2] = 0
				g.PRO[i2][i] = 0
				g.DEAD[i2][i] = 0
				g.TROOTSUM = 0
			}
		}
	}
	if !(g.DAUERKULT && g.AKF.Num > 2 && g.FRUCHT[g.AKF.Index] == g.FRUCHT[g.AKF.Index-1]) {
		g.GEHOB = ValAsFloat(LINE06[65:], PARANAM, LINE06) / 100
		g.WUGEH = ValAsFloat(LINE06b[65:], PARANAM, LINE06b) / 100
	}
	LINE1C := LineInut(scanner)
	for i := 0; i < g.NRKOM; i++ {
		if !(g.DAUERKULT && g.AKF.Num > 2 && g.FRUCHT[g.AKF.Index] == g.FRUCHT[g.AKF.Index-1]) {
			g.WORG[i] = ValAsFloat(string(LINE1b[25+8*(i+1):30+8*(i+1)]), PARANAM, LINE1bVal)
		}
		//  Maintenancerate of Organ I
		g.MAIRT[i] = ValAsFloat(LINE1C[25+8*(i+1):30+8*(i+1)], PARANAM, LINE1C)
		g.WDORG[i] = 0
	}
	LINE1d := LineInut(scanner)
	// kc Faktor unbedeckter Boden
	l.kcini = ValAsFloat(LINE1d[65:], PARANAM, LINE1d)
	LINE2 := LineInut(scanner)
	// Anzahl der Entwicklungsstufen (max 10)
	l.NRENTW = int(ValAsInt(LINE2[65:], PARANAM, LINE2))
	l.tendsum = 0
	l.useBBCH = false
	for i := 0; i < l.NRENTW; i++ {
		developmentStageHeadline := LineInut(scanner)
		bbch := 0
		if len(developmentStageHeadline) > 65 {
			// if the line is longer than 65 characters, it may have a BBCH code
			val, err := TryValAsFloat(developmentStageHeadline[65:])
			if err == nil && val >= 0 && val < 100 {
				// if the last part of the line is a number between 00 and 99, it is a BBCH code
				bbch = int(val)
			}
		}
		l.ENDBBCH[i] = float64(bbch)
		l.useBBCH = l.useBBCH || l.ENDBBCH[i] > 0

		LINE4 := LineInut(scanner)
		// Temperatursumme Entwicklungsstufe I (°C days)
		g.TSUM[i] = ValAsFloat(LINE4[65:], PARANAM, LINE4)
		LINE5 := LineInut(scanner)
		// Basisitemperatur Entwicklungsstufe I (°C)
		g.BAS[i] = ValAsFloat(LINE5[65:], PARANAM, LINE5)
		LINE6 := LineInut(scanner)
		// Benötigte Anzahl Vernalisationstage Entwicklungsstufe I (Tage)
		g.VSCHWELL[i] = ValAsFloat(LINE6[65:], PARANAM, LINE6)
		LINE7 := LineInut(scanner)
		// Tageslängenbedarf Entwicklungsstufe I (h)
		g.DAYL[i] = ValAsFloat(LINE7[65:], PARANAM, LINE7)
		LINE7b := LineInut(scanner)
		// Basistageslänge Entwicklungsstufe I (h)
		g.DLBAS[i] = ValAsFloat(LINE7b[65:], PARANAM, LINE7b)
		LINE8 := LineInut(scanner)
		// Schwelle für Trockenstress (Ta/Tp) Entwicklungsstufe I (0-1)
		g.DRYSWELL[i] = ValAsFloat(LINE8[65:], PARANAM, LINE8)
		LINE8b := LineInut(scanner)
		// kritischer Luftporenanteil Entwicklungsstufe I (cm^3/cm^3)
		g.LUKRIT[i] = ValAsFloat(LINE8b[65:], PARANAM, LINE8b)
		LINE8c := LineInut(scanner)
		// SLA specific leave area (area per mass) (m2/m2/kg TM) in I
		g.LAIFKT[i] = ValAsFloat(LINE8c[65:], PARANAM, LINE8c)
		LINE8d := LineInut(scanner)
		// N-content root end of phase I
		g.WGMAX[i] = ValAsFloat(LINE8d[65:], PARANAM, LINE8d)
		LINE9 := LineInut(scanner)
		LINE9b := LineInut(scanner)
		for L := 0; L < g.NRKOM; L++ {
			// Partitioning at end of phase I(fraction)
			g.PRO[i][L] = ValAsFloat(LINE9[25+8*(L+1):30+8*(L+1)], PARANAM, LINE9)
			// death rate at end of phase I (coefficient)
			g.DEAD[i][L] = ValAsFloat(LINE9b[25+8*(L+1):30+8*(L+1)], PARANAM, LINE9b)
		}
		l.tendsum = l.tendsum + g.TSUM[i]
		LINE9c := LineInut(scanner)
		//kc factor for evapotranspiration at end of phase I
		l.kc[i] = ValAsFloat(LINE9c[65:], PARANAM, LINE9c)
	}

}

func ConvertCropParamClassicToYml(PARANAM string) (CropParam, error) {

	cropParam := CropParam{
		CropName:              "",
		ABBr:                  "",
		Variety:               "",
		MAXAMAX:               0,
		TempTyp:               0,
		MINTMP:                0,
		WUMAXPF:               0,
		VELOC:                 0,
		NGEFKT:                0,
		RGA:                   0,
		RGB:                   0,
		SubOrgan:              0,
		AboveGroundOrgans:     []int{},
		YORGAN:                0,
		YIFAK:                 0,
		INITCONCNBIOM:         0,
		INITCONCNROOT:         0,
		NRKOM:                 0,
		CompartimentNames:     []string{},
		DAUERKULT:             false,
		LEGUM:                 false,
		WORG:                  []float64{},
		MAIRT:                 []float64{},
		KcIni:                 0,
		SGEFKT:                1,
		SFunctExp:             -0.169,
		SCritContent:          0.37,
		NRENTW:                0,
		CropDevelopmentStages: []CropDevelopmentStage{},
	}

	_, scanner, err := Open(&FileDescriptior{FilePath: PARANAM, FileDescription: "crop file", UseFilePool: true})
	if err != nil {
		return cropParam, err
	}

	// get variaty and abbreviation from filename
	cropParam.ABBr = filepath.Ext(PARANAM)
	cropParam.ABBr = strings.TrimPrefix(cropParam.ABBr, ".")
	filename := filepath.Base(PARANAM)
	// remove extension
	filename = strings.TrimSuffix(filename, "."+cropParam.ABBr)
	parts := strings.Split(filename, "_")
	if len(parts) > 1 {
		cropParam.Variety = parts[1]
	}

	LineInut(scanner)
	cropName := LineInut(scanner)
	cropName = strings.TrimPrefix(cropName, "crop:")
	cropName = strings.TrimPrefix(cropName, "Frucht:")
	cropParam.CropName = strings.TrimSpace(cropName)
	LineInut(scanner)

	LINE0 := LineInut(scanner)
	cropParam.MAXAMAX = ValAsFloat(LINE0[65:], PARANAM, LINE0)

	LINE0b := LineInut(scanner)
	cropParam.TempTyp = int(ValAsInt(LINE0b[65:], PARANAM, LINE0b))

	LINE01 := LineInut(scanner)
	cropParam.MINTMP = ValAsFloat(LINE01[65:], PARANAM, LINE01)

	LINE02 := LineInut(scanner)
	cropParam.WUMAXPF = ValAsFloat(LINE02[65:], PARANAM, LINE02)

	LINE03 := LineInut(scanner)
	cropParam.VELOC = ValAsFloat(LINE03[65:], PARANAM, LINE03)

	LINE04 := LineInut(scanner)
	cropParam.NGEFKT = int(ValAsInt(LINE04[65:], PARANAM, LINE04))

	if cropParam.NGEFKT == 5 {
		line04Token := strings.Fields(LINE04)
		for _, token := range line04Token {
			if strings.HasPrefix(token, "a=") {
				subToken := strings.Split(token, "=")
				cropParam.RGA = ValAsFloat(subToken[1], PARANAM, LINE04)
			}
			if strings.HasPrefix(token, "b=") {
				subToken := strings.Split(token, "=")
				cropParam.RGB = ValAsFloat(subToken[1], PARANAM, LINE04)
			}
			if strings.HasPrefix(token, "org=") {
				subToken := strings.Split(token, "=")
				cropParam.SubOrgan = int(ValAsInt(subToken[1][1:], PARANAM, LINE04))
			}
		}
	}

	LINE05 := LineInut(scanner)
	progip1Trimed := strings.TrimSpace(LINE05[65:])
	nrkob := len(progip1Trimed)
	for i := 0; i < nrkob; i++ {
		komp := int(ValAsInt(progip1Trimed[i:i+1], "none", LINE05))
		cropParam.AboveGroundOrgans = append(cropParam.AboveGroundOrgans, komp)
	}

	Line05b := LineInut(scanner)
	cropParam.YORGAN = int(ValAsInt(Line05b[65:66], PARANAM, Line05b))
	cropParam.YIFAK = ValAsFloat(Line05b[66:], PARANAM, Line05b)

	LINE06 := LineInut(scanner)
	cropParam.INITCONCNBIOM = ValAsFloat(LINE06[65:], PARANAM, LINE06)

	LINE06b := LineInut(scanner)
	cropParam.INITCONCNROOT = ValAsFloat(LINE06b[65:], PARANAM, LINE06b)

	LINE1 := LineInut(scanner)
	cropParam.NRKOM = int(ValAsInt(LINE1[65:], PARANAM, LINE1))
	compNames := LineInut(scanner)
	cropParam.CompartimentNames = strings.Fields(compNames)
	// drop first element and last element if it is empty
	cropParam.CompartimentNames = cropParam.CompartimentNames[1 : cropParam.NRKOM+1]

	LINE1bVal := LineInut(scanner)
	LINE1b := []rune(LINE1bVal)
	cropParam.DAUERKULT = LINE1b[32] == 'D'
	cropParam.LEGUM = LINE1b[40] == 'L'

	LINE1C := LineInut(scanner)
	cropParam.WORG = make([]float64, cropParam.NRKOM)
	cropParam.MAIRT = make([]float64, cropParam.NRKOM)
	for i := 0; i < cropParam.NRKOM; i++ {
		cropParam.WORG[i] = ValAsFloat(string(LINE1b[25+8*(i+1):30+8*(i+1)]), PARANAM, LINE1bVal)
		//  Maintenancerate of Organ I
		cropParam.MAIRT[i] = ValAsFloat(LINE1C[25+8*(i+1):30+8*(i+1)], PARANAM, LINE1C)
	}
	LINE1d := LineInut(scanner)
	// kc Faktor unbedeckter Boden
	cropParam.KcIni = ValAsFloat(LINE1d[65:], PARANAM, LINE1d)
	LINE2 := LineInut(scanner)
	// Anzahl der Entwicklungsstufen (max 10)
	cropParam.NRENTW = int(ValAsInt(LINE2[65:], PARANAM, LINE2))

	for i := 0; i < cropParam.NRENTW; i++ {
		text := LineInut(scanner)
		bbch := 0
		if len(text) > 65 {
			// if the line is longer than 65 characters, it may have a BBCH code
			val, err := TryValAsFloat(text[65:])
			if err == nil {
				// if the last part of the line is a number, it is a BBCH code
				bbch = int(val)
			}
		}
		developmentStageToString := strings.TrimSpace(text)
		// trim all leading and trailing '-' characters
		developmentStageToString = strings.Trim(developmentStageToString, "- ")
		developmentStage := CropDevelopmentStage{
			DevelopmentStageName: developmentStageToString,
			ENDBBCH:              bbch,
		}
		developmentStage.WGSMAX = 0.001
		LINE4 := LineInut(scanner)
		// Temperatursumme Entwicklungsstufe I (°C days)
		developmentStage.TSUM = ValAsFloat(LINE4[65:], PARANAM, LINE4)
		LINE5 := LineInut(scanner)
		// Basisitemperatur Entwicklungsstufe I (°C)
		developmentStage.BAS = ValAsFloat(LINE5[65:], PARANAM, LINE5)
		LINE6 := LineInut(scanner)
		// Benötigte Anzahl Vernalisationstage Entwicklungsstufe I (Tage)
		developmentStage.VSCHWELL = ValAsFloat(LINE6[65:], PARANAM, LINE6)
		LINE7 := LineInut(scanner)
		// Tageslängenbedarf Entwicklungsstufe I (h)
		developmentStage.DAYL = ValAsFloat(LINE7[65:], PARANAM, LINE7)
		LINE7b := LineInut(scanner)
		// Basistageslänge Entwicklungsstufe I (h)
		developmentStage.DLBAS = ValAsFloat(LINE7b[65:], PARANAM, LINE7b)
		LINE8 := LineInut(scanner)
		// Schwelle für Trockenstress (Ta/Tp) Entwicklungsstufe I (0-1)
		developmentStage.DRYSWELL = ValAsFloat(LINE8[65:], PARANAM, LINE8)
		LINE8b := LineInut(scanner)
		// kritischer Luftporenanteil Entwicklungsstufe I (cm^3/cm^3)
		developmentStage.LUKRIT = ValAsFloat(LINE8b[65:], PARANAM, LINE8b)
		LINE8c := LineInut(scanner)
		// SLA specific leave area (area per mass) (m2/m2/kg TM) in I
		developmentStage.LAIFKT = ValAsFloat(LINE8c[65:], PARANAM, LINE8c)
		LINE8d := LineInut(scanner)
		// N-content root end of phase I
		developmentStage.WGMAX = ValAsFloat(LINE8d[65:], PARANAM, LINE8d)
		LINE9 := LineInut(scanner)
		LINE9b := LineInut(scanner)
		developmentStage.PRO = make([]float64, cropParam.NRKOM)
		developmentStage.DEAD = make([]float64, cropParam.NRKOM)
		for L := 0; L < cropParam.NRKOM; L++ {
			// Partitioning at end of phase I(fraction)
			developmentStage.PRO[L] = ValAsFloat(LINE9[25+8*(L+1):30+8*(L+1)], PARANAM, LINE9)
			// death rate at end of phase I (coefficient)
			developmentStage.DEAD[L] = ValAsFloat(LINE9b[25+8*(L+1):30+8*(L+1)], PARANAM, LINE9b)
		}
		LINE9c := LineInut(scanner)
		//kc factor for evapotranspiration at end of phase I
		developmentStage.Kc = ValAsFloat(LINE9c[65:], PARANAM, LINE9c)
		cropParam.CropDevelopmentStages = append(cropParam.CropDevelopmentStages, developmentStage)
	}
	return cropParam, nil
}
