package hermes

import (
	"bufio"
	"fmt"
	"math"
	"strings"
)

// InputSharedVars is a struct of shared variables for this module
type InputSharedVars struct {
	NORG    [300]float64
	NDGMG   [300]float64 // amount of nitrogen in fertilizer
	SDGMG   [300]float64 // amount of sulfur in fertiliser
	NGEHALT [10]float64
	Jstr    string
	MK      [70]string
	FK      [10]float64
	KONZ1   float64
	KONZ3   float64
	KONZ4   float64
	KONZ5   float64
	KONZ6   float64
	KONZ7   float64
	IRRIGAT bool
	ANZBREG int
	SSAND   [10]float64
	SLUF    [10]float64
	TON     [10]float64
}

// Input modul for reading soil data, crop rotation, cultivation data (Fertilization, tillage) of fields and ploygon units
func Input(scanner *bufio.Scanner, l *InputSharedVars, g *GlobalVarsMain, hPath *HFilePath, driConfig *config, soilID string) error {
	//! ------Modul zum Einlesen von Boden-, Fruchtfolge und Bewirtschaftungsdaten (Duengung, Bodenbearbeitung) von Feldern und Polygonen ---------
	var ERNT, SAT string
	var winit [6]float64

	//!  Einleseprogramm für Schlagdaten
	g.WRED = 0
	// ! ----------------------- Beginn Lesen der Polygondatei ------------------------
	// ! Inputs:
	// ! FLAEID$        = Polygon-ID
	// ! GRHI			= Grundwassserhöchststand (dm u. Flur)
	// ! GRLO			= Grundwasssertiefststand (dm u. Flur)
	// ! PKT$           = Feld-ID für Feldbezogene Daten
	// ! IRRIGAT        = Trigger zum Einlesen von Bewässerungsdaten
	// ! BOF$           = Boden-ID
	// ! ------------------------------------------------------------------------------

	for scanner.Scan() {
		wa := scanner.Text()
		punr := int(ValAsInt(wa[0:5], "none", wa))
		if punr == g.SLNR {
			FLAEID := wa[0:5]
			if g.GROUNDWATERFROM == Polygonfile {
				g.GRHI = int(ValAsInt(wa[20:22], "none", wa))
				g.GRLO = int(ValAsInt(wa[23:25], "none", wa))
				g.GRW = float64(g.GRLO+g.GRHI) / 2
				g.GW = float64(g.GRLO+g.GRHI) / 2
				g.AMPL = g.GRLO - g.GRHI
			}

			g.PKT = wa[10:19]                            // Feld_ID / Field_ID
			g.PKT = strings.TrimSpace(g.PKT)             // remove spaces
			l.IRRIGAT = ValAsBool(wa[26:27], "none", wa) // irrigation on/off 1/0

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
			sid := soilID
			if soilID == "" {
				sid = wa[6:9] // second entry SID in poly file
			}
			var currentSoil soilFileData
			var soilLoadError error
			groundwaterFormSoilfile := g.GROUNDWATERFROM == Soilfile

			if strings.HasSuffix(hPath.bofile, "csv") {
				currentSoil, soilLoadError = loadSoilCSV(groundwaterFormSoilfile, g.LOGID, hPath, sid)
			} else {
				currentSoil, soilLoadError = loadSoil(groundwaterFormSoilfile, g.LOGID, hPath, sid)
			}
			if soilLoadError != nil {
				return soilLoadError
			}

			g.SoilID = currentSoil.SoilID
			g.N = currentSoil.N
			g.AZHO = currentSoil.AZHO
			g.WURZMAX = currentSoil.WURZMAX

			if currentSoil.useGroundwaterFromSoilfile {
				g.GRHI = currentSoil.GRHI
				g.GRLO = currentSoil.GRLO
				g.GRW = currentSoil.GRW
				g.GW = currentSoil.GW
				g.AMPL = currentSoil.AMPL
			}

			g.DRAIDEP = currentSoil.DRAIDEP
			g.DRAIFAK = currentSoil.DRAIFAK
			g.UKT = currentSoil.UKT
			g.BART = currentSoil.BART
			g.LD = currentSoil.LD
			g.BULK = currentSoil.BULK
			g.CGEHALT = currentSoil.CGEHALT
			g.SGEHALT = currentSoil.SGEHALT
			g.CNRAT1 = currentSoil.CNRAT1
			l.NGEHALT = currentSoil.NGEHALT
			g.HUMUS = currentSoil.HUMUS
			g.STEIN = currentSoil.STEIN
			g.FKA = currentSoil.FKA
			g.WP = currentSoil.WP
			g.GPV = currentSoil.GPV
			l.SSAND = currentSoil.SSAND
			l.SLUF = currentSoil.SLUF
			l.TON = currentSoil.TON

			g.DT.SetByIndex(1)

			// ! *************************** Bodenparameter zuweisen ***********************
			// ! Inputs aus HYDRO:(I=L=Horizontzähler)
			// ! FELDW(I)           = Wassergehalt bei Feldkapazität (cm^3/cm^3)
			// ! LIM(I)             = Wassergehalt bei PWP (cm^3/cm^3)
			// ! PRGES(I)	        = Gesamtporenvolumen (cm^3/cm^3)
			// ! NORMFK(I)          = Wassergehalt bei Feldkapazität unkorrigiert (cm^3/cm^3)
			// ! Ableitungen:
			// ! LT                 = Zähler 10cm Schichten
			// ! W(LT)			    = Wassergehalt bei Feldkapazität (cm^3/cm^3)
			// ! WMIN(LT)           = Wassergehalt bei PWP (cm^3/cm^3)
			// ! PORGES(LT)         = Gesamtporenvolumen (cm^3/cm^3)
			// ! WNOR(LT)           = Wassergehalt bei Feldkapazität unkorrigiert (cm^3/cm^3)

			for L := 1; L <= g.AZHO; L++ {
				lindex := L - 1
				BDART := Hydro(L, g, l, hPath)
				if g.FELDW[lindex] == 0 {
					g.FELDW[lindex] = g.FELDW[lindex-1]
				}
				// for every 10 cm in this layer
				for LT := g.UKT[L-1] + 1; LT <= g.UKT[L]; LT++ {
					LTindex := LT - 1
					if g.PTF == 0 {
						if g.FKA[lindex] > 0 {
							g.CAPPAR = 1
							if LT < g.N+1 {
								g.BD[LTindex] = g.BULK[lindex]
								g.W[LTindex] = g.FKA[lindex] / 100
								g.WMIN[LTindex] = g.WP[lindex] / 100
								g.PORGES[LTindex] = g.GPV[lindex] / 100
								g.WNOR[LTindex] = g.FKA[lindex] / 100

								if BDART[0] == 'S' { // if main soil component is sand
									g.WRED = g.WP[lindex] + 0.6*(g.FKA[lindex]-g.WP[lindex])
								} else {
									g.WRED = g.WP[lindex] + 0.66*(g.FKA[lindex]-g.WP[lindex])
								}
							}
						} else {
							g.CAPPAR = 0
							if LT < g.N+1 {
								g.BD[LTindex] = g.BULK[lindex]
								g.W[LTindex] = g.FELDW[lindex] * (1 - g.STEIN[lindex])
								g.WMIN[LTindex] = g.LIM[lindex] * (1 - g.STEIN[lindex])
								g.PORGES[LTindex] = g.PRGES[lindex] * (1 - g.STEIN[lindex])
								g.WNOR[LTindex] = g.NORMFK[lindex] * (1 - g.STEIN[lindex])
								g.SAND[LTindex] = l.SSAND[lindex]
								g.SILT[LTindex] = l.SLUF[lindex]
								g.CLAY[LTindex] = l.TON[lindex]
							}
						}
					} else {
						g.CAPPAR = 1
						if LT < g.N+1 {
							// check if sand, silt, clay are valid
							soilSum := l.TON[lindex] + l.SLUF[lindex] + l.SSAND[lindex]
							if soilSum > 103 || soilSum < 97 { // rounding issue
								return fmt.Errorf("sand: %f, Silt: %f, Clay: %f does not sum up to 100 percent", l.SSAND[lindex], l.SLUF[lindex], l.TON[lindex])
							}
							if l.TON[lindex] == 0 {
								return fmt.Errorf("clay content is 0")
							}
							if l.SLUF[lindex] == 0 {
								return fmt.Errorf("silt content is 0")
							}
							if l.SSAND[lindex] == 0 {
								return fmt.Errorf("sand content is 0")
							}
							if g.PTF == 1 {
								// PTF by Toth 2015
								fk := 0.2449 - 0.1887*(1/(g.CGEHALT[lindex]+1)) + 0.004527*l.TON[lindex] + 0.001535*l.SLUF[lindex] + 0.001442*l.SLUF[lindex]*(1/(g.CGEHALT[lindex]+1)) - 0.0000511*l.SLUF[lindex]*l.TON[lindex] + 0.0008676*l.TON[lindex]*(1/(g.CGEHALT[lindex]+1))
								g.W[LTindex] = fk
								pwp := 0.09878 + 0.002127*l.TON[lindex] - 0.0008366*l.SLUF[lindex] - 0.0767*(1/(g.CGEHALT[lindex]+1)) + 0.00003853*l.SLUF[lindex]*l.TON[lindex] + 0.00233*l.SLUF[lindex]*(1/(g.CGEHALT[lindex]+1)) + 0.0009498*l.SLUF[lindex]*(1/(g.CGEHALT[lindex]+1))
								g.WMIN[LTindex] = pwp
							} else if g.PTF == 2 {
								// PTF by Batjes for pF 2.5
								g.W[LTindex] = (0.46*l.TON[lindex] + 0.3045*l.SLUF[lindex] + 2.0703*g.CGEHALT[lindex]) / 100
								g.WMIN[LTindex] = (0.3624*l.TON[lindex] + 0.117*l.SLUF[lindex] + 1.6054*g.CGEHALT[lindex]) / 100
							} else if g.PTF == 3 {
								// PTF by Batjes for pF 1.7
								g.W[LTindex] = (0.6681*l.TON[lindex] + 0.2614*l.SLUF[lindex] + 2.215*g.CGEHALT[lindex]) / 100
								g.WMIN[LTindex] = (0.3624*l.TON[lindex] + 0.117*l.SLUF[lindex] + 1.6054*g.CGEHALT[lindex]) / 100
							} else if g.PTF == 4 {
								// PTF by Rawls et al. 2003 for pF 2.5
								ix := -0.837531 + 0.430183*g.CGEHALT[lindex]
								ix2 := math.Pow(ix, 2)
								ix3 := math.Pow(ix, 3)
								yps := -1.40744 + 0.0661969*l.TON[lindex]
								yps2 := math.Pow(yps, 2)
								yps3 := math.Pow(yps, 3)
								zet := -1.51866 + 0.0393284*l.SSAND[lindex]
								zet2 := math.Pow(zet, 2)
								zet3 := math.Pow(zet, 3)

								g.W[LTindex] = (29.7528 + 10.3544*(0.0461615+0.290955*ix-0.0496845*ix2+0.00704802*ix3+0.269101*yps-0.176528*ix*yps+0.0543138*ix2*yps+0.1982*yps2-0.060699*yps3-0.320249*zet-0.0111693*ix2*zet+0.14104*yps*zet+0.0657345*ix*yps*zet-0.102026*yps2*zet-0.04012*zet2+0.160838*ix*zet2-0.121392*yps*zet2-0.061667*zet3)) / 100
								g.WMIN[LTindex] = (14.2568 + 7.36318*(0.06865+0.108713*ix-0.0157225*ix2+0.00102805*ix3+0.886569*yps-0.223581*ix*yps+0.0126379*ix2*yps+0.0135266*ix*yps2-0.0334434*yps3-0.0535182*zet-0.0354271*ix*zet-0.00261313*ix2*zet-0.154563*yps*zet-0.0160219*ix*yps*zet-0.0400606*yps2*zet-0.104875*zet2*0.0159857*ix*zet2-0.0671656*yps*zet2-0.0260699*zet3)) / 100
							}
							g.PORGES[LTindex] = g.GPV[lindex] / 100
							g.BD[LTindex] = g.BULK[lindex]
							g.WNOR[LTindex] = g.W[LTindex]
							if BDART[0] == 'S' { // if main soil component is sand
								g.WRED = (g.WMIN[LTindex] + 0.6*(g.W[LTindex]-g.WMIN[LTindex])) * 100
							} else {
								g.WRED = (g.WMIN[LTindex] + 0.66*(g.W[LTindex]-g.WMIN[LTindex])) * 100
							}
						}
					}
				}
			}
			g.WRED = g.WRED / 100
			// ! -- Unterhalb Grundwasserspiegel wird FK auf GPV gesetzt --
			// below groundwater level FK will be set to GPV
			if g.GW < float64(g.N) {
				maxVal := math.Max(g.GW, 1)
				for l := int(math.Round(maxVal)); l <= g.N; l++ {
					index := l - 1
					g.W[index] = g.PORGES[index]
				}
			}
			if !g.AUTOIRRI {
				// ! +++++++++++++++++++ Einlesen der schlagspez. Beregnung ++++++++++++++++++++++++
				// ! Inputs:
				// ! ANZBREG                   = Zähler für Bewaesserungsmassnahmen
				// ! BREGDAT$(ANZBREG)         = Datum der Massnahme (TTMMJJ)
				// ! BREG(ANZBREG)             = Bewaesserungsmenge (mm)
				// ! BRKZ(ANBREG)              = N-Konzentration des Bewaesserungswassers (ppm)
				// ! abgeleitet
				// ! ZTBR(ANBREG)              = Tag der Massnahme in ZEIT Einheit (ab 1.1.1900)
				// ! -------------------------------------------------------------------------------

				l.ANZBREG = 0
				if l.IRRIGAT {
					Bereg := hPath.irrigation
					_, scannerIrrFile, _ := Open(&FileDescriptior{FilePath: Bereg, FileDescription: "irrigation file", UseFilePool: true})
					LineInut(scannerIrrFile)

					for SCHLAG, SLAGtoken, ok := NextLineInut(0, scannerIrrFile, strings.Fields); ok; SCHLAG, SLAGtoken, ok = NextLineInut(0, scannerIrrFile, strings.Fields) {
						valid := true
						dateTokenIndex := 3
						for ok := SCHLAG == g.PKT; ok; ok = SCHLAG == g.PKT && valid {
							l.ANZBREG++
							g.BREG[l.ANZBREG-1] = ValAsFloat(SLAGtoken[1], Bereg, SLAGtoken[1])
							g.BRKZn[l.ANZBREG-1] = ValAsFloat(SLAGtoken[2], Bereg, SLAGtoken[2])
							if len(SLAGtoken) > 4 {
								g.BRKZs[l.ANZBREG-1] = ValAsFloat(SLAGtoken[3], Bereg, SLAGtoken[3])
								dateTokenIndex = 4
							}
							BREGDAT := SLAGtoken[dateTokenIndex]
							_, g.ZTBR[l.ANZBREG-1] = g.Datum(BREGDAT)

							///!warning may Beginn not yet initialized
							if g.ZTBR[l.ANZBREG-1] < g.BEGINN {
								l.ANZBREG--
							}
							SCHLAG, SLAGtoken, valid = NextLineInut(0, scannerIrrFile, strings.Fields)
						}

					}
					for i := l.ANZBREG; i < 500; i++ {
						g.ZTBR[i] = 0
						g.BREG[i] = 0
						g.BRKZn[i] = 0
						g.BRKZs[i] = 0
					}
				} else {
					l.ANZBREG = 0
					for i := 0; i < 500; i++ {
						g.ZTBR[i] = 0
						g.BREG[i] = 0
						g.BRKZn[i] = 0
						g.BRKZs[i] = 0
					}
				}
			}
			// ! -----------------------------Ende Bewaessungsdaten -------------------------------
			// ! ***************************************** Fruchtfolge und Bestellungsdatei lesen ************************************
			// ! INPUTS:
			// ! SLFIND                  = Zähler für Fruchtfolgeelement (Frucht und Bestellungstermine)
			// ! FRUCHT$(SLFIND)         = Anbaufrucht (3 stelliger Fruchtkuerzel)
			// ! SAT$		             = Aussaatdatum Frucht (TTMMJJ)
			// ! ERNT$(SLFIND)           = Erntedatum Frucht (TTMMJJ)
			// ! ERTR(1)                 = Ertrag 1. Frucht (dt/ha) (nur Vorfrucht, Ernte = Beginn der Simulation)
			// ! JN(SLFIND)              = Anteil exportierte Ernterückstände (%) (100 = alles abgefahren, 0 = Verbleib auf dem Feld)
			// ! abgeleitet:
			// ! Saat(SLFIND)            = Tag der Aussaat (seit 1.1.1900)
			// ! ERNTE(SLFIND)           = Tag der Ernte   (seit 1.1.1900)
			// ! ITAG                    = DOY des Simulationsstarts
			// !----------------------------------------------------------------------------------------------------------------------

			ROTA := hPath.crop
			_, scannerRotation, _ := Open(&FileDescriptior{FilePath: ROTA, FileDescription: "rotation file", UseFilePool: true})
			cropHeader := LineInut(scannerRotation)

			hSchlag := 0
			hCrop := 1
			hSow := 2
			hHarvest := 3
			hJN := 4
			hHarvestResidue := 5
			hOrgDung := 6
			hVariety := 7
			splitLine := func(s string) []string {
				return strings.Fields(s)
			}
			if driConfig.CropFileFormat == "csv" {
				splitLine = func(s string) []string {
					return strings.Split(s, ",")
				}
				headlineTokens := strings.Split(cropHeader, ",")
				for i, t := range headlineTokens {
					if t == "Field_ID" {
						hSchlag = i
					} else if t == "crop" {
						hCrop = i
					} else if t == "sowing" {
						hSow = i
					} else if t == "harvest" {
						hHarvest = i
					} else if t == "Rex" {
						hJN = i
					} else if t == "yld" {
						hHarvestResidue = i
					} else if t == "autorg" {
						hOrgDung = i
					} else if t == "variety" {
						hVariety = i
					}
				}
			}

			checkDate := func() func(date string) {
				currentDate := 0
				currentDateStr := ""
				return func(date string) {
					if currentDate == 0 {
						_, currentDate = g.Datum(date)
						currentDateStr = date
					} else {
						if _, dateValue := g.Datum(date); dateValue <= currentDate {
							panic(fmt.Sprintf("Date %s is before %s", currentDateStr, date))
						} else {
							currentDate = dateValue
							currentDateStr = date
						}
					}
				}
			}()
			SLFIND := 0
			for SCHLAG, ROtoken, valid := NextLineInut(hSchlag, scannerRotation, splitLine); valid; SCHLAG, ROtoken, valid = NextLineInut(hSchlag, scannerRotation, splitLine) {
				for ok := SCHLAG == g.PKT; ok; ok = SCHLAG == g.PKT && valid {
					SLFIND++
					SLFINDindex := SLFIND - 1
					g.FRUCHT[SLFINDindex] = g.ToCropType(ROtoken[hCrop])
					if len(ROtoken) > hVariety {
						g.CVARIETY[SLFINDindex] = ROtoken[hVariety]
					}
					if SLFIND > 1 {
						SAT = ROtoken[hSow]
						checkDate(SAT)
					}

					ERNT = ROtoken[hHarvest]
					checkDate(ERNT)
					if len(ROtoken) > hOrgDung {
						g.ODU[SLFINDindex] = ValAsFloat(ROtoken[hOrgDung], ROTA, ROtoken[hOrgDung])
					} else {
						g.ODU[SLFINDindex] = 0
					}
					if !g.AUTOMAN {
						if SLFIND > 1 {
							_, g.SAAT[SLFINDindex] = g.Datum(SAT)
						}
					}
					if !g.AUTOHAR {
						var ERNDAT int
						ERNDAT, g.ERNTE[SLFINDindex] = g.Datum(ERNT)
						if SLFIND == 1 {
							g.ITAG = ERNDAT
						}
						g.ERNTE2[SLFINDindex] = g.ERNTE[SLFINDindex]
					}

					if SLFIND > 1 {
						if g.AUTOIRRI || g.AUTOFERT || g.AUTOHAR || g.AUTOMAN {
							autfil := hPath.auto
							_, autoScanner, _ := Open(&FileDescriptior{FilePath: autfil, FileDescription: "automated file", UseFilePool: true})
							LineInut(autoScanner)
							for autoScanner.Scan() {
								crpman := autoScanner.Text()
								if g.ToCropType(crpman[0:3]) == g.FRUCHT[SLFINDindex] {
									if g.AUTOMAN {
										if ValAsInt(crpman[4:8], autfil, crpman) == 0 {
											SAT = ROtoken[hSow]
											_, g.SAAT[SLFINDindex] = g.Datum(SAT)
											g.SAAT1[SLFINDindex] = g.SAAT[SLFINDindex] - 1
											g.SAAT2[SLFINDindex] = g.SAAT[SLFINDindex]
										} else {
											sat1 := crpman[4:8] + SAT[4:]
											sat2 := crpman[9:13] + SAT[4:]
											if crpman[24:25] == "x" {
												g.TSLMAX[SLFINDindex] = ValAsFloat(crpman[19:24], autfil, crpman)
												g.TSLMIN[SLFINDindex] = -1
											} else {
												g.TSLMAX[SLFINDindex] = -1
												g.TSLMIN[SLFINDindex] = ValAsFloat(crpman[19:24], autfil, crpman)
											}
											g.MINMOI[SLFINDindex] = ValAsFloat(crpman[25:30], autfil, crpman)
											g.MAXMOI[SLFINDindex] = ValAsFloat(crpman[32:37], autfil, crpman)
											_, g.SAAT1[SLFINDindex] = g.Datum(sat1)
											_, g.SAAT2[SLFINDindex] = g.Datum(sat2)
											g.SAAT[SLFINDindex] = 0
											g.TJAHR[SLFINDindex] = ValAsFloat(crpman[68:71], autfil, crpman)
											g.TJBAS[SLFINDindex] = ValAsFloat(crpman[74:76], autfil, crpman)
											g.TSLWINDOW[SLFINDindex] = ValAsFloat(crpman[135:137], autfil, crpman)
										}
									}
									if g.AUTOHAR {
										if ValAsInt(crpman[14:18], autfil, crpman) == 0 {
											_, g.ERNTE2[SLFINDindex] = g.Datum(ERNT)
										} else {
											har2 := crpman[14:18] + ERNT[4:]
											g.MINHMOI[SLFINDindex] = ValAsFloat(crpman[39:44], autfil, crpman)
											g.MAXHMOI[SLFINDindex] = ValAsFloat(crpman[46:51], autfil, crpman)
											g.RAINLIM[SLFINDindex] = ValAsFloat(crpman[53:57], autfil, crpman)
											g.RAINACT[SLFINDindex] = ValAsFloat(crpman[60:64], autfil, crpman)
											_, g.ERNTE2[SLFINDindex] = g.Datum(har2)
											g.ERNTE[SLFINDindex] = 0
										}
									}
									if g.AUTOIRRI {
										g.IRRST1[SLFINDindex] = ValAsFloat(crpman[80:81], autfil, crpman)
										g.IRRST2[SLFINDindex] = ValAsFloat(crpman[87:88], autfil, crpman)
										g.IRRLOW[SLFINDindex] = ValAsFloat(crpman[163:166], autfil, crpman) / 100
										g.IRRDEP[SLFINDindex] = ValAsFloat(crpman[170:173], autfil, crpman) / 10
										g.IRRMAX[SLFINDindex] = ValAsFloat(crpman[177:180], autfil, crpman)
									}
									if g.AUTOFERT {
										g.NDEM1[SLFINDindex] = ValAsFloat(crpman[94:97], autfil, crpman)
										if g.ODU[SLFINDindex] == 1 {
											g.DGART[SLFINDindex] = strings.TrimSpace(crpman[143:146])
											l.NDGMG[SLFINDindex] = ValAsFloat(crpman[149:152], autfil, crpman)
											g.ORGTIME[SLFINDindex] = crpman[156:157]
											g.ORGDOY[SLFINDindex] = int(ValAsInt(crpman[157:159], autfil, crpman))
											dueng(SLFINDindex, g, l, hPath)
										} else {
											g.ORGTIME[SLFINDindex] = "0"
											g.ORGDOY[SLFINDindex] = 0
											g.NDIR[SLFINDindex] = 0
											g.NLAS[SLFINDindex] = 0
											g.NSAS[SLFINDindex] = 0
										}
										g.NDEM2[SLFINDindex] = ValAsFloat(crpman[100:103], autfil, crpman)
										g.NDEM3[SLFINDindex] = ValAsFloat(crpman[106:109], autfil, crpman)
										if crpman[112:113] == "S" {
											g.NDOY1[SLFINDindex] = ValAsFloat(crpman[113:115], autfil, crpman)
										} else {
											g.NDOY1[SLFINDindex] = ValAsFloat(crpman[112:115], autfil, crpman)
										}
										if crpman[119:120] == "S" {
											g.NDOY2[SLFINDindex] = ValAsFloat(crpman[120:122], autfil, crpman)
										} else {
											g.NDOY2[SLFINDindex] = ValAsFloat(crpman[119:122], autfil, crpman)
										}
										if crpman[127:128] == "S" {
											g.NDOY3[SLFINDindex] = ValAsFloat(crpman[128:130], autfil, crpman)
										} else {
											g.NDOY3[SLFINDindex] = ValAsFloat(crpman[127:130], autfil, crpman)
										}
										g.TSLWINDOW[SLFINDindex] = ValAsFloat(crpman[135:137], autfil, crpman)
									}
									break
								}
							}
						}
					} else {
						var ERNDAT int
						ERNDAT, g.ERNTE[SLFINDindex] = g.Datum(ERNT)
						g.ITAG = ERNDAT
						if g.AUTOHAR || g.AUTOFERT {
							autfil := hPath.auto
							_, autoScanner, _ := Open(&FileDescriptior{FilePath: autfil, FileDescription: "automated file", UseFilePool: true})
							LineInut(autoScanner)
							for autoScanner.Scan() {
								crpman := autoScanner.Text()
								if g.ToCropType(crpman[0:3]) == g.FRUCHT[SLFINDindex] {
									if g.ODU[SLFINDindex] == 1 {
										g.DGART[SLFINDindex] = strings.TrimSpace(crpman[143:146])
										l.NDGMG[SLFINDindex] = ValAsFloat(crpman[149:152], autfil, crpman)
										g.ORGTIME[SLFINDindex] = crpman[156:157]
										g.ORGDOY[SLFINDindex] = int(ValAsInt(crpman[157:159], autfil, crpman))
										dueng(SLFIND, g, l, hPath)
									} else {
										g.ORGTIME[SLFINDindex] = "0"
										g.ORGDOY[SLFINDindex] = 0
										g.NDIR[SLFINDindex] = 0
										g.NLAS[SLFINDindex] = 0
										g.NSAS[SLFINDindex] = 0
									}
									break
								}
							}
						}
					}
					g.JN[SLFINDindex] = ValAsFloat(ROtoken[hJN], ROTA, ROtoken[hJN]) / 100
					if SLFIND == 1 {
						g.ERTR[SLFINDindex] = ValAsFloat(ROtoken[hHarvestResidue], ROTA, ROtoken[hHarvestResidue])
					}
					SCHLAG, ROtoken, valid = NextLineInut(hSchlag, scannerRotation, splitLine)

					if SCHLAG != g.PKT {
						g.FRUCHT[SLFINDindex+1] = SM // TODO: Why hardcoded SM?
						g.ERTR[SLFINDindex+1] = 0
						g.SAAT1[SLFINDindex+1] = g.SAAT[SLFINDindex] + 365
						g.SAAT2[SLFINDindex+1] = g.SAAT[SLFINDindex] + 365
						g.TSLWINDOW[SLFINDindex+1] = 5
					}
				}
			}

			err := sReadCropData(g, hPath)
			if err != nil {
				// failed to read crop data
				return err
			}
			// ! -- Setzen des Simulationsbeginns für Zeitschleife
			// set simulation start for time loop
			g.BEGINN = g.ERNTE[0]
			// ! Ernte der 1. Frucht = Düngung Nr. 1 mit Ernterückständen
			// Harvest of first crop = Fertilization nr. 1 with harvest residue
			g.ZTDG[0] = g.ERNTE[0]
			// ! Einlesen der Monatsfaktoren für HAUDE
			// loading of input data for HAUDE
			verdun(g, hPath)
			// ! ---- Ableitung der Anfangs-Nmin-Verteilung in Abhängigkeit von Vorfrucht -----
			// Deriving of start-N-min-Distribution in relation to previous crop
			for m := 1; m <= g.N+1; m++ {
				if g.FRUCHT[0] == ZR {
					g.CN[0][m-1] = 20. * 5 / 10 / (float64(m) + 1)
				} else if g.FRUCHT[0] == WRA || g.FRUCHT[0] == AB {
					g.CN[0][m-1] = 45. * 5 / 10 / (float64(m) + 1)
				} else if g.FRUCHT[0] == CCM || g.FRUCHT[0] == M || g.FRUCHT[0] == SM {
					g.CN[0][m-1] = 95. * 5 / 10 / (float64(m) + 1)
				} else if g.FRUCHT[0] == K {
					g.CN[0][m-1] = 50. * 5 / 10 / (float64(m) + 1)
				} else {
					g.CN[0][m-1] = 35. * 5 / 10 / (float64(m) + 1)
				}
			}
			// ! ********************** Messwertdatei lesen ***********************************
			// ! INPUTS:
			// ! NMESS                     = Zähler für Messereignisse
			// ! MES$(NMESS)               = Datum der Messung (TTMMJJ)
			// ! KONZ1
			// LET OBS$ = PATH$ & "init_" & locid$ & ".txt"
			obs := hPath.obs
			_, scannerObserv, _ := Open(&FileDescriptior{FilePath: obs, FileDescription: "observation file", UseFilePool: true})
			Fident := getFident(g, FLAEID)
			LineInut(scannerObserv)
			g.NMESS = 0
			for SCHLAG, OBSERtoken, valid := NextLineInut(0, scannerObserv, strings.Fields); valid; SCHLAG, OBSERtoken, valid = NextLineInut(0, scannerObserv, strings.Fields) {
				for ok := SCHLAG == Fident; ok; ok = SCHLAG == Fident && valid {
					g.NMESS++
					l.MK[g.NMESS-1] = OBSERtoken[1]
					if g.NMESS == 1 {
						g.MES[g.NMESS-1] = l.MK[g.NMESS-1]
						_, g.MESS[g.NMESS-1] = g.Datum(g.MES[g.NMESS-1])

						//! +++++++++++++++ Ueberschreiben des Erntedatums der Vorfrucht aus der Rotationsdatei ++++++++++++++++++++
						//if g.AUTOHAR {
						// commentented out by Christians newest version
						// var ERNDAT int
						// ERNDAT, gloInput.ERNTE[0] = Datum(gloInput.MES[gloInput.NMESS-1], gloInput.CENT)
						// gloInput.ITAG = ERNDAT
						// gloInput.BEGINN = gloInput.ERNTE[0]
						//}
						if g.AUTOFERT {
							if g.ORGTIME[0] == "H" {
								g.ZTDG[0] = g.ERNTE[0] + 1
							}
						}

						l.KONZ1 = ValAsFloat(OBSERtoken[2], obs, OBSERtoken[2])
						l.KONZ3 = ValAsFloat(OBSERtoken[3], obs, OBSERtoken[3])
						l.KONZ4 = ValAsFloat(OBSERtoken[4], obs, OBSERtoken[4])
						if len(OBSERtoken) > 9 {
							l.KONZ5 = ValAsFloat(OBSERtoken[9], obs, OBSERtoken[9])
							l.KONZ6 = ValAsFloat(OBSERtoken[10], obs, OBSERtoken[10])
							l.KONZ7 = ValAsFloat(OBSERtoken[11], obs, OBSERtoken[11])
						}
						l.Jstr = OBSERtoken[5]
						winit[0] = ValAsFloat(OBSERtoken[6], obs, OBSERtoken[6])
						winit[1] = ValAsFloat(OBSERtoken[7], obs, OBSERtoken[7])
						winit[2] = ValAsFloat(OBSERtoken[8], obs, OBSERtoken[8])

						if len(OBSERtoken) > 9 {
							winit[3] = ValAsFloat(OBSERtoken[12], obs, OBSERtoken[12])
							winit[4] = ValAsFloat(OBSERtoken[13], obs, OBSERtoken[13])
							winit[5] = ValAsFloat(OBSERtoken[14], obs, OBSERtoken[14])
						}
						if g.MES[0] != "------" {
							for zi := 1; zi <= g.N; zi++ {
								ziIndex := zi - 1
								if zi < 4 {
									if l.Jstr == "3" {
										g.WG[g.NMESS+1][ziIndex] = winit[0]
									} else if l.Jstr == "2" {
										g.WG[g.NMESS+1][ziIndex] = winit[0] * 1.4
									} else {
										g.WG[g.NMESS+1][ziIndex] = g.WMIN[ziIndex] + (g.W[ziIndex]-g.WMIN[ziIndex])*winit[0]
									}
								} else if zi > 3 && zi < 7 {
									if l.Jstr == "3" {
										g.WG[g.NMESS+1][ziIndex] = winit[1]
									} else if l.Jstr == "2" {
										g.WG[g.NMESS+1][ziIndex] = winit[1] * 1.5
									} else {
										g.WG[g.NMESS+1][ziIndex] = g.WMIN[ziIndex] + (g.W[ziIndex]-g.WMIN[ziIndex])*winit[1]
									}
								} else if zi > 6 && zi < 10 {
									if l.Jstr == "3" {
										g.WG[g.NMESS+1][ziIndex] = winit[2]
									} else if l.Jstr == "2" {
										g.WG[g.NMESS+1][ziIndex] = winit[2] * 1.6
									} else {
										g.WG[g.NMESS+1][ziIndex] = g.WMIN[ziIndex] + (g.W[ziIndex]-g.WMIN[ziIndex])*winit[2]
									}
								} else if zi > 9 && zi < 13 {
									if l.Jstr == "3" {
										g.WG[g.NMESS+1][ziIndex] = winit[3]
									} else if l.Jstr == "2" {
										g.WG[g.NMESS+1][ziIndex] = winit[3] * 1.6
									} else {
										g.WG[g.NMESS+1][ziIndex] = g.WMIN[ziIndex] + (g.W[ziIndex]-g.WMIN[ziIndex])*winit[3]
									}
								} else if zi > 12 && zi < 16 {
									if l.Jstr == "3" {
										g.WG[g.NMESS+1][ziIndex] = winit[4]
									} else if l.Jstr == "2" {
										g.WG[g.NMESS+1][ziIndex] = winit[4] * 1.6
									} else {
										g.WG[g.NMESS+1][ziIndex] = g.WMIN[ziIndex] + (g.W[ziIndex]-g.WMIN[ziIndex])*winit[4]
									}
								} else if zi > 15 {
									if l.Jstr == "3" {
										g.WG[g.NMESS+1][ziIndex] = winit[5]
									} else if l.Jstr == "2" {
										g.WG[g.NMESS+1][ziIndex] = winit[5] * 1.6
									} else {
										g.WG[g.NMESS+1][ziIndex] = g.WMIN[ziIndex] + (g.W[ziIndex]-g.WMIN[ziIndex])*winit[5]
									}

									//g.WG[g.NMESS+1][ziIndex] = winit[5]
								}
							}
							g.WG[g.NMESS+1][g.N] = g.WG[g.NMESS+1][g.N-1]
							if g.NMESS == 1 {
								if l.Jstr == "3" {
									g.WNZ[0] = (winit[0] + winit[1] + winit[2]) * 300
								} else if l.Jstr == "2" {
									g.WNZ[0] = (winit[0]*1.4 + winit[1]*1.5 + winit[2]*1.6) * 300
								} else {
									g.WNZ[0] = (g.WG[2][0] + g.WG[2][1] + g.WG[2][2] + g.WG[2][3] + g.WG[2][4] + g.WG[2][5] + g.WG[2][6] + g.WG[2][7] + g.WG[2][8]) * 100
								}
							}
							g.KNZ1[0] = l.KONZ1
							g.KNZ2[0] = l.KONZ3
							g.KNZ3[0] = l.KONZ4
							g.KNZ4[0] = l.KONZ5
							g.KNZ5[0] = l.KONZ6
							g.KNZ6[0] = l.KONZ7
							for i := 1; i <= g.N; i++ {
								iIndex := i - 1
								if i < 4 {
									g.CN[g.NMESS][iIndex] = g.KNZ1[0] / 3
								} else if i > 3 && i < 7 {
									g.CN[g.NMESS][iIndex] = g.KNZ2[0] / 3
								} else if i > 6 && i < 10 {
									g.CN[g.NMESS][iIndex] = g.KNZ3[0] / 3
								} else if i > 9 && i < 13 {
									g.CN[g.NMESS][iIndex] = g.KNZ4[0] / 3
								} else if i > 12 && i < 16 {
									g.CN[g.NMESS][iIndex] = g.KNZ5[0] / 3
								} else {
									g.CN[g.NMESS][iIndex] = g.KNZ6[0] / 5
								}
							}
						} else {
							g.MESS[0] = g.BEGINN
							for i := 0; i < g.N; i++ {
								g.CN[1][i] = g.CN[0][i]
							}
						}
					}

					SCHLAG, OBSERtoken, valid = NextLineInut(0, scannerObserv, strings.Fields)

				}
			}

			for i := 0; i < g.N+1; i++ {
				if g.CN[0][i] < 0 {
					g.CN[0][i] = .1
				}
			}
			// read Smin observed data
			readSmin(g, FLAEID, hPath)

			// ! ********************** Bodenbearbeitungsmassnahmen lesen ***********************************
			til := hPath.til
			_, scannertilage, err := Open(&FileDescriptior{FilePath: til, FileDescription: "tillage file", UseFilePool: false, ContinueOnError: true})
			if err == nil {
				LineInut(scannertilage)
				LineInut(scannertilage)
				NRTIL := 0

				for SCHLAG, tilageTokens, valid := NextLineInut(0, scannertilage, strings.Fields); valid; SCHLAG, tilageTokens, valid = NextLineInut(0, scannertilage, strings.Fields) {
					if SCHLAG == g.PKT {
						for ok := true; ok; ok = SCHLAG == g.PKT && valid {
							// Tokens: Schlag/FieldID(0) depth(1) type(2) date(3)
							NRTIL++
							NRTILindex := NRTIL - 1
							g.TILDAT[NRTILindex] = tilageTokens[3]
							g.EINT[NRTILindex] = ValAsFloat(tilageTokens[1], til, tilageTokens[1])
							g.TILART[NRTILindex] = int(ValAsInt(tilageTokens[2], til, tilageTokens[2]))
							_, valEinte := g.Datum(g.TILDAT[NRTILindex])
							g.EINTE[NRTIL] = valEinte
							if g.EINTE[NRTIL] < g.BEGINN {
								NRTIL--
							}
							SCHLAG, tilageTokens, valid = NextLineInut(0, scannertilage, strings.Fields)
						}
					}
				}
				// move overlapping dates
				for i := 1; i <= NRTIL; i++ {
					if g.EINTE[i+1] == g.EINTE[i] {
						g.EINTE[i+1] = g.EINTE[i+1] + 1
					}
				}
			}
			// calculate initial n residue
			residi(g, hPath)
			// calculate inital S residue
			sResidi(g, hPath)
			if !g.AUTOFERT {
				// ! ********************** Düngungsmassnahmen lesen ***********************************
				dun := hPath.dun
				_, scannerFert, _ := Open(&FileDescriptior{FilePath: dun, FileDescription: "fertilization file", UseFilePool: true})
				LineInut(scannerFert)
				NDu := 1
				for SCHLAG, fertilizerToken, valid := NextLineInut(0, scannerFert, strings.Fields); valid; SCHLAG, fertilizerToken, valid = NextLineInut(0, scannerFert, strings.Fields) {

					for ok := SCHLAG == g.PKT; ok; ok = SCHLAG == g.PKT && valid {

						// try to stay backwards compatible with old format
						//Field_ID(0)  N(1)   Frt(2) date(3)
						headerN := 1
						headerS := -1
						headerFrt := 2
						headerDate := 3
						// or Field_ID(0)  N(1) S(2)  Frt(3) date(4)
						if len(fertilizerToken) == 5 {
							headerN = 1
							headerS = 2
							headerFrt = 3
							headerDate = 4
						}
						NDu++
						NDuindex := NDu - 1

						l.NDGMG[NDuindex] = ValAsFloat(fertilizerToken[headerN], dun, fertilizerToken[headerN]) * g.DUNGSZEN
						l.SDGMG[NDuindex] = 0
						if headerS > 0 {
							l.SDGMG[NDuindex] = ValAsFloat(fertilizerToken[headerS], dun, fertilizerToken[headerS]) * g.DUNGSZEN
						}
						g.DGART[NDuindex] = fertilizerToken[headerFrt]
						DGDAT := fertilizerToken[headerDate]

						_, valztdg := g.Datum(DGDAT)
						g.ZTDG[NDuindex] = valztdg
						if g.ZTDG[NDuindex] < g.BEGINN {
							NDu--
						}
						SCHLAG, fertilizerToken, valid = NextLineInut(0, scannerFert, strings.Fields)
					}

				}
				for i := 1; i <= NDu; i++ {
					index := i - 1
					if g.ZTDG[index+1] == g.ZTDG[index] {
						g.ZTDG[index+1] = g.ZTDG[index+1] + 1
					}
				}
				for i := 1; i < NDu; i++ {
					dueng(i, g, l, hPath)
				}
			}
			break
		}
	}
	nPotMin(g, l)
	sPotMin(g)
	return nil
}

