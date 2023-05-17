package hermes

import (
	"math"
)

// WaterSharedVars is a struct of shared variables for this module
type WaterSharedVars struct {
	LIMIT    [21]float64
	EVA      [366]float64
	NFK      [21]float64
	GWAUF    float64
	EV       [21]float64
	SATDEF   float64
	Maxinfil float64
}

// Evatra handles evapotansiration and hydrological processes in soil
func Evatra(l *WaterSharedVars, g *GlobalVarsMain, hPath *HFilePath, zeit int) {
	TPAKT := 0.
	if zeit > g.BEGINN {
		//! Uebernahme des Voretageswassergehalts am Anfang des Zeitschritts
		for i := 0; i < g.N; i++ {
			g.WG[0][i] = g.WG[1][i]
			g.Q1[i+1] = 0
		}
		g.WG[0][g.N] = g.WG[1][g.N-1]
	}
	var WUEFF [21]float64
	var VERDU [366]float64
	var VAR [21]float64
	var TRRED [21]float64
	var RADn float64
	var RADRatio float64
	var FKM int

	// ! Berücksichtigung der verschiedenen Verdunstungsfaktoren n. Heger
	if g.TAG.Num > 212 && g.TAG.Num < 244 {
		FKM = 8
	} else if g.TAG.Num > 243 && g.TAG.Num < 274 {
		FKM = 9
	} else if g.TAG.Num > 273 && g.TAG.Num < 305 {
		FKM = 10
	} else if g.TAG.Num > 304 && g.TAG.Num < 335 {
		FKM = 11
	} else if g.TAG.Num > 334 {
		FKM = 12
	} else if g.TAG.Num < 32 {
		FKM = 1
	} else if g.TAG.Num > 31 && g.TAG.Num < 60 {
		FKM = 2
	} else if g.TAG.Num > 59 && g.TAG.Num < 91 {
		FKM = 3
	} else if g.TAG.Num > 90 && g.TAG.Num < 121 {
		FKM = 4
	} else if g.TAG.Num > 120 && g.TAG.Num < 152 {
		FKM = 5
	} else if g.TAG.Num > 151 && g.TAG.Num < 182 {
		FKM = 6
	} else {
		FKM = 7
	}
	// ! ******** Berechnung der Wasserverhaeltnisse am Anfang des Zeitschritts *****************
	// ! WG(0,1)    = Wassergehalt 0-10 cm (cm^3/cm^3)
	// ! REGEN(TAG) = Niederschlag an TAG x (cm)
	// ! DZ         = Schichtdicke (10 cm)
	// ! WMIN(1)    = Wassergehalt PWP (cm^3/cm^3) aus Bodenprofildatei 1. Schicht
	// !              Austrocknungsgrenze Oberboden 1/3 PWP
	// ! W(1)       = Feldkapazität 1. Schicht ((cm^3/cm^3) (inkl. Stauwasser)
	// ! PROZ       = Anteil verfügbares Wasser für Verdunstung
	// ! REDEV      = Reduktionsfaktor für Oberflächenevaporation
	// ! NFK(I)     = Fraktion an nutzbarer Feldkapazität in Schicht I
	// ! WNOR(i)    = Norm-FK in Schicht I

	WOB := g.WG[0][0] + g.REGEN[g.TAG.Index]/g.DZ.Num

	if WOB < g.WMIN[0]/3 {
		WOB = g.WMIN[0] / 3
	}
	PROZ := (WOB - g.WMIN[0]/3) / (g.W[0] - g.WMIN[0]/3)
	if PROZ > 1 {
		PROZ = 1
	}
	var REDEV float64
	if PROZ > .33 {
		REDEV = 1 - (.1 * (1 - PROZ) / (1 - .33))
	} else if PROZ > .22 {
		REDEV = .9 - (.625 * (.33 - PROZ) / (.33 - .22))
	} else if PROZ > .2 {
		REDEV = .275 - (.225 * (.22 - PROZ) / (.22 - .2))
	} else {
		REDEV = .05 - (.05 * (.2 - PROZ) / .2)
	}
	for i := 0; i < g.N; i++ {
		l.NFK[i] = (g.WG[0][i] - g.WMIN[i]) / (g.WNOR[i] - g.WMIN[i])
		l.NFK[0] = (g.WG[0][0] + g.REGEN[g.TAG.Index]/g.DZ.Num - g.WMIN[0]) / (g.WNOR[0] - g.WMIN[0])
		if l.NFK[i] < 0 {
			l.NFK[i] = 0
		}
	}
	// ! Aufruf HAUDE Verdunstungsfaktoren für aktuelle Frucht
	if zeit == g.SAAT[g.AKF.Index] {
		verdun(g, hPath)
		g.RDTSUM, g.REDSUM = 1, 1 //TODO: obsolete remove
	}
	var TRAMAX float64
	var FK float64
	var EVMAX float64
	var ETCP float64
	// ! ----------- ET Berechnung unter Vegetation -----------------
	// ! SAAT(AKF)     = Aussaatzeitpunkt der aktuellen anstehenden Frucht(AKF)
	// ! Ernte(AKF)    = Erntezeitt der aktuellen anstehenden Frucht (AKF)
	// ! FKM           = aktueller Monat
	// ! FKF(FKM)      = Fruchtartspezifischer Verdunstungsfaktor für HAUDE Formel
	// ! ETMETH        = Auswahlnr. Verdunstungsformel
	// !                   1 = HAUDE
	// !                   2 = Turc-Wendling
	// !                   3 = Penman-Monteith
	// !                   4 = Priestley-Taylor
	// !                   5 = direktes Einlesen von ETref aus Wetterdatei
	// ! Inputs:
	// ! TEMP(TAG)        = Tagesmitteltemperatur von TAG
	// ! RAD(TAG)         = PAR von TAG (Mjoule/m^2)
	// ! LAT              = Breitengrad (, decimal)
	// ! SUND(TAG)        = Sonnenscheindauer vom TAG (hrs)
	// ! ETNULL(TAG)      = Gras-Referenzverdunstung von TAG aus Wetterdatei (mm)
	// ! VERD(TAG)	   = Sättigungsdefizit der Luft (mm Hg) um 14:00 (nur für HAUDE)
	// ! kcoa             = Faktor für Küstennähe
	// ! FKC              = Kc-faktor (Pflanzen und Entwicklungsabh�ngig aus Pflanzenparameterdatei)
	// !
	// ! Output:
	// ! VERDU(Tag)       = ETpot (pflanzenspezifisch)
	if zeit > g.SAAT[g.AKF.Index] && g.INTWICK.Num > 1 &&
		((g.ERNTE[g.AKF.Index] > 0 && zeit < g.ERNTE[g.AKF.Index]) || (g.ERNTE[g.AKF.Index] == 0 && zeit < g.ERNTE2[g.AKF.Index])) {
		FK = g.FKF[FKM-1] //!(HAUDE (Heger)-Faktor für Frucht)
		if g.ETMETH == 1 {
			VERDU[g.TAG.Index] = g.VERD[g.TAG.Index] * FK * .1 //ETp (cm) für Frucht
		} else if g.ETMETH == 2 {
			//! ETP nach Turc-Wendling
			if g.RAD[g.TAG.Index] > 0 {
				VERDU[g.TAG.Index] = (g.RAD[g.TAG.Index]*200 + 93*g.KCOA) * (g.TEMP[g.TAG.Index] + 22) / (150 * (g.TEMP[g.TAG.Index] + 123)) * g.FKC * .1
			} else {
				DL, _, _, EXT, _, _, _ := CalculateDayLenght(g.TAG.Num, g.LAT)
				EXT = EXT * 100 // ETP by Turc-Wendling requires extraterrestic radiation in J cm-2, so we multiply with 100
				var GLOB float64
				if DL > 0 {
					GLOB = EXT * (0.19 + 0.55*g.SUND[g.TAG.Index]/DL)
				} else {
					GLOB = EXT * 0.19
				}
				VERDU[g.TAG.Index] = (GLOB + 93*g.KCOA) * (g.TEMP[g.TAG.Index] + 22) / (150 * (g.TEMP[g.TAG.Index] - 1 + 123)) * g.FKC * .1
			}
		} else if g.ETMETH == 5 {
			//! Einlesen Referenzverdustung aus Datei (Spalte Saettdef)
			VERDU[g.TAG.Index] = g.ETNULL[g.TAG.Index] * g.FKC * .1
		} else if g.ETMETH == 4 {
			//  ! ----------------------- Berechnung der Referenzverdunstung Gras nach Priestley Taylor --------------
			//  ! -- Notwendige inputs: Temp(Tag) (Tagesmitteltemperatur)
			//  !                       Tmin(TAG) (Tagesminimumtemperatur)
			//  !                       Tmax(TAG) (Tagesmaximumtemperatur)
			//  !                       Alti (Standorthöhe über Normal Null m)
			//  !                       SUND(TAG) ( Sonnenscheindauer, h), alternativ zu Strahlung
			//  !                       LAT       (Breitengrad, °)
			//  !                       RAD(TAG)  (PAR Einstrahlung MJ/m^2/d)
			//  ! -- Vordefinierte Konstante
			Albedo := 0.23
			Bolz := 0.0000000049
			DL, _, _, EXT, _, _, _ := CalculateDayLenght(g.TAG.Num, g.LAT)
			RS0 := (0.75 + 0.00002*g.ALTI) * EXT
			//  ! ---- Berechnung Atmosphärendruck in kPa
			ATMPress := 101.3 * math.Pow(((293-(0.0065*g.ALTI))/293), 5.26)
			//  ! ---- Berechnung Psychrometer-Konstante
			Psych := 0.000665 * ATMPress
			//  ! ---- Berechnung Sättigungsdampfdruck bei Tmax
			Vapres := 0.6108 * math.Exp(17.27*g.TMIN[g.TAG.Index]/(g.TMIN[g.TAG.Index]+237.3))
			//  ! Deltsat = Steigung der Sättigungsdampfdruck-Temperatur Beziehung
			Deltsat := (4098. * (0.6108 * math.Exp((17.27*g.TEMP[g.TAG.Index])/(g.TEMP[g.TAG.Index]+237.3)))) / math.Pow((g.TEMP[g.TAG.Index]+237.3), 2)
			if g.RAD[g.TAG.Index] > 0 {
				//     ! Berechnung der Nettostrahlung aus Globalstrahlung
				RADRatio = g.RAD[g.TAG.Index] * 2 / RS0
				if RADRatio > 1 {
					RADRatio = 1
				}
				RADn = (1-Albedo)*g.RAD[g.TAG.Index]*2 - Bolz*(math.Pow((g.TMIN[g.TAG.Index]+273.16), 4)+math.Pow((g.TMAX[g.TAG.Index]+273.16), 4))/2*(1.35*RADRatio-0.35)*(0.34-0.14*math.Sqrt(Vapres))
			} else {
				var Glob float64
				if DL > 0 {
					Glob = EXT * (0.19 + 0.55*g.SUND[g.TAG.Index]/DL)
				} else {
					Glob = EXT * 0.19
				}
				RADRatio = 1
				if RS0 > 0 {
					RADRatio = Glob / RS0
				}
				if RADRatio > 1 {
					RADRatio = 1
				}
				RADn = (1-Albedo)*Glob - Bolz*(math.Pow((g.TMIN[g.TAG.Index]+273.16), 4)+math.Pow((g.TMAX[g.TAG.Index]+273.16), 4))/2*(1.35*RADRatio-0.35)*(0.34-0.14*math.Sqrt(Vapres))
			}
			g.ET0 = (0.408 * Deltsat * RADn) / (Deltsat + Psych) * 1.26
			if g.ET0 < 0 {
				g.ET0 = 0
			}
			VERDU[g.TAG.Index] = g.ET0 * g.FKC * 0.1
		} else if g.ETMETH == 3 {
			//  ! ----------------------- Berechnung der Referenzverdunstung Gras nach Penman-Monteith --------------
			//  ! -- Notwendige inputs: Temp(Tag) (Tagesmitteltemperatur)
			//  !                       Tmin(TAG) (Tagesminimumtemperatur)
			//  !                       Tmax(TAG) (Tagesmaximumtemperatur)
			//  !                       RH(TAG)   (relative Luftfeuchte %)
			//  !                       Alti (Standorthöhe über Normal Null m)
			//  !                       Wind(TAG) (mittl. Windgeschw. m/s)
			//  !                       WINDHI    (Messhöhe Wind, m)
			//  !                       SUND(TAG) ( Sonnenscheindauer, h), alternativ zu Strahlung
			//  !                       LAT       (Breitengrad, )
			//  !                       RAD(TAG)  (PAR Einstrahlung MJ/m^2/d)
			//  ! -- Vordefinierte Konstante
			RSTOM0 := 100.0
			g.RSTOM = 100.
			Albedo := 0.23
			Bolz := 0.0000000049
			DL, _, _, EXT, _, _, _ := CalculateDayLenght(g.TAG.Num, g.LAT)
			RS0 := (0.75 + 0.00002*g.ALTI) * EXT
			//  ! ---- Berechnung Atmosphärendruck in kPa
			ATMPress := 101.3 * math.Pow(((293-(0.0065*g.ALTI))/293), 5.26)
			//  ! ---- Berechnung Psychrometer-Konstante
			Psych := 0.000665 * ATMPress
			//  ! ---- Berechnung Sättigungsdampfdruck bei Tmax
			SatPmax := 0.6108 * math.Exp(((17.27 * g.TMAX[g.TAG.Index]) / (237.3 + g.TMAX[g.TAG.Index])))
			//  ! ---- Berechnung Sättigungsdampfdruck bei Tmin
			SatPmin := 0.6108 * math.Exp(((17.27 * g.TMIN[g.TAG.Index]) / (237.3 + g.TMIN[g.TAG.Index])))
			//  ! ---- Berechnung mittlerer Sättigungsdampfdruck
			SatP := (SatPmin + SatPmax) / 2
			//  ! ---- Berechnung aktueller Dampfdruck
			Vapres := SatP * g.RH[g.TAG.Index] / 100
			//  ! ---- Berechnung Sättigungsdefizit der Luft
			l.SATDEF = SatP * (1 - g.RH[g.TAG.Index]/100)
			//  ! Deltsat = Steigung der Sättigungsdampfdruck-Temperatur Beziehung
			Deltsat := (4098 * (0.6108 * math.Exp((17.27*g.TEMP[g.TAG.Index])/(g.TEMP[g.TAG.Index]+237.3)))) / math.Pow((g.TEMP[g.TAG.Index]+237.3), 2)
			//  ! ----- wenn Windgeschwindigkeit Messung nicht in 2m Höhe, dann Umrechnung auf 2m
			//  ! --- Berechnung des Stomatawiderstands in Abh.von CO2 ---
			stomat(l, zeit, g)
			//  ! --------------------------------------------------------
			//  ! ----------- Umrechnung Wind auf 2 m, wenn Messhöhe Wind <> 2 m -------------
			if g.WINDHI != 2 {
				g.WIND[g.TAG.Index] = g.WIND[g.TAG.Index] * (4.87 / (math.Log(67.8*g.WINDHI - 5.42)))
			}
			//  ! ----------------------------------------------------------------------------
			//  ! ----- Berechnung des aerodynamischen Widerstands ra (Raero)
			if g.WIND[g.TAG.Index] < 0.5 { // fixed lower bound for wind speed to 0.5
				g.WIND[g.TAG.Index] = 0.5
			}
			//  ! ----- Berechnung des Oberflächenwiderstands rs mit Stomatawiderstand = 100 s/m (Rsurf)
			Rsurf0 := RSTOM0 / 1.44
			Rsurf := g.RSTOM / 1.44
			if g.RAD[g.TAG.Index] > 0 {
				// ! Berechnung der Nettostrahlung aus Globalstrahlung
				RADRatio = g.RAD[g.TAG.Index] * 2 / RS0
				if RADRatio > 1 {
					RADRatio = 1
				}
				RADn = (1-Albedo)*g.RAD[g.TAG.Index]*2 - Bolz*(math.Pow((g.TMIN[g.TAG.Index]+273.16), 4)+math.Pow((g.TMAX[g.TAG.Index]+273.16), 4))/2*(1.35*RADRatio-0.35)*(0.34-0.14*math.Sqrt(Vapres))
			} else {
				var Glob float64
				if DL > 0 {
					Glob = EXT * (0.19 + 0.55*g.SUND[g.TAG.Index]/DL)
				} else {
					Glob = EXT * 0.19
				}
				RADRatio = 1
				if RS0 > 0 {
					RADRatio = Glob / RS0
				}
				if RADRatio > 1 {
					RADRatio = 1
				}

				RADn = (1-Albedo)*Glob - Bolz*(math.Pow((g.TMIN[g.TAG.Index]+273.16), 4)+math.Pow((g.TMAX[g.TAG.Index]+273.16), 4))/2*(1.35*RADRatio-0.35)*(0.34-0.14*math.Sqrt(Vapres))
			}
			//  ! ------------------------ Berechnung der Referenzevapotranspiration --------------------
			if !g.CTRANS {
				g.ET0 = ((0.408 * Deltsat * RADn) + (Psych * (900 / (g.TEMP[g.TAG.Index] + 273)) * g.WIND[g.TAG.Index] * l.SATDEF)) / (Deltsat + Psych*(1+(Rsurf0/208)*g.WIND[g.TAG.Index]))
			} else {
				g.ET0 = ((0.408 * Deltsat * RADn) + (Psych * (900 / (g.TEMP[g.TAG.Index] + 273)) * g.WIND[g.TAG.Index] * l.SATDEF)) / (Deltsat + Psych*(1+(Rsurf/208)*g.WIND[g.TAG.Index]))
			}
			if g.ET0 < 0 {
				g.ET0 = 0
			}
			VERDU[g.TAG.Index] = g.ET0 * g.FKC * 0.1
		}
		// ! -- Begrenzung Verdunstung auf 6.5 mm/Tag --
		if VERDU[g.TAG.Index] > 0.65 {
			VERDU[g.TAG.Index] = 0.65
		}
		// ! Aufteilung ETp in Ep und Tp abh. von LAI (aus Pflanzenmodell)
		EVMAX = VERDU[g.TAG.Index] * math.Exp(-.5*g.LAI) // g.LAI may be 0, in the first run
		TRAMAX = VERDU[g.TAG.Index] - EVMAX
		ETCP = VERDU[g.TAG.Index]
	} else {
		// ! Berechnungen wie oben für unbedeckten Boden
		// ! FKU(FKM)     = Haudefaktor unbedeckt
		// ! FKB          = kc Faktor unbedeckt
		// LET FK = FKU(FKM)
		FK = g.FKU[FKM-1]
		// IF ETMETH = 1 then
		if g.ETMETH == 1 {
			//LET VERDU(TAG) = VERD(TAG)*FK*.1
			VERDU[g.TAG.Index] = g.VERD[g.TAG.Index] * FK * .1
			// ELSE IF ETMETH = 2 then
		} else if g.ETMETH == 2 {
			//LET FKC = FKB       !0.65
			g.FKC = g.FKB
			//! ETP nach Turc-Wendling
			//!LET FKUE = 1
			//IF RAD(TAG) > 0 then
			if g.RAD[g.TAG.Index] > 0 {
				//LET VERDU(TAG) = (RAD(TAG)*200+93 * kcoa) *(TEMP(TAG)+22)/(150*(TEMP(TAG)+123)) *FKC*.1
				VERDU[g.TAG.Index] = (g.RAD[g.TAG.Index]*200 + 93*g.KCOA) * (g.TEMP[g.TAG.Index] + 22) / (150 * (g.TEMP[g.TAG.Index] + 123)) * g.FKC * .1
				//ELSE
			} else {
				DL, _, _, EXT, _, _, _ := CalculateDayLenght(g.TAG.Num, g.LAT)
				EXT = EXT * 100 // ETP by Turc-Wendling requires extraterrestic radiation in J cm-2, so we multiply with 100
				var Glob float64
				if DL > 0 {
					Glob = EXT * (0.19 + 0.55*g.SUND[g.TAG.Index]/DL)
				} else {
					Glob = EXT * 0.19
				}
				VERDU[g.TAG.Index] = (Glob + 93*g.KCOA) * (g.TEMP[g.TAG.Index] + 22) / (150 * (g.TEMP[g.TAG.Index] + 123)) * g.FKC * .1
			}
		} else if g.ETMETH == 5 { // external read from weather file
			g.FKC = g.FKB
			//! ETP nach Penman (unfertig)
			VERDU[g.TAG.Index] = g.ETNULL[g.TAG.Index] * g.FKC * .1
		} else if g.ETMETH == 4 {
			g.FKC = g.FKB
			//! ----------------------- Berechnung der Referenzverdunstung Gras nach Priestley Taylor --------------
			//! -- Notwendige inputs: Temp(Tag) (Tagesmitteltemperatur)
			//!                       Tmin(TAG) (Tagesminimumtemperatur)
			//!                       Tmax(TAG) (Tagesmaximumtemperatur)
			//!                       Alti (Standorthöhe über Normal Null m)
			//!                       SUND(TAG) ( Sonnenscheindauer, h), alternativ zu Strahlung
			//!                       LAT       (Breitengrad)
			//!                       RAD(TAG)  (PAR Einstrahlung MJ/m^2/d)
			//! -- Vordefinierte Konstante
			Albedo := 0.23
			Bolz := 0.0000000049
			DL, _, _, EXT, _, _, _ := CalculateDayLenght(g.TAG.Num, g.LAT)
			RS0 := (0.75 + 0.00002*g.ALTI) * EXT
			Vapres := 0.6108 * math.Exp(17.27*g.TMIN[g.TAG.Index]/(g.TMIN[g.TAG.Index]+237.3))
			Deltsat := (4098. * (0.6108 * math.Exp((17.27*g.TEMP[g.TAG.Index])/(g.TEMP[g.TAG.Index]+237.3)))) / math.Pow((g.TEMP[g.TAG.Index]+237.3), 2)
			if g.RAD[g.TAG.Index] > 0 {
				// ! Berechnung der Nettostrahlung aus Globalstrahlung
				RADRatio = g.RAD[g.TAG.Index] * 2 / RS0
				if RADRatio > 1 {
					RADRatio = 1
				}
				RADn = (1-Albedo)*g.RAD[g.TAG.Index]*2 - Bolz*(math.Pow((g.TMIN[g.TAG.Index]+273.16), 4)+math.Pow((g.TMAX[g.TAG.Index]+273.16), 4))/2*(1.35*RADRatio-0.35)*(0.34-0.14*math.Sqrt(Vapres))
			} else {
				var Glob float64
				if DL > 0 {
					Glob = EXT * (0.19 + 0.55*g.SUND[g.TAG.Index]/DL)
				} else {
					Glob = EXT * 0.19
				}
				RADRatio = 1
				if RS0 > 0 {
					RADRatio = Glob / RS0
				}
				if RADRatio > 1 {
					RADRatio = 1
				}

				RADn = (1-Albedo)*Glob - Bolz*(math.Pow((g.TMIN[g.TAG.Index]+273.16), 4)+math.Pow((g.TMAX[g.TAG.Index]+273.16), 4))/2*(1.35*RADRatio-0.35)*(0.34-0.14*math.Sqrt(Vapres))
			}
			g.ET0 = (0.408 * Deltsat * RADn)
			if g.ET0 < 0 {
				g.ET0 = 0
			}
			VERDU[g.TAG.Index] = g.ET0 * g.FKC * 0.1
		} else if g.ETMETH == 3 {
			// ! ----------------------- Berechnung der Referenzverdunstung Gras nach Penman-Monteith --------------
			// ! -- Notwendige inputs: Temp(Tag) (Tagesmitteltemperatur)
			// !                       Tmin(TAG) (Tagesminimumtemperatur)
			// !                       Tmax(TAG) (Tagesmaximumtemperatur)
			// !                       RH(TAG)   (relative Luftfeuchte %)
			// !                       Alti (Standorthöhe über Normal Null m)
			// !                       Wind(TAG) (mittl. Windgeschw. m/s)
			// !                       WINDHI    (Messhöhe Wind, m)
			// !                       SUND(TAG) ( Sonnenscheindauer, h), alternativ zu Strahlung
			// !                       LAT       (Breitengrad, )
			// !                       RAD(TAG)  (PAR Einstrahlung MJ/m^2/d)
			// ! -- Vordefinierte Konstante
			g.RSTOM = 100.
			Albedo := 0.23
			Bolz := 0.0000000049
			g.FKC = g.FKB
			DL, _, _, EXT, _, _, _ := CalculateDayLenght(g.TAG.Num, g.LAT)
			RS0 := (0.75 + 0.00002*g.ALTI) * EXT
			//! ---- Berechnung Atmosphärendruck
			ATMPress := 101.3 * math.Pow(((293-(0.0065*g.ALTI))/293), 5.26)
			//! ---- Berechnung Psychrometer-Konstante
			Psych := 0.000665 * ATMPress
			//! ---- Berechnung Sättigungsdampfdruck bei Tmax
			SatPmax := 0.6108 * math.Exp(((17.27 * g.TMAX[g.TAG.Index]) / (237.3 + g.TMAX[g.TAG.Index])))
			//! ---- Berechnung Sättigungsdampfdruck bei Tmin
			SatPmin := 0.6108 * math.Exp(((17.27 * g.TMIN[g.TAG.Index]) / (237.3 + g.TMIN[g.TAG.Index])))
			//! ---- Berechnung mittlerer Sättigungsdampfdruck
			SatP := (SatPmin + SatPmax) / 2
			//! ---- Berechnung aktueller Dampfdruck
			Vapres := SatP * g.RH[g.TAG.Index] / 100
			//! ---- Berechnung Sättigungsdefizit der Luft
			l.SATDEF = SatP * (1 - g.RH[g.TAG.Index]/100)
			//! Deltsat = Steigung der Sättigungsdampfdruck-Temperatur Beziehung
			Deltsat := (4098. * (0.6108 * math.Exp((17.27*g.TEMP[g.TAG.Index])/(g.TEMP[g.TAG.Index]+237.3)))) / math.Pow((g.TEMP[g.TAG.Index]+237.3), 2)
			//! ----- wenn Windgeschwindigkeit Messung nicht in 2m Höhe, dann Umrechnung auf 2m
			if g.WINDHI != 2 {
				g.WIND[g.TAG.Index] = g.WIND[g.TAG.Index] * (4.87 / (math.Log(67.8*g.WINDHI - 5.42)))
			}
			//    ! ----- Berechnung des aerodynamischen Widerstands ra (Raero)
			if g.WIND[g.TAG.Index] < 0.5 { // fixed lower bound for wind speed to 0.5
				g.WIND[g.TAG.Index] = 0.5
			}
			//    ! ----- Berechnung des Oberflächenwiderstands rs mit Stomatawiderstand = 100 s/m (Rsurf)
			Rsurf := g.RSTOM / 1.44
			if g.RAD[g.TAG.Index] > 0 {
				// ! Berechnung der Nettostrahlung aus Globalstrahlung und Sonnenscheindauer
				if RS0 > 0 {
					RADRatio = g.RAD[g.TAG.Index] * 2 / RS0
					if RADRatio > 1 {
						RADRatio = 1
					}
				} else {
					RADRatio = 1
				}
				RADn = (1-Albedo)*g.RAD[g.TAG.Index]*2 - Bolz*(math.Pow((g.TMIN[g.TAG.Index]+273.16), 4)+math.Pow((g.TMAX[g.TAG.Index]+273.16), 4))/2*(1.35*RADRatio-0.35)*(0.34-0.14*math.Sqrt(Vapres))
			} else {
				var Glob float64
				if DL > 0 {
					Glob = EXT * (0.19 + 0.55*g.SUND[g.TAG.Index]/DL)
				} else {
					Glob = EXT * 0.19
				}
				RADRatio = 1
				if RS0 > 0 {
					RADRatio = Glob / RS0
				}
				if RADRatio > 1 {
					RADRatio = 1
				}
				RADn = (1-Albedo)*Glob - Bolz*(math.Pow((g.TMIN[g.TAG.Index]+273.16), 4)+math.Pow((g.TMAX[g.TAG.Index]+273.16), 4))/2*(1.35*RADRatio-0.35)*(0.34-0.14*math.Sqrt(Vapres))
			}
			//    ! ------------------------ Berechnung der Referenzevapotranspiration --------------------
			g.ET0 = ((0.408 * Deltsat * RADn) + (Psych * (900 / (g.TEMP[g.TAG.Index] + 273)) * g.WIND[g.TAG.Index] * l.SATDEF)) / (Deltsat + Psych*(1+(Rsurf/208)*g.WIND[g.TAG.Index]))
			if g.ET0 < 0 {
				g.ET0 = 0
			}
			VERDU[g.TAG.Index] = g.ET0 * g.FKC * 0.1
		}
		if VERDU[g.TAG.Index] > 0.6 {
			VERDU[g.TAG.Index] = 0.6
		}

		EVMAX = VERDU[g.TAG.Index]
		TRAMAX = 0
		g.WURZ = 0
	}

	if EVMAX > .65 {
		EVMAX = .65
	}
	g.VERDUNST = g.VERDUNST + VERDU[g.TAG.Index]*g.DT.Num

	g.ETC0 = g.ETC0 + VERDU[g.TAG.Index]*g.DT.Num

	// ! Berechnung der aktuellen Evaporation abh. von Oberbodenfeuchte
	l.EVA[g.TAG.Index] = EVMAX*REDEV - g.REGEN[g.TAG.Index]
	g.ETA = EVMAX * REDEV
	// ! --------- Verteilung der Evaporation über die Tiefe (Mimic procedure) --------------
	// ! Input: PROP       = Verteilungsfaktor Exponentialfunktion (abh. von Hauptbodenart)
	if l.EVA[g.TAG.Index] > 0 {
		SUMVAR := 0.
		// ! ++++++++++++ Neu: Bei Überstau Entnahme aus Ueberstauwasser bzw. 1. Schicht ++++++++++++++
		// IF Storage > 0 then
		if g.STORAGE > 0 {
			//    IF storage > EVa(tag)*dt then
			if g.STORAGE > l.EVA[g.TAG.Index]*g.DT.Num {
				// 	  LET storage = Storage - Eva(Tag) *dt
				g.STORAGE = g.STORAGE - l.EVA[g.TAG.Index]*g.DT.Num
				// 	  LET EVA(Tag) = 0
				// l.EVA[g.TAG.Index] = 0
				l.EVA[g.TAG.Index] = -g.STORAGE
				//    ELSE
			} else {
				// 	  LET EVA(TAG) = eva(tag)*dt-storage
				l.EVA[g.TAG.Index] = l.EVA[g.TAG.Index]*g.DT.Num - g.STORAGE
				// 	  LET storage = 0
				g.STORAGE = 0
				// 	  !LET fluss0 = -EV(1)
				//    END IF
			}
			//    FOR I = 2 to n
			for i := 1; i < g.N; i++ {
				// 	   LET ev(i) = 0
				l.EV[i] = 0
				//    NEXT i
			}
			// ELSE
		} else {
			for i := 0; i < g.N; i++ {
				if g.WG[0][i]-g.WMIN[i]/3 > 0 {
					VAR[i] = (g.WG[0][i] - g.WMIN[i]/3) * math.Exp(-g.PROP*.1*(float64((i+1)*10)-g.DZ.Num/2))
				} else {
					VAR[i] = 0
				}
				SUMVAR = SUMVAR + VAR[i]
			}
			for i := 0; i < g.N; i++ {
				if SUMVAR > 0 {
					l.EV[i] = l.EVA[g.TAG.Index] * VAR[i] / SUMVAR * g.DT.Num
				} else {
					l.EV[i] = 0
				}
			}
		}
	} else {
		g.STORAGE = g.STORAGE - l.EVA[g.TAG.Index]*g.DT.Num
		l.EVA[g.TAG.Index] = -g.STORAGE

		for i := 0; i < g.N; i++ {
			l.EV[i] = 0
		}
	}
	// ! --------------------------------------------------------------------------------------
	// Wasserfluss durch die Bodenoberfläche
	g.FLUSS0 = -l.EVA[g.TAG.Index]
	WEFF := 0.
	l.GWAUF = 0.
	var WEFFREST float64
	// ! Verteilung der Transpiration über die Tiefe abh. von Wasserverfügbarkeit und Durchwurzelungsdichte
	// ! Inputs:
	// ! WUDICH(I)        = Wurzellängendichte in cm/cm^3 (aus Pflanzenmodell)
	// ! NFK(I)           = nutzbare FK, siehe oben
	// ! Outputs:
	// ! TRRED(I)         = Reduktion der Wasseraufnahme in Tiefe I
	// ! WUEFF(I)         = Wurzelaktivitätsfaktor in Abh. Wasserverfügbarkeit in Tiefe I
	// IF ZEIT > SAAT(AKF) AND ZEIT < ERNTE(AKF) AND Intwick > 1 THEN
	if zeit > g.SAAT[g.AKF.Index] && g.INTWICK.Num > 1 &&
		((g.ERNTE[g.AKF.Index] > 0 && zeit < g.ERNTE[g.AKF.Index]) || (g.ERNTE[g.AKF.Index] == 0 && zeit < g.ERNTE2[g.AKF.Index])) {
		for i := 0; i < g.WURZ; i++ {
			if l.NFK[i] < .15 {
				TRRED[i] = l.NFK[i] * 3
			} else if l.NFK[i] < .3 {
				TRRED[i] = .45 + (.25 * (l.NFK[i] - .15) / .15)
			} else if l.NFK[i] < .5 {
				TRRED[i] = .7 + (.275 * (l.NFK[i] - .3) / .2)
			} else if l.NFK[i] < .75 {
				TRRED[i] = .975 + (.025 * (l.NFK[i] - .5) / .25)
			} else {
				TRRED[i] = 1
			}
			if TRRED[i] < 0 {
				TRRED[i] = 0
			}
			if l.NFK[i] < .15 {
				WUEFF[i] = .15 + .45*l.NFK[i]/.15
			} else if l.NFK[i] < .3 {
				WUEFF[i] = .6 + (.2 * (l.NFK[i] - .15) / .15)
			} else if l.NFK[i] < .5 {
				WUEFF[i] = .8 + (.2 * (l.NFK[i] - .3) / .2)
			} else {
				WUEFF[i] = 1
			}
			if WUEFF[i] < 0 {
				WUEFF[i] = 0
			}
			if float64(i+1) > g.GRW {
				WUEFF[i] = 0
			}
			WEFF = WEFF + WUEFF[i]*g.WUDICH[i]
			WEFFREST = WEFF
		}
		// ! Reduktion bei Luftmangel (fuer obere 30 cm)
		// ! Inputs:
		// ! LUKRIT(INTWICK)       = kritisches Luftporenvolumen (cm3/cm3) in Entwicklungsstatdium INTWICK
		// ! INTWICK               = aktuelles Entwicklungsstadium der Pflanze
		// ! PORGES(I)             = Gesamtporenvolumen Schicht I (cm3/cm3)
		// ! Variable
		// ! LURMAX                = Ausmaß Luftmangel (s.u.)
		// ! LUMDAY                = kumulative Dauer des Luftmangels (Tage), maximum 4
		// ! LURED                 = Reduktionsfaktor fuer Transpiration
		LUPOR := (g.PORGES[0] + g.PORGES[1] + g.PORGES[2] - g.WG[0][0] - g.WG[0][1] - g.WG[0][2]) / 3
		if LUPOR < g.LUKRIT[g.INTWICK.Index] {
			g.LUMDAY = g.LUMDAY + g.DT.Index
			if g.LUMDAY > 4 {
				g.LUMDAY = 4
			}
			if LUPOR < 0 {
				LUPOR = 0.
			}
			LURMAX := LUPOR / g.LUKRIT[g.INTWICK.Index]
			g.LURED = 1 - float64(g.LUMDAY)/4*(1-LURMAX)
		} else {
			g.LUMDAY = 0
			g.LURED = 1
		}
		if g.LURED > 1 {
			g.LURED = 1
		}
		// ! Verteilung der pot. Tranpiration und Einschränkung der Transpiration bei Luftmangel
		// ! WURZ = Wurzeltiefe
		// ! GRW  = Grundwasserstand (Wasseraufnahme nur bis zur ersten GW Schicht)
		// ! TP(I)= Wasseraufnahme in Schicht I
		for i := 0; i < g.N; i++ {
			if float64(i+1) > math.Min(float64(g.WURZ), g.GRW) {
				g.TP[i] = 0
			} else {
				if WUEFF[i]*g.WUDICH[i] > 0 {
					g.TP[i] = TRAMAX * WUEFF[i] * g.WUDICH[i] / WEFF * g.LURED
				} else {
					g.TP[i] = 0
				}
			}
		}
		// ! Defizit bei Wasseraufnahme TP(I) kann nach unten verteilt werden bis WURZ
		min := math.Min(float64(g.WURZ), g.GRW)
		for i := 1; i <= int(min); i++ {
			index := i - 1
			TREST := 0.
			WEFFREST = WEFFREST - WUEFF[index]*g.WUDICH[index]
			TDEFT := 0.
			if g.TP[index]/g.DZ.Num > (g.WG[0][index] - g.WMIN[index]) {
				TDEFT = (g.TP[index]/g.DZ.Num - (g.WG[0][index] - g.WMIN[index])) * g.DZ.Num
				if TDEFT < 0 {
					TDEFT = 0
				}
				if TDEFT > g.TP[index]/g.DZ.Num {
					TDEFT = g.TP[index] / g.DZ.Num
				}
			} else {
				TDEFT = 0
			}
			TDRED := g.TP[index] * (1 - TRRED[index])
			TREST = math.Max(TDRED, TDEFT)
			if TREST > g.TP[index] {
				TREST = g.TP[index]
			}
			if TREST > 0 {
				if float64(i) < min {
					for I2 := i + 1; I2 <= int(min); I2++ {
						I2index := I2 - 1
						if WEFFREST > 0 {
							g.TP[I2index] = g.TP[I2index] + TREST*WUEFF[I2index]*g.WUDICH[I2index]/WEFFREST
						}
					}
				}
			}
			g.TP[index] = g.TP[index] - TREST
			if g.TP[index] < 0 {
				g.TP[index] = 0
			}
			TPAKT = TPAKT + g.TP[index]
			if float64(i) == g.GRW {
				l.GWAUF = g.TP[index]
			}
		}

		if ETCP > 0 {
			g.ETREL = (TPAKT + g.ETA) / ETCP
		} else {
			g.ETREL = 1
		}
		if g.ETREL > 1 {
			g.ETREL = 1
		}

		if TRAMAX > 0 {
			g.TRREL = TPAKT / TRAMAX
		}
	} else {
		for i := 0; i < g.N; i++ {
			g.TP[i] = 0
		}
		g.TRREL = 1
	}
}

