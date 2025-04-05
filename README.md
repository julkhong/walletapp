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
![alt text](static/wallet_service_architecture.png)

## Notes


