package template

import (
	"bytes"
	_ "embed"
	"html/template"
	"log"
	"time"
)

//go:embed approved.html
var approvedTemplate string

type ApprovedData struct {
	EmailAdmin  string
	Name        string
	HomePage    string
	Time        time.Time
	PhoneCenter string
}

func Approved(data ApprovedData) (string, string) {
	subject := "Pendaftaran Telah Disetujui"
	tmp := template.Must(template.New("approved").Parse(approvedTemplate))
	var body bytes.Buffer
	err := tmp.Execute(&body, data)
	if err != nil {
		log.Println(err)
		return "", ""
	}
	return subject, body.String()
}
