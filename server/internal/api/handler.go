package api

import (
    "net/http"
    "github.com/go-chi/chi/v5"
    "github.com/julkhong/walletapp/server/internal/service"
    "github.com/julkhong/walletapp/server/internal/config"
)

func SetupRouter(cfg *config.Config) http.Handler {
    r := chi.NewRouter()

    walletService := service.NewWalletService(cfg)

    r.Post("/wallets/{id}/deposit", walletService.DepositHandler)
    r.Post("/wallets/{id}/withdraw", walletService.WithdrawHandler)
    r.Post("/wallets/{id}/transfer", walletService.TransferHandler)
    r.Get("/wallets/{id}/balance", walletService.BalanceHandler)
    r.Get("/wallets/{id}/transactions", walletService.TransactionHistoryHandler)

    return r
}
