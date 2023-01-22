package models

import (
	"time"

	"gorm.io/gorm"
)

type ChargebackRequest struct {
	gorm.Model
	PaymentId int     `json:"payment_id"`
	Reason    string  `json:"chargeback_reason"`
	Amount    float32 `json:"chargeback_amount"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Chargeback struct {
	ChargebackRequest
}
