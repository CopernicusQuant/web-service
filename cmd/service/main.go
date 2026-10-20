package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/CopernicusQuant/web-service/internal/repositories"
	"github.com/CopernicusQuant/web-service/internal/transport"
)

const (
	portNumber = 8800
	goEnv      = "GO_ENV"
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
	runtimeEnv := os.Getenv(goEnv)
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
	stockHandler := transport.NewStockStoreHandler(stockStore, logger)

	// create a Gin router
	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// define routers
	if runtimeEnv == runtimeDev {
		router.GET("/ping", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "pong"})
		})
	}
	transport.SetupStoreRoutes(router, stockHandler)

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", portNumber),
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}
	if err := server.ListenAndServe(); err != nil {
		log.Fatal("failed to start application %w", err)
		os.Exit(1)
	}
}
