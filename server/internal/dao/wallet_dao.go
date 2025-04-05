package dao

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"github.com/julkhong/walletapp/server/internal/api"
	"github.com/julkhong/walletapp/server/internal/config"
)



type WalletDao struct {
	db *gorm.DB
}

func NewWalletDao(cfg *config.Config) (*WalletDao, error) {
	db, err := gorm.Open(postgres.Open(cfg.DBURL), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return &WalletDao{db: db}, nil
}

func (dao *WalletDao) GetWalletByID(walletID string) (*api.WalletDTO, error) {
	var wallet api.WalletDTO
	result := dao.db.Table("wallets").Where("id = ?", walletID).First(&wallet)
	if result.Error != nil {
		return nil, result.Error
	}
	return &wallet, nil
}

func (dao *WalletDao) CreateWallet(dto *api.WalletDTO) error {
	return dao.db.Table("wallets").Create(dto).Error
}

func (dao *WalletDao) UpdateBalance(dto *api.UpdateBalanceDTO) error {
	return dao.db.Table("wallets").Where("id = ?", dto.WalletID).Update("balance", dto.Amount).Error
}

func (dao *WalletDao) GetBalance(walletID string) (int64, error) {
	var balance int64
	err := dao.db.Table("wallets").Select("balance").Where("id = ?", walletID).Scan(&balance).Error
	if err != nil {
		return 0, err
	}
	return balance, nil
}