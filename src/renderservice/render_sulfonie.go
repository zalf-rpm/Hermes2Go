package main

import (
	"io"
	"net/http"
	"os"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/components"
	"github.com/go-echarts/go-echarts/v2/opts"
)

func sulfoniehttpserver(w http.ResponseWriter, _ *http.Request) {

	page := components.NewPage()
	keys := extractSortedKeys()
	dates := keysAsDate(Kalender, keys)
	errKeys := generateErrorItems(keys)

	// SGEHOB
	// SREDUK
	// SGEHMAX
	// SGEHMIN
	// S1
	page.AddCharts(
		lineMultiS1(keys, errKeys, dates),
		lineMultiSGEHMIN(keys, errKeys, dates),
		lineMultiSREDUK(keys, errKeys, dates),
		lineMultiSGEHOB(keys, errKeys, dates),
	)

	page.Render(w)
	f, err := os.Create("sulfonie_last_run.html")
	if err != nil {
		panic(err)
	}
	page.Render(io.MultiWriter(f))

}

func lineMultiS1(keys, errKeys []int, dates []string) *charts.Line {
	line := makeMultiLine("Smin Content")

	line.SetXAxis(dates).
		AddSeries("S1 Schicht 1", generateS1Items(keys, 0), errorMarker(errKeys, 60)).
		AddSeries("S1 Schicht 2", generateS1Items(keys, 1)).
		AddSeries("S1 Schicht 3", generateS1Items(keys, 2))
	return line
}

// generateC1Items generate Nmin (C1) items
func generateS1Items(keys []int, index int) []opts.LineData {
	globalHandler.mux.Lock()
	items := make([]opts.LineData, 0, len(keys))

	for _, key := range keys {
		val := globalHandler.receivedDumps[key].Global.S1[index]
		items = append(items, opts.LineData{Value: val})
	}
	globalHandler.mux.Unlock()
	return items
}

func lineMultiSGEHMIN(keys, errKeys []int, dates []string) *charts.Line {

	line := makeMultiLine("S GEHALT MIN/MAX")

	line.SetXAxis(dates).
		AddSeries("SGEHMIN", generateSGEHMINItems(keys), errorMarker(errKeys, 23)).
		AddSeries("SGEHMAX", generateSGEHMAXItems(keys))
	return line
}

func generateSGEHMINItems(keys []int) []opts.LineData {
	globalHandler.mux.Lock()
	items := make([]opts.LineData, 0, len(keys))

	for _, key := range keys {
		val := globalHandler.receivedDumps[key].Global.SGEHMIN
		items = append(items, opts.LineData{Value: val})
	}
	globalHandler.mux.Unlock()
	return items
}

func generateSGEHMAXItems(keys []int) []opts.LineData {
	globalHandler.mux.Lock()
	items := make([]opts.LineData, 0, len(keys))

	for _, key := range keys {
		val := globalHandler.receivedDumps[key].Global.SGEHMAX
		items = append(items, opts.LineData{Value: val})
	}
	globalHandler.mux.Unlock()
	return items
}

func lineMultiSREDUK(keys, errKeys []int, dates []string) *charts.Line {

	line := makeMultiLine("S REDUKTION")

	line.SetXAxis(dates).
		AddSeries("SREDUK", generateSREDUKItems(keys), errorMarker(errKeys, 23))
	return line
}

func generateSREDUKItems(keys []int) []opts.LineData {
	globalHandler.mux.Lock()
	items := make([]opts.LineData, 0, len(keys))

	for _, key := range keys {
		val := globalHandler.receivedDumps[key].Global.SREDUK
		items = append(items, opts.LineData{Value: val})
	}
	globalHandler.mux.Unlock()
	return items
}

func lineMultiSGEHOB(keys, errKeys []int, dates []string) *charts.Line {
	line := makeMultiLine("S GEHALT OBEN")

	line.SetXAxis(dates).
		AddSeries("SGEHOB", generateSGEHOBItems(keys), errorMarker(errKeys, 23))
	return line
}

func generateSGEHOBItems(keys []int) []opts.LineData {
	globalHandler.mux.Lock()
	items := make([]opts.LineData, 0, len(keys))

	for _, key := range keys {
		val := globalHandler.receivedDumps[key].Global.SGEHOB
		items = append(items, opts.LineData{Value: val})
	}
	globalHandler.mux.Unlock()
	return items
}
