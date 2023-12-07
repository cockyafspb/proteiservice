package models

type AbsenceRequestBody struct {
	Status string    `json:"status"`
	Data   []Absence `json:"data"`
}

type Absence struct {
	CreatedDate string `json:"createdDate"`
	DateFrom    string `json:"dateFrom"`
	DateTo      string `json:"dateTo"`
	Id          int    `json:"id"`
	PersonId    int    `json:"personId"`
	ReasonId    int    `json:"reasonId"`
}
