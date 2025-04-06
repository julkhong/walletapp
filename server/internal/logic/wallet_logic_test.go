package logic_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/julkhong/walletapp/server/internal/common"
	"github.com/julkhong/walletapp/server/internal/dao"
	"github.com/julkhong/walletapp/server/internal/dao/mocks"
	"github.com/julkhong/walletapp/server/internal/logic"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupLogicTest() (*logic.WalletImpl, *mocks.WalletDaoInterface) {
	mockDao := new(mocks.WalletDaoInterface)
	logger := logrus.New()
	impl := logic.NewWalletImpl(mockDao, logger)
	return impl, mockDao
}

func TestDeposit(t *testing.T) {
	impl, mockDao := setupLogicTest()
	ctx := context.TODO()
	walletID := "wallet-1"
	amount := 10.12345

	t.Run("successful deposit", func(t *testing.T) {
		mockDao.On("GetBalance", walletID).Return(90.0, nil).Once()
		mockDao.On("UpdateBalance", mock.Anything).Return(nil).Once()
		mockDao.On("CreateTransaction", mock.Anything).Return(nil).Once()

		err := impl.Deposit(ctx, walletID, amount)
		assert.NoError(t, err)
		mockDao.AssertExpectations(t)
	})

	t.Run("wallet not found", func(t *testing.T) {
		mockDao.On("GetBalance", walletID).Return(0.0, dao.ErrWalletNotFound).Once()

		err := impl.Deposit(ctx, walletID, amount)
		assert.Error(t, err)
		mockDao.AssertExpectations(t)
	})

	t.Run("update failed", func(t *testing.T) {
		mockDao.On("GetBalance", walletID).Return(100.0, nil).Once()
		mockDao.On("UpdateBalance", mock.Anything).Return(errors.New("update failed")).Once()
		mockDao.On("CreateTransaction", mock.Anything).Return(nil).Once()

		err := impl.Deposit(ctx, walletID, amount)
		assert.Error(t, err)
		mockDao.AssertExpectations(t)
	})
}

func TestWithdraw(t *testing.T) {
	impl, mockDao := setupLogicTest()
	ctx := context.TODO()
	walletID := "wallet-2"
	amount := 20.00

	t.Run("successful withdraw", func(t *testing.T) {
		mockDao.On("GetBalance", walletID).Return(100.0, nil).Once()
		mockDao.On("UpdateBalance", mock.Anything).Return(nil).Once()
		mockDao.On("CreateTransaction", mock.Anything).Return(nil).Once()

		err := impl.Withdraw(ctx, walletID, amount)
		assert.NoError(t, err)
		mockDao.AssertExpectations(t)
	})

	t.Run("insufficient balance", func(t *testing.T) {
		mockDao.On("GetBalance", walletID).Return(10.0, nil).Once()

		err := impl.Withdraw(ctx, walletID, amount)
		assert.ErrorIs(t, err, logic.ErrInsufficientBalance)
		mockDao.AssertExpectations(t)
	})

	t.Run("wallet not found", func(t *testing.T) {
		mockDao.On("GetBalance", walletID).Return(0.0, dao.ErrWalletNotFound).Once()

		err := impl.Withdraw(ctx, walletID, amount)
		assert.ErrorIs(t, err, logic.ErrWalletNotFound)
		mockDao.AssertExpectations(t)
	})
}

func TestTransfer(t *testing.T) {
	impl, mockDao := setupLogicTest()
	ctx := context.TODO()
	fromWallet := "wallet-from"
	toWallet := "wallet-to"
	amount := 25.0

	t.Run("successful transfer", func(t *testing.T) {
		mockDao.On("GetBalance", fromWallet).Return(100.0, nil).Once()
		mockDao.On("GetBalance", toWallet).Return(50.0, nil).Once()
		mockDao.On("UpdateBalance", mock.Anything).Return(nil).Times(2)
		mockDao.On("CreateTransaction", mock.Anything).Return(nil).Times(2)

		err := impl.Transfer(ctx, fromWallet, toWallet, amount)
		assert.NoError(t, err)
		mockDao.AssertExpectations(t)
	})

	t.Run("insufficient funds", func(t *testing.T) {
		mockDao.On("GetBalance", fromWallet).Return(10.0, nil).Once()

		err := impl.Transfer(ctx, fromWallet, toWallet, amount)
		assert.ErrorIs(t, err, logic.ErrInsufficientBalance)
		mockDao.AssertExpectations(t)
	})

	t.Run("receiver wallet not found", func(t *testing.T) {
		mockDao.On("GetBalance", fromWallet).Return(100.0, nil).Once()
		mockDao.On("GetBalance", toWallet).Return(0.0, dao.ErrWalletNotFound).Once()

		err := impl.Transfer(ctx, fromWallet, toWallet, amount)
		assert.Error(t, err)
		mockDao.AssertExpectations(t)
	})
}

func TestGetBalance(t *testing.T) {
	impl, mockDao := setupLogicTest()
	ctx := context.TODO()
	walletID := "wallet-3"

	t.Run("successful balance fetch", func(t *testing.T) {
		mockDao.On("GetBalance", walletID).Return(45.6789, nil).Once()

		balance, err := impl.GetBalance(ctx, walletID)
		assert.NoError(t, err)
		assert.Equal(t, common.RoundToNDecimals(45.6789, 4), balance)
		mockDao.AssertExpectations(t)
	})

	t.Run("wallet not found", func(t *testing.T) {
		mockDao.On("GetBalance", walletID).Return(0.0, dao.ErrWalletNotFound).Once()

		balance, err := impl.GetBalance(ctx, walletID)
		assert.ErrorIs(t, err, logic.ErrWalletNotFound)
		assert.Equal(t, 0.0, balance)
		mockDao.AssertExpectations(t)
	})
}

func TestGetTransactionHistory(t *testing.T) {
	impl, mockDao := setupLogicTest()
	ctx := context.TODO()

	mockDao.On("GetTransactionHistory", "wallet-4", "deposit", "2024-01-01", "2024-01-02", 10, 0).
		Return([]dao.Transaction{{ID: uuid.NewString()}}, nil).Once()

	txns, err := impl.GetTransactionHistory(ctx, "wallet-4", "deposit", "2024-01-01", "2024-01-02", 10, 0)
	assert.NoError(t, err)
	assert.Len(t, txns, 1)
	mockDao.AssertExpectations(t)
}
