package models

import "time"

// Listing модель объявления
type Listing struct {
	ID          int       `json:"id" db:"id"`
	Title       string    `json:"title" db:"title"`
	Description string    `json:"description" db:"description"`
	ImageURL    *string   `json:"image_url" db:"image_url"` // pointer для nullable поля
	Price       float64   `json:"price" db:"price"`
	UserID      int       `json:"user_id" db:"user_id"`
	UserLogin   string    `json:"user_login,omitempty" db:"user_login"` // для joined запросов
	IsOwner     bool      `json:"is_owner,omitempty"`                   // признак принадлежности текущему пользователю
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// CreateListingRequest структура для создания объявления
type CreateListingRequest struct {
	Title       string  `json:"title" binding:"required,min=1,max=255"`
	Description string  `json:"description" binding:"required,min=1"`
	ImageURL    *string `json:"image_url,omitempty" binding:"omitempty,url,max=500"`
	Price       float64 `json:"price" binding:"required,gt=0"`
}

// UpdateListingRequest структура для обновления объявления
type UpdateListingRequest struct {
	Title       *string  `json:"title,omitempty" binding:"omitempty,min=1,max=255"`
	Description *string  `json:"description,omitempty" binding:"omitempty,min=1"`
	ImageURL    *string  `json:"image_url,omitempty" binding:"omitempty,url,max=500"`
	Price       *float64 `json:"price,omitempty" binding:"omitempty,gt=0"`
}

// ListingsFilter параметры фильтрации объявлений
type ListingsFilter struct {
	MinPrice *float64 `form:"min_price" binding:"omitempty,gte=0"`
	MaxPrice *float64 `form:"max_price" binding:"omitempty,gte=0"`
	SortBy   string   `form:"sort_by" binding:"omitempty,oneof=created_at price"`
	SortDir  string   `form:"sort_dir" binding:"omitempty,oneof=asc desc"`
	Page     int      `form:"page" binding:"omitempty,min=1"`
	Limit    int      `form:"limit" binding:"omitempty,min=1,max=100"`
}

// SetDefaults устанавливает значения по умолчанию для фильтра
func (f *ListingsFilter) SetDefaults() {
	if f.SortBy == "" {
		f.SortBy = "created_at"
	}
	if f.SortDir == "" {
		f.SortDir = "desc"
	}
	if f.Page == 0 {
		f.Page = 1
	}
	if f.Limit == 0 {
		f.Limit = 20
	}
}

// GetOffset возвращает offset для пагинации
func (f *ListingsFilter) GetOffset() int {
	return (f.Page - 1) * f.Limit
}

// PaginatedListings результат с пагинацией
type PaginatedListings struct {
	Data       []Listing `json:"data"`
	Total      int       `json:"total"`
	Page       int       `json:"page"`
	Limit      int       `json:"limit"`
	TotalPages int       `json:"total_pages"`
}
