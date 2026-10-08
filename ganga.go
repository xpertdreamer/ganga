package main

import (
	"ganga/back"
	"ganga/ui"
	"ganga/util"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	"io"
	"strings"
)

func main() {
	a := app.NewWithID("com.ganga")
	w := a.NewWindow("ganga")

	mascotPaths, err := ui.LoadMaskotCache()
	if err != nil {
		util.Error("initialisation error: %s", err.Error())
		a.Quit()
		return
	}

	content := widget.NewLabel("")
	var tabs* container.AppTabs

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

	maskot := canvas.NewImageFromResource(mascotPaths.Happy)
	maskot.FillMode = canvas.ImageFillContain
	maskot.SetMinSize(fyne.NewSize(250, 250))

	numQuestion := widget.NewLabel("1")
	question :=	widget.NewLabel("Question")

	gameScreenWrapper := container.NewStack()
	showResults := func() {
		winners, err := engine.GetResult()
		if err != nil {
			util.Error("error in result")
			a.Quit()
			return
		}
		var resultText string

		if len(winners) == 0 {
			resultText = "Не удалось точно определить объект"
		} else {
			var builder strings.Builder
			builder.WriteString("Я думаю, это: ")
			for _, w := range winners {
				builder.WriteString(engine.Base.Objects[w].Description)
			}
			resultText = builder.String()
			// resultText = "Я думаю, это: " + strings.Join(engine.Base.Objects[], ", ")
		}

		resultLabel := widget.NewLabelWithStyle(
			resultText,
			fyne.TextAlignCenter,
			fyne.TextStyle{Bold: true, Italic: true},
		)

		restartBtn := widget.NewButton("На главную", func() {
			tabs.SelectIndex(0)
		})
		// nextBtn := widget.NewButton("Продолжить", func)

		resultsLayout := container.NewCenter(
			container.NewVBox(
				widget.NewLabelWithStyle("Результат тестирования", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
				container.NewGridWithColumns(2, container.NewCenter(resultLabel), container.NewCenter(maskot)),
				container.NewCenter(restartBtn),
			),
		)

		gameScreenWrapper.Objects = []fyne.CanvasObject{resultsLayout}
		gameScreenWrapper.Refresh()
	}

	yesButton := widget.NewButton("Yes", func(){
		if ui.HandleButtonNext(engine, question, numQuestion, back.AnswerYes) {
			showResults()
		}
		ui.ChooseMaskotRand(maskot, mascotPaths)
	})
	noButton := widget.NewButton("No", func(){
		if ui.HandleButtonNext(engine, question, numQuestion, back.AnswerNo) {
			showResults()
		}
		ui.ChooseMaskotRand(maskot, mascotPaths)
	})
	idkButton := widget.NewButton("IDK", func(){
		if ui.HandleButtonNext(engine, question, numQuestion, back.AnswerUnknown) {
			showResults()
		}
		ui.ChooseMaskotRand(maskot, mascotPaths)
	})
	prevButton := widget.NewButton("Previous", func() {
		ui.HandleButtonPrev(engine, question, numQuestion)
		ui.ChooseMaskotRand(maskot, mascotPaths)
	})

	questionBox := container.NewCenter(container.NewHBox(numQuestion, question))
	middleRow := container.NewGridWithColumns(2,
		questionBox,
		container.NewCenter(maskot),
	)
	mainLayout := container.NewVBox(
		middleRow,
		container.NewCenter(container.NewHBox(yesButton, idkButton, noButton)),
		container.NewCenter(container.NewHBox(prevButton)),
	)
	gameTabContent := container.NewCenter(mainLayout)
	gameScreenWrapper.Add(gameTabContent)

	startBtn := widget.NewButton("Start", func() {
		util.Debug("TODO: start button")
		if len(engine.Base.Properties) == 0 {
			popUp.Show()
			return
		}
		numQuestion.SetText("1")
		engine.Start()
		gameScreenWrapper.Objects = []fyne.CanvasObject{gameTabContent}
		gameScreenWrapper.Refresh()
		tabs.EnableIndex(2)
		tabs.SelectIndex(2)
		question.SetText((engine.Base.Properties[engine.Questions[engine.StatePoiner]].Description) + "?")})

	mainTabContent := container.NewVBox(
		welcomeLabel,
		container.NewCenter(startBtn),
	)
	mainTabContentCentered := container.NewCenter(mainTabContent)

	scroll := container.NewScroll(content)
	scroll.SetMinSize(fyne.NewSize(600, 400))

	dbTabContent := container.NewVBox(
		widget.NewLabelWithStyle("Управление файлами знаний", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		container.NewCenter(openBtn),
		container.NewCenter(scroll),
		// TODO: db viewer
	)
	dbTabContentCentered := container.NewCenter(dbTabContent)

	tabs = container.NewAppTabs(
		container.NewTabItem("Главная", mainTabContentCentered),
		container.NewTabItem("Загрузка базы знаний", dbTabContentCentered),
		container.NewTabItem("Игра", gameScreenWrapper),
	)
	tabs.DisableIndex(2)

	w.SetContent(tabs)
	w.Resize(fyne.NewSize(800, 600))
	w.Show()
	a.Run()
}
