package transport

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func SetupStoreRoutes(router *gin.Engine, stockHandler *StockHandler) {
	storeGroup := router.Group("/data")

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

	storeGroup.GET("/stock/:ticker", func(c *gin.Context) {
		ticker := c.Param("ticker")
		window := c.Query("window")
		stockPrices, err := stockHandler.GetStockPrices(c.Request.Context(), ticker, window)
		if err != nil {
			if strings.HasPrefix(err.Error(), ClientError) {
				Fail(c, http.StatusBadRequest, err.Error())
			} else {
				Fail(c, http.StatusInternalServerError, err.Error())
			}
			return
		}
		OK(c, stockPrices)
	})
}
