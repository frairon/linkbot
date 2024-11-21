package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/frairon/botty"
	"github.com/frairon/linkbot/internal"
	"github.com/frairon/linkbot/internal/link"
	"github.com/frairon/linkbot/internal/storage"
)

func main() {
	cfg, err := internal.ReadConfig()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}

	st, err := storage.NewStorage(fmt.Sprintf("file:%s", cfg.Database))
	if err != nil {
		log.Fatalf("error creating storage: %v", err)
	}

	if cfg.InitialUserID != 0 && cfg.InitialUserName != "" {
		log.Printf("adding initial user: %+v", cfg.InitialUserName)
		err := st.AddUser(cfg.InitialUserID, cfg.InitialUserName)
		if err != nil {
			log.Fatalf("error adding initial user: %v", err)
		}

		_, err = st.GetUser(cfg.InitialUserID)
		fatalOnError("cannot find user that I just added: %v", err)
	}
	sessionManager, userManager := link.NewManagers(st)
	botCfg := botty.NewConfig(cfg.Token, sessionManager, userManager, link.Home)
	botCfg.RootState = link.Home

	b, err := botty.New(botCfg)
	fatalOnError("erorr creating bot: %w", err)

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		defer cancel()
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
		<-c
		signal.Stop(c)
	}()

	if err := b.Run(ctx); err != nil {
		log.Fatalf("error running bot: %v", err)
	}
}

func fatalOnError(msg string, err error) {
	if err != nil {
		log.Fatalf(msg, err)
	}
}
