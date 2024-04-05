package main

import (
	"fmt"
	"io"
	"log"
	"math"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"sort"
	"sync"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/components"
	"github.com/go-echarts/go-echarts/v2/opts"

	"github.com/zalf-rpm/Hermes2Go/hermes"
)

// this is a quick and dirty debug service to render the content of hermes variable dumps on a timeline.
// It is an RPC Service. It has to be started before a hermes run
// Hermes needs to be started with adding "-rpc localhost:8082" to command line
// it works only for one setup
// and it is not stable

func main() {
	StartRPCHandler()
	Kalender = hermes.KalenderConverter(hermes.DateDEshort, ".")
	fmt.Println("Open a browser and connet to: http://localhost:8081 ")
	fmt.Println("For N2O connet to: http://localhost:8081/n2o ")
	fmt.Println("For ground water connet to: http://localhost:8081/groundwater ")
	fmt.Println("For plant connet to: http://localhost:8081/crop ")
	http.HandleFunc("/", c1debughttpserver)
	http.HandleFunc("/n2o", n2odebughttpserver)
	http.HandleFunc("/groundwater", GroundwaterDebugHttpServer)
	http.HandleFunc("/crop", cropdebughttpserver)
	http.ListenAndServe("localhost:8081", nil)
}

// RPCHandler for receiving and storing data from a run
type RPCHandler struct {
	receivedDumps      map[int]hermes.TransferEnvGlobal // GlobalVarsMain - global vars
	receivedNitroDumps map[int]hermes.TransferEnvNitro  // NitroSharedVars - local Nitro vars
	wdtDumps           map[int]hermes.TransferEnvWdt    // wdt calc dump
	mux                sync.Mutex
}

var globalHandler *RPCHandler
var Kalender hermes.KalenderConverterFunc

// StartRPCHandler start RPC handler in another go routine
func StartRPCHandler() {

	l, err := net.Listen("tcp", "localhost:8082")
	if err != nil {
		log.Fatalf("Error while starting rpc server: %+v", err)
	}

	globalHandler = &RPCHandler{
		receivedDumps:      make(map[int]hermes.TransferEnvGlobal, 11000),
		receivedNitroDumps: make(map[int]hermes.TransferEnvNitro, 11000),
		wdtDumps:           make(map[int]hermes.TransferEnvWdt, 11000),
		mux:                sync.Mutex{},
	}
	err = rpc.Register(globalHandler)
	if err != nil {
		panic(err)
	}

	go func() {
		for {
			c, err := l.Accept()
			fmt.Printf("request from %v\n", c.RemoteAddr())
			if err != nil {
				continue
			}
			go rpc.ServeConn(c)
		}
	}()
}

// DumpGlobalVar handle global Var Dump
func (rh *RPCHandler) DumpGlobalVar(payload hermes.TransferEnvGlobal, reply *string) error {

	id := payload.Zeit*100 + payload.Step
	rh.mux.Lock()
	rh.receivedDumps[id] = payload
	rh.mux.Unlock()
	return nil
}

// DumpNitroVar handle Nitro Local var dump
func (rh *RPCHandler) DumpNitroVar(payload hermes.TransferEnvNitro, reply *string) error {

	id := payload.Zeit*100 + payload.Step
	rh.mux.Lock()
	rh.receivedNitroDumps[id] = payload
	rh.mux.Unlock()
	return nil
}
func (rh *RPCHandler) DumpWdtCalc(payload hermes.TransferEnvWdt, reply *string) error {
	id := payload.Zeit
	rh.mux.Lock()
	rh.wdtDumps[id] = payload
	rh.mux.Unlock()
	return nil

}

// c1debughttpserver for web interface, to render the stuff that has been recieved
func c1debughttpserver(w http.ResponseWriter, _ *http.Request) {

	page := components.NewPage()
	keys := extractSortedKeys()
	dates := keysAsDate(Kalender, keys)
	errKeys := generateErrorItems(keys)

	page.AddCharts(
		lineMultiC1(keys, errKeys, dates),
		lineMultiQ1(keys, errKeys, dates),
		lineMultiDISP(keys, errKeys, dates),
		lineMultiKONV(keys, errKeys, dates),
		lineMultiDB(keys, errKeys, dates),
		lineMultiV(keys, errKeys, dates),
		lineMultiWDTCalc(keys, errKeys, dates),
		lineMultiWDT(),
	)

	page.Render(w)
	f, err := os.Create("last_run.html")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	page.Render(io.MultiWriter(f))

}

