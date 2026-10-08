package main

import (
	"ganga/back"
	"ganga/ui"
	"ganga/util"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	"os"
	"io"
	"strconv"
	"strings"
)

func main() {
	a := app.NewWithID("com.ganga")
	w := a.NewWindow("ganga")

	var argNum uint8 = 0
	if len(os.Args) > 1 {
		arg := os.Args[1]
		if arg == "cat" {
			argNum = 1
		}
	}
	mascotPaths, err := ui.LoadMaskotCache(argNum)
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

	exitBtn := widget.NewButton("Выход", func() {
		dialog.ShowConfirm("Подтверждение выхода", "Вы уверены, что хотите выйти?", func(confirmed bool) {
			if confirmed {
				a.Quit()
			}
		}, w)
	})

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
		"Добро пожаловать в ganga!",
		fyne.TextAlignCenter,
		fyne.TextStyle{Bold: true},
	)
	openBtn := widget.NewButton("Загрузить из...", func() {
		fd.Show()
	})

	maskot := canvas.NewImageFromResource(mascotPaths.Happy)
	maskot.FillMode = canvas.ImageFillContain
	maskot.SetMinSize(fyne.NewSize(250, 250))

	numQuestion := widget.NewLabel("1")
	question :=	widget.NewLabel("Question")

	restartBtn := widget.NewButton("На главную", func() {
		tabs.SelectIndex(0)
		tabs.DisableIndex(2)
		tabs.DisableIndex(3)
	})
	continueBtn := widget.NewButton("Продолжить", func() {
		addedQuest := engine.GetMoreQuestions(len(engine.Questions) / 3)
		if addedQuest > 0 {
			tabs.SelectIndex(2)
			tabs.DisableIndex(3)
			nextKey := engine.Questions[engine.StatePoiner]
			nextProp := engine.Base.Properties[nextKey]
			question.SetText(nextProp.Description + "?")
			numQuestion.SetText(strconv.FormatUint(engine.StatePoiner+1, 10))
		} else {
			dialog.ShowInformation("Внимание", "В базе знаний больше нет доступных вопросов для уточнения!", w)
		}
	})
	resultLabel := widget.NewLabelWithStyle(
		"",
		fyne.TextAlignCenter,
		fyne.TextStyle{Bold: true, Italic: true},
	)
	resultTabContent := container.NewCenter(
		container.NewVBox(
			widget.NewLabelWithStyle("Результат тестирования", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
			container.NewGridWithColumns(2, container.NewCenter(resultLabel), container.NewCenter(maskot)),
			container.NewCenter(container.NewHBox(restartBtn, continueBtn)),
		),
	)

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
			for i, w := range winners {
				builder.WriteString(engine.Base.Objects[w].Description)
				if i != len(winners) - 1 {
					builder.WriteString(", ")
				}
			}
			resultText = builder.String()
		}
		resultLabel.SetText(resultText)

		resultTabContent.Refresh()

		tabs.EnableIndex(3)
		tabs.SelectIndex(3)
	}

	yesButton := widget.NewButton("Да", func(){
		if ui.HandleButtonNext(engine, question, numQuestion, back.AnswerYes) {
			showResults()
		}
		ui.ChooseMaskotRand(maskot, mascotPaths)
	})
	noButton := widget.NewButton("Нет", func(){
		if ui.HandleButtonNext(engine, question, numQuestion, back.AnswerNo) {
			showResults()
		}
		ui.ChooseMaskotRand(maskot, mascotPaths)
	})
	idkButton := widget.NewButton("Не знаю", func(){
		if ui.HandleButtonNext(engine, question, numQuestion, back.AnswerUnknown) {
			showResults()
		}
		ui.ChooseMaskotRand(maskot, mascotPaths)
	})
	prevButton := widget.NewButton("Назад", func() {
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

	startBtn := widget.NewButton("Начать", func() {
		if len(engine.Base.Properties) == 0 {
			dialog.ShowInformation("Внимание!", "Сначала загрузите, или создайте базу знаний", w)
			return
		}
		numQuestion.SetText("1")
		engine.Start()
		tabs.EnableIndex(2)
		tabs.SelectIndex(2)
		question.SetText((engine.Base.Properties[engine.Questions[engine.StatePoiner]].Description) + "?")})

	mainTabContent := container.NewVBox(
		welcomeLabel,
		container.NewCenter(startBtn),
		container.NewCenter(exitBtn),
	)
	mainTabContentCentered := container.NewCenter(mainTabContent)

	scroll := container.NewScroll(content)
	scroll.SetMinSize(fyne.NewSize(600, 400))

	dbTabContent := container.NewVBox(
		widget.NewLabelWithStyle("Управление файлами знаний", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		container.NewCenter(openBtn),
		container.NewCenter(scroll),
		// TODO: db viewer, or base editor as sub-apllication
	)
	dbTabContentCentered := container.NewCenter(dbTabContent)

	tabs = container.NewAppTabs(
		container.NewTabItem("Главная", mainTabContentCentered),
		container.NewTabItem("Загрузка базы знаний", dbTabContentCentered),
		container.NewTabItem("Игра", gameTabContent),
		container.NewTabItem("Результаты", resultTabContent),
	)

	tabs.DisableIndex(2)
	tabs.DisableIndex(3)

	w.SetContent(tabs)
	w.Resize(fyne.NewSize(800, 600))
	w.Show()
	a.Run()
}
