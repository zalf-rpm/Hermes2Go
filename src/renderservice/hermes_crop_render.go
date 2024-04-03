package main

import (
	"io"
	"net/http"
	"os"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/components"
	"github.com/go-echarts/go-echarts/v2/opts"
)

func cropdebughttpserver(w http.ResponseWriter, _ *http.Request) {

	page := components.NewPage()
	keys := extractSortedKeys()
	dates := keysAsDate(Kalender, keys)
	errKeys := generateErrorItems(keys)

	// GEHMIN
	// GEHMAX
	// PE
	page.AddCharts(
		lineMultiGEHALT(keys, errKeys, dates),
	)

	page.Render(w)
	f, err := os.Create("crop_last_run.html")
	if err != nil {
		panic(err)
	}
	page.Render(io.MultiWriter(f))

}

func lineMultiGEHALT(keys, errKeys []int, dates []string) *charts.Line {

	line := makeMultiLine("GEHALT")

	line.SetXAxis(dates).
		AddSeries("GEHMIN", generateGehMinItems(keys), errorMarker(errKeys, 23)).
		AddSeries("GEHMAX", generateGehMaxItems(keys)).
		AddSeries("WUGEH", generateWuGehItems(keys)).
		AddSeries("GEHOB", generateGehobItems(keys)).
		AddSeries("PE", generatePEItems(keys))
	return line
}

// max Gehalt
func generateGehMaxItems(keys []int) []opts.LineData {
	globalHandler.mux.Lock()
	items := make([]opts.LineData, 0, len(keys))

	for _, key := range keys {
		val := globalHandler.receivedDumps[key].Global.GEHMAX
		items = append(items, opts.LineData{Value: val})
	}
	globalHandler.mux.Unlock()
	return items
}

// gehalt wurzel
func generateWuGehItems(keys []int) []opts.LineData {
	globalHandler.mux.Lock()
	items := make([]opts.LineData, 0, len(keys))

	for _, key := range keys {
		val := globalHandler.receivedDumps[key].Global.WUGEH
		items = append(items, opts.LineData{Value: val})
	}
	globalHandler.mux.Unlock()
	return items
}

// gehalt obermasse
func generateGehobItems(keys []int) []opts.LineData {
	globalHandler.mux.Lock()
	items := make([]opts.LineData, 0, len(keys))

	for _, key := range keys {
		val := globalHandler.receivedDumps[key].Global.GEHOB
		items = append(items, opts.LineData{Value: val})
	}
	globalHandler.mux.Unlock()
	return items
}

// min Gehalt
func generateGehMinItems(keys []int) []opts.LineData {
	globalHandler.mux.Lock()
	items := make([]opts.LineData, 0, len(keys))

	for _, key := range keys {
		val := globalHandler.receivedDumps[key].Global.GEHMIN
		items = append(items, opts.LineData{Value: val})
	}
	globalHandler.mux.Unlock()
	return items
}
func generatePEItems(keys []int) []opts.LineData {
	globalHandler.mux.Lock()
	items := make([]opts.LineData, 0, len(keys))

	for _, key := range keys {
		wurzLen := globalHandler.receivedDumps[key].Global.WURZ
		uptake := 0.0
		for i := 0; i < wurzLen; i++ {
			uptake = uptake + globalHandler.receivedDumps[key].Global.PE[i]
		}

		items = append(items, opts.LineData{Value: uptake})
	}
	globalHandler.mux.Unlock()
	return items
}
