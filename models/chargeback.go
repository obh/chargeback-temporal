package models

import (
	"gorm.io/gorm"
)

type ChargebackRequest struct {
	PaymentId int     `json:"payment_id"`
	Reason    string  `json:"chargeback_reason"`
	Amount    float32 `json:"chargeback_amount"`
}

type Chargeback struct {
	gorm.Model
	ChargebackRequest
}
type ChargebackNotifyMerchant struct {
	ChargebackID int `json:"chargeback_id"`
	PaymentId    int `json:"payment_id"`
}
