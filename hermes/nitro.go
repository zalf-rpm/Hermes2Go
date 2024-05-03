package hermes

import (
	"fmt"
	"math"
)

// NitroBBBSharedVars shared variables for this module
type NitroBBBSharedVars struct {
	DDAT1   string
	DDAT2   string
	DDAT3   string
	DODAT   string
	DUNGART string
	DMENG1  float64
	DMENG2  float64
	DMENG3  float64
	DOMENG1 float64
	NRESID  float64
	NAGB    float64
	NUPTAKE float64
}

// NitroSharedVars shared variables for this module
type NitroSharedVars struct {
	DUMS    [4]float64
	D       [21]float64
	V       [21]float64
	KONV    [21]float64
	DISP    [21]float64
	DB      [21]float64
	DM      [21]float64
	DNH4UMS [4]float64
}

// Nitro ...
func Nitro(wdt float64, subd int, zeit int, g *GlobalVarsMain, l *NitroSharedVars, ln *NitroBBBSharedVars, hPath *HFilePath, output *CropOutputVars) (finishedCycle bool, runErr error) {
	finishedCycle = false
	runErr = nil
	writeFertilizerEvent := func(fertName string, ndir, nh4n float64) error {

		fertList := map[string]interface{}{
			"Fertilizer": fertName,
			"Ndirect":    ndir,
			"NH4":        nh4n,
		}
		if ndir == 0 {
			delete(fertList, "Ndirect")
		} else if nh4n == 0 {
			delete(fertList, "NH4")
		}
		err := g.managementConfig.WriteManagementEvent(NewManagementEvent(Fertilization, zeit, fertList, g))
		return err
	}

	if !g.AUTOFERT {
		//! +++++++++++++++++++++++++++++++++++++ Option real fertilization +++++++++++++++++++++++++++++++++++++++++++++++
		if zeit == g.ZTDG[g.NDG.Index]+1 && subd == 1 {
			g.NFOS[0] = g.NFOS[0] + g.NSAS[g.NDG.Index]
			g.NAOS[0] = g.NAOS[0] + g.NLAS[g.NDG.Index]
			g.DSUMM = g.DSUMM + g.NDIR[g.NDG.Index] //! Summe miner. Duengung
			g.NFERTSIM = g.NFERTSIM + g.NDIR[g.NDG.Index]
			g.NH4Sum = g.NH4Sum + g.NH4N[g.NDG.Index] // Summe min. Ammoniakalische Düngung

			if runErr = writeFertilizerEvent(g.DGART[g.NDG.Index], g.NDIR[g.NDG.Index], g.NH4N[g.NDG.Index]); runErr != nil {
				return finishedCycle, runErr
			}
			g.NDG.Inc()
		}
		//! ---------------------------------------------------------------------------------------------------------------
	} else if g.AUTOFERT {
		//! still to be defined
		if subd == 1 {
			if g.AKF.Num > 1 {
				if g.ODU[g.AKF.Index-1] == 1 && g.ORGTIME[g.AKF.Index-1] == "H" {
					if zeit == g.ZTDG[g.AKF.Index-1] {
						g.NFOS[0] = g.NFOS[0] + g.NSAS[g.AKF.Index-1]
						g.NAOS[0] = g.NAOS[0] + g.NLAS[g.AKF.Index-1]
						ln.DODAT = g.Kalender(zeit)
						ln.DOMENG1 = g.NSAS[g.AKF.Index-1] + g.NLAS[g.AKF.Index-1] + g.NDIR[g.AKF.Index-1]
						ln.DUNGART = g.DGART[g.AKF.Index-1]
						g.DSUMM = g.DSUMM + g.NDIR[g.AKF.Index-1] // ! Summe miner. Duengung

						if runErr = writeFertilizerEvent(g.DGART[g.AKF.Index-1], g.NDIR[g.AKF.Index-1], 0); runErr != nil {
							return finishedCycle, runErr
						}
					}
				}
			}
			if g.SAAT[g.AKF.Index] > 0 {
				if zeit >= g.SAAT[g.AKF.Index] {
					if g.ODU[g.AKF.Index] == 1 && g.ORGTIME[g.AKF.Index-1] == "S" {
						if zeit == g.SAAT[g.AKF.Index] {
							g.ZTDG[g.AKF.Index] = zeit + g.ORGDOY[g.AKF.Index]
						}
						if zeit == g.ZTDG[g.AKF.Index] {
							g.NFOS[0] = g.NFOS[0] + g.NSAS[g.AKF.Index]
							g.NAOS[0] = g.NAOS[0] + g.NLAS[g.AKF.Index]
							ln.DODAT = g.Kalender(zeit)
							ln.DOMENG1 = g.NSAS[g.AKF.Index] + g.NLAS[g.AKF.Index] + g.NDIR[g.AKF.Index]
							ln.DUNGART = g.DGART[g.AKF.Index]
							g.C1[0] = g.C1[0] + g.NDIR[g.AKF.Index] //! Summe miner. Duengung
							if g.C1[0] < 0 {
								g.C1[0] = 0
							}
						}
					}

					if g.NDOY1[g.AKF.Index] < 10 {
						if g.NDOY1[g.AKF.Index] == 0 {
							if zeit == g.SAAT[g.AKF.Index] {
								nmin30 := 0.0
								for i := 0; i < 3; i++ {
									nmin30 = nmin30 + g.C1[i]
								}
								ndung := math.Max(g.NDEM1[g.AKF.Index]-nmin30, 0)
								g.NFERTSIM = g.NFERTSIM + ndung
								ln.DDAT1 = g.Kalender(zeit)
								ln.DMENG1 = ndung
								g.DSUMM = g.DSUMM + ndung
								if runErr = writeFertilizerEvent(ln.DUNGART, ndung, 0); runErr != nil {
									return finishedCycle, runErr
								}
							}
						} else {
							if g.INTWICK.Num == g.NDOY1[g.AKF.Index] {
								nmin30 := 0.0
								for i := 0; i < 3; i++ {
									nmin30 = nmin30 + g.C1[i]
								}
								ndung := math.Max(g.NDEM1[g.AKF.Index]-nmin30, 0)
								g.NFERTSIM = g.NFERTSIM + ndung
								ln.DDAT1 = g.Kalender(zeit)
								ln.DMENG1 = ndung
								g.DSUMM = g.DSUMM + ndung
								g.NDOY1[g.AKF.Index] = 0
								if runErr = writeFertilizerEvent(ln.DUNGART, ndung, 0); runErr != nil {
									return finishedCycle, runErr
								}
							}
						}
					} else {
						if g.TAG.Num > g.NDOY1[g.AKF.Index] && g.TAG.Num < 210 && g.NDOY1[g.AKF.Index] < 365 {
							if (g.TEMP[g.TAG.Index] + g.TEMP[g.TAG.Index-1] + g.TEMP[g.TAG.Index-2] + g.TEMP[g.TAG.Index-3] + g.TEMP[g.TAG.Index-4]) > 20 {
								if (g.REGEN[g.TAG.Index]+g.REGEN[g.TAG.Index-1]) < 0.4 && g.REGEN[g.TAG.Index+1] < 4 {
									nmin30 := 0.0
									for i := 0; i < 3; i++ {
										nmin30 = nmin30 + g.C1[i]
									}
									ndung := math.Max(g.NDEM1[g.AKF.Index]-nmin30, 0)
									g.NFERTSIM = g.NFERTSIM + ndung
									ln.DDAT1 = g.Kalender(zeit)
									ln.DMENG1 = ndung
									g.DSUMM = g.DSUMM + ndung
									g.NDOY1[g.AKF.Index] = 370
									if runErr = writeFertilizerEvent(ln.DUNGART, ndung, 0); runErr != nil {
										return finishedCycle, runErr
									}
								}
							}
						}
					}
					if g.NDOY2[g.AKF.Index] < 10 {
						if g.INTWICK.Num == g.NDOY2[g.AKF.Index] {
							nminw2 := 0.0
							for i := 0; i < min(g.WURZ, 9); i++ {
								nminw2 = nminw2 + g.C1[i]
							}
							ndung := math.Max(g.NDEM2[g.AKF.Index]-nminw2, 0)
							g.NFERTSIM = g.NFERTSIM + ndung
							ln.DDAT2 = g.Kalender(zeit)
							ln.DMENG2 = ndung
							g.DSUMM = g.DSUMM + ndung
							g.NDOY2[g.AKF.Index] = 0
							if runErr = writeFertilizerEvent(ln.DUNGART, ndung, 0); runErr != nil {
								return finishedCycle, runErr
							}
						}
					} else {
						if g.TAG.Num == g.NDOY2[g.AKF.Index] {
							nminw2 := 0.0
							for i := 0; i < min(g.WURZ, 9); i++ {
								nminw2 = nminw2 + g.C1[i]
							}
							ndung := math.Max(g.NDEM2[g.AKF.Index]-nminw2, 0)
							g.NFERTSIM = g.NFERTSIM + ndung
							ln.DDAT2 = g.Kalender(zeit)
							ln.DMENG2 = ndung
							g.DSUMM = g.DSUMM + ndung
							if runErr = writeFertilizerEvent(ln.DUNGART, ndung, 0); runErr != nil {
								return finishedCycle, runErr
							}
						}
					}
					if g.NDOY3[g.AKF.Index] < 10 {
						if g.INTWICK.Num == g.NDOY3[g.AKF.Index] {
							nminw3 := 0.0
							for i := 0; i < min(g.WURZ, 9); i++ {
								nminw3 = nminw3 + g.C1[i]
							}
							ndung := math.Max(g.NDEM3[g.AKF.Index]-nminw3, 0)
							g.NFERTSIM = g.NFERTSIM + ndung
							ln.DDAT3 = g.Kalender(zeit)
							ln.DMENG3 = ndung
							g.DSUMM = g.DSUMM + ndung
							g.NDOY3[g.AKF.Index] = 0
							if runErr = writeFertilizerEvent(ln.DUNGART, ndung, 0); runErr != nil {
								return finishedCycle, runErr
							}
						}
					} else {
						if g.TAG.Num == g.NDOY3[g.AKF.Index] {
							nminw3 := 0.0
							for i := 0; i < min(g.WURZ, 9); i++ {
								nminw3 = nminw3 + g.C1[i]
							}
							ndung := math.Max(g.NDEM3[g.AKF.Index]-nminw3, 0)
							g.NFERTSIM = g.NFERTSIM + ndung
							ln.DDAT3 = g.Kalender(zeit)
							ln.DMENG3 = ndung
							g.DSUMM = g.DSUMM + ndung
							if runErr = writeFertilizerEvent(ln.DUNGART, ndung, 0); runErr != nil {
								return finishedCycle, runErr
							}
						}
					}
				}
			}
		}
	}
	//! ++++++++++++++++++++++ Adaptation tillage to automatic sowing/harvest (if harvest later) +++++++++++++++++++++++
	if subd == 1 {
		if zeit == g.EINTE[g.NTIL.Index+1] {
			if g.SAAT[g.AKF.Index] > 0 && g.ERNTE[g.AKF.Index] == 0 {
				g.EINTE[g.NTIL.Index+1] = g.EINTE[g.NTIL.Index+1] + 2
			}
		}
	}
	if g.SAAT[g.AKF.Index] > 0 && g.EINTE[g.NTIL.Index+1] > g.SAAT[g.AKF.Index] && g.EINTE[g.NTIL.Index+1] <= g.ERNTE[g.AKF.Index] {
		//invalid tillage date
		return finishedCycle, fmt.Errorf("tillage date %s before harvest %s at %s", g.Kalender(g.EINTE[g.NTIL.Index+1]), g.Kalender(g.ERNTE[g.AKF.Index]+1), g.PKT)
	}

	// ----------------------------------------------------------------------------------------------------------------
	if zeit == g.EINTE[g.NTIL.Index+1]+1 && subd == 1 {
		var NFOSUM, NAOSUM, nmifosum, nmiaosum, CSUM float64
		if g.EINT[g.NTIL.Index] > 0 {
			mixtief := math.Round(g.EINT[g.NTIL.Index] / g.DZ.Num)

			layerList := make(map[string]interface{})
			for z := 0; z < int(mixtief); z++ {
				// Vollstaendige Durchmischung bis Bearbeitungstiefe
				NFOSUM = NFOSUM + g.NFOS[z]
				NAOSUM = NAOSUM + g.NAOS[z]
				nmifosum = nmifosum + g.MINFOS[z]
				nmiaosum = nmiaosum + g.MINAOS[z]
				CSUM = CSUM + g.C1[z]
			}
			if g.TILART[g.NTIL.Index] == 1 {
				for z := 0; z < int(mixtief); z++ {
					g.NFOS[z] = NFOSUM / mixtief
					g.NAOS[z] = NAOSUM / mixtief
					g.MINFOS[z] = nmifosum / mixtief
					g.MINAOS[z] = nmiaosum / mixtief
					g.C1[z] = CSUM / mixtief
					if g.C1[z] < 0 {
						g.C1[z] = 0
					}
					layerList[fmt.Sprintf("NminLayer%d", z+1)] = g.C1[z]
					layerList[fmt.Sprintf("NFOSLayer%d", z+1)] = g.NFOS[z]
					layerList[fmt.Sprintf("NAOSLayer%d", z+1)] = g.NAOS[z]
				}
			}

			runErr = g.managementConfig.WriteManagementEvent(NewManagementEvent(Tillage, zeit, layerList, g))
			if runErr != nil {
				return finishedCycle, runErr
			}
		}
		g.NTIL.Inc()
	}
	if subd == 1 {
		// Aufruf Mineralisations Subroutine
		mineral(g, l)
	}
	if zeit == g.ERNTE[g.AKF.Index] && subd == 1 {

		//NSA(NDG)  = Oberird. Zufuhr schnell mineralisierbarer org. Substanz aus Ernterückständen (kg N/ha)
		//NLA(NDG)  = Oberird. Zufuhr langsam mineralisierbarer org. Substanz aus Ernterückständen (kg N/ha)
		//NDI (NDG) = Unterird. Zufuhr sofort verfügbarer N aus Ernterückständen (kg N/ha)
		//NUSA(NDG) = Zufuhr schnell mineralisierbarer org. Substanz aus Ernterückständen (kg N/ha) (wird entspr. Wurzelverteilung verteilt)
		//NULA(NDG) = Unterird. Zufuhr langsam mineralisierbarer org. Substanz aus Ernterückständen (kg N/ha) (wird entspr. Wurzelverteilung verteilt)
		var NDI, NSA, NLA, NUSA, NULA float64
		if g.AKF.Num != 1 {
			NDI, NSA, NLA, NUSA, NULA, ln.NRESID = resid(g, ln, hPath)
			runErr = g.managementConfig.WriteManagementEvent(NewManagementEvent(Harvest, zeit, map[string]interface{}{
				"Residue": ln.NRESID,
			}, g))
			if runErr != nil {
				return finishedCycle, runErr
			}
		}
		g.NFOS[0] = g.NFOS[0] + NSA
		g.NAOS[0] = g.NAOS[0] + NLA
		for i := 0; i < g.WURZ; i++ {
			g.NFOS[i] = g.NFOS[i] + NUSA*g.WUANT[i]
			g.NAOS[i] = g.NAOS[i] + NULA*g.WUANT[i]
		}

		ln.NUPTAKE = g.PESUM
		g.DSUMM = g.DSUMM + NDI
		g.PESUM = g.PESUM - (NSA + NLA + NDI)

		if g.YORGAN == 0 {
			if g.YIFAK == 0.99 {
				g.YIELD = g.OBMAS - 820
			} else {
				g.YIELD = g.OBMAS * g.YIFAK
			}
		} else {
			g.YIELD = g.WORG[g.YORGAN-1] * g.YIFAK
		}

		// +++++++++++++++  fuer langzeitlauf und Ausgabe +++++++++++++++
		hyear, _, _ := KalenderDate(zeit)

		REDUKAV := g.REDUKSUM / float64(g.ERNTE[g.AKF.Index]-g.SAAT[g.AKF.Index])
		TRRELAV := g.TRRELSUM / float64(g.ERNTE[g.AKF.Index]-g.SAAT[g.AKF.Index])
		if ln.DOMENG1 == 0 {
			ln.DODAT = "----------"
			ln.DUNGART = "---"
		}
		// +++++++++++ VERLEGUNG DES TILLAGETERMINS BEI ORGANISCHER DüNGUNG NACH ERNTE +++++++++++++

		TDAT2 := g.Kalender(g.EINTE[g.NTIL.Index+1])

		if g.AKF.Num > 1 {
			NAOSAKT := (g.NAOS[0] + g.NAOS[1] + g.NAOS[2])
			NFOSAKT := (g.NFOS[0] + g.NFOS[1] + g.NFOS[2])
			NMIN1 := g.C1[0] + g.C1[1] + g.C1[2]
			NMIN2 := 0.0
			for i := 0; i < 15; i++ {
				NMIN2 = NMIN2 + g.C1[i]
			}
			if ln.DMENG1 == 0 {
				ln.DDAT1 = "------------"
			} else {
				ln.DDAT1 = " " + ln.DDAT1 + " "
			}
			if ln.DMENG2 == 0 {
				ln.DDAT2 = "------------"
			} else {
				ln.DDAT2 = " " + ln.DDAT2 + " "
			}

			if ln.DMENG3 == 0 {
				ln.DDAT3 = "------------"
			} else {
				ln.DDAT3 = " " + ln.DDAT3 + " "
			}
			// fill gaps in BBCH_DOY array, for the case that some stages are not reached, or were skipped
			fillBBCHgaps := func(doyArr *[100]int) {
				lastValue := 0
				for i := 0; i < 100; i++ {
					if doyArr[i] == 0 {
						doyArr[i] = lastValue
					} else {
						lastValue = doyArr[i]
					}
				}
			}

			output.EmergDOY = g.DEV[1]
			output.AnthDOY = g.DEV[4]
			output.MatDOY = g.DEV[5]
			output.HarvestYear = hyear
			output.HarvestDOY = g.TAG.Index + 1
			fillBBCHgaps(&g.BBCH_DOY)
			output.BBCH_DOY = g.BBCH_DOY
			fillBBCHgaps(&g.BBCH_TIME)
			output.BBCH_DATE = convertToDate(g.BBCH_TIME, g)
			output.Crop = g.CropTypeToString(g.FRUCHT[g.AKF.Index], true)
			output.Yield = g.YIELD
			output.Biomass = g.OBMAS
			output.Roots = g.WORG[0]
			output.LAImax = g.LAIMAX
			output.Nfertil = g.NFERTSIM
			output.Irrig = g.IRRISIM
			output.Nuptake = ln.NUPTAKE
			output.Suptake = g.SUPTAKE
			output.Nagb = ln.NAGB
			output.ETcG = g.ETC0
			output.ETaG = g.ETAG
			output.TraG = g.TRAG
			output.PerG = g.PERG
			output.SWCS1 = g.SWCS1
			output.SWCS2 = g.SWCS2
			output.SWCA1 = g.SWCA1
			output.SWCA2 = g.SWCA2
			output.SWCM1 = g.SWCM1
			output.SWCM2 = g.SWCM2
			output.SoilN1 = g.NALTOS/g.NAKT*(1-g.NAKT) + NAOSAKT + NFOSAKT
			output.Nmin1 = NMIN1
			output.Nmin2 = NMIN2
			//Smin1 and Smin2 and Smin3
			output.Smin1 = g.S1[0] + g.S1[1] + g.S1[2]
			output.Smin2 = g.S1[3] + g.S1[4] + g.S1[5]
			output.Smin3 = g.S1[6] + g.S1[7] + g.S1[8]
			output.NLeaG = g.NLEAG
			output.TRRel = TRRELAV
			output.Reduk = REDUKAV
			output.DryD1 = g.DRYD1
			output.DryD2 = g.DRYD2
			output.Nresid = ln.NRESID
			output.Orgdat = ln.DODAT
			output.Type = ln.DUNGART
			output.OrgN = ln.DOMENG1
			output.NDat1 = ln.DDAT1
			output.N1 = ln.DMENG1
			output.Ndat2 = ln.DDAT2
			output.N2 = ln.DMENG2
			output.Ndat3 = ln.DDAT3
			output.N3 = ln.DMENG3
			output.Tdat = TDAT2
			output.Code = g.POLYD
			output.NotStableErr = g.C1NotStableErr
			output.PARSUM = g.PARSUM
			finishedCycle = true
		}
		g.NFERTSIM = 0
		g.IRRISIM = 0
		g.ETAG = 0
		g.TRAG = 0
		g.PERG = 0
		g.NLEAG = 0
		ln.DMENG1 = 0
		ln.DMENG2 = 0
		ln.DMENG3 = 0
		ln.DOMENG1 = 0
		ln.NRESID = 0
		g.PARSUM = 0

		if g.DAUERKULT {
			if g.JN[g.AKF.Index] == 0 || g.JN[g.AKF.Index] == 1 {
				g.WORG[3], g.WORG[4] = 0, 0
				g.WORG[1] = math.Max(g.WORG[1]*(1-g.YIFAK), 720)
				g.WORG[2] = math.Max(g.WORG[2]*(1-g.YIFAK), 100)
				g.PESUM = math.Max((g.PESUM-g.WORG[0]*g.WUGEH)*(1-g.YIFAK), 820*g.GEHOB+g.WORG[0]*g.WUGEH)
				g.OBMAS = g.WORG[1] + g.WORG[2]
				g.WUMAS = g.WORG[0]
			} else {
				for i := range g.WORG {
					g.WORG[i] = 0
				}
			}
		} else {
			for i := range g.WORG {
				g.WORG[i] = 0
			}
		}
		if g.ODU[g.AKF.Index] == 1 && g.ORGTIME[g.AKF.Index] == "H" {
			g.ZTDG[g.AKF.Index] = zeit + g.ORGDOY[g.AKF.Index]
		}
		g.AKF.Inc()

		if g.SAAT2[g.AKF.Index] <= zeit && g.AUTOMAN {
			if g.ODU[g.AKF.Index-1] == 1 && g.ORGTIME[g.AKF.Index-1] == "H" {
				g.NAOS[0] = g.NAOS[0] + g.NLAS[g.AKF.Index-1]
				ln.DODAT = g.Kalender(zeit)
				ln.DOMENG1 = g.NSAS[g.AKF.Index-1] + g.NLAS[g.AKF.Index-1] + g.NDIR[g.AKF.Index-1]
				ln.DUNGART = g.DGART[g.AKF.Index-1]
				g.DSUMM = g.DSUMM + g.NDIR[g.AKF.Index-1] //! SUMME MINER. DUENGUNG
				g.EINTE[g.NTIL.Index+1] = zeit + 1
				TDAT2 = g.Kalender(g.EINTE[g.NTIL.Index+1])
				ln.DDAT1 = "------------"
				ln.DDAT2 = "------------"
				ln.DDAT2 = "------------"
				output.SowDate = "SKIPPED"
				output.SowDOY = 0
				output.EmergDOY = 0
				output.AnthDOY = 0
				output.MatDOY = 0
				output.HarvestYear = 0
				output.HarvestDOY = 0
				output.Crop = "000"
				output.Yield = 0
				output.Biomass = 0
				output.Roots = 0
				output.LAImax = 0
				output.Nfertil = 0
				output.Irrig = 0
				output.Nuptake = 0
				output.Nagb = 0
				output.ETcG = 0
				output.ETaG = 0
				output.TraG = 0
				output.PerG = 0
				output.SWCS1 = 0
				output.SWCS2 = 0
				output.SWCA1 = 0
				output.SWCA2 = 0
				output.SWCM1 = 0
				output.SWCM2 = 0
				output.SoilN1 = 0
				output.Nmin1 = 0
				output.Nmin2 = 0
				output.NLeaG = 0
				output.TRRel = 0
				output.Reduk = 0
				output.DryD1 = 0
				output.DryD2 = 0
				output.Nresid = 0
				output.Orgdat = ln.DODAT
				output.Type = ln.DUNGART
				output.OrgN = ln.DOMENG1
				output.NDat1 = ln.DDAT1
				output.N1 = ln.DMENG1
				output.Ndat2 = ln.DDAT2
				output.N2 = ln.DMENG2
				output.Ndat3 = ln.DDAT3
				output.N3 = ln.DMENG3
				output.Tdat = TDAT2
				output.Code = g.FCODE
				output.NotStableErr = g.C1NotStableErr
				output.PARSUM = 0

				ln.DOMENG1 = 0
				g.AKF.Inc()
				finishedCycle = true
			}
		}
		pinit(g)

		resetBBCHPermaCulture := func(cp CropType, bbch int) {
			if g.FRUCHT[g.AKF.Index] == cp {
				g.BBCH = bbch
				for i := bbch + 1; i < 100; i++ {
					g.BBCH_DOY[i] = 0
					g.BBCH_TIME[i] = 0
				}
			}
		}
		if g.FRUCHT[g.AKF.Index] != GRE &&
			g.FRUCHT[g.AKF.Index] != GR &&
			g.FRUCHT[g.AKF.Index] != AA {
			g.PESUM = 0
			g.WURZ = 0
			g.LAI = 0
			g.OBMAS = 0
			g.WUMAS = 0
			g.INTWICK.SetByIndex(-1)
			g.BBCH = 0
			g.ASPOO = 0
			g.VERNTAGE = 0
			// clear g.BBCH_DOY
			for i := 0; i < 100; i++ {
				g.BBCH_DOY[i] = 0
				g.BBCH_TIME[i] = 0
			}
		} else {
			resetBBCHPermaCulture(GRE, 45)
			resetBBCHPermaCulture(GR, 45)
			resetBBCHPermaCulture(AA, 39)
		}
	}
	// ---------- Aufruf N-Verlagerung -------------------------
	nmove(wdt, subd, zeit, g, l)
	return finishedCycle, nil
}

