package utils

import (
	"chargebackapp/models"
	"time"

	"gorm.io/gorm"
)

func InsertPayment(db *gorm.DB) {
	var t time.Time
	t = t.AddDate(2022, 12, 5)
	t = t.Add(time.Hour*time.Duration(12) + time.Minute*time.Duration(43) + time.Second + time.Duration(45))

	payment := &models.Payment{
		Reference: "REF_1212399812312",
		PaidOn:    t,
		Currency:  "INR",
		Amount:    1000,
		Status:    "SUCCESS",
		Customer: models.Customer{
			Name:  "Rohit S",
			Email: "rohit@cashfree.com",
			Phone: "9909912345",
		},
	}
	db.Create(payment)
}

func GetMerchant() models.Merchant {
	return models.Merchant{
		Name:         "Test merchant",
		PrimaryEmail: "rohit+merchant@cashfree.com",
	}
}
