package main

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"

	"invoiceGenerator/invoice"
	"invoiceGenerator/summary"
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

	app.Commands = []*cli.Command{
		{
			Name:    "summary-statement",
			Aliases: []string{"ss"},
			Usage:   "generate summary statement",
			Action:  summary.Action,
			Hidden:  true,
			Flags: []cli.Flag{
				&cli.StringSliceFlag{
					Name:    "invoice",
					Aliases: []string{"i"},
					Usage:   "path to the invoice file",
				},
			},
		},
	}

	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:     "config-file",
			Aliases:  []string{"c"},
			Usage:    "path to the configuration file",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "timesheet-path",
			Aliases:  []string{"t"},
			Usage:    "path to the timesheet file",
			Required: true,
		},
		},
		&cli.StringFlag{
			Name:    "output-file",
			Aliases: []string{"o"},
			Usage:   "path to the output file",
		},
	}

	app.Action = invoice.Action

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
