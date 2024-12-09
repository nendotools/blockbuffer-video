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

type AVOptionProfile struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type AudioEncoderProfile struct {
	Name       string            `json:"name"`
	SampleRate string            `json:"sampleRate"`
	Options    []AVOptionProfile `json:"options"`
}

type EncoderProfile struct {
	Name         string              `json:"name"`
	Format       string              `json:"format"`
	Options      []AVOptionProfile   `json:"options"`
	AudioEncoder AudioEncoderProfile `json:"audioEncoder"`
}

var DefaultEncoder EncoderProfile

func init() {
	var videoOptions = []AVOptionProfile{}
	videoOptions = append(videoOptions, AVOptionProfile{
		Name:  "profile",
		Value: "dnxhr_hq",
	})

	DefaultEncoder = EncoderProfile{
		Name:    "dnxhd",
		Format:  "yuv420p",
		Options: videoOptions,
		AudioEncoder: AudioEncoderProfile{
			Name:       "pcm_s16le",
			SampleRate: "48000",
			Options:    []AVOptionProfile{},
		},
	}
}
