package main

import (
	"context"
	"fmt"

	"github.com/mangudaigb/dhauli-base/config"
	"github.com/mangudaigb/dhauli-base/discover"
	"github.com/mangudaigb/dhauli-base/logger"
	"github.com/mangudaigb/dhauli-base/tracing"
	"github.com/mangudaigb/state-service/pkg"
)

func main() {
	cfg, err := config.GetConfig()
	if err != nil {
		fmt.Println("Error reading the config file", err)
		panic(err)
	}

	log, err := logger.NewLogger(cfg)
	if err != nil {
		fmt.Println("Error creating logger", err)
		panic(err)
	}

	tp := tracing.InitTracerProvider(cfg, log)
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Errorf("Error shutting down tracer provider: %v", err)
		}
	}()
	tr := tp.Tracer("context-service")

	registry := discover.NewRegistryInfo(cfg, log)
	registry.Register(discover.SERVICE)

	//StartConsumer(context.Background(), cfg, tr, log)

	server := pkg.NewStateServer(cfg, log, tr)
	server.Start()

}
