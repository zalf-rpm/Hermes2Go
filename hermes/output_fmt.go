package hermes

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
	"unicode/utf8"

	yaml "gopkg.in/yaml.v2"
)

// OutputConfig defines output parameters in form of a csv file
type OutputConfig struct {
	numHeadLines       int
	numDataColumns     int
	FillCharacter      string `yaml:"FillCharacter"`
	fillRune           rune
	SeperatorCharacter string `yaml:"SeperatorCharacter"`
	seperatorRune      rune
	NotAvailableValue  string                   `yaml:"NaValue"`
	DataColumns        []OutputDataColum        `yaml:"DataColumns"`
	Headlines          map[int][]OutHeaderColum `yaml:"Headlines"`
}

// OutputDataColum describes data format and reference variable
type OutputDataColum struct {
	FormatStr     string    `yaml:"Format"`
	DataAlignment Alignment `yaml:"DataAlignment"`
	Width         int       `yaml:"Width"`
	Modifier      float64   `yaml:"Modifier,omitempty"`
	VarName       string    `yaml:"VariableName"`
	VarIndex1     int       `yaml:"VarIndex1,omitempty"`
	VarIndex2     int       `yaml:"VarIndex2,omitempty"`
	valueRef      interface{}
}

// OutHeaderColum describes header text, format and position
type OutHeaderColum struct {
	Text              string    `yaml:"ColumnName"`
	ColumnAlignment   Alignment `yaml:"TextAlignment"`
	ColStart          int       `yaml:"StartColumn,omitempty"`
	ColEnd            int       `yaml:"EndColumn,omitempty"`
	fillWithRune      rune
	FillWithCharacter string `yaml:"FillCharacter,omitempty"`
}

// WriteHeader of an output file
func (c *OutputConfig) WriteHeader(file *Fout, formatType OutputFileFormat) error {
	arrStartIndex := make([]int, len(c.DataColumns)+1)
	for i, col := range c.DataColumns {
		arrStartIndex[i+1] = arrStartIndex[i] + col.Width + 1
	}
	for i := 1; i < c.numHeadLines+1; i++ {
		fillerRune := c.fillRune
		if columnC, ok := c.Headlines[i]; ok {
			currentIndex := 0
			lastIndex := 0
			for idxCol, col := range columnC {
				if formatType == csvOut {
					_, err := file.Write(col.Text)
					if err != nil {
						return err
					}
					if idxCol < len(columnC)-1 {
						_, err = file.WriteRune(c.seperatorRune)
						if err != nil {
							return err
						}
					}
				} else if formatType == hermesOut {
					lastIndex = arrStartIndex[col.ColEnd]
					// move cursor to next column start index
					for currentIndex < arrStartIndex[col.ColStart-1]-1 {
						currentIndex++
						_, err := file.WriteRune(fillerRune)
						if err != nil {
							return err
						}
					}
					//reset fill rune
					fillerRune = c.fillRune
					for currentIndex < arrStartIndex[col.ColStart-1] {
						currentIndex++
						_, err := file.WriteRune(fillerRune)
						if err != nil {
							return err
						}
					}
					// set col fill rune
					fillerRune = col.fillWithRune
					if col.ColumnAlignment == leftAlignment {
						numRunesWritten := 0
						runesToWrite := arrStartIndex[col.ColEnd] - arrStartIndex[col.ColStart-1] - 1
						for _, r := range col.Text {
							if numRunesWritten < runesToWrite {
								_, err := file.WriteRune(r)
								if err != nil {
									return err
								}
								numRunesWritten++
								currentIndex++
							} else {
								break
							}
						}
					} else if col.ColumnAlignment == rightAlignment {
						runesToWrite := arrStartIndex[col.ColEnd] - arrStartIndex[col.ColStart-1] - 1
						runesInStr := utf8.RuneCountInString(col.Text)
						numfillRunes := runesToWrite - runesInStr
						for numfillRunes > 0 {
							numfillRunes--
							_, err := file.WriteRune(col.fillWithRune)
							if err != nil {
								return err
							}
							runesToWrite--
							currentIndex++
						}
						for _, r := range col.Text {
							if runesToWrite > 0 {
								_, err := file.WriteRune(r)
								if err != nil {
									return err
								}
								runesToWrite--
								currentIndex++
							} else {
								break
							}
						}
					} else if col.ColumnAlignment == centerAlignment || col.ColumnAlignment == noneAlignment {
						runesToWrite := arrStartIndex[col.ColEnd] - arrStartIndex[col.ColStart-1] - 1
						runesInStr := utf8.RuneCountInString(col.Text)
						numfillRunes := (runesToWrite - runesInStr) / 2
						for numfillRunes > 0 {
							numfillRunes--
							_, err := file.WriteRune(col.fillWithRune)
							if err != nil {
								return err
							}
							currentIndex++
							runesToWrite--
						}
						for _, r := range col.Text {
							if runesToWrite > 0 {
								_, err := file.WriteRune(r)
								if err != nil {
									return err
								}
								runesToWrite--
								currentIndex++
							} else {
								break
							}
						}
					}
				}
			}
			if formatType == hermesOut {
				// move cursor to next column start index
				for currentIndex < lastIndex {
					currentIndex++
					_, err := file.WriteRune(fillerRune)
					if err != nil {
						return err
					}
				}
			}
		}
		file.Write("\r\n")
	}

	return nil
}

