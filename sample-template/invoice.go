package template

import (
	"bytes"
	"html/template"
	"time"
)

const invoiceTemplate = `<p>Invoice Template</p>`

type InvoiceData struct {
	Name        string
	HomePage    string
	EmailAdmin  string
	Time        time.Time
	PhoneCenter string
	GuideBook   string
}

func Invoice(data InvoiceData) (string, string) {
	subject := "Invoice " + data.Name
	tmp := template.Must(template.New("invoice").Parse(invoiceTemplate))
	var body bytes.Buffer
	err := tmp.Execute(&body, data)
	if err != nil {
		return "", ""
	}
	return subject, body.String()
}
