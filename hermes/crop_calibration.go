package hermes

import (
	"fmt"
	"path/filepath"
	"strings"
)

// Overwrite crop parameters from command line
// in order to calibrate crop parameters for a specific crop file
// use the following command line parameters:
// Example:
// CropFile=PARAM_0.SOY c_TSum_1=73 c_TSum_2=55 c_TSum_3=240 c_TSum_4=330
// CropFile=PARAM_0.SOY c_MAXAMAX=80 c_MINTMP=5 c_WUMAXPF=10 c_VELOC=0.005 c_YIFAK=0.8

// names of attributes relevant to crop overwriting

// MAXAMAX  // Amax (C-Assimilation bei Lichtsättigung) bei Optimaltemperatur (kg CO2/ha leave/h)
// MINTMP 	// Minimumtemperatur für Wachstum (°C)
// WUMAXPF 	// Pflanzenspezifische effektive Durchwurzelungstiefe (dm)
// VELOC (RTVELOC/200) // root depth increase in mm/C°
// YIFAK 	// fraction of organ.(organ 4 80% =4.80)
// INITCONCNBIOM // initial N concentration in above ground biomass (%)
// INITCONCNROOT // initial N concentration in roots (%)

// in stages:

// TSUM[i] 		// Temperatursumme für Entwicklungsstufe i (°C)
// BAS[i] 		// Basiswert für Temperatursumme (°C)
// VSCHWELL[i] 	// Benötigte Anzahl Vernalisationstage Entwicklungsstufe I
// DAYL[i] 		// Tageslänge für Entwicklungsstufe i (h)
// DLBAS[i] 	// Basiswert für Tageslänge (h)
// DRYSWELL[i] 	// Schwelle für Trockenstress (Ta/Tp) Entwicklungsstufe I (0-1)
// LUKRIT[i] 	// kritischer Luftporenanteil Entwicklungsstufe I (cm^3/cm^3)
// LAIFKT[i] 	// SLA specific leave area (area per mass) (m2/m2/kg TM) in I
// WGMAX[i] 	// N-content root end of phase I
// KC[i] 		// crop factor for evapotranspiration (0-1)

// in partitions:
// PRO[i][j] 	// fraction of production for organ j in stage i (0-1)
// DEAD[i][j] 	// fraction of dead material for organ j in stage i (0-1)

// check if parameter is supported
func isValidCropParameter(param string) bool {
	if param == "MAXAMAX" || param == "MINTMP" || param == "WUMAXPF" || param == "VELOC" || param == "YIFAK" ||
		param == "TSUM" || param == "BAS" || param == "VSCHWELL" || param == "DAYL" || param == "DLBAS" ||
		param == "DRYSWELL" || param == "LUKRIT" || param == "LAIFKT" || param == "WGMAX" || param == "KC" ||
		param == "INITCONCNBIOM" || param == "INITCONCNROOT" || param == "PRO" || param == "DEAD" {
		return true
	}
	return false
}

type CropOverwrite struct {
	CropFile                   string
	BaseFloatParameters        map[string]float64
	DevelopmentStageParameters map[string]map[int]float64
	PartitioningParameters     map[string]map[PartPair]float64
}
type PartPair struct {
	Stage int
	Part  int
}

