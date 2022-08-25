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
		themeRiverTime(keys),
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

func themeRiverTime(keys []int) *charts.ThemeRiver {
	tr := charts.NewThemeRiver()
	tr.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title: "Q1 soil Layer",
		}),
		charts.WithSingleAxisOpts(opts.SingleAxis{
			Type:   "time",
			Bottom: "10%",
		}),
		charts.WithTooltipOpts(opts.Tooltip{
			Trigger: "axis",
		}),
	)
	timeKeys := make(map[int]bool, len(keys))
	for k := range keys {
		zeit := k / 100
		timeKeys[zeit] = true
	}
	extractRiverData := func(n int, items []opts.ThemeRiverData, prevTime int, prevQ1 [22]float64) []opts.ThemeRiverData {
		for q := 0; q < n; q++ {
			items = append(items, opts.ThemeRiverData{
				Date:  KalenderLong(prevTime),
				Value: prevQ1[q],
				Name:  fmt.Sprintf("Q1 Layer %d", q),
			})
		}
		return items
	}

	items := make([]opts.ThemeRiverData, 0, len(timeKeys))
	globalHandler.mux.Lock()
	prevTime := 0
	N := 0
	var prevQ1 [22]float64
	for _, key := range keys {

		q1 := globalHandler.receivedDumps[key].Global.Q1
		if prevTime != globalHandler.receivedDumps[key].Zeit {
			// only take the last entry of the day
			if prevTime > 0 {
				N = globalHandler.receivedDumps[key].Global.N
				items = extractRiverData(N, items, prevTime, prevQ1)
			}
			prevTime = globalHandler.receivedDumps[key].Zeit
			prevQ1 = q1
		} else {
			for i := range q1 {
				prevQ1[i] += q1[i]
			}
		}
	}
	if prevTime > 0 {
		items = extractRiverData(N, items, prevTime, prevQ1)
	}

	globalHandler.mux.Unlock()

	tr.AddSeries("Q1 sums", items)
	return tr
}
