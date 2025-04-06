package api

// TransferDTO represents a transfer request from one user to another
type TransferDTO struct {
	FromWalletID string  `json:"from_wallet_id"`
	ToWalletID   string  `json:"to_wallet_id"`
	Amount       float64 `json:"amount"`
}

type DepositRequest struct {
	Amount float64 `json:"amount" binding:"required,gt=0"`
}

type WithdrawRequest struct {
	Amount float64 `json:"amount" binding:"required,gt=0"`
}

type TransferRequest struct {
	FromWalletID string  `json:"from_wallet_id" binding:"required"`
	ToWalletID   string  `json:"to_wallet_id" binding:"required"`
	Amount       float64 `json:"amount" binding:"required,gt=0"`
}

type TransferResponse struct {
	Message  string  `json:"message"`
	WalletID string  `json:"wallet_id"`
	Balance  float64 `json:"balance"`
}

type BalanceResponse struct {
	WalletID string  `json:"wallet_id"`
	Balance  float64 `json:"balance"`
}

type SuccessResponse struct {
	Message string `json:"message"`
}

type TransactionHistoryResponse struct {
	WalletID     string   `json:"wallet_id"`
	Transactions []string `json:"transactions"`
}

type GenericResponse[T any] struct {
	Status string `json:"status"`
	Data   T      `json:"data"`
}
