package core_http_response

type ErrorResponse struct {
	Error   string `json:"erorr" example:"error"`
	Message string `json:"message" example:"error message"`
}
