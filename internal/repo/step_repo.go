package repo

import (
	"github.com/jibitesh/state-service/pkg/v1/runtime"
	"github.com/mangudaigb/dhauli-base/config"
	"github.com/mangudaigb/dhauli-base/db"
	"github.com/mangudaigb/dhauli-base/logger"
	"go.opentelemetry.io/otel/trace"
)

type StepRepo interface {
	Get(interactionId, workflowId, executionId, stepId string) (*runtime.Step, error)
	Save(interactionId, workflowId, executionId string, step *runtime.Step) error
	Update(interactionId, workflowId, executionId string, step *runtime.Step) error
	Delete(interactionId, workflowId, executionId string, stepId string) error
	Close()
}

type RedisStepRepo struct {
	cfg   *config.Config
	log   *logger.Logger
	tr    trace.Tracer
	store db.RedisStore[runtime.Step]
}

func (sr *RedisStepRepo) Get(interactionId, workflowId, executionId, stepId string) (*runtime.Step, error) {
	return sr.store.Get("interaction:" + interactionId + ":workflow:" + workflowId + ":execution:" + executionId + ":step:" + stepId)
}

func (sr *RedisStepRepo) Save(interactionId, workflowId, executionId string, step *runtime.Step) error {
	return sr.store.Set("interaction:"+interactionId+":workflow:"+workflowId+":execution:"+executionId+":step:"+step.ID, step)
}

func (sr *RedisStepRepo) Update(interactionId, workflowId, executionId string, step *runtime.Step) error {
	return sr.store.Set("interaction:"+interactionId+":workflow:"+workflowId+":execution:"+executionId+":step:"+step.ID, step)
}

func (sr *RedisStepRepo) Delete(interactionId, workflowId, executionId string, stepId string) error {
	return sr.store.Delete("interaction:" + interactionId + ":workflow:" + workflowId + ":execution:" + executionId + ":step:" + stepId)
}

func (sr *RedisStepRepo) Close() {
	err := sr.store.Close()
	if err != nil {
		sr.log.Errorf("Error while closing redis store: %v", err)
	}
}

func NewStepRepo(cfg *config.Config, log *logger.Logger, tr trace.Tracer) (StepRepo, error) {
	stepStore, err := db.NewRedisStore[runtime.Step](cfg, log)
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
