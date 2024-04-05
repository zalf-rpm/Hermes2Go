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
		lineMultiGW_WaterContentTop(keys, dates),
		lineMultiWTop(keys, dates),
		lineMultiWSub(keys, dates),
		// lineMultiWGTop(keys, dates),
		// lineMultiWGSub(keys, dates),
		// lineMultiWGinPercentTop(keys, dates),
		// lineMultiWGinPercentSub(keys, dates),
	)

	page.Render(w)
	f, err := os.Create("groundwater_last_run.html")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	page.Render(io.MultiWriter(f))

}

func lineMultiGW(keys []int, dates []string) *charts.Line {

	line := makeMultiLine("GW/GRW")

	line.SetXAxis(dates).
		AddSeries("GW", generateGWItems(keys)).
		AddSeries("GRW", generateGRWItems(keys, 1, 99))
	return line
}

func lineMultiGW_WaterContentTop(keys []int, dates []string) *charts.Line {

	line := makeMultiLine("GRW & WG in layer")

	line.SetXAxis(dates).
		AddSeries("GRW", generateGRWItems(keys, 100, 2100)).
		AddSeries("1", generateWGinPercOfFCItems(keys, 0, 0)).
		AddSeries("2", generateWGinPercOfFCItems(keys, 1, 100)).
		AddSeries("3", generateWGinPercOfFCItems(keys, 2, 200)).
		AddSeries("4", generateWGinPercOfFCItems(keys, 3, 300)).
		AddSeries("5", generateWGinPercOfFCItems(keys, 4, 400)).
		AddSeries("6", generateWGinPercOfFCItems(keys, 5, 500)).
		AddSeries("7", generateWGinPercOfFCItems(keys, 6, 600)).
		AddSeries("8", generateWGinPercOfFCItems(keys, 7, 700)).
		AddSeries("9", generateWGinPercOfFCItems(keys, 8, 800)).
		AddSeries("10", generateWGinPercOfFCItems(keys, 9, 900)).
		AddSeries("11", generateWGinPercOfFCItems(keys, 10, 1000)).
		AddSeries("12", generateWGinPercOfFCItems(keys, 11, 1100)).
		AddSeries("13", generateWGinPercOfFCItems(keys, 12, 1200)).
		AddSeries("14", generateWGinPercOfFCItems(keys, 13, 1300)).
		AddSeries("15", generateWGinPercOfFCItems(keys, 14, 1400)).
		AddSeries("16", generateWGinPercOfFCItems(keys, 15, 1500)).
		AddSeries("17", generateWGinPercOfFCItems(keys, 16, 1600)).
		AddSeries("18", generateWGinPercOfFCItems(keys, 17, 1700)).
		AddSeries("19", generateWGinPercOfFCItems(keys, 18, 1800)).
		AddSeries("20", generateWGinPercOfFCItems(keys, 19, 1900))
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

func generateGRWItems(keys []int, offset, cap float64) []opts.LineData {
	globalHandler.mux.Lock()
	items := make([]opts.LineData, 0, len(keys))

	for _, key := range keys {
		val := globalHandler.receivedDumps[key].Global.GRW * offset
		if val > cap {
			val = cap
		}
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

func lineMultiWGinPercentTop(keys []int, dates []string) *charts.Line {
	line := makeMultiLine("WG in % Top Layer")

	line.SetXAxis(dates).
		AddSeries("1", generateWGinPercOfFCItems(keys, 0, 0)).
		AddSeries("2", generateWGinPercOfFCItems(keys, 1, 0)).
		AddSeries("3", generateWGinPercOfFCItems(keys, 2, 0)).
		AddSeries("4", generateWGinPercOfFCItems(keys, 3, 0)).
		AddSeries("5", generateWGinPercOfFCItems(keys, 4, 0)).
		AddSeries("6", generateWGinPercOfFCItems(keys, 5, 0)).
		AddSeries("7", generateWGinPercOfFCItems(keys, 6, 0)).
		AddSeries("8", generateWGinPercOfFCItems(keys, 7, 0)).
		AddSeries("9", generateWGinPercOfFCItems(keys, 8, 0)).
		AddSeries("10", generateWGinPercOfFCItems(keys, 9, 0))
	return line
}

func lineMultiWGinPercentSub(keys []int, dates []string) *charts.Line {
	line := makeMultiLine("WG in % Sub Layer")

	line.SetXAxis(dates).
		AddSeries("11", generateWGinPercOfFCItems(keys, 10, 0)).
		AddSeries("12", generateWGinPercOfFCItems(keys, 11, 0)).
		AddSeries("13", generateWGinPercOfFCItems(keys, 12, 0)).
		AddSeries("14", generateWGinPercOfFCItems(keys, 13, 0)).
		AddSeries("15", generateWGinPercOfFCItems(keys, 14, 0)).
		AddSeries("16", generateWGinPercOfFCItems(keys, 15, 0)).
		AddSeries("17", generateWGinPercOfFCItems(keys, 16, 0)).
		AddSeries("18", generateWGinPercOfFCItems(keys, 17, 0)).
		AddSeries("19", generateWGinPercOfFCItems(keys, 18, 0)).
		AddSeries("20", generateWGinPercOfFCItems(keys, 19, 0))
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
func generateWGinPercOfFCItems(keys []int, layer int, offset float64) []opts.LineData {
	globalHandler.mux.Lock()
	items := make([]opts.LineData, 0, len(keys))

	for _, key := range keys {
		val := globalHandler.receivedDumps[key].Global.WG[1][layer]
		val = (val / globalHandler.receivedDumps[key].Global.W[layer] * 100) + offset
		items = append(items, opts.LineData{Value: val})
	}
	globalHandler.mux.Unlock()
	return items
}
