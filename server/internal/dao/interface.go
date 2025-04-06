package dao

//go:generate mockery --name=WalletDaoInterface --output=mocks --outpkg=mocks
type WalletDaoInterface interface {
	GetBalance(walletID string) (float64, error)
	UpdateBalance(input *UpdateBalance) error
	CreateTransaction(tx *Transaction) error
	GetWalletByID(walletID string) (*Wallet, error)
	GetTransactionHistory(walletID, txType, start, end string, limit, offset int) ([]Transaction, error)
	SaveIdempotencyKey(record *IdempotencyRecord) error
	CheckIdempotencyKey(key, method, path string) (*IdempotencyRecord, bool)
}
