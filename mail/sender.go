package mail

import (
	"fmt"
	"mailapi/env"
	"os"
	"strconv"

	"gopkg.in/gomail.v2"
)

func SendEmail(to, cc, bcc string, subject string, typeMessage string, body string, attach []string) (string, error) {
	err := env.Load(".env")
	if err != nil {
		return "Failed to load .env file", err
	}

	host := os.Getenv("HOST")
	if host == "" {
		return "Failed to get host from .env file", err
	}
	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		return "Failed to convert port to int", err
	}
	if port == 0 {
		return "Failed to get port from .env file", err
	}
	user := os.Getenv("USER")
	if user == "" {
		return "Failed to get user from .env file", err
	}
	pass := os.Getenv("PASS")
	if pass == "" {
		return "Failed to get pass from .env file", err
	}

	senderName := os.Getenv("SENDER_NAME")
	if senderName == "" {
		return "Failed to get 'from' value from .env file", err
	}
	d := gomail.NewDialer(host, port, user, pass)

	m := gomail.NewMessage()
	m.SetHeader("From", fmt.Sprintf("%s <%s>", senderName, user))
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
