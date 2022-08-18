package hermes

import (
	"math"
)

// Init adopts the start conditions
func Init(g *GlobalVarsMain) {

	g.TAG.SetByIndex(g.ITAG - 2)
	g.GRW = g.GW - (float64(g.AMPL) * math.Sin((g.TAG.Num+80)*math.Pi/180))
	g.TSOIL[0][0] = (g.TMIN[g.ITAG-1] + g.TMAX[g.ITAG-1]) / 2
	initp := (g.TSOIL[0][0] - g.TBASE) / float64(g.N)
	for i := 1; i <= g.N; i++ {
		g.TSOIL[0][i] = g.TSOIL[0][0] - initp*float64(i)
	}
	g.ALBEDO = 0.2
	for l := int(g.GRW + 1); l <= g.N; l++ {
		if l == int(g.GRW+1) {
			g.W[l-1] = (1-math.Mod(g.GRW+1, 1))*g.PORGES[l-1] + g.W[l-1]*(math.Mod(g.GRW+1, 1))
		} else {
			g.W[l-1] = g.PORGES[l-1]
		}
	}

	var FKPROZ float64
	if g.TAG.Num < 275 {
		if g.GW < 10 {
			FKPROZ = .5
		} else {
			FKPROZ = .4
		}
	} else {
		if g.GW < 10 {
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
	_, applyInitialS := g.SI[0]
	for z0 := 0; z0 < g.N; z0++ {
		z1 := float64(z0) + 1
		if z1 > 15 {
			g.WG[0][z0] = g.W[z0] - (g.W[z0]-g.WMIN[z0])*(1-0.95)
			if g.WG[0][z0] < g.WMIN[z0] {
				g.WG[0][z0] = g.WMIN[z0]
			}
		} else {
			g.WG[0][z0] = g.W[z0] - (g.W[z0]-g.WMIN[z0])*(1-FKPROZ)
			if g.WG[0][z0] < g.WMIN[z0] {
				g.WG[0][z0] = g.WMIN[z0]
			}
		}
		PG := g.PORGES[z0]
		if z1 >= g.GW {
			g.WG[0][z0] = PG
		}
		g.C1[z0] = g.CN[0][z0]
		if z1 > 0 && z1 < 40./g.DZ.Num {
			g.NAOS[z0] = g.NALTOS / 30 * g.DZ.Num
			g.NFOS[z0] = 0
			g.MINAOS[z0] = 0
			g.MINFOS[z0] = 0
		}
		// LET S1(z) = SI(0,z)
		// check if initial smin data exists
		if applyInitialS {
			g.S1[z0] = g.SI[0][z0]
		}

		// make sure that initial smin data is > 0
		if g.S1[z0] <= 0 {
			g.S1[z0] = 0.01
		}
		// IF Z > 0 AND Z < 40/DZ THEN
		if z1 > 0 && z1 < 40./g.DZ.Num {
			//LET SAOS(Z) = SALTOS/30*DZ
			g.SAOS[z0] = g.SALTOS / 30 * g.DZ.Num
			//LET SFOS(Z) = 0
			g.SFOS[z0] = 0
			//LET Sminaos(z) = 0
			g.Sminaos[z0] = 0
			//LET Sminfos(z) = 0
			g.Sminfos[z0] = 0
		}
		g.CA[z0] = 0
	}
	if applyInitialS {
		g.sMessIdx++
	}
	g.WG[0][10] = g.WG[0][9]
	g.DG.SetByIndex(0) // Nitrogen fertilization counter
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
	g.SDSUMM = 0
	g.UMS = 0
	g.OUTSUM = 0
	g.NFIXSUM = 0
	g.DRAISUM = 0
	g.DRAINLOSS = 0
	g.NFIX = 0
	g.SCHNORR = 0
}
