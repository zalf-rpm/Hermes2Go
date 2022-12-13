package main

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/components"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/zalf-rpm/Hermes2Go/hermes"
)

func main() {

	fmt.Println("Open a browser and connet to: http://localhost:8081 ")
	http.HandleFunc("/4", ptf4httpserver)
	http.HandleFunc("/3", ptf3httpserver)
	http.HandleFunc("/2", ptf2httpserver)
	http.HandleFunc("/1", ptf1httpserver)
	http.ListenAndServe("localhost:8081", nil)
}

func ptf4httpserver(w http.ResponseWriter, _ *http.Request) {
	page := components.NewPage()

	page.AddCharts(
		sandClay(0.0, 4),
		sandClay(0.1, 4),
		sandClay(0.5, 4),
		sandClay(1.0, 4),
		sandClay(10.0, 4),
		sandClay(20.0, 4),
		sandClay(30.0, 4),
	)

	page.Render(w)
	f, err := os.Create("ptf4.html")
	if err != nil {
		panic(err)
	}
	page.Render(io.MultiWriter(f))
}

func ptf3httpserver(w http.ResponseWriter, _ *http.Request) {

	page := components.NewPage()

	page.AddCharts(
		sandClay(0.0, 3),
		sandClay(0.1, 3),
		sandClay(0.5, 3),
		sandClay(1.0, 3),
		sandClay(10.0, 3),
		sandClay(20.0, 3),
		sandClay(30.0, 3),
	)

	page.Render(w)
	f, err := os.Create("ptf3.html")
	if err != nil {
		panic(err)
	}
	page.Render(io.MultiWriter(f))

}
func ptf2httpserver(w http.ResponseWriter, _ *http.Request) {

	page := components.NewPage()

	page.AddCharts(
		sandClay(0.0, 2),
		sandClay(0.1, 2),
		sandClay(0.5, 2),
		sandClay(1.0, 2),
		sandClay(10.0, 2),
		sandClay(20.0, 2),
		sandClay(30.0, 2),
	)

	page.Render(w)
	f, err := os.Create("ptf2.html")
	if err != nil {
		panic(err)
	}
	page.Render(io.MultiWriter(f))

}
func ptf1httpserver(w http.ResponseWriter, _ *http.Request) {

	page := components.NewPage()

	page.AddCharts(
		sandClay(0.0, 1),
		sandClay(0.1, 1),
		sandClay(0.5, 1),
		sandClay(1.0, 1),
		sandClay(10.0, 1),
		sandClay(20.0, 1),
		sandClay(30.0, 1),
	)

	page.Render(w)
	f, err := os.Create("ptp1.html")
	if err != nil {
		panic(err)
	}
	page.Render(io.MultiWriter(f))

}

func sandClay(cContent float64, ptf int) *charts.Surface3D {
	surface3d := charts.NewSurface3D()
	surface3d.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{Title: fmt.Sprintf("PTF: %d FC with Corg %3.1f %%", ptf, cContent)}),

		charts.WithXAxis3DOpts(opts.XAxis3D{
			Show: true,
			Name: "Sand %",
			Type: "value",
		}),
		charts.WithYAxis3DOpts(opts.YAxis3D{Name: "Clay %"}),
		charts.WithZAxis3DOpts(opts.ZAxis3D{Name: "FC %"}),
	)
	fc, pwp := genSurface3dFC(cContent, ptf)
	surface3d.AddSeries("PWP_3D", pwp)
	surface3d.AddSeries("FC_3D", fc)

	return surface3d
}

