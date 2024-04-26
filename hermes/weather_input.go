package hermes

import (
	"fmt"
	"log"
	"strconv"
	"time"

	yaml "gopkg.in/yaml.v3"
)

// Input units for weather files
// temperature average  				°C 			(required)
// temperature minimum 					°C			(required)
// temperature maximum 					°C			(required)
// global radiation 					MJ m-2		(required, if 0 it will be calculated)
// precipitation  						mm			(required)
// relative humidity 					%			(required)
// wind 								m s-1		(required)
// water vapor saturation deficit		mmHg (Torr) (optional, required for ETpot=1 Haude formula)
// sun shine hours 						h			(optional)
// measurement height for wind			m			(optional, default 2m)
// evapo transpiration ET0				mm			(optional, required for ETpot=5 )
// altitude  							m			(optional)
// CO2 concentration 					ppm			(optional)

// WeatherDataShared all weather split in years
type WeatherDataShared struct {
	JAR      []int
	TMP      [][366]float64 // temperature avarage 						°C
	TMI      [][366]float64 // temperature minimum 						°C
	TMA      [][366]float64 // temperature maximum 						°C
	RADI     [][366]float64 // photosynthetic active radiation 			MJ m-2
	REG      [][366]float64 // preciptation (optional on ground level) 	cm
	RELF     [][366]float64 // relative humidity 						%
	WIN      [][366]float64 // wind (capped to minimum 0.5)				m s-1
	VERD     [][366]float64 // water vapour saturation deficit			mmHg (Torr)
	SUND     [][366]float64 // sun shine hours 							h
	ETNULL   [][366]float64 // evapo transpiration ET0					mm
	WINDHI   float64        // measurment height for wind				m
	ALTITUDE float64        // altidude 								m
	CO2KONZ  []float64      // CO2 concentration 						ppm

	MaxYearDays []int // days in each year (365 or 366)
	// flags for optional parameters (if true the corresponting arrays contain valid values)
	hasWINDHI   bool
	hasALTITUDE bool
	hasCO2KONZ  bool
	hasVERD     bool
	hasSUND     bool
	hasETNULL   bool
}

// NewWeatherDataShared returns a new WeatherDataShared struct
func NewWeatherDataShared(years int, baseCO2 float64) WeatherDataShared {

	s := WeatherDataShared{
		JAR:         make([]int, years),
		TMP:         make([][366]float64, years),
		TMI:         make([][366]float64, years),
		TMA:         make([][366]float64, years),
		RADI:        make([][366]float64, years),
		REG:         make([][366]float64, years),
		RELF:        make([][366]float64, years),
		WIN:         make([][366]float64, years),
		VERD:        make([][366]float64, years),
		SUND:        make([][366]float64, years),
		ETNULL:      make([][366]float64, years),
		WINDHI:      2,
		ALTITUDE:    0,
		CO2KONZ:     make([]float64, years),
		MaxYearDays: make([]int, years),
		hasWINDHI:   false,
		hasALTITUDE: false,
		hasCO2KONZ:  false,
		hasVERD:     false,
		hasSUND:     false,
		hasETNULL:   false,
	}
	s.fillCO2Value(baseCO2)
	return s
}

func (s *WeatherDataShared) fillCO2Value(co2 float64) {
	years := len(s.CO2KONZ)
	for y := 0; y < years; y++ {
		s.CO2KONZ[y] = co2
	}
}

// SUB WETTERK(VWDAT$)

