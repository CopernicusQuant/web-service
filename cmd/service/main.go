package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/CopernicusQuant/web-service/internal/handlers"
	"github.com/CopernicusQuant/web-service/internal/repositories"
	"github.com/CopernicusQuant/web-service/internal/transport"
)

const (
	PORT_NUMBER = 8800
	GO_ENV      = "GO_ENV"
)

const (
	runtimeDev  = "dev"
	runtimeProd = "prod"
)

func main() {
	// init logger and context
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	ctx := context.Background()
	// get runtime env
	runtimeEnv := os.Getenv(GO_ENV)
	if runtimeEnv != runtimeDev && runtimeEnv != runtimeProd {
		logger.Warn("invalid runtime env", "env", runtimeEnv)
		runtimeEnv = runtimeDev
	}

	// init repositories
	stockStore, err := repositories.NewStockStore(ctx, runtimeEnv, logger)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	// init handlers
	stockHandler := handlers.NewStockStoreHandler(stockStore, logger)

	// create a Gin router
	router := gin.Default()

	// define routers
	if runtimeEnv == runtimeDev {
		router.GET("/ping", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "pong"})
		})
	}
	transport.SetupStoreRoutes(router, stockHandler)

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", PORT_NUMBER),
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}
	server.ListenAndServe()
}
