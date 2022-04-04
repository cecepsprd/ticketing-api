package model

import (
	"time"
)

type Product struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float32   `json:"price"`
	Stock       int64     `json:"stock"`
	ImageURL    string    `json:"image_url"`
	StartDate   string    `json:"start_date"`
	EndDate     string    `json:"end_date"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type ProductRequest struct {
	ID          string  `json:"_id"`
	Name        string  `json:"name" validate:"required,min=3,max=45"`
	Description string  `json:"description" validate:"required"`
	Price       float32 `json:"price"`
	Stock       int64   `json:"stock"`
	ImageURL    string  `json:"image_url"`
	StartDate   string  `json:"start_date"`
	EndDate     string  `json:"end_date"`
}
