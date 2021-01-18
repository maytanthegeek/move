package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gordonklaus/portaudio"
	"github.com/maytanthegeek/move/pkg/player"
)

func getAudioFileArg() (filename string) {
	if len(os.Args) < 2 {
		log.Fatal("Missing argument: input file name")
	}
	filename = os.Args[1]
	return
}

func chk(err error) {
	if err != nil {
		panic(err)
	}
}

func playAudioFile(filename string) {
	portaudio.Initialize()
	defer portaudio.Terminate()

	audiobuf := make([]int16, 8192)
	p := player.CreatePlayer(filename, &audiobuf)
	stream, err := portaudio.OpenDefaultStream(0, 2, 44100, len(audiobuf), &audiobuf)
	chk(err)
	defer stream.Close()

	go p.Play(stream)

	scanner := bufio.NewScanner(os.Stdin)
	var action string
UserAction:
	for {
		fmt.Print("move > ")
		scanner.Scan()
		action = scanner.Text()
		switch action {
		case "play":
			go p.Play(stream)
		case "pause":
			p.Pause()
		case "stop":
			p.Stop()
		case "change":
			var filename string
			fmt.Print("song > ")
			scanner.Scan()
			filename = scanner.Text()
			p.Stop()
			p.ChangeSong(filename)
		case "quit":
			p.Stop()
			break UserAction
		}
	}

	fmt.Println("Bye")
	time.Sleep(1 * time.Second)
}

func main() {
	filename := getAudioFileArg()
	playAudioFile(filename)
}
