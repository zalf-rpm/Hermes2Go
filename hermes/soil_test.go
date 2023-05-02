package hermes

import (
	"encoding/csv"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"testing"
)

func TestSandAndClayToKa5Texture(t *testing.T) {
	type args struct {
		sand int
		clay int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"SS", args{92, 4}, "SS "},
		{"ST2", args{83, 12}, "ST2"},
		{"ST3", args{74, 18}, "ST3"},
		{"LTS", args{40, 35}, "LTS"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			img := generatePic()

			// Set color for each pixel.
			for clayIdx := 0; clayIdx < 100; clayIdx++ {
				for sandIdx := 99; sandIdx >= 0; sandIdx-- {
					texture := SandAndClayToKa5Texture(sandIdx, clayIdx)
					silt := 100 - sandIdx - clayIdx
					img.Set(clayIdx, 100-silt, textureToColor(texture))
				}
			}
			saveImg(img, "test_data/"+tt.name+".png")
			if got := SandAndClayToKa5Texture(tt.args.sand, tt.args.clay); got != tt.want {
				t.Errorf("SandAndClayToHa5Texture() = %v, want %v", got, tt.want)
			}
		})
	}
}

func textureToColor(texture string) color.RGBA {

	reinsande := color.RGBA{255, 255, 219, 0xff}
	lehmsande := color.RGBA{255, 255, 0, 0xff}
	schluffsande := color.RGBA{255, 231, 1, 0xff}
	sandlehme := color.RGBA{229, 187, 43, 0xff}
	normallehme := color.RGBA{192, 138, 23, 0xff}
	tonlehme := color.RGBA{154, 86, 23, 0xff}
	sandschluffe := color.RGBA{255, 215, 186, 0xff}
	lehmschluffe := color.RGBA{247, 176, 107, 0xff}
	tonschluffe := color.RGBA{232, 141, 70, 0xff}
	schlufftone := color.RGBA{233, 140, 226, 0xff}
	lehmtone := color.RGBA{203, 157, 224, 0xff}

	switch texture {

	case "SS ":
		return reinsande
	case "ST2":
		return lehmsande
	case "ST3":
		return sandlehme
	case "SU2":
		return lehmsande
	case "SU3":
		return schluffsande
	case "SU4":
		return schluffsande
	case "SL2":
		return lehmsande
	case "SL3":
		return lehmsande
	case "SL4":
		return sandlehme
	case "SLU":
		return sandlehme
	case "LS2":
		return normallehme
	case "LS3":
		return normallehme
	case "LS4":
		return normallehme
	case "LT2":
		return normallehme
	case "LT3":
		return schlufftone
	case "LTS":
		return tonlehme
	case "LU ":
		return tonschluffe
	case "ULS":
		return lehmschluffe
	case "US ":
		return sandschluffe
	case "UU ":
		return sandschluffe
	case "UT2":
		return lehmschluffe
	case "UT3":
		return lehmschluffe
	case "UT4":
		return tonschluffe
	case "TS2":
		return lehmtone
	case "TS3":
		return tonlehme
	case "TS4":
		return tonlehme
	case "TL ":
		return lehmtone
	case "TU3":
		return schlufftone
	case "TU2":
		return lehmtone
	case "TU4":
		return schlufftone
	case "TT ":
		return lehmtone
	default:
		return color.RGBA{0, 0, 0, 0xff}
	}
}

func generatePic() *image.RGBA {
	width := 100
	height := 100

	upLeft := image.Point{0, 0}
	lowRight := image.Point{width, height}

	img := image.NewRGBA(image.Rectangle{upLeft, lowRight})
	return img
}

func saveImg(img *image.RGBA, imgName string) {
	// Encode as PNG.
	f, _ := os.Create(imgName)
	png.Encode(f, img)
}