// NewDefaultOutputConfigYearly create yearly hermes output configuration
func NewDefaultOutputConfigYearly(g *GlobalVarsMain) OutputConfig {
	dataColumns := []OutputDataColum{
		OutputDataColum{
			FormatStr:     "%8s",
			DataAlignment: leftAlignment,
			Width:         10,
			VarName:       "AKTUELL",
			valueRef:      &g.AKTUELL,
		},
		OutputDataColum{
			FormatStr:     "%04.f",
			DataAlignment: rightAlignment,
			Width:         4,
			Modifier:      10.0,
			VarName:       "VERDUNST",
			valueRef:      &g.VERDUNST,
		},
		OutputDataColum{
			FormatStr:     "%04.f",
			DataAlignment: rightAlignment,
			Width:         4,
			Modifier:      10.0,
			VarName:       "PFTRANS",
			valueRef:      &g.PFTRANS,
		},
		OutputDataColum{
			FormatStr:     "%04.f",
			DataAlignment: rightAlignment,
			Width:         4,
			Modifier:      10.0,
			VarName:       "TRAY",
			valueRef:      &g.TRAY,
		},
		OutputDataColum{
			FormatStr:     "%04.f",
			DataAlignment: rightAlignment,
			Width:         4,
			VarName:       "PerY",
			valueRef:      &g.PerY,
		},
		OutputDataColum{
			FormatStr:     "%04.f",
			DataAlignment: leftAlignment,
			Width:         5,
			VarName:       "SWCY1",
			valueRef:      &g.SWCY1,
		},
		OutputDataColum{
			FormatStr:     "%04.f",
			DataAlignment: leftAlignment,
			Width:         5,
			VarName:       "SWCY2",
			valueRef:      &g.SWCY2,
		},
		OutputDataColum{
			FormatStr:     "%s",
			DataAlignment: leftAlignment,
			Width:         6,
			VarName:       "NA",
			valueRef:      "n.a.",
		},
		OutputDataColum{
			FormatStr:     "%06.1f",
			DataAlignment: leftAlignment,
			Width:         7,
			VarName:       "OUTSUM",
			valueRef:      &g.OUTSUM,
		},
		OutputDataColum{
			FormatStr:     "%07.1f",
			DataAlignment: rightAlignment,
			Width:         7,
			VarName:       "MINSUM",
			valueRef:      &g.MINSUM,
		},
		OutputDataColum{
			FormatStr:     "%07.3f",
			DataAlignment: rightAlignment,
			Width:         7,
			VarName:       "CUMDENIT",
			valueRef:      &g.CUMDENIT,
		},
		OutputDataColum{
			FormatStr:     "%s",
			DataAlignment: leftAlignment,
			Width:         7,
			VarName:       "NA",
			valueRef:      "n.a.",
		},
		OutputDataColum{
			FormatStr:     "%07.f",
			DataAlignment: leftAlignment,
			Width:         7,
			VarName:       "SOC1",
			valueRef:      &g.SOC1,
		},
		OutputDataColum{
			FormatStr:     "%s",
			DataAlignment: leftAlignment,
			Width:         7,
			VarName:       "NA",
			valueRef:      "n.a.",
		},
		OutputDataColum{
			FormatStr:     "%s",
			DataAlignment: leftAlignment,
			Width:         21,
			VarName:       "POLYD",
			valueRef:      &g.POLYD,
		},
		OutputDataColum{
			FormatStr:     "%s",
			DataAlignment: leftAlignment,
			Width:         11,
			VarName:       "C1NotStableErr",
			valueRef:      &g.C1NotStableErr,
		},
	}
	headlines := map[int][]OutHeaderColum{
		1: []OutHeaderColum{
			OutHeaderColum{
				Text:            "date",
				ColumnAlignment: leftAlignment,
			},
			OutHeaderColum{
				Text:            "ETcy",
				ColumnAlignment: rightAlignment,
			},
			OutHeaderColum{
				Text:            "ETaY",
				ColumnAlignment: rightAlignment,
			},
			OutHeaderColum{
				Text:            "TraY",
				ColumnAlignment: rightAlignment,
			},
			OutHeaderColum{
				Text:            "PerY",
				ColumnAlignment: rightAlignment,
			},
			OutHeaderColum{
				Text:            "SWCY1",
				ColumnAlignment: leftAlignment,
			},
			OutHeaderColum{
				Text:            "SWCY2",
				ColumnAlignment: leftAlignment,
			},
			OutHeaderColum{
				Text:            "Runoff",
				ColumnAlignment: rightAlignment,
			},
			OutHeaderColum{
				Text:            "NleaY",
				ColumnAlignment: leftAlignment,
			},
			OutHeaderColum{
				Text:            "MINY",
				ColumnAlignment: leftAlignment,
			},
			OutHeaderColum{
				Text:            "DENY",
				ColumnAlignment: leftAlignment,
			},
			OutHeaderColum{
				Text:            "VOLAT",
				ColumnAlignment: leftAlignment,
			},
			OutHeaderColum{
				Text:            "SOC1",
				ColumnAlignment: leftAlignment,
			},
			OutHeaderColum{
				Text:            "SOC2",
				ColumnAlignment: leftAlignment,
			},
			OutHeaderColum{
				Text:            "code",
				ColumnAlignment: leftAlignment,
			},
			OutHeaderColum{
				Text:            "Error",
				ColumnAlignment: leftAlignment,
			},
		},
		2: []OutHeaderColum{
			OutHeaderColum{
				Text:            "mm",
				ColumnAlignment: centerAlignment,
				ColStart:        2,
			},
			OutHeaderColum{
				Text:            "mm",
				ColumnAlignment: centerAlignment,
				ColStart:        3,
			},
			OutHeaderColum{
				Text:            "mm",
				ColumnAlignment: centerAlignment,
				ColStart:        4,
			},
			OutHeaderColum{
				Text:            "mm",
				ColumnAlignment: centerAlignment,
				ColStart:        5,
			},
			OutHeaderColum{
				Text:            "mm",
				ColumnAlignment: centerAlignment,
				ColStart:        6,
			},
			OutHeaderColum{
				Text:            "mm",
				ColumnAlignment: centerAlignment,
				ColStart:        7,
			},
			OutHeaderColum{
				Text:            "mm",
				ColumnAlignment: centerAlignment,
				ColStart:        8,
			},
			OutHeaderColum{
				Text:            "kg N/ha",
				ColumnAlignment: centerAlignment,
				ColStart:        9,
			},
			OutHeaderColum{
				Text:            "kg N/ha",
				ColumnAlignment: centerAlignment,
				ColStart:        10,
			},
			OutHeaderColum{
				Text:            "kg N/ha",
				ColumnAlignment: centerAlignment,
				ColStart:        11,
			},
			OutHeaderColum{
				Text:            "kg N/ha",
				ColumnAlignment: centerAlignment,
				ColStart:        12,
			},
			OutHeaderColum{
				Text:            "kg C/ha",
				ColumnAlignment: centerAlignment,
				ColStart:        13,
			},
			OutHeaderColum{
				Text:            "kg C/ha",
				ColumnAlignment: centerAlignment,
				ColStart:        14,
			},
		},
	}
	for i := range headlines[1] {
		headlines[1][i].fillWithRune = ' '
		headlines[1][i].FillWithCharacter = " "
		headlines[1][i].ColStart = i + 1
		headlines[1][i].ColEnd = i + 1
	}

	for i := range headlines[2] {
		headlines[2][i].fillWithRune = ' '
		headlines[2][i].FillWithCharacter = " "
		if headlines[2][i].ColEnd == 0 {
			headlines[2][i].ColEnd = headlines[2][i].ColStart
		}
	}

	return OutputConfig{
		numHeadLines:       len(headlines),
		numDataColumns:     len(dataColumns),
		Headlines:          headlines,
		DataColumns:        dataColumns,
		FillCharacter:      " ",
		fillRune:           ' ',
		SeperatorCharacter: ",",
		seperatorRune:      ',',
		NotAvailableValue:  "n.a.",
	}

}

