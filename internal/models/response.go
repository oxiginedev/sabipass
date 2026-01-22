package models

type Response struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
	Errors  any    `json:"errors,omitempty"`
}

func NewSuccessResponse(message string, data any) Response {
	return newResponse(true, message, data, nil)
}

func NewErrorResponse(message string, errors any) Response {
	return newResponse(false, message, nil, errors)
}

func newResponse(status bool, message string, data, errors any) Response {
	return Response{
		Status:  status,
		Message: message,
		Data:    data,
		Errors:  errors,
	}
}
