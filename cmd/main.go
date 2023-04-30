package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/urfave/cli/v2"

	"invoiceGenerator"
	"invoiceGenerator/chrome"
	"invoiceGenerator/template"
	"invoiceGenerator/timesheet"
)

func main() {
	app := cli.NewApp()

	app.Authors = []*cli.Author{
		{
			Name:  "Ajitem Sahasrabuddhe",
			Email: "ajitem.s@outlook.com",
		},
	}

	app.Usage = "igen gnereates invoices from monday.com timesheets"

	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:    "config-file",
			Aliases: []string{"c"},
			Usage:   "path to the configuration file",
		},
		&cli.StringFlag{
			Name:    "timesheet-path",
			Aliases: []string{"t"},
			Usage:   "path to the timesheet file",
		},
		&cli.StringFlag{
			Name:    "output-file",
			Aliases: []string{"o"},
			Usage:   "path to the output file",
		},
	}

	app.Action = Action

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func Action(c *cli.Context) error {
	configFilePath := c.String("config-file")
	outFilePath := c.String("output-file")

	if outFilePath == "" {
		outFilePath = filepath.Dir(configFilePath)
	}

	configFile, err := os.Open(configFilePath)
	if err != nil {
		return err
	}

	invoice, err := invoiceGenerator.NewInvoice(configFile)
	if err != nil {
		return err
	}

	err = configFile.Close()
	if err != nil {
		return err
	}

	timesheetFile, err := os.Open(c.String("timesheet-path"))
	if err != nil {
		return err
	}

	err = timesheet.Parse(timesheetFile, invoice)
	if err != nil {
		return err
	}

	tpl, err := template.Get()
	if err != nil {
		return err
	}

	f, err := os.Create(filepath.Join(outFilePath, GetFileName(invoice)+".html"))
	if err != nil {
		return err
	}

	err = tpl.Execute(f, invoice)
	if err != nil {
		return err
	}

	err = f.Close()
	if err != nil {
		return err
	}

	path := chrome.Locate()

	err = exec.
		CommandContext(
			c.Context, path, "--headless", "--disable-gpu", "--no-pdf-header-footer",
			"--print-to-pdf="+filepath.Join(outFilePath, GetFileName(invoice)),
			filepath.Join(outFilePath, GetFileName(invoice)+".html"),
		).Run()
	if err != nil {
		return err
	}

	err = os.Remove(f.Name())
	if err != nil {
		return err
	}

	return nil
}

func GetFileName(invoice *invoiceGenerator.Invoice) string {
	extension := ".pdf"

	return fmt.Sprintf("%s - %s %d%s", invoice.Number, invoice.Start.Month().String(), invoice.Start.Year(), extension)
}
