package hermes

import (
	"fmt"
	"log"
	"os"

	yaml "gopkg.in/yaml.v2"
)

// open management output file

// write management event to the output

// close management output file

func (c *ManagementConfig) WriteManagementEvent(event *ManagementEvent) error {

	if c.file != nil {

		if c.eventFormats[event.eventName].Enabled {
			var err error
			// write event to file
			_, err = c.file.Write(event.hermesDate)
			if err != nil {
				return err
			}
			_, err = c.file.WriteRune(c.seperatorRune)
			if err != nil {
				return err
			}
			_, err = c.file.Write(event.eventName.String())
			if err != nil {
				return err
			}

			for k, v := range event.additionalFields {

				if formatStr, ok := c.eventFormats[event.eventName].AdditionalFields[k]; ok {
					_, err = c.file.Write(fmt.Sprintf(formatStr, k, v))
				}
			}

			if err != nil {
				return err
			}
		}
	}
	return nil
}

type ManagementConfig struct {
	// Management output configuration
	eventFormats  map[ManagementEventType]ManagementEventConfig
	seperatorRune rune
	file          *Fout
}

func (s *ManagementConfig) AnyOutputEnabled() bool {
	for _, v := range s.eventFormats {
		if v.Enabled {
			return true
		}
	}
	return false
}

func NewManagentConfig() *ManagementConfig {

	eventFormats := map[ManagementEventType]ManagementEventConfig{
		Tillage: {
			EventName: Tillage,
			Enabled:   false,
			AdditionalFields: map[string]string{
				"Depth": "%d",
				"Type":  "%d",
			},
		},
		Irrigation: {
			EventName: Irrigation,
			Enabled:   false,
			AdditionalFields: map[string]string{
				"Amount":     "%d",
				"Fertilizer": "%s",
			},
		},
		Fertilization: {
			EventName: Fertilization,
			Enabled:   false,
			AdditionalFields: map[string]string{
				"Fertilizer": "%s",
				"Amount":     "%d",
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
		eventFormats:  eventFormats,
		seperatorRune: ' ',
		file:          nil,
	}
}

func ReadManagementConfig(hp *HFilePath) *ManagementConfig {
	config := NewManagentConfig()
	// if config files exists, read it into hconfig
	if _, err := os.Stat(hp.managementOutput); err == nil {
		byteData := HermesFilePool.Get(&FileDescriptior{FilePath: hp.managementOutput, ContinueOnError: true, UseFilePool: true})
		err := yaml.Unmarshal(byteData, &config)
		if err != nil {
			log.Fatalf("error: %v", err)
		}
	} else {
		// no config exist, generate default config (if project is not fitting default setup, execution will fail)
		config = NewManagentConfig()
	}

	if anyOutPut := config.AnyOutputEnabled(); !anyOutPut {
		// open management output file
		config.file = OpenResultFile(hp.mnam, false)
	}

	return config
}

// new mamagement event handler

type ManagementEvent struct {
	eventName        ManagementEventType
	hermesDate       string
	additionalFields map[string]interface{}
}

type ManagementEventConfig struct {
	EventName        ManagementEventType
	Enabled          bool
	AdditionalFields map[string]string
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
