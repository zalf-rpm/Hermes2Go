package hermes

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net/rpc"
	"os"
	"path"
	"strconv"
	"strings"
	"sync"
)

// HermesFilePool file pool for shared files
var HermesFilePool FilePool
var HermesRPCService RPCService

// Modfil default module filename
const Modfil = "modinp.txt"

// HFilePath list of hermes file pathes and path template
type HFilePath struct {
	path         string
	locid        string
	parameter    string
	outputfolder string

	config                 string
	enam                   string // configuration file // daily output for single polygone
	vnam                   string // daily output for single polygone
	tnam                   string // output PEST
	tnnam                  string // other output PEST
	pfnam                  string // output ground temperature
	pnam                   string // output yearly
	bofile                 string // soil file e.g soil_<project>.txt
	polnam                 string // polygon file e.g poly_<project>.txt
	irrigation             string // irrigation file
	crop                   string // crop file
	obs                    string // observations file
	til                    string // tillage times file
	dun                    string // fertilization times file
	fert                   string // output fertilization suggestion
	auto                   string // automated processes file
	hypar                  string
	precorr                string
	cropn                  string
	evapo                  string
	parcap                 string
	dung                   string
	vwdatnrm               string
	pnamTemplate           string
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
}

// NewHermesFilePath create an initialized HermesFilePath struct
func NewHermesFilePath(root, locid, snam, parameterOverride, resultOverride string) HFilePath {
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
		locid:                  locid,
		path:                   pathToProject,
		parameter:              parameter,
		outputfolder:           out,
		enam:                   path.Join(pathToProject, locid+".dri"),
		vnam:                   path.Join(out, "v"+snam+".res"),
		tnam:                   path.Join(out, "t"+snam+".res"),
		tnnam:                  path.Join(out, "n"+snam+".res"),
		pfnam:                  path.Join(out, "p"+snam+".res"),
		fert:                   path.Join(out, "d_"+snam+".txt"),
		irrigation:             path.Join(pathToProject, "irr_"+locid+".txt"),
		crop:                   path.Join(pathToProject, "crop_"+locid+".txt"),
		obs:                    path.Join(pathToProject, "init_"+locid+".txt"),
		til:                    path.Join(pathToProject, "til_"+locid+".txt"),
		dun:                    path.Join(pathToProject, "fert_"+locid+".txt"),
		auto:                   path.Join(pathToProject, "automan.txt"),
		precorr:                path.Join(pathToProject, "Weather", "preco.txt"),
		parcap:                 path.Join(parameter, "PARCAP.TRU"),
		hypar:                  path.Join(parameter, "HYPAR.TRU"),
		evapo:                  path.Join(parameter, "EVAPO.HAU"),
		cropn:                  path.Join(parameter, "CROP_N.TXT"),
		dung:                   path.Join(parameter, "FERTILIZ.TXT"),
		pnamTemplate:           path.Join(out, "%s.%s"),
		paranamTemplate:        path.Join(parameter, "PARAM.%s"),
		paranamVarietyTemplate: path.Join(parameter, "PARAM_%s.%s"),
		bofileTemplate:         path.Join(pathToProject, "%s_"+locid+".%s"),
		polnamTemplate:         path.Join(pathToProject, "%s_"+locid+".txt"),
		vwdatTemplate:          path.Join(pathToProject, "Weather", "%s_"+locid+"."),
		config:                 path.Join(pathToProject, "config.yml"),
		yearlyOutput:           path.Join(pathToProject, "yearlyout_conf.yml"),
		dailyOutput:            path.Join(pathToProject, "dailyout_conf.yml"),
		cropOutput:             path.Join(pathToProject, "cropout_conf.yml"),
		pfOutput:               path.Join(pathToProject, "pfout_conf.yml"),
	}
}

// SetPnam completes pnam filename
func (hp *HFilePath) SetPnam(ins, ext string) {
	hp.pnam = fmt.Sprintf(hp.pnamTemplate, strings.TrimSpace(ins), ext)
}

// SetBofile completes bofile filename
func (hp *HFilePath) SetBofile(prefix, extension string) {
	hp.bofile = fmt.Sprintf(hp.bofileTemplate, strings.TrimSpace(prefix), strings.TrimSpace(extension))
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

// GetParanam returns the full filename for the choosen fruit
func (hp *HFilePath) GetParanam(fruit, variety string) string {
	if len(variety) > 0 {
		return fmt.Sprintf(hp.paranamVarietyTemplate, variety, fruit)
	}
	return fmt.Sprintf(hp.paranamTemplate, fruit)
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
		data, err := ioutil.ReadFile(fd.FilePath)
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

//Fout bufferd file writer
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
