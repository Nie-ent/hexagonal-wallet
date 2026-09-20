package model

import "time"

type Wallet struct {
	ID        string
	UserID    string
	Balance   float64
	CreatedAt time.Time
	UpdatedAt time.Time
}
