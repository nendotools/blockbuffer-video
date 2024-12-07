package types

type EncoderType string

const (
	Video EncoderType = "V"
	Audio EncoderType = "A"
)

type AVOptionEnum struct {
	Option      string
	ID          string
	Description string
}

type AVOption struct {
	Name        string
	Description string
	Type        string         // bool, int, enum
	Options     []AVOptionEnum // empty unless enum
}

type Encoder struct {
	Type        EncoderType
	Name        string
	Description string
	Formats     []string
	SampleRates []string // audio sample
	Options     []AVOption
}

var Encoders []Encoder

type AVOptionProfile struct {
	Name  string
	Value string
}

type AudioEncoderProfile struct {
	Name       string
	SampleRate string
	Options    []AVOptionProfile
}

type EncoderProfile struct {
	Name         string
	Format       string
	Options      []AVOptionProfile
	AudioEncoder AudioEncoderProfile
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
