package dto

type DequeueRequest struct {
	Type string `json:"type" san:"trim,xss"`
}
