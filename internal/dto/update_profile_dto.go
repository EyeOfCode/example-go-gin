package dto

type UpdateProfileRequest struct {
	Name string `json:"name" binding:"required,min=3,max=30"`
}
