package dto

type PingRequest struct {
	Url string `json:"url" binding:"required,url"`
}