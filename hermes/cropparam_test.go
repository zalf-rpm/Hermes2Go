package hermes

import (
	"reflect"
	"testing"
)

func TestWriteCropParam(t *testing.T) {
	type args struct {
		filename  string
		cropParam CropParam
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"testSoy", args{"test_data/PARAM_test1.SOY.yml", CropParam{
			CropName:          "soybean",
			ABBr:              "SOY",
			Variety:           "test",
			MAXAMAX:           28,
			TempTyp:           1,
			MINTMP:            6,
			WUMAXPF:           9,
			VELOC:             0.5574,
			NGEFKT:            1,
			RGA:               0,
			RGB:               0,
			SubOrgan:          0,
			AboveGroundOrgans: []int{2, 3, 4},
			YORGAN:            4,
			YIFAK:             0.8,
			INITCONCNBIOM:     6.0,
			INITCONCNROOT:     2.0,
			NRKOM:             4,
			CompartimentNames: []string{"root", "leave", "stem", "ears"},
			DAUERKULT:         false,
			LEGUM:             true,
			WORG:              []float64{153, 153, 0, 0},
			MAIRT:             []float64{0.01, 0.03, 0.015, 0.01},
			KcIni:             0.65,
			NRENTW:            7,
			CropDevelopmentStages: []CropDevelopmentStage{
				{
					DevelopmentStageName: "development phase 1: sowing til emergence",
					TSUM:                 82,
					BAS:                  7,
					VSCHWELL:             0,
					DAYL:                 0,
					DLBAS:                0,
					DRYSWELL:             1,
					LUKRIT:               0.08,
					LAIFKT:               0.0025,
					WGMAX:                0.02,
					PRO: []float64{
						0.5,
						0.5,
						0,
						0,
					},
					DEAD: []float64{
						0, 0, 0, 0,
					},
					Kc: 0.65,
				},
				{
					DevelopmentStageName: "development phase 2: emergence to end juvenile phase",
					TSUM:                 50,
					BAS:                  7,
					VSCHWELL:             0,
					DAYL:                 -15.0,
					DLBAS:                -23.5,
					DRYSWELL:             0.6,
					LUKRIT:               0.08,
					LAIFKT:               0.0030,
					WGMAX:                0.012,
					PRO: []float64{
						0.200,
						0.600,
						0.200,
						0.000,
					},
					DEAD: []float64{
						0.000,
						0.000,
						0.000,
						0.000,
					},
					Kc: 0.7,
				},
				{
					DevelopmentStageName: "development phase 3: end juvenile phase to flower appearance",
					TSUM:                 220,
					BAS:                  7,
					VSCHWELL:             0,
					DAYL:                 -15.0,
					DLBAS:                -23.5,
					DRYSWELL:             0.8,
					LUKRIT:               0.08,
					LAIFKT:               0.002,
					WGMAX:                0.01,
					PRO: []float64{
						0.130,
						0.330,
						0.540,
						0.000,
					},
					DEAD: []float64{
						0.000,
						0.000,
						0.000,
						0.000,
					},
					Kc: 1.00,
				},
				{
					DevelopmentStageName: "development phase 4: flower appearance to first pod",
					TSUM:                 280,
					BAS:                  7,
					VSCHWELL:             0,
					DAYL:                 0,
					DLBAS:                0,
					DRYSWELL:             0.8,
					LUKRIT:               0.08,
					LAIFKT:               0.002,
					WGMAX:                0.01,
					PRO: []float64{
						0.100,
						0.400,
						0.500,
						0.000,
					},
					DEAD: []float64{
						0.000,
						0.000,
						0.000,
						0.000,
					},
					Kc: 1.15,
				},
				{
					DevelopmentStageName: "development phase 5: first pod to last pod",
					TSUM:                 172,
					BAS:                  7,
					VSCHWELL:             0,
					DAYL:                 0,
					DLBAS:                0,
					DRYSWELL:             0.8,
					LUKRIT:               0.08,
					LAIFKT:               0.002,
					WGMAX:                0.01,
					PRO: []float64{
						0.000,
						0.000,
						0.400,
						0.600,
					},
					DEAD: []float64{
						0.000,
						0.050,
						0.000,
						0.000,
					},
					Kc: 1.15,
				},
				{
					DevelopmentStageName: "development phase 6: last pod to maturity",
					TSUM:                 400,
					BAS:                  4,
					VSCHWELL:             0,
					DAYL:                 0,
					DLBAS:                0,
					DRYSWELL:             0.7,
					LUKRIT:               0.08,
					LAIFKT:               0.002,
					WGMAX:                0.01,
					PRO: []float64{
						0.000,
						0.000,
						0.000,
						1.000,
					},
					DEAD: []float64{
						0.000,
						0.050,
						0.000,
						0.000,
					},
					Kc: 1.15,
				},
				{
					DevelopmentStageName: "development phase 7: senescence",
					TSUM:                 125,
					BAS:                  9,
					VSCHWELL:             0,
					DAYL:                 0,
					DLBAS:                0,
					DRYSWELL:             0.8,
					LUKRIT:               0.08,
					LAIFKT:               0.002,
					WGMAX:                0.01,
					PRO: []float64{
						0.000,
						0.000,
						0.000,
						0.000,
					},
					DEAD: []float64{
						0.000,
						0.050,
						0.000,
						0.000,
					},
					Kc: 0.50,
				},
			},
		}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := WriteCropParam(tt.args.filename, tt.args.cropParam); (err != nil) != tt.wantErr {
				t.Errorf("WriteCropParam() error = %v, wantErr %v", err, tt.wantErr)
			}
			// read back the file and compare
			cropParam, err := ReadCropParamFromFile(tt.args.filename)
			if err != nil {
				t.Errorf("ReadCropParam() error = %v", err)
			}
			if cropParam.CropName != tt.args.cropParam.CropName {
				t.Errorf("ReadCropParam() CropName = %v, want %v", cropParam.CropName, tt.args.cropParam.CropName)
			}
		})
	}
}

