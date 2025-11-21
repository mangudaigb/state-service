package svc

import (
	"context"

	"github.com/google/uuid"
	"github.com/mangudaigb/dhauli-base/logger"
	"github.com/mangudaigb/dhauli-base/types/runtime"
	"github.com/mangudaigb/state-service/internal/repo"
	"go.opentelemetry.io/otel/trace"
)

type McpService interface {
	GetByInteractionIdAndWorkflowIdAndId(ctx context.Context, interactionId, workflowId, mcpId string) (*runtime.MCP, error)
	CreateByInteractionIdAndWorkflowId(ctx context.Context, interactionId, workflowId string, mcp *runtime.MCP) (*runtime.MCP, error)
	UpdateByInteractionIdAndWorkflowId(ctx context.Context, interactionId, workflowId string, mcp *runtime.MCP) (*runtime.MCP, error)
	AddTool(ctx context.Context, interactionId, workflowId, mcpId string, tool *runtime.Tool) (*runtime.MCP, error)
	DeleteByInteractionIdAndWorkflowIdAndId(ctx context.Context, interactionId, workflowId, mcpId string) error
}

type mcpService struct {
	log             *logger.Logger
	tr              trace.Tracer
	mcpRepo         repo.MCPRepo
	interactionRepo repo.InteractionRepo
}

func (ms *mcpService) GetByInteractionIdAndWorkflowIdAndId(ctx context.Context, interactionId, workflowId, mcpId string) (*runtime.MCP, error) {
	return ms.mcpRepo.Get(ctx, interactionId, workflowId, mcpId)
}

// CreateByInteractionIdAndWorkflowId Saves the mcp and updates the reference in workflow
func (ms *mcpService) CreateByInteractionIdAndWorkflowId(ctx context.Context, interactionId, workflowId string, mcp *runtime.MCP) (*runtime.MCP, error) {
	if mcp.ID == "" {
		mcp.ID = uuid.NewString()
	}
	err := ms.mcpRepo.Save(ctx, interactionId, workflowId, mcp)
	if err != nil {
		ms.log.Errorf("Error while saving mcp: %v", err)
		return nil, err
	}
	interaction, err := ms.interactionRepo.Get(ctx, interactionId)
	if err != nil {
		ms.log.Errorf("Error while getting interaction id:%s to update MCP by error: %v", interactionId, err)
		return nil, err
	}
	if interaction.ExecutionFlow.ID == workflowId {
		interaction.ExecutionFlow.AvailableMcpRefs = append(interaction.ExecutionFlow.AvailableMcpRefs, mcp.ID)
		err = ms.interactionRepo.Update(ctx, interaction)
		if err != nil {
			ms.log.Errorf("Error while updating interaction id:%s with mcpId: %s by error: %v", interactionId, mcp.ID, err)
			return nil, err
		}
	}
	return mcp, nil
}

func (ms *mcpService) UpdateByInteractionIdAndWorkflowId(ctx context.Context, interactionId, workflowId string, mcp *runtime.MCP) (*runtime.MCP, error) {
	err := ms.mcpRepo.Update(ctx, interactionId, workflowId, mcp)
	if err != nil {
		ms.log.Errorf("Error while updating mcp: %v", err)
		return nil, err
	}
	return mcp, nil
}

func (ms *mcpService) AddTool(ctx context.Context, interactionId, workflowId, mcpId string, tool *runtime.Tool) (*runtime.MCP, error) {
	mcp, err := ms.mcpRepo.Get(ctx, interactionId, workflowId, mcpId)
	if err != nil {
		ms.log.Errorf("Error while getting interactionId: %s, workflowId: %s, mcpId: %s by id: %v", interactionId, workflowId, mcpId, err)
		return nil, err
	}
	mcp.Tools = append(mcp.Tools, *tool)
	if err = ms.mcpRepo.Update(ctx, interactionId, workflowId, mcp); err != nil {
		ms.log.Errorf("Error while updating mcp: %v", err)
		return nil, err
	}
	return mcp, nil
}

func (ms *mcpService) DeleteByInteractionIdAndWorkflowIdAndId(ctx context.Context, interactionId, workflowId, mcpId string) error {
	return ms.mcpRepo.Delete(ctx, interactionId, workflowId, mcpId)
}

func NewMcpService(log *logger.Logger, tr trace.Tracer, repo repo.MCPRepo) McpService {
	return &mcpService{
		log:     log,
		tr:      tr,
		mcpRepo: repo,
	}
}
