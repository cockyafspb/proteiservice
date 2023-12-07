package models

type EmployeeGetter interface {
	GetEmployee(email string) (*Employee, error)
}

type AbsenceGetter interface {
	GetAbsence(employee *Employee) (int, error)
}

type Request struct {
	EmployeeGetter EmployeeGetter
	AbsenceGetter  AbsenceGetter
	Email          string
}

type ResultRequest struct {
	Name string
	Ok   bool
	Err  error
}
