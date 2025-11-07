package svc

import (
	"time"

	"github.com/google/uuid"
	"github.com/jibitesh/state-service/internal/repo"
	"github.com/jibitesh/state-service/pkg/v1/runtime"
	"github.com/mangudaigb/dhauli-base/logger"
	"go.opentelemetry.io/otel/trace"
)

type InteractionService interface {
	GetById(interactionId string) (*runtime.Interaction, error)
	Create(interaction *runtime.Interaction) (*runtime.Interaction, error)
	Update(interaction *runtime.Interaction) (*runtime.Interaction, error)
	DeleteById(interactionId string) error

	UpdatePlan(interactionId, planId string, plan *runtime.Plan) (*runtime.Interaction, error)
	UpdateWorkflow(interactionId, workflowId string, workflow *runtime.ExecutionFlow) (*runtime.Interaction, error)
	UpdateExecutionGraph(interactionId, workflowId, executionId string, graph *runtime.ExecutionGraph) (*runtime.Interaction, error)
}

type interactionService struct {
	log  *logger.Logger
	tr   trace.Tracer
	repo repo.InteractionRepo
}

func (is *interactionService) GetById(iid string) (*runtime.Interaction, error) {
	return is.repo.Get(iid)
}

func (is *interactionService) Create(interaction *runtime.Interaction) (*runtime.Interaction, error) {
	if interaction.ID == "" {
		interaction.ID = uuid.NewString()
	}
	interaction.CreatedAt = time.Now()
	if err := is.repo.Save(interaction); err != nil {
		is.log.Errorf("Error while saving interaction: %v", err)
		return nil, err
	}
	return interaction, nil
}

func (is *interactionService) Update(interaction *runtime.Interaction) (*runtime.Interaction, error) {
	if err := is.repo.Update(interaction); err != nil {
		is.log.Errorf("Error while updating interaction: %v", err)
		return nil, err
	}
	return interaction, nil
}

func (is *interactionService) DeleteById(iid string) error {
	return is.repo.Delete(iid)
}

func (is *interactionService) UpdatePlan(interactionId, planId string, plan *runtime.Plan) (*runtime.Interaction, error) {
	interaction, err := is.repo.Get(interactionId)
	if err != nil {
		is.log.Errorf("Error while getting interaction id:%s by error: %v", interactionId, err)
		return nil, err
	}
	if interaction.Plan.ID == planId {
		interaction.Plan = plan
	}
	if err = is.repo.Update(interaction); err != nil {
		is.log.Errorf("Error while updating interaction id:%s with plan: %s by error: %v", interactionId, plan.ID, err)
		return nil, err
	}
	return interaction, nil
}

func (is *interactionService) UpdateWorkflow(interactionId, workflowId string, workflow *runtime.ExecutionFlow) (*runtime.Interaction, error) {
	interaction, err := is.repo.Get(interactionId)
	if err != nil {
		is.log.Errorf("Error while getting interaction id:%s by error: %v", interactionId, err)
		return nil, err
	}
	if interaction.ExecutionFlow.ID == workflowId {
		interaction.ExecutionFlow = workflow
	}
	if err = is.repo.Update(interaction); err != nil {
		is.log.Errorf("Error while updating interaction id:%s with workflow: %s by error: %v", interactionId, workflow.ID, err)
		return nil, err
	}
	return interaction, nil
}

func (is *interactionService) UpdateExecutionGraph(interactionId, workflowId, executionId string, graph *runtime.ExecutionGraph) (*runtime.Interaction, error) {
	interaction, err := is.repo.Get(interactionId)
	if err != nil {
		is.log.Errorf("Error while getting interaction id:%s by error: %v", interactionId, err)
		return nil, err
	}
	if interaction.ExecutionFlow.ID == workflowId && interaction.ExecutionFlow.ExecutionGraph.ID == executionId {
		interaction.ExecutionFlow.ExecutionGraph = graph
	}
	if err = is.repo.Update(interaction); err != nil {
		is.log.Errorf("Error while updating interaction id:%s with execution graph: %s by error: %v", interactionId, graph.ID, err)
		return nil, err
	}
	return interaction, nil
}

func NewInteractionService(log *logger.Logger, tr trace.Tracer, repo repo.InteractionRepo) InteractionService {
	return &interactionService{
		log:  log,
		tr:   tr,
		repo: repo,
	}
}
