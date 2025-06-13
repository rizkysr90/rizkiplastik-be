package category

import "time"

type CategoryBaseModel struct {
	CategoryID   string    `json:"category_id"`
	CategoryCode string    `json:"category_code"`
	CategoryName string    `json:"category_name"`
	IsActive     bool      `json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type CategoryDetailModel struct {
	CategoryBaseModel
	Description string `json:"description"`
}