// NewDefaultCropOutputConfig create crop output configuration
func NewDefaultCropOutputConfig(c *CropOutputVars) OutputConfig {
	dataColumns := []OutputDataColum{
		OutputDataColum{
			FormatStr:     "%s",
			DataAlignment: leftAlignment,
			Width:         12,
			VarName:       "SowDate",
			valueRef:      &c.SowDate,
		},
		OutputDataColum{
			FormatStr:     "%3d",
			DataAlignment: leftAlignment,
			Width:         4,
			VarName:       "SowDOY",
			valueRef:      &c.SowDOY,
		},
		OutputDataColum{
			FormatStr:     "%3d",
			DataAlignment: leftAlignment,
			Width:         5,
			VarName:       "EmergDOY",
			valueRef:      &c.EmergDOY,
		},
		OutputDataColum{
			FormatStr:     "%3d",
			DataAlignment: leftAlignment,
			Width:         5,
			VarName:       "AnthDOY",
			valueRef:      &c.AnthDOY,
		},
		OutputDataColum{
			FormatStr:     "%3d",
			DataAlignment: leftAlignment,
			Width:         3,
			VarName:       "MatDOY",
			valueRef:      &c.MatDOY,
		},
		OutputDataColum{
			FormatStr:     "%4d",
			DataAlignment: leftAlignment,
			Width:         4,
			VarName:       "HarvestYear",
			valueRef:      &c.HarvestYear,
		},
		OutputDataColum{
			FormatStr:     "%3d",
			DataAlignment: leftAlignment,
			Width:         4,
			VarName:       "HarvestDOY",
			valueRef:      &c.HarvestDOY,
		},
		OutputDataColum{
			FormatStr:     "%3s",
			DataAlignment: leftAlignment,
			Width:         5,
			VarName:       "Crop",
			valueRef:      &c.Crop,
		},
		OutputDataColum{
			FormatStr:     "%05.f",
			DataAlignment: leftAlignment,
			Width:         6,
			VarName:       "Yield",
			valueRef:      &c.Yield,
		},
		OutputDataColum{
			FormatStr:     "%05.f",
			DataAlignment: leftAlignment,
			Width:         7,
			VarName:       "Biomass",
			valueRef:      &c.Biomass,
		},
		OutputDataColum{
			FormatStr:     "%05.f",
			DataAlignment: leftAlignment,
			Width:         5,
			VarName:       "Roots",
			valueRef:      &c.Roots,
		},
		OutputDataColum{
			FormatStr:     "%04.1f",
			DataAlignment: leftAlignment,
			Width:         6,
			VarName:       "LAImax",
			valueRef:      &c.LAImax,
		},

		OutputDataColum{
			FormatStr:     "%05.1f",
			DataAlignment: leftAlignment,
			Width:         7,
			VarName:       "Nfertil",
			valueRef:      &c.Nfertil,
		},
		OutputDataColum{
			FormatStr:     "%04.f",
			DataAlignment: leftAlignment,
			Width:         5,
			VarName:       "Irrig",
			valueRef:      &c.Irrig,
		},
		OutputDataColum{
			FormatStr:     "%05.1f",
			DataAlignment: leftAlignment,
			Width:         8,
			VarName:       "Nuptake",
			valueRef:      &c.Nuptake,
		},
		OutputDataColum{
			FormatStr:     "%05.1f",
			DataAlignment: leftAlignment,
			Width:         7,
			VarName:       "Nagb",
			valueRef:      &c.Nagb,
		},
		OutputDataColum{
			FormatStr:     "%04.f",
			DataAlignment: leftAlignment,
			Width:         4,
			Modifier:      10,
			VarName:       "ETcG",
			valueRef:      &c.ETcG,
		},
		OutputDataColum{
			FormatStr:     "%04.f",
			DataAlignment: leftAlignment,
			Width:         4,
			Modifier:      10,
			VarName:       "ETaG",
			valueRef:      &c.ETaG,
		},
		OutputDataColum{
			FormatStr:     "%04.f",
			DataAlignment: leftAlignment,
			Width:         4,
			Modifier:      10,
			VarName:       "TraG",
			valueRef:      &c.TraG,
		},
		OutputDataColum{
			FormatStr:     "%06.1f",
			DataAlignment: leftAlignment,
			Width:         6,
			VarName:       "PerG",
			valueRef:      &c.PerG,
		},

		OutputDataColum{
			FormatStr:     "%04.f",
			DataAlignment: leftAlignment,
			Width:         5,
			VarName:       "SWCS1",
			valueRef:      &c.SWCS1,
		},
		OutputDataColum{
			FormatStr:     "%04.f",
			DataAlignment: leftAlignment,
			Width:         5,
			VarName:       "SWCS2",
			valueRef:      &c.SWCS2,
		},
		OutputDataColum{
			FormatStr:     "%04.f",
			DataAlignment: leftAlignment,
			Width:         5,
			VarName:       "SWCA1",
			valueRef:      &c.SWCA1,
		},
		OutputDataColum{
			FormatStr:     "%04.f",
			DataAlignment: leftAlignment,
			Width:         5,
			VarName:       "SWCA2",
			valueRef:      &c.SWCA2,
		},
		OutputDataColum{
			FormatStr:     "%04.f",
			DataAlignment: leftAlignment,
			Width:         5,
			VarName:       "SWCM1",
			valueRef:      &c.SWCM1,
		},
		OutputDataColum{
			FormatStr:     "%04.f",
			DataAlignment: leftAlignment,
			Width:         5,
			VarName:       "SWCM2",
			valueRef:      &c.SWCM2,
		},

		OutputDataColum{
			FormatStr:     "%05.f",
			DataAlignment: leftAlignment,
			Width:         7,
			VarName:       "SoilN1",
			valueRef:      &c.SoilN1,
		},
		OutputDataColum{
			FormatStr:     "%s",
			DataAlignment: leftAlignment,
			Width:         7,
			VarName:       "NA",
			valueRef:      "n.a.",
		},
		OutputDataColum{
			FormatStr:     "%04.f",
			DataAlignment: leftAlignment,
			Width:         7,
			VarName:       "Nmin1",
			valueRef:      &c.Nmin1,
		},
		OutputDataColum{
			FormatStr:     "%04.f",
			DataAlignment: leftAlignment,
			Width:         7,
			VarName:       "Nmin2",
			valueRef:      &c.Nmin2,
		},

		OutputDataColum{
			FormatStr:     "%05.1f",
			DataAlignment: leftAlignment,
			Width:         7,
			VarName:       "NLeaG",
			valueRef:      &c.NLeaG,
		},
		OutputDataColum{
			FormatStr:     "%5.3f",
			DataAlignment: leftAlignment,
			Width:         5,
			VarName:       "TRRel",
			valueRef:      &c.TRRel,
		},
		OutputDataColum{
			FormatStr:     "%5.3f",
			DataAlignment: leftAlignment,
			Width:         5,
			VarName:       "Reduk",
			valueRef:      &c.Reduk,
		},
		OutputDataColum{
			FormatStr:     "%03.f",
			DataAlignment: leftAlignment,
			Width:         5,
			VarName:       "DryD1",
			valueRef:      &c.DryD1,
		},
		OutputDataColum{
			FormatStr:     "%03.f",
			DataAlignment: leftAlignment,
			Width:         5,
			VarName:       "DryD2",
			valueRef:      &c.DryD2,
		},
		OutputDataColum{
			FormatStr:     "%04.f",
			DataAlignment: leftAlignment,
			Width:         8,
			VarName:       "Nresid",
			valueRef:      &c.Nresid,
		},
		OutputDataColum{
			FormatStr:     "%s",
			DataAlignment: leftAlignment,
			Width:         10,
			VarName:       "Orgdat",
			valueRef:      &c.Orgdat,
		},
		OutputDataColum{
			FormatStr:     "%3s",
			DataAlignment: rightAlignment,
			Width:         4,
			VarName:       "Type",
			valueRef:      &c.Type,
		},

		OutputDataColum{
			FormatStr:     "%05.1f",
			DataAlignment: leftAlignment,
			Width:         5,
			VarName:       "OrgN",
			valueRef:      &c.OrgN,
		},
		OutputDataColum{
			FormatStr:     "%s",
			DataAlignment: leftAlignment,
			Width:         12,
			VarName:       "NDat1",
			valueRef:      &c.NDat1,
		},

		OutputDataColum{
			FormatStr:     "%05.1f",
			DataAlignment: leftAlignment,
			Width:         5,
			VarName:       "N1",
			valueRef:      &c.N1,
		},
		OutputDataColum{
			FormatStr:     "%s",
			DataAlignment: leftAlignment,
			Width:         12,
			VarName:       "Ndat2",
			valueRef:      &c.Ndat2,
		},

		OutputDataColum{
			FormatStr:     "%05.1f",
			DataAlignment: leftAlignment,
			Width:         5,
			VarName:       "N2",
			valueRef:      &c.N2,
		},
		OutputDataColum{
			FormatStr:     "%s",
			DataAlignment: leftAlignment,
			Width:         12,
			VarName:       "Ndat3",
			valueRef:      &c.Ndat3,
		},
		OutputDataColum{
			FormatStr:     "%05.1f",
			DataAlignment: leftAlignment,
			Width:         5,
			VarName:       "N3",
			valueRef:      &c.N3,
		},
		OutputDataColum{
			FormatStr:     "%s",
			DataAlignment: centerAlignment,
			Width:         12,
			VarName:       "Tdat",
			valueRef:      &c.Tdat,
		},
		OutputDataColum{
			FormatStr:     "%s",
			DataAlignment: centerAlignment,
			Width:         35,
			VarName:       "Code",
			valueRef:      &c.Code,
		},
		OutputDataColum{
			FormatStr:     "%s",
			DataAlignment: centerAlignment,
			Width:         10,
			VarName:       "NotStableErr",
			valueRef:      &c.NotStableErr,
		},
	}
	headlines := map[int][]OutHeaderColum{
		1: []OutHeaderColum{
			OutHeaderColum{
				Text:            "date",
				ColumnAlignment: leftAlignment,
			},
			OutHeaderColum{
				Text:            "DOY",
				ColumnAlignment: leftAlignment,
			},
			OutHeaderColum{
				Text:            "DOY",
				ColumnAlignment: leftAlignment,
			},
			OutHeaderColum{
				Text:            "DOY",
				ColumnAlignment: leftAlignment,
			},
			OutHeaderColum{
				Text:            "DOY",
				ColumnAlignment: leftAlignment,
			},
			OutHeaderColum{
				Text:            "Year",
				ColumnAlignment: leftAlignment,
			},
			OutHeaderColum{
				Text:            "DOY",
				ColumnAlignment: leftAlignment,
			},
			OutHeaderColum{
				Text:            "crop",
				ColumnAlignment: leftAlignment,
			},
			OutHeaderColum{
				Text:            "yield",
				ColumnAlignment: leftAlignment,
			},
			OutHeaderColum{
				Text:            "biomass",
				ColumnAlignment: leftAlignment,
			},
			OutHeaderColum{
				Text:            "roots",
				ColumnAlignment: leftAlignment,
			},
			OutHeaderColum{
				Text:            "LAImax",
				ColumnAlignment: leftAlignment,
			},
			OutHeaderColum{
				Text:            "Nfertil",
				ColumnAlignment: leftAlignment,
			},
			OutHeaderColum{
				Text:            "irrig",
				ColumnAlignment: leftAlignment,
			},
			OutHeaderColum{
				Text:            "N-uptake",
				ColumnAlignment: leftAlignment,
			},
			OutHeaderColum{
				Text:            "Nagb",
				ColumnAlignment: leftAlignment,
			},
			OutHeaderColum{
				Text:            "ETcG",
				ColumnAlignment: leftAlignment,
			},
			OutHeaderColum{
				Text:            "ETaG",
				ColumnAlignment: leftAlignment,
			},
			OutHeaderColum{
				Text:            "TraG",
				ColumnAlignment: leftAlignment,
			},
			OutHeaderColum{
				Text:            "PerG",
				ColumnAlignment: leftAlignment,
			},
			OutHeaderColum{
				Text:            "SWCS1",
				ColumnAlignment: leftAlignment,
			},
			OutHeaderColum{
				Text:            "SWCS2",
				ColumnAlignment: leftAlignment,
			},
			OutHeaderColum{
				Text:            "SWCA1",
				ColumnAlignment: leftAlignment,
			},
			OutHeaderColum{
				Text:            "SWCA2",
				ColumnAlignment: leftAlignment,
			},
			OutHeaderColum{
				Text:            "SWCM1",
				ColumnAlignment: leftAlignment,
			},
			OutHeaderColum{
				Text:            "SWCM2",
				ColumnAlignment: leftAlignment,
			},

			OutHeaderColum{
				Text:            "soilN1",
				ColumnAlignment: leftAlignment,
			},
			OutHeaderColum{
				Text:            "soilN2",
				ColumnAlignment: leftAlignment,
			},
			OutHeaderColum{
				Text:            "Nmin1",
				ColumnAlignment: leftAlignment,
			},
			OutHeaderColum{
				Text:            "Nmin2",
				ColumnAlignment: leftAlignment,
			},
			OutHeaderColum{
				Text:            "NLeaG",
				ColumnAlignment: leftAlignment,
			},
			OutHeaderColum{
				Text:            "TRRel",
				ColumnAlignment: leftAlignment,
			},
			OutHeaderColum{
				Text:            "Reduk",
				ColumnAlignment: leftAlignment,
			},
			OutHeaderColum{
				Text:            "DryD1",
				ColumnAlignment: leftAlignment,
			},
			OutHeaderColum{
				Text:            "DryD2",
				ColumnAlignment: leftAlignment,
			},
			OutHeaderColum{
				Text:            "Nresid",
				ColumnAlignment: leftAlignment,
			},
			OutHeaderColum{
				Text:            "Orgdat",
				ColumnAlignment: leftAlignment,
			},
			OutHeaderColum{
				Text:            "Type",
				ColumnAlignment: leftAlignment,
			},
			OutHeaderColum{
				Text:            "OrgN",
				ColumnAlignment: leftAlignment,
			},
			OutHeaderColum{
				Text:            "NDat1",
				ColumnAlignment: leftAlignment,
			},
			OutHeaderColum{
				Text:            "N1",
				ColumnAlignment: leftAlignment,
			},
			OutHeaderColum{
				Text:            "Ndat2",
				ColumnAlignment: leftAlignment,
			},
			OutHeaderColum{
				Text:            "N2",
				ColumnAlignment: leftAlignment,
			},
			OutHeaderColum{
				Text:            "Ndat3",
				ColumnAlignment: leftAlignment,
			},
			OutHeaderColum{
				Text:            "N3",
				ColumnAlignment: leftAlignment,
			},
			OutHeaderColum{
				Text:            "T_dat",
				ColumnAlignment: centerAlignment,
			},
			OutHeaderColum{
				Text:            "code",
				ColumnAlignment: centerAlignment,
			},
			OutHeaderColum{
				Text:            "Error",
				ColumnAlignment: leftAlignment,
			},
		},
		2: []OutHeaderColum{
			OutHeaderColum{
				Text:            "sowing",
				ColumnAlignment: leftAlignment,
				ColStart:        1,
			},
			OutHeaderColum{
				Text:            "emerg",
				ColumnAlignment: leftAlignment,
				ColStart:        3,
			},
			OutHeaderColum{
				Text:            "anth",
				ColumnAlignment: leftAlignment,
				ColStart:        4,
			},
			OutHeaderColum{
				Text:            "mat",
				ColumnAlignment: leftAlignment,
				ColStart:        5,
			},
			OutHeaderColum{
				Text:            "harvest",
				ColumnAlignment: leftAlignment,
				ColStart:        6,
				ColEnd:          7,
			},
			OutHeaderColum{
				Text:            "kg/ha",
				ColumnAlignment: leftAlignment,
				ColStart:        9,
			},
			OutHeaderColum{
				Text:            "kg/ha",
				ColumnAlignment: leftAlignment,
				ColStart:        10,
			},
			OutHeaderColum{
				Text:            "kg/ha",
				ColumnAlignment: leftAlignment,
				ColStart:        11,
			},
			OutHeaderColum{
				Text:            "kg N/ha",
				ColumnAlignment: leftAlignment,
				ColStart:        13,
			},
			OutHeaderColum{
				Text:            "mm",
				ColumnAlignment: leftAlignment,
				ColStart:        14,
			},
			OutHeaderColum{
				Text:            "kg N/ha",
				ColumnAlignment: leftAlignment,
				ColStart:        15,
			},
			OutHeaderColum{
				Text:            "kg N/ha",
				ColumnAlignment: leftAlignment,
				ColStart:        16,
			},
			OutHeaderColum{
				Text:            "mm",
				ColumnAlignment: leftAlignment,
				ColStart:        17,
			},
			OutHeaderColum{
				Text:            "mm",
				ColumnAlignment: leftAlignment,
				ColStart:        18,
			},
			OutHeaderColum{
				Text:            "mm",
				ColumnAlignment: leftAlignment,
				ColStart:        19,
			},
			OutHeaderColum{
				Text:            "mm",
				ColumnAlignment: leftAlignment,
				ColStart:        20,
			},
			OutHeaderColum{
				Text:            "mm",
				ColumnAlignment: leftAlignment,
				ColStart:        21,
			},
			OutHeaderColum{
				Text:            "mm",
				ColumnAlignment: leftAlignment,
				ColStart:        22,
			},
			OutHeaderColum{
				Text:            "mm",
				ColumnAlignment: leftAlignment,
				ColStart:        23,
			},
			OutHeaderColum{
				Text:            "mm",
				ColumnAlignment: leftAlignment,
				ColStart:        24,
			},
			OutHeaderColum{
				Text:            "mm",
				ColumnAlignment: leftAlignment,
				ColStart:        25,
			},
			OutHeaderColum{
				Text:            "mm",
				ColumnAlignment: leftAlignment,
				ColStart:        26,
			},
			OutHeaderColum{
				Text:            "kg N/ha",
				ColumnAlignment: leftAlignment,
				ColStart:        27,
			},
			OutHeaderColum{
				Text:            "kg N/ha",
				ColumnAlignment: leftAlignment,
				ColStart:        28,
			},
			OutHeaderColum{
				Text:            "kg N/ha",
				ColumnAlignment: leftAlignment,
				ColStart:        29,
			},
			OutHeaderColum{
				Text:            "kg N/ha",
				ColumnAlignment: leftAlignment,
				ColStart:        30,
			},
			OutHeaderColum{
				Text:            "kg N/ha",
				ColumnAlignment: leftAlignment,
				ColStart:        31,
			},
			OutHeaderColum{
				Text:            "kg N/ha",
				ColumnAlignment: leftAlignment,
				ColStart:        36,
			},
		},
	}
	for i := range headlines[1] {
		headlines[1][i].fillWithRune = ' '
		headlines[1][i].FillWithCharacter = " "
		headlines[1][i].ColStart = i + 1
		headlines[1][i].ColEnd = i + 1
	}

	for i := range headlines[2] {
		headlines[2][i].fillWithRune = ' '
		headlines[2][i].FillWithCharacter = " "
		if headlines[2][i].ColEnd == 0 {
			headlines[2][i].ColEnd = headlines[2][i].ColStart
		}
	}

	return OutputConfig{
		numHeadLines:       len(headlines),
		numDataColumns:     len(dataColumns),
		Headlines:          headlines,
		DataColumns:        dataColumns,
		FillCharacter:      " ",
		fillRune:           ' ',
		SeperatorCharacter: ",",
		seperatorRune:      ',',
		NotAvailableValue:  "n.a.",
	}
}

