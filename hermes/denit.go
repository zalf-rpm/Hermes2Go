package hermes

import "math"

//Denitr Denitrification
func Denitr(g *GlobalVarsMain, thetasatFromPorges bool) {
	// !   /*  Denitrifikation Teilmodell                                    */
	// !   /*  Nach U. Schneider, Diss. 1991, und Richter & Sndgerat        */
	// !   /*  Jan. 1992                                                     */
	// ! Inputs:
	// ! WG(1,I) 		= Wassergehalt Schicht I (cm^3/cm^3)
	// ! PORGES(I)		= Gesamtporenvolumen (cm^3/cm^3)
	// ! C1(I)		= Nitratgehalt Schicht I (kg N/ha)
	// ! Tsoil(0,I)	= Bodentemperatur in Schicht I (°C)
	// !
	thetaOb30 := (g.WG[1][0] + g.WG[1][1] + g.WG[1][2]) / 3
	var thetasat float64
	if thetasatFromPorges {
		thetasat = (g.PORGES[0] + g.PORGES[1] + g.PORGES[2]) / 3
	} else {
		thetasat = 1 - (1.45 / 2.65)
	}
	thetarel := thetaOb30 / thetasat
	nitratOb30 := g.C1[0] + g.C1[1] + g.C1[2]
	if nitratOb30 > 0 {
		layerFraction := []float64{
			g.C1[0] / nitratOb30,
			g.C1[1] / nitratOb30,
			g.C1[2] / nitratOb30}

		tempOb30 := (g.TSOIL[0][0] + g.TSOIL[0][1] + g.TSOIL[0][2] + g.TSOIL[0][3]) / 4
		if tempOb30 < 0 {
			tempOb30 = 0
		}
		// !
		// !   /*  Schätzwerte von U. Schneider, Diss. '91. S.57, 0-30 cm        */
		// !   /*  Vmax   =  1274    (g/ha/tag)                                  */
		// !   /*  KNO3   =  74      (kg/ha/30 cm Tiefe)                         */
		// !   /*  Tkrt   =  15.5    (degrees Celsius)                           */
		// !   /*  Okrt   =  0.766   (relatieve volumetrische Wassergehalt)      */
		Vmax := 1274.
		KNO3 := 74.
		Tkrt := 15.5
		Okrt := 0.766
		Nquadrat := math.Pow(nitratOb30, 2)
		michment := (Vmax * Nquadrat) / (Nquadrat + KNO3)
		Ftheta := 1 - math.Exp(-1*math.Pow((thetarel/Okrt), 6))
		Ftemp := 1 - math.Exp(-1*math.Pow((tempOb30/Tkrt), 4.6))
		DENIT := michment * Ftheta * Ftemp
		DENIT = DENIT / 1000
		// !   /* Denitrifizierte N von Nitrat Pool wegnehmen                    */
		for z := 0; z < 3; z++ {
			if layerFraction[z] > 0 {
				g.C1[z] = g.C1[z] - DENIT*layerFraction[z]
				if g.C1[z] < 0 {
					g.C1[z] = 0
				}
			}
		}

		//Let MaxN2O = 0.63
		MaxN2O := 0.63
		//LET FO = 1 - 2.05 * Max(0,Thetarel-0.62)
		FO := 1 - 2.05*math.Max(0, thetarel-0.62)
		//Let DNO = (0.44 + 0.0015*3)/3
		DNO := (0.44 + 0.0015*3) / 3
		//Let FN = Min(DNO*nitratOb30*0.667,(0.44+0.0015 * 0.67*nitratOB30))
		FN := math.Min(DNO*nitratOb30*0.667, (0.44 + 0.0015*0.67*nitratOb30))
		//IF FN > 1 then Let FN = 1
		if FN > 1 {
			FN = 1
		}
		//Let FN2Oden = FN*FO * MaxN2O
		FN2Oden := FN * FO * MaxN2O
		//Let N2Oden = Denit * FN2Oden
		N2Oden := DENIT * FN2Oden
		//Let N2Odencum = N2Odencum + N2Oden
		g.N2Odencum = g.N2Odencum + N2Oden

		g.CUMDENIT = g.CUMDENIT + DENIT
	}

}

