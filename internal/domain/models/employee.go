package models

type EmployeesRequestBody struct {
	Status string     `json:"status"`
	Data   []Employee `json:"data"`
}

type Employee struct {
	Name      string `json:"displayName"`
	Email     string `json:"email"`
	WorkPhone string `json:"workPhone"`
	Id        int    `json:"id"`
}
