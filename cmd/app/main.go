package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/jtprogru/owl_clerk_bot/internal/config"
	owlhttp "github.com/jtprogru/owl_clerk_bot/internal/http"
	"github.com/jtprogru/owl_clerk_bot/internal/service/intake"
	"github.com/jtprogru/owl_clerk_bot/internal/service/notify"
	"github.com/jtprogru/owl_clerk_bot/internal/service/outbound"
	"github.com/jtprogru/owl_clerk_bot/internal/service/sm"
	"github.com/jtprogru/owl_clerk_bot/internal/storage/sqlite"
	"github.com/jtprogru/owl_clerk_bot/internal/transport/tg"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	cfg, err := config.Load()
	if err != nil {
		logger.Error("config", "err", err)
		os.Exit(1)
	}
	if cfg.Debug {
		logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
		slog.SetDefault(logger)
	}

	db, err := sqlite.Open(cfg.DBPath)
	if err != nil {
		logger.Error("open db", "err", err)
		os.Exit(1)
	}
	defer func() { _ = db.Close() }()

	profiles := sqlite.NewProfileRepo(db)
	messages := sqlite.NewMessageRepo(db)
	states := sqlite.NewStateRepo(db)

	bot, err := tg.NewBot(cfg.BotToken, cfg.Debug, logger)
	if err != nil {
		logger.Error("create bot", "err", err)
		os.Exit(1)
	}
	client := tg.NewClient(bot, logger)

	notifier := notify.NewOwner(cfg.OwnerID, client)
	outboundSvc := outbound.New(client, messages)
	machine := sm.New()
	intakeSvc := intake.New(profiles, messages, states, machine, notifier, logger)

	handlers := tg.NewHandlers(cfg.OwnerID, logger, intakeSvc, outboundSvc, profiles, messages)
	handlers.Register(bot)

	web, err := owlhttp.New(cfg.ServeAddr(), cfg.WebUser, cfg.WebPass, logger, profiles, messages, outboundSvc)
	if err != nil {
		logger.Error("create web server", "err", err)
		os.Exit(1)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		logger.Info("bot starting", "owner_id", cfg.OwnerID)
		bot.Start()
	}()

	webErr := make(chan error, 1)
	go func() {
		webErr <- web.Run(ctx)
	}()

	<-ctx.Done()
	logger.Info("shutting down")
	bot.Stop()
	if err := <-webErr; err != nil {
		logger.Error("web shutdown", "err", err)
	}
}
