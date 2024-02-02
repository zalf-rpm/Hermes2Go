package hermes

import (
	"log"
	"os"
	"strings"

	yaml "gopkg.in/yaml.v3"
)

// crop parameter IO

type CropParam struct {
	// CropParam is a struct to hold the crop parameters
	CropName string `yaml:"CropName" comment:"name of the crop"`
	ABBr     string `yaml:"CropAbbreviation" comment:"crop no./ abbreviation"` // Abbreviation of the crop
	Variety  string `yaml:"Variety" comment:"variaty of the crop"`             // Variaty of the crop

	MAXAMAX           float64   `yaml:"MAXAMAX" comment:"AMAX Max. CO2 assimilation rate (kg CO2/ha leave/h)"`                               // AMAX Max. CO2 assimilation rate (kg CO2/ha leave/h)
	TempTyp           int       `yaml:"TempTyp" comment:"type of temperature dependency (C3 = 1/ C4 = 2)"`                                   // type of temperature dependency (C3 = 1/ C4 = 2)
	MINTMP            float64   `yaml:"MINTMP" comment:"minimum temperature crop growth (in C°)"`                                            // minimum temperature crop growth (in C°)
	WUMAXPF           float64   `yaml:"WUMAXPF" comment:"crop specific maximum effective rooting depth(dm)"`                                 // crop specific maximum effective rooting depth(dm)
	VELOC             float64   `yaml:"VELOC" comment:"root depth increase in mm/C°"`                                                        // root depth increase in mm/C°
	NGEFKT            int       `yaml:"NGEFKT" comment:"crop N-content function number for critical and max. N-contents"`                    // crop N-content function number for critical and max. N-contents
	RGA               float64   `yaml:"RGA,omitempty" comment:" RGA parameter for crop N-content function number 5"`                         // RGA parameter for crop N-content function number 5
	RGB               float64   `yaml:"RGB,omitempty" comment:"RGB parameter for crop N-content function number 5"`                          // RGB parameter for crop N-content function number 5
	SubOrgan          int       `yaml:"SubOrgan,omitempty" comment:"SubOrgan parameter for crop N-content function number"`                  // SubOrgan parameter for crop N-content function number 5
	AboveGroundOrgans []int     `yaml:"AboveGroundOrgans" comment:"list of above ground organs (numbers of compartiments increasing order)"` // SubOrgan parameter for crop N-content function number
	YORGAN            int       `yaml:"YORGAN" comment:"organ number for yield"`                                                             // organ number for yield
	YIFAK             float64   `yaml:"YIFAK" comment:"fraction of yield organ (90% = 0.90)"`                                                // fraction of yield organ (90% = 0.90)
	INITCONCNBIOM     float64   `yaml:"INITCONCNBIOM" comment:"start conzentration N in above ground biomass (% i. d.m.)"`                   // start conzentration N in above ground biomass (% i. d.m.)
	INITCONCNROOT     float64   `yaml:"INITCONCNROOT" comment:"start concentration N in roots (% i. d.m.)"`                                  // start concentration N in roots (% i. d.m.)
	NRKOM             int       `yaml:"NRKOM" comment:"Number of crop compartiments"`                                                        // Number of crop compartiments
	CompartimentNames []string  `yaml:"CompartimentNames" comment:"list of compartiment names"`                                              // list of compartiment names
	DAUERKULT         rune      `yaml:"DAUERKULT" comment:"Dauerkultur - Permaculture D / Non Permaculture 0"`                               // Dauerkultur - Permaculture D / Non Permaculture 0
	LEGUM             rune      `yaml:"LEGUM" comment:"Legume L / Non Legume 0"`                                                             // Legume L / Non Legume 0
	WORG              []float64 `yaml:"WORG" comment:"initial weight kg d.m./ha of organ I"`                                                 // initial weight kg d.m./ha of organ I
	MAIRT             []float64 `yaml:"MAIRT" comment:"maintainance rates of organ I"`                                                       // Maintainance rates of organ I
	KcIni             float64   `yaml:"KcIni" comment:"initial kc factor for evapotranspiration (uncovered soil)"`                           // initial kc factor for evapotranspiration (uncovered soil)

	NRENTW                int                    `yaml:"NRENTW" comment:"number of development phases(max 10)"`   // number of development phases(max 10)
	CropDevelopmentStages []CropDevelopmentStage `yaml:"CropDevelopmentStages" comment:"development stage/phase"` // development stage/phase
}
type CropDevelopmentStage struct {
	DevelopmentStageName string    `yaml:"DevelopmentStageName" comment:"name of the development stage/phase"` // name of the development stage/phase
	TSUM                 float64   `yaml:"TSUM" comment:"development phase temperatur sum (°C days)"`          // development phase temperatur sum (°C days)
	BAS                  float64   `yaml:"BAS" comment:"base temperature in phase (°C)"`                       // base temperature in phase (°C)
	VSCHWELL             float64   `yaml:"VSCHWELL" comment:"vernalisation requirements (days)"`               // vernalisation requirements (days)
	DAYL                 float64   `yaml:"DAYL" comment:"day length requirements (hours)"`                     // day length requirements (hours)
	DLBAS                float64   `yaml:"DLBAS" comment:"base day length in phase (hours)"`                   // base day length in phase (hours)
	DRYSWELL             float64   `yaml:"DRYSWELL" comment:"drought stress below ETA/ETP-quotient"`           // drought stress below ETA/ETP-quotient
	LUKRIT               float64   `yaml:"LUKRIT" comment:"critical aircontent in topsoil (cm^3/cm^3)"`        // critical aircontent in topsoil (cm^3/cm^3)
	LAIFKT               float64   `yaml:"LAIFKT" comment:"specific leave area (area per mass) (m2/m2/kg TM)"` // specific leave area (area per mass) (m2/m2/kg TM)
	WGMAX                float64   `yaml:"WGMAX" comment:"N-content root end at the of phase"`                 // N-content root end at the of phase
	PRO                  []float64 `yaml:"PRO" comment:"Partitioning at end of phase (fraction)"`              // Partitioning at end of phase (fraction)
	DEAD                 []float64 `yaml:"DEAD" comment:"death rate at end of phase (coefficient)"`            // death rate at end of phase (coefficient)
	Kc                   float64   `yaml:"Kc" comment:"kc factor for evapotranspiration at end of phase"`      // kc factor for evapotranspiration at end of phase
}

