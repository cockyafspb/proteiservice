package absences

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"io"
	"net/http"
	"proteiservice/internal/domain/models"
	"strconv"
	"time"
)

type Absences struct {
	log      *zap.Logger
	ip       string
	port     int
	login    string
	password string
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

func (a *Absences) GetUser(ctx context.Context, email string) (fullName string, err error) {
	client := http.Client{Timeout: 5 * time.Second}
	type reqBody struct {
		Email string `json:"email"`
	}
	reqBodyByteArray, err := json.Marshal(reqBody{Email: email})
	if err != nil {
		a.log.Error("json.Marshal")
		return "", err
	}
	req, err := http.NewRequest("POST", "https://"+a.ip+":"+strconv.Itoa(a.port)+"/api", bytes.NewReader(reqBodyByteArray))
	if err != nil {
		return "", err
	}
	a.log.Info("Creating request")
	req.SetBasicAuth(a.login, a.password)
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	var body models.EmployeesRequestBody
	resBody, err := io.ReadAll(resp.Body)
	err = json.Unmarshal(resBody, &body)
	if err != nil {
		return "", err
	}
	var employee *models.Employee
	for _, v := range body.Data {
		employee = &v
	}
	if employee == nil {
		return "", fmt.Errorf("no data by this email")
	}
	a.log.Info(employee.Email)
	type personIds struct {
		PersonIds []int `json:"personIds"`
	}
	reqBodyByteArray, err = json.Marshal(personIds{PersonIds: []int{employee.Id}})
	if err != nil {
		return "", err
	}
	req, err = http.NewRequest("POST", "https://"+a.ip+":"+strconv.Itoa(a.port)+"/Portal/springApi/absences", bytes.NewReader(reqBodyByteArray))
	req.SetBasicAuth(a.login, a.password)
	resp, err = client.Do(req)
	defer resp.Body.Close()
	var bodyy models.AbsenceRequestBody
	resBody, err = io.ReadAll(resp.Body)
	err = json.Unmarshal(resBody, &bodyy)
	if err != nil {
		return "", err
	}
	var absence *models.Absence
	for _, v := range bodyy.Data {
		absence = &v
	}
	if absence == nil {
		return employee.Name, nil
	}
	// TODO: добавить получение здесь смайлика
	return employee.Name + strconv.Itoa(absence.ReasonId), nil
	return "", nil
}