// WetterK reads a climate file
// Format:
// Tp_av;Tpmin;Tpmax;ET0;rH;vappd14;wind;sundu;radia pr;prec;jday
// C_deg;C_deg;C_deg;mm;%;mm_Hg;m/sec;hours;MJ/m^2;mm ;
// 50;02;-----;-----;-----;-----;-----;-----;------;-- -;-
// 4.1;1;5.6;-99;90;0.2;3.1;0;64;1;1
// 4.8;3.9;6.3;-99;86;1.4;3.6;1.6;170;0;2
// ...
// seperator can be ',' or ';'
// character for decimal point is '.'
// the naming of the colums is irrelevant
// column content should be in following order:
// Temperature average;Temperature min;Temperture max;ET0;relative humidity; water vapour saturation deficit;wind;sun hours;global radiation;prepitation;year day
func WetterK(VWDAT string, year int, g *GlobalVarsMain, s *WeatherDataShared, hPath *HFilePath, driConfig *Config) error {
	//DIM CORRK(12),Wettin$(12),high$(2)
	var high, Wettin []string
	CORRK, err := ReadPreco(g, hPath)
	if err != nil {
		return err
	}
	// open weather file
	vwDatfile, scanner, err := Open(&FileDescriptior{
		FilePath:        VWDAT,
		FileDescription: "weather file",
		debugOut:        g.DEBUGCHANNEL,
		logID:           g.LOGID,
		ContinueOnError: true})
	if scanner == nil || vwDatfile == nil {
		return fmt.Errorf("failed to load file: %s! %v", VWDAT, err)
	}
	defer vwDatfile.Close()

	// if header consists of 3 lines (1. column names, 2. units, 3. global values)
	if driConfig.WeatherNumHeader == 3 {
		LineInut(scanner) // skip column names
		LineInut(scanner) // skip units
		heights := LineInut(scanner)
		high = Explode(heights, []rune{',', ';'})
		s.ALTITUDE = ValAsFloat(high[0], VWDAT, heights)
		s.hasALTITUDE = true
		s.WINDHI = ValAsFloat(high[1], VWDAT, heights)
		s.hasWINDHI = true

		if len(high) > 2 && high[2][0] != '-' {

			baseCO2 := ValAsFloat(high[2], VWDAT, high[2])
			s.fillCO2Value(baseCO2)
			s.hasCO2KONZ = true
		}
	} else {
		// skip any other type of header
		for i := 0; i < driConfig.WeatherNumHeader; i++ {
			LineInut(scanner)
		}
	}
	s.JAR[0] = year
	Tlast := 0
	for scanner.Scan() {
		WETTER := scanner.Text()
		Wettin = Explode(WETTER, []rune{',', ';'})
		// T = day of year
		T := int(ValAsInt(Wettin[10], VWDAT, WETTER))
		if Tlast+1 != T {
			return fmt.Errorf("%s Failed to parse file: %s, error: missing days", g.LOGID, VWDAT)
		}
		Tlast = T
		Tindex := T - 1
		s.TMP[0][Tindex] = ValAsFloat(Wettin[0], VWDAT, WETTER)
		s.TMI[0][Tindex] = ValAsFloat(Wettin[1], VWDAT, WETTER)
		s.TMA[0][Tindex] = ValAsFloat(Wettin[2], VWDAT, WETTER)
		s.ETNULL[0][Tindex] = ValAsFloat(Wettin[3], VWDAT, WETTER)
		s.hasETNULL = s.hasETNULL || driConfig.WeatherNoneValue != s.ETNULL[0][Tindex]
		s.RELF[0][Tindex] = ValAsFloat(Wettin[4], VWDAT, WETTER)
		s.VERD[0][Tindex] = ValAsFloat(Wettin[5], VWDAT, WETTER)
		s.hasVERD = s.hasVERD || driConfig.WeatherNoneValue != s.VERD[0][Tindex]
		s.WIN[0][Tindex] = ValAsFloat(Wettin[6], VWDAT, WETTER)
		s.SUND[0][Tindex] = ValAsFloat(Wettin[7], VWDAT, WETTER)
		s.hasSUND = s.hasSUND || driConfig.WeatherNoneValue != s.SUND[0][Tindex]
		s.RADI[0][Tindex] = ValAsFloat(Wettin[8], VWDAT, WETTER)
		s.REG[0][Tindex] = ValAsFloat(Wettin[9], VWDAT, WETTER)
		s.MaxYearDays[0] = T
	}

	s.replaceMissingValues(1, driConfig.WeatherNoneValue)
	s.transformWeatherData(1, CORRK[:])

	// END SUB
	return nil
}

type corrArr []float64

