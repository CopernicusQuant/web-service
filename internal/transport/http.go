package transport

import (
	"net/http"

	"github.com/CopernicusQuant/web-service/internal/handlers"
	"github.com/gin-gonic/gin"
)

func SetupStoreRoutes(router *gin.Engine, stockHandler *handlers.StockHandler) {
	storeGroup := router.Group("/store")

	storeGroup.GET("/stock/list", func(c *gin.Context) {
		stockList := stockHandler.GetStockList()
		OK(c, stockList)
	})

	storeGroup.GET("/stock/info/:ticker", func(c *gin.Context) {
		ticker := c.Param("ticker")
		stockInfo, err := stockHandler.GetStockInfo(ticker)
		if err != nil {
			Fail(c, http.StatusNotFound, err.Error())
			return
		}
		OK(c, stockInfo)

	})

	// storeGroup.GET("/stock/:ticker", func(c *gin.Context) {
	// 	ticker := c.Param("ticker")
	// 	window := c.Query("window")
	// 	c.JSON(http.StatusOK, gin.H{"ticker": ticker, "window": window})
	// })
}
