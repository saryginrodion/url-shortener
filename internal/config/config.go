package config

import "time"

type Config struct {
	Env              string `yaml:"env" env:"ENV" env-default:"local"`
	PostgresConfig   `yaml:"postgres"`
	HTTPServerConfig `yaml:"http_server"`
	RedisConfig      `yaml:"redis"`
	TokensConfig     `yaml:"tokens"`
}

type TokensConfig struct {
	AccessTTL  time.Duration `yaml:"access_ttl" env:"TOKENS_ACCESS_TTL_MINUTES" env-default:"15m"`
	RefreshTTL time.Duration `yaml:"refresh_ttl" env:"TOKENS_REFRESH_TTL_MINUTES" env-default:"43200m"`
	SecretKey  string        `yaml:"secret_key" env:"TOKENS_SECRET_KEY"`
}

type RedisConfig struct {
	Host     string `yaml:"host" env:"REDIS_HOST"`
	Password string `yaml:"password" env:"REDIS_PASSWORD"`
	Username string `yaml:"username" env:"REDIS_USERNAME"`
	DB       int    `yaml:"db" env:"REDIS_DB"`
}

type PostgresConfig struct {
	DSN string `yaml:"dsn" env:"POSTGRES_DSN"`
}

type HTTPServerConfig struct {
	Addr         string        `yaml:"addr" env:"HTTP_SERVER_ADDR" env-default:"localhost:8000"`
	ReadTimeout  time.Duration `yaml:"read_timeout" env:"HTTP_SERVER_READ_TIMEOUT" env-default:"10s"`
	WriteTimeout time.Duration `yaml:"write_timeout" env:"HTTP_SERVER_WRITE_TIMEOUT" env-default:"10s"`
	IdleTimeout  time.Duration `yaml:"idle_timeout" env:"HTTP_SERVER_IDLE_TIMEOUT" env-default:"10s"`
}

var cfg *Config = nil

func SetCfg(v *Config) {
	cfg = v
}

func Cfg() *Config {
	if cfg == nil {
		panic("Config is not loaded")
	}

	return cfg
}

func (c *Config) IsProd() bool {
	return c.Env == "prod"
}
