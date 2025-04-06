package service

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"

	"github.com/julkhong/walletapp/server/internal/common"
	"github.com/julkhong/walletapp/server/internal/config"
	"github.com/julkhong/walletapp/server/internal/dao"
	dto "github.com/julkhong/walletapp/server/internal/dto"
	"github.com/julkhong/walletapp/server/internal/logic"
)

type WalletService struct {
	logger *logrus.Logger
	Dao    dao.WalletDaoInterface
	Impl   logic.WalletImplInterface
}

func NewWalletService(cfg *config.Config, logger *logrus.Logger) *WalletService {
	dao, err := dao.NewWalletDao(cfg, logger)
	if err != nil {
		logger.Error("failed to init wallet service", err)
	}
	impl := logic.NewWalletImpl(dao, logger)

	return &WalletService{Dao: dao, logger: logger, Impl: impl}
}

func isUUID(w http.ResponseWriter, id, label string) bool {
	if !common.IsValidUUID(id) {
		common.WriteError(w, http.StatusBadRequest, common.ErrInvalidRequest, "Invalid "+label)
		return false
	}
	return true
}
func (s *WalletService) DepositHandler(w http.ResponseWriter, r *http.Request) {
	idempotencyKey := r.Header.Get("Idempotency-Key")
	if idempotencyKey == "" {
		common.WriteError(w, http.StatusBadRequest, common.ErrInvalidRequest, "Missing Idempotency-Key")
		return
	}

	walletID := chi.URLParam(r, "id")
	if walletID == "" || !isUUID(w, walletID, "wallet_id") {
		return
	}

	if record, found := s.Dao.CheckIdempotencyKey(idempotencyKey, r.Method, r.URL.Path); found {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(record.StatusCode)
		_, _ = w.Write([]byte(record.Response))
		return
	}

	var req dto.DepositRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		common.WriteError(w, http.StatusBadRequest, common.ErrInvalidRequest, "Invalid request")
		return
	}

	if err := s.Impl.Deposit(r.Context(), walletID, req.Amount); err != nil {
		s.logger.WithError(err).Error("Deposit failed")
		switch err {
		case logic.ErrWalletNotFound:
			common.WriteError(w, http.StatusNotFound, common.ErrWalletNotFound, err.Error())
		default:
			common.WriteError(w, http.StatusInternalServerError, common.ErrUnknown, "Deposit failed")
		}
		return
	}

	resp := dto.GenericResponse[dto.SuccessResponse]{
		Status: "success",
		Data:   dto.SuccessResponse{Message: "deposit success"},
	}
	respJSON, _ := json.Marshal(resp)

	_ = s.Dao.SaveIdempotencyKey(&dao.IdempotencyRecord{
		Key:        idempotencyKey,
		Method:     r.Method,
		Path:       r.URL.Path,
		Response:   string(respJSON),
		StatusCode: http.StatusOK,
	})

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(respJSON)
}

func (s *WalletService) WithdrawHandler(w http.ResponseWriter, r *http.Request) {
	idempotencyKey := r.Header.Get("Idempotency-Key")
	if idempotencyKey == "" {
		common.WriteError(w, http.StatusBadRequest, common.ErrInvalidRequest, "Missing Idempotency-Key")
		return
	}

	walletID := chi.URLParam(r, "id")
	if walletID == "" || !isUUID(w, walletID, "wallet_id") {
		return
	}

	if record, found := s.Dao.CheckIdempotencyKey(idempotencyKey, r.Method, r.URL.Path); found {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(record.StatusCode)
		_, _ = w.Write([]byte(record.Response))
		return
	}

	var req dto.WithdrawRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		common.WriteError(w, http.StatusBadRequest, common.ErrInvalidRequest, "Invalid request")
		return
	}

	if err := s.Impl.Withdraw(r.Context(), walletID, req.Amount); err != nil {
		s.logger.WithError(err).Error("Withdraw failed")
		switch err {
		case logic.ErrWalletNotFound:
			common.WriteError(w, http.StatusNotFound, common.ErrWalletNotFound, err.Error())
		case logic.ErrInsufficientBalance:
			common.WriteError(w, http.StatusBadRequest, common.ErrInsufficientBalance, err.Error())
		default:
			common.WriteError(w, http.StatusInternalServerError, common.ErrUnknown, "Withdraw failed")
		}
		return
	}

	resp := dto.GenericResponse[dto.SuccessResponse]{
		Status: "success",
		Data:   dto.SuccessResponse{Message: "withdraw success"},
	}
	respJSON, _ := json.Marshal(resp)

	_ = s.Dao.SaveIdempotencyKey(&dao.IdempotencyRecord{
		Key:        idempotencyKey,
		Method:     r.Method,
		Path:       r.URL.Path,
		Response:   string(respJSON),
		StatusCode: http.StatusOK,
	})

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(respJSON)
}