// -------- Unterprogramm Stomatawiderstand in Abh. C-Assimilation, Sättigungsdefizit und CO2 -------------------------

func stomat(l *WaterSharedVars, zeit int, g *GlobalVarsMain) {

	// ! Inputs:
	// ! AMAX            = maximale C-Assimilation bei Lichtsättigung
	// ! TEMP(TAG)       = Tagesmitteltemperatur (°C)
	// ! RAD(TAG)        = PAR (MJ/m^2/d)
	// ! SUND(TAG)       = Sonnenscheindauer (h)
	// ! LAT             = Breitengrad
	// ! CO2Konz         = CO2-Konzentration (ppm)
	// ! SATDEF          = Sättigungsdefizit der Luft (kPa?)

	DL, DLE, _, _, RDN, DRC, DEC := CalculateDayLenght(g.TAG.Num, g.LAT)
	if DLE <= 0 {
		return
	}
	DRO := .2 * DRC
	EFF0 := .5
	// ! ++++++++++++++  Auswahl mehrerer Methoden zum CO2 Effect +++++++++++++++
	var EFF float64
	var COcomp float64
	if g.CO2METH == 1 {
		COcomp := 17.5 * math.Pow(2, ((g.TEMP[g.TAG.Index]-10)/10))
		EFF = (g.CO2KONZ - COcomp) / (g.CO2KONZ + 2*COcomp) * EFF0
	} else {
		EFF = EFF0
	}
	MAXAMAXG := 30.
	var amax float64
	if g.TEMP[g.TAG.Index] < g.MINTMP {
		amax = 0
	} else if g.TEMP[g.TAG.Index] < 10 {
		amax = MAXAMAXG * g.TEMP[g.TAG.Index] / 10 * .4
	} else if g.TEMP[g.TAG.Index] < 15 {
		amax = MAXAMAXG * (.4 + (g.TEMP[g.TAG.Index]-10)/5*.5)
	} else if g.TEMP[g.TAG.Index] < 25 {
		amax = MAXAMAXG * (.9 + (g.TEMP[g.TAG.Index]-15)/10*.1)
	} else if g.TEMP[g.TAG.Index] < 35 {
		amax = MAXAMAXG * (1 - (g.TEMP[g.TAG.Index]-25)/10)
	} else {
		amax = 0
	}
	if g.CO2METH == 1 {
		amax = amax * (g.CO2KONZ - COcomp) / (350 - COcomp)
	} else if g.CO2METH == 2 {
		var KCo1 float64
		var coco float64
		if g.RAD[g.TAG.Index] > 0 {
			KCo1 = 220 + 0.158*g.RAD[g.TAG.Index]*20
			coco = 80 - 0.0036*g.RAD[g.TAG.Index]*20
		} else {
			SC := 1367 * (1 + 0.033*math.Cos(2*math.Pi*g.TAG.Num/365))
			EXT := SC * RDN / 10000
			var Glob float64
			if DL > 0 {
				Glob = EXT * (0.19 + 0.55*g.SUND[g.TAG.Index]/DL)
			} else {
				Glob = EXT * 0.19
			}
			KCo1 = 220 + 0.158*Glob
			coco = 80 - 0.0036*Glob
		}
		KCO2 := ((g.CO2KONZ - coco) / (KCo1 + g.CO2KONZ - coco)) / ((350 - coco) / (KCo1 + 350 - coco))
		amax = amax * KCO2
	}
	// ! ----------------------- Strahlungsinterception nach Penning de Vries ---------------------
	if amax < .1 {
		amax = .1
	}
	if DLE == 0 && DL > 0 {
		DLE = 0.1
	}
	REFLC := .08
	EFFE := (1. - REFLC) * EFF
	SSLAE := math.Sin((90. + DEC - g.LAT) * math.Pi / 180.)
	X := math.Log(1. + .45*DRC/(DLE*3600.)*EFFE/(SSLAE*amax))
	PHCH1 := SSLAE * amax * DLE * X / (1. + X)
	// ! Aenderung nach P.d. Vries am 25.5.93
	Y := math.Log(1. + .55*DRC/(DLE*3600.)*EFFE/((5-SSLAE)*amax))
	PHCH2 := (5. - SSLAE) * amax * DLE * Y / (1. + Y)
	PHCH := 0.95*(PHCH1+PHCH2) + 20.5
	// 1.44 = LAI kurzgeschnittenes Gras
	PHC3 := PHCH * (1. - math.Exp(-.8*1.44))

	//  1.44 = LAI kurzgeschnittenes Gras
	PHC4 := DL * 1.44 * amax
	var MIPHC float64
	var MAPHC float64
	if PHC3 < PHC4 {
		MIPHC = PHC3
		MAPHC = PHC4
	} else {
		MIPHC = PHC4
		MAPHC = PHC3
	}
	PHCL := MIPHC * (1. - math.Exp(-MAPHC/MIPHC))
	Z := DRO / (DLE * 3600.) * EFFE / (5. * amax)
	PHOH1 := 5. * amax * DLE * Z / (1. + Z)
	PHOH := 0.9935*PHOH1 + 1.1
	// ! Aenderung nach P.d. Vries am 25.5.93
	PHO3 := PHOH * (1. - math.Exp(-.8*1.44))
	var MIPHO float64
	var MAPHO float64
	if PHO3 < PHC4 {
		MIPHO = PHO3
		MAPHO = PHC4
	} else {
		MIPHO = PHC4
		MAPHO = PHO3
	}
	PHOL := MIPHO * (1. - math.Exp(-MAPHO/MIPHO))
	DGAC := PHCL
	DGAO := PHOL
	var DTGA float64
	// !----------- BERUECKSICHTIGUNG DER SONNENSCHEINDAUER -------
	if g.RAD[g.TAG.Index] == 0 {
		if g.SUND[g.TAG.Index] > DLE {
			g.SUND[g.TAG.Index] = DLE
		}
		DTGA = g.SUND[g.TAG.Index]/DLE*DGAC + (1.-g.SUND[g.TAG.Index]/DLE)*DGAO
	} else {
		KOREK := 1.
		g.RADSUM = g.RADSUM + g.RAD[g.TAG.Index]*g.DT.Num*KOREK
		//! Fraktion bedeckter Tag, DRC = Fraktion klarer tag
		FOV := (DRC - 1000000*g.RAD[g.TAG.Index]*KOREK) / (.8 * DRC)
		if FOV > 1 {
			FOV = 1
		}
		if FOV < 0 {
			FOV = 0
		}
		DTGA = FOV*DGAO + (1-FOV)*DGAC
	}
	// ------- PHOTOSYNTHESERATE IN KG GLUCOSE/HA BLATT/TAG------
	Agross := DTGA / (10 * 3600 * 24 * 44) * 22414
	g.RSTOM = 1 / (g.ALPH * Agross / (g.CO2KONZ * (1 + l.SATDEF/g.SATBETA)))
}

