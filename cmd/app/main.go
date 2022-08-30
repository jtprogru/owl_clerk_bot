package main

import (
	"github.com/jtprogru/owl_clerk_bot/internal/conf"
	"github.com/jtprogru/owl_clerk_bot/internal/service"
	"github.com/jtprogru/owl_clerk_bot/internal/transport/tg"
)

func main() {
	conf := conf.New()
	sm := service.MockSM{}
	client := tg.NewTG(sm, conf)

	client.Run()
	// logger.Info("bot stopped...")
}
