package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func main() {
	a := app.New()
	w := a.NewWindow("Secure PDF")

	text := widget.NewLabel("Hello Yeye!")
	text.TextStyle.Bold = true
	w.SetContent(container.NewVBox(
		text,
		widget.NewButton("Choose your PDF file", func() {
			dialog.NewFileOpen(func(uc fyne.URIReadCloser, err error) {
				defer uc.Close()
				if err != nil {
					text.SetText("Failed to select PDF file")
					return
				}

				path := uc.URI().Path()

				text.SetText(path)
			}, w).Show()
		}),
	))
	w.Resize(fyne.NewSize(800, 600))

	w.ShowAndRun()
}
