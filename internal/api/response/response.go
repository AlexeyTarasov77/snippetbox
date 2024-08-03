package response

import "net/http"

type Response struct {
	Success bool `json:"success"`
	Status int `json:"status"`
	Error string `json:"error,omitempty"`
}

func processStatus(status ...int) int {
	if len(status) == 0 {
		status = append(status, http.StatusInternalServerError)
	}
	return status[0]
}

// HttpError returns a Response struct with the specified error message and HTTP status code.
//
// Parameters:
//
// - msg: the error message to be included in the Response struct (optional) if not provided, the default error message for the status will be used.
// 
// - _status: the HTTP status code to be included in the Response struct (default: http.StatusInternalServerError).
//
// Returns:
// - Response: a struct containing the success status, HTTP status code, and error message.
func Error(msg string, _status ...int) Response {
	status := processStatus(_status...)
	if msg == "" {
		msg = http.StatusText(status)
	}
	return Response{
		Success: false,
		Status: status,
		Error: msg,
	}
}

func HttpError(w http.ResponseWriter, msg string, _status ...int) {
	resp := Error(msg, _status...)
	http.Error(w, resp.Error, resp.Status)
}

func Success() Response {
	return Response{
		Success: true,
		Status: http.StatusOK,
	}
}