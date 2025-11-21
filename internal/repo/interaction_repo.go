package repo

import (
	"context"

	"github.com/mangudaigb/dhauli-base/config"
	"github.com/mangudaigb/dhauli-base/db"
	"github.com/mangudaigb/dhauli-base/logger"
	"github.com/mangudaigb/dhauli-base/types/runtime"
	"go.opentelemetry.io/otel/trace"
)

type InteractionRepo interface {
	Get(ctx context.Context, iid string) (*runtime.Interaction, error)
	Save(ctx context.Context, interaction *runtime.Interaction) error
	Update(ctx context.Context, interaction *runtime.Interaction) error
	Delete(ctx context.Context, iid string) error
	Close()
}

type RedisInteractionRepo struct {
	cfg   *config.Config
	log   *logger.Logger
	tr    trace.Tracer
	store db.RedisStore[runtime.Interaction]
}

func (ir *RedisInteractionRepo) Get(ctx context.Context, iid string) (*runtime.Interaction, error) {
	return ir.store.Get(ctx, "interaction:"+iid)
}

func (ir *RedisInteractionRepo) Save(ctx context.Context, interaction *runtime.Interaction) error {
	return ir.store.Set(ctx, "interaction:"+interaction.ID, interaction)
}

func (ir *RedisInteractionRepo) Update(ctx context.Context, interaction *runtime.Interaction) error {
	return ir.store.Set(ctx, "interaction:"+interaction.ID, interaction)
}

func (ir *RedisInteractionRepo) Delete(ctx context.Context, iid string) error {
	return ir.store.Delete(ctx, "interaction:"+iid)
}

func (ir *RedisInteractionRepo) Close() {
	err := ir.store.Close()
	if err != nil {
		ir.log.Errorf("Error while closing redis store: %v", err)
	}
}

func NewInteractionRepo(ctx context.Context, cfg *config.Config, log *logger.Logger, tr trace.Tracer) (InteractionRepo, error) {
	interactionStore, err := db.NewRedisStore[runtime.Interaction](ctx, cfg, log)
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
