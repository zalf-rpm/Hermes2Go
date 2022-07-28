package hermes

import "math"

func PotMin(g *GlobalVarsMain) {
	//IF CGEHALT(1) > 14 THEN
	if g.CGEHALT[0] > 14 {
		g.SALTOS = 5000 * g.SGEHALT[0] * g.SAKT * float64(g.UKT[0])
		//ELSE IF CGEHALT(1) > 5 THEN
	} else if g.CGEHALT[0] > 5 {
		g.SALTOS = 11000 * g.SGEHALT[0] * g.SAKT * float64(g.UKT[0])
		//ELSE IF CGEHALT(1) < 1 THEN
	} else if g.CGEHALT[0] < 1 {
		g.SALTOS = 15000 * g.SGEHALT[0] * g.SAKT * float64(g.UKT[0])
		//ELSE
	} else {
		g.SALTOS = 15000 * g.SGEHALT[0] * g.SAKT * float64(g.UKT[0])
	}
}

// SUB SULFO() call before nitro()
func Sulfo(wdt float64, subd, zeit int, g *GlobalVarsMain, hPath *HFilePath) {

	if subd == 1 {
		if zeit == g.MESS[0] {
			for z := 0; z < g.N; z++ {
				g.S1[z] = g.SI[1][z]
			}
		}
		if !g.AUTOFERT {
			//! +++++++++++++++++++++++++++++++++++++ Option real fertilization +++++++++++++++++++++++++++++++++++++++++++++++
			// IF ZEIT = ZTDG(NDG)+1 THEN
			if zeit == g.ZTDG[g.SDG.Index]+1 {
				// 		   LET SFOS(1) = SFOS(1) + SSAS(NDG)
				g.SFOS[0] += g.SSAS[g.SDG.Index]
				// 		   LET SAOS(1) = SAOS(1) + SLAS(NDG)
				g.SAOS[0] += g.SLAS[g.SDG.Index]
				// 		   LET DSUMM = DSUMM + SDIR(NDG)    ! Summe miner. Düngung
				g.SDSUMM += g.SDIR[g.SDG.Index]
				// 		   LET NDG = NDG + 1
				g.SDG.Inc()
			}
		}
		// TODO: Auto-fertilization

		// 		IF Zeit = EINTE(1)+1 then
		if zeit == g.EINTE[0]+1 {
			g.SFOSUM, g.SAOSUM, g.SSUM, g.SFSUM = 0, 0, 0, 0
			// 		   IF EINT(1) > 0 then
			if g.EINT[0] > 0 {
				// 		 FOR Z = 1 TO EINT(1)/dz
				for z := 0; z < int(g.EINT[0]/g.DZ.Num); z++ {
					//! Vollständige Durchmischung bis Bearbeitungstiefe für Anfangsforfrucht
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
				for z := 0; z < int(g.EINT[0]/g.DZ.Num); z++ {
					//LET SFOS(Z) = SFOSUM/EINT(1)*dz
					g.SFOS[z] = g.SFOSUM / g.EINT[0] * g.DZ.Num
					//LET SAOS(Z) = SAOSUM/EINT(1)*dz
					g.SAOS[z] = g.SAOSUM / g.EINT[0] * g.DZ.Num
					//LET S1(z)   = SSUM/EINT(1)*dz
					g.S1[z] = g.SSUM / g.EINT[0] * g.DZ.Num
					//LET SF(z)   = SFSUM/EINT(1)*dz
					g.SF[z] = g.SFSUM / g.EINT[0] * g.DZ.Num
				}
			}
		}

		sMineral(g)

		//IF ZEIT = ERNTE(AKF) THEN
		if zeit == g.ERNTE[g.AKF.Index] {
			var SSA, SLA, SDI float64
			//IF akf <> 1 then
			if g.AKF.Num != 1 {
				//CALL SRESID(SSA,SLA)
				sResid(g, hPath)
			}
			//LET SFOS(1) = SFOS(1) + SSA
			g.SFOS[0] += SSA
			//LET SAOS(1) = SAOS(1) + SLA
			g.SAOS[0] += SLA
			//LET DSUMM = DSUMM + SDI
			g.SDSUMM += SDI
			// 		   LET AKF = AKF+1
			// don't increment AKF here, it is used in nitro()
		}
		// 		! ---- Homogene Vermischung von mineralischem u. organ. N bei Bearbeitung ----
		// 		IF ZEIT = EINTE(AKF-1) THEN
		if zeit == g.EINTE[g.AKF.Index] {
			g.SFOSUM, g.SAOSUM, g.SSUM, g.SFSUM = 0, 0, 0, 0
			//IF EINT(AKF-1) > 0 then
			if g.EINT[g.AKF.Index] > 0 {
				//FOR Z = 1 TO EINT(AKF-1)/dz
				for z := 0; z < int(g.EINT[g.AKF.Index]/g.DZ.Num); z++ {
					//LET SFOSUM = SFOSUM + SFOS(Z)
					g.SFOSUM += g.SFOS[z]
					//LET SAOSUM = SAOSUM + SAOS(Z)
					g.SAOSUM += g.SAOS[z]
					//LET SSUM   = SSUM   + S1(Z)
					g.SSUM += g.S1[z]
					//LET SFSUM = SFSUM + SF(Z)
					g.SFSUM += g.SF[z]

				}
				//FOR Z = 1 TO EINT(AKF-1)/dz
				for z := 0; z < int(g.EINT[g.AKF.Index]/g.DZ.Num); z++ {

					//LET SFOS(Z) = SFOSUM/EINT(AKF-1)*dz
					g.SFOS[z] = g.SFOSUM / g.EINT[g.AKF.Index] * g.DZ.Num
					//LET SAOS(Z) = SAOSUM/EINT(AKF-1)*dz
					g.SAOS[z] = g.SAOSUM / g.EINT[g.AKF.Index] * g.DZ.Num
					//LET S1(z)   = SSUM/EINT(AKF-1)*dz
					g.S1[z] = g.SSUM / g.EINT[g.AKF.Index] * g.DZ.Num
					//LET SF(z)   = SFSUM/EINT(AKF-1)*dz
					g.SF[z] = g.SFSUM / g.EINT[g.AKF.Index] * g.DZ.Num
				}
			}
		}
	}
	sMove(wdt, subd, zeit, g)
}

// 	SUB SMINERAL
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

// 	SUB SMOVE(#5)
func sMove(wdt float64, subd, zeit int, g *GlobalVarsMain) {
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
		Sarray[z1] = (g.S1[z] - g.SF[z]) / (g.WG[0][z] * g.DZ.Num * 100)
		//IF S(Z) >= SKSAT then
		if Sarray[z1] >= g.SKSAT {
			//LET S(Z) = SKSAT
			Sarray[z1] = g.SKSAT
			//LET SF(Z) = S1(Z) - S(Z)*(wg(0,z)*dz*100)
			g.SF[z] = g.S1[z] - Sarray[z1]*(g.WG[0][z]*g.DZ.Num*100)
		} else {
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
			//LET AUFNASUM = AUFNASUM + PES(Z)/2
			g.SAUFNASUM += g.PES[z] / 2
		}
		// Umrechnung in Bodenloesungskonzentration (kg/ha --> g/l)
		// Quellen und Senken jeweils halb bei Beginn/Ende Zeitschritt
		// LET S(Z) = (S1(Z)-SF(Z) + DNS(Z)/2 - PES(Z)/2)/(wg(0,z)*DZ*100)
		Sarray[z1] = (g.S1[z] - g.SF[z] + g.DNS[z]*wdt/2 - g.PES[z]/2) / (g.WG[0][z] * g.DZ.Num * 100)

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
		DB[z0] = (g.WG[0][z0]+g.WG[0][z0+1])/2*(DiffCoeff[z0]+g.DV*V[z0]) - .5*g.DZ.Num*math.Abs(g.Q1[z1]) + .5*wdt*math.Abs((g.Q1[z1]+g.Q1[z1-1])/2)*V[z0]
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

	// 	! 2. Haelfte Quellen und Senken nach Verlagerung
	// 	FOR Z = 1 TO N
	for z0 := 0; z0 < g.N; z0++ {

		//IF PES(Z)/2 > S1(Z)-SF(Z)-.02 THEN
		if g.PES[z0]/2 > g.S1[z0]-g.SF[z0]-0.02 {
			//LET PES(Z) = (S1(Z)-SF(Z)-.02) * 2
			g.PES[z0] = (g.S1[z0] - g.SF[z0] - 0.02) * 2
		}
		//IF PES(Z) < 0 THEN LET PES(Z) = 0
		if g.PES[z0] < 0 {
			g.PES[z0] = 0
		}
		//LET S1(Z) = S1(Z) + DNS(Z)/2 - PES(Z)/2
		g.S1[z0] = g.S1[z0] + g.DNS[z0]/2 - g.PES[z0]/2
		//LET PESUMS = PESUMS + PES(Z)/2
		g.PESUMS = g.PESUMS + g.PES[z0]/2
		//LET AUFNASUM = AUFNASUM + PES(Z)/2
		g.SAUFNASUM = g.SAUFNASUM + g.PES[z0]/2
		//IF S1(Z) < 0 THEN
		if g.S1[z0] < 0 {
			//LET S1(Z) = 0.0001
			g.S1[z0] = 0.0001
		}
	}
}

// 	SUB SRESID(SDI,SSA,SLA)
func sResid(g *GlobalVarsMain, hPath *HFilePath) {
	// 		!Mineralisationspotentiale aus Vorfruchtresiduen
	CRONAM := hPath.cropn
	_, scanner, _ := Open(&FileDescriptior{FilePath: CRONAM, UseFilePool: true})
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
	// 		IF JN(AKF) = 0 THEN
	if g.JN[g.AKF.Index] == 0 {
		//IF EINT(AKF) = 0 then
		if g.EINT[g.AKF.Index] == 0 {
			//! LET DGM = 0
			//LET DGM = PESUMS * SWURA
			DGM = g.PESUMS * SWURA
		} else {
			//LET DGM = PESUMS - (ERTR(AKF) * SERNT)
			DGM = g.PESUMS - (g.ERTR[g.AKF.Index] * SERNT)
		}
		//ELSE IF JN(AKF) = 1  then
	} else if g.JN[g.AKF.Index] == 1 {
		//LET DGM = PESUMS * SWURA
		DGM = g.PESUMS * SWURA
	} else {
		//LET DGM = PESUMS *SWURA + (1-JN(AKF)) * (PESUMS - ERTR(AKF) * SERNT - PESUMS*SWURA)
		DGM = g.PESUMS*SWURA + (1-g.JN[g.AKF.Index])*(g.PESUMS-g.ERTR[g.AKF.Index]*SERNT-g.PESUMS*SWURA)
	}
	//IF DGM < 0 then LET DGM = 0
	if DGM < 0 {
		DGM = 0
	}
	//LET SSAS(1) = DGM * SFAST
	g.SSAS[0] = DGM * SFAST
	//LET SLAS(1) = DGM * (1-SFAST)
	g.SLAS[0] = DGM * (1 - SFAST)
	//LET SDIR(1) = 0.0
	g.SDIR[0] = 0.0 // SDIR is set to 0 afer harvest.. what about quick following cultures?
}
