package hermes

import (
	"bufio"
	"fmt"
	"log"
	"net/rpc"
	"os"
	"path"
	"strconv"
	"strings"
	"sync"
)

// // HermesFilePool file pool for shared files
// var HermesFilePool FilePool
// var HermesRPCService RPCService
// var HermesOutWriter OutWriterGenerator

// Modfil default module filename
const Modfil = "modinp.txt"

// HFilePath list of hermes file pathes and path template
type HFilePath struct {
	path         string
	locid        string
	parameter    string
	outputfolder string

	config string
	//enam   string // configuration file // daily output for single polygone
	vnam string // daily output for single polygone
	// tnam                   string // output PEST
	// tnnam                  string // other output PEST
	pfnam                  string // output ground temperature
	pnam                   string // output yearly
	mnam                   string // output management
	bofile                 string // soil file e.g soil_<project>.txt
	polnam                 string // polygon file e.g poly_<project>.txt
	irrigation             string // irrigation file
	crop                   string // crop file
	obs                    string // observations file
	til                    string // tillage times file
	dun                    string // fertilization times file
	fert                   string // output fertilization suggestion
	auto                   string // automated processes file
	gwtimeseries           string // groundwater timeseries file
	hypar                  string
	precorr                string
	cropn                  string
	evapo                  string
	parcap                 string
	dung                   string
	vwdatnrm               string
	paranamTemplate        string
	paranamVarietyTemplate string
	bofileTemplate         string
	polnamTemplate         string
	vwdatTemplate          string
	vwdatNoExt             string
	yearlyOutput           string
	dailyOutput            string
	cropOutput             string
	pfOutput               string
	cnam                   string
	managementOutput       string
}

// NewHermesFilePath create an initialized HermesFilePath struct
func NewHermesFilePath(root, locid, uniqueOutputId, parameterOverride, resultOverride string) HFilePath {
	pathToProject := path.Join(root, "project", locid)
	parameter := path.Join(root, "parameter")
	if len(parameterOverride) > 0 {
		parameter = path.Join(root, parameterOverride)
	}
	var out string
	if len(resultOverride) > 0 {
		out = resultOverride
	} else {
		out = path.Join(pathToProject, "RESULT")
	}
	return HFilePath{
		locid:        locid,         // location id, equals project folder name
		path:         pathToProject, // project folder
		parameter:    parameter,     // parameter folder
		outputfolder: out,           // output folder
		// project input files
		irrigation:       path.Join(pathToProject, "irr_"+locid+".txt"),        // irrigation file
		crop:             path.Join(pathToProject, "crop_"+locid+".txt"),       // crop file
		obs:              path.Join(pathToProject, "init_"+locid+".txt"),       // observations file
		til:              path.Join(pathToProject, "til_"+locid+".txt"),        // tillage times file
		dun:              path.Join(pathToProject, "fert_"+locid+".txt"),       // fertilization times file
		auto:             path.Join(pathToProject, "automan.txt"),              // automated processes file
		gwtimeseries:     path.Join(pathToProject, "gw_"+locid+".csv"),         // groundwater timeseries file
		precorr:          path.Join(pathToProject, "Weather", "preco.txt"),     // weather precorrection file
		bofileTemplate:   path.Join(pathToProject, "%s_"+locid+".%s"),          // soil file e.g soil_<project>.txt
		polnamTemplate:   path.Join(pathToProject, "%s_"+locid+".txt"),         // polygon file e.g poly_<project>.txt
		vwdatTemplate:    path.Join(pathToProject, "Weather", "%s_"+locid+"."), // weather data file
		config:           path.Join(pathToProject, "config.yml"),               // configuration file
		yearlyOutput:     path.Join(pathToProject, "yearlyout_conf.yml"),       // yearly output configuration
		dailyOutput:      path.Join(pathToProject, "dailyout_conf.yml"),        // daily output configuration
		cropOutput:       path.Join(pathToProject, "cropout_conf.yml"),         // crop output configuration
		pfOutput:         path.Join(pathToProject, "pfout_conf.yml"),           // harvest output configuration
		managementOutput: path.Join(pathToProject, "managementout_conf.yml"),   // management output configuration
		// output files
		vnam:  path.Join(out, "V"+uniqueOutputId+".%s"),  // daily output template for single polygone
		pfnam: path.Join(out, "P"+uniqueOutputId+".%s"),  // harvest output template
		pnam:  path.Join(out, "Y"+uniqueOutputId+".%s"),  // output yearly template
		cnam:  path.Join(out, "C"+uniqueOutputId+".%s"),  // output crop template
		mnam:  path.Join(out, "M"+uniqueOutputId+".txt"), // management output
		fert:  path.Join(out, "D"+uniqueOutputId+".txt"), // output fertilization suggestion

		// parameter files
		parcap:                 path.Join(parameter, "PARCAP.TRU"),   // capacity parameters
		hypar:                  path.Join(parameter, "HYPAR.TRU"),    // hydraulic parameters
		evapo:                  path.Join(parameter, "EVAPO.HAU"),    // evaporation parameters
		cropn:                  path.Join(parameter, "CROP_N.TXT"),   // crop parameters
		dung:                   path.Join(parameter, "FERTILIZ.TXT"), // fertilization parameters
		paranamTemplate:        path.Join(parameter, "PARAM.%s"),     // parameter file for crop
		paranamVarietyTemplate: path.Join(parameter, "PARAM_%s.%s"),  // parameter file for crop variety

	}
}