// generateC1Items generate Nmin (C1) items
func generateC1Items(keys []int, index int) []opts.LineData {
	globalHandler.mux.Lock()
	items := make([]opts.LineData, 0, len(keys))

	for _, key := range keys {
		val := globalHandler.receivedDumps[key].Global.C1[index]
		items = append(items, opts.LineData{Value: val})
	}
	globalHandler.mux.Unlock()
	return items
}

func generateQ1Items(keys []int, index int) []opts.LineData {
	globalHandler.mux.Lock()
	items := make([]opts.LineData, 0, len(keys))

	for _, key := range keys {
		val := globalHandler.receivedDumps[key].Global.Q1[index]
		items = append(items, opts.LineData{Value: val})
	}
	globalHandler.mux.Unlock()
	return items
}

func generateRegenItems(keys []int) []opts.LineData {
	globalHandler.mux.Lock()
	items := make([]opts.LineData, 0, len(keys))

	for _, key := range keys {
		tag := globalHandler.receivedDumps[key].Global.TAG.Index
		val := globalHandler.receivedDumps[key].Global.REGEN[tag]
		items = append(items, opts.LineData{Value: val})
	}
	globalHandler.mux.Unlock()
	return items
}

func generateDISPItems(keys []int, index int) []opts.LineData {
	globalHandler.mux.Lock()
	items := make([]opts.LineData, 0, len(keys))

	for _, key := range keys {
		val := globalHandler.receivedNitroDumps[key].Nitro.DISP[index]
		items = append(items, opts.LineData{Value: val})
	}
	globalHandler.mux.Unlock()
	return items
}

func generateKONVItems(keys []int, index int) []opts.LineData {
	globalHandler.mux.Lock()
	items := make([]opts.LineData, 0, len(keys))

	for _, key := range keys {
		val := globalHandler.receivedNitroDumps[key].Nitro.KONV[index]
		items = append(items, opts.LineData{Value: val})
	}
	globalHandler.mux.Unlock()
	return items
}

func generateVItems(keys []int, index int) []opts.LineData {
	globalHandler.mux.Lock()
	items := make([]opts.LineData, 0, len(keys))

	for _, key := range keys {
		val := globalHandler.receivedNitroDumps[key].Nitro.V[index]
		items = append(items, opts.LineData{Value: val})
	}
	globalHandler.mux.Unlock()
	return items
}

func generateWdtVisItems(keys []int) map[string][]opts.LineData {
	globalHandler.mux.Lock()

	items := map[string][]opts.LineData{}

	items["regen"] = make([]opts.LineData, 0, len(keys))
	items["wdt calc"] = make([]opts.LineData, 0, len(keys))
	items["W"] = make([]opts.LineData, 0, len(keys))
	items["WG1"] = make([]opts.LineData, 0, len(keys))
	items["WG2"] = make([]opts.LineData, 0, len(keys))
	items["WG3"] = make([]opts.LineData, 0, len(keys))
	items["FSCSUM1"] = make([]opts.LineData, 0, len(keys))
	items["FSCSUM2"] = make([]opts.LineData, 0, len(keys))
	items["FSCSUM3"] = make([]opts.LineData, 0, len(keys))

	var FSCSUM [20]float64
	for _, key := range keys {

		N := globalHandler.wdtDumps[key].N
		WG := globalHandler.wdtDumps[key].WG
		W := globalHandler.wdtDumps[key].W
		DZ := globalHandler.wdtDumps[key].DZ
		REGEN := globalHandler.wdtDumps[key].REGEN
		FSCS := 0.0
		ZSR := 1.0
		//WDT = g.DT.Num
		for I := 1; I <= N; I++ {
			index := I - 1
			FSC := (W[index] - WG[0][index]) * DZ
			FSCS = FSCS + FSC
			FSCSUM[index] = FSCS
		}

		for I := 1; I <= N; I++ {
			index := I - 1

			if REGEN-FSCSUM[index] > W[index]*DZ/3 {
				ZSR = math.Max(ZSR, (REGEN-FSCSUM[index])/(W[index]*DZ/3))
			}
		}
		calcWDT := 1 / math.Ceil(ZSR)
		items["wdt calc"] = append(items["wdt calc"], opts.LineData{Value: calcWDT})
		items["W"] = append(items["W"], opts.LineData{Value: W[0] * DZ})
		items["WG1"] = append(items["WG1"], opts.LineData{Value: WG[1][0] * DZ})
		items["WG2"] = append(items["WG2"], opts.LineData{Value: WG[1][1] * DZ})
		items["WG3"] = append(items["WG3"], opts.LineData{Value: WG[1][2] * DZ})
		items["FSCSUM1"] = append(items["FSCSUM1"], opts.LineData{Value: FSCSUM[0]})
		items["FSCSUM2"] = append(items["FSCSUM2"], opts.LineData{Value: FSCSUM[1]})
		items["FSCSUM3"] = append(items["FSCSUM3"], opts.LineData{Value: FSCSUM[2]})
		items["regen"] = append(items["regen"], opts.LineData{Value: REGEN})
	}

	globalHandler.mux.Unlock()
	return items
}

