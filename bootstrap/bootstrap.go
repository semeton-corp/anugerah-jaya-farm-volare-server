package bootstrap

import (
	"context"
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	_email "github.com/semeton-corp/anugerah-jaya-farm-volare/infra/email"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/infra/env"
	_logger "github.com/semeton-corp/anugerah-jaya-farm-volare/infra/logger"
	_persistence "github.com/semeton-corp/anugerah-jaya-farm-volare/infra/persistence"
	_router "github.com/semeton-corp/anugerah-jaya-farm-volare/infra/router"
	_scheduler "github.com/semeton-corp/anugerah-jaya-farm-volare/infra/scheduler"
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
	email     *_email.Email
	scheduler *_scheduler.Scheduler
	validator *validator.Validate
}

type Handler interface {
	SetEndpoint(router *fiber.App)
}

func New() *Bootstrap {
	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		zap.L().Fatal(fmt.Sprintf("failed to load timezone: %v", err))
	}

	time.Local = loc // -> this is setting the global timezone

	env.Load()
	router := _router.New()
	logger := _logger.New()
	db := _persistence.New(logger)
	validator := _validator.New()
	email := _email.New()
	scheduler := _scheduler.New(db, logger)

	return &Bootstrap{
		router:    router,
		log:       logger,
		db:        db,
		handlers:  []Handler{},
		validator: validator,
		email:     email,
		scheduler: scheduler,
	}
}

func (b *Bootstrap) DepedencyInjection() {
	staffRepository := repository.NewStaffRepository(b.db)
	staffService := service.NewStaffService(b.log, staffRepository)

	authRepository := repository.NewAuthenticationRepository(b.db)
	authService := service.NewAuthenticationService(b.log, authRepository, b.email, staffService)
	authenticationHandler := rest.NewAuthenticationHandler(b.log, authService, b.validator)

	roleRepository := repository.NewRoleRepository(b.db)
	roleService := service.NewRoleService(b.log, roleRepository)
	roleHandler := rest.NewRoleHandler(b.log, roleService, b.validator)

	cageRepository := repository.NewCageRepository(b.db)
	cageService := service.NewCageService(b.log, cageRepository)
	cageHandler := rest.NewCageHandler(b.log, cageService, b.validator)

	chickenRepository := repository.NewChickenRepository(b.db)
	chickenService := service.NewChickenService(b.log, chickenRepository)
	chickenHandler := rest.NewChickenHandler(b.log, chickenService, b.validator)

	eggRepository := repository.NewEggRepository(b.db)
	eggService := service.NewEggService(b.log, eggRepository)
	eggHandler := rest.NewEggHandler(b.log, eggService, b.validator)

	storeRepository := repository.NewStoreRepository(b.db)
	storeService := service.NewStoreService(b.log, storeRepository)
	storeHandler := rest.NewStoreHandler(b.log, storeService, b.validator)

	warehouseRepository := repository.NewWarehouseRepository(b.db)
	warehouseService := service.NewWarehouseService(b.log, warehouseRepository, storeService)
	warehouseHandler := rest.NewWarehouseHandler(b.log, warehouseService, b.validator)

	workRepository := repository.NewWorkRepository(b.db)
	workService := service.NewWorkService(b.log, workRepository, roleService)
	workHandler := rest.NewWorkHandler(b.log, workService, b.validator)

	presenceRepository := repository.NewPresenceRepository(b.db)
	presenceService := service.NewPresenceService(b.log, presenceRepository)
	presenceHandler := rest.NewPresenceHandler(b.log, presenceService, b.validator)

	supplierRepository := repository.NewSupplierRepository(b.db)
	supplierService := service.NewSupplierService(b.log, supplierRepository)
	supplierHandler := rest.NewSupplierHandler(b.log, supplierService, b.validator)

	b.handlers = []Handler{
		authenticationHandler,
		roleHandler,
		cageHandler,
		chickenHandler,
		eggHandler,
		warehouseHandler,
		storeHandler,
		workHandler,
		presenceHandler,
		supplierHandler,
	}
}

func (b *Bootstrap) Run() {
	b.DepedencyInjection()
	b.Health()

	b.scheduler.InitScheduler()
	b.scheduler.Start()

	_persistence.Migrate(b.db)
	// _persistence.Rollback(b.db)

	b.router.Use(cors.New(cors.Config{
		AllowOrigins:  viper.GetString("server.cors.allow_origins"),
		AllowMethods:  viper.GetString("server.cors.allow_methods"),
		AllowHeaders:  viper.GetString("server.cors.allow_headers"),
		ExposeHeaders: viper.GetString("server.cors.expose_headers"),
		MaxAge:        viper.GetInt("server.cors.max_age"),
	}))

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

	b.scheduler.Stop()

	b.log.Info("server shutdown gracefully...")
}

func (b *Bootstrap) Health() {
	b.router.Get("/health", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON("OK")
	})
}
