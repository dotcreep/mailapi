package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mailapi/mail"
	"mailapi/template"
	"mailapi/utils"
	"net/http"
	"os"
	"path/filepath"
	"time"
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
	NameClient  string `json:"client_name" form:"client_name"`
	EmailClient string `json:"client_email" form:"client_email"`
	AppMobile   string `json:"app_mobile_name" form:"app_mobile_name"`
}

type UserCredential struct {
	OTP                string `json:"otp" form:"otp"`
	UserMerchant       string `json:"user_merchant" form:"user_merchant"`
	PasswordMerchant   string `json:"password_merchant" form:"password_merchant"`
	UserSuperadmin     string `json:"user_superadmin" form:"user_superadmin"`
	PasswordSuperadmin string `json:"password_superadmin" form:"password_superadmin"`
}

func SendEmailTemplate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.RespondJSON(w, http.StatusMethodNotAllowed, map[string]interface{}{
			"success": false,
			"message": "Method not allowed",
		})
		return
	}
	tmpPath := os.Getenv("STORAGE")
	mailAdmin := os.Getenv("MAIL_ADMIN")
	homePage := os.Getenv("HOME_PAGE")

	var userData UserData
	var userCredential UserCredential
	var email Email
	if r.Header.Get("Content-Type") == "application/json" {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			utils.RespondJSON(w, http.StatusBadRequest, map[string]interface{}{
				"success": false,
				"message": "Unable to read body",
			})
			return
		}
		r.Body = io.NopCloser(bytes.NewBuffer(body))

		err = json.Unmarshal(body, &email)
		if err != nil {
			utils.RespondJSON(w, http.StatusBadRequest, map[string]interface{}{
				"success": false,
				"message": "Invalid email JSON format",
			})
			return
		}

		err = json.Unmarshal(body, &userData)
		if err != nil {
			utils.RespondJSON(w, http.StatusBadRequest, map[string]interface{}{
				"success": false,
				"message": "Invalid userData JSON format",
			})
			return
		}

		err = json.Unmarshal(body, &userCredential)
		if err != nil {
			utils.RespondJSON(w, http.StatusBadRequest, map[string]interface{}{
				"success": false,
				"message": "Invalid userCredential JSON format",
			})
			return
		}
	} else {
		err := r.ParseMultipartForm(10 << 20)
		if err != nil {
			utils.RespondJSON(w, http.StatusBadRequest, map[string]interface{}{
				"success": false,
				"message": "Unable to parse form",
			})
			return
		}
		email.To = r.FormValue("to")
		email.Cc = r.FormValue("cc")
		email.Bcc = r.FormValue("bcc")
		email.Subject = r.FormValue("subject")
		email.TypeMessage = r.FormValue("type_message")
		email.Body = r.FormValue("body")
		files := r.MultipartForm.File["attach"]
		for _, fileHeader := range files {
			size := fileHeader.Size
			if size > 10*1024*1024 {
				utils.RespondJSON(w, http.StatusBadRequest, map[string]interface{}{
					"success": false,
					"message": "File size is too large. Max size is 10MB",
				})
				return
			}
			file, err := fileHeader.Open()
			if err != nil {
				utils.RespondJSON(w, http.StatusBadRequest, map[string]interface{}{
					"success": false,
					"message": "Unable open file",
				})
				return
			}
			defer file.Close()
			tempFilePath := filepath.Join(tmpPath, fileHeader.Filename)
			out, err := os.Create(tempFilePath)
			if err != nil {
				utils.RespondJSON(w, http.StatusBadRequest, map[string]interface{}{
					"success": false,
					"message": "Unable create file",
				})
				return
			}
			defer out.Close()
			if _, err = io.Copy(out, file); err != nil {
				utils.RespondJSON(w, http.StatusBadRequest, map[string]interface{}{
					"success": false,
					"message": "Unable save file",
				})
				return
			}
			email.Attach = append(email.Attach, tempFilePath)
		}
	}

	if email.To == "" {
		utils.RespondJSON(w, http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": "To field is required",
		})
		return
	}
	templateName := r.PathValue("name")
	if templateName == "" {
		utils.RespondJSON(w, http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": "Template name is required",
		})
		return
	}
	getYear := time.Now().Year()

	phoneCenter := os.Getenv("PHONE_CENTER")
	if phoneCenter == "" {
		phoneCenter = "+620123456789"
	}

	guideBook := os.Getenv("GUIDE_BOOK")
	if guideBook == "" {
		guideBook = homePage
	}

	switch templateName {
	case "welcome":
		if r.Header.Get("Content-Type") != "application/json" {
			userCredential.OTP = r.FormValue("otp")
		}
		if userCredential.OTP == "" {
			utils.RespondJSON(w, http.StatusBadRequest, map[string]interface{}{
				"success": false,
				"message": "OTP is required",
			})
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
			utils.RespondJSON(w, http.StatusBadRequest, map[string]interface{}{
				"success": false,
				"message": "Client name is required",
			})
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
			userData.AppMobile = r.FormValue("app_mobile_name")
		}
		if userCredential.UserMerchant == "" {
			utils.RespondJSON(w, http.StatusBadRequest, map[string]interface{}{
				"success": false,
				"message": "User Merchant is required",
			})
			return
		}
		if userCredential.PasswordMerchant == "" {
			utils.RespondJSON(w, http.StatusBadRequest, map[string]interface{}{
				"success": false,
				"message": "Password merchant is required",
			})
			return
		}
		if userCredential.UserSuperadmin == "" {
			utils.RespondJSON(w, http.StatusBadRequest, map[string]interface{}{
				"success": false,
				"message": "User superadmin is required",
			})
			return
		}
		if userCredential.PasswordSuperadmin == "" {
			utils.RespondJSON(w, http.StatusBadRequest, map[string]interface{}{
				"success": false,
				"message": "Password superadmin is required",
			})
			return
		}
		if userData.NameClient == "" {
			utils.RespondJSON(w, http.StatusBadRequest, map[string]interface{}{
				"success": false,
				"message": "Client name is required",
			})
			return
		}
		if userData.AppMobile == "" {
			utils.RespondJSON(w, http.StatusBadRequest, map[string]interface{}{
				"success": false,
				"message": "Mobile name is required",
			})
			return
		}
		accountData := template.AccountInfo{
			Name:               userData.NameClient,
			UserMerchant:       userCredential.UserMerchant,
			PasswordMerchant:   userCredential.PasswordMerchant,
			UserSuperadmin:     userCredential.UserSuperadmin,
			PasswordSuperadmin: userCredential.PasswordSuperadmin,
			AppMobile:          userData.AppMobile,
			HomePage:           homePage,
			EmailAdmin:         mailAdmin,
			Time:               time.Date(getYear, time.January, 1, 0, 0, 0, 0, time.Local),
			PhoneCenter:        phoneCenter,
			GuideBook:          guideBook,
		}
		email.Subject, email.Body = template.Credentials(accountData)
	case "rejected":
		if r.Header.Get("Content-Type") != "application/json" {
			userData.NameClient = r.FormValue("client_name")
		}
		if userData.NameClient == "" {
			utils.RespondJSON(w, http.StatusBadRequest, map[string]interface{}{
				"success": false,
				"message": "Client name is required",
			})
			return
		}
		rejectData := template.RejectTemplate{
			Name:        userData.NameClient,
			HomePage:    homePage,
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
			utils.RespondJSON(w, http.StatusBadRequest, map[string]interface{}{
				"success": false,
				"message": "Client name is required",
			})
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
	default:
		utils.RespondJSON(w, http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": "Template not found",
		})
		return
	}
	if email.TypeMessage == "" {
		email.TypeMessage = "text/html"
	}
	info, err := mail.SendEmail(email.To, email.Cc, email.Bcc, email.Subject, email.TypeMessage, email.Body, email.Attach)
	if err != nil {
		utils.RespondJSON(w, http.StatusUnprocessableEntity, map[string]interface{}{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	for _, filePath := range email.Attach {
		os.Remove(filePath)
	}

	utils.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": info,
	})
}
