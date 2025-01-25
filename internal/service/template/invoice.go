package template

import (
	"bytes"
	_ "embed"
	"fmt"
	"html/template"
	"log"
	"time"
)

//go:embed invoice.html
var invoiceTemplate string

type InvoiceData struct {
	Name        string
	HomePage    string
	EmailAdmin  string
	Time        time.Time
	Month       time.Month
	Year        int
	PhoneCenter string
	ConfirmURL  string
}

func Invoice(data InvoiceData) (string, string) {
	subject := fmt.Sprintf("Invoice %s Periode %v %d", data.Name, data.Month, data.Year)
	tmp := template.Must(template.New("invoice").Parse(invoiceTemplate))
	var body bytes.Buffer
	err := tmp.Execute(&body, data)
	if err != nil {
		log.Println(err)
		return "", ""
	}
	return subject, body.String()
}
