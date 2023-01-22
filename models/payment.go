package models

import "time"

type Payment struct {
	Id        int
	Amount    float64
	Reference string
	PaidOn    time.Time
	Status    string
}
