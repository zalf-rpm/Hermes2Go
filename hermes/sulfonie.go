package hermes

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
func Sulfo(subd, zeit int, g *GlobalVarsMain) {

	if zeit == g.MESS[0] {
		for z := 0; z < g.N; z++ {
			g.S1[z] = g.SI[1][z]
		}
	}
	if !g.AUTOFERT {
		//! +++++++++++++++++++++++++++++++++++++ Option real fertilization +++++++++++++++++++++++++++++++++++++++++++++++
		// IF ZEIT = ZTDG(NDG)+1 THEN
		if zeit == g.ZTDG[g.SDG.Index]+1 && subd == 1 {
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
	if subd == 1 {
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
	}

	if subd == 1 {
		sMineral(g)
	}

	//IF ZEIT = ERNTE(AKF) THEN
	if zeit == g.ERNTE[g.AKF.Index] && subd == 1 {
		var SSA, SLA, SDI float64
		//IF akf <> 1 then
		if g.AKF.Num != 1 {
			//CALL SRESID(SDI,SSA,SLA)
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
	}
	// 		! ---- Homogene Vermischung von mineralischem u. organ. N bei Bearbeitung ----
	// 		IF ZEIT = EINTE(AKF-1) THEN
	if zeit == g.EINTE[g.AKF.Index] && subd == 1 {
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
	sMove(subd, zeit, g)
}

// 	SUB SMINERAL
// 		DIM DSAOS(4),DSFOS(4),MIRED(4)
// 		REM ---------------------  Mineralisation  --------------------
// 		FOR z = 1 to izm/dz
// 			IF Z = 1 THEN
// 			   IF TEMPBO1(TAG) <> 0 THEN
// 				  LET TEMPBO = TEMPBO1(TAG)
// 			   ELSE IF TEMPBO2(TAG) <> 0 THEN
// 				  LET TEMPBO = (TEMPBO2(TAG)+TEMP(TAG))/2
// 			   ELSE
// 				  LET TEMPBO = TEMP(TAG)
// 			   END IF
// 			ELSE
// 			   IF TEMPBO2(TAG) <> 0 THEN
// 				  LET TEMPBO = TEMPBO2(TAG)
// 			   ELSE IF TEMPBO1(TAG) <> 0 THEN
// 				  LET TEMPBO = 2*TEMPBO1(TAG)-TEMP(TAG)
// 			   ELSE
// 				  LET TEMPBO = TEMP(TAG)
// 			   END IF
// 			END IF
// 			REM --------- Berechnung Mineralisationskoeffizienten ---------
// 			REM ----------- in Abh�ngigkeit von TEMP UND WASSER -----------
// 			REM - Umsetzung von mineralischen D�ngern
// 			LET KTD = .4
// 			! Reduktion bei suboptimalem Wassergehalt
// 			IF WG(0,z) <= WNOR(Z) AND WG(0,Z) >= WRED THEN
// 			   LET MIRED(Z) = 1
// 			ELSE IF WG(0,Z) < WRED AND WG(0,Z) > WMIN(Z) THEN
// 			   LET MIRED(Z) = (WG(0,Z) - WMIN(z)) / (WRED - WMIN(z))
// 			ELSE IF WG(0,Z) > WNOR(Z) THEN
// 			   LET MIRED(Z) = (PORGES(Z)-WG(0,Z)) / (PORGES(Z) - WNOR(Z))
// 			ELSE
// 			   LET MIRED(Z) =0
// 			END IF
// 			IF MIRED(Z) < 0 THEN LET MIRED(Z) = 0
// 			IF MIRED(Z) > 1 THEN LET MIRED(Z) = 1
// 			! Temperaturabh�ngigkeit der Mineralisationskoeffizienten nur >0 Grad
// 			IF TEMPBO > 0 THEN
// 			   !
// 			   REM - Reaktionskoeffizient der schwer abbaubaren Fraktion
// 			   !
// 			   LET kt0 = 4000000000*exp(-8400/(TEMPBO+273.16))*dt
// 			   !
// 			   REM - Reaktionskoeffizient der leicht abbaubaren Fraktion
// 			   !
// 			   LET kt1 = 5.6e+12*exp(-9800/(TEMPBO+273.16))*dt
// 			   !
// 			ELSE
// 			   LET kt0,kt1 = 0
// 			END IF
// 			! - Mineralisation der schwer abbaubaren Fraktion
// 			LET DSAOS(z) = kt0*SAOS(z)*MIRED(Z)
// 			! Reduktion des Pools um die mineralisierte Menge
// 			LET SAOS(z) = SAOS(z) - DSAOS(z)
// 			!
// 			REM - Mineralisation der leicht abbaubaren Fraktion
// 			!
// 			LET DSFOS(z) = kt1*SFOS(z)*MIRED(Z)
// 			! Reduktion des Pools um die mineralisierte Menge
// 			LET SFOS(z) = SFOS(z) - DSFOS(z)
// 			!
// 			IF z = 1 THEN
// 			   LET DUMS(Z) = KTD * MIRED(Z) * DSUMM *dt
// 			   LET DSUMM = DSUMM - DUMS(z)
// 			ELSE
// 			   LET DUMS(Z) = 0
// 			END IF
// 			REM - Mineralisationssumme => Quellterm ( DNS(z) )
// 			!
// 			!print z,dsaos(z),DSFOS(Z),SAOS(z),SFOS(Z)
// 			LET DNS(z) = DSAOS(z)+DSFOS(z)+DUMS(Z)
// 			LET SMINAOS(z) = SMINAOS(z)+DSAOS(z)
// 			LET SMINFOS(z) = SMINFOS(z)+DSFOS(z)
// 			LET SMINSUM = SMINSUM + DNS(Z)-DUMS(Z)
// 		NEXT z
// 		!get key wart
// 	END SUB

// 	SUB SMOVE(#5)
// 		DIM S(0:21)
// 		FOR Z = 1 TO N
// 			! --- Berechnung des Diffusionskoeffizienten am unteren Kompartimentrand ---
// 			LET D(Z) = D0S * (AD*EXP((WG(0,Z)+WG(0,Z+1))*5)/((WG(0,Z)+WG(0,Z+1))/2))*DT
// 			! ** Loeslichkeitsobergrenzen und Loesungs-/Faellungsreaktion **
// 			! SKSAT ist die Saettigungs-Loesungskonzentration in Gramm S/Liter
// 			! SF(z) ist die nicht gel�ste Smin-Menge in kg S/ha
// 			LET S(Z) = (S1(Z)-SF(Z))/(wg(0,z)*DZ*100)
// 			IF S(Z) >= SKSAT then
// 			 LET S(Z) = SKSAT
// 			 LET SF(Z) = S1(Z) - S(Z)*(wg(0,z)*dz*100)
// 			ELSE
// 			 LET S(Z)  = S(Z) + (Sksat-S(Z)) * (1-EXP(-klos*(SF(Z)/(WG(0,z)*dz*100))))
// 			 !LET S(Z) = S(Z) + (SF(Z)/(Wg(0,z)*dz*100)*(SKSAT-S(Z))) !*lok
// 			 LET SF(Z) = S1(Z) - S(Z)*(wg(0,z)*dz*100)
// 			END IF
// 			! Untere Begrenzung f�r Entleerung pro Schicht (0.02 kg S/ha)
// 			IF PES(Z) > S1(Z)-SF(Z)-.02 THEN
// 			   LET PES(Z) = (S1(Z)-SF(Z)-.02)
// 			END IF
// 			IF PES(Z) < 0 THEN LET PES(Z) = 0
// 			LET PESUMS = PESUMS + PES(Z)/2
// 			LET AUFNASUM = AUFNASUM + PES(Z)/2
// 			! Umrechnung in Bodenloesungskonzentration (kg/ha --> g/l)
// 			! Quellen und Senken jeweils halb bei Beginn/Ende Zeitschritt
// 			LET S(Z) = (S1(Z)-SF(Z) + DNS(Z)/2 - PES(Z)/2)/(wg(0,z)*DZ*100)
// 		NEXT Z

// 		!--------------------- Verlagerung nach unten ---------------------
// 		LET Q1(0) = FLUSS0*DT
// 		! Berechnung der Dispersion in Abh�ngigkeit der Porenwassergeschwindigkeit
// 		FOR Z = 1 TO N
// 			!LET V(Z) = ABS(Q1(Z)/((WG(0,Z)+WG(0,Z+1))*.5))
// 			LET V(Z) = ABS(Q1(Z)/((W(Z)+W(Z+1))*.5))
// 			LET DB(Z) = (WG(0,Z)+WG(0,Z+1))/2*(D(Z) + DV*V(Z))-.5*dz*ABS(q1(z))+.5*dt*ABS((q1(z)+q1(z-1))/2)*v(z)
// 			IF Z = 1 THEN
// 			   LET DISP(Z) = - DB(Z) * (S(Z)-S(Z+1))/DZ^2
// 			ELSE IF Z < N THEN
// 			   LET DISP(Z) = DB(Z-1)*(S(Z-1)-S(Z))/DZ^2-DB(Z)*(S(Z)-S(Z+1))/DZ^2
// 			ELSE
// 			   LET DISP(Z) = DB(Z-1)*(S(Z-1)-S(Z))/DZ^2
// 			END IF
// 		NEXT Z
// 		! Berechnung der Konvektion f�r unterschiedliche Flie�richtungsfaelle
// 		FOR Z = 1 TO N
// 			IF Q1(Z) >= 0 AND Q1(Z-1) >= 0 then
// 			   LET KONV(Z) = (S(Z)*Q1(Z) - S(Z-1)*Q1(Z-1))/dz
// 			ELSE IF Q1(Z) >= 0 AND Q1(Z-1) < 0 then
// 			   IF Z > 1 then
// 				  LET KONV(Z) = (S(Z)*Q1(Z) - S(Z)*Q1(Z-1))/dz
// 			   ELSE
// 				  LET KONV(Z) = S(Z)*Q1(Z)/dz
// 			   END IF
// 			ELSE IF Q1(Z) < 0 AND Q1(Z-1) < 0 then
// 			   IF Z > 1 then
// 				  LET KONV(Z) = (S(Z+1)*Q1(Z) - S(Z)*Q1(Z-1))/dz
// 			   ELSE
// 				  LET KONV(Z) = S(Z+1)*Q1(Z)/dz
// 			   END IF
// 			ELSE IF Q1(Z) < 0 AND Q1(Z-1) >= 0 then
// 			   LET KONV(Z) = (S(Z+1)*Q1(Z) - S(Z-1)*Q1(Z-1))/dz
// 			END IF
// 		NEXT Z
// 		! Neuberechnung der Sulfatverteilung nach Transport
// 		!        einschliesslich Umrechnung in kg/ha
// 		FOR Z = 1 TO N
// 			LET S1(Z) = (S(Z)*WG(0,Z) + DISP(Z) - KONV(Z))*DZ*100 + SF(Z)
// 		NEXT Z
// 		IF Q1(outn) > 0 then
// If Outn < n then
// 		   LET SOUTSUM = SOUTSUM + Q1(outn)*S(outn)/dz*100*DZ + DB(outn) * (S(outn)-S(outn+1))/DZ^2*100*dz
// else
// 		   LET SOUTSUM = SOUTSUM + Q1(outn)*S(outn)/dz*100*DZ
// end if
// 		ELSE
// If Outn < n then
// 		   LET SOUTSUM = Soutsum + (Q1(outn)*S(outn+1)/dz*100*DZ) + (DB(outn) * (S(outn)-S(outn+1))/DZ^2*100*dz)
// else
// 		   Let SOutsum = Soutsum
// end if
// 		END IF
// 		! 2. Haelfte Quellen und Senken nach Verlagerung
// 		FOR Z = 1 TO N
// 			IF PES(Z)/2 > S1(Z)-SF(Z)-.02 THEN
// 			   LET PES(Z) = (S1(Z)-SF(Z)-.02) * 2
// 			END IF
// 			IF PES(Z) < 0 THEN LET PES(Z) = 0
// 			LET S1(Z) = S1(Z) + DNS(Z)/2 - PES(Z)/2
// 			LET PESUMS = PESUMS + PES(Z)/2
// 			LET AUFNASUM = AUFNASUM + PES(Z)/2
// 			IF S1(Z) < 0 THEN
// 			   LET S1(Z) = 0.0001
// 			END IF
// 		NEXT Z
// 	END SUB

// 	SUB SRESID(SDI,SSA,SLA)
// 		!             Mineralisationspotentiale aus Vorfruchtresiduen
// 		LET CRONAM$ = PATH$ & "CROP_S.TXT"
// 		OPEN #4:Name CROnam$,ACCESS INPUT,ORGANIZATION TEXT
// 		DO while more #4
// 		   LINE INPUT #4: CROP$
// 		   IF CROP$(1:3) = FRUCHT$(AKF) then
// 			  LET KOSTRO = VAL(CROP$(5:7))
// 			  LET TM     = VAL(CROP$(9:12))
// 			  LET SERNT  = VAL(CROP$(20:24))     ! N = (14:18))
// 			  LET SKOPP  = VAL(CROP$(32:35))     ! N = (26:30))
// 			  LET SWURA  = VAL(CROP$(37:40))
// 			  LET SFAST  = VAL(CROP$(47:50))     ! N = (42:45))
// 			  EXIT DO
// 		   END IF
// 		LOOP

// 		IF JN(AKF) = 0 THEN
// 		   IF EINT(AKF) = 0 then
// 			  ! LET DGM = 0
// 			  LET DGM = PESUMS * SWURA
// 		   ELSE
// 			  LET DGM = PESUMS - (ERTR(AKF) * SERNT)
// 		   END IF
// 		   !LET DGM = PESUMS - (ERTR(AKF) * TM * SERNT)
// 		ELSE IF JN(AKF) = 1  then      !"A" or JN$(AKF) = "V" then
// 		   LET DGM = PESUMS * SWURA
// 		   ! END IF
// 		ELSE
// 		   LET DGM = PESUMS *SWURA + (1-JN(AKF)) * (PESUMS - ERTR(AKF) * SERNT - PESUMS*SWURA)
// 		   !LET DGM = PESUMS *SWURA + JN(AKF) * (PESUMS - ERTR(AKF) * TM * SERNT - PESUMS*SWURA)
// 		END IF
// 		IF DGM < 0 then LET DGM = 0
// 		LET SSA = DGM * SFAST
// 		LET SLA = DGM * (1-SFAST)
// 		LET SDI = 0.0
// !print akf,FRUCHT$(AKF),ERTR(AKF),SSA,SLA,PESUMS
// !get key wart
// 	END SUB
