package svc

import (
	"github.com/google/uuid"
	"github.com/jibitesh/state-service/internal/repo"
	"github.com/jibitesh/state-service/pkg/v1/runtime"
	"github.com/mangudaigb/dhauli-base/logger"
	"go.opentelemetry.io/otel/trace"
)

type McpService interface {
	GetByInteractionIdAndWorkflowIdAndId(interactionId, workflowId, mcpId string) (*runtime.MCP, error)
	CreateByInteractionIdAndWorkflowId(interactionId, workflowId string, mcp *runtime.MCP) (*runtime.MCP, error)
	UpdateByInteractionIdAndWorkflowId(interactionId, workflowId string, mcp *runtime.MCP) (*runtime.MCP, error)
	AddTool(interactionId, workflowId, mcpId string, tool *runtime.Tool) (*runtime.MCP, error)
	DeleteByInteractionIdAndWorkflowIdAndId(interactionId, workflowId, mcpId string) error
}

type mcpService struct {
	log             *logger.Logger
	tr              trace.Tracer
	mcpRepo         repo.MCPRepo
	interactionRepo repo.InteractionRepo
}

func (ms *mcpService) GetByInteractionIdAndWorkflowIdAndId(interactionId, workflowId, mcpId string) (*runtime.MCP, error) {
	return ms.mcpRepo.Get(interactionId, workflowId, mcpId)
}

// CreateByInteractionIdAndWorkflowId Saves the mcp and updates the reference in workflow
func (ms *mcpService) CreateByInteractionIdAndWorkflowId(interactionId, workflowId string, mcp *runtime.MCP) (*runtime.MCP, error) {
	if mcp.ID == "" {
		mcp.ID = uuid.NewString()
	}
	err := ms.mcpRepo.Save(interactionId, workflowId, mcp)
	if err != nil {
		ms.log.Errorf("Error while saving mcp: %v", err)
		return nil, err
	}
	interaction, err := ms.interactionRepo.Get(interactionId)
	if err != nil {
		ms.log.Errorf("Error while getting interaction id:%s to update MCP by error: %v", interactionId, err)
		return nil, err
	}
	if interaction.ExecutionFlow.ID == workflowId {
		interaction.ExecutionFlow.AvailableMcpRefs = append(interaction.ExecutionFlow.AvailableMcpRefs, mcp.ID)
		err = ms.interactionRepo.Update(interaction)
		if err != nil {
			ms.log.Errorf("Error while updating interaction id:%s with mcpId: %s by error: %v", interactionId, mcp.ID, err)
			return nil, err
		}
	}
	return mcp, nil
}

func (ms *mcpService) UpdateByInteractionIdAndWorkflowId(interactionId, workflowId string, mcp *runtime.MCP) (*runtime.MCP, error) {
	err := ms.mcpRepo.Update(interactionId, workflowId, mcp)
	if err != nil {
		ms.log.Errorf("Error while updating mcp: %v", err)
		return nil, err
	}
	return mcp, nil
}

func (ms *mcpService) AddTool(interactionId, workflowId, mcpId string, tool *runtime.Tool) (*runtime.MCP, error) {
	mcp, err := ms.mcpRepo.Get(interactionId, workflowId, mcpId)
	if err != nil {
		ms.log.Errorf("Error while getting interactionId: %s, workflowId: %s, mcpId: %s by id: %v", interactionId, workflowId, mcpId, err)
		return nil, err
	}
	mcp.Tools = append(mcp.Tools, *tool)
	if err = ms.mcpRepo.Update(interactionId, workflowId, mcp); err != nil {
		ms.log.Errorf("Error while updating mcp: %v", err)
		return nil, err
	}
	return mcp, nil
}

func (ms *mcpService) DeleteByInteractionIdAndWorkflowIdAndId(interactionId, workflowId, mcpId string) error {
	return ms.mcpRepo.Delete(interactionId, workflowId, mcpId)
}

func NewMcpService(log *logger.Logger, tr trace.Tracer, repo repo.MCPRepo) McpService {
	return &mcpService{
		log:     log,
		tr:      tr,
		mcpRepo: repo,
	}
}
