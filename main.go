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

				upwd, opwd := generatePasswords()
				err = handlPDF(path, "/Users/haoxilu/Downloads/test_encrypted.pdf", upwd, opwd)
				if err != nil {
					dialog.ShowError(err, w)
					return
				}

				pwd_text := generatePasswordsText(upwd, opwd)
				setTextToClipboard(pwd_text, w)
				alert(pwd_text, w)
			}, w).Show()
		}),
	))
	w.Resize(fyne.NewSize(800, 600))

	w.ShowAndRun()
}

func alert(message string, w fyne.Window) {
	dialog.ShowInformation("Alert", message, w)
}

func setTextToClipboard(text string, w fyne.Window) {
	w.Clipboard().SetContent(text)
}

func generatePasswords() (string, string) {
	return "upwd", "opwd"
}

func generatePasswordsText(upwd, opwd string) string {
	str_builder := strings.Builder{}
	str_builder.WriteString("You have copied password to your clipboard.\n")
	str_builder.WriteString("User Password: ")
	str_builder.WriteString(upwd)
	str_builder.WriteString("\n")
	str_builder.WriteString("Owner Password: ")
	str_builder.WriteString(opwd)
	return str_builder.String()
}

func isPDF(filename string) bool {
	return strings.HasSuffix(strings.ToLower(filename), ".pdf")
}

func handlPDF(inputFile, outputFile, userPassword, ownerPassword string) error {
	// create temporary file (with watermark)
	temp_file, err := os.CreateTemp(os.TempDir(), "secure_pdf_input_file")
	if err != nil {
		return fmt.Errorf("create temporary input file: %w", err)
	}
	temp_file_name := temp_file.Name()
	defer os.Remove(temp_file_name)

	// add watermark
	if err := addWatermarksToPDF(inputFile, temp_file_name); err != nil {
		return fmt.Errorf("add watermark: %w", err)
	}

	// add password
	if err := addPasswordToPDF(temp_file_name, outputFile, userPassword, ownerPassword); err != nil {
		return fmt.Errorf("add password: %w", err)
	}

	return nil
}

func addWatermarksToPDF(inputFile, outputFile string) error {
	// read the input PDF file.
	ctx, err := api.ReadContextFile(inputFile)
	if err != nil {
		return fmt.Errorf("error reading PDF file: %w", err)
	}

	// add text watermarks
	watermark_text := "野行Yeye"
	wm, err := api.TextWatermark(watermark_text, "font:Times-Italic, scale:.5, opacity:0.35", true, false, types.POINTS)
	if err != nil {
		return fmt.Errorf("parse text watermark: %w", err)
	}
	if err := pdfcpu.AddWatermarks(ctx, nil, wm); err != nil {
		return fmt.Errorf("add watermark to pdf: %w", err)
	}

	// write output PDF file (with watermark)
	if err := api.WriteContextFile(ctx, outputFile); err != nil {
		return fmt.Errorf("write output PDF file: %w", err)
	}

	return nil
}

func addPasswordToPDF(inputFile, outputFile, userPassword, ownerPassword string) error {
	encrypt_conf := model.NewAESConfiguration(userPassword, ownerPassword, 256)
	encrypt_conf.Permissions = model.PermissionsNone
	if err := api.EncryptFile(inputFile, inputFile, encrypt_conf); err != nil {
		return fmt.Errorf("encrypt file: %w", err)
	}
	if err := api.SetPermissionsFile(inputFile, outputFile, encrypt_conf); err != nil {
		return fmt.Errorf("set permissions: %w", err)
	}

	return nil
}
