package hermes

import (
	"fmt"
	"math"
	"strings"
)

// potential sulfur mineralization
func sPotMin(g *GlobalVarsMain) {
	if g.CGEHALT[0] > 14 {
		g.SALTOS = 5000 * g.SGEHALT[0] * g.SAKT * float64(g.UKT[1])
	} else if g.CGEHALT[0] > 5 {
		g.SALTOS = 11000 * g.SGEHALT[0] * g.SAKT * float64(g.UKT[1])
	} else if g.CGEHALT[0] < 1 {
		g.SALTOS = 15000 * g.SGEHALT[0] * g.SAKT * float64(g.UKT[1])
	} else {
		g.SALTOS = 15000 * g.SGEHALT[0] * g.SAKT * float64(g.UKT[1])
	}
}

// sulfur fertilization event
func writeSulfurFertilizationEvent(fertName string, s interface{}, zeit int, g *GlobalVarsMain) error {
	err := g.managementConfig.WriteManagementEvent(NewManagementEvent(Fertilization, zeit, map[string]interface{}{
		"Fertilizer": fertName,
		"S":          s,
	}, g))
	return err
}

// Smin calculation, must be called before nitro()
func Sulfo(wdt float64, subd, zeit int, g *GlobalVarsMain, hPath *HFilePath) error {

	if subd == 1 {
		// apply observed values
		// if observed date is 0, apply as inital value
		if len(g.sMESS) > g.sMessIdx && (zeit == g.sMESS[g.sMessIdx]) {
			messDates := g.SI[g.sMESS[g.sMessIdx]]
			for z := 0; z < g.N; z++ {
				currS1 := g.S1[z]
				// apply only values > 0 (negative values count as not set)
				if messDates[z] > 0 {
					g.S1[z] = messDates[z]
				}
				g.SDiff += currS1 - g.S1[z]
			}
			g.sMessIdx++
			g.SDSUMM = 0 // but why?
		}
		if !g.AUTOFERT {
			//! +++++++++++++++++++++++++++++++++++++ Option real fertilization +++++++++++++++++++++++++++++++++++++++++++++++
			// IF ZEIT = ZTDG(NDG)+1 THEN
			if zeit == g.ZTDG[g.NDG.Index]+1 {
				// 		   LET SFOS(1) = SFOS(1) + SSAS(NDG)
				g.SFOS[0] += g.SSAS[g.NDG.Index]
				// 		   LET SAOS(1) = SAOS(1) + SLAS(NDG)
				g.SAOS[0] += g.SLAS[g.NDG.Index]
				// 		   LET DSUMM = DSUMM + SDIR(NDG)    ! Summe miner. Düngung
				g.SDSUMM += g.SDIR[g.NDG.Index]
				// 		   LET NDG = NDG + 1
				// do not increase fertilization index here ... this happens in nitro()
				if runErr := writeSulfurFertilizationEvent(g.DGART[g.NDG.Index], g.SDIR[g.NDG.Index], zeit, g); runErr != nil {
					return runErr
				}
			}
		}
		// TODO: Auto-fertilization
		// ---- Homogene Vermischung von mineralischem u. organ. S bei Bearbeitung ----
		// 		IF Zeit = EINTE(1)+1 then
		if zeit == g.EINTE[g.NTIL.Index+1]+1 {
			g.SFOSUM, g.SAOSUM, g.SSUM, g.SFSUM = 0, 0, 0, 0
			// 		   IF EINT(1) > 0 then
			if g.EINT[g.NTIL.Index] > 0 {
				mixDepth := math.Round(g.EINT[g.NTIL.Index] / g.DZ.Num)
				// 		 FOR Z = 1 TO EINT(1)/dz
				for z := 0; z < int(mixDepth); z++ {
					//! Vollständige Durchmischung bis Bearbeitungstiefe für Anfangsforfrucht
					// mix everything up until the mix depth for the initial crop
					//LET SFOSUM = SFOSUM + SFOS(Z)
					g.SFOSUM += g.SFOS[z]
					//LET SAOSUM = SAOSUM + SAOS(Z)
					g.SAOSUM += g.SAOS[z]
					//LET SSUM   = SSUM   + S1(Z)
					g.SSUM += g.S1[z]
					//LET SFSUM  = SFSUM  + SF(Z)
					g.SFSUM += g.SF[z]
					//NEXT Z
				}
				//FOR Z = 1 TO EINT(1)/dz
				for z := 0; z < int(mixDepth); z++ {
					//LET SFOS(Z) = SFOSUM/EINT(1)*dz
					g.SFOS[z] = g.SFOSUM / mixDepth
					//LET SAOS(Z) = SAOSUM/EINT(1)*dz
					g.SAOS[z] = g.SAOSUM / mixDepth
					//LET S1(z)   = SSUM/EINT(1)*dz
					g.S1[z] = g.SSUM / mixDepth
					//LET SF(z)   = SFSUM/EINT(1)*dz
					g.SF[z] = g.SFSUM / mixDepth
				}
			}
			// don't increment NTIL here, it is used in nitro()
		}

		sMineral(g)

		//IF ZEIT = ERNTE(AKF) THEN
		if zeit == g.ERNTE[g.AKF.Index] {
			var SSA, SLA, SDI float64
			//IF akf <> 1 then
			if g.AKF.Num != 1 {
				//CALL SRESID(SSA,SLA)
				SSA, SLA, SDI = sResid(g, hPath)
			}
			//LET SFOS(1) = SFOS(1) + SSA
			g.SFOS[0] += SSA
			//LET SAOS(1) = SAOS(1) + SLA
			g.SAOS[0] += SLA
			//LET DSUMM = DSUMM + SDI
			g.SDSUMM += SDI
			// 		   LET AKF = AKF+1
			// don't increment AKF here, it is used in nitro()
			g.SUPTAKE = g.PESUMS
		}
	}
	sMove(wdt, subd, g)
	return nil
}

