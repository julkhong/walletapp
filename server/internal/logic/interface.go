package logic

import (
	"context"

	"github.com/julkhong/walletapp/server/internal/dao"
)

//go:generate mockery --name=WalletImplInterface --output=./mocks --outpkg=mocks
type WalletImplInterface interface {
	Deposit(ctx context.Context, walletID string, amount float64) error
	Withdraw(ctx context.Context, walletID string, amount float64) error
	Transfer(ctx context.Context, fromWalletID, toWalletID string, amount float64) error
	GetBalance(ctx context.Context, walletID string) (float64, error)
	GetTransactionHistory(ctx context.Context, walletID, txType, start, end string, limit, offset int) ([]dao.Transaction, error)
}
