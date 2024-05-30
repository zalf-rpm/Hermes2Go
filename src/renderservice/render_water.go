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
	dates := keysAsDate(Kalender, keys, false)
	errKeys := generateErrorItems(keys)
	dailyKeys := extractDailyKeys()
	dailyKeysDates := keysAsDate(Kalender, dailyKeys, true)

	page.AddCharts(
		lineMultiDailyWater(dailyKeys, dailyKeysDates),
		lineMultiDailyIO(dailyKeys, dailyKeysDates),
		lineMultiSTORAGE(keys, dates),
		lineMultiQ1(keys, errKeys, dates),
		lineMultiQ1End(keys, dates),
		lineMultiSickerDaily(keys, dates),
		//themeRiverTime(keys),
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
func lineMultiDailyWater(keys []int, dates []string) *charts.Line {
	line := makeMultiLine("Water Balcance Daily")

	// WDT:             wdt,
	// N:               g.N,
	// WG:              g.WG,
	// W:               g.W, // field capacity
	// DZ:              g.DZ.Num,
	// REGEN:           g.REGEN[g.TAG.Index],
	// SumWaterContent: g.SumWaterContent, // sum of all layers
	// Irrigation:      g.EffectiveIRRIG,
	// CapillaryRise:   g.CAPSUM,
	// Evaporation:     g.ETA,
	// Drainage:        g.DRAISUM,
	// Storage:         g.STORAGE,

	line.SetXAxis(dates).
		AddSeries("SumWaterContent", generateSumWaterContentItems(keys)).
		AddSeries("REGEN", generateREGENItems(keys)).
		AddSeries("Evaporation", generateEvaporationItems(keys)).
		AddSeries("Storage", generateStorageItems(keys)).
		AddSeries("SickerLoss", generateSickerLossItems(keys)).
		AddSeries("Infiltration", generateInfiltrationItems(keys)).
		AddSeries("TPSumDaily", generateTPSumDailyItems(keys))

	return line
}
func lineMultiDailyIO(keys []int, dates []string) *charts.Line {
	line := makeMultiLine("Daily Water I/O")

	line.SetXAxis(dates).
		AddSeries("SumWaterContentDiff", generateSumWaterContentDiffItems(keys)).
		AddSeries("Loss", generateWaterLossItems(keys)).
		AddSeries("Infiltration", generateInfiltrationItems(keys)).
		AddSeries("IODiff", generateWaterDiffItems(keys)).
		AddSeries("WaterDiff", generateWaterDiff2Items(keys))

	return line
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

func generateSumWaterContentItems(keys []int) []opts.LineData {
	globalHandler.mux.Lock()
	items := make([]opts.LineData, 0, len(keys))

	for _, key := range keys {
		val := globalHandler.waterBalanceDumps[key].SumWaterContent
		items = append(items, opts.LineData{Value: val})
	}
	globalHandler.mux.Unlock()
	return items
}
func generateSumWaterContentDiffItems(keys []int) []opts.LineData {
	globalHandler.mux.Lock()
	items := make([]opts.LineData, 0, len(keys))

	prevDay := 0.0
	first := true
	for _, key := range keys {
		if first {
			first = false
			prevDay = globalHandler.waterBalanceDumps[key].SumWaterContent
		}
		val := globalHandler.waterBalanceDumps[key].SumWaterContent
		items = append(items, opts.LineData{Value: val - prevDay})
	}
	globalHandler.mux.Unlock()
	return items
}

func generateREGENItems(keys []int) []opts.LineData {
	globalHandler.mux.Lock()
	items := make([]opts.LineData, 0, len(keys))

	for _, key := range keys {
		val := globalHandler.waterBalanceDumps[key].REGEN
		items = append(items, opts.LineData{Value: val})
	}
	globalHandler.mux.Unlock()
	return items
}
func generateEvaporationItems(keys []int) []opts.LineData {
	globalHandler.mux.Lock()
	items := make([]opts.LineData, 0, len(keys))

	for _, key := range keys {
		val := globalHandler.waterBalanceDumps[key].Evaporation
		items = append(items, opts.LineData{Value: val})
	}
	globalHandler.mux.Unlock()
	return items
}
func generateStorageItems(keys []int) []opts.LineData {
	globalHandler.mux.Lock()
	items := make([]opts.LineData, 0, len(keys))

	for _, key := range keys {
		val := globalHandler.waterBalanceDumps[key].Storage
		items = append(items, opts.LineData{Value: val})
	}
	globalHandler.mux.Unlock()
	return items
}

func generateTPSumDailyItems(keys []int) []opts.LineData {
	globalHandler.mux.Lock()
	items := make([]opts.LineData, 0, len(keys))

	for _, key := range keys {
		val := globalHandler.waterBalanceDumps[key].TPSumDaily
		items = append(items, opts.LineData{Value: val})
	}
	globalHandler.mux.Unlock()
	return items
}

func generateWaterLossItems(keys []int) []opts.LineData {
	globalHandler.mux.Lock()
	items := make([]opts.LineData, 0, len(keys))

	for _, key := range keys {
		val1 := globalHandler.waterBalanceDumps[key].SickerLoss
		val2 := globalHandler.waterBalanceDumps[key].Drainage
		val3 := globalHandler.waterBalanceDumps[key].Evaporation
		val4 := globalHandler.waterBalanceDumps[key].TPSumDaily

		items = append(items, opts.LineData{Value: (val1 + val2 + val3 + val4) * -1})
	}
	globalHandler.mux.Unlock()
	return items
}

func generateWaterDiffItems(keys []int) []opts.LineData {
	globalHandler.mux.Lock()
	items := make([]opts.LineData, 0, len(keys))

	for _, key := range keys {
		out1 := globalHandler.waterBalanceDumps[key].SickerLoss
		out2 := globalHandler.waterBalanceDumps[key].Drainage
		out3 := globalHandler.waterBalanceDumps[key].Evaporation
		out4 := globalHandler.waterBalanceDumps[key].TPSumDaily
		in := globalHandler.waterBalanceDumps[key].Infiltration

		items = append(items, opts.LineData{Value: in - (out1 + out2 + out3 + out4)})
	}
	globalHandler.mux.Unlock()
	return items
}
func generateWaterDiff2Items(keys []int) []opts.LineData {
	globalHandler.mux.Lock()
	items := make([]opts.LineData, 0, len(keys))

	for _, key := range keys {
		out1 := globalHandler.waterBalanceDumps[key].WaterDiff
		items = append(items, opts.LineData{Value: out1})
	}
	globalHandler.mux.Unlock()
	return items
}
func generateSickerLossItems(keys []int) []opts.LineData {
	globalHandler.mux.Lock()
	items := make([]opts.LineData, 0, len(keys))

	for _, key := range keys {
		val := globalHandler.waterBalanceDumps[key].SickerLoss
		items = append(items, opts.LineData{Value: val})
	}
	globalHandler.mux.Unlock()
	return items
}
func generateInfiltrationItems(keys []int) []opts.LineData {
	globalHandler.mux.Lock()
	items := make([]opts.LineData, 0, len(keys))

	for _, key := range keys {
		val := globalHandler.waterBalanceDumps[key].Infiltration
		items = append(items, opts.LineData{Value: val})
	}
	globalHandler.mux.Unlock()
	return items
}
