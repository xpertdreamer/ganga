package main

import (
	"ganga/back"
	"ganga/util"
	"ganga/ui"

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

	content := widget.NewLabel("")

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
		content.SetText(ui.CreatePreview(kb))
	}, w)

	welcomeLabel := widget.NewLabelWithStyle(
		"Добро пожаловать в Ganga!",
		fyne.TextAlignCenter,
		fyne.TextStyle{Bold: true},
	)
	openBtn := widget.NewButton("File Manager", func() {
		fd.Show()
	})
	var popUp *widget.PopUp
	popUpContent := container.NewVBox(
		widget.NewLabel("TODO: start button"),
		widget.NewButton("Close", func() {
			popUp.Hide()
		}),
	)
	popUp = widget.NewModalPopUp(popUpContent, w.Canvas())
	startBtn := widget.NewButton("Start", func() { util.Debug("TODO: start button"); popUp.Show()})


	mainTabContent := container.NewVBox(
		welcomeLabel,
		widget.NewSeparator(),
		container.NewCenter(startBtn),
	)

	scroll := container.NewScroll(content)
	scroll.SetMinSize(fyne.NewSize(600, 400))

	dbTabContent := container.NewVBox(
		widget.NewLabelWithStyle("Управление файлами знаний", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewSeparator(),
		container.NewCenter(openBtn),
		container.NewCenter(scroll),
		// TODO: db viewer
	)

	tabs := container.NewAppTabs(
		container.NewTabItem("Главная", mainTabContent),
		container.NewTabItem("Загрузка базы знаний", dbTabContent))

	w.SetContent(tabs)
	w.Resize(fyne.NewSize(800, 600))
	w.Show()
	a.Run()
}
