package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/dotcreep/mailapi/internal/service/mail"
	"github.com/dotcreep/mailapi/internal/service/template"
	"github.com/dotcreep/mailapi/internal/utils"
)

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
	AppMobile     string    `json:"app_mobile_name" form:"app_mobile_name"`
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

func SendEmailTemplate(w http.ResponseWriter, r *http.Request) {
	Json := utils.Json{}
	cfg, err := utils.OpenYAML()
	if err != nil {
		Json.NewResponse(false, w, nil, "internal server error", http.StatusInternalServerError, err.Error())
		return
	}
	if r.Method != http.MethodPost {
		Json.NewResponse(false, w, nil, "Method not allowed", http.StatusMethodNotAllowed, nil)
		return
	}
	tmpPath := cfg.Config.Storage
	mailAdmin := cfg.DataUser.EmailAdmin
	homePage := cfg.DataUser.Homepage

	var userData UserData
	var userCredential UserCredential
	var email Email
	if strings.Contains(r.Header.Get("Content-Type"), "application/json") {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			Json.NewResponse(false, w, nil, "Unable to read body", http.StatusBadRequest, err.Error())
			return
		}
		r.Body = io.NopCloser(bytes.NewBuffer(body))

		err = json.Unmarshal(body, &email)
		if err != nil {
			Json.NewResponse(false, w, nil, "Invalid email JSON format", http.StatusBadRequest, err.Error())
			return
		}

		if email.Attach != nil || len(email.Attach) > 0 {
			err = errors.New("attach not allowed")
			Json.NewResponse(false, w, nil, "attach not allowed", http.StatusBadRequest, err.Error())
			return
		}

		err = json.Unmarshal(body, &userData)
		if err != nil {
			Json.NewResponse(false, w, nil, "Invalid userData JSON format", http.StatusBadRequest, err.Error())
			return
		}

		err = json.Unmarshal(body, &userCredential)
		if err != nil {
			Json.NewResponse(false, w, nil, "Invalid userCredential JSON format", http.StatusBadRequest, err.Error())
			return
		}
	} else {
		err := r.ParseMultipartForm(10 << 20)
		if err != nil {
			Json.NewResponse(false, w, nil, "Unable to parse form", http.StatusBadRequest, err.Error())
			return
		}
		email.To = r.FormValue("to")
		email.Cc = r.FormValue("cc")
		email.Bcc = r.FormValue("bcc")
		email.Subject = r.FormValue("subject")
		email.TypeMessage = r.FormValue("type_message")
		email.Body = r.FormValue("body")
		userData.Reason = r.FormValue("Reason")
		files := r.MultipartForm.File["attach"]
		for _, fileHeader := range files {
			size := fileHeader.Size
			if size > 10*1024*1024 {
				err := errors.New("max size is 10MB")
				Json.NewResponse(false, w, nil, "File size is too large.", http.StatusBadRequest, err.Error())
				return
			}
			file, err := fileHeader.Open()
			if err != nil {
				Json.NewResponse(false, w, nil, "Unable open file", http.StatusBadRequest, err.Error())
				return
			}
			defer file.Close()
			if _, err := os.Stat(tmpPath); os.IsNotExist(err) {
				os.MkdirAll(tmpPath, 0755)
			}
			tempFilePath := filepath.Join(tmpPath, fileHeader.Filename)
			out, err := os.Create(tempFilePath)
			if err != nil {
				Json.NewResponse(false, w, nil, "Unable create file", http.StatusBadRequest, err.Error())
				return
			}
			defer out.Close()
			if _, err = io.Copy(out, file); err != nil {
				Json.NewResponse(false, w, nil, "Unable save file", http.StatusBadRequest, err.Error())
				return
			}
			email.Attach = append(email.Attach, tempFilePath)
		}
	}

	if email.To == "" {
		Json.NewResponse(false, w, nil, "To field is required", http.StatusBadRequest, nil)
		return
	}
	templateName := r.PathValue("name")
	if templateName == "" {
		Json.NewResponse(false, w, nil, "Template name is required", http.StatusBadRequest, nil)
		return
	}
	getYear := time.Now().Year()

	phoneCenter := cfg.DataUser.Phone
	if phoneCenter == "" {
		phoneCenter = "+620123456789"
	}

	guideBookAdminMerchant := cfg.DataUser.Guide.AdminMerchant
	if guideBookAdminMerchant == "" {
		guideBookAdminMerchant = homePage
	}

	guideBookMerchant := cfg.DataUser.Guide.Merchant
	if guideBookMerchant == "" {
		guideBookMerchant = homePage
	}
	switch templateName {
	case "welcome":
		if r.Header.Get("Content-Type") != "application/json" {
			userCredential.OTP = r.FormValue("otp")
		}
		if userCredential.OTP == "" {
			Json.NewResponse(false, w, nil, "OTP is required", http.StatusBadRequest, nil)
			return
		}
		getOTP := userCredential.OTP
		otp := fmt.Sprintf("%v %v %v %v", getOTP[0:1], getOTP[1:2], getOTP[2:3], getOTP[3:4])
		welcomeData := template.WelcomeData{
			Name:          userData.NameClient,
			EmailCustomer: email.To,
			EmailAdmin:    mailAdmin,
			HomePage:      homePage,
			OTP:           otp,
			Time:          time.Date(getYear, time.January, 1, 0, 0, 0, 0, time.Local),
			PhoneCenter:   phoneCenter,
		}
		email.Subject, email.Body = template.Welcome(welcomeData)
	case "registered":
		if r.Header.Get("Content-Type") != "application/json" {
			userData.NameClient = r.FormValue("client_name")
		}
		if userData.NameClient == "" {
			Json.NewResponse(false, w, nil, "Client name is required", http.StatusBadRequest, nil)
			return
		}
		regData := template.RegisteredData{
			Name:        userData.NameClient,
			EmailAdmin:  mailAdmin,
			HomePage:    homePage,
			Time:        time.Date(getYear, time.January, 1, 0, 0, 0, 0, time.Local),
			PhoneCenter: phoneCenter,
		}
		email.Subject, email.Body = template.Registered(regData)
	case "credentials":
		if r.Header.Get("Content-Type") != "application/json" {
			userCredential.UserMerchant = r.FormValue("user_merchant")
			userCredential.PasswordMerchant = r.FormValue("password_merchant")
			userCredential.UserSuperadmin = r.FormValue("user_superadmin")
			userCredential.PasswordSuperadmin = r.FormValue("password_superadmin")
			userData.NameClient = r.FormValue("client_name")
			userData.AppMobile = r.FormValue("app_mobile_url")
			userData.Website = r.FormValue("website")
		}
		if userCredential.UserMerchant == "" {
			Json.NewResponse(false, w, nil, "User Merchant is required", http.StatusBadRequest, nil)
			return
		}
		if userCredential.PasswordMerchant == "" {
			Json.NewResponse(false, w, nil, "Password Merchant is required", http.StatusBadRequest, nil)
			return
		}
		if userCredential.UserSuperadmin == "" {
			Json.NewResponse(false, w, nil, "User superadmin is required", http.StatusBadRequest, nil)
			return
		}
		if userCredential.PasswordSuperadmin == "" {
			Json.NewResponse(false, w, nil, "Password superadmin is required", http.StatusBadRequest, nil)
			return
		}
		if userData.NameClient == "" {
			Json.NewResponse(false, w, nil, "Client name is required", http.StatusBadRequest, nil)
			return
		}
		if userData.AppMobile == "" {
			Json.NewResponse(false, w, nil, "Mobile url is required", http.StatusBadRequest, nil)
			return
		}
		if userData.Website == "" {
			Json.NewResponse(false, w, nil, "Website url is required", http.StatusBadRequest, nil)
			return
		}
		accountData := template.AccountInfo{
			Name:                   userData.NameClient,
			UserMerchant:           userCredential.UserMerchant,
			PasswordMerchant:       userCredential.PasswordMerchant,
			UserSuperadmin:         userCredential.UserSuperadmin,
			PasswordSuperadmin:     userCredential.PasswordSuperadmin,
			AppMobile:              userData.AppMobile,
			HomePage:               homePage,
			EmailAdmin:             mailAdmin,
			Time:                   time.Date(getYear, time.January, 1, 0, 0, 0, 0, time.Local),
			PhoneCenter:            phoneCenter,
			GuideBookAdminMerchant: guideBookAdminMerchant,
			GuideBookMerchant:      guideBookMerchant,
			Website:                userData.Website,
		}
		email.Subject, email.Body = template.Credentials(accountData)
	case "rejected":
		if r.Header.Get("Content-Type") != "application/json" {
			userData.NameClient = r.FormValue("client_name")
			userData.Reason = r.FormValue("reason")
		}
		if userData.NameClient == "" {
			Json.NewResponse(false, w, nil, "Client name is required", http.StatusBadRequest, nil)
			return
		}
		if userData.Reason == "" {
			Json.NewResponse(false, w, nil, "Reason is required", http.StatusBadRequest, nil)
			return
		}
		rejectData := template.RejectTemplate{
			Name:        userData.NameClient,
			HomePage:    homePage,
			Reason:      strings.ToLower(userData.Reason),
			EmailAdmin:  mailAdmin,
			Time:        time.Date(getYear, time.January, 1, 0, 0, 0, 0, time.Local),
			PhoneCenter: phoneCenter,
		}
		email.Subject, email.Body = template.Rejected(rejectData)
	case "approved":
		if r.Header.Get("Content-Type") != "application/json" {
			userData.NameClient = r.FormValue("client_name")
		}
		if userData.NameClient == "" {
			Json.NewResponse(false, w, nil, "Client name is required", http.StatusBadRequest, nil)
			return
		}
		approvedData := template.ApprovedData{
			EmailAdmin:  mailAdmin,
			Name:        userData.NameClient,
			HomePage:    homePage,
			Time:        time.Date(getYear, time.January, 1, 0, 0, 0, 0, time.Local),
			PhoneCenter: phoneCenter,
		}
		email.Subject, email.Body = template.Approved(approvedData)
	case "invoice":
		if strings.Contains(r.Header.Get("Content-Type"), "application/json") {
			Json.NewResponse(false, w, nil, "Header is not allowed", http.StatusBadRequest, nil)
			return
		}
		if r.Header.Get("Content-Type") != "application/json" {
			userData.NameClient = r.FormValue("client_name")
			userData.Files = email.Attach
			userData.ConfirmURL = r.FormValue("url_upload")
		}
		if userData.NameClient == "" {
			Json.NewResponse(false, w, nil, "Client name is required", http.StatusBadRequest, nil)
			return
		}
		if len(userData.Files) == 0 {
			Json.NewResponse(false, w, nil, "Attachment is required", http.StatusBadRequest, nil)
			return
		}
		if userData.ConfirmURL == "" {
			Json.NewResponse(false, w, nil, "URL upload is required", http.StatusBadRequest, nil)
			return
		}
		invoiceData := template.InvoiceData{
			Name:        userData.NameClient,
			HomePage:    homePage,
			EmailAdmin:  mailAdmin,
			Time:        time.Date(getYear, time.January, 1, 0, 0, 0, 0, time.Local),
			Month:       time.Now().Month(),
			Year:        time.Now().Year(),
			PhoneCenter: phoneCenter,
			ConfirmURL:  userData.ConfirmURL,
		}
		email.Subject, email.Body = template.Invoice(invoiceData)
	case "invoice-paid":
		if strings.Contains(r.Header.Get("Content-Type"), "application/json") {
			Json.NewResponse(false, w, nil, "Header is not allowed", http.StatusBadRequest, nil)
			return
		}
		if r.Header.Get("Content-Type") != "application/json" {
			userData.NameClient = r.FormValue("client_name")
			userData.Files = email.Attach
			userData.ConfirmURL = r.FormValue("url_upload")
			strInvoiceNumber := r.FormValue("invoice_number")
			strTotal := r.FormValue("total")
			userData.InvoiceNumber, _ = strconv.Atoi(strInvoiceNumber)
			userData.Total, _ = strconv.Atoi(strTotal)
		}

		if userData.NameClient == "" {
			Json.NewResponse(false, w, nil, "Client name is required", http.StatusBadRequest, nil)
			return
		}
		if len(userData.Files) == 0 {
			Json.NewResponse(false, w, nil, "Attachment is required", http.StatusBadRequest, nil)
			return
		}
		if userData.ConfirmURL == "" {
			Json.NewResponse(false, w, nil, "URL upload is required", http.StatusBadRequest, nil)
			return
		}
		if userData.InvoiceNumber == 0 {
			Json.NewResponse(false, w, nil, "Invoice number is required", http.StatusBadRequest, nil)
			return
		}
		if userData.Total == 0 {
			Json.NewResponse(false, w, nil, "Total is required", http.StatusBadRequest, nil)
			return
		}
		invoicePaidData := template.InvoicePaidData{
			Name:          userData.NameClient,
			HomePage:      homePage,
			InvoiceNumber: userData.InvoiceNumber,
			Total:         userData.Total,
			DatePaid:      time.Now(),
			EmailAdmin:    mailAdmin,
			Time:          time.Date(getYear, time.January, 1, 0, 0, 0, 0, time.Local),
			Month:         time.Now().Month(),
			Year:          time.Now().Year(),
			PhoneCenter:   phoneCenter,
			ConfirmURL:    userData.ConfirmURL,
		}
		email.Subject, email.Body = template.InvoicePaid(invoicePaidData)
	default:
		Json.NewResponse(false, w, nil, "Template not found", http.StatusBadRequest, nil)
		return
	}
	if email.TypeMessage == "" {
		email.TypeMessage = "text/html"
	}
	info, err := mail.SendEmail(email.To, email.Cc, email.Bcc, email.Subject, email.TypeMessage, email.Body, email.Attach)
	if err != nil {
		Json.NewResponse(false, w, nil, "Failed to send email", http.StatusUnprocessableEntity, err)
		return
	}

	for _, filePath := range email.Attach {
		os.Remove(filePath)
	}

	Json.NewResponse(true, w, info, "Success", http.StatusOK, nil)
}
