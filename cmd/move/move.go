package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gordonklaus/portaudio"
	"github.com/maytanthegeek/move/pkg/keyboard"
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
	keysEvents, err := keyboard.GetKeys(10)
	chk(err)
	defer func() {
		_ = keyboard.Close()
	}()

UserAction:
	for {
		fmt.Print("move > ")

		event := <-keysEvents
		chk(event.Err)

		switch event.Rune {
		case 'p':
			if p.GetStatus() == player.Paused {
				go p.Play(stream)
			} else {
				p.Pause()
			}
		case 's':
			p.Stop()
		case 'c':
			fmt.Print("song > ")

			var filename string
			scanner.Scan()
			filename = scanner.Text()

			p.Stop()
			p.ChangeSong(filename)
		default:
			if event.Key == keyboard.KeyCtrlC {
				p.Stop()
				break UserAction
			}
		}

		fmt.Printf("%q\n", event.Rune)
	}

	fmt.Println("Bye")
	time.Sleep(1 * time.Second)
}

func main() {
	filename := getAudioFileArg()
	playAudioFile(filename)
}
