package errors

// ApiError ...
type ApiError struct {
	Message  string `json:"message"`
	ErrorStr string `json:"error"`
	Status   int    `json:"status"`
}

func (e ApiError) Error() string {
	return e.Message
}
