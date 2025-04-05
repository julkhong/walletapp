package repository

import "time"

type User struct {
    ID        string
    Name      string
    Email     string
    CreatedAt time.Time
}

type Wallet struct {
    ID        string
    UserID    string
    Balance   int64
    CreatedAt time.Time
}

type Transaction struct {
    ID            string
    WalletID      string
    Type          string
    Amount        int64
    RelatedUserID *string
    CreatedAt     time.Time
}
