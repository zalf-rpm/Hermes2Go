package main

import (
	"io"
	"net/http"
	"os"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/components"
	"github.com/go-echarts/go-echarts/v2/opts"
)

func GroundwaterDebugHttpServer(w http.ResponseWriter, _ *http.Request) {

	page := components.NewPage()
	keys := extractSortedKeys()
	dates := keysAsDate(Kalender, keys)

	// GW
	// GRW
	// W
	// WMIN
	// PORGES
	// WG
	page.AddCharts(
		lineMultiGW(keys, dates),
		lineMultiWTop(keys, dates),
		lineMultiWSub(keys, dates),
		lineMultiWGTop(keys, dates),
		lineMultiWGSub(keys, dates),
		// lineMultiWMIN(keys, errKeys, dates),
		// lineMultiPORGES(keys, errKeys, dates),
	)

	page.Render(w)
	f, err := os.Create("groundwater_last_run.html")
	if err != nil {
		panic(err)
	}
	page.Render(io.MultiWriter(f))

}

func lineMultiGW(keys []int, dates []string) *charts.Line {

	line := makeMultiLine("GW/GRW")

	line.SetXAxis(dates).
		AddSeries("GW", generateGWItems(keys)).
		AddSeries("GRW", generateGRWItems(keys))
	return line
}

func lineMultiWTop(keys []int, dates []string) *charts.Line {

	line := makeMultiLine("FC Top Layer")

	line.SetXAxis(dates).
		AddSeries("1", generateFCItems(keys, 0)).
		AddSeries("2", generateFCItems(keys, 1)).
		AddSeries("3", generateFCItems(keys, 2)).
		AddSeries("4", generateFCItems(keys, 3)).
		AddSeries("5", generateFCItems(keys, 4)).
		AddSeries("6", generateFCItems(keys, 5)).
		AddSeries("7", generateFCItems(keys, 6)).
		AddSeries("8", generateFCItems(keys, 7)).
		AddSeries("9", generateFCItems(keys, 8)).
		AddSeries("10", generateFCItems(keys, 9))
	return line
}

func lineMultiWSub(keys []int, dates []string) *charts.Line {

	line := makeMultiLine("FC Sub Layer")

	line.SetXAxis(dates).
		AddSeries("11", generateFCItems(keys, 10)).
		AddSeries("12", generateFCItems(keys, 11)).
		AddSeries("13", generateFCItems(keys, 12)).
		AddSeries("14", generateFCItems(keys, 13)).
		AddSeries("15", generateFCItems(keys, 14)).
		AddSeries("16", generateFCItems(keys, 15)).
		AddSeries("17", generateFCItems(keys, 16)).
		AddSeries("18", generateFCItems(keys, 17)).
		AddSeries("19", generateFCItems(keys, 18)).
		AddSeries("20", generateFCItems(keys, 19))

	return line
}

func generateGRWItems(keys []int) []opts.LineData {
	globalHandler.mux.Lock()
	items := make([]opts.LineData, 0, len(keys))

	for _, key := range keys {
		val := globalHandler.receivedDumps[key].Global.GRW
		items = append(items, opts.LineData{Value: val})
	}
	globalHandler.mux.Unlock()
	return items
}

func generateGWItems(keys []int) []opts.LineData {
	globalHandler.mux.Lock()
	items := make([]opts.LineData, 0, len(keys))

	for _, key := range keys {
		val := globalHandler.receivedDumps[key].Global.GW
		items = append(items, opts.LineData{Value: val})
	}
	globalHandler.mux.Unlock()
	return items
}
func lineMultiWGTop(keys []int, dates []string) *charts.Line {
	line := makeMultiLine("WG Top Layer")

	line.SetXAxis(dates).
		AddSeries("1", generateWGItems(keys, 0)).
		AddSeries("2", generateWGItems(keys, 1)).
		AddSeries("3", generateWGItems(keys, 2)).
		AddSeries("4", generateWGItems(keys, 3)).
		AddSeries("5", generateWGItems(keys, 4)).
		AddSeries("6", generateWGItems(keys, 5)).
		AddSeries("7", generateWGItems(keys, 6)).
		AddSeries("8", generateWGItems(keys, 7)).
		AddSeries("9", generateWGItems(keys, 8)).
		AddSeries("10", generateWGItems(keys, 9))
	return line
}

func lineMultiWGSub(keys []int, dates []string) *charts.Line {
	line := makeMultiLine("WG Sub Layer")

	line.SetXAxis(dates).
		AddSeries("11", generateWGItems(keys, 10)).
		AddSeries("12", generateWGItems(keys, 11)).
		AddSeries("13", generateWGItems(keys, 12)).
		AddSeries("14", generateWGItems(keys, 13)).
		AddSeries("15", generateWGItems(keys, 14)).
		AddSeries("16", generateWGItems(keys, 15)).
		AddSeries("17", generateWGItems(keys, 16)).
		AddSeries("18", generateWGItems(keys, 17)).
		AddSeries("19", generateWGItems(keys, 18)).
		AddSeries("20", generateWGItems(keys, 19))
	return line
}

func generateFCItems(keys []int, layer int) []opts.LineData {
	globalHandler.mux.Lock()
	items := make([]opts.LineData, 0, len(keys))

	for _, key := range keys {
		val := globalHandler.receivedDumps[key].Global.W[layer]
		items = append(items, opts.LineData{Value: val})
	}
	globalHandler.mux.Unlock()
	return items
}
func generateWGItems(keys []int, layer int) []opts.LineData {
	globalHandler.mux.Lock()
	items := make([]opts.LineData, 0, len(keys))

	for _, key := range keys {
		val := globalHandler.receivedDumps[key].Global.WG[1][layer]
		items = append(items, opts.LineData{Value: val})
	}
	globalHandler.mux.Unlock()
	return items
}
