package main

import (
	"ganga/ui"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func main() {
	a := app.NewWithID("com.ganga")
	w := a.NewWindow("ganga")
	fd := ui.CreateFileDialog(w, a)
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
	w.ShowAndRun()
}