// NewDefaultDailyOutputConfig create daily output configuration
func NewDefaultDailyOutputConfig(g *GlobalVarsMain) OutputConfig {
	dataColumns := []OutputDataColum{
		OutputDataColum{
			FormatStr:     "%s",
			DataAlignment: leftAlignment,
			Width:         10,
			VarName:       "AKTUELL",
			valueRef:      &g.AKTUELL,
		},
		OutputDataColum{
			FormatStr:     "%04.1f",
			DataAlignment: leftAlignment,
			Width:         4,
			VarName:       "WG",
			Modifier:      100,
			VarIndex1:     1,
			VarIndex2:     0,
			valueRef:      &g.WG,
		},
		OutputDataColum{
			FormatStr:     "%04.1f",
			DataAlignment: leftAlignment,
			Width:         4,
			VarName:       "WG",
			Modifier:      100,
			VarIndex1:     1,
			VarIndex2:     1,
			valueRef:      &g.WG,
		},
		OutputDataColum{
			FormatStr:     "%04.1f",
			DataAlignment: leftAlignment,
			Width:         4,
			VarName:       "WG",
			Modifier:      100,
			VarIndex1:     1,
			VarIndex2:     2,
			valueRef:      &g.WG,
		},
		OutputDataColum{
			FormatStr:     "%04.1f",
			DataAlignment: leftAlignment,
			Width:         4,
			VarName:       "WG",
			Modifier:      100,
			VarIndex1:     1,
			VarIndex2:     3,
			valueRef:      &g.WG,
		},
		OutputDataColum{
			FormatStr:     "%04.1f",
			DataAlignment: leftAlignment,
			Width:         4,
			VarName:       "WG",
			Modifier:      100,
			VarIndex1:     1,
			VarIndex2:     4,
			valueRef:      &g.WG,
		},
		OutputDataColum{
			FormatStr:     "%04.1f",
			DataAlignment: leftAlignment,
			Width:         4,
			VarName:       "WG",
			Modifier:      100,
			VarIndex1:     1,
			VarIndex2:     5,
			valueRef:      &g.WG,
		},
		OutputDataColum{
			FormatStr:     "%04.1f",
			DataAlignment: leftAlignment,
			Width:         4,
			VarName:       "WG",
			Modifier:      100,
			VarIndex1:     1,
			VarIndex2:     6,
			valueRef:      &g.WG,
		},
		OutputDataColum{
			FormatStr:     "%04.1f",
			DataAlignment: leftAlignment,
			Width:         4,
			VarName:       "WG",
			Modifier:      100,
			VarIndex1:     1,
			VarIndex2:     7,
			valueRef:      &g.WG,
		},
		OutputDataColum{
			FormatStr:     "%04.1f",
			DataAlignment: leftAlignment,
			Width:         4,
			VarName:       "WG",
			Modifier:      100,
			VarIndex1:     1,
			VarIndex2:     8,
			valueRef:      &g.WG,
		},

		OutputDataColum{
			FormatStr:     "%04.f",
			DataAlignment: leftAlignment,
			Width:         4,
			VarName:       "C1",
			VarIndex1:     0,
			valueRef:      &g.C1,
		},
		OutputDataColum{
			FormatStr:     "%04.f",
			DataAlignment: leftAlignment,
			Width:         4,
			VarName:       "C1",
			VarIndex1:     1,
			valueRef:      &g.C1,
		},
		OutputDataColum{
			FormatStr:     "%04.f",
			DataAlignment: leftAlignment,
			Width:         4,
			VarName:       "C1",
			VarIndex1:     2,
			valueRef:      &g.C1,
		},
		OutputDataColum{
			FormatStr:     "%04.f",
			DataAlignment: leftAlignment,
			Width:         4,
			VarName:       "C1",
			VarIndex1:     3,
			valueRef:      &g.C1,
		},
		OutputDataColum{
			FormatStr:     "%04.f",
			DataAlignment: leftAlignment,
			Width:         4,
			VarName:       "C1",
			VarIndex1:     4,
			valueRef:      &g.C1,
		},
		OutputDataColum{
			FormatStr:     "%04.f",
			DataAlignment: leftAlignment,
			Width:         4,
			VarName:       "C1",
			VarIndex1:     5,
			valueRef:      &g.C1,
		},
		OutputDataColum{
			FormatStr:     "%04.f",
			DataAlignment: leftAlignment,
			Width:         4,
			VarName:       "C1",
			VarIndex1:     6,
			valueRef:      &g.C1,
		},
		OutputDataColum{
			FormatStr:     "%04.f",
			DataAlignment: leftAlignment,
			Width:         4,
			VarName:       "C1",
			VarIndex1:     7,
			valueRef:      &g.C1,
		},
		OutputDataColum{
			FormatStr:     "%04.f",
			DataAlignment: leftAlignment,
			Width:         4,
			VarName:       "C1",
			VarIndex1:     8,
			valueRef:      &g.C1,
		},
		OutputDataColum{
			FormatStr:     "%04.f",
			DataAlignment: leftAlignment,
			Width:         4,
			VarName:       "Nmin9to20",
			valueRef:      &g.Nmin9to20,
		},
		OutputDataColum{
			FormatStr:     "%05.f",
			DataAlignment: leftAlignment,
			Width:         5,
			VarName:       "PESUM",
			valueRef:      &g.PESUM,
		},
		OutputDataColum{
			FormatStr:     "%05.f",
			DataAlignment: leftAlignment,
			Width:         5,
			VarName:       "OUTSUM",
			valueRef:      &g.OUTSUM,
		},
		OutputDataColum{
			FormatStr:     "%05.f",
			DataAlignment: leftAlignment,
			Width:         5,
			VarName:       "SickerDaily",
			valueRef:      &g.SickerDaily,
		},
		OutputDataColum{
			FormatStr:     "%05.1f",
			DataAlignment: leftAlignment,
			Width:         5,
			VarName:       "CUMDENIT",
			valueRef:      &g.CUMDENIT,
		},
		OutputDataColum{
			FormatStr:     "%05.f",
			DataAlignment: leftAlignment,
			Width:         5,
			VarName:       "OBMAS",
			valueRef:      &g.OBMAS,
		},
		OutputDataColum{
			FormatStr:     "%04.f",
			DataAlignment: leftAlignment,
			Width:         5,
			VarName:       "MINSUM",
			valueRef:      &g.MINSUM,
		},
		OutputDataColum{
			FormatStr:     "%05.1f",
			DataAlignment: leftAlignment,
			Width:         6,
			VarName:       "HARVEST",
			valueRef:      &g.HARVEST,
		},
		OutputDataColum{
			FormatStr:     "%04.f",
			DataAlignment: leftAlignment,
			Width:         4,
			VarName:       "PHYLLO",
			valueRef:      &g.PHYLLO,
		},
		OutputDataColum{
			FormatStr:     "%05.f",
			DataAlignment: leftAlignment,
			Width:         5,
			VarName:       "WORG",
			VarIndex1:     0,
			valueRef:      &g.WORG,
		},
		OutputDataColum{
			FormatStr:     "%05.f",
			DataAlignment: leftAlignment,
			Width:         5,
			VarName:       "WORG",
			VarIndex1:     1,
			valueRef:      &g.WORG,
		},
		OutputDataColum{
			FormatStr:     "%05.f",
			DataAlignment: leftAlignment,
			Width:         5,
			VarName:       "WORG",
			VarIndex1:     2,
			valueRef:      &g.WORG,
		},
		OutputDataColum{
			FormatStr:     "%05.f",
			DataAlignment: leftAlignment,
			Width:         5,
			VarName:       "WORG",
			VarIndex1:     3,
			valueRef:      &g.WORG,
		},
		OutputDataColum{
			FormatStr:     "%05.f",
			DataAlignment: leftAlignment,
			Width:         5,
			VarName:       "WORG",
			VarIndex1:     4,
			valueRef:      &g.WORG,
		},
		OutputDataColum{
			FormatStr:     "%04.f",
			DataAlignment: leftAlignment,
			Width:         4,
			VarName:       "ASPOO",
			valueRef:      &g.ASPOO,
		},
		OutputDataColum{
			FormatStr:     "%05.2f",
			DataAlignment: leftAlignment,
			Width:         5,
			VarName:       "LAI",
			valueRef:      &g.LAI,
		},
		OutputDataColum{
			FormatStr:     "%5.3f",
			DataAlignment: leftAlignment,
			Width:         5,
			VarName:       "GEHMIN",
			valueRef:      &g.GEHMIN,
		},
		OutputDataColum{
			FormatStr:     "%5.3f",
			DataAlignment: leftAlignment,
			Width:         5,
			VarName:       "GEHMAX",
			valueRef:      &g.GEHMAX,
		},
		OutputDataColum{
			FormatStr:     "%4.f",
			DataAlignment: leftAlignment,
			Width:         4,
			VarName:       "NAOSAKT",
			valueRef:      &g.NAOSAKT,
		},
		OutputDataColum{
			FormatStr:     "%2d",
			DataAlignment: leftAlignment,
			Width:         2,
			VarName:       "WURZ",
			valueRef:      &g.WURZ,
		},
		OutputDataColum{
			FormatStr:     "%2.f",
			DataAlignment: leftAlignment,
			Width:         3,
			VarName:       "INTWICK.Num",
			valueRef:      &g.INTWICK.Num,
		},
		OutputDataColum{
			FormatStr:     "%06.2f",
			DataAlignment: leftAlignment,
			Width:         6,
			VarName:       "BLATTSUM",
			valueRef:      &g.BLATTSUM,
		},
		OutputDataColum{
			FormatStr:     "%05.1f",
			DataAlignment: leftAlignment,
			Width:         5,
			VarName:       "NFOSAKT",
			valueRef:      &g.NFOSAKT,
		},

		OutputDataColum{
			FormatStr:     "%06.1f",
			DataAlignment: leftAlignment,
			Width:         6,
			VarName:       "SumMINAOS",
			valueRef:      &g.SumMINAOS,
		},
		OutputDataColum{
			FormatStr:     "%05.1f",
			DataAlignment: leftAlignment,
			Width:         5,
			VarName:       "SumMINFOS",
			valueRef:      &g.SumMINFOS,
		},
		OutputDataColum{
			FormatStr:     "%5.f",
			DataAlignment: leftAlignment,
			Width:         5,
			Modifier:      10,
			VarName:       "VERDUNST",
			valueRef:      &g.VERDUNST,
		},
		OutputDataColum{
			FormatStr:     "%5.f",
			DataAlignment: leftAlignment,
			Width:         5,
			Modifier:      10,
			VarName:       "PFTRANS",
			valueRef:      &g.PFTRANS,
		},
		OutputDataColum{
			FormatStr:     "%04.2f",
			DataAlignment: leftAlignment,
			Width:         4,
			VarName:       "TRREL",
			valueRef:      &g.TRREL,
		},
		OutputDataColum{
			FormatStr:     "%04.2f",
			DataAlignment: leftAlignment,
			Width:         7,
			VarName:       "REDUK",
			valueRef:      &g.REDUK,
		},
		OutputDataColum{
			FormatStr:     "%s",
			DataAlignment: leftAlignment,
			Width:         7,
			VarName:       "NA",
			valueRef:      "n.a.",
		},
		OutputDataColum{
			FormatStr:     "%s",
			DataAlignment: leftAlignment,
			Width:         7,
			VarName:       "NA",
			valueRef:      "n.a.",
		},
		OutputDataColum{
			FormatStr:     "%05.f",
			DataAlignment: leftAlignment,
			Width:         5,
			VarName:       "ASPOO",
			valueRef:      &g.ASPOO,
		},
		OutputDataColum{
			FormatStr:     "%4.f",
			DataAlignment: leftAlignment,
			Width:         5,
			VarName:       "DRAISUM",
			valueRef:      &g.DRAISUM,
		},
		OutputDataColum{
			FormatStr:     "%6.f",
			DataAlignment: leftAlignment,
			Width:         6,
			VarName:       "DRAINLOSS",
			valueRef:      &g.DRAINLOSS,
		},
		OutputDataColum{
			FormatStr:     "%06.4f",
			DataAlignment: leftAlignment,
			Width:         6,
			VarName:       "GEHOB",
			valueRef:      &g.GEHOB,
		},
		OutputDataColum{
			FormatStr:     "%05.f",
			DataAlignment: leftAlignment,
			Width:         5,
			VarName:       "NFIXSUM",
			valueRef:      &g.NFIXSUM,
		},
		OutputDataColum{
			FormatStr:     "%05.f",
			DataAlignment: leftAlignment,
			Width:         5,
			VarName:       "NAOSAKT",
			valueRef:      &g.NAOSAKT,
		},
		OutputDataColum{
			FormatStr:     "%05.f",
			DataAlignment: leftAlignment,
			Width:         5,
			VarName:       "NFOSAKT",
			valueRef:      &g.NFOSAKT,
		},
		OutputDataColum{
			FormatStr:     "%+05.1f",
			DataAlignment: rightAlignment,
			Width:         7,
			VarIndex1:     5,
			VarName:       "AvgTSoil",
			valueRef:      &g.AvgTSoil,
		},
		OutputDataColum{
			FormatStr:     "%+05.1f",
			DataAlignment: leftAlignment,
			Width:         6,
			VarIndex1:     5,
			VarName:       "TD",
			valueRef:      &g.TD,
		},
		OutputDataColum{
			FormatStr:     "%s",
			DataAlignment: leftAlignment,
			Width:         10,
			VarName:       "C1NotStable",
			valueRef:      &g.C1NotStable,
		},
	}
	headlines := map[int][]OutHeaderColum{
		1: []OutHeaderColum{
			OutHeaderColum{
				Text:            "Date",
				ColumnAlignment: leftAlignment,
				ColStart:        1,
				ColEnd:          1,
			},
			OutHeaderColum{
				Text:            "water contents",
				ColumnAlignment: centerAlignment,
				ColStart:        2,
				ColEnd:          10,
			},
			OutHeaderColum{
				Text:            "-Nmin-content-",
				ColumnAlignment: centerAlignment,
				ColStart:        11,
				ColEnd:          20,
			},
			OutHeaderColum{
				Text:            "N crp",
				ColumnAlignment: centerAlignment,
			},
			OutHeaderColum{
				Text:            "Nleac",
				ColumnAlignment: centerAlignment,
			},
			OutHeaderColum{
				Text:            "perco",
				ColumnAlignment: centerAlignment,
			},
			OutHeaderColum{
				Text:            "denit",
				ColumnAlignment: centerAlignment,
			},
			OutHeaderColum{
				Text:            "agDM",
				ColumnAlignment: centerAlignment,
				ColStart:        25,
			},
			OutHeaderColum{
				Text:            "miner",
				ColumnAlignment: centerAlignment,
			},
			OutHeaderColum{
				Text:            "yield",
				ColumnAlignment: centerAlignment,
			},
			OutHeaderColum{
				Text:            "Phyl",
				ColumnAlignment: centerAlignment,
			},
			OutHeaderColum{
				Text:            "ORG1",
				ColumnAlignment: centerAlignment,
			},
			OutHeaderColum{
				Text:            "ORG2",
				ColumnAlignment: centerAlignment,
			},
			OutHeaderColum{
				Text:            "ORG3",
				ColumnAlignment: centerAlignment,
			},
			OutHeaderColum{
				Text:            "ORG4",
				ColumnAlignment: centerAlignment,
			},
			OutHeaderColum{
				Text:            "ORG5",
				ColumnAlignment: centerAlignment,
			},
			OutHeaderColum{
				Text:            "ASPO",
				ColumnAlignment: centerAlignment,
			},
			OutHeaderColum{
				Text:            "LAI",
				ColumnAlignment: leftAlignment,
			},
			OutHeaderColum{
				Text:            "ghmin",
				ColumnAlignment: centerAlignment,
			},
			OutHeaderColum{
				Text:            "ghmax",
				ColumnAlignment: centerAlignment,
			},
			OutHeaderColum{
				Text:            "NAOS",
				ColumnAlignment: centerAlignment,
			},
			OutHeaderColum{
				Text:            "rt",
				ColumnAlignment: centerAlignment,
			},
			OutHeaderColum{
				Text:            "stg",
				ColumnAlignment: centerAlignment,
			},
			OutHeaderColum{
				Text:            "leave",
				ColumnAlignment: centerAlignment,
			},
			OutHeaderColum{
				Text:            "NFOS",
				ColumnAlignment: centerAlignment,
			},
			OutHeaderColum{
				Text:            "MAOS",
				ColumnAlignment: centerAlignment,
			},
			OutHeaderColum{
				Text:            "MFOS",
				ColumnAlignment: centerAlignment,
			},
			OutHeaderColum{
				Text:            "ETP",
				ColumnAlignment: centerAlignment,
			},
			OutHeaderColum{
				Text:            "ETA",
				ColumnAlignment: centerAlignment,
			},
			OutHeaderColum{
				Text:            "TREL",
				ColumnAlignment: centerAlignment,
			},
			OutHeaderColum{
				Text:            "REDUK",
				ColumnAlignment: centerAlignment,
			},
			OutHeaderColum{
				Text:            "GPHOT",
				ColumnAlignment: centerAlignment,
			},
			OutHeaderColum{
				Text:            "MAINT",
				ColumnAlignment: centerAlignment,
			},
			OutHeaderColum{
				Text:            "ASPOO",
				ColumnAlignment: centerAlignment,
			},
			OutHeaderColum{
				Text:            "Wdrai",
				ColumnAlignment: centerAlignment,
			},
			OutHeaderColum{
				Text:            "NDrain",
				ColumnAlignment: centerAlignment,
			},
			OutHeaderColum{
				Text:            "Nagrb",
				ColumnAlignment: centerAlignment,
			},
			OutHeaderColum{
				Text:            "Nfix",
				ColumnAlignment: centerAlignment,
			},
			OutHeaderColum{
				Text:            "AOSAK",
				ColumnAlignment: centerAlignment,
			},
			OutHeaderColum{
				Text:            "FOSAK",
				ColumnAlignment: centerAlignment,
			},
			OutHeaderColum{
				Text:            "TS15",
				ColumnAlignment: centerAlignment,
			},
			OutHeaderColum{
				Text:            "TS50",
				ColumnAlignment: centerAlignment,
			},
			OutHeaderColum{
				Text:            "Error",
				ColumnAlignment: leftAlignment,
			},
		},
		2: []OutHeaderColum{
			OutHeaderColum{
				Text:              " Vol% ",
				ColumnAlignment:   centerAlignment,
				ColStart:          2,
				ColEnd:            10,
				FillWithCharacter: "-",
			},
			OutHeaderColum{
				Text:              " kg N/ha ",
				ColumnAlignment:   centerAlignment,
				ColStart:          11,
				ColEnd:            20,
				FillWithCharacter: "-",
			},
			OutHeaderColum{
				Text:              " kg/ha ",
				ColumnAlignment:   centerAlignment,
				ColStart:          21,
				ColEnd:            22,
				FillWithCharacter: "-",
			},
			OutHeaderColum{
				Text:              "mm",
				ColumnAlignment:   centerAlignment,
				ColStart:          23,
				ColEnd:            23,
				FillWithCharacter: " ",
			},
			OutHeaderColum{
				Text:              "----",
				ColumnAlignment:   leftAlignment,
				ColStart:          24,
				ColEnd:            24,
				FillWithCharacter: " ",
			},
			OutHeaderColum{
				Text:              "kg/ha",
				ColumnAlignment:   leftAlignment,
				ColStart:          25,
				ColEnd:            25,
				FillWithCharacter: " ",
			},
			OutHeaderColum{
				Text:              "--",
				ColumnAlignment:   leftAlignment,
				ColStart:          26,
				ColEnd:            26,
				FillWithCharacter: " ",
			},
		},
		3: []OutHeaderColum{
			OutHeaderColum{
				Text:            "0_1",
				ColumnAlignment: centerAlignment,
				ColStart:        2,
			},
			OutHeaderColum{
				Text:            "1_2",
				ColumnAlignment: centerAlignment,
				ColStart:        3,
			},
			OutHeaderColum{
				Text:            "2_3",
				ColumnAlignment: centerAlignment,
				ColStart:        4,
			},
			OutHeaderColum{
				Text:            "3_4",
				ColumnAlignment: centerAlignment,
				ColStart:        5,
			},
			OutHeaderColum{
				Text:            "4_5",
				ColumnAlignment: centerAlignment,
				ColStart:        6,
			},
			OutHeaderColum{
				Text:            "5_6",
				ColumnAlignment: centerAlignment,
				ColStart:        7,
			},
			OutHeaderColum{
				Text:            "6_7",
				ColumnAlignment: centerAlignment,
				ColStart:        8,
			},
			OutHeaderColum{
				Text:            "7_8",
				ColumnAlignment: centerAlignment,
				ColStart:        9,
			},
			OutHeaderColum{
				Text:            "8_9",
				ColumnAlignment: centerAlignment,
				ColStart:        10,
			},
			OutHeaderColum{
				Text:            "0_1",
				ColumnAlignment: centerAlignment,
				ColStart:        11,
			},
			OutHeaderColum{
				Text:            "1_2",
				ColumnAlignment: centerAlignment,
				ColStart:        12,
			},
			OutHeaderColum{
				Text:            "2_3",
				ColumnAlignment: centerAlignment,
				ColStart:        13,
			},
			OutHeaderColum{
				Text:            "3_4",
				ColumnAlignment: centerAlignment,
				ColStart:        14,
			},
			OutHeaderColum{
				Text:            "4_5",
				ColumnAlignment: centerAlignment,
				ColStart:        15,
			},
			OutHeaderColum{
				Text:            "5_6",
				ColumnAlignment: centerAlignment,
				ColStart:        16,
			},
			OutHeaderColum{
				Text:            "6_7",
				ColumnAlignment: centerAlignment,
				ColStart:        17,
			},
			OutHeaderColum{
				Text:            "7_8",
				ColumnAlignment: centerAlignment,
				ColStart:        18,
			},
			OutHeaderColum{
				Text:            "8_9",
				ColumnAlignment: centerAlignment,
				ColStart:        19,
			},
			OutHeaderColum{
				Text:            "9_20",
				ColumnAlignment: centerAlignment,
				ColStart:        20,
			},
		},
		4: []OutHeaderColum{
			OutHeaderColum{
				Text:              "_",
				ColumnAlignment:   leftAlignment,
				ColStart:          1,
				ColEnd:            15,
				FillWithCharacter: "_",
			},
		},
		// PrintTo(VNAMfile, "Date                         water contents                                -Nmin-content-                 N crp Nleac perco denit agDM  miner yield  Phyl ORG1  ORG2  ORG3  ORG4  ORG5  ASPO LAI   ghmin ghmax NAOS rt stg leave  NFOS  MAOS  MFOS   ETP   ETA TREL REDUK GPHOT  MAINT   ASPOO date_no  Wdrai NDrain Nagrb  Nfix  AOSAK FOSAK   TS15  TS50   Error\r\n")
		// PrintTo(VNAMfile, "           --------------------  Vol% ----------------- -------------------- kg N/ha -------------------- -- kg/ha --  mm   ---- kg/ha --\r\n")
		// line := fmt.Sprin("           0_1  1_2  2_3  3_4  4_5  5_6  6_7  7_8  8_9  0_1  1_2  2_3  3_4  4_5  5_6  6_7  7_8  8_9  9_20         %02d   %02d (in dm)  \r\n", g.OUTN, g.OUTN)
		// PrintTo(VNAMfile, line)
		// PrintTo(VNAMfile, "____________________________________________________________________________\r\n")

	}
	for i := range headlines[1] {
		headlines[1][i].fillWithRune = ' '
		headlines[1][i].FillWithCharacter = " "
		if headlines[1][i].ColStart == 0 && i > 0 {
			headlines[1][i].ColStart = headlines[1][i-1].ColEnd + 1
		}
		if headlines[1][i].ColEnd == 0 {
			headlines[1][i].ColEnd = headlines[1][i].ColStart
		}
	}

	for i := range headlines[2] {
		if headlines[2][i].ColEnd == 0 {
			headlines[2][i].ColEnd = headlines[2][i].ColStart
		}
	}
	for i := range headlines[3] {
		headlines[3][i].fillWithRune = ' '
		headlines[3][i].FillWithCharacter = " "
		if headlines[3][i].ColEnd == 0 {
			headlines[3][i].ColEnd = headlines[3][i].ColStart
		}
	}

	return OutputConfig{
		numHeadLines:       len(headlines),
		numDataColumns:     len(dataColumns),
		Headlines:          headlines,
		DataColumns:        dataColumns,
		FillCharacter:      " ",
		fillRune:           ' ',
		SeperatorCharacter: ",",
		seperatorRune:      ',',
		NotAvailableValue:  "n.a.",
	}
}

