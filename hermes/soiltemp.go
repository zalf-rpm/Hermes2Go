package hermes

import (
	"math"
)

// Soiltemp calculates and write the soil temperature data
func Soiltemp(g *GlobalVarsMain) {
	// ! Inputs:
	// ! LAI  		= Blattflächenindex
	// ! RAD(TAG)	= PAR (MJoule/m^2)
	// ! TEMP(TAG)	= Tagesmitteltemperatur (°C)
	// ! TMIN(TAG) 	= Tagesminimumtemperatur (°C)
	// ! TMAX(TAG	= Tagesmaximumtemp ((°C)
	// ! HUMUS(I)	= HUMUSgehalt schicht I (%)
	// ! BD(I)		= Lagerungsdichte I (g/cm^3)
	// ! WG(0,I)	= Wassergehakt (cm^3/cm^3)
	// ! WATABIL	= Wasserbilanz
	// ! TSOIL(0,I)	= Bodentemperatur am Anfang des Zeitschritts
	// ! Variable:
	// ! HEATCOND(I)= Wärmeleitfähigkeit
	// ! HEATCAP(I)	= Wärmekapazität
	// ! OUTPUT:
	// ! TSOIL(0,I)=TSOIL(1,I)	= Bodentemperatur am Ende des Zeitschritts
	var scov float64
	var radiat float64
	if g.LAI < 3 {
		scov = 1 - math.Exp(-g.LAI)
		if scov < 0 {
			scov = 0
		}
		radiat = g.RAD[g.TAG.Index]*200*(1-scov) - (g.ETA * 10 * (2.498 - 0.00242*g.TEMP[g.TAG.Index]) * 10)
	} else {
		scov = 1
		radiat = 0
	}

	g.ALBEDO = 0.31
	g.TSOIL[0][g.N] = g.TBASE
	g.TSOIL[1][g.N] = g.TBASE
	if radiat > 833 {
		g.TSOIL[1][0] = ((1 - g.ALBEDO) * (g.TMIN[g.TAG.Index] + ((g.TMAX[g.TAG.Index] - g.TMIN[g.TAG.Index]) * math.Sqrt(0.0003*radiat)))) + (g.ALBEDO * g.TSOIL[0][0])
	} else {
		g.TSOIL[1][0] = (g.TMIN[g.TAG.Index] + g.TMAX[g.TAG.Index]) / 2
	}
	for i := 0; i < g.N; i++ {
		g.HEATCOND[i] = ((3*g.BD[i] - 1.7) * 0.001) / (1.0 + (11.5-5.0*g.BD[i])*math.Exp((-50)*math.Pow((g.WG[0][i]/g.BD[i]), 1.5))) * 86400 * g.DT.Num * 4.189
		g.HEATCAP[i] = (g.WG[0][i]*1*1 + (1-g.BD[i]/2.65-g.WG[0][i])*0.0013*0.23 + g.HUMUS[i]*1.3*1.3*0.45 + (g.BD[i]/2.65-g.HUMUS[i]*1.3)*2.65*0.18) * 4.189
		g.TDSUM[i] = 0
	}
	for std := 0; std < 24; std++ {
		for i := 1; i <= g.N-1; i++ {
			alpha := g.HEATCOND[i-1] / g.HEATCAP[i-1]
			g.TSOIL[1][i] = g.TSOIL[0][i] + alpha*(g.TSOIL[0][i+1]-2*g.TSOIL[0][i]+g.TSOIL[0][i-1])*g.DT.Num/24/math.Pow(g.DZ.Num, 2)
			g.TDSUM[i-1] = g.TDSUM[i-1] + g.TSOIL[1][i]
		}
		g.TD[0] = g.TSOIL[1][0]
		for i := 0; i <= g.N; i++ {
			g.TSOIL[0][i] = g.TSOIL[1][i]
		}
	}
	for i := 1; i <= g.N-1; i++ {
		g.TD[i] = g.TDSUM[i-1] / 24
	}
	g.TD[g.N] = g.TSOIL[0][g.N]
	for i := 1; i <= g.N; i++ {
		g.TSOIL[0][i] = g.TD[i]
	}
}
