package send_template

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/dotcreep/mailapi/internal/service/mail"
	"github.com/dotcreep/mailapi/internal/service/template"
	"github.com/dotcreep/mailapi/internal/utils"
)

// @Summary		Send email template - Credentials
// @Description	Using template Credentials
// @Tags			Send Template
// @Accept			json
// @Produce		json
// @Security		X-Auth-Key
// @Param			body	body		utils.Credentials			true	"Body"
// @Success		200		{object}	utils.Success				"Success"
// @Failure		400		{object}	utils.BadRequest			"Bad request"
// @Failure		500		{object}	utils.InternalServerError	"Internal server error"
// @Router			/send/template/credentials [post]
func Credentials(w http.ResponseWriter, r *http.Request) {
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

	var userData utils.UserData
	var userCredential utils.UserCredential
	var email utils.Email
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