func generateWdtCalcItems(keys []int) [][]opts.LineData {
	globalHandler.mux.Lock()
	numLines := 4
	items := make([][]opts.LineData, 0, numLines)
	for i := 0; i < numLines; i++ {
		items = append(items, make([]opts.LineData, 0, len(keys)))
	}

	for _, key := range keys {
		g := globalHandler.receivedDumps[key].Global
		items[0] = append(items[0], opts.LineData{Value: g.REGEN[g.TAG.Index]})
		items[1] = append(items[1], opts.LineData{Value: globalHandler.receivedDumps[key].Wdt})

		// try a test with Monica variante and Fluss0
		pri := g.FLUSS0 * g.DZ.Num
		items[2] = append(items[2], opts.LineData{Value: g.FLUSS0})
		wdt := 1.0
		timeStepFactorCurrentLayer := 1.0
		if -5.0 <= pri && pri <= 5.0 && wdt > 1.0 {
			timeStepFactorCurrentLayer = 1.0
		} else if (-10.0 <= pri && pri < -5.0) || (5.0 < pri && pri <= 10.0) {
			timeStepFactorCurrentLayer = 0.5
		} else if (-15.0 <= pri && pri < -10.0) || (10.0 < pri && pri <= 15.0) {
			timeStepFactorCurrentLayer = 0.25
		} else if pri < -15.0 || pri > 15.0 {
			timeStepFactorCurrentLayer = 0.125
		}
		wdt = math.Min(wdt, timeStepFactorCurrentLayer)
		items[3] = append(items[3], opts.LineData{Value: wdt})

	}
	globalHandler.mux.Unlock()
	return items
}

func generateDBItems(keys []int, index int) []opts.LineData {
	globalHandler.mux.Lock()
	items := make([]opts.LineData, 0, len(keys))

	for _, key := range keys {
		val := globalHandler.receivedNitroDumps[key].Nitro.DB[index]
		items = append(items, opts.LineData{Value: val})
	}
	globalHandler.mux.Unlock()
	return items
}

func extractSortedKeys() []int {
	globalHandler.mux.Lock()
	keys := make([]int, 0, len(globalHandler.receivedDumps))
	for k := range globalHandler.receivedDumps {
		keys = append(keys, k)
	}
	globalHandler.mux.Unlock()
	sort.Ints(keys)
	return keys
}

func makeMultiLine(title string) *charts.Line {
	line := charts.NewLine()
	line.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title: title,
		}),
		charts.WithInitializationOpts(opts.Initialization{
			Theme: "shine",
		}),
		charts.WithLegendOpts(opts.Legend{Show: true}),
		charts.WithDataZoomOpts(opts.DataZoom{
			Type:       "inside",
			Start:      50,
			End:        100,
			XAxisIndex: []int{0},
		}),
		charts.WithDataZoomOpts(opts.DataZoom{
			Type:       "slider",
			Start:      50,
			End:        100,
			XAxisIndex: []int{0},
		}),
	)
	return line
}

func lineMultiC1(keys, errKeys []int, dates []string) *charts.Line {
	line := makeMultiLine("Nmin Content")

	line.SetXAxis(dates).
		AddSeries("C1 Schicht 1", generateC1Items(keys, 0), errorMarker(errKeys, 60)).
		AddSeries("C1 Schicht 2", generateC1Items(keys, 1)).
		AddSeries("C1 Schicht 3", generateC1Items(keys, 2))
	return line
}
func errorMarker(errKeys []int, offset float64) charts.SeriesOpts {
	dates := keysAsDate(Kalender, errKeys)

	marker := make([]opts.MarkPointNameCoordItem, 0, len(dates))
	for _, date := range dates {
		marker = append(marker, opts.MarkPointNameCoordItem{
			Name:       "err",
			Coordinate: []interface{}{date, offset},
			Label:      &opts.Label{Show: true, Color: "white", Position: "inside", Formatter: "{b}"},
		})
	}
	options := charts.WithMarkPointNameCoordItemOpts(marker...)

	return options
}