func TestConvertCropParamClassicToYml(t *testing.T) {
	type args struct {
		PARANAM string
	}
	tests := []struct {
		name    string
		args    args
		want    CropParam
		wantErr bool
	}{
		{"testSoy", args{"test_data/PARAM_test.SOY"}, CropParam{
			CropName:          "soybean",
			ABBr:              "SOY",
			Variety:           "test",
			MAXAMAX:           28,
			TempTyp:           1,
			MINTMP:            6,
			WUMAXPF:           9,
			VELOC:             0.5574,
			NGEFKT:            1,
			RGA:               0,
			RGB:               0,
			SubOrgan:          0,
			AboveGroundOrgans: []int{2, 3, 4},
			YORGAN:            4,
			YIFAK:             0.8,
			INITCONCNBIOM:     6.0,
			INITCONCNROOT:     2.0,
			NRKOM:             4,
			CompartimentNames: []string{"root", "leave", "stem", "ears"},
			DAUERKULT:         false,
			LEGUM:             true,
			WORG:              []float64{153, 153, 0, 0},
			MAIRT:             []float64{0.01, 0.03, 0.015, 0.01},
			KcIni:             0.65,
			NRENTW:            7,
			CropDevelopmentStages: []CropDevelopmentStage{
				{
					DevelopmentStageName: "development phase 1: sowing til emergence",
					TSUM:                 82,
					BAS:                  7,
					VSCHWELL:             0,
					DAYL:                 0,
					DLBAS:                0,
					DRYSWELL:             1,
					LUKRIT:               0.08,
					LAIFKT:               0.0025,
					WGMAX:                0.02,
					PRO: []float64{
						0.5,
						0.5,
						0,
						0,
					},
					DEAD: []float64{
						0, 0, 0, 0,
					},
					Kc: 0.65,
				},
				{
					DevelopmentStageName: "development phase 2: emergence to end juvenile phase",
					TSUM:                 50,
					BAS:                  7,
					VSCHWELL:             0,
					DAYL:                 -15.0,
					DLBAS:                -23.5,
					DRYSWELL:             0.6,
					LUKRIT:               0.08,
					LAIFKT:               0.0030,
					WGMAX:                0.012,
					PRO: []float64{
						0.200,
						0.600,
						0.200,
						0.000,
					},
					DEAD: []float64{
						0.000,
						0.000,
						0.000,
						0.000,
					},
					Kc: 0.7,
				},
				{
					DevelopmentStageName: "development phase 3: end juvenile phase to flower appearance",
					TSUM:                 220,
					BAS:                  7,
					VSCHWELL:             0,
					DAYL:                 -15.0,
					DLBAS:                -23.5,
					DRYSWELL:             0.8,
					LUKRIT:               0.08,
					LAIFKT:               0.002,
					WGMAX:                0.01,
					PRO: []float64{
						0.130,
						0.330,
						0.540,
						0.000,
					},
					DEAD: []float64{
						0.000,
						0.000,
						0.000,
						0.000,
					},
					Kc: 1.00,
				},
				{
					DevelopmentStageName: "development phase 4: flower appearance to first pod",
					TSUM:                 280,
					BAS:                  7,
					VSCHWELL:             0,
					DAYL:                 0,
					DLBAS:                0,
					DRYSWELL:             0.8,
					LUKRIT:               0.08,
					LAIFKT:               0.002,
					WGMAX:                0.01,
					PRO: []float64{
						0.100,
						0.400,
						0.500,
						0.000,
					},
					DEAD: []float64{
						0.000,
						0.000,
						0.000,
						0.000,
					},
					Kc: 1.15,
				},
				{
					DevelopmentStageName: "development phase 5: first pod to last pod",
					TSUM:                 172,
					BAS:                  7,
					VSCHWELL:             0,
					DAYL:                 0,
					DLBAS:                0,
					DRYSWELL:             0.8,
					LUKRIT:               0.08,
					LAIFKT:               0.002,
					WGMAX:                0.01,
					PRO: []float64{
						0.000,
						0.000,
						0.400,
						0.600,
					},
					DEAD: []float64{
						0.000,
						0.050,
						0.000,
						0.000,
					},
					Kc: 1.15,
				},
				{
					DevelopmentStageName: "development phase 6: last pod to maturity",
					TSUM:                 400,
					BAS:                  4,
					VSCHWELL:             0,
					DAYL:                 0,
					DLBAS:                0,
					DRYSWELL:             0.7,
					LUKRIT:               0.08,
					LAIFKT:               0.002,
					WGMAX:                0.01,
					PRO: []float64{
						0.000,
						0.000,
						0.000,
						1.000,
					},
					DEAD: []float64{
						0.000,
						0.050,
						0.000,
						0.000,
					},
					Kc: 1.15,
				},
				{
					DevelopmentStageName: "development phase 7: senescence",
					TSUM:                 125,
					BAS:                  9,
					VSCHWELL:             0,
					DAYL:                 0,
					DLBAS:                0,
					DRYSWELL:             0.8,
					LUKRIT:               0.08,
					LAIFKT:               0.002,
					WGMAX:                0.01,
					PRO: []float64{
						0.000,
						0.000,
						0.000,
						0.000,
					},
					DEAD: []float64{
						0.000,
						0.050,
						0.000,
						0.000,
					},
					Kc: 0.50,
				},
			},
		}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ConvertCropParamClassicToYml(tt.args.PARANAM)
			if (err != nil) != tt.wantErr {
				t.Errorf("ConvertCropParamClassicToYml() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err := WriteCropParam(tt.args.PARANAM+".yml", got); (err != nil) != tt.wantErr {
				t.Errorf("WriteCropParam() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ConvertCropParamClassicToYml() = %v, want %v", got, tt.want)
			}
		})
	}
}
