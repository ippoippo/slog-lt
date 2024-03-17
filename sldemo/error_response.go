package sldemo

type ErrorResponse struct {
	Message string `json:"message"`
	// Would normally be other useful things here
}

func ErrWithMsg(msg string) *ErrorResponse {
	return &ErrorResponse{
		Message: msg,
	}
}
