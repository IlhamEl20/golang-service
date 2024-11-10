package models

// Product defines the structure of the product model
type Product struct {
	ID    uint    `json:"id" gorm:"primaryKey"` // Primary key for the product
	Name  string  `json:"name"`                 // Name of the product
	Price float64 `json:"price"`                // Price of the product
}