func generateErrorItems(keys []int) []int {
	globalHandler.mux.Lock()
	listOfErrors := []int{}
	for _, key := range keys {
		tag := globalHandler.receivedDumps[key].Global.C1NotStable
		if len(tag) > 0 {
			listOfErrors = append(listOfErrors, key)
		}
	}
	globalHandler.mux.Unlock()
	return listOfErrors
}

func lineMultiQ1(keys, errKeys []int, dates []string) *charts.Line {
	line := makeMultiLine("Q1 Fluss durch Untergrenze")

	line.SetXAxis(dates)
	line.AddSeries("Regen ", generateRegenItems(keys))
	line.AddSeries("Q1 Schicht 0", generateQ1Items(keys, 0), errorMarker(errKeys, 9))
	line.AddSeries("Q1 Schicht 1", generateQ1Items(keys, 1))
	line.AddSeries("Q1 Schicht 2", generateQ1Items(keys, 2))
	line.AddSeries("Q1 Schicht 3", generateQ1Items(keys, 3))

	return line
}

func lineMultiDISP(keys, errKeys []int, dates []string) *charts.Line {
	line := makeMultiLine("Dispersion")

	line.SetXAxis(dates).
		AddSeries("DISP 1", generateDISPItems(keys, 0), errorMarker(errKeys, 0.035)).
		AddSeries("DISP 2", generateDISPItems(keys, 1)).
		AddSeries("DISP 3", generateDISPItems(keys, 2))
	return line
}

func lineMultiKONV(keys, errKeys []int, dates []string) *charts.Line {
	line := makeMultiLine("Konvection")

	line.SetXAxis(dates).
		AddSeries("KONV 1", generateKONVItems(keys, 0), errorMarker(errKeys, 0.025)).
		AddSeries("KONV 2", generateKONVItems(keys, 1)).
		AddSeries("KONV 3", generateKONVItems(keys, 2))
	return line
}

func lineMultiDB(keys, errKeys []int, dates []string) *charts.Line {
	line := makeMultiLine("DB")

	line.SetXAxis(dates).
		AddSeries("DB 1", generateDBItems(keys, 0), errorMarker(errKeys, 23)).
		AddSeries("DB 2", generateDBItems(keys, 1)).
		AddSeries("DB 3", generateDBItems(keys, 2))
	return line
}

func lineMultiV(keys, errKeys []int, dates []string) *charts.Line {
	line := makeMultiLine("Porenwassergeschwindigkeit")

	line.SetXAxis(dates).
		AddSeries("V 1", generateVItems(keys, 0), errorMarker(errKeys, 1)).
		AddSeries("V 2", generateVItems(keys, 1)).
		AddSeries("V 3", generateVItems(keys, 2))
	return line
}

func lineMultiWDTCalc(keys, errKeys []int, dates []string) *charts.Line {
	line := makeMultiLine("Zeitschritt Berechnung")

	linesContent := generateWdtCalcItems(keys)

	line.SetXAxis(dates)
	line.AddSeries("regen", linesContent[0])
	line.AddSeries("wdt", linesContent[1], errorMarker(errKeys, 1))
	line.AddSeries("Fluss0", linesContent[2])
	line.AddSeries("wdt Fluss0/Monica", linesContent[3])
	//AddSeries("Q1 Schicht 0", generateQ1Items(keys, 0))

	return line
}
func lineMultiWDT() *charts.Line {
	line := makeMultiLine("Zeitschritt Berechnung Details")

	globalHandler.mux.Lock()
	keys := make([]int, 0, len(globalHandler.receivedDumps))
	for k := range globalHandler.wdtDumps {
		keys = append(keys, k)
	}
	globalHandler.mux.Unlock()
	sort.Ints(keys)

	dates := make([]string, 0, len(keys))
	for _, key := range keys {
		dates = append(dates, Kalender(key))
	}
	linesContent := generateWdtVisItems(keys)

	line.SetXAxis(dates)
	for entry, value := range linesContent {
		line.AddSeries(entry, value)
	}

	return line
}

func keysAsDate(dateConverter func(int) string, keys []int) []string {
	asDate := make([]string, 0, len(keys))

	for _, key := range keys {
		date := key / 100
		steps := key % 100

		asDate = append(asDate, fmt.Sprintf("%s_%d", dateConverter(date), steps))
	}
	return asDate
}
