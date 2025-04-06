package logic

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"github.com/julkhong/walletapp/server/internal/common"
	dao "github.com/julkhong/walletapp/server/internal/dao"
)

var (
	ErrInsufficientBalance = errors.New("insufficient balance")
	ErrWalletNotFound      = errors.New("wallet not found")
)

type WalletImpl struct {
	dao    dao.WalletDaoInterface
	logger *logrus.Entry
}

func NewWalletImpl(dao dao.WalletDaoInterface, baseLogger *logrus.Logger) *WalletImpl {
	logger := baseLogger.WithField("tag", "WALLET-LOGIC")
	return &WalletImpl{dao: dao, logger: logger}
}

func (l *WalletImpl) Deposit(ctx context.Context, walletID string, amount float64) error {
	amount = common.RoundToNDecimals(amount, 4)
	l.logger.Infof("Depositing %.4f into wallet %s", amount, walletID)

	current, err := l.dao.GetBalance(walletID)
	if err != nil {
		l.logger.WithError(err).Errorf("Failed to fetch wallet balance")
		return fmt.Errorf("deposit failed: %w", err)
	}

	err = l.dao.UpdateBalance(&dao.UpdateBalance{
		WalletID: walletID,
		Amount:   common.RoundToNDecimals(current+amount, 4),
	})
	var related *string
	_ = l.dao.CreateTransaction(&dao.Transaction{
		ID:            uuid.NewString(),
		WalletID:      walletID,
		Type:          common.TransactionTypeDeposit,
		Amount:        amount,
		RelatedUserID: related,
		CreatedAt:     time.Now(),
	})

	if err != nil {
		l.logger.WithError(err).Errorf("Failed to update balance")
		return fmt.Errorf("deposit failed: %w", err)
	}

	return nil
}

func (l *WalletImpl) Withdraw(ctx context.Context, walletID string, amount float64) error {
	amount = common.RoundToNDecimals(amount, 4)
	l.logger.Infof("Withdrawing %.4f from wallet %s", amount, walletID)

	current, err := l.dao.GetBalance(walletID)
	if err != nil {
		l.logger.WithError(err).Errorf("Failed to fetch wallet balance")
		if errors.Is(err, dao.ErrWalletNotFound) {
			return ErrWalletNotFound
		}
		return fmt.Errorf("withdraw failed: %w", err)
	}

	if current < amount {
		l.logger.Warnf("Insufficient balance: current=%.4f, requested=%.4f", current, amount)
		return ErrInsufficientBalance
	}

	err = l.dao.UpdateBalance(&dao.UpdateBalance{
		WalletID: walletID,
		Amount:   common.RoundToNDecimals(current-amount, 4),
	})
	var related *string
	_ = l.dao.CreateTransaction(&dao.Transaction{
		ID:            uuid.NewString(),
		WalletID:      walletID,
		Type:          common.TransactionTypeWithdraw,
		Amount:        amount,
		RelatedUserID: related,
		CreatedAt:     time.Now(),
	})

	if err != nil {
		l.logger.WithError(err).Errorf("Failed to update balance after withdraw")
		return fmt.Errorf("withdraw failed: %w", err)
	}

	return nil
}

func (l *WalletImpl) Transfer(ctx context.Context, fromWalletID, toWalletID string, amount float64) error {
	amount = common.RoundToNDecimals(amount, 4)
	l.logger.Infof("Transferring %.4f from wallet %s to %s", amount, fromWalletID, toWalletID)

	fromBalance, err := l.dao.GetBalance(fromWalletID)
	if err != nil {
		l.logger.WithError(err).Error("Failed to get sender balance")
		if errors.Is(err, dao.ErrWalletNotFound) {
			return fmt.Errorf("sender wallet not found: %w", err)
		}
		return fmt.Errorf("transfer failed: %w", err)
	}

	if fromBalance < amount {
		l.logger.Warnf("Insufficient funds for transfer: current=%.4f, requested=%.4f", fromBalance, amount)
		return ErrInsufficientBalance
	}

	toBalance, err := l.dao.GetBalance(toWalletID)
	if err != nil {
		l.logger.WithError(err).Error("Failed to get receiver balance")
		if errors.Is(err, dao.ErrWalletNotFound) {
			return fmt.Errorf("receiver wallet not found: %w", err)
		}
		return fmt.Errorf("transfer failed: %w", err)
	}

	// Update balances
	if err := l.dao.UpdateBalance(&dao.UpdateBalance{
		WalletID: fromWalletID,
		Amount:   common.RoundToNDecimals(fromBalance-amount, 4),
	}); err != nil {
		l.logger.WithError(err).Error("Failed to update sender balance")
		return fmt.Errorf("transfer failed: %w", err)
	}

	if err := l.dao.UpdateBalance(&dao.UpdateBalance{
		WalletID: toWalletID,
		Amount:   common.RoundToNDecimals(toBalance+amount, 4),
	}); err != nil {
		l.logger.WithError(err).Error("Failed to update receiver balance")
		return fmt.Errorf("transfer failed: %w", err)
	}

	// Create transactions
	_ = l.dao.CreateTransaction(&dao.Transaction{
		ID:            uuid.NewString(),
		WalletID:      fromWalletID,
		Type:          common.TransactionTypeTransfer,
		Amount:        -amount,
		RelatedUserID: &toWalletID,
		CreatedAt:     time.Now(),
	})

	_ = l.dao.CreateTransaction(&dao.Transaction{
		ID:            uuid.NewString(),
		WalletID:      toWalletID,
		Type:          common.TransactionTypeTransfer,
		Amount:        amount,
		RelatedUserID: &fromWalletID,
		CreatedAt:     time.Now(),
	})

	return nil
}

func (l *WalletImpl) GetBalance(ctx context.Context, walletID string) (float64, error) {
	balance, err := l.dao.GetBalance(walletID)
	if err != nil {
		l.logger.WithError(err).Errorf("Failed to get balance for wallet %s", walletID)
		if errors.Is(err, dao.ErrWalletNotFound) {
			return 0, ErrWalletNotFound
		}
		return 0, err
	}
	return common.RoundToNDecimals(balance, 4), nil
}

func (l *WalletImpl) GetTransactionHistory(ctx context.Context, walletID, txType, start, end string, limit, offset int) ([]dao.Transaction, error) {
	return l.dao.GetTransactionHistory(walletID, txType, start, end, limit, offset)
}
