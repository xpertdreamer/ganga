package main

import (
	"ganga/back"
	"ganga/util"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	"io"
)

func main() {
	a := app.NewWithID("com.ganga")
	w := a.NewWindow("ganga")

	kb := &back.KBase{}
	fd := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
		if err != nil {
			util.Error("%s", err)
			a.Quit()
			return
		}
		if reader == nil {
			util.Error("No file selected or dialog closed")
			return
		}
		defer reader.Close()

		util.Debug("Selected: %s", reader.URI().Path())

		data, err := io.ReadAll(reader)
		if err != nil {
			util.Error("Error reading file: %s", err)
			a.Quit()
			return
		}

		util.Debug("Content: %s", string(data))
		if err := kb.Parse(data); err != nil {
			util.Error("cannot parse data")
			a.Quit()
			return
		}
	}, w)

	welcomeLabel := widget.NewLabelWithStyle(
		"Добро пожаловать в Ganga!",
		fyne.TextAlignCenter,
		fyne.TextStyle{Bold: true},
	)
	openBtn := widget.NewButton("File Manager", func() {
		fd.Show()

	})
	mainTabContent := container.NewVBox(
		welcomeLabel,
		widget.NewSeparator(),
	)
	dbTabContent := container.NewVBox(
		widget.NewLabelWithStyle("Управление файлами знаний", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewSeparator(),
		container.NewCenter(openBtn),
	)
	tabs := container.NewAppTabs(
		container.NewTabItem("Главная", mainTabContent),
		container.NewTabItem("Загрузка базы знаний", dbTabContent))
	w.SetContent(tabs)
	w.Resize(fyne.NewSize(800, 600))
	w.Show()
	a.Run()
}
