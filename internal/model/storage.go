package model

type Storage struct {
    ID          uint    `json:"id" gorm:"primaryKey"`
    Name        string  `json:"name"`
    Price       float64 `json:"price"`
    ImageURL    string  `json:"image_url"`
    UserID      uint    `json:"user_id"`
}