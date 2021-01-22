package hermes

import "math"

// CalculateDayLenght calculate day lenght, effective day lenght(DLE, DLP), extra terrestial Radiation
func CalculateDayLenght(tag float64, lat float64) (DL, DLE, DLP, EXT, RDN, DRC, DEC float64) {

	// -------- BERECHNUNG VON TAGLAENGE UND EINSTRAHLUNG -----
	// calculation of day length an radiation
	// ----------------------- DECLINATION -----------------------
	DEC = 0.409 * math.Sin(2*math.Pi/365*tag-1.39) * 180 / math.Pi
	SINLD := math.Sin(DEC*math.Pi/180.) * math.Sin(lat*math.Pi/180.)
	COSLD := math.Cos(DEC*math.Pi/180.) * math.Cos(lat*math.Pi/180.)
	// -------------------- ASTRONOMISCHE TAGESLäNGE ------------
	// astronomical daylenth
	DL = 12. * (math.Pi + 2.*math.Asin(Limit(SINLD/COSLD, 1, -1))) / math.Pi
	// -------------------- EFFEKTIVE TAGESLäNGE ----------------
	// effective day length
	DLE = 12. * (math.Pi + 2.*math.Asin(Limit((-math.Sin(8.*math.Pi/180.)+SINLD)/COSLD, 1, -1))) / math.Pi
	DLP = 12. * (math.Pi + 2.*math.Asin(Limit((-math.Sin(-6.*math.Pi/180.)+SINLD)/COSLD, 1, -1))) / math.Pi

	SC := 24. * 60. / math.Pi * 8.20 * (1 + 0.033*math.Cos(2*math.Pi*tag/365.))
	SHA := math.Acos(Limit(-math.Tan(lat*math.Pi/180)*math.Tan(DEC*math.Pi/180), 1, -1))
	EXT = SC * (SHA*SINLD + COSLD*math.Sin(SHA)) / 100.0 // from J cm-2 to MJ m-2

	// ----- MITTLERE PHOTOSYNTHETISCH AKTIVE EINSTRAHLUNG ------
	// average photosynthetic active radiation
	if DL > 0 {
		RDN = 3600. * (SINLD*DL + 24./math.Pi*COSLD*math.Sqrt(1.-math.Pow(Limit(SINLD/COSLD, 1, -1), 2)))
		// ------------Strahlung klarer Tag (in Joule/m^2)------------
		// radiation on a clear day (in Joule/m^2)
		DRC = 0.5 * 1300. * RDN * math.Exp(-.14/(RDN/(DL*3600.)))
	}

	return DL, DLE, DLP, EXT, RDN, DRC, DEC
}

// Limit a value to upper and lower bounds(inclusive)
func Limit(val, upper, lower float64) float64 {
	if val > upper {
		return upper
	} else if val < lower {
		return lower
	}
	return val
}
