package player

import (
	"io"
	"os/exec"
)

func createFfmpegPipe(filename string) (output io.ReadCloser) {
	cmd := exec.Command("ffmpeg", "-i", filename, "-f", "s16le", "-")
	output, err := cmd.StdoutPipe()
	chk(err)

	err = cmd.Start()
	chk(err)

	return
}
