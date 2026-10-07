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
	engine := &back.Engine{}
	if err := engine.Create(kb); err != nil {
		util.Error("%s", err.Error())
		a.Quit()
		return
	}

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
	// var popUp *widget.PopUp
	// popUpContent := container.NewVBox(
	// 	widget.NewLabel("TODO: start button"),
	// 	widget.NewButton("Close", func() {
	// 		popUp.Hide()
	// 	}),
	// )
	// popUp = widget.NewModalPopUp(popUpContent, w.Canvas())

	question :=	widget.NewLabel("Question")
	prevButton := widget.NewButton("Previous", func() {util.Debug("TODO: previous")})
	nextButton := widget.NewButton("Next", func() {util.Debug("TODO: next")})

	gameTabContent := container.NewVBox(
		container.NewCenter(question),
		widget.NewSeparator(),
		container.NewCenter(container.NewHBox(prevButton, nextButton)),
	)

	var tabs* container.AppTabs

	startBtn := widget.NewButton("Start", func() { util.Debug("TODO: start button"); engine.Start(); tabs.EnableIndex(2); tabs.SelectIndex(2)})

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

	tabs = container.NewAppTabs(
		container.NewTabItem("Главная", mainTabContent),
		container.NewTabItem("Загрузка базы знаний", dbTabContent),
		container.NewTabItem("Игра", gameTabContent),
	)
	tabs.DisableIndex(2)

	w.SetContent(tabs)
	w.Resize(fyne.NewSize(800, 600))
	w.Show()
	a.Run()
}