func genSurface3dFC(cGehalt float64, ptf int) ([]opts.Chart3DData, []opts.Chart3DData) {

	dataFC := make([][4]interface{}, 0)
	dataPWP := make([][3]interface{}, 0)
	for sand := 1; sand < 99; sand++ {
		for ton := 1; ton < 99; ton++ {
			if sand+ton < 100 {
				schluf := float64(100 - sand - ton)
				var fc, pwp float64
				if ptf == 4 {
					fc, pwp = hermes.PTF4(cGehalt, float64(ton), float64(sand))
					// ix := -0.837531 + 0.430183*cGehalt
					// ix2 := math.Pow(ix, 2)
					// ix3 := math.Pow(ix, 3)
					// yps := -1.40744 + 0.0661969*float64(ton)
					// yps2 := math.Pow(yps, 2)
					// yps3 := math.Pow(yps, 3)
					// zet := -1.51866 + 0.0393284*float64(sand)
					// zet2 := math.Pow(zet, 2)
					// zet3 := math.Pow(zet, 3)

					// fc := (29.7528 + 10.3544*(0.0461615+0.290955*ix-0.0496845*ix2+0.00704802*ix3+0.269101*yps-0.176528*ix*yps+0.0543138*ix2*yps+0.1982*yps2-0.060699*yps3-0.320249*zet-0.0111693*ix2*zet+0.14104*yps*zet+0.0657345*ix*yps*zet-0.102026*yps2*zet-0.04012*zet2+0.160838*ix*zet2-0.121392*yps*zet2-0.061667*zet3)) / 100
					// pwp := (14.2568 + 7.36318*(0.06865+0.108713*ix-0.0157225*ix2+0.00102805*ix3+0.886569*yps-0.223581*ix*yps+0.0126379*ix2*yps+0.0135266*ix*yps2-0.0334434*yps3-0.0535182*zet-0.0354271*ix*zet-0.00261313*ix2*zet-0.154563*yps*zet-0.0160219*ix*yps*zet-0.0400606*yps2*zet-0.104875*zet2*0.0159857*ix*zet2-0.0671656*yps*zet2-0.0260699*zet3)) / 100

					// dataFC = append(dataFC, [4]interface{}{sand, ton, fc * 100, fc > pwp})
					// dataPWP = append(dataPWP, [3]interface{}{sand, ton, pwp * 100})
				} else if ptf == 3 {

					// PTF by Batjes for pF 1.7
					fc, pwp = hermes.PTF3(cGehalt, float64(ton), float64(schluf))
					// fc := (0.6681*float64(ton) + 0.2614*schluf + 2.215*cGehalt) / 100
					// pwp := (0.3624*float64(ton) + 0.117*schluf + 1.6054*cGehalt) / 100

					// dataFC = append(dataFC, [4]interface{}{sand, ton, fc * 100, fc > pwp})
					// dataPWP = append(dataPWP, [3]interface{}{sand, ton, pwp * 100})
				} else if ptf == 1 {

					fc, pwp = hermes.PTF1(cGehalt, float64(ton), float64(schluf))
					// fc := 0.2449 - 0.1887*(1/(cGehalt+1)) + 0.004527*float64(ton) + 0.001535*schluf + 0.001442*schluf*(1/(cGehalt+1)) - 0.0000511*schluf*float64(ton) + 0.0008676*float64(ton)*(1/(cGehalt+1))
					// pwp := 0.09878 + 0.002127*float64(ton) - 0.0008366*schluf - 0.0767*(1/(cGehalt+1)) + 0.00003853*schluf*float64(ton) + 0.00233*schluf*(1/(cGehalt+1)) + 0.0009498*schluf*(1/(cGehalt+1))

					// dataFC = append(dataFC, [4]interface{}{sand, ton, fc * 100, fc > pwp})
					// dataPWP = append(dataPWP, [3]interface{}{sand, ton, pwp * 100})
				} else if ptf == 2 {
					// PTF by Batjes for pF 2.5
					fc, pwp = hermes.PTF2(cGehalt, float64(ton), float64(schluf))
					// fc := (0.46*float64(ton) + 0.3045*schluf + 2.0703*cGehalt) / 100
					// pwp := (0.3624*float64(ton) + 0.117*schluf + 1.6054*cGehalt) / 100

				}
				dataFC = append(dataFC, [4]interface{}{sand, ton, fc * 100, fc > pwp})
				dataPWP = append(dataPWP, [3]interface{}{sand, ton, pwp * 100})
			}
		}
	}
	retPWP := make([]opts.Chart3DData, 0, len(dataPWP))
	for _, d := range dataPWP {
		retPWP = append(retPWP, opts.Chart3DData{
			Value: []interface{}{d[0], d[1], d[2]},
		})
	}
	ret := make([]opts.Chart3DData, 0, len(dataFC))
	for _, d := range dataFC {
		if b, ok := d[3].(bool); ok && b {
			color := "green"
			if d[2].(float64) > 99 {
				color = "purple"
			}
			ret = append(ret, opts.Chart3DData{
				Value: []interface{}{d[0], d[1], d[2]},
				ItemStyle: &opts.ItemStyle{
					Color: color,
				},
			})
		} else {
			ret = append(ret, opts.Chart3DData{
				Value:     []interface{}{d[0], d[1], d[2]},
				ItemStyle: &opts.ItemStyle{Color: "red"},
			})

		}
	}

	return ret, retPWP
}
