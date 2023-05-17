package invoice

import (
	"encoding/json"
	"fmt"
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
	Tax        Tax       `json:"tax"`
	Total      float64   `json:"-"`
	TotalHours float64   `json:"-"`
	Date       string    `json:"invoiceDate"`
	Start      time.Time `json:"-"`
	End        time.Time `json:"-"`
	Layout     string    `json:"-"`
}

func (i Invoice) FileName() string {
	return fmt.Sprintf("%s - %s %d.pdf", i.Number, i.Start.Month().String(), i.Start.Year())
}

type Contact struct {
	Email        string   `json:"email"`
	Name         string   `json:"name"`
	Phone        []string `json:"phone"`
	AddressLines []string `json:"addressLines"`
}

type Line struct {
	StartDate   time.Time
	Description string  `json:"description"`
	Amount      float64 `json:"amount"`
}

type Tax struct {
	Type          string  `json:"type"`
	AccountNumber string  `json:"accountNumber"`
	Rate          float64 `json:"rate"`
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
