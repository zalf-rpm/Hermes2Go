package hermes

import (
	"log"
	"math"
	"strings"
)

// CropSharedVars includes the shared variables for this module
type CropSharedVars struct {
	DGORG   [10]float64
	GORG    [10]float64
	FV      float64
	FP      float64
	NRENTW  int
	Progip1 string
	temptyp int64
	kc      [10]float64
	kcini   float64
	tendsum float64
	EZIEL   float64
	PROZIEL float64
	MANT    [5]float64
}

// PhytoOut calculates plant mass, creates output
func PhytoOut(g *GlobalVarsMain, l *CropSharedVars, hPath *HFilePath, zeit int, output *CropOutputVars) {
	// ! --------------------------------------- sub-Modul für Pflanzenentwicklung und Wachstum
	// !      ---------------------------------------------------------
	// !      --------- BERECHNUNG DER PHYTOMASSE (EXPLIZIT) ----------
	// !      ---------------------------------------------------------
	// ! Inputs:
	// ! AKF                       = aktuelle Frucht
	// ! SAAT (AKF)                = Aussaattermin aktuelle Frucht
	// ! Ernte(AKF)                = Erntetermin aktuelle Frucht
	// ! TEMP(TAG)                 = Tagesmitteltemperatur vom TAG (°C)
	// ! TSOIL(0,Z)                = Bodentemperatur am Anfang Zeitschritt in Schicht Z  (°C)
	// ! RAD(TAG)                  = PAR von TAG (Mjoule/m^2)
	// ! TRREL                     = Quotient Ta/Tp (aus Wassermodell)
	// ! DLP                       = photoperiodisch aktive Tageslänge (incl. 6? b?rgerl. Dämmerung) aus Wassermodel (h)
	// ! Variable:
	// ! INTWICK                   = Nr. Entwicklungsstadium
	// ! SUM(INTWICK)              = entwicklungswirksame Temperatursumme in Stadium INTWICK
	// ! WUMAS                     = Wurzeltrockenmasse (kg/ha)
	// ! OBMAS                     = oberirdische Trockenmasse (kg/ha)Gesamt N-Aufnahme der Pflanze (kg N/ha)
	// ! GEHOB                     = N-Gehalt in oberird. Biomasse (kg N/kg OBMAS)
	// ! LAI                       = Blattflächenindex
	// ! WORG(I)                   = TMasse von Organ I (kg/ha)
	// ! WDORG(I)                  = abgestorbene Masse von Organ I (kg/ha)
	// ! PHYLLO                    = kumulative entwicklungswirksame Temperatursumme (°C days)
	// ! WULAEN                    = Gesamtwurzellänge (cm/cm^2)
	// ! WUDICH(z)                 = Wurzellängendichte in Tiefenkompartiment z (cm Wurzel/cm^3 Boden)
	// ! REDUK                     = Stickstoffstressfaktor (0-1)
	// ! Pflanzenparameter siehe Einleseliste unten

	//! ------------------------- Einlesen der Parameter für neue Frucht bei deren Aussaat ------------------
	if zeit == g.SAAT[g.AKF.Index] {
		output.SowDate = g.Kalender(zeit)
		output.SowDOY = g.TAG.Index + 1
		g.managementConfig.WriteManagementEvent(NewManagementEvent(Sowing, zeit, make(map[string]interface{}), g))

		PARANAM := hPath.GetParanam(g.CropTypeToString(g.FRUCHT[g.AKF.Index], false), g.CVARIETY[g.AKF.Index])
		g.TRRELSUM = 0
		g.REDUKSUM = 0
		g.ETAG = 0
		g.TRAG = 0
		g.PERG = 0
		g.NLEAG = 0
		g.DRYD1 = 0
		g.DRYD2 = 0
		g.SWCA1 = 0
		g.SWCA2 = 0
		g.SWCM1 = 0
		g.SWCM2 = 0
		g.LAIMAX = 0

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
		if g.DAUERKULT == 'D' && g.AKF.Num > 2 && g.FRUCHT[g.AKF.Index] == g.FRUCHT[g.AKF.Index-1] { // && g.AKF.Num > 2
			g.INTWICK.SetByIndex(1)
			g.SUM[0] = g.TSUM[0] + 1
			g.SUM[1] = math.Min(g.PHYLLO, g.TSUM[1]*0.5)
			g.PHYLLO = math.Min(g.PHYLLO, g.TSUM[1]*0.5)
			// ! ------------- no vernalisation required after first winter inserted on 30.11.2017 ----------
			g.OBMAS = 0
		} else {
			g.INTWICK.SetByIndex(0)
			g.OBMAS = 0
		}
		progip1Trimed := strings.TrimSpace(l.Progip1)
		nrkob := len(progip1Trimed)
		for i := 0; i < nrkob; i++ {
			komp := ValAsInt(progip1Trimed[i:i+1], "none", l.Progip1)
			//  Berechnung der oberirdischen Masse (kg/ha)
			g.OBMAS = g.OBMAS + g.WORG[komp-1]
		}
		g.ASPOO = 0 // Assimilationspool reset for new crop
		// Berechnung Anfangs-LAI aus Blattgewicht und SLA
		g.LAI = g.WORG[1] * g.LAIFKT[g.INTWICK.Index]
		if g.LAI < 0 {
			g.LAI = 0
		}

		g.WUMAS = g.WORG[0]
		// Berechnung Wurzellänge aus Wurzelmasse
		//WULAEN = (g.WUMAS * 100000 * 100.0 / 7.0) / (math.Pow(0.015, 2.0) * math.Pi)
		if g.FRUCHT[g.AKF.Index] == ZR || g.FRUCHT[g.AKF.Index] == K {
			g.PESUM = (g.OBMAS*g.GEHOB + (g.WUMAS+g.WORG[3])*g.WUGEH)
		} else {
			g.PESUM = (g.OBMAS*g.GEHOB + g.WUMAS*g.WUGEH)
		}
		for i := 0; i < g.N; i++ {
			g.WUDICH[i] = 0
		}
	}
	// ! ********************************* end of reading crop parameter file *************************************************

	if g.INTWICK.Num == 1 {
		if g.TEMP[g.TAG.Index] > g.BAS[g.INTWICK.Index] {
			if g.WG[0][0] > 0.3*(g.W[0]-g.WMIN[0])+g.WMIN[0] {
				g.SUM[g.INTWICK.Index] = g.SUM[g.INTWICK.Index] + (g.TEMP[g.TAG.Index]-g.BAS[g.INTWICK.Index])*g.DT.Num
			} else {
				g.SUM[g.INTWICK.Index] = g.SUM[g.INTWICK.Index] + (g.TEMP[g.TAG.Index]-g.BAS[g.INTWICK.Index])*(g.WG[0][0]/(0.3*(g.W[0]-g.WMIN[0])+g.WMIN[0]))*g.DT.Num
			}
		}
		g.FKC = l.kcini + (l.kc[g.INTWICK.Index]-l.kcini)*g.SUM[0]/g.TSUM[0]
	}
	var DTGESN float64
	//var DTGESS float64
	var DTGESS2 float64
	var WUMALT float64
	var OBALT float64
	var GEHALT float64
	if g.SUM[0] >= g.TSUM[0] {
		if g.SUM[g.INTWICK.Index] >= g.TSUM[g.INTWICK.Index] {
			if int(g.INTWICK.Num) < l.NRENTW {
				//! Improvement 1: Consideration of oversum for next stage
				g.SUM[g.INTWICK.Index+1] = g.SUM[g.INTWICK.Index] - g.TSUM[g.INTWICK.Index]
				g.INTWICK.Inc()
				g.DEV[g.INTWICK.Index] = g.TAG.Index + 1
				g.DevStateDate[g.INTWICK.Index] = g.Kalender(zeit)

				SWC := 0.0
				SWC1 := 0.0
				if g.INTWICK.Num == 5 {
					for i := 1; i <= 15; i++ {
						if i < 4 {
							SWC1 = SWC1 + g.WG[0][i-1]*100
						}
						SWC = SWC + g.WG[0][i-1]*100
					}
					g.SWCA1 = SWC1
					g.SWCA2 = SWC
				} else if g.INTWICK.Num == 6 {
					for i := 1; i <= 15; i++ {
						if i < 4 {
							SWC1 = SWC1 + g.WG[0][i-1]*100
						}
						SWC = SWC + g.WG[0][i-1]*100
					}
					g.SWCM1 = SWC1
					g.SWCM2 = SWC
				}
			}
		}
		//! +++++++++++++++++++++++++++++++++++++ Automatic harvest ++++++++++++++++++++++++++++++++++++++++
		if g.ERNTE[g.AKF.Index] == 0 {
			if g.INTWICK.Index+1 == l.NRENTW {
				if g.SUM[g.INTWICK.Index] > 0.6*g.TSUM[g.INTWICK.Index] {
					NFK1 := (g.WG[0][0] + g.REGEN[g.TAG.Index]/g.DZ.Num - g.WMIN[0]) / (g.WNOR[0] - g.WMIN[0]) * 100
					if NFK1 <= g.MAXHMOI[g.AKF.Index] && NFK1 >= g.MINHMOI[g.AKF.Index] {
						if g.TAG.Num > 3 {
							if g.REGEN[g.TAG.Index]+g.REGEN[g.TAG.Index-1]+g.REGEN[g.TAG.Index-2]+g.REGEN[g.TAG.Index-3] <= g.RAINLIM[g.AKF.Index] && g.REGEN[g.TAG.Index] <= g.RAINACT[g.AKF.Index] {
								g.ERNTE[g.AKF.Index] = zeit
								g.ERNTE2[g.AKF.Index] = g.ERNTE[g.AKF.Index]
								if g.SAAT[g.AKF.Index+1] > 0 && g.SAAT[g.AKF.Index+1] < zeit {
									g.SAAT[g.AKF.Index+1] = zeit + 4
									g.SAAT2[g.AKF.Index+1] = zeit + 4
								}
							}
						}
					}
				}
			}
			if zeit == g.ERNTE2[g.AKF.Index]-1 && g.ERNTE[g.AKF.Index] == 0 {
				g.ERNTE[g.AKF.Index] = zeit + 1
			}
		}
		//! ----------------------- End of automatic harvest --------------------------
		if g.LAI <= 0 {
			g.LAI = 0.001
		}
		//---------------------------------------------

		// Aufruf Modul für Stahlungsinterception und Photosynthese nach Penning de Vries 1982 ----
		_, DLP, GPHOT, MAINT := radia(g, l, zeit)
		// ----------------------------------------------------------------------------------------
		//  Netto-Assimilation kg C/ha
		GTW := GPHOT + g.ASPOO

		// Assimilationspool
		g.ASPOO = 0
		// Beruecksichtigung des Proteinziels im Korn für Duengerbedarfsanalyse
		if g.FRUCHT[g.AKF.Index] == WW {
			l.EZIEL = 95
			l.PROZIEL = 0.0224
		} else if g.FRUCHT[g.AKF.Index] == WR {
			l.PROZIEL = 0.0175
			l.EZIEL = 80
		} else if g.FRUCHT[g.AKF.Index] == WG {
			l.PROZIEL = 0.0192
			l.EZIEL = 80
		}
		if g.VSCHWELL[g.INTWICK.Index] == 0 {
			l.FV = 1
		} else {
			//----------- Aufruf Sub-Modul Vernalisation -----------
			vern(l, g)
		}
		if g.DAYL[g.INTWICK.Index] > 0 {
			// FP = Tageslängenkorrekturfaktor für SUM(Intwick)
			l.FP = (DLP - g.DLBAS[g.INTWICK.Index]) / (g.DAYL[g.INTWICK.Index] - g.DLBAS[g.INTWICK.Index])
		} else if g.DAYL[g.INTWICK.Index] < 0 {
			if DLP <= math.Abs(g.DAYL[g.INTWICK.Index]) {
				l.FP = 1
			} else {
				Daycrit := math.Abs(g.DAYL[g.INTWICK.Index])
				Maxdaylength := math.Abs(g.DLBAS[g.INTWICK.Index])
				l.FP = (DLP - Maxdaylength) / (Daycrit - Maxdaylength)
			}
		} else {
			l.FP = 1
		}
		if l.FP > 1 {
			l.FP = 1
		}
		if l.FP < 0 {
			l.FP = 0
		}

		// Entwicklungsbeschleunigung durch Wasser bzw Stickstoffstress (7.11.07 aus Ottawa)
		var Nprog float64
		if g.FRUCHT[g.AKF.Index] == ZR || g.FRUCHT[g.AKF.Index] == SM {
			// keine Beschleunigung durch N-Stress
			Nprog = 1
		} else {
			Nprog = 1 + math.Pow((1-g.REDUK), 2)
		}
		var WPROG float64
		// keine Entwicklungsbeschleunigung, wenn Ta durch Luftmangel indiziert
		if g.TRREL < g.DRYSWELL[g.INTWICK.Index] {
			if g.LURED < 1 {
				WPROG = 1
			} else {
				WPROG = 1 + 0.2*math.Pow((1-g.TRREL), 2)
			}
		} else {
			WPROG = 1
		}
		devprog := math.Max(Nprog, WPROG)
		// ---------- Berechnung der entwicklungswirksamen Temperatursummen (SUM(INTWICK) ----------
		if g.TEMP[g.TAG.Index] >= g.BAS[g.INTWICK.Index] {
			g.SUM[g.INTWICK.Index] = g.SUM[g.INTWICK.Index] + (g.TEMP[g.TAG.Index]-g.BAS[g.INTWICK.Index])*l.FV*l.FP*devprog*g.DT.Num
			g.PHYLLO = g.PHYLLO + (g.TEMP[g.TAG.Index]-g.BAS[g.INTWICK.Index])*l.FV*l.FP*devprog*g.DT.Num
			CalulateDevelopmentStages(zeit, l.FV, l.FP, g)
		}
		// -- Interpolsation kc Faktor aus Entwicklungsfortschritt --
		if g.INTWICK.Num < 2 {
			relint := g.SUM[g.INTWICK.Index] / g.TSUM[g.INTWICK.Index]
			if relint > 1 {
				relint = 1
			}
			g.FKC = l.kcini + (l.kc[g.INTWICK.Index]-l.kcini)*relint
		} else {
			relint := g.SUM[g.INTWICK.Index] / g.TSUM[g.INTWICK.Index]
			if relint > 1 {
				relint = 1
			}
			g.FKC = l.kc[g.INTWICK.Index-1] + (l.kc[g.INTWICK.Index]-l.kc[g.INTWICK.Index-1])*relint
		}
		// +++++++++++++++++++  N-Gehaltsfunktionen  +++++++++++++++++++++++++
		// Funtionen für GEHMAX und GEHMIN in Abhängigkeit der Entwicklung (PHYLLO) oder der oberird. Biomasse (OBMAS)
		// GEHMAX   = maximal möglicher N-Gehalt (Treiber für N-Aufnahme)(kg N/kg Biomasse)
		// GEHMIN   = kritischer N-Gehalt der Biomasse (Beginn N-Stress) (kg N/kg Biomasse)
		if g.NGEFKT == 1 {
			if g.PHYLLO < 200 {
				g.GEHMIN = .0415
				g.GEHMAX = .06
			} else {
				if g.FRUCHT[g.AKF.Index] == WR || g.FRUCHT[g.AKF.Index] == SG {
					g.GEHMIN = 5.1 * math.Exp(-.00165*g.PHYLLO) / 100
					g.GEHMAX = 8.0 * math.Exp(-.0017*g.PHYLLO) / 100
				} else {
					g.GEHMIN = 5.5 * math.Exp(-.0014*g.PHYLLO) / 100
					g.GEHMAX = 8.1 * math.Exp(-.00147*g.PHYLLO) / 100
				}
			}
		} else if g.NGEFKT == 2 {
			if g.PHYLLO < 263 {
				g.GEHMIN = 0.035
			} else {
				g.GEHMIN = 0.035 - 0.024645*math.Pow((1-math.Exp(-(g.PHYLLO-152.30391*math.Log(1-math.Sqrt2/2)-438.63545)/152.30391)), 2)
			}
			if g.PHYLLO < 142 {
				g.GEHMAX = 0.049
			} else {
				g.GEHMAX = 0.049 - 0.037883841*math.Pow((1-math.Exp(-(g.PHYLLO-201.50354*math.Log(1-math.Sqrt2/2)-385.8318)/201.50354)), 2)
			}
		} else if g.NGEFKT == 3 {
			if g.OBMAS < 1000 {
				g.GEHMAX = 0.06
				g.GEHMIN = 0.045
			} else {
				g.GEHMAX = 0.06 * math.Pow((g.OBMAS/1000), (-0.25))
				g.GEHMIN = 0.045 * math.Pow((g.OBMAS/1000), (-0.25))
			}
		} else if g.NGEFKT == 4 {
			if (g.OBMAS + g.WORG[3]) < 1000 {
				g.GEHMAX = 0.06
				g.GEHMIN = 0.045
			} else {
				g.GEHMAX = 0.0285 + 0.0403*math.Exp(-0.26*(g.OBMAS+g.WORG[3])/1000)
				g.GEHMIN = 0.0135 + 0.0403*math.Exp(-0.26*(g.OBMAS+g.WORG[3])/1000)
			}
			// } else if g.NGEFKT == 5 {

		} else if g.NGEFKT == 5 {
			// new variables RGA and RGB read from line N-content function
			org := 0.0
			if g.SubOrgan > 0 {
				org = g.WORG[g.SubOrgan-1]
			}
			if (g.OBMAS + org) < 1100 {
				g.GEHMAX = 0.06
				g.GEHMIN = g.RGA
			} else {
				g.GEHMAX = 0.06 * math.Pow(((g.OBMAS+org)/1000), g.RGB)
				g.GEHMIN = g.RGA * math.Pow(((g.OBMAS+org)/1000), g.RGB)
			}

			// if (g.OBMAS + g.WORG[3]) < 1100 {
			// 	g.GEHMAX = 0.06
			// 	g.GEHMIN = 0.045
			// } else {
			// 	g.GEHMAX = 0.06 * math.Pow(((g.OBMAS+g.WORG[3])/1000), 0.5294)
			// 	g.GEHMIN = 0.046694 * math.Pow(((g.OBMAS+g.WORG[3])/1000), 0.5294)
			// }
		} else if g.NGEFKT == 6 {
			if g.PHYLLO < 400 {
				g.GEHMIN = .0415
				g.GEHMAX = .06
			} else {
				g.GEHMIN = 5.5 * math.Exp(-.0007*g.PHYLLO) / 100
				g.GEHMAX = 8.1 * math.Exp(-.0007*g.PHYLLO) / 100
			}
		} else if g.NGEFKT == 7 {
			if g.OBMAS < 1000 {
				g.GEHMAX = 0.0615
				g.GEHMIN = 0.0448
			} else {
				g.GEHMAX = 0.0615 * math.Pow((g.OBMAS/1000), (-0.25))
				g.GEHMIN = 0.0448 * math.Pow((g.OBMAS/1000), (-0.25))
			}
		} else if g.NGEFKT == 8 {
			if g.PHYLLO < 200*l.tendsum/1260 {
				g.GEHMIN = .0415
				g.GEHMAX = .06
			} else {
				// Korrekturfaktor für Entwicklungsfunktion bei sortenspez. Temperatursummen
				dvkor := 1 / ((l.tendsum - 200) / (1260 - 200))
				if g.FRUCHT[g.AKF.Index] == WR || g.FRUCHT[g.AKF.Index] == SG {
					g.GEHMIN = 5.1 * math.Exp(-.00165*dvkor*g.PHYLLO) / 100
					g.GEHMAX = 8.0 * math.Exp(-.0017*dvkor*g.PHYLLO) / 100
				} else {
					g.GEHMIN = 5.5 * math.Exp(-.0014*dvkor*g.PHYLLO) / 100
					g.GEHMAX = 8.1 * math.Exp(-.00147*dvkor*g.PHYLLO) / 100
				}
			}
		} else if g.NGEFKT == 9 {
			if (g.OBMAS + g.WORG[3]) < 1000 {
				g.GEHMAX = 0.06
				g.GEHMIN = 0.045
			} else {
				g.GEHMAX = 0.0285 + 0.0403*math.Exp(-0.26*g.OBMAS/1000)
				g.GEHMIN = 0.0135 + 0.0403*math.Exp(-0.26*g.OBMAS/1000)
			}
		}

		if g.Sulfonie {
			// calc biomass
			BM := g.OBMAS / 1000
			// for _, bioMassPart := range g.WORG {
			// 	BM += bioMassPart
			// }
			SC := g.CRITSGEHALT[g.FRUCHT[g.AKF.Index]]
			exp := g.CRITSEXP[g.FRUCHT[g.AKF.Index]]
			// wheat, maize, soybean
			if g.SGEFKT[g.FRUCHT[g.AKF.Index]] == 1 {
				if BM > 1.0 {
					SC = SC * math.Pow((BM), exp)
				}
				// oilseed rape
			} else if g.SGEFKT[g.FRUCHT[g.AKF.Index]] == 2 {
				if BM > 1.0 {
					SC = SC * math.Exp(exp*BM)
				}
			}

			// SGEHMAX   = maximal möglicher S-Gehalt (Treiber für S-Aufnahme)(kg S/kg Biomasse)
			// SGEHMIN   = kritischer S-Gehalt der Biomasse (Beginn S-Stress) (kg S/kg Biomasse)
			g.SGEHMAX = SC * 1.3
			g.SGEHMIN = SC

		}

		// -------------------------------------------------------
		//              Trockenmassenproduktion
		// -------------------------------------------------------
		GEHALT = g.GEHOB
		if g.INTWICK.Num > 0 {
			if g.GEHOB < g.GEHMIN {
				var MININ float64
				if g.NGEFKT == 1 {
					MININ = 0.005
				} else {
					MININ = 0.004
				}
				if g.GEHOB <= MININ {
					g.REDUK = 0.0
				} else {
					AUX := (g.GEHOB - MININ) / (g.GEHMIN - MININ)
					g.REDUK = math.Pow((1 - math.Exp(1+1/(AUX-1))), 2)
				}
			} else {
				g.REDUK = 1.
			}
			if g.SGEHOB < g.SGEHMIN {
				AUX := g.SGEHOB / g.GEHMIN
				g.SREDUK = math.Pow((1 - math.Exp(1+1/(AUX-1))), 2)
			}
			g.REDUKSUM = g.REDUKSUM + g.REDUK
			g.SREDUKSUM = g.SREDUKSUM + g.SREDUK

			g.TRRELSUM = g.TRRELSUM + g.TRREL
			// ************************ DAYS WITH ETA/ETP < 0.4 UNTIL ANTHESIS UND ANTHESIS TO MATURITY
			if g.ETREL < 0.4 {
				if g.INTWICK.Num < 5 {
					g.DRYD1 = g.DRYD1 + 1
				} else if g.INTWICK.Num > 4 && g.INTWICK.Num < 6 {
					g.DRYD2 = g.DRYD2 + 1
				}
			}

			for i := 0; i < g.NRKOM; i++ {
				if g.SUM[g.INTWICK.Index]/g.TSUM[g.INTWICK.Index] > 1 {
					l.GORG[i] = 0
				} else {
					// Berücksichtigung der unterschiedlichen organspezifischen Maintenanceraten
					l.GORG[i] = GTW*0.7*(g.PRO[g.INTWICK.Index-1][i]+(g.PRO[g.INTWICK.Index][i]-g.PRO[g.INTWICK.Index-1][i])*g.SUM[g.INTWICK.Index]/g.TSUM[g.INTWICK.Index])*g.REDUK - (MAINT * l.MANT[i] * 0.7)
					l.DGORG[i] = g.WORG[i] * (g.DEAD[g.INTWICK.Index-1][i] + (g.DEAD[g.INTWICK.Index][i]-g.DEAD[g.INTWICK.Index-1][i])*(math.Min(1, g.SUM[g.INTWICK.Index]/g.TSUM[g.INTWICK.Index])))
				}
				if i+1 < 4 {
					if g.WORG[i]+(l.GORG[i]-l.DGORG[i])*g.DT.Num > 0.0000000000001 { // almost 0
						g.WORG[i] = g.WORG[i] + l.GORG[i]*g.DT.Num - l.DGORG[i]*g.DT.Num
					} else {
						l.DGORG[i] = g.WORG[i]/g.DT.Num + l.GORG[i]
						g.WORG[i] = 0.1
					}
				} else {
					if int(g.INTWICK.Num) < l.NRENTW {
						g.WORG[i] = g.WORG[i] + l.GORG[i]*g.DT.Num - l.DGORG[i]*g.DT.Num + 0.3*(l.DGORG[i-1]*g.DT.Num+l.DGORG[i-2]*g.DT.Num+l.DGORG[i-3]*g.DT.Num)
					} else {
						g.WORG[i] = g.WORG[i] + l.GORG[i]*g.DT.Num - l.DGORG[i]*g.DT.Num
					}
					if g.WORG[i] < 0 {
						l.DGORG[i] = l.DGORG[i] + g.WORG[i]/g.DT.Num
						g.WORG[i] = 0
					}
				}
				if i+1 == 2 {
					laialt := g.LAI
					g.LAI = g.LAI + l.GORG[i]*(g.LAIFKT[g.INTWICK.Index-1]+(g.SUM[g.INTWICK.Index]/g.TSUM[g.INTWICK.Index]*(g.LAIFKT[g.INTWICK.Index]-g.LAIFKT[g.INTWICK.Index-1])))*g.DT.Num - l.DGORG[i]*g.LAIFKT[0]*g.DT.Num
					// when LAI goes negativ, capped at 0 to prevent side effects
					if g.LAI < 0 {
						g.LAI = 0
					}
					if g.LAI > laialt {
						g.LAIMAX = g.LAI
					}
				}
				if i+1 == 1 {
					//  Fluss abgestorbener Wurzeln in org. Pools (inaktiviert)
				} else if i+1 < 4 {
					g.NFOS[0] = g.NFOS[0] + 0.7*0.8*l.DGORG[i]*GEHALT*g.DT.Num
					g.NAOS[0] = g.NAOS[0] + 0.7*0.2*l.DGORG[i]*GEHALT*g.DT.Num
					g.PESUM = g.PESUM - 0.7*l.DGORG[i]*GEHALT*g.DT.Num
				}
				g.WDORG[i] = g.WDORG[i] + l.DGORG[i]*g.DT.Num
				if (g.WORG[i] - g.WDORG[i]) <= 0 {
					g.WDORG[i] = g.WORG[i] - 0.001
				}

			}
			g.ASPOO = g.ASPOO + GTW*(1.-g.REDUK)
			if g.FRUCHT[g.AKF.Index] == ZR || g.FRUCHT[g.AKF.Index] == K {
				OBALT = g.OBMAS + g.WORG[3]
			} else {
				OBALT = g.OBMAS
			}
			// ! Definieren Oberirdische Masse
			g.OBMAS = 0
			progip1Trimed := strings.TrimSpace(l.Progip1)
			nrkob := len(progip1Trimed)
			for i := 0; i < nrkob; i++ {
				komp := ValAsInt(progip1Trimed[i:i+1], "none", l.Progip1)
				// Berechnung der oberirdischen Masse (kg/ha)
				g.OBMAS = g.OBMAS + g.WORG[komp-1]
			}
			WUMALT = g.WUMAS
			g.WUMAS = g.WORG[0]
			// ++++++++++++++++ Einschub Gras automatischer Wiederaustrieb ++++++++++++++++++++++
			if g.DAUERKULT == 'D' && g.INTWICK.Num > 4 {
				if g.OBMAS <= OBALT && g.WORG[1] < 100 {
					g.INTWICK.SetByIndex(0)
					g.SUM[0] = 5
					for i := 1; i <= l.NRENTW; i++ {
						g.SUM[i] = 0
					}
					g.PHYLLO = 0
					g.NAOS[0] = g.NAOS[0] + (g.WORG[2]+g.WORG[3]*g.GEHOB)*g.DT.Num
					g.PESUM = g.PESUM - ((g.WORG[2] + g.WORG[3]) * g.GEHOB)
					for i := 0; i < g.NRKOM; i++ {
						if !(i+1 < 3) {
							g.WORG[i] = 0
						}
						g.WDORG[i] = 0
					}
					g.OBMAS = g.WORG[1]
					OBALT = g.WORG[1]
				}
			}
			if g.FRUCHT[g.AKF.Index] == ZR || g.FRUCHT[g.AKF.Index] == K {
				DTGESN = (g.GEHMAX*g.OBMAS + (g.WUMAS+g.WORG[3])*g.WGMAX[g.INTWICK.Index] - g.PESUM) * g.DT.Num
			} else {
				DTGESN = (g.GEHMAX*g.OBMAS + g.WUMAS*g.WGMAX[g.INTWICK.Index] - g.PESUM) * g.DT.Num
			}

			// NMAX := 2.5
			// // !*******************  S-Aufnahmefunktion  ********************************
			// // LET SUP = Nmax * 10^(-ZF * (log10(Tempsum/Warmsum))^2)
			// SUP := NMAX * math.Pow(10, -g.ZF[g.FRUCHT[g.AKF.Index]]) * math.Pow(math.Log10(g.PHYLLO+g.SUM[0]/g.TSUM[g.INTWICK.Index]), 2)
			// // !*************************************************************************
			// // LET DTGESS = (SUP - PESUMS)*DT
			// DTGESS = (SUP - g.PESUMS) * g.DT.Num
			// // IF DTGESS > 1.5*DT THEN LET DTGESS = 1.5*DT
			// if DTGESS > 1.5*g.DT.Num {
			// 	DTGESS = 1.5 * g.DT.Num
			// }
			// if DTGESS < 0 {
			// 	DTGESS = 0.0
			// }
			// TODO: S-uptake depending on crop parameters
			// !*******************  S-Aufnahmefunktion  ********************************
			//WGSMax := g.WGMAX[g.INTWICK.Index] / g.SNRatio[g.FRUCHT[g.AKF.Index]]
			//WGSMax := g.WGMAX[g.INTWICK.Index] * g.SWura[g.FRUCHT[g.AKF.Index]]
			WGSMax := g.WUGEH * g.SWura[g.FRUCHT[g.AKF.Index]]
			if g.FRUCHT[g.AKF.Index] == ZR || g.FRUCHT[g.AKF.Index] == K {
				DTGESS2 = (g.SGEHMAX*g.OBMAS + (g.WUMAS+g.WORG[3])*WGSMax - g.PESUMS) * g.DT.Num
			} else {
				DTGESS2 = (g.SGEHMAX*g.OBMAS + g.WUMAS*WGSMax - g.PESUMS) * g.DT.Num
			}
			if DTGESS2 < 0 {
				DTGESS2 = 0.0
			}
		}
	}
	if zeit == g.ERNTE2[g.AKF.Index]-1 && g.ERNTE[g.AKF.Index] == 0 {
		g.ERNTE[g.AKF.Index] = zeit + 1
		if g.SAAT[g.AKF.Index+1] > 0 && g.SAAT[g.AKF.Index+1] < zeit {
			g.SAAT[g.AKF.Index+1] = g.ERNTE[g.AKF.Index] + 4
			g.SAAT2[g.AKF.Index+1] = g.ERNTE[g.AKF.Index] + 4
		}
	}

	if DTGESN > 6*g.DT.Num {
		DTGESN = 6 * g.DT.Num
	}
	if DTGESN < 0 {
		DTGESN = 0.0
	}
	var WUMM float64
	if g.WUMAS < WUMALT {
		WUMM = (WUMALT - g.WUMAS) * g.WUGEH
	} else {
		WUMM = 0
	}
	//---------------------------------------------------------------
	//------ Berechnung der Wurzeldichte (-laenge/Dichte Boden) -----
	//---------------------------------------------------------------
	WURM := math.Round(float64(g.WURZMAX) * (g.WUMAXPF / 11.))
	if WURM > float64(g.N) {
		WURM = float64(g.N)
	}
	if WURM < 1 {
		WURM = 1
	}
	// new Qrez TODO: use new root funtion
	Qrez, potentialRootingDepth, _ := root(g.VELOC, g.PHYLLO+g.SUM[0], g.DZ.Num)
	g.POTROOTINGDEPTH = potentialRootingDepth

	//var Qrez float64
	// if g.FRUCHT[g.AKF.Index] == "ORF" || g.FRUCHT[g.AKF.Index] == "ORH" || g.FRUCHT[g.AKF.Index] == "WRA" || g.FRUCHT[g.AKF.Index] == "ZR " {
	// 	Qrez = math.Pow((0.081476 + math.Exp(-.004*(g.PHYLLO+g.SUM[0]+185.))), 1.8)
	// } else if g.FRUCHT[g.AKF.Index] == "SM " || g.FRUCHT[g.AKF.Index] == "K  " {
	// 	Qrez = math.Pow((0.081476 + math.Exp(-.0035*(g.PHYLLO+g.SUM[0]+211.))), 1.8)
	// } else if g.FRUCHT[g.AKF.Index] == "GR " && g.AKF.Num > 2 {
	// 	Qrez = math.Pow((0.081476 + math.Exp(-.002787*(math.Max(g.PHYLLO+g.SUM[0], 1500)))), 1.8)
	// } else if g.FRUCHT[g.AKF.Index] == "AA " && g.AKF.Num > 2 {
	// 	Qrez = math.Pow((0.081476 + math.Exp(-.002787*(math.Max(g.PHYLLO+g.SUM[0], 1500)))), 1.8)
	// } else if g.FRUCHT[g.AKF.Index] == "CLU" && g.AKF.Num > 2 {
	// 	Qrez = math.Pow((0.081476 + math.Exp(-.002787*(math.Max(g.PHYLLO+g.SUM[0], 1500)))), 1.8)
	// } else {
	// 	Qrez = math.Pow((0.081476 + math.Exp(-.002787*(g.PHYLLO+g.SUM[0]+265.))), 1.8)
	// }
	if Qrez > .35 {
		Qrez = .35
	}
	// make sure rooting depth is not deeper than WURZMAX in soil
	if Qrez < 4.5/(WURM*g.DZ.Num) {
		Qrez = 4.5 / (WURM * g.DZ.Num)
	}

	// root layer depth. adapted to WURZMAX
	g.WURZ = int(4.5 / Qrez / g.DZ.Num)

	// assumption: root radius decreases with depth: radius (cm)  RRAD(I) =  .02 - I*.001,
	// WRAD root radius
	WRAD := make([]float64, g.WURZ)
	for i := 1; i <= g.WURZ; i++ {
		if g.FRUCHT[g.AKF.Index] == ZR || g.FRUCHT[g.AKF.Index] == K {
			WRAD[i-1] = .01
		} else {
			WRAD[i-1] = .020 - float64(i)*.001
			if WRAD[i-1] <= 0 {
				WRAD[i-1] = (.020 - float64(19)*.001) / 2
			}
		}
	}
	//to estimate root surface and root length density per layer you need to convert root dry matter to fresh weight and scale from ha to cm^3:
	//dry matter content fresh root 7%, density fesh roo 1 gr/cm^3
	// root fresh mass
	rFreshWeight := make([]float64, g.WURZ)
	// root density
	rDense := make([]float64, g.WURZ)
	// root surface
	rSurface := make([]float64, g.WURZ)
	for i := 1; i <= g.WURZ; i++ {
		index := i - 1
		Tiefe := float64(i) * g.DZ.Num
		//Root fresh mass
		rFreshWeight[index] = (g.WUMAS * (1 - math.Exp(-Qrez*Tiefe)) / 100000 * 100 / 7)
		if i > 1 {
			rDense[index] = math.Abs(rFreshWeight[index]-rFreshWeight[index-1]) / (math.Pow(WRAD[index], 2) * math.Pi) / g.DZ.Num
		} else {
			rDense[index] = math.Abs(rFreshWeight[index]) / (math.Pow(WRAD[index], 2) * math.Pi) / g.DZ.Num
		}

		// rFreshWeight(i) = g/cm^2 from 0 to lower boundary of layer I
		//
		//  Root density /Volume soil
		// 	cm root/cm^3 soil
		g.WUDICH[index] = rDense[index]
		// ------------------------------------------------------------
		// ---------- root area cm^2/cm^3 ----------
		// ------------------------------------------------------------
		rSurface[index] = g.WUDICH[index] * WRAD[index] * 2 * math.Pi
	}
	WULAEN := 0.0
	for i := 0; i < g.WURZ; i++ {
		// ---------------  WURZELLÄNGE in cm/cm^2 -----------------------
		WULAEN = WULAEN + g.WUDICH[i]*g.DZ.Num
	}

	// ------------------------------------------------------------
	// S-uptake from root
	var SminSUM, TRNSUMS float64
	var MASS_S [20]float64
	// FOR I = 1 TO WURZ
	for i := 0; i < g.WURZ; i++ {
		//IF I < 11 THEN
		if i < 10 {
			//LET NMINSUM = NMINSUM + (S1(I)-.01)
			SminSUM = SminSUM + (g.S1[i] - 0.01)
			//LET MASS(I) = TP(I)*(S1(I)/(WG(0,1)*dz))*dt
			MASS_S[i] = g.TP[i] * (g.S1[i] / (g.WG[0][i] * g.DZ.Num)) * g.DT.Num
			// TODO: ask Christian about WG(0,1) or WG(0,i)

			//LET TRNSUM = TRNSUM + TP(I)*(S1(I)/(WG(0,I)*dz))*dt
			TRNSUMS = TRNSUMS + g.TP[i]*(g.S1[i]/(g.WG[0][i]*g.DZ.Num))*g.DT.Num

		}

	}
	// ------------------------------------------------------------

	for i := 0; i < 3; i++ {
		g.NFOS[i] = g.NFOS[i] + 0.5*WUMM/3
		g.NAOS[i] = g.NAOS[i] + 0.5*WUMM/3
	}
	// Limitieren der maximalen N-Aufnahme auf 26-13*10^-14 mol/cm W./sec
	var maxup float64
	if g.FRUCHT[g.AKF.Index] == ORH || g.FRUCHT[g.AKF.Index] == WRA || g.FRUCHT[g.AKF.Index] == SE {
		maxup = .09145 - .015725*(g.PHYLLO/1300)
	} else if g.FRUCHT[g.AKF.Index] == SM {
		maxup = .074 - .01*(g.PHYLLO/l.tendsum)
	} else if g.FRUCHT[g.AKF.Index] == ZR {
		maxup = .05645 - .01*(g.PHYLLO/l.tendsum)
	} else {
		maxup = .03145 - .015725*(g.PHYLLO/1300)
	}
	if DTGESN > WULAEN*maxup*g.DT.Num {
		if g.LEGUM != 'L' {
			DTGESN = WULAEN * maxup * g.DT.Num
		}
	}
	var MASS, D, DIFF [20]float64
	var NMINSUM, TRNSUM, SUMDIFF float64
	min := math.Min(float64(g.WURZ), g.GRW)
	for index := 0; index < int(min); index++ {
		if index+1 < 11 {
			NMINSUM = NMINSUM + (g.C1[index] - .75)
			MASS[index] = g.TP[index] * (g.C1[index] / (g.WG[0][index] * g.DZ.Num)) * g.DT.Num
			TRNSUM = TRNSUM + g.TP[index]*(g.C1[index]/(g.WG[0][index]*g.DZ.Num))*g.DT.Num
			D[index] = 2.14 * (g.AD * math.Exp(g.WG[0][index]*10)) / g.WG[0][index]
			DIFF[index] = (D[index] * g.WG[0][index] * 2 * math.Pi * WRAD[index] * (g.C1[index]/1000/g.WG[0][index] - .000014) * math.Sqrt(math.Pi*g.WUDICH[index])) * g.WUDICH[index] * 1000 * g.DT.Num
			SUMDIFF = SUMDIFF + DIFF[index]
		}
	}
	// Einschub Duengeermittlung
	SimulateFertilizationAfterPrognose(zeit, DTGESN, SUMDIFF, TRNSUM, g)
	// Ende Einschub

	var SUMPE float64
	min = math.Min(float64(g.WURZ), g.GRW)
	for index := 0; index < int(min); index++ {
		if DTGESN > 0 {
			if TRNSUM >= DTGESN {
				g.PE[index] = DTGESN * MASS[index] / TRNSUM
			} else {
				if DTGESN-TRNSUM < SUMDIFF {
					g.PE[index] = MASS[index] + (DTGESN-TRNSUM)*DIFF[index]/SUMDIFF
				} else {
					g.PE[index] = MASS[index] + DIFF[index]
				}
			}
			g.MASSUM = g.MASSUM + MASS[index]
			g.DIFFSUM = g.DIFFSUM + DIFF[index]
			if g.PE[index] > g.C1[index]-.75 {
				g.PE[index] = g.C1[index] - .75
			}
			if g.PE[index] < 0 {
				g.PE[index] = 0
			}
		} else {
			g.PE[index] = 0
		}
		SUMPE = SUMPE + g.PE[index]
	}
	// ------------------------------------------------------------
	// S-uptake from root
	var SUMPES float64
	for index := 0; index < int(min); index++ {
		if DTGESS2 > 0 {
			if TRNSUMS >= DTGESS2 {
				g.PES[index] = DTGESS2 * MASS_S[index] / TRNSUMS
			} else {
				if SminSUM > TRNSUMS {
					//LET PES(I) = MASS(I) + (DTGESS - TRNSUM) * (S1(I)-0.01-MASS(I))/(NMINSUM-TRNSUM)
					g.PES[index] = MASS_S[index] + (DTGESS2-TRNSUMS)*(g.S1[index]-.01-MASS_S[index])/(SminSUM-TRNSUMS)
				} else {
					g.PES[index] = MASS_S[index]
				}
			}
			//IF PES(I) > S1(I)-.01 THEN LET PES(I) = S1(I)-.01
			if g.PES[index] > g.S1[index]-.01 {
				g.PES[index] = g.S1[index] - .01
			}
			if g.PES[index] < 0 {
				g.PES[index] = 0
			}
		} else {
			g.PES[index] = 0
		}
		SUMPES = SUMPES + g.PES[index]
	}

	// ------------------------------------------------------------

	if g.LEGUM == 'L' {
		if DTGESN-SUMPE > 0.74*DTGESN {
			g.NFIX = 0.74 * DTGESN
		} else {
			g.NFIX = DTGESN - SUMPE
		}
	} else {
		g.NFIX = 0
	}
	g.SCHNORR = g.NFIX
	g.NFIXSUM = g.NFIXSUM + g.NFIX
	if g.WUMAS > WUMALT {
		if g.FRUCHT[g.AKF.Index] == ZR || g.FRUCHT[g.AKF.Index] == K {
			if (g.OBMAS - OBALT + g.WUMAS - WUMALT) > 0 {
				g.WUGEH = (WUMALT*g.WUGEH + ((g.WUMAS - WUMALT) / (g.OBMAS + g.WORG[3] - OBALT + g.WUMAS - WUMALT) * SUMPE)) / g.WUMAS
			}
		} else {
			if (g.OBMAS - OBALT + g.WUMAS - WUMALT) > 0 {
				g.WUGEH = (WUMALT*g.WUGEH + (g.WUMAS-WUMALT)/(g.OBMAS-OBALT+g.WUMAS-WUMALT)*(SUMPE+g.NFIX)) / g.WUMAS
			}
		}
		g.WUGEH = math.Min(g.WUGEH, g.WGMAX[g.INTWICK.Index])
		if g.WUGEH < 0.005 {
			g.WUGEH = 0.005
		}
	}
	if g.FRUCHT[g.AKF.Index] == ZR || g.FRUCHT[g.AKF.Index] == K {
		g.GEHOB = (g.PESUM + SUMPE - g.WUMAS*g.WUGEH) / (g.OBMAS + g.WORG[3])
		if g.GEHOB*(g.OBMAS+g.WORG[3]) < OBALT*GEHALT {
			g.WUGEH = (g.PESUM + SUMPE - (g.OBMAS+g.WORG[3])*g.GEHOB) / (g.WUMAS)
		}
	} else {
		g.GEHOB = (g.PESUM + SUMPE + g.NFIX - g.WUMAS*g.WUGEH) / g.OBMAS
	}

	// ------------------------------------------------------------
	//SWUGEH := (g.PESUMS + SUMPES) * g.SWura[g.FRUCHT[g.AKF.Index]]
	SWUGEH := g.WUGEH / g.SNRatio[g.FRUCHT[g.AKF.Index]]
	if g.FRUCHT[g.AKF.Index] == ZR || g.FRUCHT[g.AKF.Index] == K {

		g.SGEHOB = (g.PESUMS + SUMPES - SWUGEH) / (g.OBMAS + g.WORG[3])
	} else {
		g.SGEHOB = (g.PESUMS + SUMPES - SWUGEH) / g.OBMAS
	}

}

