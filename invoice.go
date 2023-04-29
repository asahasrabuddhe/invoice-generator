package main

import (
	"encoding/json"
	"io"
	"time"
)

type Invoice struct {
	Number     string    `json:"invoiceNumber"`
	Rate       float64   `json:"rate"`
	From       Contact   `json:"from"`
	To         Contact   `json:"to"`
	Lines      []Line    `json:"-"`
	ExtraLines []Line    `json:"extraLines"`
	Total      float64   `json:"-"`
	Date       string    `json:"invoiceDate"`
	Start      time.Time `json:"-"`
	End        time.Time `json:"-"`
}

type Contact struct {
	Email        string   `json:"email"`
	Name         string   `json:"name"`
	Phone        []string `json:"phone"`
	AddressLines []string `json:"addressLines"`
}

type Line struct {
	Description string  `json:"description"`
	Amount      float64 `json:"amount"`
}

func NewInvoice(r io.Reader) (*Invoice, error) {
	i := &Invoice{}
	i.Date = time.Now().Format("02-01-2006")

	err := json.NewDecoder(r).Decode(i)
	if err != nil {
		return nil, err
	}

	return i, nil
}
