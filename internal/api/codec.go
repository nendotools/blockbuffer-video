package api

import (
	"fmt"
	"net/http"
	"os/exec"
	"regexp"
	"strings"

	"blockbuffer/internal/io"
	"blockbuffer/internal/types"
	// opts "blockbuffer/internal/settings"
)

var Cmd *exec.Cmd

func HandleCodec(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		io.ErrorJSON(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// If Encoders is not nil, send the JSON response
	if types.Encoders != nil {
		io.SuccessJSON(w, types.Encoders)
		return
	}

	InitializeCodecs()
	io.SuccessJSON(w, types.Encoders)
}

func InitializeCodecs() {
	// Get video and audio codecs from ffmpeg
	Cmd = exec.Command("ffmpeg", "-encoders")
	encoders, err := Cmd.Output()
	if err != nil {
		io.Logf("Error getting encoders: %v", io.Error, err)
		return
	}

	// format codecs to JSON list
	/*
		* The output includes a section formatted like this:
		Decoders:
		V..... = Video
		A..... = Audio
		S..... = Subtitle
		.F.... = Frame-level multithreading
		..S... = Slice-level multithreading
		...X.. = Codec is experimental
		....B. = Supports draw_horiz_band
		.....D = Supports direct rendering method 1
		------
		*
		* This section is followed by a list of codecs, each formatted like this:
		V....D av1                  Alliance for Open Media AV1
		VFS..D dnxhd                VC3/DNxHD
		VFS..D h264                 H.264 / AVC / MPEG-4 AVC / MPEG-4 part 10
		VFS..D hevc                 HEVC (High Efficiency Video Coding)
		A....D aac                  AAC (Advanced Audio Coding)
		AF...D flac                 FLAC (Free Lossless Audio Codec)
		A....D wavesynth            Wave synthesis pseudo-codec
		*
		* We need to extract the codecs, group them by video or audio, and include their name and description from each line and group them by type.
	**/
	var ready = false
	for _, line := range strings.Split(string(encoders), "\n") {
		line = strings.TrimSpace(line)
		if !ready || line == "" {
			if strings.Contains(line, "------") {
				ready = true
			}
			continue
		}

		if line[0] == 'V' || line[0] == 'A' {
			name := strings.Split(line, " ")[1]
			desc := strings.Split(line, " ")[2:]
			codec, err := buildOptions(name, types.EncoderType(line[0]), strings.TrimSpace(strings.Join(desc, "")))
			if err != nil {
				continue
			}
			types.Encoders = append(types.Encoders, codec)
		}
	}
}

func buildOptions(encoderName string, encType types.EncoderType, desc string) (types.Encoder, error) {
	// Get video and audio codecs from ffmpeg
	Cmd = exec.Command("ffmpeg", "-help", "encoder="+encoderName)
	encoderOptions, err := Cmd.Output()
	if err != nil {
		return types.Encoder{}, err
	}

	var encoder = types.Encoder{
		Type:        encType,
		Name:        encoderName,
		Description: desc,
		Formats:     []string{},
		SampleRates: []string{},
		Options:     []types.AVOption{},
	}

	// Parse the options...
	// Audio and Video options are somewhat different, so we need to parse them differently
	if encType == types.Video {
		// Parse video options
		formats, options := processVideoOptions(string(encoderOptions))
		encoder.Formats = formats
		encoder.Options = options
		return encoder, nil
	}
	formats, sampleRates, options := processAudioOptions(string(encoderOptions))
	encoder.Formats = formats
	encoder.SampleRates = sampleRates
	encoder.Options = options
	return encoder, nil
}

func processVideoOptions(rawOPtions string) ([]string, []types.AVOption) {
	var formats []string
	var options []types.AVOption

	processingOptions := false
	var currentOption types.AVOption
	for _, line := range strings.Split(string(rawOPtions), "\n") {
		// skip empty lines
		if line == "" {
			continue
		}
		// Parse formats
		if strings.Contains(line, "Supported pixel formats:") {
			formats = strings.Split(strings.TrimSpace(strings.Split(line, ":")[1]), " ")
		}

		// Parse AVOptions
		if strings.Contains(line, " encoder AVOptions:") {
			processingOptions = true
			continue
		}
		if processingOptions {
			// Check if this is the start of a new option, row starts with a dash
			// If it is, commit the current option to the encoder and start a new one
			// If it is not, append the line to the current option Options enum field
			// If the line is empty, skip it
			//
			// options are formatted like this:
			// -{option}    {type}    E...A...... ({description})
			//    {options} {int}     E...A...... ({description}) <- if type is enum
			// -{option}    {boolean} E...A...... ({description}) <- bool types have no options and can be processed as flags
			// spaces could split each section by any number, so we need to ignore multiple spaces between each section
			// we can split the line by spaces and remove empty strings from the resulting slice
			var parts = strings.Split(line, " ")
			var cleanedParts []string
			for _, part := range parts {
				if part != "" {
					if len(cleanedParts) < 4 {
						// default or auto enum values do not have a type, so we need add it manually
						if len(cleanedParts) == 1 && (cleanedParts[0] == "default" || cleanedParts[0] == "auto") {
							cleanedParts = append(cleanedParts, "default")
						}
						cleanedParts = append(cleanedParts, part)
					} else {
						cleanedParts[3] = cleanedParts[3] + " " + part
					}
				}
			}
			cleanedParts = append(cleanedParts, "") // add an empty string to the end of the slice to avoid index error if description is missing

			// if starts with a dash, it is a new option
			if cleanedParts[0][0] == '-' {
				if currentOption.Name != "" {
					options = append(options, currentOption)
				}
				currentOption = types.AVOption{
					Name:        cleanedParts[0][1:],
					Type:        regexp.MustCompile("[^a-zA-Z]").ReplaceAllString(cleanedParts[1], ""),
					Description: cleanedParts[3],
					Options:     []types.AVOptionEnum{},
				}
			} else {
				id := cleanedParts[1]
				desc := cleanedParts[3]
				if strings.HasPrefix(cleanedParts[1], "E..") {
					id = fmt.Sprintf("%d", len(currentOption.Options))
					desc = cleanedParts[2] + " " + cleanedParts[3]
				}

				// if it is not a new option, it is an option enum
				currentOption.Options = append(currentOption.Options, types.AVOptionEnum{
					Option:      cleanedParts[0],
					ID:          id,
					Description: desc,
				})
			}
		}
	}

	// append the last option
	if currentOption.Name != "" {
		options = append(options, currentOption)
	}
	return formats, options
}

func processAudioOptions(rawOptions string) ([]string, []string, []types.AVOption) {
	var formats []string
	var sampleRates []string
	var options []types.AVOption

	processingOptions := false
	var currentOption types.AVOption
	for _, line := range strings.Split(string(rawOptions), "\n") {
		// skip empty lines
		if line == "" {
			continue
		}
		// Parse formats
		if strings.Contains(line, "Supported sample formats:") {
			formats = strings.Split(strings.TrimSpace(strings.Split(line, ":")[1]), " ")
		}

		// Parse sample rates
		if strings.Contains(line, "Supported sample rates:") {
			sampleRates = strings.Split(strings.TrimSpace(strings.Split(line, ":")[1]), " ")
		}

		// Parse AVOptions
		if strings.Contains(line, " encoder AVOptions:") {
			processingOptions = true
			continue
		}

		if processingOptions {
			// Check if this is the start of a new option, row starts with a dash
			// If it is, commit the current option to the encoder and start a new one
			// If it is not, append the line to the current option Options enum field
			// If the line is empty, skip it
			//
			// options are formatted like this:
			// -{option}    {type}    E...A...... ({description})
			//    {options} {int}     E...A...... ({description}) <- if type is enum
			// -{option}    {boolean} E...A...... ({description}) <- bool types have no options and can be processed as flags
			// spaces could split each section by any number, so we need to ignore multiple spaces between each section
			// we can split the line by spaces and remove empty strings from the resulting slice
			var parts = strings.Split(line, " ")
			var cleanedParts []string
			for _, part := range parts {
				if part != "" {
					if len(cleanedParts) < 4 {
						// default or auto enum values do not have a type, so we need add it manually
						if len(cleanedParts) == 1 && (cleanedParts[0] == "default" || cleanedParts[0] == "auto") {
							cleanedParts = append(cleanedParts, "default")
						}
						cleanedParts = append(cleanedParts, part)
					} else {
						cleanedParts[3] = cleanedParts[3] + " " + part
					}
				}
			}
			cleanedParts = append(cleanedParts, "") // add an empty string to the end of the slice to avoid index error if description is missing

			// if starts with a dash, it is a new option
			if cleanedParts[0][0] == '-' {
				if currentOption.Name != "" {
					options = append(options, currentOption)
				}
				currentOption = types.AVOption{
					// remove the dash from the option name
					Name:        cleanedParts[0][1:],
					Type:        regexp.MustCompile("[^a-zA-Z]").ReplaceAllString(cleanedParts[1], ""),
					Description: cleanedParts[3],
					Options:     []types.AVOptionEnum{},
				}
			} else {
				id := cleanedParts[1]
				desc := cleanedParts[3]
				if strings.HasPrefix(cleanedParts[1], "E..") {
					id = fmt.Sprintf("%d", len(currentOption.Options))
					desc = cleanedParts[2] + " " + cleanedParts[3]
				}

				// if it is not a new option, it is an option enum
				currentOption.Options = append(currentOption.Options, types.AVOptionEnum{
					Option:      cleanedParts[0],
					ID:          id,
					Description: desc,
				})
			}
		}
	}

	// append the last option
	if currentOption.Name != "" {
		options = append(options, currentOption)
	}
	return formats, sampleRates, options
}
