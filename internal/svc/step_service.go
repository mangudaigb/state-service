package svc

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/mangudaigb/dhauli-base/logger"
	"github.com/mangudaigb/dhauli-base/types/runtime"
	"github.com/mangudaigb/state-service/internal/repo"
	"go.opentelemetry.io/otel/trace"
)

type StepService interface {
	GetByInteractionIdAndExecutionIdAndId(ctx context.Context, interactionId, workflowId, executionId string, stepId string) (*runtime.Step, error)
	CreateByInteractionIdAndExecutionId(ctx context.Context, interactionId, workflowId, executionId string, step *runtime.Step) (*runtime.Step, error)
	UpdateByInteractionIdAndExecutionId(ctx context.Context, interactionId, workflowId, executionId string, step *runtime.Step) (*runtime.Step, error)
	UpdateStatusByInteractionIdAndExecutionIdAndId(ctx context.Context, interactionId, workflowId, executionId, stepId string, status runtime.Status) (*runtime.Step, error)
	DeleteByInteractionIdAndExecutionIdAndId(ctx context.Context, interactionId, workflowId, executionId, stepId string) error
}

type stepService struct {
	log             *logger.Logger
	tr              trace.Tracer
	stepRepo        repo.StepRepo
	interactionRepo repo.InteractionRepo
}

func (ss *stepService) GetByInteractionIdAndExecutionIdAndId(ctx context.Context, interactionId, workflowId, executionId, stepId string) (*runtime.Step, error) {
	return ss.stepRepo.Get(ctx, interactionId, workflowId, executionId, stepId)
}

// CreateByInteractionIdAndExecutionId Saves the step and updates the reference in execution graph
func (ss *stepService) CreateByInteractionIdAndExecutionId(ctx context.Context, interactionId, workflowId, executionId string, step *runtime.Step) (*runtime.Step, error) {
	if step.ID == "" {
		step.ID = uuid.NewString()
	}
	step.Status = runtime.StatusPending
	interaction, err := ss.interactionRepo.Get(ctx, interactionId)
	if err != nil {
		ss.log.Errorf("Error while getting interaction id: %s to update Step: %s by error: %v", interactionId, step.ID, err)
		return nil, err
	}
	if interaction.ExecutionFlow.ID != workflowId || interaction.ExecutionFlow.ExecutionGraph.ID != executionId {
		ss.log.Errorf("Mismatch in ids, i.wf.Id/wfid:%s/%s, i.wf.eg.Id/egId:%s/%s", workflowId, interaction.ExecutionFlow.ID, executionId, interaction.ExecutionFlow.ExecutionGraph.ID)
		return nil, errors.New("mismatch in workflow or execution graph or both ids")
	}
	err = ss.stepRepo.Save(ctx, interactionId, workflowId, executionId, step)
	if err != nil {
		ss.log.Errorf("Error while creating step: %v for interaction id: %s", err, step.ID)
		return nil, err
	}
	en := runtime.ExecutionNode{
		StepId: step.ID,
		Name:   step.Name,
		Status: step.Status,
	}
	interaction.ExecutionFlow.ExecutionGraph.Nodes = append(interaction.ExecutionFlow.ExecutionGraph.Nodes, en)
	err = ss.interactionRepo.Update(ctx, interaction)
	if err != nil {
		ss.log.Errorf("Error while updating interaction id: %s with step: %s by error: %v", interactionId, step.ID, err)
		return nil, err
	}
	return step, nil
}

func (ss *stepService) UpdateByInteractionIdAndExecutionId(ctx context.Context, interactionId, workflowId, executionId string, step *runtime.Step) (*runtime.Step, error) {
	err := ss.stepRepo.Update(ctx, interactionId, workflowId, executionId, step)
	if err != nil {
		ss.log.Errorf("Error while updating step: %v", err)
		return nil, err
	}
	return step, nil
}

// UpdateStatusByInteractionIdAndExecutionIdAndId Saves the step and updates the status in the step reference in execution graph
func (ss *stepService) UpdateStatusByInteractionIdAndExecutionIdAndId(ctx context.Context, interactionId, workflowId, executionId, stepId string, status runtime.Status) (*runtime.Step, error) {
	interaction, err := ss.interactionRepo.Get(ctx, interactionId)
	if err != nil {
		ss.log.Errorf("Error while getting interaction id: %s to update Step: %s by error: %v", interactionId, stepId, err)
		return nil, err
	}
	if interaction.ExecutionFlow.ID != workflowId || interaction.ExecutionFlow.ExecutionGraph.ID != executionId {
		ss.log.Errorf("Mismatch in ids, i.wf.Id/wfid:%s/%s, i.wf.eg.Id/egId:%s/%s", workflowId, interaction.ExecutionFlow.ID, executionId, interaction.ExecutionFlow.ExecutionGraph.ID)
		return nil, errors.New("mismatch in workflow or execution graph or both ids")
	}

	step, err := ss.stepRepo.Get(ctx, interactionId, workflowId, executionId, stepId)
	if err != nil {
		ss.log.Errorf("Error while getting step by id: %v", err)
		return nil, err
	}
	if status == runtime.StatusStop || status == runtime.StatusError || status == runtime.StatusSuccess {
		step.FinishedAt = time.Now()
	}
	step.Status = status
	err = ss.stepRepo.Update(ctx, interactionId, workflowId, executionId, step)
	if err != nil {
		ss.log.Errorf("Error while updating status of step: %v", err)
		return nil, err
	}
	for i, node := range interaction.ExecutionFlow.ExecutionGraph.Nodes {
		if node.StepId == step.ID {
			interaction.ExecutionFlow.ExecutionGraph.Nodes[i].Status = step.Status
			break
		}
	}
	err = ss.interactionRepo.Update(ctx, interaction)
	if err != nil {
		ss.log.Errorf("Error while updating interaction id: %s with step: %s by error: %v", interactionId, step.ID, err)
		return nil, err
	}
	return step, nil
}

func (ss *stepService) DeleteByInteractionIdAndExecutionIdAndId(ctx context.Context, interactionId, workflowId, executionId, stepId string) error {
	return ss.stepRepo.Delete(ctx, interactionId, workflowId, executionId, stepId)
}

func NewStepService(log *logger.Logger, tr trace.Tracer, repo repo.StepRepo) StepService {
	return &stepService{
		log:      log,
		tr:       tr,
		stepRepo: repo,
	}
}