func (CORRK corrArr) getCorrValue(T int) float64 {
	var cor float64
	if T < 32 {
		cor = CORRK[0]
	} else if T < 60 {
		cor = CORRK[1]
	} else if T < 91 {
		cor = CORRK[2]
	} else if T < 121 {
		cor = CORRK[3]
	} else if T < 152 {
		cor = CORRK[4]
	} else if T < 182 {
		cor = CORRK[5]
	} else if T < 213 {
		cor = CORRK[6]
	} else if T < 244 {
		cor = CORRK[7]
	} else if T < 274 {
		cor = CORRK[8]
	} else if T < 305 {
		cor = CORRK[9]
	} else if T < 335 {
		cor = CORRK[10]
	} else {
		cor = CORRK[11]
	}

	return cor
}

// ReadPreco reads the pre correction file for precipitaion
func ReadPreco(g *GlobalVarsMain, hPath *HFilePath) ([12]float64, error) {
	var CORRK [12]float64
	if g.PRECO {
		PRECORR := hPath.precorr
		_, scanner, err := Open(&FileDescriptior{
			FilePath:        PRECORR,
			FileDescription: "preco file",
			UseFilePool:     true,
			debugOut:        g.DEBUGCHANNEL,
			logID:           g.LOGID,
			ContinueOnError: true})
		if scanner == nil {
			return CORRK, fmt.Errorf("failed to load file: %s! %v", PRECORR, err)
		}
		LineInut(scanner) // skip headline
		for scanner.Scan() {
			PKO := scanner.Text()
			M := int(ValAsInt(PKO[0:2], PRECORR, PKO))
			CORRK[M-1] = ValAsFloat(PKO[3:7], PRECORR, PKO)
		}
	} else {
		// no correction
		for m := 0; m < 12; m++ {
			CORRK[m] = 1
		}
	}
	return CORRK, nil
}

// Header for csv weather files
type Header int

const (
	isodate Header = iota
	doydate
	tmin
	tavg
	tmax
	precip
	globrad
	wind
	relhumid
	co2
	sunhours
	verd
)

var headerNames = map[string]Header{
	"iso-date": isodate,
	"tmin":     tmin,
	"tavg":     tavg,
	"tmax":     tmax,
	"precip":   precip,
	"globrad":  globrad,
	"wind":     wind,
	"relhumid": relhumid,
	"@YYYYJJJ": doydate,
	"RAD":      globrad,
	"TMAX":     tmax,
	"TMIN":     tmin,
	"RH":       relhumid,
	"WIND":     wind,
	"PREC":     precip,
	"CO2":      co2,
	"sunhours": sunhours,
	"SUNH":     sunhours,
	"sun":      sunhours,
	"VERD":     verd,
	"verd":     verd,
}

func readHeader(line string) map[Header]int {
	tokens := Explode(line, []rune{',', ';', '\t', ' '})
	headers := make(map[Header]int)
	for kHeader, vHeader := range headerNames {
		for i, token := range tokens {
			if token == kHeader {
				headers[vHeader] = i
				break
			}
		}
	}
	return headers
}

