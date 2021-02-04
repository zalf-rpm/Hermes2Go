package hermes

import (
	"log"
)

// DualType for values that have been used as int index and float for calculations
type DualType struct {
	Index  int
	Num    float64
	Offset int
}

// Inc increments DualType
func (d *DualType) Inc() {
	d.Index++
	d.Num++
	if d.Index+d.Offset != int(d.Num) {
		log.Fatal("Float to Integer conversion error > 1")
	}
}

// SetByIndex DualType to specific value
func (d *DualType) SetByIndex(val int) {
	d.Index = val
	d.Num = float64(val + d.Offset)
	if d.Index+d.Offset != int(d.Num) {
		log.Fatal("Float to Integer conversion error > 1")
	}
}

// Add add to DualType
func (d *DualType) Add(val int) {
	d.Index += val
	d.Num = float64(d.Index + d.Offset)
	if d.Index+d.Offset != int(d.Num) {
		log.Fatal("Float to Integer conversion error > 1")
	}
}

// NewDualType create a new DualType
func NewDualType(baseIndex int, offset int) DualType {
	return DualType{
		Index:  baseIndex,
		Num:    float64(baseIndex + offset),
		Offset: offset,
	}
}

// GlobalVarsMain contains all variables that are use by multiple sub modules
type GlobalVarsMain struct {
	IZM             int
	DT              DualType
	DZ              DualType
	N               int // max number of layer
	DV              float64
	ALPH            float64
	SATBETA         float64
	AKF             DualType    // current crop index (aktuelle frucht)
	SLNR            int         // Schlag Nummer, plot number
	NFOS            [21]float64 // Nitrogen in fast decomposable fraction (kg N ha-1)
	W, WMIN, PORGES [21]float64
	NAKT            float64
	ETMETH          int
	INIWAHL         int
	DUNGSZEN        float64
	AZHO            int         // number of layer in soil profile (Anzahl Horizonte des Bodenprofils)
	WURZMAX         int         // effective root depth in profile (effektive Wurzeltiefe des Profils)
	DRAIDEP         int         // drainage depth (Tiefe der Drainung)
	DRAIFAK         float64     // part of drainage water in soakage (Anteil des Drainwassers am Sickerwasseranfakk (fraction))
	UKT             [11]int     //(0:10)
	BART            [10]string  // soil type by KA5(Bodenkundlichen Kartieranleitung 5. Auflage) (special spelling convertions)
	LD              [10]int     // bulk density KA5 (1-5) (Lagerungsdichtestufe nach KA5 (1-5))
	BULK            [10]float64 `yaml:"-"` // avarage bulk density (Zuweisung mittlere Lagerungsdichte von LD(I) (g/cm^3))
	CGEHALT         [10]float64 // C organic content in soil layer (Corg-Gehalt in Horizont I (Gew.%))
	HUMUS           [21]float64 // humus content in soil layer (Humusgehalt in Hor. I (Gew.%))
	STEIN           [10]float64 // stone content in soil layer (%)
	FKA             [10]float64 `yaml:"-"` // water content at field capacity (Wassergehalt bei Feldkapazität) (Vol. %)
	WP              [10]float64 `yaml:"-"`
	GPV             [10]float64 `yaml:"-"` // total interstice percentage (Gesamtporenvolumen) (Vol%)
	CAPS            [21]float64
	LIM             [10]float64
	PRGES           [10]float64
	WUMAX           [10]float64 // obsolete
	AD              float64
	GRLO            int
	GRHI            int
	GRW             float64
	GW              float64
	AMPL            int
	PKT             string
	WRED            float64
	PROP            float64
	NORMFK          [10]float64
	FELDW           [10]float64
	CAPPAR          int // not used
	BD              [21]float64
	WNOR            [21]float64
	SAND            [21]float64
	SILT            [21]float64
	CLAY            [21]float64
	NALTOS          float64
	BREG            []float64
	BRKZ            []float64
	ZTBR            []int
	BEGINN          int
	ENDE            int
	FRUCHT          [300]string
	SAAT            [300]int
	JN              [300]float64
	ERNTE           [300]int
	ERTR            [300]float64
	ITAG            int
	TAG             DualType // current day
	JTAG            int      // number of days in current year
	ZTDG            [300]int
	FKU             [12]float64
	CN              [2][21]float64 // (0:1,21) all 21 slots used!
	WG              [3][21]float64 //(0:2,21) all 21 slots used!
	NMESS           int
	MES             [100]string // Should be a local array
	MESS            [100]int
	WNZ             [100]float64
	KNZ1            [100]float64
	KNZ2            [100]float64
	KNZ3            [100]float64
	TILDAT          [200]string
	EINT            [300]float64
	TILART          [200]int
	EINTE           [201]int //(0:200)
	DGART           [300]string
	NDIR            [300]float64
	NSAS            [300]float64
	NLAS            [300]float64
	TSOIL           [2][22]float64 //(0:1,0:21)
	TMIN            [367]float64
	TMAX            [367]float64
	TBASE           float64
	ETNULL          [367]float64
	TEMP            [367]float64
	//TEMPBO1, TEMPBO2 [367]float64 // not initialized, obsolete?
	RH            [367]float64 // relative humidity
	VERD          [367]float64 // Verdunstung, Evaporation, required for ETMETH = 1
	WIND          [367]float64 // wind
	REGEN         [368]float64 // TODO: set to 368 for irrigation calculation (FIXME: last value is 0, load data from next year)
	SUND          [367]float64 // Sun shine hours, required if RAD is 0
	RAD           [367]float64 // photosynthetic active radiation
	ALBEDO        float64
	FEU           int `yaml:"-"`
	C1            [21]float64
	NAOS          [21]float64 // Nitrogen in slowly decomposable pool (kg N ha-1)
	MINAOS        [4]float64  //  should be size of N
	MINFOS        [4]float64  // should be size of N
	CA            [21]float64
	NDG           DualType
	MZ            int
	NBR           int
	NTIL          DualType
	REGENSUM      float64
	MINSUM        float64
	RADSUM        float64
	BLATTSUM      float64
	DSUMM         float64
	UMS           float64
	OUTSUM        float64
	NFIXSUM       float64
	DRAISUM       float64
	DRAINLOSS     float64
	NFIX          float64
	SCHNORR       float64
	PRECO         bool // enable/disable correction factor of rain fall data
	KCOA          float64
	CO2KONZ       float64
	CO2METH       int
	CTRANS        bool
	OUTN          int
	DEPOS         float64
	PROGNOS       int
	FKB           float64
	ANJAHR        int
	J             int
	WINDHI        float64
	ALTI          float64
	SICKER        float64
	CAPSUM        float64
	Q1            [22]float64 //(0:21)
	INTWICK       DualType    // crop development state
	FKF           [12]float64
	FKC           float64
	LAT           float64
	MINTMP        float64
	RSTOM         float64
	LAI           float64 // Leaf area index
	WURZ          int
	VERDUNST      float64
	FLUSS0        float64
	WUDICH        [21]float64
	LUKRIT        [10]float64
	LUMDAY        int
	TP            [21]float64
	TRREL         float64 // Water stress factor (1 = no stress, 0 = full stress)
	REDUK         float64 // Nitrogen stress factor (1 = no stress, 0 = full stress)
	ETA           float64 // Potential/actual Evapotranspiration (mm)
	HEATCOND      [21]float64
	HEATCAP       [21]float64
	TDSUM         [20]float64
	TD            [22]float64 // starts with 0 BBB
	QDRAIN        float64
	TP3, TP6, TP9 float64
	PFTRANS       float64
	INFILT        float64
	ET0           float64
	PE            [21]float64
	MAXAMAX       float64
	WUMAXPF       float64
	WUFKT         int
	NGEFKT        int
	YORGAN        int
	YIFAK         float64
	NRKOM         int
	DAUERKULT     rune
	LEGUM         rune
	DOUBLE        int //day of  double ridge stage / Doppelringstadium
	ASIP          int
	BLUET         int
	REIF          int
	ENDPRO        int // TODO: obsolete?
	PHYLLO        float64
	VERNTAGE      float64
	SUM           [10]float64
	PRO           [10][5]float64
	DEAD          [10][5]float64
	TROOTSUM      float64 `yaml:"-"`
	GEHOB         float64
	WUGEH         float64
	WORG          [5]float64
	WDORG         [10]float64
	MAIRT         [10]float64
	TSUM          [10]float64
	BAS           [10]float64
	VSCHWELL      [10]float64
	DAYL          [10]float64
	DLBAS         [10]float64
	DRYSWELL      [10]float64
	LAIFKT        [10]float64
	WGMAX         [10]float64
	OBMAS         float64
	ASPOO         float64 // Assimilation pool in crops
	WUMAS         float64
	PESUM         float64
	LURED         float64
	DOPP          string // obsolete? Kalender date of double ridge stage / Doppelringstadium
	P1, P2        int
	SUMAE         float64
	AEHR          string
	BLUEH         string
	REIFE         string
	GEHMAX        float64
	GEHMIN        float64
	DUNGBED       float64
	DEFDAT        int
	ENDSTADIUM    DevelopmentStage
	DIFFSUM       float64
	MASSUM        float64
	DN            [21]float64
	YIELD         float64 //Grain yield (only for cereals) (kg ha-1)
	AUFNASUM      float64
	NDRAINTAG     float64 `yaml:"-"`
	CUMDENIT      float64
	AUFNA         [131]float64 //(0:130)
	SIC           [131]float64 //(0:130)
	AUS           [131]float64 // (0:130)
	MINA          float64      // TODO: obsolete
	PLANA         float64
	OUTA          float64 // TODO: obsolete?
	NAPPDAT       string
	PROGDAT       string
	SLNAM         string // not assigned
	// new
	TJBAS     [300]float64
	IRRST1    [300]float64
	IRRST2    [300]float64
	IRRDEP    [300]float64
	IRRLOW    [300]float64
	IRRMAX    [300]float64
	IRRISIM   float64
	TSLWINDOW [300]float64
	TSLMIN    [300]float64
	TSLMAX    [300]float64
	SAAT1     [300]int
	SAAT2     [300]int
	TJAHRSUM  float64
	TJAHR     [300]float64
	MAXMOI    [300]float64
	MINMOI    [300]float64
	ETAG      float64
	SWCS1     float64 // sum of water content for upper 3 layers on sowing date
	SWCS2     float64 // sum of water content for upper 15 layers on sowing date
	SWCA1     float64 // sum of water content for upper 3 layers start of fruit growing
	SWCA2     float64 // sum of water content for upper 15 layers start of fruit growing
	SWCM1     float64 // sum of water content for upper 3 layers on maturity
	SWCM2     float64 // sum of water content for upper 15 layers on maturity
	DRYD1     float64
	DRYD2     float64
	ERNTE2    [300]int
	ETC0      float64
	RDTSUM    float64
	REDSUM    float64
	TRAG      float64
	TRAY      float64
	AUTOMAN   bool // automatic management
	AUTOFERT  bool // automatic fertilization
	AUTOIRRI  bool // automatic irrigation
	AUTOHAR   bool // automatic harvest
	CNRAT1    float64
	PERG      float64
	ETREL     float64
	MAXHMOI   [300]float64
	MINHMOI   [300]float64
	RAINLIM   [300]float64
	RAINACT   [300]float64
	DEV       [10]int // day of year (like sowing, maturity, harvest)
	REDUKSUM  float64
	TRRELSUM  float64
	LAIMAX    float64
	ODU       [300]float64
	NDEM1     [300]float64
	NDEM2     [300]float64
	NDEM3     [300]float64
	ORGDOY    [300]int
	ORGTIME   [300]string
	NDOY1     [300]float64
	NDOY2     [300]float64
	NDOY3     [300]float64
	KNZ4      [100]float64
	KNZ5      [100]float64
	KNZ6      [100]float64
	NLEAG     float64
	NFERTSIM  float64
	AKTUELL   string // current Date string

	// output parameters
	PerY        float64 // accumulated output
	SWCY1       float64 // accumulated output
	SWCY2       float64 // accumulated output
	SOC1        float64 // accumulated output
	Nmin9to20   float64 // sum of C1 from layer 9 to 20
	SickerDaily float64 // sicker - capsum daily update
	HARVEST     float64 // potential harvest daily
	NAOSAKT     float64 // sum NAOS
	NFOSAKT     float64 // sum NFOS
	SumMINAOS   float64 // sum MINAOS
	SumMINFOS   float64 // sum MINFOS
	AvgTSoil    float64 // average TD soil temperature upper 2 layers

	DRflowsum    float64
	Ndrflow      float64
	Nleach       float64
	Percsum      float64
	NfixP        float64
	NAbgbio      float64
	Crop         string
	DevStateDate [10]string

	SNAM            string
	POLYD           string
	FCODE           string
	C1NotStable     string
	C1NotStableErr  string
	C1stabilityVal  float64
	GROUNDWATERFROM GroundWaterFrom        `yaml:"-"`
	DATEFORMAT      DateFormat             `yaml:"-"`
	DEBUGOUT        func(int, interface{}) `yaml:"-"`
	DEBUGCHANNEL    chan<- string          `yaml:"-"`
	LOGID           string                 `yaml:"-"`
	Datum           DateConverterFunc      `yaml:"-"`
	Kalender        KalenderConverterFunc  `yaml:"-"`
	LangTag         LangTagConverterFunc   `yaml:"-"`
}

