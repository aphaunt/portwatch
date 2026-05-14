package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/example/portwatch/internal/alert"
	"github.com/example/portwatch/internal/config"
	"github.com/example/portwatch/internal/monitor"
	"github.com/example/portwatch/internal/scanner"
	"github.com/example/portwatch/internal/state"
)

func main() {
	configPath := flag.String("config", "portwatch.json", "path to config file")
	flag.Parse()

	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	store, err := state.NewStore(cfg.StateFile)
	if err != nil {
		log.Fatalf("failed to init state store: %v", err)
	}

	sc := scanner.NewScanner(cfg.Timeout)
	notifier := alert.NewLogNotifier(log.New(os.Stdout, "", log.LstdFlags))

	m := monitor.New(sc, store, notifier, cfg)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	log.Printf("portwatch starting, config: %s", *configPath)
	if err := m.Run(ctx); err != nil {
		log.Fatalf("monitor exited with error: %v", err)
	}
	log.Println("portwatch stopped")
}
