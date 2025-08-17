package response

import (
	"encoding/json"
	"net/http"
)

type APIResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Type    string      `json:"type"`
	Data    interface{} `json:"data"`
}

func JSON(w http.ResponseWriter, code int, resp APIResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(resp)
}

func Error(w http.ResponseWriter, code int, errType, message string) {
	JSON(w, code, APIResponse{
		Code:    code,
		Type:    errType,
		Message: message,
	})
}

func Success(w http.ResponseWriter, data interface{}) {
	JSON(w, http.StatusOK, APIResponse{
		Code:    http.StatusOK,
		Type:    "Success",
		Message: "OK",
		Data:    data,
	})
}
