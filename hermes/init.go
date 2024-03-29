package hermes

import (
	"math"
)

// Init adopts the start conditions
func Init(g *GlobalVarsMain) {

	g.TAG.SetByIndex(g.ITAG - 2)
	if g.GROUNDWATERFROM == Polygonfile {
		g.GRW = g.GW - (g.AMPL * math.Sin((g.TAG.Num+float64(g.GWPhase))*math.Pi/180))
	} else if g.GROUNDWATERFROM == GWTimeSeries {
		g.GRW, _ = GetGroundWaterLevel(g, g.BEGINN-2)
	}
	g.TSOIL[0][0] = (g.TMIN[g.ITAG-1] + g.TMAX[g.ITAG-1]) / 2
	initp := (g.TSOIL[0][0] - g.TBASE) / float64(g.N)
	for i := 1; i <= g.N; i++ {
		g.TSOIL[0][i] = g.TSOIL[0][0] - initp*float64(i)
	}
	g.ALBEDO = 0.2
	setFieldCapacityWithGW(g)

	var FKPROZ float64
	if g.TAG.Num < 275 {
		if g.GRW < 10 {
			FKPROZ = .5
		} else {
			FKPROZ = .4
		}
	} else {
		if g.GRW < 10 {
			FKPROZ = .65
		} else {
			FKPROZ = .6
		}
	}
	if g.FEU == 1 {
		FKPROZ = FKPROZ - 0.3
	} else if g.FEU == 3 {
		FKPROZ = FKPROZ + .3
	}
	for z := 0; z < g.N; z++ {
		zNum := float64(z) + 1
		if zNum > 15 {
			g.WG[0][z] = g.W[z] - (g.W[z]-g.WMIN[z])*(1-0.95)
			if g.WG[0][z] < g.WMIN[z] {
				g.WG[0][z] = g.WMIN[z]
			}
		} else {
			g.WG[0][z] = g.W[z] - (g.W[z]-g.WMIN[z])*(1-FKPROZ)
			if g.WG[0][z] < g.WMIN[z] {
				g.WG[0][z] = g.WMIN[z]
			}
		}
		PG := g.PORGES[z]
		if zNum >= g.GRW {
			g.WG[0][z] = PG
		}
		g.C1[z] = g.CN[0][z]
		if zNum > 0 && zNum < 40./g.DZ.Num {
			g.NAOS[z] = g.NALTOS / 30 * g.DZ.Num
			g.NFOS[z] = 0
			g.MINAOS[z] = 0
			g.MINFOS[z] = 0
		}
		g.CA[z] = 0
	}
	g.WG[0][10] = g.WG[0][9]
	g.NDG.SetByIndex(0)
	g.MZ = 1
	g.NBR = 1
	g.NTIL.SetByIndex(0)
	// obsolete since everything is already 0 initialized
	g.REGENSUM = 0
	g.MINSUM = 0
	g.CAPSUM = 0
	g.RADSUM = 0
	g.BLATTSUM = 0
	g.DSUMM = 0
	g.UMS = 0
	g.OUTSUM = 0
	g.NFIXSUM = 0
	g.DRAISUM = 0
	g.DRAINLOSS = 0
	g.NFIX = 0
	g.SCHNORR = 0
}

func setFieldCapacityWithGW(g *GlobalVarsMain) {
	for l := int(g.GRW + 1); l <= g.N; l++ {
		if l == int(g.GRW+1) {
			g.W[l-1] = (1-math.Mod(g.GRW+1, 1))*g.PORGES[l-1] + g.W[l-1]*(math.Mod(g.GRW+1, 1))
		} else {
			g.W[l-1] = g.PORGES[l-1]
		}
	}
}

// setWRed calc reduced Field Capacity
// takes soil texture, wilting point and field capacity from top layer
func calcWRed(wiltingPoint float64, fieldCapacity float64, g *GlobalVarsMain) {
	if g.BART[0][0] == 'S' { // if main soil component is sand
		g.WRED = wiltingPoint + 0.6*(fieldCapacity-wiltingPoint)
	} else {
		g.WRED = wiltingPoint + 0.66*(fieldCapacity-wiltingPoint)
	}
	g.WRED = g.WRED / 100
}
