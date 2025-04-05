package repository

type WalletRepository struct {
    DB string
}

func NewWalletRepository(db string) *WalletRepository {
    return &WalletRepository{DB: db}
}
