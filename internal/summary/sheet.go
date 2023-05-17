package summary

type Sheet struct {
	Month string `json:"month"`
	Lines []Line `json:"lines"`
	Total float64
}

type Line struct {
	Resource      string
	InvoiceNumber string
	Amount        float64
}