// SUB SMINERAL
func sMineral(g *GlobalVarsMain) {
	// 		DIM DSAOS(4),DSFOS(4),MIRED(4)
	var DSAOS, DSFOS, MIRED, DUMS [4]float64

	// 	---------------------  Mineralisation  --------------------
	// FOR z = 1 to izm/dz
	num := g.IZM / g.DZ.Index
	for z := 0; z < num; z++ {
		TEMPBO := (g.TD[z+1] + g.TD[z]) / 2

		// --------- Berechnung Mineralisationskoeffizienten ---------
		// ----------- in Abhängigkeit von TEMP UND WASSER -----------
		//- Umsetzung von mineralischen Düngern
		ktd := 0.4
		// 	! Reduktion bei suboptimalem Wassergehalt
		if g.WG[0][z] <= g.WNOR[z] && g.WG[0][z] >= g.WRED {
			MIRED[z] = 1
		} else if g.WG[0][z] < g.WRED && g.WG[0][z] > g.WMIN[z] {
			MIRED[z] = (g.WG[0][z] - g.WMIN[z]) / (g.WRED - g.WMIN[z])
		} else if g.WG[0][z] > g.WNOR[z] {
			MIRED[z] = (g.PORGES[z] - g.WG[0][z]) / (g.PORGES[z] - g.WNOR[z])
		} else {
			MIRED[z] = 0
		}
		if MIRED[z] < 0 {
			MIRED[z] = 0
		}
		if MIRED[z] > 1 {
			MIRED[z] = 1
		}

		// 	! Temperaturabhängigkeit der Mineralisationskoeffizienten nur >0 Grad
		var kt0, kt1 float64
		if TEMPBO > 0 {
			//  Reaktionskoeffizient der schwer abbaubaren Fraktion
			//LET kt0 = 4000000000*exp(-8400/(TEMPBO+273.16))*dt
			kt0 = 4000000000 * math.Exp(-8400/(TEMPBO+273.16)) * g.DT.Num
			// Reaktionskoeffizient der leicht abbaubaren Fraktion
			//LET kt1 = 5.6e+12*exp(-9800/(TEMPBO+273.16))*dt
			kt1 = 5.6e+12 * math.Exp(-9800/(TEMPBO+273.16)) * g.DT.Num
		}

		//  Mineralisation der schwer abbaubaren Fraktion
		// 	LET DSAOS(z) = kt0*SAOS(z)*MIRED(Z)
		DSAOS[z] = kt0 * g.SAOS[z] * MIRED[z]
		// 	Reduktion des Pools um die mineralisierte Menge
		// 	LET SAOS(z) = SAOS(z) - DSAOS(z)
		g.SAOS[z] -= DSAOS[z]

		//  Mineralisation der leicht abbaubaren Fraktion
		// 	LET DSFOS(z) = kt1*SFOS(z)*MIRED(Z)
		DSFOS[z] = kt1 * g.SFOS[z] * MIRED[z]
		// 	Reduktion des Pools um die mineralisierte Menge
		// 	LET SFOS(z) = SFOS(z) - DSFOS(z)
		g.SFOS[z] -= DSFOS[z]

		// 	IF z = 1 THEN
		if z == 0 {
			// 	   LET DUMS(Z) = KTD * MIRED(Z) * DSUMM *dt
			DUMS[z] = ktd * MIRED[z] * g.SDSUMM * g.DT.Num
			// 	   LET DSUMM = DSUMM - DUMS(z)
			g.SDSUMM -= DUMS[z]

			// 	ELSE
		} else {
			// 	   LET DUMS(Z) = 0
			DUMS[z] = 0
			// 	END IF
		}
		// 	 Mineralisationssumme => Quellterm ( DNS(z) )
		// 	LET DNS(z) = DSAOS(z)+DSFOS(z)+DUMS(Z)
		g.DNS[z] = DSAOS[z] + DSFOS[z] + DUMS[z]
		// 	LET SMINAOS(z) = SMINAOS(z)+DSAOS(z)
		g.Sminaos[z] += DSAOS[z]
		// 	LET SMINFOS(z) = SMINFOS(z)+DSFOS(z)
		g.Sminfos[z] += DSFOS[z]
		// 	LET SMINSUM = SMINSUM + DNS(Z)-DUMS(Z)
		g.SMINSUM += g.DNS[z] - DUMS[z]
	}
}

