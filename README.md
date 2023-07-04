# Invoice Generator

`igen` is a command line tool to generate invoices from monday.com time-sheets.

## Installation

Check for the latest version on the [releases page](https://github.com/asahasrabuddhe/invoice-generator/releases/latest) page and download the relevant binary for your OS and CPU architecture.

## Basic Usage

1. Create a working directory on your computer and move the downloaded binary from the previous step to this directory.
2. Create a configuration file, config.json, in the working directory. See [sample configuration file](#sample-configuration-file) for reference.
3. Download your timesheet from monday.com and save it in the working directory and rename it to timesheet.xlsx.
4. Run the following command in the terminal / command prompt / powershell to generate the invoice.
    ### Linux / MacOS
    ```bash
    ./igen -c config.json -t timesheet.xlsx
    ```
    
    ### Windows
    ```bash
    igen.exe -c config.json -t timesheet.xlsx
    ```
5. The generated invoice pdf will be saved in the working directory.

## Sample Configuration File

```json
{
  "invoiceNumber": "12345",
  "rate": 50,
  "currency": "US$",
  "from": {
    "email": "john@work.com",
    "name": "John Wick",
    "phone": [
      "+91 123 456 7890",
      "+91 123 456 7890"
    ],
    "addressLines": [
      "Address Line 1,",
      "Address Line 2,",
      "City - Zip Code"
    ]
  },
  "to": {
    "name": "Company Inc.",
    "addressLines": [
      "Address Line 1,",
      "Address Line 2,",
      "City - Zip Code"
    ]
  },
  "extraLines": [
    {
      "description": "expenses",
      "amount": 1234
    }
  ],
  "tax": {
    "name": "GST",
    "rate": 18,
    "accountNumber": "27AABCU9603R1ZP"
  }
}
```

## Reference
```
NAME:
   igen - igen gnereates invoices from monday.com timesheets

USAGE:
   igen [global options] command [command options] [arguments...]

AUTHOR:
   Ajitem Sahasrabuddhe <ajitem.s@outlook.com>

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --config-file value, -c value     path to the configuration file
   --timesheet-path value, -t value  path to the timesheet file
   --layout value, -l value          configures the layout of the generated invoice. possible values: [weekly, monthly] (default: "monthly")
   --output-file value, -o value     path to the output file
   --help, -h                        show help

```
