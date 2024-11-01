package template

import (
	"bytes"
	"html/template"
	"time"
)

const templateWelcome = `<p>welcome</p>`

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
	subject := "welcomes"
	tmp := template.Must(template.New("welcome").Parse(templateWelcome))
	var body bytes.Buffer
	err := tmp.Execute(&body, data)
	if err != nil {
		return "", ""
	}
	return subject, body.String()
}