// ReadCropParam reads the crop parameters from a yml file
func ReadCropParam(filename string) CropParam {

	return CropParam{}
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
	l.temptyp = ValAsInt(LINE0b[65:], PARANAM, LINE0b)
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
	l.Progip1 = LINE05[65:]
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
	g.DAUERKULT = LINE1b[32]
	// Kennzeichnung Leguminose  = "L"
	g.LEGUM = LINE1b[40]
	if g.DAUERKULT != 'D' {
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
	if !(g.DAUERKULT == 'D' && g.AKF.Num > 2 && g.FRUCHT[g.AKF.Index] == g.FRUCHT[g.AKF.Index-1]) {
		g.GEHOB = ValAsFloat(LINE06[65:], PARANAM, LINE06) / 100
		g.WUGEH = ValAsFloat(LINE06b[65:], PARANAM, LINE06b) / 100
	}
	LINE1C := LineInut(scanner)
	for i := 0; i < g.NRKOM; i++ {
		if !(g.DAUERKULT == 'D' && g.AKF.Num > 2 && g.FRUCHT[g.AKF.Index] == g.FRUCHT[g.AKF.Index-1]) {
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
	iForLater := 0
	for i := 0; i < l.NRENTW; i++ {
		LineInut(scanner)
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
		iForLater++
	}

}
