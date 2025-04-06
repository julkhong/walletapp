package api

import (
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"

	"github.com/julkhong/walletapp/server/internal/config"
	"github.com/julkhong/walletapp/server/internal/service"
)

const (
	serviceLogTag = "WALLET-SERVICE"
)

func SetupRouter(cfg *config.Config) http.Handler {
	r := chi.NewRouter()
	// recovers panic
	r.Use(middleware.Recoverer)

	logger := logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	logger.SetLevel(logrus.InfoLevel)
	logger.WithField("tag", serviceLogTag)

	walletService := service.NewWalletService(cfg, logger)

	r.Post("/wallets/{id}/deposit", walletService.DepositHandler)
	r.Post("/wallets/{id}/withdraw", walletService.WithdrawHandler)
	r.Post("/wallets/transfer", walletService.TransferHandler)
	r.Get("/wallets/{id}/balance", walletService.BalanceHandler)
	r.Get("/wallets/{id}/transactions", walletService.TransactionHistoryHandler)

	return r
}
