package api

// WalletDTO is used for input/output in HTTP handlers
type WalletDTO struct {
	UserID  string `json:"user_id"`
	Balance int64  `json:"balance"`
}

// UpdateBalanceDTO is used when updating wallet balance
type UpdateBalanceDTO struct {
	WalletID string `json:"wallet_id"`
	Amount   int64  `json:"amount"`
}

// TransferDTO represents a transfer request from one user to another
type TransferDTO struct {
	FromWalletID string `json:"from_wallet_id"`
	ToWalletID   string `json:"to_wallet_id"`
	Amount       int64  `json:"amount"`
}
