package response

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func Ok(data interface{}) Response {
	return Response{
		Code:    0,
		Message: "success",
		Data:    data,
	}
}

func OkMsg(message string) Response {
	return Response{
		Code:    0,
		Message: message,
	}
}

func Fail(code int, message string) Response {
	return Response{
		Code:    code,
		Message: message,
	}
}

func Error(code int, message string) Response {
	return Response{
		Code:    code,
		Message: message,
	}
}

func (r Response) WriteTo(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(r)
}
