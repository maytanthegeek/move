package main

import (
	"fyne.io/fyne/v2/app"

	"github.com/maytanthegeek/move/internal/move"
)

func main() {
	myApp := app.New()
	moveTheme := &move.Theme{}
	myApp.Settings().SetTheme(moveTheme)
	w := myApp.NewWindow("Move")

	mainContainer := getAppContainer()
	w.SetContent(mainContainer)

	w.ShowAndRun()
}
