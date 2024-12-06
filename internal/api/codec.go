package api

import (
	"net/http"
	"os/exec"
	"strings"

	"blockbuffer/internal/io"
	// opts "blockbuffer/internal/settings"
)

var Cmd *exec.Cmd

func HandleCodec(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		io.ErrorJSON(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get video and audio codecs from ffmpeg
	Cmd = exec.Command("ffmpeg", "-encoders")
	encoders, err := Cmd.Output()
	if err != nil {
		io.ErrorJSON(w, "Failed to get codecs", http.StatusInternalServerError)
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
	codecs := []map[string]string{}
	for _, line := range strings.Split(string(encoders), "\n") {
		line = strings.TrimSpace(line)
		if !ready || line == "" {
			if strings.Contains(line, "------") {
				ready = true
			}
			continue
		}

		io.Logf("line: %s", io.Info, line)
		if line[0] == 'V' || line[0] == 'A' {
			name := strings.Split(line, " ")[1]
			desc := strings.Split(line, " ")[2:]
			codec := map[string]string{
				"type": string(line[0]),
				"name": name,
				"desc": strings.TrimSpace(strings.Join(desc, "")),
			}
			codecs = append(codecs, codec)
		}
	}

	io.SuccessJSON(w, codecs)
}
