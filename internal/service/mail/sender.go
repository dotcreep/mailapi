package mail

import (
	"fmt"
	"os"
	"strconv"

	"github.com/dotcreep/mailapi/internal/utils"
	"gopkg.in/gomail.v2"
)

func SendEmail(to, cc, bcc, subject, typeMessage, body string, attach []string) (string, error) {
	cfg, err := utils.OpenYAML()
	if err != nil {
		return "Failed to open yaml file", err
	}
	host := cfg.Server.Host
	if host == "" {
		return "Failed to get host from .env file", err
	}
	portint := cfg.Server.Port
	port, err := strconv.Atoi(portint)
	if err != nil {
		return "Failed to get port from .env file", err
	}
	if port == 0 {
		return "Failed to get port from .env file", err
	}
	alias := cfg.Account.Alias
	if alias == "" {
		return "Failed to get alias from .env file", err
	}
	user := cfg.Account.Username
	if user == "" {
		return "Failed to get user from .env file", err
	}
	pass := cfg.Account.Password
	if pass == "" {
		return "Failed to get pass from .env file", err
	}

	senderName := cfg.Account.UsernameSender
	if senderName == "" {
		return "Failed to get 'from' value from .env file", err
	}
	d := gomail.NewDialer(host, port, user, pass)

	m := gomail.NewMessage()
	m.SetHeader("From", fmt.Sprintf("%s <%s>", senderName, alias))
	m.SetHeader("To", to)
	if cc != "" {
		m.SetHeader("Cc", cc)
	}
	if bcc != "" {
		m.SetHeader("Bcc", bcc)
	}
	m.SetHeader("Subject", subject)
	m.SetBody(typeMessage, body)

	if len(attach) != 0 {
		for _, file := range attach {
			if file == "qrcode.png" {
				m.Embed(file)
			}
			if _, err := os.Stat(file); os.IsNotExist(err) {
				return "Failed to attach file", fmt.Errorf("file %s does not exist", file)
			}
			m.Attach(file)
		}
	}

	if err := d.DialAndSend(m); err != nil {
		return "Failed to send email", err
	}
	return "Email sent successfully", nil
}
