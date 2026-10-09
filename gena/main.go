package main

import (
	//	"ganga/util"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	// "fyne.io/fyne/v2/widget"
)

func main() {
	a := app.NewWithID("com.ganga.gena")
	w := a.NewWindow("gena")

	m := newModel()
	var currentPath string
	if len(os.Args) > 1 {
		if err := m.load(os.Args[1]); err == nil {
			currentPath = os.Args[1]
			w.SetTitle("gena - " + currentPath)
		}
	}



	w.Resize(fyne.NewSize(800, 600))
	w.Show()
	a.Run()
}