// ReadWeatherCSV read a weather file
func ReadWeatherCSV(VWDAT string, startyear int, g *GlobalVarsMain, s *WeatherDataShared, hPath *HFilePath, driConfig *Config) error {

	// read pre correction file for precipitation
	CORRK, err := ReadPreco(g, hPath)
	if err != nil {
		return err
	}
	// open weather file with multible years
	vwDatfile, scanner, _ := Open(&FileDescriptior{
		FilePath:        VWDAT,
		FileDescription: "weather file",
		debugOut:        g.DEBUGCHANNEL,
		logID:           g.LOGID,
		ContinueOnError: true})
	if scanner == nil || vwDatfile == nil {
		return fmt.Errorf("failed to load file: %s! %v", VWDAT, err)
	}
	defer vwDatfile.Close()

	line := LineInut(scanner)
	h := readHeader(line)

	// if header consists of 3 lines (1. column names, 2. units, 3. global values)
	if driConfig.WeatherNumHeader == 3 {
		LineInut(scanner) // skip units
		heights := LineInut(scanner)
		high := Explode(heights, []rune{',', ';'})
		s.ALTITUDE = ValAsFloat(high[0], VWDAT, heights)
		s.hasALTITUDE = true
		s.WINDHI = ValAsFloat(high[1], VWDAT, heights)
		s.hasWINDHI = true

		if len(high) > 2 && high[2][0] != '-' {
			baseCO2 := ValAsFloat(high[2], VWDAT, high[2])
			s.fillCO2Value(baseCO2)
			s.hasCO2KONZ = true
		}
	} else {
		// skip other header lines
		for i := 1; i < driConfig.WeatherNumHeader; i++ {
			LineInut(scanner)
		}
	}

	T := 0
	yrz := 0
	first := true
	for scanner.Scan() {
		line := scanner.Text()
		T++
		tokens := Explode(line, []rune{',', ';', '\t'})

		type weatherdate struct {
			wind     float64
			precip   float64
			globrad  float64
			tmax     float64
			tmin     float64
			tavg     float64
			relhumid float64
			sunh     float64
			verd     float64
			datetime time.Time
		}
		d := weatherdate{
			wind:     0,
			precip:   0,
			globrad:  driConfig.WeatherNoneValue,
			tmax:     0,
			tmin:     0,
			tavg:     0,
			relhumid: 0,
			sunh:     driConfig.WeatherNoneValue,
			verd:     driConfig.WeatherNoneValue,
			datetime: time.Time{},
		}
		err := make([]error, 10)
		isodate := tokens[h[isodate]]
		d.datetime, err[7] = time.Parse("2006-01-02", isodate)
		// skip years before start year
		if d.datetime.Year() < startyear {
			continue
		}
		d.wind, err[0] = strconv.ParseFloat(tokens[h[wind]], 64)
		d.precip, err[1] = strconv.ParseFloat(tokens[h[precip]], 64)
		if _, ok := h[globrad]; ok {
			d.globrad, err[2] = strconv.ParseFloat(tokens[h[globrad]], 64)
		}
		d.tmax, err[3] = strconv.ParseFloat(tokens[h[tmax]], 64)
		d.tmin, err[4] = strconv.ParseFloat(tokens[h[tmin]], 64)
		d.tavg, err[5] = strconv.ParseFloat(tokens[h[tavg]], 64)
		d.relhumid, err[6] = strconv.ParseFloat(tokens[h[relhumid]], 64)
		if _, ok := h[sunhours]; ok {
			d.sunh, err[8] = strconv.ParseFloat(tokens[h[sunhours]], 64)
			if (d.sunh > 24 || d.sunh < 0) && d.sunh != driConfig.WeatherNoneValue {
				err[8] = fmt.Errorf("sunhours should be a value between 0-24 (on earth!) ->'%s' ", tokens[h[sunhours]])
			}
			s.hasSUND = true
		}
		if _, ok := h[verd]; ok {
			d.verd, err[9] = strconv.ParseFloat(tokens[h[verd]], 64)
			if d.verd != driConfig.WeatherNoneValue && d.verd > 0 {
				s.hasVERD = true
			}
		}

		if anyError := anyWeatherError(err, g.LOGID); anyError != nil {
			return anyError
		}
		if first {
			// failsave if the first date is not 1.Jan
			first = false
			T = d.datetime.YearDay()
			yrz = 1
		} else if d.datetime.Day() == 1 && d.datetime.Month() == time.January {
			T = 1
			yrz = yrz + 1
		}
		if d.datetime.YearDay() != T {
			return fmt.Errorf("%s Failed to parse file: %s, error: missing days", g.LOGID, VWDAT)
		}
		if yrz > len(s.JAR) {
			yrz--
			break
		}
		s.JAR[yrz-1] = d.datetime.Year()
		s.TMP[yrz-1][T-1] = d.tavg
		s.TMI[yrz-1][T-1] = d.tmin
		s.TMA[yrz-1][T-1] = d.tmax
		s.RELF[yrz-1][T-1] = d.relhumid
		s.RADI[yrz-1][T-1] = d.globrad
		s.WIN[yrz-1][T-1] = d.wind
		s.REG[yrz-1][T-1] = d.precip
		s.SUND[yrz-1][T-1] = d.sunh
		s.VERD[yrz-1][T-1] = d.verd
		s.MaxYearDays[yrz-1] = T
	}
	s.replaceMissingValues(yrz, driConfig.WeatherNoneValue)
	// apply value changes
	s.transformWeatherData(yrz, CORRK[:])
	return nil
}