func (s *WalletService) TransferHandler(w http.ResponseWriter, r *http.Request) {
	idempotencyKey := r.Header.Get("Idempotency-Key")
	if idempotencyKey == "" {
		common.WriteError(w, http.StatusBadRequest, common.ErrInvalidRequest, "Missing Idempotency-Key")
		return
	}

	if record, found := s.Dao.CheckIdempotencyKey(idempotencyKey, r.Method, r.URL.Path); found {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(record.StatusCode)
		_, _ = w.Write([]byte(record.Response))
		return
	}

	var req dto.TransferRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		common.WriteError(w, http.StatusBadRequest, common.ErrInvalidRequest, "Invalid request")
		return
	}

	if !isUUID(w, req.FromWalletID, "from_wallet_id") || !isUUID(w, req.ToWalletID, "to_wallet_id") {
		return
	}

	if req.Amount <= 0.1 {
		common.WriteError(w, http.StatusBadRequest, common.ErrInvalidRequest, "Amount must be more than 0.1")
		return
	}

	if err := s.Impl.Transfer(r.Context(), req.FromWalletID, req.ToWalletID, req.Amount); err != nil {
		s.logger.WithError(err).Error("Transfer failed")
		switch err {
		case logic.ErrWalletNotFound:
			common.WriteError(w, http.StatusNotFound, common.ErrWalletNotFound, err.Error())
		case logic.ErrInsufficientBalance:
			common.WriteError(w, http.StatusBadRequest, common.ErrInsufficientBalance, err.Error())
		default:
			common.WriteError(w, http.StatusInternalServerError, common.ErrUnknown, "Transfer failed")
		}
		return
	}

	balance, err := s.Impl.GetBalance(r.Context(), req.FromWalletID)
	if err != nil {
		s.logger.WithError(err).Warn("Failed to fetch updated sender balance")
	}

	resp := dto.GenericResponse[dto.TransferResponse]{
		Status: "success",
		Data: dto.TransferResponse{
			Message:  "transfer success",
			WalletID: req.FromWalletID,
			Balance:  balance,
		},
	}
	respJSON, _ := json.Marshal(resp)

	_ = s.Dao.SaveIdempotencyKey(&dao.IdempotencyRecord{
		Key:        idempotencyKey,
		Method:     r.Method,
		Path:       r.URL.Path,
		Response:   string(respJSON),
		StatusCode: http.StatusOK,
	})

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(respJSON)
}

func (s *WalletService) BalanceHandler(w http.ResponseWriter, r *http.Request) {
	walletID := chi.URLParam(r, "id")
	if walletID == "" || !isUUID(w, walletID, "wallet_id") {
		return
	}

	balance, err := s.Impl.GetBalance(r.Context(), walletID)
	if err != nil {
		s.logger.WithError(err).Error("Balance fetch failed")
		if err == logic.ErrWalletNotFound {
			common.WriteError(w, http.StatusNotFound, common.ErrWalletNotFound, err.Error())
		} else {
			common.WriteError(w, http.StatusInternalServerError, common.ErrDatabase, "Failed to get balance")
		}
		return
	}

	common.WriteJSON(w, http.StatusOK, dto.GenericResponse[dto.BalanceResponse]{
		Status: "success",
		Data: dto.BalanceResponse{
			WalletID: walletID,
			Balance:  balance,
		},
	})
}

func (s *WalletService) TransactionHistoryHandler(w http.ResponseWriter, r *http.Request) {
	walletID := chi.URLParam(r, "id")
	if walletID == "" || !isUUID(w, walletID, "wallet_id") {
		return
	}

	query := r.URL.Query()
	txType := query.Get("type")
	start := query.Get("start")
	end := query.Get("end")
	limit, _ := strconv.Atoi(query.Get("limit"))
	offset, _ := strconv.Atoi(query.Get("offset"))

	txs, err := s.Impl.GetTransactionHistory(r.Context(), walletID, txType, start, end, limit, offset)
	if err != nil {
		s.logger.WithError(err).Error("Failed to fetch transaction history")
		common.WriteError(w, http.StatusInternalServerError, common.ErrDatabase, "Failed to fetch transaction history")
		return
	}

	common.WriteJSON(w, http.StatusOK, dto.GenericResponse[[]dao.Transaction]{
		Status: "success",
		Data:   txs,
	})
}
