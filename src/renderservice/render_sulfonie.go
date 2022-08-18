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

func sulfoniehttpserver(w http.ResponseWriter, _ *http.Request) {

	page := components.NewPage()
	keys := extractSortedKeys()
	dates := keysAsDate(Kalender, keys)
	measurementKeys := generateMeasurementItems(keys)
	eventKeys, eventTypes := generateEventItems(keys)

	// SGEHOB
	// SREDUK
	// SGEHMAX
	// SGEHMIN
	// S1
	page.AddCharts(
		lineMultiSmin(keys, measurementKeys, dates),
		lineMultiS1(keys, measurementKeys, dates),
		lineMultiSGEHMIN(keys, measurementKeys, dates),
		lineMultiSREDUK(keys, measurementKeys, dates),
		lineMultiSGEHOB(keys, eventKeys, eventTypes, dates),
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
		AddSeries("S1 Layer 1", generateS1Items(keys, 0), measurementMarker(errKeys, 1)).
		AddSeries("S1 Layer 2", generateS1Items(keys, 1)).
		AddSeries("S1 Layer 3", generateS1Items(keys, 2))
	return line
}

func lineMultiSmin(keys, errKeys []int, dates []string) *charts.Line {
	line := makeMultiLine("Smin every 30cm")

	line.SetXAxis(dates).
		AddSeries("Layer 0-3cm", generateS1SumItems(keys, 0, 3), measurementMarker(errKeys, 1)).
		AddSeries("Layer 3-6cm", generateS1SumItems(keys, 3, 6)).
		AddSeries("Layer 6-9cm", generateS1SumItems(keys, 6, 9))
	return line
}
func generateS1SumItems(keys []int, layerDepth1, layerDepth2 int) []opts.LineData {
	globalHandler.mux.Lock()
	items := make([]opts.LineData, 0, len(keys))

	for _, key := range keys {
		var val float64
		for i := layerDepth1; i < layerDepth2; i++ {
			val += globalHandler.receivedDumps[key].Global.S1[i]
		}

		items = append(items, opts.LineData{Value: val})
	}
	globalHandler.mux.Unlock()
	return items
}

// generateS1Items generate Smin (S1) items
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
		AddSeries("SGEHMIN", generateSGEHMINItems(keys), measurementMarker(errKeys, 1)).
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
		AddSeries("SREDUK", generateSREDUKItems(keys), measurementMarker(errKeys, 1))
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

func lineMultiSGEHOB(keys, eventKeys []int, eventTypes []string, dates []string) *charts.Line {
	line := makeMultiLine("S GEHALT OBEN")

	line.SetXAxis(dates).
		AddSeries("SGEHOB", generateSGEHOBItems(keys), eventMarker(eventKeys, eventTypes, 1))
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

func generateMeasurementItems(keys []int) []int {
	globalHandler.mux.Lock()
	listOfMeasurements := []int{}

	for _, key := range keys {
		zeit := globalHandler.receivedDumps[key].Zeit
		if _, ok := globalHandler.receivedDumps[key].Global.SI[zeit]; ok {
			listOfMeasurements = append(listOfMeasurements, key)
		}
	}
	globalHandler.mux.Unlock()
	return listOfMeasurements
}

func measurementMarker(errKeys []int, offset float64) charts.SeriesOpts {
	dates := keysAsDate(Kalender, errKeys)

	marker := make([]opts.MarkPointNameCoordItem, 0, len(dates))
	for _, date := range dates {
		marker = append(marker, opts.MarkPointNameCoordItem{
			Name:       "M",
			Coordinate: []interface{}{date, 0},
			Label:      &opts.Label{Show: true, Color: "white", Position: "inside", Formatter: "{b}"},
		})
	}
	options := charts.WithMarkPointNameCoordItemOpts(marker...)

	return options
}

func generateEventItems(keys []int) ([]int, []string) {
	globalHandler.mux.Lock()
	listOfEvents := []int{}
	listOfEventTypes := []string{}
	currentIntwick := 0
	for _, key := range keys {
		dmp := globalHandler.receivedDumps[key]
		zeit := dmp.Zeit
		akf := dmp.Global.AKF.Index
		if dmp.Global.ERNTE[akf-1] == zeit {
			listOfEvents = append(listOfEvents, key)
			listOfEventTypes = append(listOfEventTypes, "har")
		}
		if currentIntwick != dmp.Global.INTWICK.Index {
			currentIntwick = dmp.Global.INTWICK.Index
			listOfEvents = append(listOfEvents, key)
			listOfEventTypes = append(listOfEventTypes, strconv.Itoa(currentIntwick))
		}
		if dmp.Global.SAAT[akf] == zeit {
			listOfEvents = append(listOfEvents, key)
			listOfEventTypes = append(listOfEventTypes, "sow")
		}

	}
	globalHandler.mux.Unlock()
	return listOfEvents, listOfEventTypes
}

func eventMarker(eventKeys []int, eventType []string, offset float64) charts.SeriesOpts {
	dates := keysAsDate(Kalender, eventKeys)

	marker := make([]opts.MarkPointNameCoordItem, 0, len(dates))
	for idx, date := range dates {
		marker = append(marker, opts.MarkPointNameCoordItem{
			Name:       eventType[idx],
			Coordinate: []interface{}{date, 0},
			Label:      &opts.Label{Show: true, Color: "white", Position: "inside", Formatter: "{b}"},
		})
	}
	options := charts.WithMarkPointNameCoordItemOpts(marker...)

	return options
}
