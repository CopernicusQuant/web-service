package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"

	"github.com/CopernicusQuant/web-service/internal/repositories"
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
	_, err := repositories.NewStockStore(ctx, runtimeEnv, logger)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	// create a Gin router
	r := gin.Default()

	// define routers
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	if err := r.Run(fmt.Sprintf(":%d", PORT_NUMBER)); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
