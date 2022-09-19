package hermes

import "fmt"

// DevelopmentStage for wheat
type DevelopmentStage int

const (
	invalidState DevelopmentStage = iota
	schossen
	grperiode
	reife
	abreife
	aehrenschieben
)

func (s DevelopmentStage) String() string {
	return developmentStageToString[s]
}

var developmentStageToString = map[DevelopmentStage]string{
	schossen:       "SCHOSSEN",
	grperiode:      "Gr.Periode",
	reife:          "Reife",
	abreife:        "ABREIFE",
	aehrenschieben: "AEHRENSCHIEBEN",
	invalidState:   "manueller Abbruch",
}

// CalulateDevelopmentStages for wheat
func CalulateDevelopmentStages(zeit int, FV, FP float64, g *GlobalVarsMain) {
	if g.SUM[1] >= g.TSUM[1] {
		if g.DOUBLE == 0 {
			//!  Datum des Doppelringstadiums
			g.DOUBLE = zeit
			g.DOPP = g.Kalender(zeit)
		}
	}
	if zeit > g.P2+10 {
		g.SUMAE = g.SUMAE + (g.TEMP[g.TAG.Index]-g.BAS[g.INTWICK.Index])*FV*FP*g.DT.Num
		if g.SUMAE > 130 && g.ASIP == 0 {
			g.ASIP = zeit
			g.AEHR = g.Kalender(g.ASIP)
		}
	} else {
		g.SUMAE = 0
	}
	if g.SUM[3] >= g.TSUM[3] && g.BLUET == 0 {
		g.BLUET = zeit
		g.BLUEH = g.Kalender(g.BLUET)
	}
	if g.SUM[4] >= g.TSUM[4] && g.REIF == 0 {
		g.REIF = zeit
		g.REIFE = g.Kalender(g.REIF)
	}
}

// ResetStages for wheat
func ResetStages(g *GlobalVarsMain) {
	g.DOUBLE, g.ASIP, g.BLUET, g.REIF, g.ENDPRO = 0, 0, 0, 0, 0
}

// SimulateFertilizationAfterPrognose simulates fertilization as required and calcules new dates
func SimulateFertilizationAfterPrognose(zeit int, DTGESN, SUMDIFF, TRNSUM float64, g *GlobalVarsMain) {
	var ANGEBOT, BEDARF float64
	if zeit > g.PROGNOS {
		ANGEBOT = SUMDIFF + TRNSUM
		if ANGEBOT < DTGESN {
			BEDARF = DTGESN - ANGEBOT
			if (g.C1[0]+BEDARF)/(g.WG[0][0]*g.DZ.Num)*10 < 200 {
				g.C1[0] = g.C1[0] + BEDARF
			} else {
				BEDARF = 200*g.WG[0][0]*g.DZ.Num/10 - g.C1[0]
				if BEDARF < 0 {
					BEDARF = 0
				}
				g.C1[0] = g.C1[0] + BEDARF
			}
			g.DUNGBED = g.DUNGBED + BEDARF
			if g.DEFDAT == 0 {
				g.DEFDAT = zeit
			}
		}
		// check for next prognose period
		if g.ENDE >= g.ERNTE[g.AKF.Index] {
			if g.REIF != 0 {
				g.ENDE = zeit
				g.ENDPRO = zeit
				g.ENDSTADIUM = abreife
			} else if g.ASIP != 0 {
				if g.ASIP-g.PROGNOS > 7 {
					g.ENDE = zeit
					g.ENDPRO = zeit
					g.ENDSTADIUM = aehrenschieben
				}
			}
		}
	}
}

// OnDoubleRidgeStateNotReached force DoubleRidgeState, move P1 date
func OnDoubleRidgeStateNotReached(zeit int, g *GlobalVarsMain) {
	if g.DOUBLE == 0 {
		g.SUM[1] = g.TSUM[1]
		g.PHYLLO = g.TSUM[1]
		g.DOUBLE = zeit
		g.P1 = g.P1 + 4
	}
}