// SUB SMOVE(#5)
func sMove(wdt float64, subd int, g *GlobalVarsMain) {
	// ---------------------      S-Verlagerung konvektions-Dispersionsgleichung ---------------------
	//Inputs:
	// DV                        = Dispersionslänge (cm)
	// KLOS                      = Loeslichkeitskonstante (1/d)
	// SKSAT					 = Saettigungs-Loesungskonzentration (g S/L)
	// SDV                       = Dispersionslängenkoeffizient (cm^2/d)
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
	// PES(Z)                    = S-Aufnahme Pflanze in Schicht Z (kg S/ha)
	// S1(Z)                     = Smin-gehalt der Schicht Z (kg S/ha)
	// DN(Z)                     = Quellterm aus Mineralisation (kg S/ha) in Schicht Z
	// OUTN                      = Tiefe für Auswaschungsberechnung (dm)
	// 		DIM Sarray(0:21)
	var Sarray [22]float64
	var DiffCoeff [21]float64 //Diffusionskoeffizienten pro layer

	for z := 0; z < g.N; z++ {
		z1 := z + 1
		// ! --- Berechnung des Diffusionskoeffizienten am unteren Kompartimentrand ---
		// LET D(Z) = D0S * (AD*EXP((WG(0,Z)+WG(0,Z+1))*5)/((WG(0,Z)+WG(0,Z+1))/2))*DT
		DiffCoeff[z] = 2.14 * (g.AD * math.Exp((g.WG[0][z]+g.WG[0][z+1])*5) / ((g.WG[0][z] + g.WG[0][z+1]) / 2)) * wdt
		// ** Loeslichkeitsobergrenzen und Loesungs-/Faellungsreaktion **
		// SKSAT ist die Saettigungs-Loesungskonzentration in Gramm S/Liter
		//  SF(z) ist die nicht gelöste Smin-Menge in kg S/ha
		//LET S(Z) = (S1(Z)-SF(Z))/(wg(0,z)*DZ*100)
		// Sarrray is the concentration of S in solution
		Sarray[z1] = (g.S1[z] - g.SF[z]) / (g.WG[0][z] * g.DZ.Num * 100)
		//IF S(Z) >= SKSAT then
		if Sarray[z1] >= g.SKSAT {
			// satturation is reached
			// cap Smin to SKSAT
			//LET S(Z) = SKSAT
			Sarray[z1] = g.SKSAT
			//LET SF(Z) = S1(Z) - S(Z)*(wg(0,z)*dz*100)
			g.SF[z] = g.S1[z] - Sarray[z1]*(g.WG[0][z]*g.DZ.Num*100)
		} else {
			// satturation is not reached
			// add Smin from SF
			//LET S(Z)  = S(Z) + (Sksat-S(Z)) * (1-EXP(-klos*(SF(Z)/(WG(0,z)*dz*100))))
			Sarray[z1] += (g.SKSAT - Sarray[z1]) * (1 - math.Exp(-g.KLOS*(g.SF[z]/(g.WG[0][z]*g.DZ.Num*100))))
			//LET SF(Z) = S1(Z) - S(Z)*(wg(0,z)*dz*100)
			g.SF[z] = g.S1[z] - Sarray[z1]*(g.WG[0][z]*g.DZ.Num*100)
		}
		if subd == 1 {
			//! Untere Begrenzung für Entleerung pro Schicht (0.02 kg S/ha)
			//IF PES(Z) > S1(Z)-SF(Z)-.02 THEN
			if g.PES[z] > g.S1[z]-g.SF[z]-0.02 {
				//LET PES(Z) = (S1(Z)-SF(Z)-.02)
				g.PES[z] = g.S1[z] - g.SF[z] - 0.02
			}
			//IF PES(Z) < 0 THEN LET PES(Z) = 0
			if g.PES[z] < 0 {
				g.PES[z] = 0
			}

			//LET PESUMS = PESUMS + PES(Z)/2
			g.PESUMS += g.PES[z] / 2

			g.S1[z] = g.S1[z] - g.PES[z]/2

		}

		// Umrechnung in Bodenloesungskonzentration (kg/ha --> g/l)
		// Quellen und Senken jeweils halb bei Beginn/Ende Zeitschritt
		// LET S(Z) = (S1(Z)-SF(Z) + DNS(Z)/2 - PES(Z)/2)/(wg(0,z)*DZ*100)
		Sarray[z1] = (g.S1[z] - g.SF[z] + g.DNS[z]*wdt/2) / (g.WG[0][z] * g.DZ.Num * 100)

	}
	var V, DB, DISP, KONV [21]float64
	// 		!--------------------- Verlagerung nach unten ---------------------
	// 		LET Q1(0) = FLUSS0*DT
	g.Q1[0] = g.FLUSS0 * wdt
	sqrdDZ := math.Pow(g.DZ.Num, 2)
	// 	! Berechnung der Dispersion in Abhängigkeit der Porenwassergeschwindigkeit
	// 	FOR Z = 1 TO N
	for z1 := 1; z1 <= g.N; z1++ {
		z0 := z1 - 1
		// 		!LET V(Z) = ABS(Q1(Z)/((WG(0,Z)+WG(0,Z+1))*.5))
		//LET V(Z) = ABS(Q1(Z)/((W(Z)+W(Z+1))*.5))
		V[z0] = math.Abs(g.Q1[z1] / ((g.WG[0][z0] + g.WG[0][z0+1]) * .5)) // note: in nitro it is W not WG[0]
		//LET DB(Z) = (WG(0,Z)+WG(0,Z+1))/2*(D(Z) + DV*V(Z))-.5*dz*ABS(q1(z))+.5*dt*ABS((q1(z)+q1(z-1))/2)*v(z)
		DB[z0] = (g.WG[0][z0]+g.WG[0][z0+1])/2*(DiffCoeff[z0]+g.SDV*V[z0]) - .5*wdt*g.DZ.Num*math.Abs(g.Q1[z1]) + .5*wdt*math.Abs((g.Q1[z1]+g.Q1[z1-1])/2)*V[z0]
		if z1 == 1 {
			// LET DISP(Z) = - DB(Z) * (S(Z)-S(Z+1))/DZ^2
			DISP[z0] = -DB[z0] * (Sarray[z1] - Sarray[z1+1]) / sqrdDZ

		} else if z1 < g.N {
			//LET DISP(Z) = DB(Z-1)*(S(Z-1)-S(Z))/DZ^2-DB(Z)*(S(Z)-S(Z+1))/DZ^2
			DISP[z0] = DB[z0-1]*(Sarray[z1-1]-Sarray[z1])/sqrdDZ - DB[z0]*(Sarray[z1]-Sarray[z1+1])/sqrdDZ
		} else {
			//LET DISP(Z) = DB(Z-1)*(S(Z-1)-S(Z))/DZ^2
			DISP[z0] = DB[z0-1] * (Sarray[z1-1] - Sarray[z1]) / sqrdDZ
		}
	}

	// 	! Berechnung der Konvektion für unterschiedliche Fließrichtungsfälle
	// TODO: ask about the Drainage as in Nitro
	// 	FOR Z = 1 TO N
	for z1 := 1; z1 <= g.N; z1++ {
		z0 := z1 - 1
		//IF Q1(Z) >= 0 AND Q1(Z-1) >= 0 then
		if g.Q1[z1] >= 0 && g.Q1[z1-1] >= 0 {
			//LET KONV(Z) = (S(Z)*Q1(Z) - S(Z-1)*Q1(Z-1))/dz
			KONV[z0] = (Sarray[z1]*g.Q1[z1] - Sarray[z1-1]*g.Q1[z1-1]) / g.DZ.Num
			//ELSE IF Q1(Z) >= 0 AND Q1(Z-1) < 0 then
		} else if g.Q1[z1] >= 0 && g.Q1[z1-1] < 0 {
			//IF Z > 1 then
			if z1 > 1 {
				//LET KONV(Z) = (S(Z)*Q1(Z) - S(Z)*Q1(Z-1))/dz
				KONV[z0] = (Sarray[z1]*g.Q1[z1] - Sarray[z1]*g.Q1[z1-1]) / g.DZ.Num
			} else {
				//LET KONV(Z) = S(Z)*Q1(Z)/dz
				KONV[z0] = Sarray[z1] * g.Q1[z1] / g.DZ.Num
			}
			//ELSE IF Q1(Z) < 0 AND Q1(Z-1) < 0 then
		} else if g.Q1[z1] < 0 && g.Q1[z1-1] < 0 {
			//IF Z > 1 then
			if z1 > 1 {
				//LET KONV(Z) = (S(Z+1)*Q1(Z) - S(Z)*Q1(Z-1))/dz
				KONV[z0] = (Sarray[z1+1]*g.Q1[z1] - Sarray[z1]*g.Q1[z1-1]) / g.DZ.Num
			} else {
				//LET KONV(Z) = S(Z+1)*Q1(Z)/dz
				KONV[z0] = Sarray[z1+1] * g.Q1[z1] / g.DZ.Num
			}
			//ELSE IF Q1(Z) < 0 AND Q1(Z-1) >= 0 then
		} else if g.Q1[z1] < 0 && g.Q1[z1-1] >= 0 {
			//LET KONV(Z) = (S(Z+1)*Q1(Z) - S(Z-1)*Q1(Z-1))/dz
			KONV[z0] = (Sarray[z1+1]*g.Q1[z1] - Sarray[z1-1]*g.Q1[z1-1]) / g.DZ.Num
		}
	}
	// 	! Neuberechnung der Sulfatverteilung nach Transport
	// 	!        einschliesslich Umrechnung in kg/ha
	// 	FOR Z = 1 TO N
	for z0 := 0; z0 < g.N; z0++ {
		z1 := z0 + 1
		//LET S1(Z) = (S(Z)*WG(0,Z) + DISP(Z) - KONV(Z))*DZ*100 + SF(Z)
		g.S1[z0] = (Sarray[z1]*g.WG[0][z0]+DISP[z0]-KONV[z0])*g.DZ.Num*100 + g.SF[z0]
		// 	NEXT Z
	}
	// 	! Auswaschungsberechnung
	// 	IF Q1(outn) > 0 then
	if g.Q1[g.OUTN] > 0 {
		// If Outn < n then
		if g.OUTN < g.N {
			//LET SOUTSUM = SOUTSUM + Q1(outn)*S(outn)/dz*100*DZ + DB(outn) * (S(outn)-S(outn+1))/DZ^2*100*dz
			g.SOUTSUM = g.SOUTSUM + g.Q1[g.OUTN]*Sarray[g.OUTN]/g.DZ.Num*100 + DB[g.OUTN]*(Sarray[g.OUTN]-Sarray[g.OUTN+1])/sqrdDZ*100
		} else {
			//LET SOUTSUM = SOUTSUM + Q1(outn)*S(outn)/dz*100*DZ
			g.SOUTSUM = g.SOUTSUM + g.Q1[g.OUTN]*Sarray[g.OUTN]/g.DZ.Num*100*g.DZ.Num
		}
	} else {
		// If Outn < n then
		if g.OUTN < g.N {
			//LET SOUTSUM = Soutsum + (Q1(outn)*S(outn+1)/dz*100*DZ) + (DB(outn) * (S(outn)-S(outn+1))/DZ^2*100*dz)
			g.SOUTSUM = g.SOUTSUM + (g.Q1[g.OUTN] * Sarray[g.OUTN+1] / g.DZ.Num * 100 * g.DZ.Num) + (DB[g.OUTN] * (Sarray[g.OUTN] - Sarray[g.OUTN+1]) / sqrdDZ * 100 * g.DZ.Num)
		}
	}
	// out of last layer
	if g.Q1[g.N] > 0 {
		g.SoutLastLayer = g.SoutLastLayer + g.Q1[g.N]*Sarray[g.N]/g.DZ.Num*100*g.DZ.Num
	}

	// 	! 2. Haelfte Quellen und Senken nach Verlagerung
	// 	FOR Z = 1 TO N
	for z0 := 0; z0 < g.N; z0++ {

		if subd == 1 {
			//IF PES(Z)/2 > S1(Z)-SF(Z)-.02 THEN
			if g.PES[z0]/2 > g.S1[z0]-g.SF[z0]-0.02 {
				//LET PES(Z) = (S1(Z)-SF(Z)-.02) * 2
				g.PES[z0] = (g.S1[z0] - g.SF[z0] - 0.02) * 2
			}
			//IF PES(Z) < 0 THEN LET PES(Z) = 0
			if g.PES[z0] < 0 {
				g.PES[z0] = 0
			}
			//LET PESUMS = PESUMS + PES(Z)/2
			g.PESUMS = g.PESUMS + g.PES[z0]/2
			//LET S1(Z) = S1(Z) + DNS(Z)/2 - PES(Z)/2
			g.S1[z0] = g.S1[z0] + g.DNS[z0]*wdt/2 - g.PES[z0]/2
		} else {
			//LET S1(Z) = S1(Z) + DNS(Z)/2
			g.S1[z0] = g.S1[z0] + g.DNS[z0]*wdt/2
		}

		//IF S1(Z) < 0 THEN
		if g.S1[z0] < 0 {
			//LET S1(Z) = 0.0001
			g.S1[z0] = 0.0001
		}
	}
}

