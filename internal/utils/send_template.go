package utils

import "time"

type EmailTemplate struct {
	To          string   `json:"to" form:"to"`
	Cc          string   `json:"cc" form:"cc"`
	Bcc         string   `json:"bcc" form:"bcc"`
	TypeMessage string   `json:"type" form:"type"`
	Subject     string   `json:"subject" form:"subject"`
	Body        string   `json:"body" form:"body"`
	Attach      []string `json:"attach" form:"attach"`
}

type UserData struct {
	NameClient    string    `json:"client_name" form:"client_name"`
	EmailClient   string    `json:"client_email" form:"client_email"`
	AppMobile     string    `json:"app_mobile_url" form:"app_mobile_url"`
	Files         []string  `json:"file" form:"file"`
	Reason        string    `json:"reason" form:"reason"`
	ConfirmURL    string    `json:"url_upload" form:"url_upload"`
	InvoiceNumber int       `json:"invoice_number" form:"invoice_number"`
	Total         int       `json:"total" form:"total"`
	DatePaid      time.Time `json:"date_paid" form:"date_paid"`
	Website       string    `json:"website" form:"website"`
}

type UserCredential struct {
	OTP                string `json:"otp" form:"otp"`
	UserMerchant       string `json:"user_merchant" form:"user_merchant"`
	PasswordMerchant   string `json:"password_merchant" form:"password_merchant"`
	UserSuperadmin     string `json:"user_superadmin" form:"user_superadmin"`
	PasswordSuperadmin string `json:"password_superadmin" form:"password_superadmin"`
}

type Email struct {
	To          string   `json:"to" form:"to"`
	Cc          string   `json:"cc" form:"cc"`
	Bcc         string   `json:"bcc" form:"bcc"`
	TypeMessage string   `json:"type" form:"type"`
	Subject     string   `json:"subject" form:"subject"`
	Body        string   `json:"body" form:"body"`
	Attach      []string `json:"attach" form:"attach"`
}