// SetPrognoseDate set prognose date, if value was valid
func SetPrognoseDate(prog string, g *GlobalVarsMain) (PR bool) {
	g.PROGDAT = prog
	if prog[1] == '-' {
		PR = false
		g.PROGNOS = g.ENDE + 1
	} else {
		_, g.PROGNOS = g.Datum(g.PROGDAT)
		if g.PROGNOS < g.ENDE {
			PR = true
		}
	}
	return PR
}

// PrognoseTime triggers loading of weather prognose Data and prepares prognose
func PrognoseTime(ZEIT int, g *GlobalVarsMain, herPath *HFilePath, driConfig *Config) {

	// overwrite weather data with prognosed weather data (optional? if not given no)
	VWDAT := herPath.vwdatnrm
	year, _, _ := KalenderDate(ZEIT)
	s := NewWeatherDataShared(1, g.CO2KONZ)
	err := WetterK(VWDAT, year, g, &s, herPath, driConfig)
	if err != nil {
		if g.DEBUGCHANNEL != nil {
			g.DEBUGCHANNEL <- fmt.Sprintln(err)
		} else {
			fmt.Println(err)
		}
	} else {
		LoadYear(g, &s, year)
	}

	for i := 0; i < g.N; i++ {
		g.CA[i] = g.C1[i]
	}
	g.MINA = g.MINSUM //obsolete?
	g.PLANA = g.PESUM
	g.OUTA = g.OUTSUM //obsolete?

	if g.DOUBLE == 0 {
		if ZEIT < g.P1-6 {
			g.ENDE = g.P1
			g.ENDSTADIUM = schossen
			g.ENDPRO = g.P1
		} else {
			g.P1 = g.P1 + 4
			g.P2 = g.P2 + 2
			g.ENDSTADIUM = grperiode

			g.ENDE = g.P2
			g.ENDPRO = g.P2
			g.SUM[1] = g.TSUM[1] - 4
			g.PHYLLO = g.TSUM[0] + g.SUM[1]
		}
	} else if ZEIT < g.P1-15 {
		g.ENDE = g.P1
		g.ENDSTADIUM = schossen
		g.ENDPRO = g.P1
	} else if ZEIT < g.P2-5 {
		g.ENDE = g.P2
		g.ENDSTADIUM = grperiode
		g.ENDPRO = g.P2
	} else {
		g.ENDE = g.ERNTE[g.AKF.Index]
		g.ENDSTADIUM = reife
		g.ENDPRO = g.ENDE
	}
}

// FinalDungPrognose last step for fertilization recommendation
func FinalDungPrognose(g *GlobalVarsMain) (NAPP int) {
	//  --------------- Calculation of fertilization recommendation ---------------
	if g.ENDE == g.P1 {
		if g.DEFDAT > 0 {
			if g.DOUBLE == 0 || g.DOUBLE < g.PROGNOS {
				NAPP = g.DEFDAT
			} else {
				NAPP = min(g.DEFDAT, g.DOUBLE)
			}
		} else {
			NAPP = g.DOUBLE
		}
	} else if g.ENDE == g.P2 {
		if g.DEFDAT > 0 {
			if g.P1 > g.PROGNOS {
				NAPP = min(g.DEFDAT, g.P1-5)
			} else {
				NAPP = g.DEFDAT
			}
		} else {
			NAPP = g.P1 - 5
		}
	} else if g.ENDE == g.ASIP {
		if g.DEFDAT > 0 {
			if g.P2 > g.PROGNOS {
				NAPP = min(g.DEFDAT, g.P2-4)
			} else {
				NAPP = g.DEFDAT
			}
		} else {
			NAPP = g.P2
		}
	} else if g.ENDE == g.REIF {
		if g.DEFDAT > 0 {
			if g.ASIP == 0 || g.ASIP < g.PROGNOS {
				NAPP = g.DEFDAT
			} else {
				NAPP = min(g.DEFDAT, g.ASIP-3)
			}
		} else {
			NAPP = g.ASIP
		}
	} else {
		NAPP = g.DEFDAT
	}
	if NAPP < g.PROGNOS {
		NAPP = g.PROGNOS
	}
	if NAPP > 1 {
		g.NAPPDAT = g.Kalender(NAPP)
	} else {
		g.NAPPDAT = "--------"
	}
	return NAPP
}
