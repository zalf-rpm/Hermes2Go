package hermes

import (
	"log"
	"strings"
)

// crop parameter IO

type CropParam struct {
	// CropParam is a struct to hold the crop parameters
}

// ReadCropParam reads the crop parameters from a yml file
func ReadCropParam(filename string) CropParam {
	return CropParam{}
}

// WriteCropParam writes the crop parameters to a yml file (with comments?)
func WriteCropParam(filename string, cropParam CropParam) {

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
