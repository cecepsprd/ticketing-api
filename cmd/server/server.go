package server

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/cecepsprd/ticketing-api/config"
	"github.com/cecepsprd/ticketing-api/handler"
	"github.com/cecepsprd/ticketing-api/repository"
	"github.com/cecepsprd/ticketing-api/service"
	"github.com/cecepsprd/ticketing-api/utils/logger"
	"github.com/cecepsprd/ticketing-api/utils/validate"
	"github.com/labstack/echo"

	en_translations "gopkg.in/go-playground/validator.v9/translations/en"
)

func RunServer() {
	var cfg = config.NewConfig()

	db, err := cfg.MysqlConnect()
	if err != nil {
		log.Fatal("error connecting to database: ", err.Error())
	}

	if err = logger.Init(cfg.App.LogLevel, cfg.App.LogTimeFormat); err != nil {
		log.Fatal(err)
	}

	e := echo.New()

	customValidator := validate.NewValidator()
	en_translations.RegisterDefaultTranslations(customValidator.Validator, customValidator.Translator)
	e.Validator = customValidator

	timeoutContext := time.Duration(cfg.App.ContextTimeout) * time.Second

	userRepository := repository.NewUserRepository(db)
	productRepository := repository.NewProductRepository(db)
	transactionRepository := repository.NewTransactionRepository(db)

	userService := service.NewUserService(userRepository, timeoutContext)
	authService := service.NewAuthService(userService, cfg.App.JWTSecret)
	productService := service.NewProductService(productRepository, timeoutContext)
	transactionService := service.NewTransactionService(transactionRepository, productRepository, timeoutContext)

	handler.NewAuthHandler(e, authService)
	handler.NewUserHandler(e, userService)
	handler.NewProductHandler(e, productService, transactionService)

	// Starting server
	go func() {
		err := e.Start(cfg.App.HTTPPort)
		if err != nil {
			log.Fatal("error starting server: ", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	signal.Notify(quit, syscall.SIGTERM)

	// Block until a signal is received.
	<-quit

	log.Println("server shutdown of 5 second.")

	// gracefully shutdown the server, waiting max 5 seconds for current operations to complete
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	e.Shutdown(ctx)
}
