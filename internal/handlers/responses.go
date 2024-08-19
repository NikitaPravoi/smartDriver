package handlers

// SuccessOutput is a basic `"success": true` response
type SuccessOutput struct {
	Body struct {
		Success bool `json:"success" example:"true" doc:"Status of succession"`
	}
}

// SuccessErrorOutput is a basic `"success": true` response
// with error explanation added
type SuccessErrorOutput struct {
	Body struct {
		Success bool   `json:"success" example:"true" doc:"Status of succession"`
		Error   string `json:"error" example:"An exception occurred during processing operation" doc:"Error message"`
	}
}