// radia  Strahlunsinterception, Photosynthese und Erhaltungsatmung nach Penning de Vries 1982
func radia(g *GlobalVarsMain, l *CropSharedVars, zeit int) (DLE, DLP, GPHOT, MAINT float64) {
	//! Inputs:
	//! LAT              = geogr. Breite (°)
	//! TEMP(TAG)        = Tagesmitteltemperatur (°C)
	//! RAD(TAG)         = PAR (Mj/m^2/d)
	//! CO2KONZ          = CO2 Konzentration der Atmosphäre (ppm)
	//! CO2METH          = Methode für CO2 Response
	//! MAXAMAX          = maximale C-Assimilationsrate bei Lichtsättigung und Optimaltemperatur (kg CO2/ha leave/h)
	var DL, RDN, DRC, DEC float64
	DL, DLE, DLP, _, RDN, DRC, DEC = CalculateDayLenght(g.TAG.Num, g.LAT)
	if DL <= 0 {
		return DLE, DLP, 0, 0
	}

	DRO := .2 * DRC
	EFF0 := .5
	var EFF float64
	var amax float64
	var cocomp float64
	// ! ++++++++++++++  Auswahl mehrerer Methoden zum CO2 Effect +++++++++++++++
	if g.CO2METH == 1 {
		cocomp = 17.5 * math.Pow(2, ((g.TEMP[g.TAG.Index]-10)/10))
		EFF = (g.CO2KONZ - cocomp) / (g.CO2KONZ + 2*cocomp) * EFF0
	} else if g.CO2METH == 3 {
		// ********* Gleichungen von Long 1991 und Mitchel et al. 1995 **************************
		KTvmax := math.Exp(68800 * ((g.TEMP[g.TAG.Index] + 273) - 298) / (298 * (g.TEMP[g.TAG.Index] + 273) * 8.314))
		Ktkc := math.Exp(65800 * ((g.TEMP[g.TAG.Index] + 273) - 298) / (298 * (g.TEMP[g.TAG.Index] + 273) * 8.314))
		Ktko := math.Exp(1400 * ((g.TEMP[g.TAG.Index] + 273) - 298) / (298 * (g.TEMP[g.TAG.Index] + 273) * 8.314))
		// Berechnung des Transformationsfaktors für pflanzenspez. AMAX bei 25 grad *********
		Fakamax := g.MAXAMAX / 34.695
		vcmax := 98 * Fakamax * KTvmax
		// **************************************************************************************
		MKC := 460 * Ktkc
		Mko := 210 * Ktko
		Oi := 210 + (0.047-0.0013087*g.TEMP[g.TAG.Index]+0.000025603*math.Pow(g.TEMP[g.TAG.Index], 2)-0.00000021441*math.Pow(g.TEMP[g.TAG.Index], 3))/0.026934
		Ci := g.CO2KONZ * 0.7 * (1.674 - 0.061294*g.TEMP[g.TAG.Index] + 0.0011688*math.Pow(g.TEMP[g.TAG.Index], 2) - 0.0000088741*math.Pow(g.TEMP[g.TAG.Index], 3)) / 0.73547
		cocomp = 0.5 * 0.21 * vcmax * Oi / (vcmax * Mko)
		amax = (Ci - cocomp) * vcmax / (Ci + MKC*(1+Oi/Mko)) * 1.656
		if g.TEMP[g.TAG.Index] < g.MINTMP {
			amax = 0
		}
		EFF = EFF0
	} else {
		EFF = EFF0
	}
	if l.temptyp == 1 {
		if g.CO2METH != 3 {
			if g.TEMP[g.TAG.Index] < g.MINTMP {
				amax = 0
			} else if g.TEMP[g.TAG.Index] < 10 {
				amax = g.MAXAMAX * g.TEMP[g.TAG.Index] / 10 * .4
			} else if g.TEMP[g.TAG.Index] < 15 {
				amax = g.MAXAMAX * (.4 + (g.TEMP[g.TAG.Index]-10)/5*.5)
			} else if g.TEMP[g.TAG.Index] < 25 {
				amax = g.MAXAMAX * (.9 + (g.TEMP[g.TAG.Index]-15)/10*.1)
			} else if g.TEMP[g.TAG.Index] < 35 {
				amax = g.MAXAMAX * (1 - (g.TEMP[g.TAG.Index]-25)/10)
			} else {
				amax = 0
			}
		}
		if g.CO2METH == 1 {
			amax = amax * (g.CO2KONZ - cocomp) / (350 - cocomp)
		} else if g.CO2METH == 2 {
			var KCo1 float64
			var Coco float64
			if g.RAD[g.TAG.Index] > 0 {
				KCo1 = 220 + 0.158*g.RAD[g.TAG.Index]*20
				Coco = 80 - 0.0036*g.RAD[g.TAG.Index]*20
			} else {
				SC := 1367. * (1 + 0.033*math.Cos(2*math.Pi*g.TAG.Num/365))
				EXT := SC * RDN / 10000
				Glob := EXT * (0.19 + 0.55*g.SUND[g.TAG.Index]/DL)
				KCo1 = 220 + 0.158*Glob
				Coco = 80 - 0.0036*Glob
			}
			kco2 := ((g.CO2KONZ - Coco) / (KCo1 + g.CO2KONZ - Coco)) / ((350 - Coco) / (KCo1 + 350 - Coco))
			amax = amax * kco2
		}
	} else {
		if g.TEMP[g.TAG.Index] < g.MINTMP {
			amax = 0
		} else if g.TEMP[g.TAG.Index] < 9 {
			amax = g.MAXAMAX * g.TEMP[g.TAG.Index] / 10 * .0555
		} else if g.TEMP[g.TAG.Index] < 16 {
			amax = g.MAXAMAX * (.05 + (g.TEMP[g.TAG.Index]-9)/7*.75)
		} else if g.TEMP[g.TAG.Index] < 18 {
			amax = g.MAXAMAX * (.8 + (g.TEMP[g.TAG.Index]-16)*.07)
		} else if g.TEMP[g.TAG.Index] < 20 {
			amax = g.MAXAMAX * (.94 + (g.TEMP[g.TAG.Index]-18)*.03)
		} else if g.TEMP[g.TAG.Index] >= 20 && g.TEMP[g.TAG.Index] <= 30 {
			amax = g.MAXAMAX
		} else if g.TEMP[g.TAG.Index] < 36 {
			amax = g.MAXAMAX * (1 - (g.TEMP[g.TAG.Index]-30)*.0083)
		} else if g.TEMP[g.TAG.Index] < 42 {
			amax = g.MAXAMAX * (1 - (g.TEMP[g.TAG.Index]-36)*.0065)
		} else {
			amax = 0
		}
	}
	if amax < 0.1 {
		amax = 0.1
	}
	if DLE == 0 && DL > 0 {
		DLE = 0.1
	}
	REFLC := .08
	EFFE := (1. - REFLC) * EFF
	SSLAE := math.Sin((90. + DEC - g.LAT) * math.Pi / 180.)
	X := math.Log(1. + .45*DRC/(DLE*3600.)*EFFE/(SSLAE*amax))
	PHCH1 := SSLAE * amax * DLE * X / (1. + X)
	Y := math.Log(1. + .55*DRC/(DLE*3600.)*EFFE/((5-SSLAE)*amax))
	PHCH2 := (5. - SSLAE) * amax * DLE * Y / (1. + Y)
	PHCH := 0.95*(PHCH1+PHCH2) + 20.5
	PHC3 := PHCH * (1. - math.Exp(-.8*g.LAI))
	PHC4 := DL * g.LAI * amax
	var MIPHC, MAPHC float64
	if PHC3 < PHC4 {
		MIPHC = PHC3
		MAPHC = PHC4
	} else {
		MIPHC = PHC4
		MAPHC = PHC3
	}
	if MIPHC == 0 {
		MIPHC = 0.000001
	}
	PHCL := MIPHC * (1. - math.Exp(-MAPHC/MIPHC))
	Z := DRO / (DLE * 3600.) * EFFE / (5. * amax)
	PHOH1 := 5. * amax * DLE * Z / (1. + Z)
	PHOH := 0.9935*PHOH1 + 1.1
	PHO3 := PHOH * (1. - math.Exp(-.8*g.LAI))
	var MIPHO, MAPHO float64
	if PHO3 < PHC4 {
		MIPHO = PHO3
		MAPHO = PHC4
	} else {
		MIPHO = PHC4
		MAPHO = PHO3
	}
	if MIPHO == 0 {
		MIPHO = 0.000001
	}
	PHOL := MIPHO * (1. - math.Exp(-MAPHO/MIPHO))
	var DGAC, DGAO float64
	if g.LAI-5 < 0 {
		DGAC = PHCL
		DGAO = PHOL
	} else {
		DGAC = PHCH
		DGAO = PHOH
	}
	var DTGA float64
	// ----------- BERÜCKSICHTIGUNG DER SONNENSCHEINDAUER -------
	if g.RAD[g.TAG.Index] == 0 {
		if g.SUND[g.TAG.Index] > DLE {
			g.SUND[g.TAG.Index] = DLE
		}
		DTGA = g.SUND[g.TAG.Index]/DLE*DGAC + (1.-g.SUND[g.TAG.Index]/DLE)*DGAO
	} else {
		KOREK := 1.
		g.RADSUM = g.RADSUM + g.RAD[g.TAG.Index]*g.DT.Num*KOREK
		FOV := (DRC - 1000000*g.RAD[g.TAG.Index]*KOREK) / (.8 * DRC)
		if FOV > 1 {
			FOV = 1
		}
		if FOV < 0 {
			FOV = 0
		}
		DTGA = FOV*DGAO + (1-FOV)*DGAC
	}
	// !     ------- PHOTOSYNTHESERATE IN KG GLUCOSE/HA BLATT/TAG------
	GPHOT = DTGA * 30. / 44
	var vswell float64
	if g.LURED == 1 {
		vswell = g.DRYSWELL[g.INTWICK.Index]
	} else {
		if g.FRUCHT[g.AKF.Index] == SM || g.FRUCHT[g.AKF.Index] == K || g.FRUCHT[g.AKF.Index] == WR || g.FRUCHT[g.AKF.Index] == SG || g.FRUCHT[g.AKF.Index] == WW || g.FRUCHT[g.AKF.Index] == WG {
			vswell = 1
		} else {
			vswell = 0.8
		}
	}
	if g.TRREL < vswell {
		GPHOT = GPHOT * g.TRREL
	}
	// ! ----------- MAINTENANCE IN ABH. VON TEMPERATUR -----------
	TEFF := math.Pow(2., (.1*g.TEMP[g.TAG.Index] - 2.5))
	MAINORG := make([]float64, g.NRKOM)
	var MAINTS float64
	for i := 0; i < g.NRKOM; i++ {
		MAINTS = MAINTS + g.WORG[i]*g.MAIRT[i]
		MAINORG[i] = g.WORG[i] * g.MAIRT[i]
	}
	for i := 0; i < g.NRKOM; i++ {
		l.MANT[i] = MAINORG[i] / MAINTS
	}

	if GPHOT < MAINTS*TEFF {
		MAINT = GPHOT
	} else {
		MAINT = MAINTS * TEFF
	}
	if g.TEMP[g.TAG.Index] < g.MINTMP {
		GPHOT = MAINT
	}
	return DLE, DLP, GPHOT, MAINT
}

