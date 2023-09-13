package main

import (
	"image/color"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/maytanthegeek/move/internal/move"
)

func getAppContainer() *fyne.Container {
	// App Bar
	screenDescription := widget.NewLabel("PLAYING NOW")
	screenDescription.TextStyle.Bold = true

	appBarContainer := container.NewCenter(screenDescription)

	// Album Art
	albumArtBack := canvas.NewImageFromResource(move.ResourceAlbumArtPng)
	albumArtBack.FillMode = canvas.ImageFillContain
	albumArtBack.SetMinSize(fyne.NewSize(500, 500))

	// Song Info
	songName := canvas.NewText("Low Life", color.RGBA{0xA7, 0xA8, 0xAA, 0xFF})
	songName.Alignment = fyne.TextAlignCenter
	songName.TextStyle.Bold = true
	songName.TextSize = 36
	songArtist := container.NewCenter(widget.NewLabel("Future ft. The Weeknd"))
	songLengthSeek := move.NewSlider(0, 100)
	songLengthSeek.SetValue(50)

	// Action Buttons
	backwardButton := move.NewNeuButtonWithIcon("", theme.MediaSkipPreviousIcon(), func() {
		log.Println("tapped")
	})
	playPauseButton := move.NewNeuButtonWithIcon("", theme.MediaPlayIcon(), func() {
		log.Println("tapped")
	})
	forwardButton := move.NewNeuButtonWithIcon("", theme.MediaSkipNextIcon(), func() {
		log.Println("tapped")
	})

	// Main Containers
	actionButtonsContainer := container.NewHBox(
		layout.NewSpacer(),
		backwardButton,
		playPauseButton,
		forwardButton,
		layout.NewSpacer(),
	)

	mainContainer := container.NewVBox(
		layout.NewSpacer(),
		appBarContainer,
		layout.NewSpacer(),
		albumArtBack,
		layout.NewSpacer(),
		songName,
		songArtist,
		layout.NewSpacer(),
		songLengthSeek,
		layout.NewSpacer(),
		actionButtonsContainer,
		layout.NewSpacer(),
	)
	return mainContainer
}
