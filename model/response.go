package model

type ResponseError struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
}
