package hermes

import (
	"bufio"
	"math"
	"os"
	"strconv"
	"strings"
	"testing"
)

const epsilon = 0.00001

func nearEqual(a, b float64) bool {
	return math.Abs(a-b) <= epsilon
}
func nearEqualArray(a, b []float64) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if !nearEqual(a[i], b[i]) {
			return false
		}
	}
	return true
}
func Test_root(t *testing.T) {
	type args struct {
		veloc         float64
		tempsum       float64
		numberOfLayer int
	}
	type test struct {
		name                    string
		args                    args
		wantQrez                float64
		wantRootingDepth        float64
		wantCulRootPercPerLayer []float64
	}
	testfiles := map[string]float64{
		"./test_data/root_func_oilseed_rape.csv": 0.004,
		"./test_data/root_func_wheat.csv":        0.0035,
		"./test_data/root_func_gras.csv":         0.002787,
	}
	tests := make([]test, 0, 2152)
	for testfile, veloc := range testfiles {
		file, err := os.Open(testfile)
		if err != nil {
			t.Errorf("root() failed to read test file %s", testfile)
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		index := -1
		for scanner.Scan() {
			index++
			if index == 0 {
				// skip headline
				continue
			}
			tokens := strings.FieldsFunc(scanner.Text(), func(r rune) bool {
				return r == ','
			})
			tempsum := ValAsFloat(tokens[0], testfile, scanner.Text())
			wantQrez := ValAsFloat(tokens[1], testfile, scanner.Text())
			wantRootingDepth := ValAsFloat(tokens[2], testfile, scanner.Text())
			numberOfLayer := 0
			wantCulRootPercPerLayer := []float64{}
			if len(tokens) > 3 {
				numberOfLayer = len(tokens) - 3
				wantCulRootPercPerLayer = make([]float64, numberOfLayer)
				for i := 0; i < numberOfLayer; i++ {
					wantCulRootPercPerLayer[i] = ValAsFloat(tokens[i+3], testfile, scanner.Text())
				}
			}
			tests = append(tests, test{
				name: testfile + strconv.Itoa(int(tempsum)),
				args: args{
					veloc:         veloc,
					tempsum:       tempsum,
					numberOfLayer: numberOfLayer,
				},
				wantQrez:                wantQrez,
				wantRootingDepth:        wantRootingDepth,
				wantCulRootPercPerLayer: wantCulRootPercPerLayer,
			})
		}
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotQrez, gotRootingDepth, gotCulRootPercPerLayer := root(tt.args.veloc, tt.args.tempsum, tt.args.numberOfLayer)
			if !nearEqual(gotQrez, tt.wantQrez) {
				t.Errorf("root() gotQrez = %v, want %v", gotQrez, tt.wantQrez)
			}
			if !nearEqual(gotRootingDepth, tt.wantRootingDepth) {
				t.Errorf("root() gotRootingDepth = %v, want %v", gotRootingDepth, tt.wantRootingDepth)
			}
			if !nearEqualArray(gotCulRootPercPerLayer, tt.wantCulRootPercPerLayer) {
				t.Errorf("root() gotCulRootPercPerLayer = %v, want %v", gotCulRootPercPerLayer, tt.wantCulRootPercPerLayer)
			}
		})
	}
}
