package models

import (
	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model
	FirstName  string  `json:"first_name"`          // User's first name
	LastName   string  `json:"last_name"`           // User's last name
	Email      string  `json:"email" gorm:"unique"` // User's email, must be unique
	Age        int     `json:"age"`                 // User's age
	Gender     string  `json:"gender"`              // Gender (e.g., "Male", "Female", "Other")
	Department string  `json:"department"`          // User's department
	Company    string  `json:"company"`             // User's company name
	Salary     float64 `json:"salary"`              // User's salary
	DateJoined string  `json:"date_joined"`         // Date when the user joined (ISO format: YYYY-MM-DD)
	IsActive   bool    `json:"is_active"`           // Active status of the user
}