// SUB SRESID(SDI,SSA,SLA)
func sResid(g *GlobalVarsMain, hPath *HFilePath) (SSA, SLA, SDI float64) {
	// 		!Mineralisationspotentiale aus Vorfruchtresiduen
	CRONAM := hPath.cropn
	_, scanner, _ := g.Session.Open(&FileDescriptior{FilePath: CRONAM, UseFilePool: true})
	var SERNT, SWURA, SFAST float64
	for scanner.Scan() {
		CROP := scanner.Text()
		if g.ToCropType(CROP[0:3]) == g.FRUCHT[g.AKF.Index] {
			//! Korn-Stroh Verhältnis
			//LET KOSTRO = VAL(CROP$(5:7))
			//KOSTRO = ValAsFloat(CROP[4:7], CRONAM, CROP)
			//LET TM     = VAL(CROP$(9:12))
			//TM = ValAsFloat(CROP[8:12], CRONAM, CROP) //TM_
			//LET SERNT  = VAL(CROP$(20:24))     ! N = (14:18))
			SERNT = ValAsFloat(CROP[19:24], CRONAM, CROP) //S_HEG
			//LET SKOPP  = VAL(CROP$(32:35))     ! N = (26:30))
			//SKOPP = ValAsFloat(CROP[31:35], CRONAM, CROP) // SNEG
			//LET SWURA  = VAL(CROP$(37:40))
			//Wurzelanteil an Gesamt-S in Pflanze
			SWURA = ValAsFloat(CROP[36:40], CRONAM, CROP) // SWur
			//LET SFAST  = VAL(CROP$(47:50))     ! N = (42:45))
			//Schnell mineralisierbarer Anteil von S in Ernterückständen (Fraktion)
			SFAST = ValAsFloat(CROP[46:50], CRONAM, CROP) // Sfas

			break
		}
	}
	var DGM float64
	// 	IF JN(AKF) = 0 THEN
	if g.JN[g.AKF.Index] == 0 {
		// 	IF EINT(AKF) = 0 then
		if g.EINT[g.NTIL.Index] == 0 {
			// 	   ! LET DGM = 0
			// 	   LET DGM = PESUMS * SWURA
			DGM = g.PESUMS * SWURA
		} else {
			// 	   LET DGM = PESUMS - (ERTR(AKF) * SERNT)
			DGM = g.PESUMS - (g.ERTR[g.AKF.Index] * SERNT)
		}
		//  ELSE IF JN(AKF) = 1  then
	} else if g.JN[g.AKF.Index] == 1 {
		// 	LET DGM = PESUMS * SWURA
		DGM = g.PESUMS * SWURA
		//  ELSE
	} else {
		// 	LET DGM = PESUMS *SWURA + (1-JN(AKF)) * (PESUMS - ERTR(AKF) * SERNT - PESUMS*SWURA)
		DGM = g.PESUMS*SWURA + (1-g.JN[g.AKF.Index])*(g.PESUMS-g.ERTR[g.AKF.Index]*SERNT-g.PESUMS*SWURA)
		// 	!LET DGM = PESUMS *SWURA + JN(AKF) * (PESUMS - ERTR(AKF) * TM * SERNT - PESUMS*SWURA)
		//  END IF
	}
	//  IF DGM < 0 then LET DGM = 0
	if DGM < 0 {
		DGM = 0
	}
	//  LET SSA = DGM * SFAST
	SSA = DGM * SFAST
	//  LET SLA = DGM * (1-SFAST)
	SLA = DGM * (1 - SFAST)
	//  LET SDI = 0.0
	SDI = 0.0

	return SSA, SLA, SDI

}