func root(veloc, tempsum, dz float64) (qrez, potentialRootingDepth float64, culRootPercPerLayer []float64) {
	// Qrez = MAX( (0.081476+math.Exp((-Veloc*(A3+Tsumbase)))^1.8;0.0409)
	// Veloc = increase root depth(cm/°C) / 200
	// Tsumbase = LOG(0.35^(1/1.8)-0.081476;EXP(-Veloc))
	//
	// rooting depth = 4.5/Qrez
	// cumulative percentage until layer I (column H-S) = (1-EXP(-QREZ*ry(I)lower bounda))*100

	Tsumbase := math.Log(math.Pow(0.35, 1/1.8)-0.081476) / math.Log(math.Exp(-veloc))
	qrez = math.Max(math.Pow((0.081476+math.Exp(-veloc*(tempsum+Tsumbase))), 1.8), 0.022)
	//qrez = math.Max(math.Pow((0.081476+math.Exp(-veloc*(tempsum+Tsumbase))), 1.8), 0.0409)

	potentialRootingDepth = 4.5 / qrez
	rootLayer := int(potentialRootingDepth / dz) // WURZ

	// cumulative percentage until layer I (column H-S) = (1-EXP(-QREZ*lower boundary(I)))*100
	culRootPercPerLayer = make([]float64, rootLayer)
	for i := 1; i <= rootLayer; i++ {
		culRootPercPerLayer[i-1] = (1 - math.Exp((-1.0)*qrez*(float64(i)*10))) * 100
	}

	return qrez, potentialRootingDepth, culRootPercPerLayer
}

