package main

import (
	"ganga/back"
	"ganga/ui"
	"ganga/util"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	"io"
)

func handleButtonNext(e* back.Engine, quest *widget.Label, numQuest *widget.Label) {
	if e == nil {
		util.Error("cant deal with nil pointers")
		return
	}
	util.Debug("state pointer: %d", e.StatePoiner)
	util.Debug("questions: %v", e.Questions)
	util.Debug("props: %v", e.Base.Properties)
	if e.StatePoiner < uint64(len(e.Questions) - 1) {
		e.StatePoiner += 1
		quest.SetText((e.Base.Properties[e.Questions[e.StatePoiner]].Description) + "?")
		numQuest.SetText(strconv.FormatUint(e.StatePoiner + 1, 10))
	}
}

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

	var popUp *widget.PopUp
	popUpContent := container.NewVBox(
		widget.NewLabel("Upload base first"),
		widget.NewButton("Close", func() {
			popUp.Hide()
		}),
	)
	popUp = widget.NewModalPopUp(popUpContent, w.Canvas())

	numQuestion := widget.NewLabel("1")
	question :=	widget.NewLabel("Question")
	yesButton := widget.NewButton("Yes", func(){
		handleButtonNext(engine, question, numQuestion)
	})
	noButton := widget.NewButton("No", func(){
		handleButtonNext(engine, question, numQuestion)
	})
	idkButton := widget.NewButton("IDK", func(){
		handleButtonNext(engine, question, numQuestion)
	})
	prevButton := widget.NewButton("Previous", func() {util.Debug("TODO: previous")})

	gameTabContent := container.NewVBox(
		container.NewCenter(container.NewHBox(numQuestion, question)),
		widget.NewSeparator(),
		container.NewCenter(container.NewHBox(yesButton, idkButton, noButton)),
		container.NewCenter(container.NewHBox(prevButton)),
	)

	var tabs* container.AppTabs

	startBtn := widget.NewButton("Start", func() {
		util.Debug("TODO: start button")
		if len(engine.Base.Properties) == 0 {
			popUp.Show()
			return
		}
		engine.Start()
		tabs.EnableIndex(2)
		tabs.SelectIndex(2)
		question.SetText((engine.Base.Properties[engine.Questions[engine.StatePoiner]].Description) + "?")})

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
