package service

import (
    "encoding/json"
    "net/http"

    "github.com/julkhong/walletapp/server/internal/repository"
    "github.com/julkhong/walletapp/server/internal/config"
)

type WalletService struct {
    Repo *repository.WalletRepository
}


func NewWalletService(cfg *config.Config) *WalletService {
    repo := repository.NewWalletRepository(cfg.DB)
    return &WalletService{Repo: repo}
}

func (s *WalletService) DepositHandler(w http.ResponseWriter, r *http.Request) {
    // Placeholder logic
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{"message": "deposit success"})
}

func (s *WalletService) WithdrawHandler(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{"message": "withdraw success"})
}

func (s *WalletService) TransferHandler(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{"message": "transfer success"})
}

func (s *WalletService) BalanceHandler(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]int64{"balance": 10000}) // Example balance
}

func (s *WalletService) TransactionHistoryHandler(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode([]string{"tx1", "tx2"}) // Dummy history
}
