package dto

// Response represents a standard API response
type Response struct {
	Status string      `json:"status"`
	Data   interface{} `json:"data,omitempty"`
	Error  string      `json:"error,omitempty"`
}

// TokenResponse represents JWT token response
type TokenResponse struct {
	Status string `json:"status"`
	Token  string `json:"token"`
}

// SuccessResponse represents a success response
type SuccessResponse struct {
	Status      string `json:"status"`
	UpdateCount int64  `json:"update_count,omitempty"`
}