func anyWeatherError(list []error, logid string) error {
	for _, b := range list {
		if b != nil {
			return fmt.Errorf("%s Failed to parse file: %s, error :%v", logid, logid, b)
		}
	}
	return nil
}

// ReadWeatherCZ read weather file (cz format)
func ReadWeatherCZ(VWDAT string, startyear int, g *GlobalVarsMain, s *WeatherDataShared, hPath *HFilePath, driConfig *Config) error {

	// read pre correction file for precipitation
	CORRK, err := ReadPreco(g, hPath)
	if err != nil {
		return err
	}
	// open weather file with multible years
	vwDatfile, scanner, err := Open(&FileDescriptior{
		FilePath:        VWDAT,
		FileDescription: "weather file",
		debugOut:        g.DEBUGCHANNEL,
		logID:           g.LOGID,
		ContinueOnError: true})
	if scanner == nil || vwDatfile == nil {
		return fmt.Errorf("failed to load file: %s! %v", VWDAT, err)
	}
	defer vwDatfile.Close()

	line := LineInut(scanner)
	h := readHeader(line)
	//@YYYYJJJ     RAD    TMAX    TMIN      RH    WIND    PREC     CO2    SUNH

	// skip header lines
	for i := 1; i < driConfig.WeatherNumHeader; i++ {
		LineInut(scanner)
	}

	T := 0
	yrz := 0
	first := true
	currentCO2 := s.CO2KONZ[0] // baseCO2 for first year
	for scanner.Scan() {
		line := scanner.Text()
		T++
		tokens := Explode(line, []rune{',', ';', '\t', ' '})

		type weatherdate struct {
			wind     float64
			precip   float64
			globrad  float64
			sunh     float64
			tmax     float64
			tmin     float64
			tavg     float64
			relhumid float64
			verd     float64
			datetime time.Time
		}
		var d weatherdate
		err := make([]error, 10)
		doydate := tokens[h[doydate]]
		//time = yyyydoy
		d.datetime, err[0] = time.Parse("2006002", doydate)
		// skip years before start year
		if d.datetime.Year() < startyear {
			continue
		}
		if d.datetime.Year() < startyear {
			continue
		}

		d.wind, err[1] = strconv.ParseFloat(tokens[h[wind]], 64)
		d.precip, err[2] = strconv.ParseFloat(tokens[h[precip]], 64)
		d.tmax, err[3] = strconv.ParseFloat(tokens[h[tmax]], 64)
		d.tmin, err[4] = strconv.ParseFloat(tokens[h[tmin]], 64)
		d.relhumid, err[5] = strconv.ParseFloat(tokens[h[relhumid]], 64)

		// sunhours and/or rad
		if _, ok := h[sunhours]; ok {
			d.sunh, err[6] = strconv.ParseFloat(tokens[h[sunhours]], 64)
			if (d.sunh > 24 || d.sunh < 0) && d.sunh != driConfig.WeatherNoneValue {
				err[6] = fmt.Errorf("sunhours should be a value between 0-24 (on earth!) ->'%s' ", tokens[h[sunhours]])
			}
			s.hasSUND = true
		}
		if _, ok := h[globrad]; ok {
			d.globrad, err[7] = strconv.ParseFloat(tokens[h[globrad]], 64)
		}
		if _, ok := h[verd]; ok {
			d.verd, err[8] = strconv.ParseFloat(tokens[h[verd]], 64)
			if d.verd != driConfig.WeatherNoneValue && d.verd > 0 {
				s.hasVERD = true
			}
		}
		// optional co2 token, if left empty the previous co2 value will persist
		if _, ok := h[co2]; ok && len(tokens) > h[co2] {
			currentCO2, err[9] = strconv.ParseFloat(tokens[h[co2]], 64)
		}
		if anyError := anyWeatherError(err, g.LOGID); anyError != nil {
			return anyError
		}
		d.tavg = (d.tmax + d.tmin) / 2

		if first {
			// failsave if the first date is not 1.Jan
			first = false
			T = d.datetime.YearDay()
			yrz = 1
		} else if d.datetime.Day() == 1 && d.datetime.Month() == time.January {
			T = 1
			yrz = yrz + 1
		}
		if d.datetime.YearDay() != T {
			return fmt.Errorf("%s Failed to parse file: %s, error: missing days", g.LOGID, VWDAT)
		}
		if yrz > len(s.JAR) {
			yrz--
			break
		}
		s.JAR[yrz-1] = d.datetime.Year()
		s.TMP[yrz-1][T-1] = d.tavg
		s.TMI[yrz-1][T-1] = d.tmin
		s.TMA[yrz-1][T-1] = d.tmax
		s.RELF[yrz-1][T-1] = d.relhumid
		s.RADI[yrz-1][T-1] = d.globrad
		s.WIN[yrz-1][T-1] = d.wind
		s.REG[yrz-1][T-1] = d.precip
		s.SUND[yrz-1][T-1] = d.sunh
		s.VERD[yrz-1][T-1] = d.verd
		s.CO2KONZ[yrz-1] = currentCO2
		s.hasCO2KONZ = true
		s.MaxYearDays[yrz-1] = T
	}
	s.replaceMissingValues(yrz, driConfig.WeatherNoneValue)
	// apply value changes
	s.transformWeatherData(yrz, CORRK[:])

	return nil
}

