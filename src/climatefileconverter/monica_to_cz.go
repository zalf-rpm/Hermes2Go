package main

import (
	"fmt"

	"github.com/zalf-rpm/Hermes2Go/hermes"
)

func ConvertFileMonicaToCZ(in, out string) error {

	var g hermes.GlobalVarsMain
	s := hermes.NewWeatherDataShared(36)
	var hPath hermes.HFilePath
	driConfig := hermes.NewDefaultConfig()
	driConfig.WeatherNoneValue = -99
	err := hermes.ReadWeatherCSV(in, 1900, &g, &s, &hPath, &driConfig)
	if err != nil {
		return err
	}

	headline := "@YYYYJJJ   TMIN    TMAX     RAD    PREC    WIND      RH \n"
	if s.HasSunHours() {
		headline = "@YYYYJJJ   TMIN    TMAX     RAD    PREC    WIND      RH     SUNH \n"
	}

	file := hermes.OpenResultFile(out, false)
	defer file.Close()

	if _, err := file.Write(headline); err != nil {
		return err
	}

	loadedYears := len(s.MaxYearDays)
	for idxYear := 0; idxYear < loadedYears; idxYear++ {
		for doy := 0; doy < s.MaxYearDays[idxYear]; doy++ {
			if s.HasSunHours() {
				file.Write(fmt.Sprintf(" %d%03d %6.1f %7.1f %7.2f %7.1f %7.1f %7.1f %7.1f",
					s.JAR[idxYear],
					doy+1,
					s.TMI[idxYear][doy],
					s.TMA[idxYear][doy],
					s.RADI[idxYear][doy]*2,
					s.REG[idxYear][doy]*10.0,
					s.WIN[idxYear][doy],
					s.RELF[idxYear][doy],
					s.SUND[idxYear][doy]))
			} else {
				file.Write(fmt.Sprintf(" %d%03d %6.1f %7.1f %7.2f %7.1f %7.1f %7.1f",
					s.JAR[idxYear],
					doy+1,
					s.TMI[idxYear][doy],
					s.TMA[idxYear][doy],
					s.RADI[idxYear][doy]*2,
					s.REG[idxYear][doy]*10.0,
					s.WIN[idxYear][doy],
					s.RELF[idxYear][doy]))
			}
			file.WriteRune('\n')
		}
	}

	return nil
}