//LoadHermesOutputConfig loads a output file and reflects to programm variables
func LoadHermesOutputConfig(path string, g interface{}) (OutputConfig, error) {
	outConfig := OutputConfig{
		NotAvailableValue: "n.a.",
	}

	// if config files exists, read it into outConfig
	if _, err := os.Stat(path); err == nil {
		byteData := HermesFilePool.Get(&FileDescriptior{FilePath: path, ContinueOnError: true, UseFilePool: true})
		err := yaml.Unmarshal(byteData, &outConfig)
		if err != nil {
			return outConfig, err
		}
	}
	outConfig.numHeadLines = len(outConfig.Headlines)
	outConfig.numDataColumns = len(outConfig.DataColumns)
	outConfig.fillRune = getFirstRune(outConfig.FillCharacter, ' ')
	outConfig.seperatorRune = getFirstRune(outConfig.SeperatorCharacter, ',')
	for i := range outConfig.Headlines {
		lastColIndex := 0
		for idxCol := range outConfig.Headlines[i] {
			col := &outConfig.Headlines[i][idxCol]
			if col.ColStart == 0 {
				col.ColStart = lastColIndex + 1
			}
			if outConfig.Headlines[i][idxCol].ColStart < lastColIndex {
				return outConfig, fmt.Errorf("Out of order column index :%s", col.Text)
			}
			lastColIndex = col.ColStart
			if col.ColEnd == 0 {
				col.ColEnd = col.ColStart
			}
			col.fillWithRune = getFirstRune(col.FillWithCharacter, ' ')
		}
	}
	for i := range outConfig.DataColumns {
		dataCol := &outConfig.DataColumns[i]
		v := reflect.ValueOf(g)
		if v.Kind() == reflect.Ptr && !v.IsNil() {
			v = v.Elem()
		}
		varName := dataCol.VarName
		var subName string
		if strings.ContainsRune(varName, '.') {
			tokens := strings.SplitN(varName, ".", 2)
			varName = tokens[0]
			subName = tokens[1]
		}
		// handle reference to Array index
		f := v.FieldByName(varName)
		failedToBindAddr := true

		if f.IsValid() && f.Kind() == reflect.Struct {
			f = f.FieldByName(subName)
		}
		if f.IsValid() && f.Kind() == reflect.Array {
			if dataCol.VarIndex1 >= f.Len() {
				goto failed
			}
			f = f.Index(dataCol.VarIndex1)
			if f.IsValid() && f.Kind() == reflect.Array {
				if dataCol.VarIndex2 >= f.Len() {
					goto failed
				}
				f = f.Index(dataCol.VarIndex2)
			}
		}

		// all reflected Variables in the struct need to be public
		// that means: they have to start with a capital letter
		// if not it will panic here
		if f.IsValid() && f.CanAddr() {
			dataCol.valueRef = f.Addr().Interface()
			failedToBindAddr = false
		}

	failed:
		if failedToBindAddr {
			dataCol.valueRef = outConfig.NotAvailableValue
		}
	}

	return outConfig, nil
}

