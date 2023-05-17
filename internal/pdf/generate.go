package pdf

import (
	"context"
	"errors"
	"os"
	"os/exec"
	"path/filepath"

	"invoiceGenerator/internal/chrome"
	"invoiceGenerator/template"
)

func Generate(ctx context.Context, templateName, outFileName string, templateData any) error {
	// get current working directory
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	f, err := os.Create(filepath.Join(wd, outFileName+".html"))
	if err != nil {
		return err
	}

	tpl, err := template.Get(templateName)
	if err != nil {
		return err
	}

	err = tpl.Execute(f, templateData)
	if err != nil {
		return err
	}

	err = f.Close()
	if err != nil {
		return err
	}

	path := chrome.Locate()
	if path == "" {
		chrome.PromptDownload()
		return errors.New("chrome not found")
	}

	args := []string{
		"--headless", "--disable-gpu", "--no-pdf-header-footer", "--print-to-pdf=" + outFileName, outFileName + ".html",
	}
	err = exec.CommandContext(ctx, path, args...).Run()
	if err != nil {
		return err
	}

	err = os.Remove(f.Name())
	if err != nil {
		return err
	}

	return nil
}
