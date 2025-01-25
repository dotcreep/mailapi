package template

import (
	"bytes"
	_ "embed"
	"html/template"
	"log"
	"time"
)

//go:embed registered.html
var registeredTemplate string

type RegisteredData struct {
	Name        string
	EmailAdmin  string
	HomePage    string
	Time        time.Time
	PhoneCenter string
}

func Registered(data RegisteredData) (string, string) {
	subject := "Pendaftaran Berhasil"
	tmp := template.Must(template.New("registered").Parse(registeredTemplate))
	var body bytes.Buffer
	err := tmp.Execute(&body, data)
	if err != nil {
		log.Fatal(err)
		return "", ""
	}
	return subject, body.String()
}
