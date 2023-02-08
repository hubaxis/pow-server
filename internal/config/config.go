package config

import (
    "github.com/caarlos0/env/v6"
)

type Config struct {
    Port              int    `env:"PORT" envDefault:"44444"`
    ChuckNorrisEndpoint string `env:"QUOTE_ENDPOINT" envDefault:"https://api.chucknorris.io/jokes/random"`
}

func New() (*Config, error) {
    cfg := new(Config)
    err := env.Parse(cfg)
    if err != nil {
        return nil, err
    }
    return cfg, nil
}