// WriteLine to outputfile
func (c *OutputConfig) WriteLine(file *Fout, formatType OutputFileFormat) error {

	outLine := NewOutputLine(c.numDataColumns)
	for _, col := range c.DataColumns {
		switch v := col.valueRef.(type) {
		case *string:
			outLine.Add(col.FormatStr, *v)
		case string:
			outLine.Add(col.FormatStr, v)
		case *int:
			outLine.Add(col.FormatStr, *v)
		case *float64:
			val := *v
			if col.Modifier != 0 {
				val = val * col.Modifier
			}
			outLine.Add(col.FormatStr, val)
		default:
			fmt.Println("unknown")
		}
	}
	var err error
	if formatType == csvOut {
		err = outLine.writeCSVString(file, c.seperatorRune)
	} else if formatType == hermesOut {
		err = outLine.writeHermesString(file, c)
	}

	return err
}

// Alignment for column texts
type Alignment int

const (
	leftAlignment Alignment = iota
	rightAlignment
	centerAlignment
	noneAlignment
)

var alignmentToString = map[Alignment]string{
	leftAlignment:   "left",
	rightAlignment:  "right",
	centerAlignment: "center",
	noneAlignment:   "none",
}

var alignmentToID = map[string]Alignment{
	"left":   leftAlignment,
	"right":  rightAlignment,
	"center": centerAlignment,
	"none":   noneAlignment,
}

