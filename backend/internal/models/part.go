package models

import "time"

// Part представляет запасную часть на складе.
type Part struct {
	ID            uint      `json:"id" gorm:"primaryKey"`
	Name          string    `json:"name" gorm:"size:255;not null"`
	ArticleNumber string    `json:"article_number" gorm:"size:100"`
	Quantity      int       `json:"quantity" gorm:"default:0"`
	Unit          string    `json:"unit" gorm:"size:20"`
	MinQuantity   int       `json:"min_quantity" gorm:"default:0"`
	Location      string    `json:"location" gorm:"size:255"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
