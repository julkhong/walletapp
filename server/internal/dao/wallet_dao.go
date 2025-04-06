package dao

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/julkhong/walletapp/server/internal/config"
	"github.com/sirupsen/logrus"
)

var (
	ErrWalletNotFound = errors.New("wallet not found")
)

type WalletDao struct {
	db     *gorm.DB
	logger *logrus.Entry
	cfg    *config.Config
}

func NewWalletDao(cfg *config.Config, baseLogger *logrus.Logger) (*WalletDao, error) {
	logger := baseLogger.WithField("tag", "WALLET-DAO")

	db, err := gorm.Open(postgres.Open(cfg.DBURL), &gorm.Config{})
	if err != nil {
		logger.Error("Failed to connect to database", err)
		return nil, err
	}
	logger.Info("Connected to database successfully")
	return &WalletDao{db: db, logger: logger, cfg: cfg}, nil
}

func (dao *WalletDao) GetWalletByID(walletID string) (*Wallet, error) {
	dao.logger.Infof("Fetching wallet by ID: %s", walletID)

	var wallet Wallet
	result := dao.db.Table("wallets").Where("id = ?", walletID).First(&wallet)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			dao.logger.Warnf("Wallet not found: %s", walletID)
			return nil, ErrWalletNotFound
		}
		dao.logger.WithError(result.Error).Error("Failed to fetch wallet")
		return nil, result.Error
	}

	return &wallet, nil
}

func (dao *WalletDao) CreateWallet(wallet *Wallet) error {
	dao.logger.Infof("Creating wallet for user ID: %s", wallet.UserID)
	return dao.db.Table("wallets").Create(wallet).Error
}

func (dao *WalletDao) UpdateBalance(input *UpdateBalance) error {
	dao.logger.Infof("Updating balance for wallet ID: %s, amount: %.4f", input.WalletID, input.Amount)

	ctx := context.Background()
	key := fmt.Sprintf("wallet_balance:%s", input.WalletID)

	// Invalidate Redis cache first
	if err := dao.cfg.Redis.Del(ctx, key).Err(); err != nil {
		dao.logger.WithError(err).Warnf("Failed to delete Redis cache before updating balance for wallet %s", input.WalletID)
	}

	// Update in DB
	result := dao.db.Table("wallets").Where("id = ?", input.WalletID).Update("balance", input.Amount)
	if result.Error != nil {
		dao.logger.WithError(result.Error).Error("Failed to update balance in DB")
		return result.Error
	}
	if result.RowsAffected == 0 {
		dao.logger.Warnf("No wallet found for update: %s", input.WalletID)
		return ErrWalletNotFound
	}

	// Refresh Redis cache with new balance (optional but recommended)
	err := dao.cfg.Redis.Set(ctx, key, fmt.Sprintf("%.4f", input.Amount), 10*time.Minute).Err()
	if err != nil {
		dao.logger.WithError(err).Warn("Failed to update Redis cache after balance update")
	}

	return nil
}

func (dao *WalletDao) GetBalance(walletID string) (float64, error) {
	dao.logger.Infof("Getting balance for wallet ID: %s", walletID)

	key := fmt.Sprintf("wallet_balance:%s", walletID)
	cached, err := dao.cfg.Redis.Get(context.Background(), key).Result()
	if err == nil {
		val, parseErr := strconv.ParseFloat(cached, 64)
		if parseErr == nil {
			dao.logger.Infof("Cache hit for wallet %s", walletID)
			return val, nil
		}
		dao.logger.Warnf("Failed to parse cached balance: %v", parseErr)
	}

	var wallet Wallet
	result := dao.db.Table("wallets").Select("balance").Where("id = ?", walletID).First(&wallet)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			dao.logger.Warnf("Wallet not found when fetching balance: %s", walletID)
			return 0, ErrWalletNotFound
		}
		dao.logger.WithError(result.Error).Error("Failed to retrieve balance")
		return 0, result.Error
	}

	// Cache in Redis
	err = dao.cfg.Redis.Set(context.Background(), key, fmt.Sprintf("%.4f", wallet.Balance), 10*time.Minute).Err()
	if err != nil {
		dao.logger.WithError(err).Warn("Failed to cache balance in Redis")
	}

	return wallet.Balance, nil
}

func (dao *WalletDao) GetTransactionHistory(walletID string, txType string, start, end string, limit, offset int) ([]Transaction, error) {
	query := dao.db.Table("transactions").Where("wallet_id = ?", walletID)

	if txType != "" {
		query = query.Where("type = ?", txType)
	}
	if start != "" && end != "" {
		query = query.Where("created_at BETWEEN ? AND ?", start, end)
	}

	if limit > 0 {
		query = query.Limit(limit).Offset(offset)
	}

	var txs []Transaction
	err := query.Order("created_at DESC").Find(&txs).Error
	if err != nil {
		dao.logger.WithError(err).Error("Failed to fetch transaction history")
		return nil, err
	}
	return txs, nil
}

func (dao *WalletDao) CreateTransaction(tx *Transaction) error {
	dao.logger.Infof("Creating transaction for wallet %s, type: %s", tx.WalletID, tx.Type)
	return dao.db.Table("transactions").Create(tx).Error
}