func (s *WeatherDataShared) transformWeatherData(yrz int, corr corrArr) {
	for y := 0; y < yrz; y++ {
		T := s.MaxYearDays[y]
		for index := 0; index < T; index++ {
			cor := corr.getCorrValue(index + 1)
			// water model for rivers calculates in cm, so mm is transformed to cm by dividing by 10

			// correction of precipitation (turn on/off in config)
			// correction of rain in standard-Hellmann-Rainwater measurement in 1m height to what arrives on the ground.
			// Which is in average 10% higher caused by drift due to wind
			// if turned off, all 'cor' values will be 1
			s.REG[y][index] = s.REG[y][index] / 10 * cor

			// transform global radiation to PAR(photosynthetic active radiation), which is 50% of Global radiation.
			s.RADI[y][index] = s.RADI[y][index] / 2

			// correct wind to a minimum of 0.5 for ET0 calculations
			if s.WIN[yrz-1][T-1] < 0.5 {
				s.WIN[yrz-1][T-1] = 0.5
			}
		}
	}
}

func (s *WeatherDataShared) replaceMissingValues(yrz int, noneValue float64) {

	prevIndex, prevYear := -1, -1
	nextIndex, nextYear := -1, -1
	for y := 0; y < yrz; y++ {
		T := s.MaxYearDays[y]

		for index := 0; index < T; index++ {

			prevIndex, prevYear = index-1, y
			nextIndex, nextYear = index+1, y
			if nextIndex >= T {
				nextIndex = 1
				nextYear = nextYear + 1
				if nextYear >= yrz {
					nextYear = -1
					nextIndex = -1
				}
			}
			if prevIndex < 0 && y > 0 {
				prevIndex = s.MaxYearDays[y-1] - 1
				prevYear = prevYear - 1
			}

			if prevIndex >= 0 && prevYear >= 0 && nextIndex >= 0 && nextYear >= 0 {

				if s.TMP[y][index] == noneValue &&
					s.TMP[prevYear][prevIndex] != noneValue &&
					s.TMP[nextYear][nextIndex] != noneValue {
					s.TMP[y][index] = (s.TMP[prevYear][prevIndex] + s.TMP[nextYear][nextIndex]) / 2
				}

				if s.VERD[y][index] == noneValue &&
					s.VERD[prevYear][prevIndex] != noneValue &&
					s.VERD[nextYear][nextIndex] != noneValue {
					s.VERD[y][index] = (s.VERD[prevYear][prevIndex] + s.VERD[nextYear][nextIndex]) / 2
				}

				if s.SUND[y][index] == noneValue &&
					s.SUND[prevYear][prevIndex] != noneValue &&
					s.SUND[nextYear][nextIndex] != noneValue {
					s.SUND[y][index] = (s.SUND[prevYear][prevIndex] + s.SUND[nextYear][nextIndex]) / 2
				}
			} else {
				if s.TMP[y][index] == noneValue {
					s.TMP[y][index] = 0
				}
				if s.VERD[y][index] == noneValue {
					s.VERD[y][index] = 0
				}
				if s.SUND[y][index] == noneValue {
					s.SUND[y][index] = 0
				}
			}
			if s.SUND[y][index] == noneValue {
				s.SUND[y][index] = 0
			}
			if s.RADI[y][index] == noneValue {
				s.RADI[y][index] = 0
			}
			if s.REG[y][index] == noneValue {
				s.REG[y][index] = 0
			}
		}
	}
}

