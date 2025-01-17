package main

import (
	"io/ioutil"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
	"github.com/zalf-rpm/Hermes2Go/hermes"
)

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Config Editor")

	textArea := widget.NewMultiLineEntry()

	openButton := widget.NewButton("Open", func() {
		// get current directory
		currentDir, _ := os.Getwd()

		ShowConfigFileDialog(currentDir, func(reader fyne.URIReadCloser, err error) {
			if err == nil && reader != nil {
				data, _ := ioutil.ReadAll(reader)
				textArea.SetText(string(data))
				reader.Close()
			}
		}, myWindow)
	})

	saveButton := widget.NewButton("Save", func() {
		dialog.ShowFileSave(func(writer fyne.URIWriteCloser, err error) {
			if err == nil && writer != nil {
				data := []byte(textArea.Text)
				writer.Write(data)
				writer.Close()
			}
		}, myWindow)
	})

	myWindow.SetContent(container.NewVBox(
		openButton,
		saveButton,
		textArea,
	))

	myWindow.Resize(fyne.NewSize(800, 600))
	myWindow.ShowAndRun()
}

func ShowConfigFileDialog(startPath string, callback func(reader fyne.URIReadCloser, err error), parent fyne.Window) {

	locationURI, err := storage.ParseURI("file://" + startPath)
	if err != nil {
		dialog.ShowError(err, parent)
		return
	}
	listableURI, err := storage.ListerForURI(locationURI)
	if err != nil {
		dialog.ShowError(err, parent)
		return
	}
	dialog := dialog.NewFileOpen(callback, parent)

	dialog.SetLocation(listableURI)
	dialog.Show()
	dialog.SetOnClosed(func() {
		dialog = nil
	})

}

type ConfigEditorApp struct {
	loadedFileName string

	session              *hermes.HermesSession
	defaultGlobaVarlData *hermes.GlobalVarsMain
	loadedConfig         *hermes.OutputConfig
}

func NewConfigEditorApp() *ConfigEditorApp {

	return &ConfigEditorApp{
		session: hermes.NewHermesSession(),
	}
}

func (app *ConfigEditorApp) Close() {
	if app.session != nil {
		app.session.Close()
		app.session = nil
	}
}

func (app *ConfigEditorApp) LoadConfigFile(reader fyne.URIReadCloser) error {
	if reader == nil {
		return nil
	}
	defer reader.Close()
	// hmm how do I get the file name?
	return nil
}