// set extension for output files (e.g. res, csv, txt, ...)
func (hp *HFilePath) SetOutputExtension(ext string) {

	hp.vnam = fmt.Sprintf(hp.vnam, ext)
	hp.pfnam = fmt.Sprintf(hp.pfnam, ext)
	hp.cnam = fmt.Sprintf(hp.cnam, ext)
	hp.pnam = fmt.Sprintf(hp.pnam, ext)
}

// SetBofile completes bofile filename
func (hp *HFilePath) SetBofile(prefix, extension string) {
	hp.bofile = fmt.Sprintf(hp.bofileTemplate, strings.TrimSpace(prefix), strings.TrimSpace(extension))
}

// OverrideBofile overrides the complete bofile filename
func (hp *HFilePath) OverrideBofile(newPath string) {
	hp.bofile = newPath
}

// SetPolnam completes polnam filename
func (hp *HFilePath) SetPolnam(ins string) {
	hp.polnam = fmt.Sprintf(hp.polnamTemplate, strings.TrimSpace(ins))
}

// SetVwdatNoExt completes weather data filename, without extension
func (hp *HFilePath) SetVwdatNoExt(ins string) {
	hp.vwdatNoExt = fmt.Sprintf(hp.vwdatTemplate, strings.TrimSpace(ins))
	hp.vwdatnrm = hp.vwdatNoExt + "nrm"
}

func (hp *HFilePath) SetPreCorrFolder(folder string) {
	hp.precorr = path.Join(folder, "preco.txt")
}

// GetParanam returns the full filename for the choosen fruit
func (hp *HFilePath) GetParanam(fruit, variety string, yml bool) string {
	var filename string
	if len(variety) > 0 {
		filename = fmt.Sprintf(hp.paranamVarietyTemplate, variety, fruit)
	} else {
		filename = fmt.Sprintf(hp.paranamTemplate, fruit)
	}
	if yml {
		filename = filename + ".yml"
	}
	return filename
}

// VWdat returns weather data file with the correct extension for a year
func (hp *HFilePath) VWdat(year int) string {

	return hp.vwdatNoExt + yearToExtension(year)
}

func yearToExtension(year int) string {
	var EXTENT string
	strYear := strconv.Itoa(year)
	// IF j >= 100 then
	if year >= 100 {
		//    LET EXTENT$ = "0"&STR$(j)(2:3)
		EXTENT = "0" + strYear[1:3]
		// ELSE
	} else {
		//    LET extent$ = "9"& STR$(j)
		EXTENT = "9" + strYear
		// END IF
	}
	return EXTENT
}

// FilePool to store filepointer by path
type FilePool struct {
	list map[string][]byte
	mux  sync.Mutex
}

// Get file content from the filepool
func (fp *FilePool) Get(fd *FileDescriptior) []byte {
	fp.mux.Lock()
	if fp.list == nil {
		fp.list = make(map[string][]byte)
	}
	if _, ok := fp.list[fd.FilePath]; !ok {
		data, err := os.ReadFile(fd.FilePath)
		if err != nil {
			if fd.FileDescription != "" {
				log.Fatalf("Error occured while reading %s: %s   \n", fd.FileDescription, fd.FilePath)
			} else {

				log.Fatalf("Error occured while reading: %s   \n", fd.FilePath)
			}
		}
		// if fd.debugOut != nil {
		// 	fd.debugOut <- fmt.Sprintf("%s %s", fd.logID, fd.filePath)
		// }
		fp.list[fd.FilePath] = data
	}
	defer fp.mux.Unlock()
	return fp.list[fd.FilePath]
}