// mineral
func mineral(g *GlobalVarsMain, l *NitroSharedVars) {
	//! ------------------------------------- Mineralisation in Abh. von Temperatur und Wassergehalt ------------
	//! Inputs:
	//! IZM                       = bodenartspezifische Mineralisierungstiefe
	//! TEMP(TAG)                 = Tagesmitteltemperatur vom TAG (°C)
	//! TSOIL(0,Z)                = Bodentemperatur am Anfang Zeitschritt in Schicht Z  (°C)
	//! WG(0,Z)                   = Wassergehalt am Anfang Zeitschritt in Schicht Z  (cm^3/cm^3)
	//! WNOR(Z)                   = NORM-FK (ohne Wasserstau) in Schicht Z (cm^3/cm^3)
	//! WMIN(Z)                   = Wassergehalt bei PWP in Schicht Z (cm^3/cm^3)
	//! PORGES(Z)                 = Gesamtporenvolumen in Schicht Z  (cm^3/cm^3)
	//! MINAOS(Z)                 = bereits mineralisierter langsamer N-Pool in Z (kg N/ha)
	//! MINFOS(Z)                 = bereits mineralisierter schneller N-Pool in Z (kg N/ha)
	//! DSUMM                     = Summe der mineralischen Düngung (kg N/ha)
	//! NH4SUM                    = Summe ammoniakalischer N in Dünger
	//! UMS                       = Summe des bereits gelösten mineralischen N (kg N/ha)
	//! NH4UMS                    = Summe des nitrifizierten Nicht-Nitratanteils des mineralischen Düngers (kg N/ha)

	//! ----------------------------------------------------------------------------------------------------------
	var DTOTALN, DMINFOS, MIRED [4]float64

	//---------------------  Mineralisation  --------------------
	num := g.IZM / g.DZ.Index
	for z := 1; z <= num; z++ {
		zIndex := z - 1
		TEMPBO := (g.TD[z] + g.TD[z-1]) / 2
		// --------- Berechnung Mineralisationskoeffizienten ---------
		// ----------- in Abhängigkeit von TEMP UND WASSER -----------
		// - Umsetzung von mineralischen Düngern
		KTD := .4
		if TEMPBO > 0 {
			// Reaktionskoeffizient der schwer abbaubaren Fraktion
			kt0 := 4000000000. * math.Exp(-8400./(TEMPBO+273.16))
			// Reaktionskoeffizient der leicht abbaubaren Fraktion
			kt1 := 5.6e+12 * math.Exp(-9800./(TEMPBO+273.16))
			// Reduktionsfaktoren bei suboptimalem Wassergehalt
			if g.WG[0][zIndex] <= g.WNOR[zIndex] && g.WG[0][zIndex] >= g.WRED {
				MIRED[zIndex] = 1
			} else if g.WG[0][zIndex] < g.WRED && g.WG[0][zIndex] > g.WMIN[zIndex] {
				MIRED[zIndex] = (g.WG[0][zIndex] - g.WMIN[zIndex]) / (g.WRED - g.WMIN[zIndex])
			} else if g.WG[0][zIndex] > g.WNOR[zIndex] {
				MIRED[zIndex] = (g.PORGES[zIndex] - g.WG[0][zIndex]) / (g.PORGES[zIndex] - g.WNOR[zIndex])
			} else {
				MIRED[zIndex] = 0
			}
			if MIRED[zIndex] < 0 {
				MIRED[zIndex] = 0
			}
			if MIRED[zIndex] > 1 {
				MIRED[zIndex] = 1
			}
			// Mineralisation der schwer abbaubaren Fraktion
			DTOTALN[zIndex] = kt0 * g.NAOS[zIndex] * MIRED[zIndex]

			if DTOTALN[zIndex] < 0 {
				DTOTALN[zIndex] = 0
			}
			g.NAOS[zIndex] = g.NAOS[zIndex] - DTOTALN[zIndex]
			// Mineralisation der leicht abbaubaren Fraktion
			DMINFOS[zIndex] = kt1 * g.NFOS[zIndex] * MIRED[zIndex]
			if DMINFOS[zIndex] < 0 {
				DMINFOS[zIndex] = 0
			}
			g.NFOS[zIndex] = g.NFOS[zIndex] - DMINFOS[zIndex]
			if z == 1 {
				l.DUMS[zIndex] = KTD * MIRED[zIndex] * (g.DSUMM - g.UMS)
				l.DNH4UMS[zIndex] = KTD * MIRED[zIndex] * (g.NH4Sum - g.NH4UMS) //!Nitrifikation pro Zeitschritt (kg N/ha)

			} else {
				l.DUMS[zIndex] = 0
				l.DNH4UMS[zIndex] = 0

			}
			FN2oNit := (0.4*(g.WG[0][zIndex]/g.PORGES[zIndex]) - 1.04) / (g.WG[0][zIndex]/g.PORGES[zIndex] - 1.04) * 0.0016 //! Faktor N2O aus Nitrifikation
			N2oNIT := (l.DNH4UMS[zIndex] + DTOTALN[zIndex] + DMINFOS[zIndex]) * FN2oNit                                     //! N2O emission aus Nitrifikation pro Zeitschritt (kg N/ha)

			// Mineralisationssumme => Quellterm ( dn(z) )
			g.DN[zIndex] = DTOTALN[zIndex] + DMINFOS[zIndex] + l.DUMS[zIndex] - N2oNIT

			g.MINAOS[zIndex] = g.MINAOS[zIndex] + DTOTALN[zIndex]
			g.MINFOS[zIndex] = g.MINFOS[zIndex] + DMINFOS[zIndex]
			g.UMS = g.UMS + l.DUMS[zIndex]
			g.NH4UMS = g.NH4UMS + l.DNH4UMS[zIndex]
			g.N2onitsum = g.N2onitsum + N2oNIT
			g.N2onitDaily = N2oNIT
			g.MINSUM = g.MINSUM + g.DN[zIndex] - l.DUMS[zIndex]
		} else {
			if z == 1 {
				// Reduktionsfaktoren bei suboptimalem Wassergehalt
				if g.WG[0][zIndex] < g.W[zIndex] && g.WG[0][zIndex] > g.WRED {
					MIRED[zIndex] = 1
				} else if g.WG[0][zIndex] < g.WRED {
					MIRED[zIndex] = (g.WG[0][zIndex] - g.WMIN[zIndex]) / (g.WRED - g.WMIN[zIndex])
				} else if g.WG[0][zIndex] > g.W[zIndex]+.01 && g.WG[0][zIndex] < g.PORGES[0] {
					MIRED[zIndex] = (g.PORGES[0] - g.WG[0][zIndex]) / (g.PORGES[0] - g.W[zIndex])
				} else if g.WG[0][zIndex] > g.PORGES[0] {
					MIRED[zIndex] = 0
				} else {
					MIRED[zIndex] = 1
				}
				if MIRED[zIndex] < 0 {
					MIRED[zIndex] = 0
				}
				l.DUMS[zIndex] = 0.4 * MIRED[zIndex] * (g.DSUMM - g.UMS)
				l.DNH4UMS[zIndex] = 0.4 * MIRED[zIndex] * (g.NH4Sum - g.NH4UMS) //!Nitrifikation pro Zeitschritt (kg N/ha)

			} else {
				l.DUMS[zIndex] = 0
				l.DNH4UMS[zIndex] = 0
			}
			g.UMS = g.UMS + l.DUMS[zIndex]

			g.NH4UMS = g.NH4UMS + l.DNH4UMS[zIndex]
			FN2oNit := (0.4*(g.WG[0][zIndex]/g.PORGES[zIndex]) - 1.04) / (g.WG[0][zIndex]/g.PORGES[zIndex] - 1.04) * 0.0016 //! Faktor N2O aus Nitrifikation
			N2ONIT := l.DNH4UMS[zIndex] * FN2oNit                                                                           //! N2O emission aus Nitrifikation pro Zeitschritt (kg N/ha)
			g.N2onitsum = g.N2onitsum + N2ONIT
			g.N2onitDaily = N2ONIT
			g.DN[zIndex] = l.DUMS[zIndex] - N2ONIT
		}
	}
}

