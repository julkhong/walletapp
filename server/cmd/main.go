package main

import (
    "log"
    "net/http"

    "github.com/julkhong/walletapp/server/internal/api"
    "github.com/julkhong/walletapp/server/internal/config"
)

func main() {
    cfg := config.LoadConfig()
    router := api.SetupRouter(cfg)

    log.Println("Starting server on :8080")
    http.ListenAndServe(":8080", router)
}
