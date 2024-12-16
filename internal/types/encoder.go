package types

type EncoderType string

const (
	Video EncoderType = "V"
	Audio EncoderType = "A"
)

type AVOptionEnum struct {
	Option      string `json:"option"`
	ID          string `json:"id"`
	Description string `json:"description"`
}

type AVOption struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Type        string         `json:"type"`    // bool, int, enum
	Options     []AVOptionEnum `json:"options"` // empty unless enum
}

type Encoder struct {
	Type        EncoderType `json:"type"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Formats     []string    `json:"formats"`
	SampleRates []string    `json:"sampleRates"` // empty unless audio
	Options     []AVOption  `json:"options"`
}

var Encoders []Encoder

type AVProfileOption struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type AudioEncoderProfile struct {
	Name       string            `json:"name"`
	SampleRate string            `json:"sampleRate"`
	Options    []AVProfileOption `json:"options"`
}

type EncoderProfile struct {
	Name         string              `json:"name"`
	Format       string              `json:"format"`
	Options      []AVProfileOption   `json:"options"`
	AudioEncoder AudioEncoderProfile `json:"audio"`
}

var DefaultEncoder EncoderProfile

func init() {
	DefaultEncoder = EncoderProfile{
		Name:   "dnxhd",
		Format: "yuv420p",
		Options: []AVProfileOption{
			{Name: "profile", Value: "dnxhr"},
		},
		AudioEncoder: AudioEncoderProfile{
			Name:       "pcm_s16le",
			SampleRate: "48000",
			Options:    []AVProfileOption{},
		},
	}
}
