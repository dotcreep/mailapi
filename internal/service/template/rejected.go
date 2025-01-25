package template

import (
	"bytes"
	_ "embed"
	"html/template"
	"log"
	"time"
)

//go:embed rejected.html
var rejectedTemplate string

type RejectTemplate struct {
	Name        string
	HomePage    string
	Reason      string
	EmailAdmin  string
	Time        time.Time
	PhoneCenter string
}

func Rejected(data RejectTemplate) (string, string) {
	subject := "Pendaftaran Ditolak"
	tmp := template.Must(template.New("rejected").Parse(rejectedTemplate))
	var body bytes.Buffer
	err := tmp.Execute(&body, data)
	if err != nil {
		log.Println(err)
		return "", ""
	}
	return subject, body.String()
}
