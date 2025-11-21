package repo

import (
	"context"

	"github.com/mangudaigb/dhauli-base/config"
	"github.com/mangudaigb/dhauli-base/db"
	"github.com/mangudaigb/dhauli-base/logger"
	"github.com/mangudaigb/dhauli-base/types/runtime"
	"go.opentelemetry.io/otel/trace"
)

type MCPRepo interface {
	Get(ctx context.Context, interactionId, workflowId string, mcpId string) (*runtime.MCP, error)
	Save(ctx context.Context, interactionId, workflowId string, mcp *runtime.MCP) error
	Update(ctx context.Context, interactionId, workflowId string, mcp *runtime.MCP) error
	Delete(ctx context.Context, interactionId, workflowId, mcpId string) error
	Close()
}

type RedisMCPRepo struct {
	cfg   *config.Config
	log   *logger.Logger
	tr    trace.Tracer
	store db.RedisStore[runtime.MCP]
}

func (mr *RedisMCPRepo) Get(ctx context.Context, interactionId string, workflowId string, mcpId string) (*runtime.MCP, error) {
	return mr.store.Get(ctx, "interaction:"+interactionId+":workflow:"+workflowId+":mcp:"+mcpId)
}

func (mr *RedisMCPRepo) Save(ctx context.Context, interactionId, workflowId string, mcp *runtime.MCP) error {
	return mr.store.Set(ctx, "interaction:"+interactionId+":workflow:"+workflowId+":mcp:"+mcp.ID, mcp)
}

func (mr *RedisMCPRepo) Update(ctx context.Context, interactionId, workflowId string, mcp *runtime.MCP) error {
	return mr.store.Set(ctx, "interaction:"+interactionId+":workflow:"+workflowId+":mcp:"+mcp.ID, mcp)
}

func (mr *RedisMCPRepo) Delete(ctx context.Context, interactionId, workflowId string, mcpId string) error {
	return mr.store.Delete(ctx, "interaction:"+interactionId+":workflow:"+workflowId+":mcp:"+mcpId)
}

func (mr *RedisMCPRepo) Close() {
	err := mr.store.Close()
	if err != nil {
		mr.log.Errorf("Error while closing redis store: %v", err)
	}
}

func NewMcpRepo(ctx context.Context, cfg *config.Config, log *logger.Logger, tr trace.Tracer) (MCPRepo, error) {
	mcpStore, err := db.NewRedisStore[runtime.MCP](ctx, cfg, log)
	if err != nil {
		log.Errorf("Error while creating mcp redis store: %v", err)
		return nil, err
	}

	return &RedisMCPRepo{
		cfg:   cfg,
		log:   log,
		tr:    tr,
		store: mcpStore,
	}, nil
}