// MarshalYAML implement YAML Marshaler
func (s Alignment) MarshalYAML() (interface{}, error) {
	return alignmentToString[s], nil
}

// UnmarshalYAML implement YAML Unmarshaler interface
func (s *Alignment) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var j string
	err := unmarshal(&j)
	if err != nil {
		return err
	}
	*s = alignmentToID[j]
	return nil
}

// OutputFileFormat to determine in which style the output is formated
type OutputFileFormat int

const (
	hermesOut OutputFileFormat = iota
	csvOut
)

// OutputLine for tupels of value and format into one line of output
type OutputLine struct {
	len          int
	format       []string
	counter      int
	addLinebreak bool
}

// NewOutputLine create a OutputLine with max number of elements
func NewOutputLine(num int) OutputLine {
	return OutputLine{len: num,
		format:       make([]string, num),
		counter:      0,
		addLinebreak: true}
}

// Add value to OutputLine
func (l *OutputLine) Add(format string, v interface{}) {
	if l.counter < l.len {
		l.format[l.counter] = fmt.Sprintf(format, v)
		l.counter++
	}
}

// AddDate string to OutputLine
func (l *OutputLine) AddDate(format string, date string, size int) {
	s := date
	if size > len(date) {
		sp := size - len(date)
		spl := sp/2 + len(date)
		spr := sp/2 + sp%2
		strFormat := "%" + strconv.Itoa(spl) + "s%" + strconv.Itoa(spr) + "s"
		s = fmt.Sprintf(strFormat, date, "")
	}
	if l.counter < l.len {
		l.format[l.counter] = fmt.Sprintf(format, s)
		l.counter++
	}
}

