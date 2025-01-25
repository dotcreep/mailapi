package template

import (
	"bytes"
	_ "embed"
	"html/template"
	"log"
	"time"
)

//go:embed welcome.html
var templateWelcome string

type WelcomeData struct {
	Name          string
	EmailCustomer string
	EmailAdmin    string
	HomePage      string
	OTP           string
	Time          time.Time
	PhoneCenter   string
}

func Welcome(data WelcomeData) (string, string) {
	subject := "Verifikasi Akun"
	tmp := template.Must(template.New("welcome").Parse(templateWelcome))
	var body bytes.Buffer
	err := tmp.Execute(&body, data)
	if err != nil {
		log.Fatal(err)
		return "", ""
	}
	return subject, body.String()
}
