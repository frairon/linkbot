package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/frairon/linkbot/internal"
	"github.com/frairon/linkbot/internal/bot"
	"github.com/frairon/linkbot/internal/link"
	"github.com/frairon/linkbot/internal/storage"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"golang.org/x/sync/errgroup"
)

func main() {
	cfg, err := internal.ReadConfig()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}

	botApi, err := tgbotapi.NewBotAPI(cfg.Token)
	if err != nil {
		log.Fatalf("error creating bot: %v", err)
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

		user, err := st.GetUser(cfg.InitialUserID)
		if err != nil {
			log.Fatalf("cannot find user that I just added: %v", err)
		}
		if user == nil {
			log.Fatalf("cannot load user that I just added")
		}
	}

	botApi.Debug = false
	b, err := bot.New(botApi, st, link.Home)
	if err != nil {
		log.Fatalf("error creating bot: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		defer cancel()
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
		<-c
		signal.Stop(c)
	}()

	// serve metrics from home
	errg, ctx := errgroup.WithContext(ctx)
	errg.Go(func() error {
		defer cancel()
		return b.Run(ctx)
	})

	adminServer := &http.Server{Addr: ":9090"}
	errg.Go(func() error {
		http.Handle("/metrics", promhttp.HandlerFor(prometheus.DefaultGatherer, promhttp.HandlerOpts{Registry: prometheus.DefaultRegisterer}))
		err := adminServer.ListenAndServe()
		if err == http.ErrServerClosed {
			return nil
		}
		return err
	})

	errg.Go(func() error {
		<-ctx.Done()

		return adminServer.Shutdown(context.Background())
	})

	if err := errg.Wait(); err != nil {
		log.Fatalf("error running bot: %v", err)
	}
}

func fatalOnError(msg string, err error) {
	if err != nil {
		log.Fatalf(msg, err)
	}
}
