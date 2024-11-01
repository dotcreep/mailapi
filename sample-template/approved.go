package template

import (
	"bytes"
	"html/template"
	"time"
)

var approvedTemplate = `<p>Approved</p>`

type ApprovedData struct {
	EmailAdmin  string
	Name        string
	HomePage    string
	Time        time.Time
	PhoneCenter string
}

func Approved(data ApprovedData) (string, string) {
	subject := "Aprroved"
	tmp := template.Must(template.New("approved").Parse(approvedTemplate))
	var body bytes.Buffer
	err := tmp.Execute(&body, data)
	if err != nil {
		return "", ""
	}
	return subject, body.String()
}
