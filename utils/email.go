package utils

import (
	"bytes"
	"log"
	"os"
	"strconv"
	"text/template"
	"time"

	"golang.org/x/exp/rand"
	"gopkg.in/gomail.v2"
)

// Add this import

type Request struct {
	To      string `json:"to"`
	Subject string `json:"subject"`
	Text    string `json:"text"`
}

type Response struct {
	Message string `json:"message"`
}

type EmailData struct {
	Name  string
	Token string
	Url   string
}

func SendMail(to string, subject string, text string) error {

	m := gomail.NewMessage()
	m.SetHeader("From", os.Getenv("EMAIL_FROM"))
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", text)
	port, _ := strconv.Atoi(os.Getenv("EMAIL_PORT"))
	d := gomail.NewDialer(os.Getenv("EMAIL_HOST"), port, os.Getenv("EMAIL_FROM"), os.Getenv("EMAIL_PASSWORD"))
	if err := d.DialAndSend(m); err != nil {
		log.Fatal(err)
	}
	return nil
}

func RenderTemplate(templatePath string, data interface{}) (string, error) {
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return "", err
	}

	var rendered bytes.Buffer
	if err := tmpl.Execute(&rendered, data); err != nil {
		return "", err
	}

	return rendered.String(), nil
}

func GenerateRandomCode() string {
	code := ""
	for i := 0; i < 6; i++ {
		rand.Seed(uint64(time.Now().UnixNano()))
		code += strconv.Itoa(rand.Intn(10-1) + 1)
	}
	return code
}
