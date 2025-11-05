package repo

import (
	"github.com/jibitesh/state-service/pkg/v1/runtime"
	"github.com/mangudaigb/dhauli-base/config"
	"github.com/mangudaigb/dhauli-base/db"
	"github.com/mangudaigb/dhauli-base/logger"
	"go.opentelemetry.io/otel/trace"
)

type MCPRepo interface {
	Get(interactionId, workflowId string, mcpId string) (*runtime.MCP, error)
	Save(interactionId, workflowId string, mcp *runtime.MCP) error
	Update(interactionId, workflowId string, mcp *runtime.MCP) error
	Delete(interactionId, workflowId, mcpId string) error
	Close()
}

type RedisMCPRepo struct {
	cfg   *config.Config
	log   *logger.Logger
	tr    trace.Tracer
	store db.RedisStore[runtime.MCP]
}

func (mr *RedisMCPRepo) Get(interactionId string, workflowId string, mcpId string) (*runtime.MCP, error) {
	return mr.store.Get("interaction:" + interactionId + ":workflow:" + workflowId + ":mcp:" + mcpId)
}

func (mr *RedisMCPRepo) Save(interactionId, workflowId string, mcp *runtime.MCP) error {
	return mr.store.Set("interaction:"+interactionId+":workflow:"+workflowId+":mcp:"+mcp.ID, mcp)
}

func (mr *RedisMCPRepo) Update(interactionId, workflowId string, mcp *runtime.MCP) error {
	return mr.store.Set("interaction:"+interactionId+":workflow:"+workflowId+":mcp:"+mcp.ID, mcp)
}

func (mr *RedisMCPRepo) Delete(interactionId, workflowId string, mcpId string) error {
	return mr.store.Delete("interaction:" + interactionId + ":workflow:" + workflowId + ":mcp:" + mcpId)
}

func (mr *RedisMCPRepo) Close() {
	err := mr.store.Close()
	if err != nil {
		mr.log.Errorf("Error while closing redis store: %v", err)
	}
}

func NewMcpRepo(cfg *config.Config, log *logger.Logger, tr trace.Tracer) (MCPRepo, error) {
	mcpStore, err := db.NewRedisStore[runtime.MCP](cfg, log)
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