// SUB SRESIDI
func sResidi(g *GlobalVarsMain, hPath *HFilePath) {
	// ! ******  Mineralisationspotentiale aus Vorfruchtresiduen
	// CROP_S.TXT
	CRONAM := hPath.cropn
	_, scanner, _ := g.Session.Open(&FileDescriptior{FilePath: CRONAM, UseFilePool: true})
	var KOSTRO, SKOPP, SERNT, SWURA, SFAST float64
	for scanner.Scan() {
		CROP := scanner.Text()
		//    IF CROP$(1:3) = FRUCHT$(1) then
		if g.ToCropType(CROP[0:3]) == g.FRUCHT[g.AKF.Index] {
			//       LET KOSTRO = VAL(CROP$(5:7))
			KOSTRO = ValAsFloat(CROP[4:7], CRONAM, CROP)
			//       LET TM     = VAL(CROP$(9:12))
			//       LET SERNT  = VAL(CROP$(20:24))     ! N = (14:18))
			SERNT = ValAsFloat(CROP[19:24], CRONAM, CROP) //S_HEG
			//       LET SKOPP  = VAL(CROP$(32:35))     ! N = (26:30))
			SKOPP = ValAsFloat(CROP[31:35], CRONAM, CROP) // SNEG
			//       LET SWURA  = VAL(CROP$(37:40))
			SWURA = ValAsFloat(CROP[36:40], CRONAM, CROP) // SWur
			//       LET SFAST  = VAL(CROP$(47:50))     ! N = (42:45))
			SFAST = ValAsFloat(CROP[46:50], CRONAM, CROP) // Sfas
			break
		}
	}

	AUFGES := (g.ERTR[0]*SERNT + g.ERTR[0]*KOSTRO*SKOPP) / (1 - SWURA)
	var DGM float64
	if g.JN[0] == 0 {
		if g.EINT[0] == 0 {
			DGM = 0
		} else {
			DGM = AUFGES - (g.ERTR[0] * SERNT)
		}
	} else if g.JN[0] == 1 {
		DGM = AUFGES * SWURA
	} else {
		DGM = AUFGES*SWURA + (1-g.JN[0])*(AUFGES-g.ERTR[0]*SERNT-AUFGES*SWURA)
	}
	if DGM < 0 {
		DGM = 0
	}
	g.SSAS[0] = DGM * SFAST
	g.SLAS[0] = DGM * (1 - SFAST)
	g.SDIR[0] = 0.0
}

