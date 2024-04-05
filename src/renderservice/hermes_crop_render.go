package main

import (
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/components"
	"github.com/go-echarts/go-echarts/v2/opts"
)

func cropdebughttpserver(w http.ResponseWriter, _ *http.Request) {

	page := components.NewPage()
	keys := extractSortedKeys()
	dates := keysAsDate(Kalender, keys)
	stageKeys, stageValues := generateDevStageItems(keys)

	// GEHMIN
	// GEHMAX
	// PE
	page.AddCharts(
		lineMultiGEHALT(keys, stageKeys, stageValues, dates),
		lineMultiNUptake(keys, stageKeys, stageValues, dates),
		lineMultiWaterUptake(keys, stageKeys, stageValues, dates),
		lineMultiOrganBioMass(keys, stageKeys, stageValues, dates),
	)

	page.Render(w)
	f, err := os.Create("crop_last_run.html")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	page.Render(io.MultiWriter(f))

}

func generateDevStageItems(keys []int) ([]int, []int) {
	globalHandler.mux.Lock()
	listDevStageChange := []int{}
	listDevStageChangeValue := []int{}

	currentIntWICK := -1
	for _, key := range keys {
		tag := globalHandler.receivedDumps[key].Global.INTWICK.Index
		if tag >= 0 && tag != currentIntWICK {
			listDevStageChange = append(listDevStageChange, key)
			listDevStageChangeValue = append(listDevStageChangeValue, tag+1)
			currentIntWICK = tag
		}
	}
	globalHandler.mux.Unlock()
	return listDevStageChange, listDevStageChangeValue
}

func lineMultiGEHALT(keys, devStage, devStageVal []int, dates []string) *charts.Line {

	line := makeMultiLine("GEHALT")

	line.SetXAxis(dates).
		AddSeries("GEHMIN", generateGehMinItems(keys), stageMarker(devStage, devStageVal)).
		AddSeries("GEHMAX", generateGehMaxItems(keys)).
		AddSeries("WUGEH", generateWuGehItems(keys)).
		AddSeries("GEHOB", generateGehobItems(keys))

	return line
}

func lineMultiNUptake(keys, devStage, devStageVal []int, dates []string) *charts.Line {

	line := makeMultiLine("Daily N Uptake")

	line.SetXAxis(dates).
		AddSeries("PE", generatePEItems(keys), stageMarker(devStage, devStageVal))

	return line
}

func lineMultiWaterUptake(keys, devStage, devStageVal []int, dates []string) *charts.Line {

	line := makeMultiLine("Daily Water Uptake")

	line.SetXAxis(dates).
		AddSeries("TP", generateTPItems(keys), stageMarker(devStage, devStageVal))

	return line
}
func lineMultiOrganBioMass(keys, devStage, devStageVal []int, dates []string) *charts.Line {

	line := makeMultiLine("Daily Water Uptake")

	line.SetXAxis(dates).
		AddSeries("WORG0", generateWORGItems(keys, 0), stageMarker(devStage, devStageVal)).
		AddSeries("WORG1", generateWORGItems(keys, 1)).
		AddSeries("WORG2", generateWORGItems(keys, 2)).
		AddSeries("WORG3", generateWORGItems(keys, 3)).
		AddSeries("WORG4", generateWORGItems(keys, 4)).
		AddSeries("HARVEST", generateYieldItems(keys))

	return line
}

func stageMarker(stageKeys, values []int) charts.SeriesOpts {
	dates := keysAsDate(Kalender, stageKeys)

	marker := make([]opts.MarkPointNameCoordItem, 0, len(dates))
	for idx, date := range dates {
		// stage index as name
		name := strconv.Itoa(values[idx])
		marker = append(marker, opts.MarkPointNameCoordItem{
			Name:       name,
			Coordinate: []interface{}{date, 0},
			Label:      &opts.Label{Show: true, Color: "white", Position: "inside", Formatter: "{b}"},
		})
	}
	options := charts.WithMarkPointNameCoordItemOpts(marker...)

	return options
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

func generateTPItems(keys []int) []opts.LineData {
	globalHandler.mux.Lock()
	items := make([]opts.LineData, 0, len(keys))

	for _, key := range keys {
		wurzLen := globalHandler.receivedDumps[key].Global.WURZ
		uptake := 0.0
		for i := 0; i < wurzLen; i++ {
			uptake = uptake + globalHandler.receivedDumps[key].Global.TP[i]
		}

		items = append(items, opts.LineData{Value: uptake})
	}
	globalHandler.mux.Unlock()
	return items
}

func generateWORGItems(keys []int, organIndex int) []opts.LineData {
	globalHandler.mux.Lock()
	items := make([]opts.LineData, 0, len(keys))

	for _, key := range keys {
		val := globalHandler.receivedDumps[key].Global.WORG[organIndex]
		items = append(items, opts.LineData{Value: val})
	}
	globalHandler.mux.Unlock()
	return items
}

func generateYieldItems(keys []int) []opts.LineData {
	globalHandler.mux.Lock()
	items := make([]opts.LineData, 0, len(keys))

	for _, key := range keys {
		val := globalHandler.receivedDumps[key].Global.HARVEST
		items = append(items, opts.LineData{Value: val})
	}
	globalHandler.mux.Unlock()
	return items
}
