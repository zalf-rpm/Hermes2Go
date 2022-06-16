package hermes

import (
	"log"
	"os"
	"reflect"
	"strconv"

	yaml "gopkg.in/yaml.v2"
)

type config struct {

	//***** Formats *****
	// DateDEshort ddmmyy <- default format(old)
	// DateDElong ddmmyyyy
	// DateENshort mmddyy
	// DateENlong mmddyyyy
	// short format "ddmmyy", e.g. 24.01.95 -> "ddmmyy" requires input as 240195, you need to set the century devision year "DivideCentury", e.g. 1950 -> 50
	// long format "ddmmyyyy", e.g. 24.01.2066  requires input as 24012066
	Dateformat    DateFormat `yaml:"Dateformat"`
	DivideCentury int        `yaml:"DivideCentury,omitempty"` // (depends on Date format) Year to divide 20. and 21. Century (YY)

	GroundWaterFrom     GroundWaterFrom `yaml:"GroundWaterFrom"`            // ground water is read from either 'soilfile' or 'polygonfile
	ResultFileFormat    int             `yaml:"ResultFileFormat,omitempty"` // result file format (0= hermes default, 1 = csv)
	ResultFileExt       string          `yaml:"ResultFileExt,omitempty"`    // result file extensions (default RES, csv)
	OutputIntervall     int             `yaml:"OutputIntervall"`            // Output intervall (days) (0=no time serie)
	InitSelection       int             `yaml:"InitSelection"`              // Init.values all(1),Field_ID(2), Polyg(3) -> POLY_XXX.txt, Uses: 1= all (if the word ALLE is written in the file), 2= Field_ID, 3= Polyg
	SoilFile            string          `yaml:"SoilFile"`                   // soil profile file name (without projectname)
	SoilFileExtension   string          `yaml:"SoilFileExtension"`          // soil file extension (txt = hermes soil, csv = csv table format)
	CropFileFormat      string          `yaml:"CropFileFormat"`             // crop file format (txt = hermes crop, csv = csv table format)
	PolygonGridFileName string          `yaml:"PolygonGridFileName"`        // Name of Polygon resp. grid file

	//***** Weather *****
	WeatherFile              string        `yaml:"WeatherFile"`              // weather file name template (without projectname)
	WeatherFileFormat        int           `yaml:"WeatherFileFormat"`        // Weather file format (0=separator(, ; \t), 1 year per file ) (1=separator(, ; \t), multiple years per file, (1=separator(, ; \t ' '), cz format, multiple years per file)
	WeatherFolder            string        `yaml:"WeatherFolder"`            // Weather scenario folder
	WeatherRootFolder        string        `yaml:"WeatherRootFolder"`        // weather root directory without scenario folder or filename
	WeatherNoneValue         float64       `yaml:"WeatherNoneValue"`         // weather none value, default -99.9
	WeatherNumHeader         int           `yaml:"WeatherNumHeader"`         // number of header lines (min = 1, with column names)
	CorrectionPrecipitation  FeatureSwitch `yaml:"CorrectionPrecipitation"`  // correction precipitation (0= no, 1 = yes)
	AnnualAverageTemperature float64       `yaml:"AnnualAverageTemperature"` // annual average temperature (Celsius)

	//***** Atmosphere *****
	ETpot               int           `yaml:"ETpot"`               // ETpot method(1=Haude,2=Turc-Wendling,3 Penman-Monteith)
	CO2method           int           `yaml:"CO2method"`           // CO2method(1=Nonhebel,2=Hoffmann,3=Mitchell)
	CO2concentration    float64       `yaml:"CO2concentration"`    // CO2 concentration (ppm)
	CO2StomataInfluence FeatureSwitch `yaml:"CO2StomataInfluence"` // CO2 Stomata influence (1=on/0= off)
	NDeposition         float64       `yaml:"NDeposition"`         // N-Deposition (annual kg/ha)

	//***** Time *****
	StartYear                       int    `yaml:"StartYear"`                       // Starting year of simulation (YYYY)
	EndDate                         string `yaml:"EndDate"`                         // End date of simulation (DDMMYYYY)
	AnnualOutputDate                string `yaml:"AnnualOutputDate"`                // Date for annual output
	VirtualDateFertilizerPrediction string `yaml:"VirtualDateFertilizerPrediction"` // Virtual date for fertilizer prediction; '------' for no prediction

	//***** Geoography *****
	Latitude      float64 `yaml:"Latitude"`      // Latitude
	Altitude      float64 `yaml:"Altitude"`      // Altitude - height (can be overwritten weather file)
	CoastDistance float64 `yaml:"CoastDistance"` // Distance to coast (km)

	//***** Soil *****
	PTF                            int     `yaml:"PTF"`                            // PTF pedo transfer function (0 = none (from file), 1 = Toth 2015, 2 = Batjes for pF 2.5, 3 = Batjes for pF 1.7, 4 = Rawls et al. 2003 for pF 2.5 )
	LeachingDepth                  int     `yaml:"LeachingDepth"`                  // Depth for leaching/seepage calculation (dm)
	OrganicMatterMineralProportion float64 `yaml:"OrganicMatterMineralProportion"` // Mineralisable proportion of organic matter
	KcFactorBareSoil               float64 `yaml:"KcFactorBareSoil"`               // kc factor for bare soil

	//***** Management *****
	Fertilization     float64       `yaml:"Fertilization"`     // fertilization scenario (fertilization in %)
	AutoSowingHarvest FeatureSwitch `yaml:"AutoSowingHarvest"` // automatic sowing/harvest (0=no, 1 = yes)
	AutoFertilization FeatureSwitch `yaml:"AutoFertilization"` // automatic fertilization (0=no, 1=on demand)
	AutoIrrigation    FeatureSwitch `yaml:"AutoIrrigation"`    // automatic irrigation (0=no, 1= on demand)
	AutoHarvest       FeatureSwitch `yaml:"AutoHarvest"`       // automatic harvest (0=no, 1= on demand)
}

