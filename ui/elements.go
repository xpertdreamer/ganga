package ui

import (
	"io"

	"ganga/util"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
)

func CreateFileDialog(window fyne.Window, appInstance fyne.App) *dialog.FileDialog {
	fd := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
		if err != nil {
			util.Error("%s", err)
			appInstance.Quit()
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
			appInstance.Quit()
			return
		}

		util.Debug("Content: %s", string(data))
	}, window)

	return fd
}
