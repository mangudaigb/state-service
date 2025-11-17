package pkg

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mangudaigb/dhauli-base/config"
	"github.com/mangudaigb/dhauli-base/logger"
	"github.com/mangudaigb/state-service/internal"
	"github.com/mangudaigb/state-service/internal/handler"
	"github.com/mangudaigb/state-service/internal/repo"
	"github.com/mangudaigb/state-service/internal/svc"
	"go.opentelemetry.io/otel/trace"
	"golang.org/x/net/context"
)

type StateServer struct {
	log *logger.Logger
	cfg *config.Config
	tr  trace.Tracer
}

func (ss *StateServer) Start() {
	iRepo, err := repo.NewInteractionRepo(context.Background(), ss.cfg, ss.log, ss.tr)
	if err != nil {
		ss.log.Fatalf("Error while creating interaction repo: %v", err)
	}
	mRepo, err := repo.NewMcpRepo(context.Background(), ss.cfg, ss.log, ss.tr)
	if err != nil {
		ss.log.Fatalf("Error while creating mcp repo: %v", err)
	}
	sRepo, err := repo.NewStepRepo(context.Background(), ss.cfg, ss.log, ss.tr)
	if err != nil {
		ss.log.Fatalf("Error while creating step repo: %v", err)
	}

	iSvc := svc.NewInteractionService(ss.log, ss.tr, iRepo)
	mSvc := svc.NewMcpService(ss.log, ss.tr, mRepo)
	sSvc := svc.NewStepService(ss.log, ss.tr, sRepo)

	ih := handler.NewInteractionHandler(ss.log, ss.tr, iSvc)
	mh := handler.NewMcpHandler(ss.log, ss.tr, mSvc)
	sh := handler.NewStepHandler(ss.log, ss.tr, sSvc)

	gh := gin.Default()
	internal.SetupRouter(gh, ih, sh, mh)

	serverAddr := fmt.Sprintf(":%d", ss.cfg.Server.Port)

	server := &http.Server{
		Addr:         serverAddr,
		Handler:      gh,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	go func() {
		ss.log.Infof("Starting server on %s", serverAddr)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			ss.log.Fatalf("Error while starting server: %v", err)
		}
	}()
	ss.log.Infof("Server listening on %s", serverAddr)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, os.Kill)
	<-quit
	ss.log.Info("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		ss.log.Fatalf("Server forced to shutdown (timeout/error): %v", err)
	}
	ss.log.Info("Server successfully exited.")
}

func NewStateServer(cfg *config.Config, log *logger.Logger, tr trace.Tracer) *StateServer {
	return &StateServer{
		log: log,
		cfg: cfg,
		tr:  tr,
	}
}
