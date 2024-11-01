package template

import (
	"bytes"
	"html/template"
	"time"
)

const accountInfoTemplate = `<p>credentials</p>`

type AccountInfo struct {
	Name               string
	UserMerchant       string
	PasswordMerchant   string
	UserSuperadmin     string
	PasswordSuperadmin string
	AppMobile          string
	HomePage           string
	EmailAdmin         string
	Time               time.Time
	PhoneCenter        string
	GuideBook          string
}

func Credentials(data AccountInfo) (string, string) {
	subject := "Credentials"
	tmp := template.Must(template.New("credentials").Parse(accountInfoTemplate))
	var body bytes.Buffer
	err := tmp.Execute(&body, data)
	if err != nil {
		return "", ""
	}
	return subject, body.String()
}