func (l *OutputLine) writeCSVString(file *Fout, seperatorRune rune) error {
	var err error
	for i, line := range l.format {
		_, err = file.Write(line)
		if err != nil {
			return err
		}
		if i < l.counter-1 {
			_, err = file.WriteRune(seperatorRune)
			if err != nil {
				return err
			}
		}
	}
	if l.addLinebreak {
		if l.counter > 0 && !strings.HasSuffix(l.format[l.counter-1], "\r\n") {
			_, err = file.Write("\r\n")
			if err != nil {
				return err
			}
		}
	}
	return err
}

func (l *OutputLine) writeHermesString(file *Fout, c *OutputConfig) error {
	var err error
	if c.numDataColumns != l.counter {
		return fmt.Errorf("Number of output columns: %d does not match counted values: %d", c.numDataColumns, l.counter)
	}
	for i, line := range l.format {
		column := c.DataColumns[i]
		columnWith := column.Width
		lenLine := utf8.RuneCountInString(line)
		if column.DataAlignment == rightAlignment {
			writefillChar := columnWith - lenLine
			for writefillChar > 0 {
				_, err = file.WriteRune(c.fillRune)
				if err != nil {
					return err
				}
				writefillChar--
			}
			_, err = file.Write(line)
			if err != nil {
				return err
			}
		} else if column.DataAlignment == leftAlignment {
			_, err = file.Write(line)
			if err != nil {
				return err
			}
			writefillChar := columnWith - lenLine
			for writefillChar > 0 {
				_, err = file.WriteRune(c.fillRune)
				if err != nil {
					return err
				}
				writefillChar--
			}
		} else if column.DataAlignment == centerAlignment || column.DataAlignment == noneAlignment {
			startWritefillChar := (columnWith - lenLine) / 2
			endWritefillChar := startWritefillChar + ((columnWith - lenLine) % 2)
			for startWritefillChar > 0 {
				_, err = file.WriteRune(c.fillRune)
				if err != nil {
					return err
				}
				startWritefillChar--
			}
			_, err = file.Write(line)
			if err != nil {
				return err
			}
			for endWritefillChar > 0 {
				_, err = file.WriteRune(c.fillRune)
				if err != nil {
					return err
				}
				endWritefillChar--
			}
		}
		if i < l.counter {
			_, err = file.WriteRune(c.fillRune)
			if err != nil {
				return err
			}
		}
	}
	if l.addLinebreak {
		if l.counter > 0 && !strings.HasSuffix(l.format[l.counter-1], "\r\n") {
			_, err = file.Write("\r\n")
			if err != nil {
				return err
			}
		}
	}
	return err
}

func (l *OutputLine) String() string {
	outStr := ""
	for _, line := range l.format {
		outStr = outStr + line
	}
	if l.addLinebreak {
		if l.counter > 0 && !strings.HasSuffix(outStr, "\r\n") {
			outStr = outStr + "\r\n"
		}
	}
	return outStr
}

func getFirstRune(str string, defaultVal rune) rune {
	first := defaultVal
	for _, c := range str {
		first = c
		break
	}
	return first
}
