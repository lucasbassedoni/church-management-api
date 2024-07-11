package models

type FinancialRecord struct {
	ID     int     `json:"id"`
	Type   string  `json:"type"`
	Name   string  `json:"name"`
	Amount float64 `json:"amount"`
	Date   string  `json:"date"`
}
