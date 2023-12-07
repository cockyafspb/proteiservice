package models

type AbsenceRequestBody struct {
	Status string    `json:"status"`
	Data   []Absence `json:"data"`
}

type Absence struct {
	Id          int    `json:"id"`
	PersonId    int    `json:"personId"`
	CreatedDate string `json:"createdDate"`
	DateFrom    string `json:"dateFrom"`
	DateTo      string `json:"dateTo"`
	ReasonId    int    `json:"reasonId"`
}