// nmove
func nmove(wdt float64, subd int, zeit int, g *GlobalVarsMain, l *NitroSharedVars) {
	// ---------------------      N-Verlagerung konvektions-Dispersionsgleichung ---------------------
	//Inputs:
	// DV                        = Dispersionslänge (cm)
	// FLUSS=                    = Infiltration durch Bodenoberfläche (cm/d)
	// Q1(Z)                     = Fluss durch Untergrenze (cm/d)
	// QDRAIN                    = Ausfluss in Drainrohr (cm/d)
	// DRAIDEP                   = Tiefe des Drains (dm)
	// AD                        = Faktor für Diffusivität?
	// DZ                        = Schichtdicke (cm)
	// WG(0,Z)                   = Wassergehalt am Anfang Zeitschritt in Schicht Z  (cm^3/cm^3)
	// WNOR(Z)                   = NORM-FK (ohne Wasserstau) in Schicht Z (cm^3/cm^3)
	// WMIN(Z)                   = Wassergehalt bei PWP in Schicht Z (cm^3/cm^3)
	// PORGES(Z)                 = Gesamtporenvolumen in Schicht Z  (cm^3/cm^3)
	// PE(Z)                     = N-Aufnahme Pflanze in Schicht Z (kg N/ha)
	// C1(Z)                     = Nmin-gehalt der Schicht Z (kg N/ha)
	// DN(Z)                     = Quellterm aus Mineralisation (kg N/ha) in Schicht Z
	// OUTN                      = Tiefe für Auswaschungsberechnung (dm)
	var Carray [22]float64
	for z := 0; z < g.N; z++ {
		// --- Berechnung des Diffusionskoeffizienten am unteren Kompartimentrand ---
		l.D[z] = 2.14 * (g.AD * math.Exp((g.WG[0][z]+g.WG[0][z+1])*5) / ((g.WG[0][z] + g.WG[0][z+1]) / 2)) * wdt
		if subd == 1 {
			if g.PE[z] > g.C1[z]-.5 {
				g.PE[z] = (g.C1[z] - .5)
			}
			if g.PE[z] < 0 {
				g.PE[z] = 0
			}
			g.PESUM = g.PESUM + g.PE[z]
			g.AUFNASUM = g.AUFNASUM + g.PE[z]
			if g.C1[z]-g.PE[z] < 0 {
				g.C1[z] = 0
			} else {
				g.C1[z] = g.C1[z] - g.PE[z]
			}
		}
		Carray[z+1] = (g.C1[z] + g.DN[z]*wdt/2) / (g.WG[0][z] * g.DZ.Num * 100)
		if Carray[z+1] < 0 {
			Carray[z+1] = 0
		}
	}
	// --------------------- Verlagerung nach unten ---------------------
	g.Q1[0] = g.FLUSS0 * wdt
	for zIndex0 := 0; zIndex0 < g.N; zIndex0++ {
		zIndex1 := zIndex0 + 1
		// Porenwassergeschwindigkeit V
		l.V[zIndex0] = math.Abs(g.Q1[zIndex1] / ((g.W[zIndex0] + g.W[zIndex0+1]) * .5))
		l.DB[zIndex0] = (g.WG[0][zIndex0]+g.WG[0][zIndex0+1])/2*(l.D[zIndex0]+g.DV*l.V[zIndex0]) - 0.5*wdt*math.Abs(g.Q1[zIndex1]) + 0.5*wdt*math.Abs((g.Q1[zIndex1]+g.Q1[zIndex1-1])/2)*l.V[zIndex0]
		if zIndex1 == 1 {
			cVar := Carray[zIndex1] - Carray[zIndex1+1]
			dbVar := -l.DB[zIndex0]
			num100 := math.Pow(g.DZ.Num, 2)
			l.DISP[zIndex0] = dbVar * cVar / num100
		} else if zIndex1 < g.N {
			l.DISP[zIndex0] = l.DB[zIndex0-1]*(Carray[zIndex1-1]-Carray[zIndex1])/math.Pow(g.DZ.Num, 2) - l.DB[zIndex0]*(Carray[zIndex1]-Carray[zIndex1+1])/math.Pow(g.DZ.Num, 2)
		} else {
			l.DISP[zIndex0] = l.DB[zIndex0-1] * (Carray[zIndex1-1] - Carray[zIndex1]) / math.Pow(g.DZ.Num, 2)
		}
	}
	for z := 1; z <= g.N; z++ {
		z0 := z - 1
		if g.Q1[z] >= 0 && g.Q1[z-1] >= 0 {
			if z == g.DRAIDEP {
				l.KONV[z0] = (Carray[z]*g.Q1[z] + Carray[z]*g.QDRAIN - Carray[z-1]*g.Q1[z-1]) / g.DZ.Num
			} else {
				l.KONV[z0] = (Carray[z]*g.Q1[z] - Carray[z-1]*g.Q1[z-1]) / g.DZ.Num
			}
		} else if g.Q1[z] >= 0 && g.Q1[z-1] < 0 {
			if z > 1 {
				if z == g.DRAIDEP {
					l.KONV[z0] = (Carray[z]*g.Q1[z] + Carray[z]*g.QDRAIN - Carray[z]*g.Q1[z-1]) / g.DZ.Num
				} else {
					l.KONV[z0] = (Carray[z]*g.Q1[z] - Carray[z]*g.Q1[z-1]) / g.DZ.Num
				}
			} else {
				l.KONV[z0] = Carray[z] * g.Q1[z] / g.DZ.Num
			}
		} else if g.Q1[z] < 0 && g.Q1[z-1] < 0 {
			if z > 1 {
				l.KONV[z0] = (Carray[z+1]*g.Q1[z] - Carray[z]*g.Q1[z-1]) / g.DZ.Num
			} else {
				l.KONV[z0] = Carray[z+1] * g.Q1[z] / g.DZ.Num
			}
		} else if g.Q1[z] < 0 && g.Q1[z-1] >= 0 {
			l.KONV[z0] = (Carray[z+1]*g.Q1[z] - Carray[z-1]*g.Q1[z-1]) / g.DZ.Num
		}
	}
	g.DRAINLOSS = g.DRAINLOSS + g.QDRAIN*Carray[g.DRAIDEP]/g.DZ.Num*100*g.DZ.Num

	g.C1NotStable = ""
	for z := 0; z < g.N; z++ {
		cKonz := (Carray[z+1]*g.WG[0][z] + l.DISP[z] - l.KONV[z]) * g.DZ.Num * 100

		if cKonz < 0 {
			g.C1[z] = 0
			// C1 may be below 0 because of rounding issues, set it to 0
			// if C1 is significat below zero, there might be an instabily in the calculations
			if cKonz < g.C1stabilityVal {
				g.C1NotStable = "C1 unstable"
				g.C1NotStableErr = "C1 unstable"
			}
		} else {
			g.C1[z] = cKonz
		}
	}
	// this part will only be triggerd if a RPC service was connected at start
	if err := HermesRPCService.SendGV(g, zeit, wdt, subd); err != nil {
		fmt.Println(err)
	}
	if err := HermesRPCService.SendNV(l, zeit, wdt, subd); err != nil {
		fmt.Println(err)
	}

	if g.Q1[g.OUTN] > 0 {
		if g.OUTN < g.N {
			g.OUTSUM = g.OUTSUM + g.Q1[g.OUTN]*Carray[g.OUTN]/g.DZ.Num*100*g.DZ.Num + l.DB[g.OUTN-1]*(Carray[g.OUTN]-Carray[g.OUTN+1])/math.Pow(g.DZ.Num, 2)*100*g.DZ.Num
			if zeit > g.SAAT[g.AKF.Index] {
				g.NLEAG = g.NLEAG + g.Q1[g.OUTN]*Carray[g.OUTN]/g.DZ.Num*100*g.DZ.Num + l.DB[g.OUTN-1]*(Carray[g.OUTN]-Carray[g.OUTN+1])/math.Pow(g.DZ.Num, 2)*100*g.DZ.Num
			}
		} else {
			g.OUTSUM = g.OUTSUM + g.Q1[g.OUTN]*Carray[g.OUTN]/g.DZ.Num*100*g.DZ.Num
			if zeit > g.SAAT[g.AKF.Index] {
				g.NLEAG = g.NLEAG + g.Q1[g.OUTN]*Carray[g.OUTN]/g.DZ.Num*100*g.DZ.Num
			}
		}
	} else {
		if g.OUTN < g.N {
			g.OUTSUM = g.OUTSUM + g.Q1[g.OUTN]*Carray[g.OUTN+1]/g.DZ.Num*100*g.DZ.Num + l.DB[g.OUTN-1]*(Carray[g.OUTN]-Carray[g.OUTN+1])/math.Pow(g.DZ.Num, 2)*100*g.DZ.Num
			if zeit > g.SAAT[g.AKF.Index] {
				g.NLEAG = g.NLEAG + g.Q1[g.OUTN]*Carray[g.OUTN+1]/g.DZ.Num*100*g.DZ.Num + l.DB[g.OUTN-1]*(Carray[g.OUTN]-Carray[g.OUTN+1])/math.Pow(g.DZ.Num, 2)*100*g.DZ.Num
			}
		}
	}
	for z := 0; z < g.N; z++ {
		g.C1[z] = g.C1[z] + g.DN[z]*wdt/2

		if g.C1[z] < 0 {
			g.C1[z] = 0
		}
	}
	if zeit >= g.SAAT[g.AKF.Index] && zeit <= g.ERNTE2[g.AKF.Index] {
		g.PESUM = g.PESUM + g.SCHNORR
	}
}

