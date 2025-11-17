package repo

import (
	"context"

	"github.com/mangudaigb/dhauli-base/config"
	"github.com/mangudaigb/dhauli-base/db"
	"github.com/mangudaigb/dhauli-base/logger"
	"github.com/mangudaigb/state-service/pkg/v1/runtime"
	"go.opentelemetry.io/otel/trace"
)

type StepRepo interface {
	Get(ctx context.Context, interactionId, workflowId, executionId, stepId string) (*runtime.Step, error)
	Save(ctx context.Context, interactionId, workflowId, executionId string, step *runtime.Step) error
	Update(ctx context.Context, interactionId, workflowId, executionId string, step *runtime.Step) error
	Delete(ctx context.Context, interactionId, workflowId, executionId string, stepId string) error
	Close()
}

type RedisStepRepo struct {
	cfg   *config.Config
	log   *logger.Logger
	tr    trace.Tracer
	store db.RedisStore[runtime.Step]
}

func (sr *RedisStepRepo) Get(ctx context.Context, interactionId, workflowId, executionId, stepId string) (*runtime.Step, error) {
	return sr.store.Get(ctx, "interaction:"+interactionId+":workflow:"+workflowId+":execution:"+executionId+":step:"+stepId)
}

func (sr *RedisStepRepo) Save(ctx context.Context, interactionId, workflowId, executionId string, step *runtime.Step) error {
	return sr.store.Set(ctx, "interaction:"+interactionId+":workflow:"+workflowId+":execution:"+executionId+":step:"+step.ID, step)
}

func (sr *RedisStepRepo) Update(ctx context.Context, interactionId, workflowId, executionId string, step *runtime.Step) error {
	return sr.store.Set(ctx, "interaction:"+interactionId+":workflow:"+workflowId+":execution:"+executionId+":step:"+step.ID, step)
}

func (sr *RedisStepRepo) Delete(ctx context.Context, interactionId, workflowId, executionId string, stepId string) error {
	return sr.store.Delete(ctx, "interaction:"+interactionId+":workflow:"+workflowId+":execution:"+executionId+":step:"+stepId)
}

func (sr *RedisStepRepo) Close() {
	err := sr.store.Close()
	if err != nil {
		sr.log.Errorf("Error while closing redis store: %v", err)
	}
}

func NewStepRepo(ctx context.Context, cfg *config.Config, log *logger.Logger, tr trace.Tracer) (StepRepo, error) {
	stepStore, err := db.NewRedisStore[runtime.Step](ctx, cfg, log)
	if err != nil {
		log.Errorf("Error while creating step redis store: %v", err)
		return nil, err
	}

	return &RedisStepRepo{
		cfg:   cfg,
		log:   log,
		tr:    tr,
		store: stepStore,
	}, nil
}