// overwrite crop parameters in global vars
func (cropOW *CropOverwrite) OverwriteCropParameters(cropFile string, g *GlobalVarsMain, l *CropSharedVars) {
	if cropOW.CropFile != filepath.Base(cropFile) {
		return
	}
	if !cropOW.isValidCropOverwrite(g.NRKOM, l.NRENTW) {
		return
	}
	// overwrite base parameters
	for key, value := range cropOW.BaseFloatParameters {
		if key == "MAXAMAX" {
			g.MAXAMAX = value
		} else if key == "MINTMP" {
			g.MINTMP = value
		} else if key == "WUMAXPF" {
			g.WUMAXPF = value
		} else if key == "VELOC" {
			g.VELOC = value / 200
		} else if key == "YIFAK" {
			g.YIFAK = value
		}

		if key == "INITCONCNBIOM" {
			// apply only in case of initial crop
			if !(g.DAUERKULT && g.AKF.Num > 2 && g.FRUCHT[g.AKF.Index] == g.FRUCHT[g.AKF.Index-1]) {
				g.GEHOB = value / 100
			}
		}
		if key == "INITCONCNROOT" {
			// apply only in case of initial crop
			if !(g.DAUERKULT && g.AKF.Num > 2 && g.FRUCHT[g.AKF.Index] == g.FRUCHT[g.AKF.Index-1]) {
				g.WUGEH = value / 100
			}
		}
	}
	// overwrite development stage parameters
	for key, stages := range cropOW.DevelopmentStageParameters {
		if key == "TSUM" {
			for stage, value := range stages {
				stageIdx := stage - 1
				g.TSUM[stageIdx] = value
			}
		} else if key == "BAS" {
			for stage, value := range stages {
				stageIdx := stage - 1
				g.BAS[stageIdx] = value
			}
		} else if key == "VSCHWELL" {
			for stage, value := range stages {
				stageIdx := stage - 1
				g.VSCHWELL[stageIdx] = value
			}
		} else if key == "DAYL" {
			for stage, value := range stages {
				stageIdx := stage - 1
				g.DAYL[stageIdx] = value
			}
		} else if key == "DLBAS" {
			for stage, value := range stages {
				stageIdx := stage - 1
				g.DLBAS[stageIdx] = value
			}
		} else if key == "DRYSWELL" {
			for stage, value := range stages {
				stageIdx := stage - 1
				g.DRYSWELL[stageIdx] = value
			}
		} else if key == "LUKRIT" {
			for stage, value := range stages {
				stageIdx := stage - 1
				g.LUKRIT[stageIdx] = value
			}
		} else if key == "LAIFKT" {
			for stage, value := range stages {
				stageIdx := stage - 1
				g.LAIFKT[stageIdx] = value
			}
		} else if key == "WGMAX" {
			for stage, value := range stages {
				stageIdx := stage - 1
				g.WGMAX[stageIdx] = value
			}
		} else if key == "KC" {
			for stage, value := range stages {
				stageIdx := stage - 1
				l.kc[stageIdx] = value
			}
		}

	}
	// overwrite partitioning parameters
	for key, parts := range cropOW.PartitioningParameters {
		if key == "PRO" {
			for part, value := range parts {
				stageIdx := part.Stage - 1
				partIdx := part.Part - 1
				g.PRO[stageIdx][partIdx] = value
			}
		} else if key == "DEAD" {
			for part, value := range parts {
				stageIdx := part.Stage - 1
				partIdx := part.Part - 1
				g.DEAD[stageIdx][partIdx] = value
			}
		}
	}
}

// check if crop overwrite parameters are in valid range
func (cropOW *CropOverwrite) isValidCropOverwrite(numPartitions, numStages int) bool {
	if cropOW.CropFile == "" {
		return false
	}
	for key, value := range cropOW.BaseFloatParameters {
		if key == "MAXAMAX" && (value <= 0 || value > 100) {
			return false
		} else if key == "MINTMP" && (value <= -30 || value >= 50) {
			return false
		} else if key == "WUMAXPF" && (value <= 0 || value > 20) {
			return false
		} else if key == "VELOC" && (value <= 0 || value > 1) {
			return false
		} else if key == "YIFAK" && (value < 0 || value > 1) {
			return false
		} else if key == "INITCONCNBIOM" && (value < 0 || value > 100) {
			return false
		} else if key == "INITCONCNROOT" && (value < 0 || value > 100) {
			return false
		}
	}
	for key, stages := range cropOW.DevelopmentStageParameters {
		// check if stage is within valid range
		for stage := range stages {
			if stage < 1 || stage > numStages {
				return false
			}
		}
		if key == "TSUM" {
			for _, value := range stages {
				if value < 0 || value > 10000 {
					return false
				}
			}
		} else if key == "BAS" {
			for _, value := range stages {
				if value < -10 || value > 40 {
					return false
				}
			}
		} else if key == "VSCHWELL" {
			for _, value := range stages {
				if value < 0 || value > 100 {
					return false
				}
			}
		} else if key == "DAYL" {
			for _, value := range stages {
				if value < 24 || value > 24 {
					return false
				}
			}
		} else if key == "DLBAS" {
			for _, value := range stages {
				if value < 24 || value > 24 {
					return false
				}
			}
		} else if key == "DRYSWELL" {
			for _, value := range stages {
				if value < 0 || value > 1 {
					return false
				}
			}
		} else if key == "LUKRIT" {
			for _, value := range stages {
				if value < 0 || value > 1 {
					return false
				}
			}
		} else if key == "LAIFKT" {
			for _, value := range stages {
				if value < 0 || value > 100 {
					return false
				}
			}
		} else if key == "WGMAX" {
			for _, value := range stages {
				if value < 0 || value > 100 {
					return false
				}
			}
		} else if key == "KC" {
			for _, value := range stages {
				if value < 0 || value > 1 {
					return false
				}
			}
		} else {
			return false
		}
	}
	for key, parts := range cropOW.PartitioningParameters {
		// check if stage and partition are within valid range
		for pair := range parts {
			if pair.Stage < 1 || pair.Stage > numStages {
				return false
			}
			if pair.Part < 1 || pair.Part > numPartitions {
				return false
			}
		}

		if key == "PRO" {
			for _, value := range parts {
				if value < 0 || value > 1 {
					return false
				}
			}
		} else if key == "DEAD" {
			for _, value := range parts {
				if value < 0 || value > 1 {
					return false
				}
			}
		} else {
			return false
		}
	}

	return true
}

