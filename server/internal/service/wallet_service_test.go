package service

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	daoMocks "github.com/julkhong/walletapp/server/internal/dao/mocks"
	logicMocks "github.com/julkhong/walletapp/server/internal/logic/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupTestService() (*WalletService, *logicMocks.WalletImplInterface, *daoMocks.WalletDaoInterface) {
	logicMock := new(logicMocks.WalletImplInterface)
	daoMock := new(daoMocks.WalletDaoInterface)

	return &WalletService{
		Impl: logicMock,
		Dao:  daoMock,
	}, logicMock, daoMock
}

func withRouteParam(r *http.Request, key, val string) *http.Request {
	ctx := chi.NewRouteContext()
	ctx.URLParams.Add(key, val)
	return r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, ctx))
}

func TestDepositHandler(t *testing.T) {
	svc, logicMock, daoMock := setupTestService()

	t.Run("DepositHandler", func(t *testing.T) {
		t.Run("valid deposit request", func(t *testing.T) {
			reqBody := `{"amount": 100}`
			req := httptest.NewRequest(http.MethodPost, "/wallets/10000000-0000-0000-0000-000000000000/deposit", strings.NewReader(reqBody))
			req.Header.Set("Idempotency-Key", "key-123")
			req = withRouteParam(req, "id", "10000000-0000-0000-0000-000000000000")

			daoMock.On("CheckIdempotencyKey", "key-123", "POST", "/wallets/10000000-0000-0000-0000-000000000000/deposit").
				Return(nil, false)
			logicMock.On("Deposit", mock.Anything, "10000000-0000-0000-0000-000000000000", 100.0).
				Return(nil)
			daoMock.On("SaveIdempotencyKey", mock.Anything).Return(nil)

			w := httptest.NewRecorder()
			svc.DepositHandler(w, req)

			assert.Equal(t, http.StatusOK, w.Code)
			assert.Contains(t, w.Body.String(), "deposit success")
		})

		t.Run("missing idempotency key", func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/wallets/10000000-0000-0000-0000-000000000000/deposit", nil)
			req = withRouteParam(req, "id", "10000000-0000-0000-0000-000000000000")

			w := httptest.NewRecorder()
			svc.DepositHandler(w, req)

			assert.Equal(t, http.StatusBadRequest, w.Code)
			assert.Contains(t, w.Body.String(), "Missing Idempotency-Key")
		})

		t.Run("invalid request body", func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/wallets/10000000-0000-0000-0000-000000000000/deposit", strings.NewReader(`bad`))
			req.Header.Set("Idempotency-Key", "key-123")
			req = withRouteParam(req, "id", "10000000-0000-0000-0000-000000000000")

			daoMock.On("CheckIdempotencyKey", "key-123", "POST", "/wallets/10000000-0000-0000-0000-000000000000/deposit").
				Return(nil, false)

			w := httptest.NewRecorder()
			svc.DepositHandler(w, req)

			assert.Equal(t, http.StatusBadRequest, w.Code)
		})
	})
}
