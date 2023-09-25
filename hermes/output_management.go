package hermes

import (
	"fmt"
	"math"
	"os"
	"sort"

	yaml "gopkg.in/yaml.v2"
)

// open management output file

// write management event to the output

// close management output file

func (c *ManagementConfig) WriteManagementEvent(event *ManagementEvent) error {

	if c == nil {
		return nil
	}
	if c.file != nil {

		if c.EventFormats[event.eventName].Enabled {
			var err error
			// write event to file
			_, err = c.file.Write(event.hermesDate)
			if err != nil {
				return err
			}
			_, err = c.file.WriteRune(c.SeperatorRune)
			if err != nil {
				return err
			}
			_, err = c.file.Write(event.eventName.String())
			if err != nil {
				return err
			}
			_, err = c.file.WriteRune(c.SeperatorRune)
			if err != nil {
				return err
			}

			for _, attr := range c.EventFormats[event.eventName].sorted {
				formatStr := c.EventFormats[event.eventName].AdditionalFields[attr]
				if val, ok := event.additionalFields[attr]; ok {
					_, err = c.file.Write(fmt.Sprintf("%s: ", attr))
					if err != nil {
						return err
					}
					_, err = c.file.Write(fmt.Sprintf(formatStr, val))
					if err != nil {
						return err
					}
					_, err = c.file.WriteRune(c.SeperatorRune)
					if err != nil {
						return err
					}
				}
			}
			_, err = c.file.Write("\n")
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (c *ManagementConfig) Close() {
	if c.file != nil {
		c.file.Close()
	}
}

type ManagementConfig struct {
	// Management output configuration
	EventFormats  map[ManagementEventType]*ManagementEventConfig
	SeperatorRune rune
	file          *Fout
}

func (s *ManagementConfig) AnyOutputEnabled() bool {
	for _, v := range s.EventFormats {
		if v.Enabled {
			return true
		}
	}
	return false
}

func NewManagentConfig() *ManagementConfig {

	eventFormats := map[ManagementEventType]*ManagementEventConfig{
		Tillage: {
			EventName: Tillage,
			Enabled:   false,
			AdditionalFields: map[string]string{
				"Depth": "%dcm",
				"Type":  "%d",
			},
		},
		Irrigation: {
			EventName: Irrigation,
			Enabled:   false,
			AdditionalFields: map[string]string{
				"Amount": "%dmm",
				"N03":    "%2.1fmg/l",
			},
		},
		Fertilization: {
			EventName: Fertilization,
			Enabled:   false,
			AdditionalFields: map[string]string{
				"Fertilizer": "%2.1f",
			},
		},
		Sowing: {
			EventName: Sowing,
			Enabled:   false,
			AdditionalFields: map[string]string{
				"Crop": "%s",
			},
		},
		Harvest: {
			EventName: Harvest,
			Enabled:   false,
			AdditionalFields: map[string]string{
				"Crop":    "%s",
				"Residue": "%2.1f",
			},
		},
	}
	return &ManagementConfig{
		EventFormats:  eventFormats,
		SeperatorRune: ' ',
		file:          nil,
	}
}

func LoadManagementConfig(hp *HFilePath) (*ManagementConfig, error) {
	config := NewManagentConfig()
	// if config files exists, read it into hconfig
	if _, err := os.Stat(hp.managementOutput); err == nil {
		byteData := HermesFilePool.Get(&FileDescriptior{FilePath: hp.managementOutput, ContinueOnError: true, UseFilePool: true})
		err := yaml.Unmarshal(byteData, &config)
		if err != nil {
			return nil, err
		}
	} else {
		// no config exist, generate default config (if project is not fitting default setup, execution will fail)
		config = NewManagentConfig()
	}

	if anyOutPut := config.AnyOutputEnabled(); anyOutPut {
		// open management output file
		config.file = OpenResultFile(hp.mnam, false)
	}

	for _, eventConf := range config.EventFormats {

		sortedFields := make([]string, 0, len(eventConf.AdditionalFields))
		for k := range eventConf.AdditionalFields {
			sortedFields = append(sortedFields, k)
		}
		sort.Strings(sortedFields)
		eventConf.sorted = sortedFields
	}

	return config, nil
}

// new mamagement event handler

type ManagementEvent struct {
	eventName        ManagementEventType
	hermesDate       string
	additionalFields map[string]interface{}
}

func NewManagementEvent(eventType ManagementEventType, zeit int, additionalFields map[string]interface{}, g *GlobalVarsMain) *ManagementEvent {

	if eventType == Tillage {
		additionalFields["Depth"] = int(g.EINT[g.NTIL.Index])
		additionalFields["Type"] = g.TILART[g.NTIL.Index]
	} else if eventType == Irrigation {
		additionalFields["Amount"] = int(math.Round(g.EffectiveIRRIG))
		additionalFields["Fertilizer"] = g.BRKZn[g.NBR-1] * g.BREG[g.NBR-1] * 0.01 // fertilizer concentation in water
	} else if eventType == Sowing {
		additionalFields["Crop"] = g.CropTypeToString(g.FRUCHT[g.AKF.Index], false)
	} else if eventType == Harvest {
		additionalFields["Crop"] = g.CropTypeToString(g.FRUCHT[g.AKF.Index], true)
	}

	// create new management event
	return &ManagementEvent{
		eventName:        eventType,
		hermesDate:       g.Kalender(zeit),
		additionalFields: additionalFields,
	}
}

type ManagementEventConfig struct {
	EventName        ManagementEventType
	Enabled          bool
	AdditionalFields map[string]string
	sorted           []string
}

type ManagementEventType int

const (
	Tillage ManagementEventType = iota
	Irrigation
	Sowing
	Harvest
	Fertilization
)

func (s ManagementEventType) String() string {
	return meToString[s]
}

var meToString = map[ManagementEventType]string{
	Tillage:       "tillage",
	Irrigation:    "irrigation",
	Sowing:        "sowing",
	Harvest:       "harvest",
	Fertilization: "fertilization",
}

var meToID = map[string]ManagementEventType{
	"tillage":       Tillage,
	"irrigation":    Irrigation,
	"sowing":        Sowing,
	"harvest":       Harvest,
	"fertilization": Fertilization,
}

// MarshalYAML implement YAML Marshaler
func (s ManagementEventType) MarshalYAML() (interface{}, error) {
	return meToString[s], nil
}

// UnmarshalYAML implement YAML Unmarshaler interface
func (s *ManagementEventType) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var j string
	err := unmarshal(&j)
	if err != nil {
		return err
	}
	*s = meToID[j]
	return nil
}
