package template

import (
	"bytes"
	_ "embed"
	"html/template"
	"log"
	"time"
)

//go:embed credentials.html
var accountInfoTemplate string

type AccountInfo struct {
	Name                   string
	UserMerchant           string
	PasswordMerchant       string
	UserSuperadmin         string
	PasswordSuperadmin     string
	AppMobile              string
	HomePage               string
	EmailAdmin             string
	Time                   time.Time
	PhoneCenter            string
	GuideBookAdminMerchant string
	GuideBookMerchant      string
	Website                string
}

func Credentials(data AccountInfo) (string, string) {
	subject := "Informasi Akun"
	tmp := template.Must(template.New("credentials").Parse(accountInfoTemplate))
	var body bytes.Buffer
	err := tmp.Execute(&body, data)
	if err != nil {
		log.Fatal(err)
		return "", ""
	}
	return subject, body.String()
}
