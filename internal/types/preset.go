package types

import (
	"embed"
	"encoding/json"
	"os"

	"blockbuffer/internal/io"
	opts "blockbuffer/internal/settings"
)

type Options []AVProfileOption

func (o *Options) UnmarshalJSON(data []byte) error {
	// Decode the JSON into a map
	var rawOptions map[string]string
	if err := json.Unmarshal(data, &rawOptions); err != nil {
		return err
	}

	// Transform the map into a slice of Option structs
	var options []AVProfileOption
	for key, value := range rawOptions {
		options = append(options, AVProfileOption{Name: key, Value: value})
	}

	// Assign to the receiver
	*o = options
	return nil
}

func (o Options) MarshalJSON() ([]byte, error) {
	// Transform the slice of Option structs into a map
	var rawOptions = make(map[string]string)
	if len(o) > 0 {
		for _, opt := range o {
			rawOptions[opt.Name] = opt.Value
		}
	}

	// Marshal the map into JSON
	return json.Marshal(rawOptions)
}

type AudioPreset struct {
	Codec      string   `json:"codec"`
	SampleRate *string  `json:"sampleRate"`
	Options    *Options `json:"options"`
}

type VideoPreset struct {
	Codec   string   `json:"codec"`
	Format  string   `json:"format"`
	Options *Options `json:"options"`
}

type PresetBundle struct {
	Name        string `json:"name"`
	Default     *bool
	Description string      `json:"description"`
	Extension   string      `json:"extension"`
	VideoPreset VideoPreset `json:"video"`
	AudioPreset AudioPreset `json:"audio"`
}

type PresetConfig struct {
	Presets []PresetBundle `json:"presets"`
}

//go:embed defaults.json
var embededPresets embed.FS
var Presets = make(map[string]PresetBundle)
var DefaultPreset PresetBundle

func init() {
	loadDefault()
	loadConfig(*opts.PresetConfigPath)
}

func loadDefault() {
	// load presets from file as JSON
	data, err := embededPresets.ReadFile("defaults.json")
	if err != nil {
		panic(err)
	}

	var presets PresetConfig
	presets = PresetConfig{}
	err = json.Unmarshal(data, &presets)
	if err != nil {
		panic(err)
	}

	for i, p := range presets.Presets {
		p.Default = &[]bool{true}[0]
		if i == 0 {
			DefaultPreset = p
		}
		Presets[p.Name] = p
	}
}

func loadConfig(path string) {
	// load presets from file as JSON
	data, err := os.ReadFile(path)
	if err != nil {
		// presets file doesn't exist, so create a blank one
		var blank = []byte(`{"presets":[]}`)
		os.WriteFile(path, blank, 0644)
		io.Logf("Config file not found, creating blank file at %s", io.Info, path)
	}

	var presets PresetConfig
	presets = PresetConfig{}
	err = json.Unmarshal(data, &presets)
	if err != nil {
		// log error and escape
		io.Logf("Error reading preset config: %v", io.Error, err)
	}
	for _, p := range presets.Presets {
		Presets[p.Name] = p
	}
}

func AddPreset(preset PresetBundle) {
	Presets[preset.Name] = preset
	ExportPresets()
}

func ExportPresets() {
	var presets PresetConfig
	presets = PresetConfig{}
	for _, preset := range Presets {
		// only store presets with no Default field or Default set to false
		if preset.Default == nil || !*preset.Default {
			presets.Presets = append(presets.Presets, preset)
		}
	}

	var data, err = json.MarshalIndent(presets, "", "  ")
	if err != nil {
		io.Logf("Error exporting presets: %v", io.Error, err)
	} else {
		os.WriteFile(*opts.PresetConfigPath, data, 0644)
	}
}
