package bootstrap

import (
	"context"
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
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

	return &Bootstrap{
		router:    router,
		log:       logger,
		db:        db,
		handlers:  []Handler{},
		validator: validator,
	}
}

func (b *Bootstrap) DepedencyInjection() {
	authenticationHandler := rest.NewAuthenticationHandler(
		b.log,
		service.NewAuthenticationService(
			b.log,
			repository.NewAuthenticationRepository(b.db),
		),
	)

	b.handlers = []Handler{
		authenticationHandler,
	}
}

func (b *Bootstrap) Run() {
	b.DepedencyInjection()
	b.Health()

	_persistence.Migrate(b.db)

	for _, handler := range b.handlers {
		handler.SetEndpoint(b.router)
	}

	if err := b.router.Listen(fmt.Sprintf(":%d", viper.GetInt("app.port"))); err != nil {
		b.log.Fatal("failed to run server", zap.Error(err))
	}
}

func (b *Bootstrap) Shutdown(ctx context.Context) {

}

func (b *Bootstrap) Health() {
	b.router.Get("/health", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON("OK")
	})
}