func getFident(g *GlobalVarsMain, FLAEID string) (Fident string) {
	if g.INIWAHL == 1 {
		Fident = "ALLE"
	} else if g.INIWAHL == 2 {
		Fident = g.PKT
	} else if g.INIWAHL == 3 {
		Fident = FLAEID
	}
	return Fident
}

// Hydro reads hydro parameter
func Hydro(las1 int, g *GlobalVarsMain, local *InputSharedVars, hPath *HFilePath) (BDART string) {
	lIndex := las1 - 1
	BDART = g.BART[lIndex]
	if las1 == g.AZHO {
		_, scannerParCap, _ := Open(&FileDescriptior{FilePath: hPath.parcap, UseFilePool: true})

		for ok := true; ok; ok = g.CAPS[0] == 0 {
			PARA := LineInut(scannerParCap)
			PARA2 := LineInut(scannerParCap)
			if PARA[0:3] == BDART {
				for i := 1; i <= 10; i++ {
					indexStart := i*6 - 1
					endIndex := i*6 + 4
					if len(PARA) <= endIndex {
						endIndex = len(PARA)
					}
					g.CAPS[i-1] = ValAsFloat(PARA[indexStart:endIndex], "none", PARA)
					if len(PARA2) <= endIndex {
						endIndex = len(PARA2)
					}
					g.CAPS[i+10-1] = ValAsFloat(PARA2[indexStart:endIndex], "none", PARA2)
				}
			}
		}
	}
	g.IZM = 30

	hyparName := hPath.hypar
	_, scannerHyPar, _ := Open(&FileDescriptior{FilePath: hyparName, UseFilePool: true})

	for {
		wa := LineInut(scannerHyPar)
		if wa[0:3] == BDART {
			if g.LD[lIndex] == 1 || g.LD[lIndex] == 2 {
				local.FK[lIndex] = ValAsFloat(wa[4:6], hyparName, wa) / 100
				g.LIM[lIndex] = local.FK[lIndex] - ValAsFloat(wa[13:15], hyparName, wa)/100
				g.PRGES[lIndex] = ValAsFloat(wa[22:24], hyparName, wa) / 100
			} else if g.LD[lIndex] == 3 {
				local.FK[lIndex] = ValAsFloat(wa[7:9], hyparName, wa) / 100
				g.LIM[lIndex] = local.FK[lIndex] - ValAsFloat(wa[16:18], hyparName, wa)/100
				g.PRGES[lIndex] = ValAsFloat(wa[25:27], hyparName, wa) / 100
			} else if g.LD[lIndex] == 4 || g.LD[lIndex] == 5 {
				local.FK[lIndex] = ValAsFloat(wa[10:12], hyparName, wa) / 100
				g.LIM[lIndex] = local.FK[lIndex] - ValAsFloat(wa[19:21], hyparName, wa)/100
				g.PRGES[lIndex] = ValAsFloat(wa[28:30], hyparName, wa) / 100
			}

			g.WUMAX[lIndex] = ValAsFloat(wa[31:33], hyparName, wa)
			if las1 == 1 {
				if BDART[0:1] == "S" {
					g.WRED = (g.LIM[lIndex] + 0.6*(local.FK[lIndex]-g.LIM[lIndex])) * 100
				} else {
					g.WRED = (g.LIM[lIndex] + 0.66*(local.FK[lIndex]-g.LIM[lIndex])) * 100
				}
			}
			break
		}
	}
	var KRR, KRG float64
	if BDART[0] == 'S' {
		g.AD = .004
		if g.GW < 9 {
			KRR = 2
		} else if g.GW >= 20 && g.GW < 30 {
			KRR = -1
		} else if g.GW >= 30 {
			KRR = -2
		} else {
			KRR = 0
		}
		if las1 < 2 {
			g.IZM = 30
			g.PROP = 0.6
		}
		if BDART == "SL2" || BDART[1] == 'U' || BDART == "SG " || BDART == "SM " || BDART == "SF " {
			if g.CGEHALT[lIndex] > 4.6 {
				KRR = KRR + 10
				KRG = 10
			} else if g.CGEHALT[lIndex] > 2.3 {
				KRR = KRR + 7.5
				KRG = 6.5
			} else if g.CGEHALT[lIndex] > 1.16 {
				KRR = KRR + 3.5
				KRG = 2.5
			} else {
				KRG = 0
			}
		} else {
			if g.CGEHALT[lIndex] > 4.6 {
				KRR = KRR + 11.5
				KRG = 14
			} else if g.CGEHALT[lIndex] > 2.3 {
				KRR = KRR + 8
				KRG = 10
			} else if g.CGEHALT[lIndex] > 1.16 {
				KRR = KRR + 3.5
				KRG = 4.5
			} else if g.CGEHALT[lIndex] > 0.58 {
				KRR = KRR + 1.5
				KRG = 1.50
			} else {
				KRG = 0
			}
		}

		if BDART[1] == 'U' || BDART[1] == 'u' {
			g.IZM = 30
		} else if BDART[1] == 'L' || BDART[1] == 'l' {
			if BDART[2] == '2' {
				if las1 == 1 {
					g.IZM = 30
				}
			}
		} else if BDART[1] == 'F' || BDART[1] == 'f' {
			if las1 == 1 {
				g.IZM = 30
			}
		} else if BDART[1] == 'G' || BDART[1] == 'g' {
			if las1 == 1 {
				g.IZM = 30
			}
		} else if BDART[1] == 'M' || BDART[1] == 'm' {
			if las1 == 1 {
				g.IZM = 30
			}
		}
	} else if BDART[0] == 'U' {
		g.AD = .002
		if g.GW < 8 {
			KRR = 1
		} else if g.GW > 35 {
			KRR = -1
		} else {
			KRR = 0
		}
		if las1 < 2 {
			g.PROP = 0.3
		}
		if g.CGEHALT[lIndex] > 5.2 {
			KRR = KRR + 12
		} else if g.CGEHALT[lIndex] > 4.6 {
			KRR = KRR + 7
		} else if g.CGEHALT[lIndex] > 3.5 {
			KRR = KRR + 5
		} else if g.CGEHALT[lIndex] > 2.3 {
			KRR = KRR + 1
		}
		if BDART[1] == 'S' || BDART[1] == 's' {
			if las1 == 1 {
				g.IZM = 30
			}
		} else if BDART[1] == 'T' || BDART[1] == 't' {
			if las1 == 1 {
				g.IZM = 20
			}
		}
	} else if BDART[0] == 'L' {
		g.AD = .005
		if g.GW < 8 {
			KRR = 1
		} else {
			KRR = 0
		}

		if las1 < 2 {
			g.PROP = 0.3
		}
		if g.CGEHALT[lIndex] > 4.6 {
			KRR = KRR + 7
		} else if g.CGEHALT[lIndex] > 3.5 {
			KRR = KRR + 4
		} else if g.CGEHALT[lIndex] > 2.3 {
			KRR = KRR + 1
		}
		if BDART[1] == 'S' || BDART[1] == 's' {
			if las1 == 1 {
				g.IZM = 30
			}
		} else if BDART[1] == 'T' || BDART[1] == 't' {
			if BDART[2] == '2' {
				if las1 == 1 {
					g.IZM = 30
				}
			} else if BDART[2] == '3' {
				if las1 == 1 {
					g.IZM = 20
				}
			} else if BDART[2] == 'U' {
				if las1 == 1 {
					g.IZM = 20
				}
			} else if BDART[2] == 'S' {
				if las1 == 1 {
					g.IZM = 20
				}
			}
		}
	} else if BDART == "T  " || BDART == " T " || BDART == "  T" {
		if las1 == 1 {
			g.IZM = 20
		}
		g.AD = .001
		if g.GW < 8 {
			KRR = 1
		} else {
			KRR = 0
		}
		if las1 < 2 {
			g.PROP = 0.4
		}
	} else if BDART[0] == 'T' {
		if las1 == 1 {
			g.IZM = 20
		}
		g.AD = .001
		if g.GW < 8 {
			KRR = 1
		} else {
			KRR = 0
		}
		if las1 < 2 {
			g.PROP = 0.4
		}
		if BDART[1] == 'U' || BDART[1] == 'u' {
			if BDART[2] == '2' {
				if g.CGEHALT[lIndex] > 4.6 {
					KRR = KRR + 4
				} else if g.CGEHALT[lIndex] > 3.5 {
					KRR = KRR + 2
				}
			} else if BDART[2] == '3' {
				if g.CGEHALT[lIndex] > 4.6 {
					KRR = KRR + 4
				} else if g.CGEHALT[lIndex] > 3.5 {
					KRR = KRR + 2
				}
			} else if BDART[2] == '4' {
				if g.CGEHALT[lIndex] > 4.6 {
					KRR = KRR + 7
				} else if g.CGEHALT[lIndex] > 3.5 {
					KRR = KRR + 4
				} else if g.CGEHALT[lIndex] > 2.3 {
					KRR = KRR + 1
				}
			}
		}
	} else if BDART[0] == 'H' {
		if las1 < 2 {
			g.IZM = 10
			g.PROP = 0.1
		}
	}
	g.FELDW[lIndex] = local.FK[lIndex] + KRR/100
	g.NORMFK[lIndex] = local.FK[lIndex]
	g.PRGES[lIndex] = g.PRGES[lIndex] + KRG/100
	return BDART
}

