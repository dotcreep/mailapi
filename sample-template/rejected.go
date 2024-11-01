package template

import (
	"bytes"
	"html/template"
	"time"
)

const rejectedTemplate = `<p>rejected</p>`

type RejectTemplate struct {
	Name        string
	HomePage    string
	EmailAdmin  string
	Time        time.Time
	PhoneCenter string
}

func Rejected(data RejectTemplate) (string, string) {
	subject := "reject"
	tmp := template.Must(template.New("rejected").Parse(rejectedTemplate))
	var body bytes.Buffer
	err := tmp.Execute(&body, data)
	if err != nil {
		return "", ""
	}
	return subject, body.String()
}
