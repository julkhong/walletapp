package dao

import "time"

type User struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

type Wallet struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Balance   float64   `json:"balance"`
	CreatedAt time.Time `json:"created_at"`
}

type UpdateBalance struct {
	WalletID string  `json:"wallet_id"`
	Amount   float64 `json:"amount"`
}

type Transaction struct {
	ID            string    `json:"id"`
	WalletID      string    `json:"wallet_id"`
	Type          string    `json:"type"`
	Amount        float64   `json:"amount"`
	RelatedUserID *string   `json:"related_user_id"`
	CreatedAt     time.Time `json:"created_at"`
}

type IdempotencyRecord struct {
	Key        string    `gorm:"primaryKey;column:key"`
	Method     string    `gorm:"column:method"`
	Path       string    `gorm:"column:path"`
	Response   string    `gorm:"column:response"`
	StatusCode int       `gorm:"column:status_code"`
	CreatedAt  time.Time `gorm:"column:created_at"`
}