func vern(l *CropSharedVars, g *GlobalVarsMain) {
	// Sub-Modul zur Berechnung der Vernalisation
	// Inputs:
	// TEMP(TAG)           = Tagesmitteltemperatur (°C)
	// VSCHWELL(INTWICK)   = notwendige Anzahl der Vernalisationstage (Tag) für Stadium INTWICK
	// Output:
	// FV                  = Vernalisationsfaktor zur Korrektur von SUM(INTWICK)
	var veff float64
	if g.TEMP[g.TAG.Index] < 0 && g.TEMP[g.TAG.Index] > -4 {
		veff = (g.TEMP[g.TAG.Index] + 4) / 4
	} else if g.TEMP[g.TAG.Index] < -4 {
		veff = 0
	} else if g.TEMP[g.TAG.Index] > 3 && g.TEMP[g.TAG.Index] < 7 {
		veff = 1 - .2*(g.TEMP[g.TAG.Index]-3)/4
	} else if g.TEMP[g.TAG.Index] > 7 && g.TEMP[g.TAG.Index] < 9 {
		veff = .8 - .4*(g.TEMP[g.TAG.Index]-7)/2
	} else if g.TEMP[g.TAG.Index] > 9 && g.TEMP[g.TAG.Index] < 18 {
		veff = .4 - .4*(g.TEMP[g.TAG.Index]-9)/9
	} else if g.TEMP[g.TAG.Index] < -4 || g.TEMP[g.TAG.Index] > 18 {
		veff = 0
	} else {
		veff = 1
	}
	g.VERNTAGE = g.VERNTAGE + veff*g.DT.Num
	verschwell := math.Min(g.VSCHWELL[g.INTWICK.Index], 9) - 1
	if verschwell >= 1 {
		l.FV = (g.VERNTAGE - verschwell) / (g.VSCHWELL[g.INTWICK.Index] - verschwell)
		if l.FV < 0 {
			l.FV = 0
		} else if l.FV > 1 {
			l.FV = 1
		}
	} else {
		l.FV = 1
	}
}
