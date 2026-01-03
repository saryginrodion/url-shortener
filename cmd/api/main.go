package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/golang-jwt/jwt/v4"
	"github.com/ilyakaznacheev/cleanenv"
	"roadmap.restapi/internal/api"
	"roadmap.restapi/internal/config"
	"roadmap.restapi/internal/postgres"
	"roadmap.restapi/internal/redis"
	"roadmap.restapi/internal/token"
	"roadmap.restapi/internal/url"
	"roadmap.restapi/internal/user"
)

func main() {
	// Config
	cfgPath := os.Getenv("CONFIG_PATH")
	if cfgPath == "" {
		slog.Error("failed to load $CONFIG_PATH from environment")
		os.Exit(1)
	}

	var cfg config.Config
	if err := cleanenv.ReadConfig(cfgPath, &cfg); err != nil {
		slog.Error("error reading config file", "err", err)
		os.Exit(1)
	}

	config.SetCfg(&cfg)

	// Logging
	logLevel := slog.LevelInfo
	if !config.Cfg().IsProd() {
		logLevel = slog.LevelDebug
	}

	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: logLevel,
	})).
		With("env", cfg.Env)

	// Database
	pgDB, err := postgres.Connect(cfg.PostgresConfig.DSN)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to connect to postgres. err: %s", err.Error())
		os.Exit(1)
	}

	rdb, err := redis.Connect(
		cfg.RedisConfig.Host,
		cfg.RedisConfig.Username,
		cfg.RedisConfig.Password,
		cfg.RedisConfig.DB,
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to connect to redis. err: %s", err.Error())
		os.Exit(1)
	}

	// users
	userRepo := user.NewPostgresUserRepository(pgDB)
	passwordHasher := user.NewArgon2IDPasswordHasher()
	users := user.NewUseCases(userRepo, passwordHasher)

	// tokens
	tokenExtractor := token.NewJWTClaimsExtractor(
		[]byte(cfg.TokensConfig.SecretKey),
		jwt.SigningMethodHS512,
	)
	tokenGenerator := token.NewJWTGenerator(
		[]byte(cfg.TokensConfig.SecretKey),
		jwt.SigningMethodHS512,
		config.Cfg().TokensConfig.AccessTTL,
		config.Cfg().TokensConfig.RefreshTTL,
	)
	tokenWhitelist := token.NewRedisWhitelistRepository(
		rdb,
		time.Duration(config.Cfg().TokensConfig.RefreshTTL)*time.Minute,
	)
	tokens := token.NewUseCases(tokenExtractor, tokenGenerator, tokenWhitelist)

	// urls
	urlsRepo := url.NewPostgresURLRepository(pgDB)
	urls := url.NewUseCases(urlsRepo)

	// Router
	router := api.NewRouter(
		log,
		users,
		userRepo,
		tokens,
		tokenExtractor,
		urlsRepo,
		urls,
	)
	router.Mount("/debug", middleware.Profiler())

	// Server
	server := &http.Server{
		Addr:         cfg.HTTPServerConfig.Addr,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServerConfig.ReadTimeout,
		WriteTimeout: cfg.HTTPServerConfig.WriteTimeout,
		IdleTimeout:  cfg.HTTPServerConfig.IdleTimeout,
	}

	log.Info(fmt.Sprintf("Started at %s", server.Addr))
	server.ListenAndServe()
}