// resid
func resid(g *GlobalVarsMain, ln *NitroBBBSharedVars, hPath *HFilePath) (NDI, NSA, NLA, NUSA, NULA, NRESID float64) {
	// ------------------------------- Mineralisationspotentiale aus Vorfruchtresiduen ---------------------------------------
	// Input:
	// Dauerkult$         = D = Dauerkultur / Permanent crop
	// JN(AKF)            = Anteil der exportierten Pflanzenrückstände(Fraktion) / Fraction of crop residues that are removed from the field
	// 						0 = Verbleib auf dem Feld / all residues remain on the field
	//                      1 = vollständige Entfernung / all above ground residues are removed from the field
	// 						2 = komplette Pflanze beibt auf dem Feld / complete plant remains on the field, no yield is harvested
	// PESUM              = aufgenommene N-Menge der Pflanze (kg N/ha) / N amount taken up by the crop (kg N/ha)

	CRONAM := hPath.cropn
	_, scanner, _ := Open(&FileDescriptior{FilePath: CRONAM, UseFilePool: true})
	var KOSTRO, NERNT, NKOPP, NWURA, NFAST float64
	for scanner.Scan() {
		CROP := scanner.Text()
		if g.ToCropType(CROP[0:3]) == g.FRUCHT[g.AKF.Index] {
			//! Korn-Stroh Verhältnis
			KOSTRO = ValAsFloat(CROP[4:7], CRONAM, CROP)
			// N im Erntegut (kg N/dt)
			NERNT = ValAsFloat(CROP[13:18], CRONAM, CROP)
			// N im Koppelprodukt (kg N/dt)
			NKOPP = ValAsFloat(CROP[25:30], CRONAM, CROP)
			//Wurzelanteil an Gesamt-N in Pflanze
			NWURA = ValAsFloat(CROP[36:40], CRONAM, CROP)
			//Schnell mineralisierbarer Anteil von N in Ernterückständen (Fraktion)
			NFAST = ValAsFloat(CROP[41:45], CRONAM, CROP)
			break
		}
	}
	if ln != nil {
		ln.NAGB = g.PESUM - (g.PESUM * NWURA)
	}
	var DGM, DGU float64
	// DGM = N amount from crop residues (kg N/ha)
	// DGU = N amount from roots (kg N/ha)
	if g.JN[g.AKF.Index] == 0 {
		// all residues remain on the field
		if g.DAUERKULT {
			DGM = (g.OBMAS - 820) * g.GEHOB
			DGU = 0
		} else {
			DGU = g.PESUM * NWURA
			DGM = (1 - g.JN[g.AKF.Index]) * (g.PESUM - g.PESUM*(1-NWURA)*NERNT/(NERNT+KOSTRO*NKOPP) - g.PESUM*NWURA)
		}
	} else if g.JN[g.AKF.Index] == 1 {
		// all residues are removed from the field
		if g.DAUERKULT {
			if g.FRUCHT[g.AKF.Index] == AA {
				DGM = 0
				DGU = g.PESUM * NWURA * 0.74
			} else {
				DGM = 0
				DGU = g.PESUM * NWURA * 0.2
			}
		} else {
			DGM = 0
			DGU = g.PESUM * NWURA
		}
	} else if g.JN[g.AKF.Index] == 2 {
		// complete plant remains on the field, no yield is harvested
		DGU = g.PESUM * NWURA
		DGM = g.PESUM - DGU
	} else {
		// JN is a fraction between 0 and 1 of residues that are removed from the field
		if g.DAUERKULT {
			DGU = g.PESUM * NWURA * 0.74
			DGM = g.PESUM - (g.OBMAS * g.JN[g.AKF.Index] * g.GEHOB)
		} else {
			DGU = g.PESUM * NWURA
			DGM = (1 - g.JN[g.AKF.Index]) * (g.PESUM - g.PESUM*(1-NWURA)*NERNT/(NERNT+KOSTRO*NKOPP) - g.PESUM*NWURA)
		}
	}
	if DGM < 0 {
		DGM = 0
	}
	if DGU < 0 {
		DGU = 0
	}
	NSA = DGM * NFAST        // N amount from above ground crop residues that decompose fast(kg N/ha)
	NUSA = DGU * NFAST       // N amount from roots that decompose fast(kg N/ha)
	NLA = DGM * (1 - NFAST)  // N amount from above ground crop residues that decompose slow(kg N/ha)
	NULA = DGU * (1 - NFAST) // N amount from roots that decompose slow(kg N/ha)
	NDI = 0.0
	NRESID = DGM // N residue from above ground crops(kg N/ha)
	return NDI, NSA, NLA, NUSA, NULA, NRESID
}

// pinit
func pinit(g *GlobalVarsMain) {

	if !g.DAUERKULT {
		g.PESUM = 0
		g.PESUMS = 0
		g.VERNTAGE = 0
		g.OBMAS = 0
		g.WUMAS = 0
		g.INTWICK.SetByIndex(-1)
		g.WURZ = 0
		g.PHYLLO = 0
	}
}

// END MODULE
