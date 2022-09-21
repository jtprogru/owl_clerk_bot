package main

import (
	"context"
	"fmt"
	"github.com/jtprogru/owl_clerk_bot/internal/entities/smentities"
	"github.com/jtprogru/owl_clerk_bot/internal/infrastructure/db/mem/smmemstore"
	"github.com/jtprogru/owl_clerk_bot/internal/service/ex"
	"github.com/jtprogru/owl_clerk_bot/internal/usecases/app/repo/smrepo"
	"github.com/sirupsen/logrus"
	"os"

	"github.com/jtprogru/owl_clerk_bot/internal/transport/tg"
)

//type MockAnswer struct {
//	msg string
//	kb  []string
//}
//
//func (ma MockAnswer) GetMessage() string {
//	return ma.msg
//}
//
//func (ma MockAnswer) GetKeyboard() []string {
//	return ma.kb
//}
//
type csm struct {
	sm *ex.SM
}

func (s csm) SaveOrUpdateState(ctx context.Context, p tg.IProfile, m tg.IMessage) (tg.Answer, error) {
	return s.sm.SaveOrUpdateState(ctx, p, m)
}

func main() {
	logger := logrus.New()
	logger.Formatter = &logrus.TextFormatter{
		DisableColors:   true,
		TimestampFormat: "2006-01-02T15:04:05",
		FullTimestamp:   true,
	}
	logger.Out = os.Stdout
	logger.SetReportCaller(true)
	// создаем стор
	smSt := smmemstore.NewMemStore()
	// передаем стор в репозиторий
	smRp := smrepo.NewSmRepo(smSt)

	// Для себя сделал небольшую проверку, было лень писать тест
	State0 := smentities.SM{
		Id:      0,
		NextIds: []int{1, 6, 2, 4},
		Answer:  "Test",
		Buttons: []string{"Button1", "Button6", "Button2", "Button4"},
		Handler: "State0",
	}

	State1 := smentities.SM{
		Id:      1,
		NextIds: []int{6, 2},
		Answer:  "Test",
		Buttons: []string{"Button6", "Button2"},
		Handler: "Func State1",
	}
	st0, err := smRp.Create(State0)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(smRp.Read(st0))

	st1, err := smRp.Create(State1)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(smRp.Read(st1))

	State0.Answer = "Update Test"
	smRp.Update(State0)

	fmt.Println(smRp.Read(st0))

	smRp.Delete(st0)

	fmt.Println(smRp.Read(st0))
	fmt.Println(smRp.Read(st1))

	os.Exit(1)

	botToken := os.Getenv("OWL_BOT_TOKEN")
	cfg := &tg.Config{
		BotToken: botToken,
		IsDebug:  false,
	}
	mocksm := csm{
		sm: ex.NewSM(),
	}
	client := tg.NewTG(mocksm, logger, cfg)

	client.Run()
	logger.Info("bot stopped...")
}