func readConfig(g *GlobalVarsMain, argValues map[string]string, hp *HFilePath) config {
	hconfig := NewDefaultConfig()
	// if config files exists, read it into hconfig
	if _, err := os.Stat(hp.config); err == nil {
		byteData := HermesFilePool.Get(&FileDescriptior{FilePath: hp.config, ContinueOnError: true, UseFilePool: true})
		err := yaml.Unmarshal(byteData, &hconfig)
		if err != nil {
			log.Fatalf("error: %v", err)
		}
	} else {
		// no config exist, generate default config (if project is not fitting default setup, execution will fail)
		hconfig = NewDefaultConfig()
	}
	err := commandlineOverride(argValues, &hconfig)
	if err != nil {
		log.Fatalf("error while parsing commandline: %v", err)
	}

	g.GROUNDWATERFROM = hconfig.GroundWaterFrom
	g.DATEFORMAT = hconfig.Dateformat
	g.ANJAHR = hconfig.StartYear
	g.INIWAHL = hconfig.InitSelection
	g.PRECO = bool(hconfig.CorrectionPrecipitation)
	g.CO2METH = hconfig.CO2method
	g.CO2KONZ = hconfig.CO2concentration
	g.CTRANS = bool(hconfig.CO2StomataInfluence)
	g.ETMETH = hconfig.ETpot
	g.LAT = hconfig.Latitude
	g.ALTI = hconfig.Altitude
	g.OUTN = hconfig.LeachingDepth
	g.NAKT = hconfig.OrganicMatterMineralProportion
	g.DEPOS = hconfig.NDeposition
	g.DUNGSZEN = hconfig.Fertilization / 100
	g.Datum = DateConverter(hconfig.DivideCentury, g.DATEFORMAT)
	g.Kalender = KalenderConverter(g.DATEFORMAT, ".")
	g.LangTag = LangTagConverter(hconfig.DivideCentury, g.DATEFORMAT)
	_, g.ENDE = g.Datum(hconfig.EndDate)
	g.FKB = hconfig.KcFactorBareSoil
	g.TBASE = hconfig.AnnualAverageTemperature
	g.AUTOMAN = bool(hconfig.AutoSowingHarvest)
	g.AUTOFERT = bool(hconfig.AutoFertilization)
	g.AUTOIRRI = bool(hconfig.AutoIrrigation)
	g.AUTOHAR = bool(hconfig.AutoHarvest)
	g.PTF = hconfig.PTF
	if len(hconfig.WeatherFolder) == 0 {
		hconfig.WeatherFolder = "Weather"
	}
	if len(hconfig.WeatherRootFolder) == 0 {
		hconfig.WeatherRootFolder = hp.path
	}
	if len(hconfig.ResultFileExt) == 0 {
		if OutputFileFormat(hconfig.ResultFileFormat) == csvOut {
			hconfig.ResultFileExt = "csv"
		} else {
			hconfig.ResultFileExt = "RES"
		}
	}
	return hconfig
}

// commandlineOverride will parse through the commandline and try to parse them into config
func commandlineOverride(argValues map[string]string, hconfig *config) error {
	if argValues != nil {
		v := reflect.ValueOf(hconfig)
		if v.Kind() == reflect.Ptr && !v.IsNil() {
			v = v.Elem()
		}
		for argKey, argVal := range argValues {
			f := v.FieldByName(argKey)
			if f.IsValid() {
				if f.CanSet() {
					if f.Kind() == reflect.Float64 {
						newRefVal, err := strconv.ParseFloat(argVal, 64)
						if err != nil {
							return err
						}
						if !f.OverflowFloat(newRefVal) {
							f.SetFloat(newRefVal)
						}
					}

					if f.Kind() == reflect.Int {
						newRefVal, err := strconv.ParseInt(argVal, 10, 64)
						if err != nil {
							return err
						}
						if !f.OverflowInt(newRefVal) {
							f.SetInt(newRefVal)
						}
					}
					if f.Kind() == reflect.String {
						f.SetString(argVal)
					}
					if f.Kind() == reflect.Bool {
						if fs, ok := featureSwitchStrToID[argVal]; ok {
							f.SetBool(bool(fs))
						}
					}
				}
			}
		}
	}
	return nil
}