// Water model
func Water(wdt float64, subd int, zeit int, g *GlobalVarsMain, l *WaterSharedVars) {

	// !**************************************************************
	// !***                WASSERMODELL NACH BURNS                 ***
	// !***    Versickerung bei Überschreiten der Feldkapazität    ***
	// !***  Aufstieg bei Evaporation in Anlehnung an Groot (1987) ***
	// !**************************************************************
	// !*** BENUTZTE VARIABLE:                                     ***
	// !***     SCHICHTDICKE (cm):                         dz      ***
	// !***     FELDKAPAZITAET (cm^3/cm^3):                W(z)    ***
	// !***     ANZAHL SCHICHTEN:                          N       ***
	// !***     WASSERGEHALT DER SCHICHTEN (cm^3/cm^3):    WG(,)   ***
	// !***     GESAMTWASSER PRO SCHICHT (cm):             WASSER  ***
	// !***     REGEN(+) BZW. EVAPORATION (-) (cm):        FLUSS0  ***
	// !***     Draintiefe bei Drainung (dm)               DRAIDEP ***
	// !***     Drainfaktor (Fraktion) wenn > FK           DRAIFAK ***
	// !***     maximaler Fluss bei FK                     qm(z)   ***
	// !***     Profilminimum des maximalen Flusses        qmax    ***
	// !**************************************************************
	var WATER [2][21]float64
	if subd == 1 {
		g.EvapoLoss = 0
		g.SickerLoss = 0
		g.InfilDaily = 0
		g.WaterDiff = 0
		gwtab := 0
		for i := 0; i < g.N; i++ {
			if g.TP[i] > (g.WG[0][i]-g.WMIN[i])*g.DZ.Num {
				if g.WG[0][i] < g.WMIN[i] {
					g.TP[i] = 0
				} else {
					g.TP[i] = (g.WG[0][i] - g.WMIN[i]) * g.DZ.Num
				}
			}
			WATER[0][i] = g.WG[0][i]*g.DZ.Num - g.TP[i]*wdt
			g.WaterDiff = g.WaterDiff - g.TP[i]*wdt
			//! -------- Definition des Stauwasserhorizonts ----------
			if WATER[0][i]/g.DZ.Num >= g.PORGES[i]-g.TP[i]*wdt && gwtab == 0 {
				gwtab = i
			}
		}
		// meins
		if g.STORAGE > 0 {
			qmax := g.PORGES[0]*g.DZ.Num - WATER[0][0] + g.QM[0]
			l.Maxinfil = math.Min(g.STORAGE, qmax)
		} else {
			l.Maxinfil = 0
		}
	} else {
		for i := 0; i < g.N; i++ {
			g.WG[0][i] = g.WG[1][i]
			WATER[0][i] = g.WG[0][i]*g.DZ.Num - g.TP[i]*wdt
			g.WaterDiff = g.WaterDiff - g.TP[i]*wdt
		}
	}
	g.QDRAIN = 0
	if g.FLUSS0 > 0 {

		// ! Bei Wassersättigung und Überstau Infiltration bis kleinste Profilleitfähigkeit
		//    !IF storage > 0 or gwtab = 1 then
		//        LET qmax = PORGES(1)*dz - WATER(0,1) + qm(1)
		//his
		//qmax := g.PORGES[0]*g.DZ.Num - WATER[0][0] + g.QM[0]
		//       LET storage = storage + fluss0 * wdt
		// g.STORAGE = g.FLUSS0 * wdt
		//       Let maxinfil = Min(storage,qmax*wdt)
		//his
		//maxinfil := math.Min(g.STORAGE, qmax*wdt)
		// maxinfil := math.Min(g.STORAGE*wdt, qmax)
		// if wdt < 1 {
		// 	fmt.Println(g.PORGES[0]*g.DZ.Num, WATER[0][0], g.QM[0], wdt, g.STORAGE, l.Maxinfil)
		// }
		//       Let storage = storage - maxinfil

		//his
		//g.STORAGE = g.STORAGE - maxinfil
		//a := maxinfil

		// g.STORAGE = g.STORAGE - l.Maxinfil
		// a := l.Maxinfil

		// meins
		a := l.Maxinfil * wdt
		if l.Maxinfil*wdt >= g.STORAGE {
			a = g.STORAGE
			g.STORAGE = 0
		} else {
			g.STORAGE = g.STORAGE - l.Maxinfil*wdt
		}

		//a := g.FLUSS0 * wdt
		g.Q1[0] = a
		g.InfilDaily += a
		//------------------------ Infiltration------------------------
		for k1 := 1; k1 <= g.N; k1++ {
			k1INdex := k1 - 1
			b := a + WATER[0][k1INdex]
			a = b - g.W[k1INdex]*g.DZ.Num
			if a < 0 {
				g.WaterDiff = g.WaterDiff + b - WATER[0][k1INdex]
				WATER[1][k1INdex] = b
				g.Q1[k1] = 0
				a = 0
				// for k2 := k1 + 1; k2 <= g.N; k2++ {
				// 	k2Index := k2 - 1
				// 	WATER[1][k2Index] = WATER[0][k2Index]
				// 	g.Q1[k2] = 0
				// }
				// break
			} else {

				var qma float64
				//IF k1 < n then
				if k1 < g.N {
					// LET qma = PORGES(k1+1)*dz - WATER(0,k1+1) + qm(k1+1)
					qma = g.PORGES[k1INdex+1]*g.DZ.Num - WATER[0][k1INdex+1] + g.QM[k1INdex+1]*wdt
				} else {
					// LET qma = qm(k1)
					qma = g.QM[k1INdex] * wdt
					// end if
				}
				if k1 == g.DRAIDEP {
					// g.Q1[k1] = (1 - g.DRAIFAK) * a
					// g.QDRAIN = g.DRAIFAK * a
					// a = g.Q1[k1]

					// LET Q1(k1) = MIN((1-draifak) * a,qm(k1))
					g.Q1[k1] = math.Min((1-g.DRAIFAK)*a, qma)
					// LET qdrain = a-q1(k1)  !draifak * a
					g.QDRAIN = a - g.Q1[k1]
					// LET a = q1(k1)
					a = g.Q1[k1]
					//LET WATER(1,k1) = W(K1) * dz
					WATER[1][k1INdex] = g.W[k1INdex] * g.DZ.Num

				} else {
					// g.Q1[k1] = a

					// LET Q1(k1) = MIN(a,qm(k1))
					g.Q1[k1] = math.Min(a, qma)
					// LET a = q1(k1)
					a = g.Q1[k1]
					// LET WATER(1,k1) = b-a  ! W(K1) * dz
					WATER[1][k1INdex] = b - a
					g.WaterDiff = g.WaterDiff + (b - WATER[0][k1INdex] - a)
				}
				// WATER[1][k1INdex] = g.W[k1INdex] * g.DZ.Num
				// g.Q1[k1] = a
			}
		}
		// if wdt < 1 {
		// 	fmt.Println(g.PORGES[0]*g.DZ.Num, WATER[0][0], g.STORAGE)
		// }
		//--------------------------Evaporation------------------------
	} else if g.FLUSS0 < 0 {
		a := math.Abs(g.FLUSS0) * wdt
		a1 := a
		g.Q1[0] = 0
		for k1 := 0; k1 < g.N; k1++ {
			k1IdxQ := k1 + 1
			l.LIMIT[k1] = WATER[0][k1] - l.EV[k1]*wdt
			if l.LIMIT[k1] < (g.WMIN[k1]/3)*g.DZ.Num {
				l.EV[k1+1] = l.EV[k1+1] + (l.EV[k1] - WATER[0][k1] + g.WMIN[k1]/3*g.DZ.Num)
				l.EV[k1] = WATER[0][k1] - g.WMIN[k1]/3*g.DZ.Num
				l.LIMIT[k1] = g.WMIN[k1] / 3 * g.DZ.Num
			}
			vcap := WATER[0][k1] - l.LIMIT[k1]
			var wlost float64
			if vcap > a1 {
				wlost = a1
				//WATER[1][k1] = WATER[0][k1] - wlost
				g.Q1[k1IdxQ] = 0
				//LET a1 = 0
				a1 = 0
				// for k2 := k1 + 1; k2 < g.N; k2++ {
				// 	WATER[1][k2] = WATER[0][k2]
				// 	g.Q1[k2+1] = 0
				// }
				//break
			} else {
				wlost = vcap
				a1 = a1 - vcap
				g.Q1[k1IdxQ] = -a1
			}
			WATER[1][k1] = WATER[0][k1] - wlost
			g.EvapoLoss += wlost / 10
			g.WaterDiff = g.WaterDiff - wlost
			//! ++++++++downward water flow of excessive water over FC from previous time steps+++++++++++++++++++++++++
			//    IF WATER(1,k1) > W(k1)*dz then
			if WATER[1][k1] > g.W[k1]*g.DZ.Num {
				//LET a = MIN(WATER(1,k1) - W(K1) * dz,qm(k1))
				a = math.Min(WATER[1][k1]-g.W[k1]*g.DZ.Num, g.QM[k1]*wdt)
				//LET WATER(1,k1) = WATER(1,k1)-a
				WATER[1][k1] = WATER[1][k1] - a
				g.WaterDiff = g.WaterDiff - a
				//LET Q1(k1) = q1(k1) + a
				g.Q1[k1IdxQ] = g.Q1[k1IdxQ] + a
				//IF q1(k1) > 0 then
				if g.Q1[k1IdxQ] > 0 {
					//LET a1 = 0
					a1 = 0
					//LET WATER(0,k1+1) = Water (0,k1+1) + Q1(k1)
					WATER[0][k1+1] = WATER[0][k1+1] + g.Q1[k1IdxQ]
					//ELSE
				} else {
					//LET a1 = -q1(k1)
					a1 = -g.Q1[k1IdxQ]
					//END IF
				}
				//END IF
			}
			//! +++++++++++++++++++++++++++++++++++
		}
	} else {
		for i := 0; i < g.N; i++ {
			WATER[1][i] = WATER[0][i]
			g.Q1[i+1] = 0
		}
	}

	// for i := 0; i < g.N; i++ {
	// 	if WATER[1][i]/g.DZ.Num > g.W[i] {
	// 		sink := WATER[1][i] - g.W[i]*g.DZ.Num
	// 		WATER[1][i] = g.W[i] * g.DZ.Num
	// 		WATER[1][i+1] = WATER[1][i+1] + sink
	// 		g.Q1[i+1] = g.Q1[i+1] + sink
	// 	}
	// }

	// ! -------------------  Berechnung des kapillaren Aufstiegs -----------------
	// ! GRW            = Grundwasserstand (dm unter Flur)
	// ! CAPLAY         = Tiefste Schicht mit nFK < 0.7
	// ! GWDIST         = Abstand CAPLAY zu GRW (dm)
	// ! Inputs:
	// ! CAPS           = täglicher kapillarer Aufstieg aus GW (aus KA 5 Tabelle in Abh. von Bodenart und Abstand zum GW (cm/d)
	// ! Q1(I)          = Fluss durch untere Schichtgrenze von Komartment I (cm/d)
	caplay := 0
	capdep := 1
	for i := g.N; i >= capdep; i-- {
		index := i - 1
		if caplay == 0 {
			if l.NFK[index] < 0.7 {
				caplay = i
			}
		}
	}
	if caplay > 0 {
		caplayIndex := caplay - 1
		GWDIST := g.GRW + 1 - float64(caplay)
		if GWDIST < 21 {
			if l.NFK[caplayIndex] < .7 {
				if GWDIST < 0 {
					GWDIST = 0
				}
				if GWDIST > .9 {
					GWDISTindex := int(math.Round(math.Max(GWDIST, 1))) - 1
					WATER[1][caplayIndex] = WATER[1][caplayIndex] + g.CAPS[GWDISTindex]*g.DZ.Num*wdt
					g.WaterDiff = g.WaterDiff + g.CAPS[GWDISTindex]*g.DZ.Num*wdt
					for i := caplay; i <= g.N; i++ {
						g.Q1[i] = g.Q1[i] - g.CAPS[GWDISTindex]*g.DZ.Num*wdt
					}
				}
			}
		}
	}
	// !---------- Berechnung der neuen Wasserkonzentrationen ----------
	for i := 1; i <= g.N; i++ {
		index := i - 1
		g.WG[1][index] = WATER[1][index] / g.DZ.Num
		g.PFTRANS = g.PFTRANS + g.TP[index]*wdt
		g.TRAY = g.TRAY + g.TP[index]*wdt
		if zeit > g.SAAT[g.AKF.Index] {
			g.TRAG = g.TRAG + g.TP[index]*wdt
			g.ETAG = g.ETAG + g.TP[index]*wdt
		}

		if i < 4 {
			g.TP3 = g.TP3 + g.TP[index]*wdt
		} else if i < 7 {
			g.TP6 = g.TP6 + g.TP[index]*wdt
		} else if i < 10 {
			g.TP9 = g.TP9 + g.TP[index]*wdt
		}
	}
	// ! Summierung für Ausgabe
	g.PFTRANS = g.PFTRANS + g.ETA*wdt

	if zeit > g.SAAT[g.AKF.Index] {
		g.ETAG = g.ETAG + g.ETA*wdt
	}

	g.WG[1][g.N] = g.WG[1][g.N-1]
	g.DRAISUM = g.DRAISUM + g.QDRAIN*10
	// 	IF subd = 1 then
	if subd == 1 {
		// 	LET draitag = 0
		g.DRAITAG = 0
		//  END IF
	}
	//  IF subd<= dt/wdt then
	if float64(subd) <= g.DT.Num/wdt {
		// 	LET DRAITAG = Draitag+Qdrain*10
		g.DRAITAG = g.DRAITAG + g.QDRAIN*10
		//  END IF
	}
	if g.Q1[g.OUTN] > 0 {
		g.SICKER = g.SICKER + g.Q1[g.OUTN]*10
		g.SickerLoss = g.SickerLoss + g.Q1[g.OUTN]/10
	} else {
		g.CAPSUM = g.CAPSUM + g.Q1[g.OUTN]*10
	}
	g.CAPSUM = g.CAPSUM - l.GWAUF*10*wdt

	if zeit > g.SAAT[g.AKF.Index] {
		g.PERG = g.PERG + g.Q1[g.OUTN]*10 - l.GWAUF*10*wdt
	}

	if g.FLUSS0 > 0 {
		g.INFILT = g.INFILT + g.FLUSS0*wdt
	}
}
