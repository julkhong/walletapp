package config

type Config struct {
    DB string
}

func LoadConfig() *Config {
    return &Config{
        DB: "postgres://user:pass@localhost:5432/wallet",
    }
}
