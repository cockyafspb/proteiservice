package transport

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"io"
	"net/http"
	"proteiservice/internal/domain/models"
	"time"
)

const (
	NoAbsence = -1
	NilValue  = 0
)

type HTTPManager struct {
	log      *zap.Logger
	client   *http.Client
	ip       string
	login    string
	password string
	port     int
}

func (m *HTTPManager) GetEmployee(email string) (*models.Employee, error) {
	type reqBody struct {
		Email string `json:"email"`
	}
	reqBodyByteArray, err := json.Marshal(reqBody{Email: email})
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(
		"POST",
		fmt.Sprintf("https://%s:%d/Portal/springApi/api/employees", m.ip, m.port),
		bytes.NewReader(reqBodyByteArray))
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(m.login, m.password)
	resp, err := m.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var body models.EmployeesRequestBody
	resBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(resBody, &body)
	if err != nil {
		return nil, err
	}
	var employee *models.Employee
	for i := range body.Data {
		employee = &body.Data[i]
		break
	}
	if employee == nil {
		return nil, errors.New("no employee with this email: " + email)
	}
	return employee, nil
}

// GetAbsence returns absence's id and error.
func (m *HTTPManager) GetAbsence(employee *models.Employee) (int, error) {
	type personIds struct {
		PersonIds []int `json:"personIds"`
	}
	reqBodyByteArray, err := json.Marshal(personIds{PersonIds: []int{employee.Id}})
	if err != nil {
		return NilValue, err
	}
	req, err := http.NewRequest(
		"POST",
		fmt.Sprintf("https://%s:%d/Portal/springApi/absences", m.ip, m.port),
		bytes.NewReader(reqBodyByteArray))
	if err != nil {
		return NilValue, err
	}
	req.SetBasicAuth(m.login, m.password)
	resp, err := m.client.Do(req)
	if err != nil {
		return NilValue, err
	}
	defer resp.Body.Close()
	var body models.AbsenceRequestBody
	resBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return NilValue, err
	}
	err = json.Unmarshal(resBody, &body)
	if err != nil {
		return NilValue, err
	}
	var absence *models.Absence
	for i := range body.Data {
		absence = &body.Data[i]
		break
	}
	if absence == nil {
		return NoAbsence, nil
	}
	return absence.ReasonId, nil
}

func New(log *zap.Logger, ip string, port int, login string, password string) *HTTPManager {
	return &HTTPManager{
		log:      log,
		client:   &http.Client{Timeout: 5 * time.Second},
		ip:       ip,
		port:     port,
		login:    login,
		password: password,
	}
}