func TestSandAndClayToKa5TextureInHypar(t *testing.T) {

	// generate all possible textures
	// from 0 to 100% sand and clay

	allTextures := make(map[string]bool)
	for clayIdx := 0; clayIdx < 100; clayIdx++ {
		for sandIdx := 99; sandIdx >= 0; sandIdx-- {
			if sandIdx+clayIdx <= 100 {
				texture := SandAndClayToKa5Texture(sandIdx, clayIdx)
				if texture == "" {
					continue
				}
				allTextures[texture] = true
			}
		}
	}
	type args struct {
		texture string
		path    string
	}
	tests := []struct {
		name string
		args args
		want string
	}{}
	root := AskDirectory()
	hyparName := filepath.Join(root, "../examples/parameter/HYPAR.TRU")
	hyparName, err := filepath.Abs(hyparName)
	if err != nil {
		t.Errorf("Hypar() = %v, File %v", err, hyparName)
	}
	for texture := range allTextures {
		tests = append(tests, struct {
			name string
			args args
			want string
		}{texture, args{texture, hyparName}, texture})
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FindTextureInHypar(tt.args.texture, tt.args.path); got != tt.want {
				t.Errorf("Hypar() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSandAndClayToKa5TextureInParcap(t *testing.T) {

	// generate all possible textures
	// from 0 to 100% sand and clay

	allTextures := make(map[string]bool)
	for clayIdx := 0; clayIdx < 100; clayIdx++ {
		for sandIdx := 99; sandIdx >= 0; sandIdx-- {
			if sandIdx+clayIdx <= 100 {
				texture := SandAndClayToKa5Texture(sandIdx, clayIdx)
				if texture == "" {
					continue
				}
				allTextures[texture] = true
			}
		}
	}
	type args struct {
		texture string
		path    string
	}
	tests := []struct {
		name string
		args args
		want string
	}{}
	root := AskDirectory()
	hyparName := filepath.Join(root, "../examples/parameter/PARCAP.TRU")
	hyparName, err := filepath.Abs(hyparName)
	if err != nil {
		t.Errorf("PARCAP() = %v, File %v", err, hyparName)
	}
	for texture := range allTextures {
		tests = append(tests, struct {
			name string
			args args
			want string
		}{texture, args{texture, hyparName}, texture})
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FindTextureInPARCAP(tt.args.texture, tt.args.path); got != tt.want {
				t.Errorf("PARCAP() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSoilCompressionOverTime(t *testing.T) {
	type args struct {
		sumke                  float64
		startBD                float64
		currentBD              float64
		cOrg                   float64
		fc                     float64
		layerDepth             float64
		precip                 float64
		tillagePoreSpaceFactor float64
	}
	// load csv file with test data
	path := "test_data/test_setup_tillage.csv"
	file, err := os.Open(path)
	if err != nil {
		t.Errorf("Error opening file %v", err)
	}
	defer file.Close()
	reader := csv.NewReader(file)
	reader.Comma = ','
	startBDGoodSoil := make([]float64, 4)
	startBDBadSoil := make([]float64, 4)
	layerdepth := make([]float64, 4)
	bdgood := make([][]float64, 172)
	bdbad := make([][]float64, 172)
	mineralizationGood := make([][]float64, 172)
	mineralizationBad := make([][]float64, 172)
	airVolGood := make([][]float64, 172)
	airVolBad := make([][]float64, 172)
	sumke := make([]float64, 0, 172)
	precip := make([]float64, 0, 172)
	index := -1
	for line, err := reader.Read(); err != io.EOF; line, err = reader.Read() {
		index++
		// skip header
		if index == 0 {
			for i := 0; i < 4; i++ {
				layerdepth[i], err = strconv.ParseFloat(line[19+i], 64)
				if err != nil {
					t.Errorf("Error reading layerdepth %v layer %d", err, i)
				}
			}
			continue
		}
		if index == 1 {
			for i := 0; i < 4; i++ {
				startBDGoodSoil[i], err = strconv.ParseFloat(line[3+i], 64)
				if err != nil {
					t.Errorf("Error reading good startBD %v layer %d", err, i)
				}
				startBDBadSoil[i], err = strconv.ParseFloat(line[7+i], 64)
				if err != nil {
					t.Errorf("Error reading startBD %v layer %d", err, i)
				}
			}
			continue
		}
		bdgood[index-2] = make([]float64, 4)
		bdbad[index-2] = make([]float64, 4)
		mineralizationGood[index-2] = make([]float64, 4)
		mineralizationBad[index-2] = make([]float64, 4)
		airVolGood[index-2] = make([]float64, 4)
		airVolBad[index-2] = make([]float64, 4)
		for i := 0; i < 4; i++ {

			bdgood[index-2][i], err = strconv.ParseFloat(line[3+i], 64)
			if err != nil {
				t.Errorf("Error reading good bd %v layer %d", err, i)
			}
			bdbad[index-2][i], err = strconv.ParseFloat(line[7+i], 64)
			if err != nil {
				t.Errorf("Error reading bd %v layer %d", err, i)
			}
			airVolGood[index-2][i], err = strconv.ParseFloat(line[11+i], 64)
			if err != nil {
				t.Errorf("Error reading good airVol %v layer %d", err, i)
			}
			airVolBad[index-2][i], err = strconv.ParseFloat(line[15+i], 64)
			if err != nil {
				t.Errorf("Error reading airVol %v layer %d", err, i)
			}
			mineralizationGood[index-2][i], err = strconv.ParseFloat(line[19+i], 64)
			if err != nil {
				t.Errorf("Error reading good mineralization %v layer %d", err, i)
			}
			mineralizationBad[index-2][i], err = strconv.ParseFloat(line[24+i], 64)
			if err != nil {
				t.Errorf("Error reading mineralization %v layer %d", err, i)
			}

		}
		ke, err := strconv.ParseFloat(line[2], 64)
		if err != nil {
			t.Errorf("Error reading ke %v", err)
		}
		sumke = append(sumke, ke)
		precipVal, err := strconv.ParseFloat(line[1], 64)
		if err != nil {
			t.Errorf("Error reading precip %v", err)
		}
		precip = append(precip, precipVal)

	}
	type testStruct struct {
		name                     string
		args                     args
		wantNewBD                float64
		wantNewSumke             float64
		wantAirPoreVolume        float64
		wantMineralisationFactor float64
	}
	tests := make([]testStruct, 0, len(precip)*2)
	for i := 1; i < len(precip); i++ {
		for j := 0; j < 4; j++ {
			tests = append(tests, testStruct{
				name: fmt.Sprintf("Test good soil day %d layer depth %f", i+1, layerdepth[j]),
				args: args{
					sumke:                  sumke[i-1],
					startBD:                startBDGoodSoil[j],
					currentBD:              bdgood[i-1][j],
					cOrg:                   1.17,
					fc:                     0.33,
					layerDepth:             layerdepth[j],
					precip:                 precip[i],
					tillagePoreSpaceFactor: 0.001,
				},
				wantNewBD:                bdgood[i][j],
				wantNewSumke:             sumke[i],
				wantAirPoreVolume:        airVolGood[i][j],
				wantMineralisationFactor: mineralizationGood[i][j],
			}, testStruct{
				name: fmt.Sprintf("Test bad soil day %d layer depth %f", i+1, layerdepth[j]),
				args: args{
					sumke:                  sumke[i-1],
					startBD:                startBDBadSoil[j],
					currentBD:              bdbad[i-1][j],
					cOrg:                   1.03,
					fc:                     0.24,
					layerDepth:             layerdepth[j],
					precip:                 precip[i],
					tillagePoreSpaceFactor: 0.001,
				},
				wantNewBD:                bdbad[i][j],
				wantNewSumke:             sumke[i],
				wantAirPoreVolume:        airVolBad[i][j],
				wantMineralisationFactor: mineralizationBad[i][j],
			})
		}
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotNewBD, gotNewSumke, gotAirPoreVolume, gotMineralisationFactor := SoilCompressionOverTime(tt.args.sumke, tt.args.startBD, tt.args.currentBD, tt.args.cOrg, tt.args.fc, tt.args.layerDepth, tt.args.precip, 1.0, tt.args.tillagePoreSpaceFactor, false)
			if math.Abs(gotNewBD-tt.wantNewBD) > 0.00001 {
				t.Errorf("SoilCompressionOverTime() gotNewBD = %v, want %v", gotNewBD, tt.wantNewBD)
			}
			if math.Abs(gotNewSumke-tt.wantNewSumke) > 0.00001 {
				t.Errorf("SoilCompressionOverTime() gotNewSumke = %v, want %v", gotNewSumke, tt.wantNewSumke)
			}
			if math.Abs(gotAirPoreVolume-tt.wantAirPoreVolume) > 0.00001 {
				t.Errorf("SoilCompressionOverTime() gotAirPoreVolume = %v, want %v", gotAirPoreVolume, tt.wantAirPoreVolume)
			}
			if math.Abs(gotMineralisationFactor-tt.wantMineralisationFactor) > 0.00001 {
				t.Errorf("SoilCompressionOverTime() gotMineralisationFactor = %v, want %v", gotMineralisationFactor, tt.wantMineralisationFactor)
			}
		})
	}
}
