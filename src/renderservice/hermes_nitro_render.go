package main

import (
	"io"
	"net/http"
	"os"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/components"
	"github.com/go-echarts/go-echarts/v2/opts"
)

func n2odebughttpserver(w http.ResponseWriter, _ *http.Request) {

	page := components.NewPage()
	keys := extractSortedKeys()
	dates := keysAsDate(Kalender, keys, false)
	errKeys := generateErrorItems(keys)

	//NH4N
	//N2Odencum
	//NH4Sum
	//NH4UMS
	//N2onitsum
	page.AddCharts(
		lineMultiNH4N(keys, errKeys, dates),
		lineMultiN2Odencum(keys, errKeys, dates),
		lineMultiNH4UMS(keys, dates),
		lineMultiDNH4UMS(keys, dates),
		lineMultiN2onitsum(keys, errKeys, dates),
	)

	page.Render(w)
	f, err := os.Create("n2o_last_run.html")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	page.Render(io.MultiWriter(f))

}

func lineMultiNH4N(keys, errKeys []int, dates []string) *charts.Line {

	line := makeMultiLine("NH4N")

	line.SetXAxis(dates).
		AddSeries("NH4Sum", generateNH4SumItems(keys), errorMarker(errKeys, 23)).
		AddSeries("NH4N", generateNH4NItems(keys))
	return line
}

func generateNH4NItems(keys []int) []opts.LineData {
	globalHandler.mux.Lock()
	items := make([]opts.LineData, 0, len(keys))

	for i, key := range keys {
		dmp := globalHandler.receivedDumps[key]
		if dmp.Global.NDG.Index > 0 &&
			dmp.Zeit == dmp.Global.ZTDG[dmp.Global.NDG.Index-1]+1 &&
			dmp.Step == 1 {

			//g.NH4Sum = g.NH4Sum + g.NH4N[g.NDG.Index]
			val := dmp.Global.NH4N[dmp.Global.NDG.Index-1]
			items = append(items, opts.LineData{Value: val, Symbol: "diamond", XAxisIndex: i})
		}
	}
	globalHandler.mux.Unlock()
	return items
}

func generateNH4SumItems(keys []int) []opts.LineData {
	globalHandler.mux.Lock()
	items := make([]opts.LineData, 0, len(keys))

	for _, key := range keys {
		val := globalHandler.receivedDumps[key].Global.NH4Sum
		items = append(items, opts.LineData{Value: val})
	}
	globalHandler.mux.Unlock()
	return items
}

func lineMultiN2Odencum(keys, errKeys []int, dates []string) *charts.Line {

	line := makeMultiLine("N2Odencum")

	line.SetXAxis(dates).
		AddSeries("N2Odencum", generateN2OdencumItems(keys), errorMarker(errKeys, 23))
	return line
}

func generateN2OdencumItems(keys []int) []opts.LineData {
	globalHandler.mux.Lock()
	items := make([]opts.LineData, 0, len(keys))

	for _, key := range keys {
		val := globalHandler.receivedDumps[key].Global.N2Odencum
		items = append(items, opts.LineData{Value: val})
	}
	globalHandler.mux.Unlock()
	return items
}

func lineMultiDNH4UMS(keys []int, dates []string) *charts.Line {

	line := makeMultiLine("DNH4UMS")

	globalHandler.mux.Lock()
	num := 0
	for key := range globalHandler.receivedDumps {
		g := globalHandler.receivedDumps[key].Global
		num = g.IZM / g.DZ.Index
		break
	}
	globalHandler.mux.Unlock()
	line.SetXAxis(dates)
	for i := 0; i < num; i++ {
		line.AddSeries("DNH4UMS 1", generateDNH4UMSItems(keys, 0))
		line.AddSeries("DNH4UMS 2", generateDNH4UMSItems(keys, 1))
		line.AddSeries("DNH4UMS 3", generateDNH4UMSItems(keys, 2))
	}

	return line
}
func lineMultiNH4UMS(keys []int, dates []string) *charts.Line {

	line := makeMultiLine("NH4UMS")

	line.SetXAxis(dates)
	line.AddSeries("NH4UMS 1", generateNH4UMSItems(keys))
	return line
}
func generateNH4UMSItems(keys []int) []opts.LineData {
	globalHandler.mux.Lock()
	items := make([]opts.LineData, 0, len(keys))

	for _, key := range keys {
		val := globalHandler.receivedDumps[key].Global.NH4UMS
		items = append(items, opts.LineData{Value: val})
	}
	globalHandler.mux.Unlock()
	return items
}

func generateDNH4UMSItems(keys []int, index int) []opts.LineData {
	globalHandler.mux.Lock()
	items := make([]opts.LineData, 0, len(keys))

	for _, key := range keys {
		val := globalHandler.receivedNitroDumps[key].Nitro.DNH4UMS[index]
		items = append(items, opts.LineData{Value: val})
	}
	globalHandler.mux.Unlock()
	return items
}

func lineMultiN2onitsum(keys, errKeys []int, dates []string) *charts.Line {

	line := makeMultiLine("N2onitsum")

	line.SetXAxis(dates).
		AddSeries("N2onitsum", generateN2onitsumItems(keys), errorMarker(errKeys, 23))
	return line
}
func generateN2onitsumItems(keys []int) []opts.LineData {
	globalHandler.mux.Lock()
	items := make([]opts.LineData, 0, len(keys))

	for _, key := range keys {
		val := globalHandler.receivedDumps[key].Global.N2onitsum
		items = append(items, opts.LineData{Value: val})
	}
	globalHandler.mux.Unlock()
	return items
}