// read crop data, calculate S/N Ratio
func sReadCropData(g *GlobalVarsMain, hpath *HFilePath) error {

	// Notes : https://d-nb.info/1175826456/34
	// Schwefel in Ernte/Ernteresten Mengenverhältnis zum Vergleichen
	// WW 0,22 :0,18 -> 1:0,8
	// WRA 0,50 : 0,40  -> 1:1,6
	// K 0,04 : 0,06 -> 1:0,5
	// ZR 0,03 : 0,03 -> 1:0,7

	// only read this file if sulfonie is enabled
	if !g.Sulfonie {
		return nil
	}
	cData := hpath.cropdata
	// check if file exists
	_, scannerCropDataFile, err := g.Session.Open(&FileDescriptior{FilePath: cData, FileDescription: "crop data file", UseFilePool: true})
	if err != nil {
		return err
	}

	header := LineInut(scannerCropDataFile) // skip header
	headerTokens := strings.Fields(header)

	cDHeader := map[string]int{
		"Crop":     -1,
		"HEGzuNEG": -1,
		"N_HEG":    -1,
		"S_HEG":    -1,
		"N_NEG":    -1,
		"S_NEG":    -1,
	}
	for i, token := range headerTokens {
		if _, ok := cDHeader[token]; ok {
			cDHeader[token] = i
		}
	}
	g.HEGzuNEG = make(map[CropType]float64)
	g.N_HEG = make(map[CropType]float64)
	g.S_HEG = make(map[CropType]float64)
	g.N_NEG = make(map[CropType]float64)
	g.S_NEG = make(map[CropType]float64)
	g.SNRatio = make(map[CropType]float64)

	for scannerCropDataFile.Scan() {
		line := scannerCropDataFile.Text()
		token := strings.Fields(line)

		if len(token) >= len(cDHeader) {
			crop := token[cDHeader["Crop"]]
			cropt := g.ToCropType(crop)
			g.S_HEG[cropt] = ValAsFloat(token[cDHeader["S_HEG"]], cData, line)
			g.S_NEG[cropt] = ValAsFloat(token[cDHeader["S_NEG"]], cData, line)
			g.N_HEG[cropt] = ValAsFloat(token[cDHeader["N_HEG"]], cData, line)
			g.N_NEG[cropt] = ValAsFloat(token[cDHeader["N_NEG"]], cData, line)
			g.HEGzuNEG[cropt] = ValAsFloat(token[cDHeader["HEGzuNEG"]], cData, line)

			// calculate SNRatio
			// SNRatio = (HEGzuNEG * N_HEG + N_NEG) / (HEGzuNEG * S_HEG + S_NEG)
			if g.HEGzuNEG[cropt] > 0 {
				g.SNRatio[cropt] = (g.HEGzuNEG[cropt]*g.N_HEG[cropt] + g.N_NEG[cropt]) / (g.HEGzuNEG[cropt]*g.S_HEG[cropt] + g.S_NEG[cropt])
			} else {
				g.SNRatio[cropt] = (g.N_HEG[cropt]) / (g.S_HEG[cropt])
			}
			// check if SNRatio is valid float
			if math.IsInf(g.SNRatio[cropt], 0) {
				return fmt.Errorf("SNRatio is not a number for crop %s", crop)
			}
		}
	}
	return nil
}

