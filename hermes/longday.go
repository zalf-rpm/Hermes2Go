package hermes

import (
	"log"
)

// LangTagConverterFunc get LangTagConverter
type LangTagConverterFunc func(float64, string, int) (int, int, int)

// LangTagConverter get function to get 14h and 16h Day in respect to latitude
func LangTagConverter(century int, dateFormat DateFormat) func(float64, string, int) (int, int, int) {
	cent := century
	format := dateFormat
	// LangTag calculate the 14h and 16h Day in respect to latitude
	return func(LAT float64, progDat string, anjahr int) (TAG, P1, P2 int) {
		TAG = 0
		P1 = 0
		P2 = 0
		for ok := true; ok; ok = P1 == 0 {
			TAG++
			DL, _, _, _, _, _, _ := CalculateDayLenght(float64(TAG), LAT)
			if DL > 14 {
				P1 = TAG
			}
		}
		for ok := true; ok; ok = P2 == 0 {
			TAG++
			DL, _, _, _, _, _, _ := CalculateDayLenght(float64(TAG), LAT)
			if DL > 16 {
				P2 = TAG // Beginn Große Periode
			}
		}
		if progDat[1] != '-' {
			var progja int
			if format == DateDEshort || format == DateENshort {
				_, _, progj, err := extractDate(progDat, true, false)
				if err != nil {
					log.Fatal(err)
				}
				if progj < cent {
					progja = 100 + progj
				} else {
					progja = progj
				}
			} else {
				_, _, progj, err := extractDate(progDat, false, false)
				if err != nil {
					log.Fatal(err)
				}
				progja = progj - 1900
			}

			P2 = (progja-1)*365 + (progja)/4 + P2
			P1 = P1 + 20 // Ungefährer Schossbeginn
			P1 = (progja-1)*365 + (progja)/4 + P1
		} else {
			P2 = (anjahr-1)*365 + (anjahr)/4 + P2
			P1 = P1 + 20 // Ungefährer Schossbeginn
			P1 = (anjahr-1)*365 + (anjahr)/4 + P1
		}
		return TAG, P1, P2
	}
}
