package player

import (
	"encoding/binary"
	"io"
	"time"

	"github.com/gordonklaus/portaudio"
)

func chk(err error) {
	if err != nil {
		panic(err)
	}
}

type Player struct {
	filename string
	song     *io.ReadCloser
	buffer   *[]int16
	status   SongState
}

func CreatePlayer(filename string, buffer *[]int16) (p *Player) {
	p = &Player{buffer: buffer, status: Stopped}
	p.ChangeSong(filename)
	return
}

func (p *Player) Play(stream *portaudio.Stream) {
	if p.status != Playing {
		p.status = Playing
		chk(stream.Start())

	SongPlaying:
		for err := binary.Read(*p.song, binary.LittleEndian, p.buffer); err == nil; err = binary.Read(*p.song, binary.LittleEndian, p.buffer) {
			stream.Write()

			switch p.status {
			case Paused:
				time.Sleep(1 * time.Second)
				chk(stream.Stop())
				break SongPlaying
			case Stopped:
				time.Sleep(1 * time.Second)
				chk(stream.Stop())
				p.ChangeSong(p.filename)
				break SongPlaying
			}
		}
	}
}

func (p *Player) Pause() {
	p.status = Paused
}

func (p *Player) Stop() {
	p.status = Stopped
}

func (p *Player) ChangeSong(filename string) {
	output := createFfmpegPipe(filename)
	p.filename = filename
	p.song = &output
}

func (p *Player) GetStatus() SongState {
	return p.status
}