// read Smin measurements
func readSmin(g *GlobalVarsMain, Fident string, hPath *HFilePath) {
	// only read this file if sulfonie is enabled
	if g.Sulfonie {
		_, scannerSminFile, _ := g.Session.Open(&FileDescriptior{FilePath: hPath.smin, FileDescription: "smin file", UseFilePool: true})

		g.SI = make(map[int][]float64)

		isHeader := true
		headerMap := map[string]int{"Plot_ID": -1,
			"DATE":    -1,
			"Smi03":   -1,
			"Smi3-6":  -1,
			"Smi6-9":  -1,
			"Smi912":  -1,
			"Smi1215": -1,
			"Smi1520": -1}

		for scannerSminFile.Scan() {
			line := scannerSminFile.Text()
			if isHeader {
				isHeader = false

				// check which columns are available
				//Plot_ID   DATE    Smi03  Smi3-6  Smi6-9  Smi912  Smi1215 Smi1520
				token := Explode(line, []rune{' ', ',', ';'})
				for i, t := range token {
					// check if the header is in the possibleHeader
					if _, ok := headerMap[t]; ok {
						headerMap[t] = i
					}
				}

				continue
			}

			if strings.HasPrefix(line, Fident) {
				token := Explode(line, []rune{' ', ',', ';'})
				if token[0] == Fident {
					date := 0
					if token[1] != "nan" {
						_, date = g.Datum(token[1])
					}
					siValues := make([]float64, g.N)
					Smi0_3 := 0.3
					if headerMap["Smi03"] != -1 {
						Smi0_3 = ValAsFloat(token[headerMap["Smi03"]], hPath.smin, line)
					}
					Smi3_6 := 0.3
					if headerMap["Smi3-6"] != -1 {
						Smi3_6 = ValAsFloat(token[headerMap["Smi3-6"]], hPath.smin, line)
					}
					Smi6_9 := 0.3
					if headerMap["Smi6-9"] != -1 {
						Smi6_9 = ValAsFloat(token[4], hPath.smin, line)
					}
					Smi9_12 := 0.3
					if headerMap["Smi912"] != -1 {
						Smi9_12 = ValAsFloat(token[5], hPath.smin, line)
					}
					Smi12_15 := 0.3
					if headerMap["Smi1215"] != -1 {
						Smi12_15 = ValAsFloat(token[6], hPath.smin, line)
					}
					Smi15_20 := 0.5
					if headerMap["Smi1520"] != -1 {
						Smi15_20 = ValAsFloat(token[7], hPath.smin, line)
					}

					for i := 0; i < g.N; i++ {
						var val float64
						if i < 3 {
							val = Smi0_3
						} else if i < 6 {
							val = Smi3_6
						} else if i < 9 {
							val = Smi6_9
						} else if i < 12 {
							val = Smi9_12
						} else if i < 15 {
							val = Smi12_15
						} else {
							val = Smi15_20
						}
						if val > 0 {
							if i > 16 {
								val = val / 5
							} else {
								val = val / 3
							}
						}

						siValues[i] = val
					}
					g.SI[date] = siValues
					g.sMESS = append(g.sMESS, date)
				}
			}
		}
	}
}