// Close clears the filepool
func (fp *FilePool) Close() {
	fp.mux.Lock()
	fp.list = nil
	fp.mux.Unlock()
}

// interface for Fout
type OutWriter interface {
	Write(string) (int, error)
	WriteBytes([]byte) (int, error)
	WriteRune(rune) (int, error)
	Close()
}

type OutWriterGenerator func(string, bool) (OutWriter, error)

func DefaultFoutGenerator(filePath string, append bool) (OutWriter, error) {

	var flags int
	if append {
		flags = os.O_CREATE | os.O_APPEND | os.O_WRONLY
	} else {
		flags = os.O_CREATE | os.O_TRUNC | os.O_WRONLY
	}

	file, err := os.OpenFile(filePath, flags, 0600)
	if err != nil {
		return nil, fmt.Errorf("failed to open result file: %s ", filePath)
	}

	fwriter := bufio.NewWriter(file)
	return &Fout{file, fwriter}, nil
}

// Fout bufferd file writer
type Fout struct {
	file    *os.File
	fwriter *bufio.Writer
}

// Write string to bufferd file
func (f *Fout) Write(s string) (int, error) {
	return f.fwriter.WriteString(s)
}

// WriteBytes writes a bufferd byte array
func (f *Fout) WriteBytes(s []byte) (int, error) {
	return f.fwriter.Write(s)
}

// WriteRune writes a bufferd rune
func (f *Fout) WriteRune(s rune) (int, error) {
	return f.fwriter.WriteRune(s)
}

// Close file writer
func (f *Fout) Close() {
	err := f.fwriter.Flush()
	if err != nil {
		log.Fatalln(err)
	}
	err = f.file.Close()
	if err != nil {
		log.Fatalln(err)
	}
}

type RPCService struct {
	address string // "localhost:8081"
	client  *rpc.Client
}

func NewRPCService(address string) (RPCService, error) {

	client, err := rpc.Dial("tcp", address)
	if err != nil {
		return RPCService{}, err
	}
	return RPCService{address: address, client: client}, nil
}

type TransferEnvGlobal struct {
	Global GlobalVarsMain
	Zeit   int
	Wdt    float64
	Step   int
}

type TransferEnvNitro struct {
	Nitro NitroSharedVars
	Zeit  int
	Wdt   float64
	Step  int
}

type TransferEnvWdt struct {
	Zeit  int
	WDT   float64
	N     int
	WG    [3][21]float64
	W     [21]float64
	DZ    float64
	REGEN float64
}

func (rs *RPCService) SendWdt(g *GlobalVarsMain, zeit int, wdt float64) error {
	if rs.client != nil {
		wdtData := TransferEnvWdt{
			Zeit:  zeit,
			WDT:   wdt,
			N:     g.N,
			WG:    g.WG,
			W:     g.W,
			DZ:    g.DZ.Num,
			REGEN: g.REGEN[g.TAG.Index],
		}
		if err := rs.client.Call("RPCHandler.DumpWdtCalc", wdtData, nil); err != nil {
			return fmt.Errorf("DumpWdtCalc %+v", err)
		}
	}

	return nil
}

func (rs *RPCService) SendGV(g *GlobalVarsMain, zeit int, wdt float64, step int) error {
	if rs.client != nil {
		glob := TransferEnvGlobal{
			Global: *g,
			Zeit:   zeit,
			Wdt:    wdt,
			Step:   step,
		}
		if err := rs.client.Call("RPCHandler.DumpGlobalVar", glob, nil); err != nil {
			return fmt.Errorf("DumpGlobalVar %+v", err)
		}
	}
	return nil
}

func (rs *RPCService) SendNV(n *NitroSharedVars, zeit int, wdt float64, step int) error {
	if rs.client != nil {
		nitro := TransferEnvNitro{
			Nitro: *n,
			Zeit:  zeit,
			Wdt:   wdt,
			Step:  step,
		}

		if err := rs.client.Call("RPCHandler.DumpNitroVar", nitro, nil); err != nil {
			return fmt.Errorf("DumpNitroVar %+v", err)
		}
	}
	return nil
}
