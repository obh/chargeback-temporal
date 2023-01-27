package models

import (
	"time"

	"gorm.io/gorm"
)

type Payment struct {
	gorm.Model
	// Id         int
	Currency   string
	Amount     float64
	Reference  string
	PaidOn     time.Time
	Status     string
	CustomerID int
	Customer   Customer
}
