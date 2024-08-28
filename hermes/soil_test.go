package hermes

import (
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
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
	session := NewHermesSession()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FindTextureInHypar(tt.args.texture, tt.args.path, session); got != tt.want {
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
	session := NewHermesSession()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FindTextureInPARCAP(tt.args.texture, tt.args.path, session); got != tt.want {
				t.Errorf("PARCAP() = %v, want %v", got, tt.want)
			}
		})
	}
}
