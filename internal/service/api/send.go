package api

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/dotcreep/mailapi/internal/service/mail"
	"github.com/dotcreep/mailapi/internal/utils"
)

type Email struct {
	To          string   `json:"to" form:"to"`
	Cc          string   `json:"cc" form:"cc"`
	Bcc         string   `json:"bcc" form:"bcc"`
	TypeMessage string   `json:"type" form:"type"`
	Subject     string   `json:"subject" form:"subject"`
	Body        string   `json:"body" form:"body"`
	Attach      []string `json:"attach" form:"attach"`
}

func SendEmailAPI(w http.ResponseWriter, r *http.Request) {
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
	if tmpPath == "" {
		Json.NewResponse(false, w, nil, "Storage path is required", http.StatusBadRequest, nil)
		return
	}

	var email Email
	if r.Header.Get("Content-Type") == "application/json" {
		err := json.NewDecoder(r.Body).Decode(&email)
		if err != nil {
			Json.NewResponse(false, w, nil, "Invalid JSON format", http.StatusBadRequest, nil)
			return
		}
	} else {
		err := r.ParseMultipartForm(10 << 20)
		if err != nil {
			Json.NewResponse(false, w, nil, "Unable to parse form", http.StatusBadRequest, nil)
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
				Json.NewResponse(false, w, nil, "File size is too large. Max size is 10MB", http.StatusBadRequest, nil)
				return
			}
			file, err := fileHeader.Open()
			if err != nil {
				Json.NewResponse(false, w, nil, "Unable open file", http.StatusBadRequest, nil)
				return
			}
			defer file.Close()
			tempFilePath := filepath.Join(tmpPath, fileHeader.Filename)
			out, err := os.Create(tempFilePath)
			if err != nil {
				Json.NewResponse(false, w, nil, "Unable create file", http.StatusBadRequest, nil)
				return
			}
			defer out.Close()
			if _, err = io.Copy(out, file); err != nil {
				Json.NewResponse(false, w, nil, "Unable save file", http.StatusBadRequest, nil)
				return
			}
			email.Attach = append(email.Attach, tempFilePath)
		}
	}

	if email.To == "" {
		Json.NewResponse(false, w, nil, "To field is required", http.StatusBadRequest, nil)
		return
	}
	if email.Subject == "" {
		Json.NewResponse(false, w, nil, "Subject field is required", http.StatusBadRequest, nil)
		return
	}
	if email.Body == "" {
		Json.NewResponse(false, w, nil, "Body field is required", http.StatusBadRequest, nil)
		return
	}
	if email.TypeMessage == "" {
		email.TypeMessage = "text/html"
	}
	info, err := mail.SendEmail(email.To, email.Cc, email.Bcc, email.Subject, email.TypeMessage, email.Body, email.Attach)
	if err != nil {
		Json.NewResponse(false, w, nil, err.Error(), http.StatusUnprocessableEntity, nil)
		return
	}

	for _, filePath := range email.Attach {
		os.Remove(filePath)
	}

	Json.NewResponse(true, w, nil, info, http.StatusOK, nil)
}
