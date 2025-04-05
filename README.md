# Wallet App


## TODO
1. db schema
2. architecture diagram
3. no secrets commit 
4. video demo
5. can setup redis 
6. include ut 

Include a readme
- explaining any decisions you made
- how to setup and run your code
- highlight how should reviewer view your code
- areas to be improved
- how long you spent on the test
- which features you chose not to do in the submission
- How does it satisfy functional and non-functional requirements
- Does it follow engineering best practices?

## Setup

1. Run Postgres
2. Set DB URL in `config/config.go`
3. Run migrations
4. Run the server:
```bash
go run cmd/main.go
```

## Features

- Deposit
- Withdraw
- Transfer
- Balance Check
- Transaction History

## Architecture 
Below is a simplified architecture diagram for wallet service.
<br></br>
![alt text](static/wallet_service_architecture.png)

## Service Directory
```
walletapp/
├── cmd/
│   └── main.go                    # App entry point: loads config, starts HTTP server
│
├── internal/
│   ├── api/                       # REST API layer (HTTP handlers/controllers)
│   │   └── handler.go             # Route setup and HTTP handlers
│
│   ├── config/                    # App configuration (env vars, DB URL)
│   │   └── config.go              # Loads and returns config struct
│
│   ├── repository/                # DB access layer (PostgreSQL + Redis if used)
│   │   ├── wallet_repo.go         # CRUD and query logic for wallet, user, transactions
│   │   ├── redis_repo.go          # (optional) Redis operations for caching
│   │   └── models.go              # Domain models: User, Wallet, Transaction
│
│   ├── service/                   # Business logic layer
│   │   └── wallet_service.go      # Implements deposit, withdraw, transfer, etc.
│
│   ├── middleware/                # (optional) Middlewares for logging, auth, recovery
│   │   └── logging.go             # Example logging middleware
│
│   └── utils/                     # (optional) Utility functions (e.g., ID generators)
│       └── helpers.go             # JSON response helpers, UUID, validation, etc.
│
├── migrations/                    # SQL migration files
│   └── 001_create_tables.sql      # SQL: CREATE TABLES for users, wallets, txns
│
├── test/                          # Unit and integration tests
│   ├── wallet_service_test.go     # Business logic test
│   └── api_test.go                # (optional) API integration test
│
├── .env                           # (optional) Environment variables
├── go.mod                         # Go module file
├── go.sum                         # Go dependencies checksum
└── README.md                      # Project overview, setup, assumptions, decisions

...
```
## Notes


