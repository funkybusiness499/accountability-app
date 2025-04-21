package api

// ErrorResponse represents a standard error response structure
type ErrorResponse struct {
	Error string `json:"error"`
}

// SuccessResponse represents a standard success response structure
type SuccessResponse struct {
	Message string `json:"message"`
}
