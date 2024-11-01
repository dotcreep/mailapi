package template

import (
	"bytes"
	"html/template"
	"time"
)

const registeredTemplate = `<p>registered</p>`

type RegisteredData struct {
	Name        string
	EmailAdmin  string
	HomePage    string
	Time        time.Time
	PhoneCenter string
}

func Registered(data RegisteredData) (string, string) {
	subject := "Registered"
	tmp := template.Must(template.New("registered").Parse(registeredTemplate))
	var body bytes.Buffer
	err := tmp.Execute(&body, data)
	if err != nil {
		return "", ""
	}
	return subject, body.String()
}