//Denitmo Denitrification for marsh land soils
func Denitmo(g *GlobalVarsMain) {
	// ! Vorläufige Erweiterung des Denitrifikationsansatzes für Moorböden
	// !   /*  Denitrifikation Teilmodell                                    */
	// !   /*  Nach U. Schneider, Diss. 1991, und Richter & Sndgerat        */
	// !   /*  Jan. 1992                                                     */
	thetaOb30 := (g.WG[1][0] + g.WG[1][1] + g.WG[1][2]) / 3
	thetaOb60 := (g.WG[1][3] + g.WG[1][4] + g.WG[1][5]) / 3
	thetaOb90 := (g.WG[1][6] + g.WG[1][7] + g.WG[1][8]) / 3
	thetasat1 := (g.PORGES[0] + g.PORGES[1] + g.PORGES[2]) / 3
	thetasat2 := (g.PORGES[3] + g.PORGES[4] + g.PORGES[5]) / 3
	thetasat3 := (g.PORGES[6] + g.PORGES[7] + g.PORGES[8]) / 3
	thetarel1 := thetaOb30 / thetasat1
	thetarel2 := thetaOb60 / thetasat2
	thetarel3 := thetaOb90 / thetasat3
	nitratOb30 := g.C1[0] + g.C1[1] + g.C1[2]
	nitratOb60 := g.C1[3] + g.C1[4] + g.C1[5]
	nitratOb90 := g.C1[6] + g.C1[7] + g.C1[8]
	var layerFraction30 [3]float64
	var layerFraction60 [3]float64
	var layerFraction90 [3]float64

	if nitratOb30 > 0 {
		layerFraction30[0] = g.C1[0] / nitratOb30
		layerFraction30[1] = g.C1[1] / nitratOb30
		layerFraction30[2] = g.C1[2] / nitratOb30
	}
	if nitratOb60 > 0 {
		layerFraction60[0] = g.C1[3] / nitratOb60
		layerFraction60[1] = g.C1[4] / nitratOb60
		layerFraction60[2] = g.C1[5] / nitratOb60
	}
	if nitratOb90 > 0 {
		layerFraction90[0] = g.C1[6] / nitratOb90
		layerFraction90[1] = g.C1[8] / nitratOb90
		layerFraction90[2] = g.C1[7] / nitratOb90
	}

	tempOb30 := g.TEMP[g.TAG.Index]
	if tempOb30 < 0 {
		tempOb30 = 0
	}
	tempOb60 := tempOb30
	// geschätzte Jahresmitteltemperatur
	tempOb90 := 8.
	// !   /*  Schätzwerte von U. Schneider, Diss. '91. S.57, 0-30 cm        */
	// !   /*  Vmax   =  1274 * CGehalt/3.75 (g/ha/tag)                      */
	// !   Von Neuenkirchen korrigiert um 1/3 (Dichte) u. 1/1.5 (10/15)
	// !   /*  KNO3   =  74      (kg/ha/30 cm Tiefe)                         */
	// !   /*  Tkrt   =  15.5    (degrees Celsius)                           */
	// !   /*  Okrt   =  0.766   (relativer volumetrischer Wassergehalt)     */
	Vmax := 4242. //  geändert 21.9.93  von Faktor 5
	KNO3 := 74.
	Tkrt := 15.5
	Okrt := 0.766
	var Denit1, Denit2, Denit3 float64
	if nitratOb30 > 0 {
		Nquadrat1 := math.Pow(nitratOb30, 2)
		michment1 := (Vmax * Nquadrat1) / (Nquadrat1 + KNO3)
		Ftheta1 := 1 - math.Exp(-1*math.Pow((thetarel1/Okrt), 6))
		Ftemp1 := 1 - math.Exp(-1*math.Pow((tempOb30/Tkrt), 4.6))
		Denit1 = michment1 * Ftheta1 * Ftemp1
		Denit1 = Denit1 / 1000 // (kg/ha)
	}
	if nitratOb60 > 0 {
		Nquadrat2 := math.Pow(nitratOb60, 2)
		michment2 := (Vmax * Nquadrat2) / (Nquadrat2 + KNO3)
		Ftheta2 := 1 - math.Exp(-1*math.Pow((thetarel2/Okrt), 6))
		Ftemp2 := 1 - math.Exp(-1*math.Pow((tempOb60/Tkrt), 4.6))
		Denit2 = michment2 * Ftheta2 * Ftemp2
		Denit2 = Denit2 / 1000 // (kg/ha)
	}
	if nitratOb90 > 0 {
		Nquadrat3 := math.Pow(nitratOb90, 2)
		michment3 := (Vmax * Nquadrat3) / (Nquadrat3 + KNO3)
		Ftheta3 := 1 - math.Exp(-1*math.Pow((thetarel3/Okrt), 6))
		Ftemp3 := 1 - math.Exp(-1*math.Pow((tempOb90/Tkrt), 4.6))
		Denit3 = michment3 * Ftheta3 * Ftemp3
		Denit3 = Denit3 / 1000 // (kg/ha)
	}
	//! new for N2O from denitrification ! acc. to Bessou et al. 2010
	MaxN2O := 0.63
	FO1 := 1 - 2.05*math.Max(0, thetarel1-0.62)
	FO2 := 1 - 2.05*math.Max(0, thetarel2-0.62)
	FO3 := 1 - 2.05*math.Max(0, thetarel3-0.62)
	DNO := (0.44 + 0.0015*3) / 3
	FN1 := math.Min(DNO*nitratOb30*0.667, (0.44 + 0.0015*0.67*nitratOb30))
	FN2 := math.Min(DNO*nitratOb60*0.667, (0.44 + 0.0015*0.67*nitratOb60))
	FN3 := math.Min(DNO*nitratOb90*0.667, (0.44 + 0.0015*0.67*nitratOb90))
	if FN1 > 1 {
		FN1 = 1
	}
	if FN2 > 1 {
		FN2 = 1
	}
	if FN3 > 1 {
		FN3 = 1
	}
	FN2Oden1 := FN1 * FO1 * MaxN2O
	FN2Oden2 := FN2 * FO2 * MaxN2O
	FN2Oden3 := FN3 * FO3 * MaxN2O
	N2Oden := Denit1*FN2Oden1 + Denit2*FN2Oden2 + Denit3*FN2Oden3
	g.N2Odencum = g.N2Odencum + N2Oden

	// !   /* Denitrifizierte N von Nitrat Pool wegnehmen                    */
	calcDenitLayer := func(c1Val *float64, faction float64, denit float64) {
		if faction > 0 {
			*c1Val = *c1Val - denit*faction
			if *c1Val < 0 {
				*c1Val = 0
			}
		}
	}

	calcDenitLayer(&g.C1[0], layerFraction30[0], Denit1)
	calcDenitLayer(&g.C1[1], layerFraction30[1], Denit1)
	calcDenitLayer(&g.C1[2], layerFraction30[2], Denit1)
	calcDenitLayer(&g.C1[3], layerFraction60[0], Denit2)
	calcDenitLayer(&g.C1[4], layerFraction60[1], Denit2)
	calcDenitLayer(&g.C1[5], layerFraction60[2], Denit2)
	calcDenitLayer(&g.C1[6], layerFraction90[0], Denit3)
	calcDenitLayer(&g.C1[7], layerFraction90[1], Denit3)
	calcDenitLayer(&g.C1[8], layerFraction90[2], Denit3)

	g.CUMDENIT = g.CUMDENIT + Denit1 + Denit2 + Denit3
}
