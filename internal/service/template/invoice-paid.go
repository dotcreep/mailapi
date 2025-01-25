package template

import (
	"bytes"
	_ "embed"
	"fmt"
	"html/template"
	"log"
	"strings"
	"time"
)

//go:embed invoice-paid.html
var invoicePaidTemplate string

type InvoicePaidData struct {
	Name          string
	HomePage      string
	InvoiceNumber int
	Total         int
	DatePaid      time.Time
	EmailAdmin    string
	Time          time.Time
	Month         time.Month
	Year          int
	PhoneCenter   string
	ConfirmURL    string
	FormatTotal   string
}

func formatCurrency(amount int) string {
	s := fmt.Sprintf("%d", amount)
	var result strings.Builder

	lamount := len(s)
	for i := 0; i < lamount; i++ {
		digit := rune(s[i])
		if (lamount-i)%3 == 0 && i != 0 {
			result.WriteRune('.')
		}
		result.WriteRune(digit)
	}
	return fmt.Sprintf("Rp.%s", result.String())
}

func InvoicePaid(data InvoicePaidData) (string, string) {
	subject := "Konfirmasi Pembayaran Invoice"
	data.FormatTotal = formatCurrency(data.Total)
	tmp := template.Must(template.New("invoice-paid").Parse(invoicePaidTemplate))
	var body bytes.Buffer
	err := tmp.Execute(&body, data)
	if err != nil {
		log.Println(err)
		return "", ""
	}
	return subject, body.String()
}
