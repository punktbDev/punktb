package dto

type (
	Response struct {
		Success bool `json:"success"`
	}
	ResponseData struct {
		Data any `json:"data"`
	}
)
