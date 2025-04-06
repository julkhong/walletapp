package dao

import (
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)

	dialector := postgres.New(postgres.Config{Conn: db})
	gdb, err := gorm.Open(dialector, &gorm.Config{})
	assert.NoError(t, err)

	return gdb, mock
}

func TestSaveIdempotencyKey(t *testing.T) {
	t.Run("successfully saves record", func(t *testing.T) {
		db, mock := setupMockDB(t)
		dao := &WalletDao{db: db}

		now := time.Now()
		record := &IdempotencyRecord{
			Key:        "key-123",
			Method:     "POST",
			Path:       "/wallets/123/deposit",
			Response:   `{"message":"success"}`,
			StatusCode: 200,
			CreatedAt:  now,
		}

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "idempotency_keys"`)).
			WithArgs(
				record.Key,
				record.Method,
				record.Path,
				record.Response,
				record.StatusCode,
				record.CreatedAt,
			).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := dao.SaveIdempotencyKey(record)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestCheckIdempotencyKey(t *testing.T) {
	const queryPattern = `SELECT \* FROM "idempotency_keys" WHERE key = \$1 AND method = \$2 AND path = \$3 ORDER BY "idempotency_keys"\."key" LIMIT \$4`

	t.Run("key found in DB", func(t *testing.T) {
		db, mock := setupMockDB(t)
		dao := &WalletDao{db: db}

		now := time.Now()
		row := sqlmock.NewRows([]string{
			"key", "method", "path", "response", "status_code", "created_at",
		}).AddRow("key-123", "POST", "/wallets/123/deposit", `{"message":"success"}`, 200, now)

		mock.ExpectQuery(queryPattern).
			WithArgs("key-123", "POST", "/wallets/123/deposit", 1).
			WillReturnRows(row)

		result, ok := dao.CheckIdempotencyKey("key-123", "POST", "/wallets/123/deposit")
		assert.True(t, ok)
		assert.NotNil(t, result)
		assert.Equal(t, 200, result.StatusCode)
		assert.Equal(t, `{"message":"success"}`, result.Response)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("key not found in DB", func(t *testing.T) {
		db, mock := setupMockDB(t)
		dao := &WalletDao{db: db}

		mock.ExpectQuery(queryPattern).
			WithArgs("missing-key", "POST", "/wallets/deposit", 1).
			WillReturnError(gorm.ErrRecordNotFound)

		result, ok := dao.CheckIdempotencyKey("missing-key", "POST", "/wallets/deposit")
		assert.False(t, ok)
		assert.Nil(t, result)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