// residi loads potential mineralization from previous crops
func residi(g *GlobalVarsMain, hPath *HFilePath) {
	//   Mineralisationspotentiale aus Vorfruchtresiduen
	// "CROP_N.TXT"
	cropN := hPath.cropn
	_, scanner, _ := Open(&FileDescriptior{FilePath: cropN, UseFilePool: true})

	var KOSTRO, NERNT, NKOPP, NWURA, NFAST float64
	for scanner.Scan() {
		CROP := scanner.Text()
		if g.ToCropType(CROP[0:3]) == g.FRUCHT[g.AKF.Index] {
			KOSTRO = ValAsFloat(CROP[4:7], cropN, CROP)
			NERNT = ValAsFloat(CROP[13:18], cropN, CROP)
			NKOPP = ValAsFloat(CROP[25:30], cropN, CROP)
			NWURA = ValAsFloat(CROP[36:40], cropN, CROP)
			NFAST = ValAsFloat(CROP[41:45], cropN, CROP)
			break
		}
	}

	AUFGES := (g.ERTR[0]*NERNT + g.ERTR[0]*KOSTRO*NKOPP) / (1 - NWURA)
	var DGM float64
	if g.JN[0] == 0 {
		if g.EINT[0] == 0 {
			DGM = 0
		} else {
			DGM = AUFGES - (g.ERTR[0] * NERNT)
		}
	} else if g.JN[0] == 1 {
		DGM = AUFGES * NWURA
	} else {
		DGM = AUFGES*NWURA + (1-g.JN[0])*(AUFGES-g.ERTR[0]*NERNT-AUFGES*NWURA)
	}
	if DGM < 0 {
		DGM = 0
	}
	g.NSAS[0] = DGM * NFAST
	g.NLAS[0] = DGM * (1 - NFAST)
	g.NDIR[0] = 0.0
}

