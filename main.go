package main

import (
	"fmt"
	"os"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/types"
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
					dialog.ShowError(err, w)
					return
				}

				path := uc.URI().Path()

				if !isPDF(path) {
					alert("Please select the PDF file", w)
					return
				}

				err = addPasswordToPDF(path, "/Users/haoxilu/Downloads/test_encrypted.pdf", "testabcdefg")
				if err != nil {
					dialog.ShowError(err, w)
					return
				}

				text.SetText(path)
				alert("Done.", w)
			}, w).Show()
		}),
	))
	w.Resize(fyne.NewSize(800, 600))

	w.ShowAndRun()
}

func alert(message string, w fyne.Window) {
	dialog.ShowInformation("Alert", message, w)
}

func isPDF(filename string) bool {
	return strings.HasSuffix(strings.ToLower(filename), ".pdf")
}

func addPasswordToPDF(inputFile, outputFile, password string) error {
	// r the input PDF file.
	ctx, err := api.ReadContextFile(inputFile)
	if err != nil {
		return fmt.Errorf("error reading PDF file: %w", err)
	}

	// add watermark
	watermark_text := "野行Yeye"
	wm, err := api.TextWatermark(watermark_text, "font:Times-Italic, scale:.5, opacity:0.35", true, false, types.POINTS)
	if err != nil {
		return fmt.Errorf("parse text watermark: %w", err)
	}
	if err := pdfcpu.AddWatermarks(ctx, nil, wm); err != nil {
		return fmt.Errorf("add watermark to pdf: %w", err)
	}

	// create temporary file (with watermark)
	temp_file, err := os.CreateTemp(os.TempDir(), "secure_pdf_input_file")
	if err != nil {
		return fmt.Errorf("create temporary input file: %w", err)
	}
	defer os.Remove(temp_file.Name())

	// write output PDF file (with watermark)
	if err := api.WriteContextFile(ctx, temp_file.Name()); err != nil {
		return fmt.Errorf("write output PDF file: %w", err)
	}

	// add password
	encrypt_conf := model.NewAESConfiguration("upw", "opw", 256)
	encrypt_conf.Permissions = model.PermissionsNone
	if err := api.EncryptFile(temp_file.Name(), temp_file.Name(), encrypt_conf); err != nil {
		return fmt.Errorf("encrypt file: %w", err)
	}
	if err := api.SetPermissionsFile(temp_file.Name(), outputFile, encrypt_conf); err != nil {
		return fmt.Errorf("set permissions: %w", err)
	}

	return nil
}
