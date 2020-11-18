package entity

type Response struct {
	Message string      `json:"message"`
	Status  int         `json:"status"`
	Data    interface{} `json:"data"`
}

func ToResponse(message string, status int, data interface{}) *Response {
	return &Response{
		Message: message,
		Status:  status,
		Data:    data,
	}
}
