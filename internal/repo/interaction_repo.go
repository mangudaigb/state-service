package repo

import (
	"github.com/jibitesh/state-service/pkg/v1/runtime"
	"github.com/mangudaigb/dhauli-base/config"
	"github.com/mangudaigb/dhauli-base/db"
	"github.com/mangudaigb/dhauli-base/logger"
	"go.opentelemetry.io/otel/trace"
)

type InteractionRepo interface {
	Get(iid string) (*runtime.Interaction, error)
	Save(interaction *runtime.Interaction) error
	Update(interaction *runtime.Interaction) error
	Delete(iid string) error
	Close()
}

type RedisInteractionRepo struct {
	cfg   *config.Config
	log   *logger.Logger
	tr    trace.Tracer
	store db.RedisStore[runtime.Interaction]
}

func (ir *RedisInteractionRepo) Get(iid string) (*runtime.Interaction, error) {
	return ir.store.Get(iid)
}

func (ir *RedisInteractionRepo) Save(interaction *runtime.Interaction) error {
	return ir.store.Set("interaction:"+interaction.ID, interaction)
}

func (ir *RedisInteractionRepo) Update(interaction *runtime.Interaction) error {
	return ir.store.Set("interaction:"+interaction.ID, interaction)
}

func (ir *RedisInteractionRepo) Delete(iid string) error {
	return ir.store.Delete("interaction:" + iid)
}

func (ir *RedisInteractionRepo) Close() {
	err := ir.store.Close()
	if err != nil {
		ir.log.Errorf("Error while closing redis store: %v", err)
	}
}

func NewInteractionRepo(cfg *config.Config, log *logger.Logger, tr trace.Tracer) (InteractionRepo, error) {
	interactionStore, err := db.NewRedisStore[runtime.Interaction](cfg, log)
	if err != nil {
		log.Errorf("Error while creating interaction redis store: %v", err)
		return nil, err
	}

	return &RedisInteractionRepo{
		cfg:   cfg,
		log:   log,
		tr:    tr,
		store: interactionStore,
	}, nil
}
