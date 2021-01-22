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
	g.C1[0] = g.C1[0] - DENIT/3
	g.C1[1] = g.C1[1] - DENIT/3
	g.C1[2] = g.C1[2] - DENIT/3
	g.CUMDENIT = g.CUMDENIT + DENIT
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
	Nquadrat1 := math.Pow(nitratOb30, 2)
	michment1 := (Vmax * Nquadrat1) / (Nquadrat1 + KNO3)
	Ftheta1 := 1 - math.Exp(-1*math.Pow((thetarel1/Okrt), 6))
	Ftemp1 := 1 - math.Exp(-1*math.Pow((tempOb30/Tkrt), 4.6))
	Denit1 := michment1 * Ftheta1 * Ftemp1
	Denit1 = Denit1 / 1000 // (kg/ha)
	Nquadrat2 := math.Pow(nitratOb60, 2)
	michment2 := (Vmax * Nquadrat2) / (Nquadrat2 + KNO3)
	Ftheta2 := 1 - math.Exp(-1*math.Pow((thetarel2/Okrt), 6))
	Ftemp2 := 1 - math.Exp(-1*math.Pow((tempOb60/Tkrt), 4.6))
	Denit2 := michment2 * Ftheta2 * Ftemp2
	Denit2 = Denit2 / 1000 // (kg/ha)
	Nquadrat3 := math.Pow(nitratOb90, 2)
	michment3 := (Vmax * Nquadrat3) / (Nquadrat3 + KNO3)
	Ftheta3 := 1 - math.Exp(-1*math.Pow((thetarel3/Okrt), 6))
	Ftemp3 := 1 - math.Exp(-1*math.Pow((tempOb90/Tkrt), 4.6))
	Denit3 := michment3 * Ftheta3 * Ftemp3
	Denit3 = Denit3 / 1000 // (kg/ha)

	// !   /* Denitrifizierte N von Nitrat Pool wegnehmen                    */
	g.C1[0] = g.C1[0] - Denit1/3
	g.C1[1] = g.C1[1] - Denit1/3
	g.C1[2] = g.C1[2] - Denit1/3
	g.C1[3] = g.C1[3] - Denit2/3
	g.C1[4] = g.C1[4] - Denit2/3
	g.C1[5] = g.C1[5] - Denit2/3
	g.C1[6] = g.C1[6] - Denit3/3
	g.C1[7] = g.C1[7] - Denit3/3
	g.C1[8] = g.C1[8] - Denit3/3

	g.CUMDENIT = g.CUMDENIT + Denit1 + Denit2 + Denit3
}
