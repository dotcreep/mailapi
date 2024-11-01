package api

import (
	"encoding/json"
	"io"
	"mailapi/mail"
	"mailapi/utils"
	"net/http"
	"os"
	"path/filepath"
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
	if r.Method != http.MethodPost {
		utils.RespondJSON(w, http.StatusMethodNotAllowed, map[string]interface{}{
			"success": false,
			"message": "Method not allowed",
		})
		return
	}
	tmpPath := os.Getenv("STORAGE")
	if tmpPath == "" {
		utils.RespondJSON(w, http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": "Storage path is required",
		})
		return
	}

	var email Email
	if r.Header.Get("Content-Type") == "application/json" {
		err := json.NewDecoder(r.Body).Decode(&email)
		if err != nil {
			utils.RespondJSON(w, http.StatusBadRequest, map[string]interface{}{
				"success": false,
				"message": "Invalid JSON format",
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
	if email.Subject == "" {
		utils.RespondJSON(w, http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": "Subject field is required",
		})
		return
	}
	if email.Body == "" {
		utils.RespondJSON(w, http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": "Body field is required",
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
