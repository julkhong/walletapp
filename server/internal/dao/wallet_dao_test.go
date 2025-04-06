package dao

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-redis/redismock/v9"
	"github.com/julkhong/walletapp/server/internal/config"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupTest(t *testing.T) (*WalletDao, sqlmock.Sqlmock, redismock.ClientMock) {
	db, dbMock, err := sqlmock.New()
	assert.NoError(t, err)

	gdb, err := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
	assert.NoError(t, err)

	redisClient, redisMock := redismock.NewClientMock()

	cfg := &config.Config{Redis: redisClient}
	logger := logrus.New()

	walletDao := &WalletDao{
		db:     gdb,
		logger: logger.WithField("tag", "TEST"),
		cfg:    cfg,
	}

	return walletDao, dbMock, redisMock
}

func TestGetWalletByID(t *testing.T) {
	dao, dbMock, _ := setupTest(t)
	walletID := "wallet-123"

	t.Run("wallet exists", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "user_id", "balance", "created_at"}).
			AddRow(walletID, "user-123", 100.50, time.Now())
		dbMock.ExpectQuery(`SELECT \* FROM "wallets" WHERE id = \$1 ORDER BY "wallets"\."id" LIMIT \$2`).
			WithArgs(walletID, 1).WillReturnRows(rows)

		_, err := dao.GetWalletByID(walletID)
		assert.NoError(t, err)
	})

	t.Run("wallet not found", func(t *testing.T) {
		dbMock.ExpectQuery(`SELECT \* FROM "wallets" WHERE id = \$1 ORDER BY "wallets"\."id" LIMIT \$2`).
			WithArgs(walletID, 1).WillReturnError(gorm.ErrRecordNotFound)

		wallet, err := dao.GetWalletByID(walletID)
		assert.ErrorIs(t, err, ErrWalletNotFound)
		assert.Nil(t, wallet)
	})

	t.Run("DB error", func(t *testing.T) {
		dbMock.ExpectQuery(`SELECT \* FROM "wallets" WHERE id = \$1 ORDER BY "wallets"\."id" LIMIT \$2`).
			WithArgs(walletID, 1).WillReturnError(errors.New("db error"))

		wallet, err := dao.GetWalletByID(walletID)
		assert.Error(t, err)
		assert.Nil(t, wallet)
	})
}

func TestUpdateBalance(t *testing.T) {
	dao, dbMock, redisMock := setupTest(t)
	input := &UpdateBalance{WalletID: "wallet-123", Amount: 50.1234}
	redisKey := fmt.Sprintf("wallet_balance:%s", input.WalletID)

	t.Run("successful update", func(t *testing.T) {
		redisMock.ExpectDel(redisKey).SetVal(1)

		// Begin transaction
		dbMock.ExpectBegin()

		// Expect the balance update SQL
		dbMock.ExpectExec(`UPDATE "wallets" SET "balance"=\$1 WHERE id = \$2`).
			WithArgs(input.Amount, input.WalletID).
			WillReturnResult(sqlmock.NewResult(1, 1))

		// Commit transaction
		dbMock.ExpectCommit()

		// Expect Redis set with TTL
		redisMock.ExpectSet(redisKey, "50.1234", 600*time.Second).SetVal("OK")

		err := dao.UpdateBalance(input)
		assert.NoError(t, err)
		assert.NoError(t, dbMock.ExpectationsWereMet())
		assert.NoError(t, redisMock.ExpectationsWereMet())
	})

	t.Run("no rows affected", func(t *testing.T) {
		redisMock.ExpectDel(redisKey).SetVal(1)

		// Begin transaction
		dbMock.ExpectBegin()

		// Expect balance update, but simulate no rows affected
		dbMock.ExpectExec(`UPDATE "wallets" SET "balance"=\$1 WHERE id = \$2`).
			WithArgs(input.Amount, input.WalletID).
			WillReturnResult(sqlmock.NewResult(0, 0))

		// Commit even though nothing was updated
		dbMock.ExpectCommit()

		err := dao.UpdateBalance(input)
		assert.ErrorIs(t, err, ErrWalletNotFound)
		assert.NoError(t, dbMock.ExpectationsWereMet())
		assert.NoError(t, redisMock.ExpectationsWereMet())
	})

	t.Run("db failure", func(t *testing.T) {
		redisMock.ExpectDel(redisKey).SetVal(1)
		dbMock.ExpectExec(`UPDATE "wallets" SET "balance"=\$1 WHERE id = \$2`).
			WithArgs(input.Amount, input.WalletID).
			WillReturnError(errors.New("update failed"))

		err := dao.UpdateBalance(input)
		assert.Error(t, err)
	})
}

func TestGetBalance(t *testing.T) {
	dao, dbMock, redisMock := setupTest(t)
	walletID := "wallet-456"
	redisKey := fmt.Sprintf("wallet_balance:%s", walletID)

	t.Run("cache hit", func(t *testing.T) {
		redisMock.ExpectGet(redisKey).SetVal("123.4567")

		balance, err := dao.GetBalance(walletID)
		assert.NoError(t, err)
		assert.Equal(t, 123.4567, balance)
	})

	t.Run("cache miss, fetch from DB", func(t *testing.T) {
		redisMock.ExpectGet(redisKey).RedisNil()
		rows := sqlmock.NewRows([]string{"balance"}).AddRow(88.8888)
		dbMock.ExpectQuery(`SELECT "balance" FROM "wallets" WHERE id = \$1 ORDER BY "wallets"\."id" LIMIT \$2`).
			WithArgs(walletID, 1).WillReturnRows(rows)
		redisMock.ExpectSet(redisKey, "88.8888", 600*time.Second).SetVal("OK")

		balance, err := dao.GetBalance(walletID)
		assert.NoError(t, err)
		assert.Equal(t, 88.8888, balance)
	})

	t.Run("wallet not found", func(t *testing.T) {
		redisMock.ExpectGet(redisKey).RedisNil()
		dbMock.ExpectQuery(`SELECT "balance" FROM "wallets" WHERE id = \$1 ORDER BY "wallets"\."id" LIMIT \$2`).
			WithArgs(walletID, 1).WillReturnError(gorm.ErrRecordNotFound)

		balance, err := dao.GetBalance(walletID)
		assert.ErrorIs(t, err, ErrWalletNotFound)
		assert.Equal(t, 0.0, balance)
	})
}
