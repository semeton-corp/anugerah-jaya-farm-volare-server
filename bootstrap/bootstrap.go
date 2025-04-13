package bootstrap

import (
	"context"
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/infra/email"
	_email "github.com/semeton-corp/anugerah-jaya-farm-volare/infra/email"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/infra/env"
	_logger "github.com/semeton-corp/anugerah-jaya-farm-volare/infra/logger"
	_persistence "github.com/semeton-corp/anugerah-jaya-farm-volare/infra/persistence"
	_router "github.com/semeton-corp/anugerah-jaya-farm-volare/infra/router"
	_validator "github.com/semeton-corp/anugerah-jaya-farm-volare/infra/validator"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/handler/rest"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/repository"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/service"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Bootstrap struct {
	router    *fiber.App
	log       *zap.Logger
	db        *gorm.DB
	handlers  []Handler
	email     *email.Email
	validator *validator.Validate
}

type Handler interface {
	SetEndpoint(router *fiber.App)
}

func New() *Bootstrap {
	env.Load()
	router := _router.New()
	logger := _logger.New()
	db := _persistence.New(logger)
	validator := _validator.New()
	email := _email.New()

	return &Bootstrap{
		router:    router,
		log:       logger,
		db:        db,
		handlers:  []Handler{},
		validator: validator,
		email:     email,
	}
}

func (b *Bootstrap) DepedencyInjection() {
	authenticationHandler := rest.NewAuthenticationHandler(
		b.log,
		service.NewAuthenticationService(
			b.log,
			repository.NewAuthenticationRepository(b.db),
			b.email,
		),
		b.validator,
	)

	roleHandler := rest.NewRoleHandler(
		b.log,
		service.NewRoleService(
			b.log,
			repository.NewRoleRepository(b.db),
		),
		b.validator,
	)

	cageHandler := rest.NewCageHandler(
		b.log,
		service.NewCageService(
			b.log,
			repository.NewCageRepository(b.db),
		),
		b.validator,
	)

	b.handlers = []Handler{
		authenticationHandler,
		roleHandler,
		cageHandler,
	}
}

func (b *Bootstrap) Run() {
	b.DepedencyInjection()
	b.Health()

	_persistence.Migrate(b.db)

	for _, handler := range b.handlers {
		handler.SetEndpoint(b.router)
	}

	addr := fmt.Sprintf("%s:%d", viper.GetString("app.address"), viper.GetInt("app.port"))

	if err := b.router.Listen(addr); err != nil {
		b.log.Fatal("failed to run server", zap.Error(err))
	}
}

func (b *Bootstrap) Shutdown(ctx context.Context) {
	if err := b.router.Shutdown(); err != nil {
		b.log.Error("failed to shutdown server", zap.Error(err))
	}

	db, err := b.db.DB()
	if err != nil {
		b.log.Error("failed to get database connection", zap.Error(err))
	}

	if err := db.Close(); err != nil {
		b.log.Error("failed to close database connection", zap.Error(err))
	}

	if err := b.log.Sync(); err != nil {
		b.log.Error("failed to sync logger", zap.Error(err))
	}

	b.log.Info("server shutdown gracefully")
}

func (b *Bootstrap) Health() {
	b.router.Get("/health", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON("OK")
	})
}