func nPotMin(g *GlobalVarsMain, l *InputSharedVars) {
	if g.CGEHALT[0] > 14 {
		g.NALTOS = 5000 * l.NGEHALT[0] * g.NAKT * float64(g.UKT[1])
	} else if g.CGEHALT[0] > 5 {
		g.NALTOS = 10600 * l.NGEHALT[0] * g.NAKT * float64(g.UKT[1])
	} else if g.CGEHALT[0] < 1 {
		g.NALTOS = 15000 * l.NGEHALT[0] * g.NAKT * float64(g.UKT[1])
	} else {
		g.NALTOS = 15000 * l.NGEHALT[0] * g.NAKT * float64(g.UKT[1])
	}
}

func verdun(gloInput *GlobalVarsMain, hPath *HFilePath) {
	if gloInput.ETMETH == 1 {
		//   ! Read Haude/Heger factors
		filename := hPath.evapo
		_, scanner, _ := Open(&FileDescriptior{FilePath: filename, FileDescription: "evapo file", UseFilePool: true})

		LineInut(scanner)
		if gloInput.AKF.Num == 1 {
			HAUF := LineInut(scanner)
			for i := 1; i <= 12; i++ {
				gloInput.FKU[i-1] = ValAsFloat(HAUF[5*i-1:3+5*i], filename, HAUF)
			}
		} else {
			for scanner.Scan() {
				HAUF := scanner.Text()
				if gloInput.ToCropType(HAUF[0:3]) == gloInput.FRUCHT[gloInput.AKF.Index] {
					for i := 1; i <= 12; i++ {
						gloInput.FKF[i-1] = ValAsFloat(HAUF[5*i-1:3+5*i], filename, HAUF)
					}
				}
			}
		}
	}

}