// WriteYamlConfig write a default config file
func WriteYamlConfig(filename string, structIn interface{}) {
	file := OpenResultFile(filename, false)
	defer file.Close()
	data, err := yaml.Marshal(structIn)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	if _, err := file.WriteBytes(data); err != nil {
		log.Fatal(err)
	}
}

//NewDefaultConfig creates a config file with default setup
func NewDefaultConfig() config {
	return config{
		Dateformat:                      DateDElong,
		GroundWaterFrom:                 Soilfile,
		ResultFileFormat:                0,
		OutputIntervall:                 0,
		InitSelection:                   3,
		SoilFile:                        "soil",
		SoilFileExtension:               "txt",
		CropFileFormat:                  "txt",
		WeatherFile:                     "%s.csv",
		WeatherFileFormat:               1,
		WeatherFolder:                   "Weather",
		WeatherRootFolder:               "",
		WeatherNoneValue:                -99.9,
		WeatherNumHeader:                2,
		CorrectionPrecipitation:         false,
		ETpot:                           3,
		PTF:                             0,
		CoastDistance:                   300,
		CO2method:                       2,
		CO2concentration:                360,
		CO2StomataInfluence:             true,
		StartYear:                       1980,
		EndDate:                         "31122010",
		AnnualOutputDate:                "3009",
		Latitude:                        52.52,
		Altitude:                        0,
		LeachingDepth:                   15,
		OrganicMatterMineralProportion:  0.13,
		NDeposition:                     20,
		PolygonGridFileName:             "poly",
		Fertilization:                   100,
		VirtualDateFertilizerPrediction: "--------",
		KcFactorBareSoil:                0.4,
		AnnualAverageTemperature:        8.7,
		AutoFertilization:               true,
		AutoSowingHarvest:               true,
		AutoIrrigation:                  true,
		AutoHarvest:                     true,
	}
}

// GroundWaterFrom enum from where to load ground water data
type GroundWaterFrom int

const (
	// Polygonfile load from polygon file
	Polygonfile GroundWaterFrom = iota
	// Soilfile load from soil file
	Soilfile
)

func (s GroundWaterFrom) String() string {
	return toString[s]
}

var toString = map[GroundWaterFrom]string{
	Polygonfile: "polygonfile",
	Soilfile:    "soilfile",
}

var toID = map[string]GroundWaterFrom{
	"polygonfile": Polygonfile,
	"soilfile":    Soilfile,
}

// MarshalYAML implement YAML Marshaler
func (s GroundWaterFrom) MarshalYAML() (interface{}, error) {
	return toString[s], nil
}

// UnmarshalYAML implement YAML Unmarshaler interface
func (s *GroundWaterFrom) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var j string
	err := unmarshal(&j)
	if err != nil {
		return err
	}
	// if the string cannot be found, it will be set to 'soilfile'
	*s = toID[j]
	return nil
}

// FeatureSwitch to convert 1/0 on/off to true/false
type FeatureSwitch bool

var featureSwitchToString = map[FeatureSwitch]int{
	true:  1,
	false: 0,
}
var featureSwitchStrToID = map[string]FeatureSwitch{
	"1":     true,
	"0":     false,
	"on":    true,
	"off":   false,
	"yes":   true,
	"no":    false,
	"true":  true,
	"false": false,
}

// MarshalYAML implement YAML Marshaler
func (s FeatureSwitch) MarshalYAML() (interface{}, error) {
	return featureSwitchToString[s], nil
}

// UnmarshalYAML implement YAML Unmarshaler interface
func (s *FeatureSwitch) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var j string
	err := unmarshal(&j)
	if err != nil {
		return err
	}
	*s = featureSwitchStrToID[j]
	return nil
}

// DateFormat - default date format for all hermes input
type DateFormat int

const (
	// DateDEshort ddmmyy
	DateDEshort DateFormat = iota
	// DateDElong ddmmyyyy
	DateDElong
	// DateENshort mmddyy
	DateENshort
	// DateENlong mmddyyyy
	DateENlong
)

func (s DateFormat) String() string {
	return dateFToString[s]
}

var dateFToString = map[DateFormat]string{
	DateDEshort: "DateDEshort",
	DateDElong:  "DateDElong",
	DateENshort: "DateENshort",
	DateENlong:  "DateENlong",
}

var dateStrToID = map[string]DateFormat{
	"DateDEshort": DateDEshort,
	"DateDElong":  DateDElong,
	"DateENshort": DateENshort,
	"DateENlong":  DateENlong,
}

// MarshalYAML implement YAML Marshaler
func (s DateFormat) MarshalYAML() (interface{}, error) {
	return dateFToString[s], nil
}

// UnmarshalYAML implement YAML Unmarshaler interface
func (s *DateFormat) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var j string
	err := unmarshal(&j)
	if err != nil {
		return err
	}
	// if the string cannot be found, it will be set to 'soilfile'
	*s = dateStrToID[j]
	return nil
}
