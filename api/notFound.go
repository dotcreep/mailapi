package api

import (
	"mailapi/utils"
	"net/http"
)

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	utils.RespondJSON(w, http.StatusNotFound, map[string]interface{}{
		"success": false,
		"message": "404 not found",
	})
}
