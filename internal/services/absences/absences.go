package absences

import (
	"context"
	"go.uber.org/zap"
	"proteiservice/internal/domain/models"
	"strconv"
)

const NoAbsence = -1

type Absences struct {
	log            *zap.Logger
	employeeGetter EmployeeGetter
	absenceGetter  AbsenceGetter
	ip             string
	port           int
	login          string
	password       string
}

func New(log *zap.Logger, ip string, port int, login string, password string) *Absences {
	return &Absences{
		log:      log,
		ip:       ip,
		port:     port,
		login:    login,
		password: password,
	}
}

type EmployeeGetter interface {
	GetEmployee(email string) (*models.Employee, error)
}

type AbsenceGetter interface {
	GetAbsence(employee *models.Employee) (int, error)
}

func (a *Absences) GetUser(ctx context.Context, email string) (string, error) {
	employee, err := a.employeeGetter.GetEmployee(email)
	if err != nil {
		return "", err
	}
	id, err := a.absenceGetter.GetAbsence(employee)
	if err != nil {
		return "", err
	}
	if id == NoAbsence {
		return employee.Name, nil
	}
	return employee.Name + strconv.Itoa(id), nil
}
