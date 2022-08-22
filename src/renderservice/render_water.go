package main

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/components"
	"github.com/go-echarts/go-echarts/v2/opts"
)

func waterhttpserver(w http.ResponseWriter, _ *http.Request) {

	page := components.NewPage()
	keys := extractSortedKeys()
	dates := keysAsDate(Kalender, keys)
	errKeys := generateErrorItems(keys)

	page.AddCharts(
		lineMultiSTORAGE(keys, dates),
		lineMultiQ1(keys, errKeys, dates),
		lineMultiQ1End(keys, dates),
		lineMultiSickerDaily(keys, dates),
	)

	page.Render(w)
	f, err := os.Create("water_last_run.html")
	if err != nil {
		panic(err)
	}
	page.Render(io.MultiWriter(f))

}

func lineMultiQ1End(keys []int, dates []string) *charts.Line {
	line := makeMultiLine("Q1 Fluss durch untere Ebenen")
	N := numberOfLayer(keys[0])
	line.SetXAxis(dates)
	line.AddSeries(fmt.Sprintf("Q1 Layer %d", N-1), generateQ1Items(keys, N-1))
	line.AddSeries(fmt.Sprintf("Q1 Layer %d", N), generateQ1Items(keys, N))

	return line
}

func numberOfLayer(index int) int {
	globalHandler.mux.Lock()
	defer globalHandler.mux.Unlock()
	return globalHandler.receivedDumps[index].Global.N
}

func lineMultiSTORAGE(keys []int, dates []string) *charts.Line {
	line := makeMultiLine("STORAGE")

	line.SetXAxis(dates).
		AddSeries("STORAGE", generateSTORAGEItems(keys))

	return line
}

func generateSTORAGEItems(keys []int) []opts.LineData {
	globalHandler.mux.Lock()
	items := make([]opts.LineData, 0, len(keys))

	for _, key := range keys {
		val := globalHandler.receivedDumps[key].Global.STORAGE
		items = append(items, opts.LineData{Value: val})
	}
	globalHandler.mux.Unlock()
	return items
}

func lineMultiSickerDaily(keys []int, dates []string) *charts.Line {
	line := makeMultiLine("Sicker Daily")

	line.SetXAxis(dates).
		AddSeries("SickerDaily", generateSickerDailyItems(keys))

	return line
}

func generateSickerDailyItems(keys []int) []opts.LineData {
	globalHandler.mux.Lock()
	items := make([]opts.LineData, 0, len(keys))

	for _, key := range keys {
		val := globalHandler.receivedDumps[key].Global.SickerDaily
		items = append(items, opts.LineData{Value: val})
	}
	globalHandler.mux.Unlock()
	return items
}