// crop overwrite parameters from command line
func ParseCropOverwrites(args map[string]string) (*CropOverwrite, error) {
	// args format:
	// CropFile=PARAM_0.SOY c_TSum_1=73 c_TSum_2=55 c_TSum_3=240 c_TSum_4=330

	// check if args contain a CropFile=... entry
	if _, ok := args["CropFile"]; !ok {
		// no CropFile entry, skip
		return nil, nil
	}
	cropOW := &CropOverwrite{
		CropFile:                   args["CropFile"],
		BaseFloatParameters:        map[string]float64{},
		DevelopmentStageParameters: map[string]map[int]float64{},
		PartitioningParameters:     map[string]map[PartPair]float64{},
	}

	for key, arg := range args {

		if strings.HasPrefix(key, "c_") {
			// convert value to float
			paramValue := ValAsFloat(arg, "cmd line", arg)

			// split by _
			paramSplit := strings.Split(key, "_")
			// extract parameter name
			paramName := paramSplit[1]
			if !isValidCropParameter(paramName) {
				// fail if parameter name is not valid
				return nil, fmt.Errorf("invalid crop parameter name: " + paramName)
			}
			// check if param is a base parameter
			if len(paramSplit) == 2 {
				// base parameter
				// add parameter to map
				cropOW.BaseFloatParameters[paramName] = paramValue
			} else if len(paramSplit) == 3 {
				// development stage parameter
				// extract development stage
				developmentStage := int(ValAsInt(paramSplit[2], "cmd line", arg))
				if developmentStage < 1 || developmentStage > 9 {
					// fail if development stage is not valid
					return nil, fmt.Errorf("invalid development stage index: " + paramSplit[2])
				}
				// add parameter to map
				if _, ok := cropOW.DevelopmentStageParameters[paramName]; !ok {
					cropOW.DevelopmentStageParameters[paramName] = map[int]float64{}
				}
				cropOW.DevelopmentStageParameters[paramName][developmentStage] = paramValue
			} else if len(paramSplit) == 4 {
				// partitioning parameter
				// extract development stage
				developmentStage := int(ValAsInt(paramSplit[2], "cmd line", arg))
				if developmentStage < 1 || developmentStage > 9 {
					// fail if development stage is not valid
					return nil, fmt.Errorf("invalid development stage index: " + paramSplit[2])
				}
				// extract partition
				partition := int(ValAsInt(paramSplit[3], "cmd line", arg))
				if partition < 1 || partition > 5 {
					// fail if partition is not valid
					return nil, fmt.Errorf("invalid partition index: " + paramSplit[3])
				}
				// add parameter to map
				if _, ok := cropOW.PartitioningParameters[paramName]; !ok {
					cropOW.PartitioningParameters[paramName] = map[PartPair]float64{}
				}
				cropOW.PartitioningParameters[paramName][PartPair{Stage: developmentStage, Part: partition}] = paramValue
			} else {
				// fail if parameter name is not valid
				return nil, fmt.Errorf("invalid crop parameter name: " + paramName)
			}
		}
	}
	return cropOW, nil
}
