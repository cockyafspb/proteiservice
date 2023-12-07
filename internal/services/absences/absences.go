package absences

import (
	"context"
	"go.uber.org/zap"
	"proteiservice/internal/domain/models"
)

type Absences struct {
	log            *zap.Logger
	employeeGetter EmployeeGetter
	absenceGetter  AbsenceGetter
	emojis         map[int]string
	requestQueue   chan models.Request
	resultQueue    chan models.ResultRequest
}

func New(
	log *zap.Logger,
	employeeGetter EmployeeGetter,
	absenceGetter AbsenceGetter,
	emojis map[int]string,
	requestQueue chan models.Request,
	resultQueue chan models.ResultRequest) *Absences {
	return &Absences{
		log:            log,
		employeeGetter: employeeGetter,
		absenceGetter:  absenceGetter,
		emojis:         emojis,
		requestQueue:   requestQueue,
		resultQueue:    resultQueue,
	}
}

type EmployeeGetter interface {
	GetEmployee(email string) (*models.Employee, error)
}

type AbsenceGetter interface {
	GetAbsence(employee *models.Employee) (int, error)
}

func (a *Absences) GetUser(_ context.Context, email string) (string, bool, error) {
	a.requestQueue <- models.Request{
		EmployeeGetter: a.employeeGetter,
		AbsenceGetter:  a.absenceGetter,
		Email:          email,
	}
	res := <-a.resultQueue
	return res.Name, res.Ok, res.Err
}