//CropOutputVars at harvest
type CropOutputVars struct {
	SowDate      string
	SowDOY       int
	EmergDOY     int
	AnthDOY      int
	MatDOY       int
	HarvestDOY   int
	HarvestYear  int
	Crop         string
	Yield        float64
	Biomass      float64
	Roots        float64
	LAImax       float64
	Nfertil      float64
	Irrig        float64
	Nuptake      float64
	Nagb         float64
	ETcG         float64
	ETaG         float64
	TraG         float64
	PerG         float64
	SWCS1        float64
	SWCS2        float64
	SWCA1        float64
	SWCA2        float64
	SWCM1        float64
	SWCM2        float64
	SoilN1       float64
	Nmin1        float64
	Nmin2        float64
	NLeaG        float64
	TRRel        float64
	Reduk        float64
	DryD1        float64
	DryD2        float64
	Nresid       float64
	Orgdat       string
	Type         string
	OrgN         float64
	NDat1        string
	N1           float64
	Ndat2        string
	N2           float64
	Ndat3        string
	N3           float64
	Tdat         string
	Code         string
	NotStableErr string
}

// NewGlobalVarsMain create GlobalVarsMain
func NewGlobalVarsMain() GlobalVarsMain {

	main := GlobalVarsMain{
		TAG:     NewDualType(0, 1),
		INTWICK: NewDualType(-1, 1),
		AKF:     NewDualType(0, 1),
		DT:      NewDualType(1, 0),
		NDG:     NewDualType(0, 1),
		NTIL:    NewDualType(0, 1),
		DZ:      NewDualType(10, 0),
		WINDHI:  2,
		BREG:    make([]float64, 1200),
		BRKZ:    make([]float64, 1200),
		ZTBR:    make([]int, 1200),
		IZM:     30,
		N:       20, // default, will be overwritten by soil
		DV:      4.9,
		// _______ PARAMETER FOR YU/ALLEN _________
		ALPH:           40,
		SATBETA:        2.5,
		C1stabilityVal: -1.5, // Threashold, when becomes negative C1 an error: must be a value below 0
	}
	main.DEBUGOUT = main.printToLimit(100)
	return main
}

func (g *GlobalVarsMain) setIrrigation(zeit, index int, value float64) {
	lenSL := len(g.BREG)
	if index >= lenSL-1 {
		sliceBREG := make([]float64, index+10)
		sliceBRKZ := make([]float64, index+10)
		sliceZTBR := make([]int, index+10)
		copy(sliceBREG, g.BREG)
		copy(sliceBRKZ, g.BRKZ)
		copy(sliceZTBR, g.ZTBR)
		g.BREG = sliceBREG
		g.BRKZ = sliceBRKZ
		g.ZTBR = sliceZTBR
	}
	g.BREG[index] = value
	g.ZTBR[index] = zeit
}