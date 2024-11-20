package handler

// successOut is a basic `"success": true` response
type successOut struct {
	Body struct {
		Success bool `json:"success" example:"true" doc:"Status of succession"`
	}
}

// successErrorOut is a basic `"success": true` response
// with error explanation added
type successErrorOut struct {
	Body struct {
		Success bool   `json:"success" example:"true" doc:"Status of succession"`
		Error   string `json:"error" example:"An exception occurred during processing operation" doc:"Error message"`
	}
}