// LoadYear loads weather data from WeatherDataShared of a given year into global GlobalVarsMain
func LoadYear(g *GlobalVarsMain, s *WeatherDataShared, year int) error {

	checkTminTmax := func(tmin, tmax float64) bool {
		if tmin > tmax+0.5 {
			errorStr := fmt.Sprintf("%s Error in Weather data: Tmin(%0.3f) > Tmax(%0.3f) ", g.LOGID, tmin, tmax)
			if g.DEBUGCHANNEL != nil {
				g.DEBUGCHANNEL <- errorStr
			} else {
				log.Print(errorStr)
			}
			return false
		}
		return true
	}
	loadedYears := len(s.MaxYearDays)
	for yearIdx := 0; yearIdx < loadedYears; yearIdx++ {
		days := s.MaxYearDays[yearIdx]
		if s.JAR[yearIdx] == year {
			for Tidx := 0; Tidx < days; Tidx++ {

				g.TEMP[Tidx] = s.TMP[yearIdx][Tidx]
				g.TMIN[Tidx] = s.TMI[yearIdx][Tidx]
				g.TMAX[Tidx] = s.TMA[yearIdx][Tidx]
				// check if measured temperature values are valid
				if !checkTminTmax(g.TMIN[Tidx], g.TMAX[Tidx]) {
					g.TMIN[Tidx] = s.TMA[yearIdx][Tidx]
					g.TMAX[Tidx] = s.TMI[yearIdx][Tidx]
				}

				g.RH[Tidx] = s.RELF[yearIdx][Tidx]
				g.RAD[Tidx] = s.RADI[yearIdx][Tidx]
				g.WIND[Tidx] = s.WIN[yearIdx][Tidx]
				g.REGEN[Tidx] = s.REG[yearIdx][Tidx]
				if s.hasSUND {
					g.SUND[Tidx] = s.SUND[yearIdx][Tidx]
				}
				if s.hasVERD {
					g.VERD[Tidx] = s.VERD[yearIdx][Tidx]
				}
				if s.hasETNULL {
					g.ETNULL[Tidx] = s.ETNULL[yearIdx][Tidx]
				}
			}
			if s.hasWINDHI {
				g.WINDHI = s.WINDHI
			}
			if s.hasCO2KONZ {
				g.CO2KONZ = s.CO2KONZ[yearIdx]
			}
			if s.hasALTITUDE {
				g.ALTI = s.ALTITUDE
			}

			g.JTAG = days
			return nil
		}
	}
	return fmt.Errorf(`requested year (%d) was not loaded: loaded years %d - %d `, year, s.JAR[0], s.JAR[loadedYears-1])
}

// check if sunhours are given
func (s *WeatherDataShared) HasSunHours() bool {
	return s.hasSUND
}

func DumpWeatherDataToFile(filename string, s *WeatherDataShared) {

	file := OpenResultFile(filename, false)
	defer file.Close()

	data, err := yaml.Marshal(s)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	if _, err := file.WriteBytes(data); err != nil {
		log.Fatal(err)
	}
}