func dueng(i int, g *GlobalVarsMain, l *InputSharedVars, hPath *HFilePath) {
	// "FERTILIZ.TXT"
	dungfile := hPath.dung
	_, scanner, _ := Open(&FileDescriptior{FilePath: dungfile, FileDescription: "fertilization file", UseFilePool: true})
	for scanner.Scan() {
		du := scanner.Text()
		token := strings.Fields(du)
		if token[0] == g.DGART[i] {
			l.NORG[i] = ValAsFloat(token[1], dungfile, du)                                      //Ntot
			VOL := ValAsFloat(token[6], dungfile, du)                                           // Loss
			g.NDIR[i] = l.NDGMG[i] * l.NORG[i] * ValAsFloat(token[2], dungfile, du)             // NDIR
			g.NH4N[i] = g.NDIR[i] * ValAsFloat(token[5], dungfile, du) * (1 - VOL)              // Neu: nicht Nitrat-N in Dünger NH4N
			g.NDIR[i] = g.NDIR[i] - g.NDIR[i]*ValAsFloat(token[5], dungfile, du)*VOL            // NH4
			g.NSAS[i] = (l.NDGMG[i]*l.NORG[i] - g.NDIR[i]) * ValAsFloat(token[3], dungfile, du) // Nfst
			g.NLAS[i] = (l.NDGMG[i]*l.NORG[i] - g.NDIR[i]) * ValAsFloat(token[4], dungfile, du) // Nslo

			// sulfur in fertilizer
			if g.Sulfonie {
				//LET SO4   = VAL(DUNG$(18:21))
				SO4 := ValAsFloat(token[7], dungfile, du)
				//LET SORG  = 1- SO4  !VAL(DUNG$(27:30))
				SORG := 1 - SO4
				//LET SFAST = VAL(DUNG$(23:26))
				SFAST := ValAsFloat(token[8], dungfile, du)
				//LET SDIR(I) = DGMG(I) * SO4
				g.SDIR[i] = l.SDGMG[i] * SO4
				//LET SSAS(I) = DGMG(I) * SORG * SFAST
				g.SSAS[i] = l.SDGMG[i] * SORG * SFAST
				//LET SLAS(I) = DGMG(I) * SORG * (1-SFAST)
				g.SLAS[i] = l.SDGMG[i] * SORG * (1 - SFAST)
			}
			break
		}
	}
}

//TODO: add sulfur fertilization
// SUB DUENG(I)
//         LET DU$ = "DUENGER.TXT"
//         WHEN ERROR IN
//              OPEN #4:Name DU$,ACCESS INPUT,ORGANIZATION TEXT
//         USE
//              PRINT "Datei Duengercharakteristik ";DU$;" nicht gefunden "
//              STOP
//         END WHEN
//         DO while more #4
//            LINE INPUT #4: DUNG$
//            IF DUNG$(1:1) = DGART$(I) then
//               LET SO4   = VAL(DUNG$(18:21))
//               LET SORG  = 1- SO4  !VAL(DUNG$(27:30))
//               LET SFAST = VAL(DUNG$(23:26))
//               LET SDIR(I) = DGMG(I) * SO4
//               LET SSAS(I) = DGMG(I) * SORG * SFAST
//               LET SLAS(I) = DGMG(I) * SORG * (1-SFAST)
//               EXIT DO
//            END IF
//         LOOP
//     END SUB
