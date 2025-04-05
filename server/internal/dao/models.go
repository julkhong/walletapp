package dao

import "time"

type User struct {
	ID        string    `json:"d"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

type Wallet struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Balance   int64     `json:"balance"`
	CreatedAt time.Time `json:"created_at"`
}

type Transaction struct {
	ID            string    `json:"id"`
	WalletID      string    `json:"wallet_id"`
	Type          string    `json:"type"`
	Amount        int64     `json:"amount"`
	RelatedUserID *string   `json:"related_user_id"`
	CreatedAt     time.Time `json:"created_at"`
}
