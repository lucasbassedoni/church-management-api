package models

type Event struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Date     string `json:"date"`
	Time     string `json:"time"`
	Location string `json:"location"`
	Month    string `json:"month"`
	Day      int    `json:"day"`
}
