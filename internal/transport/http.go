package transport

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetupStoreRoutes(router *gin.Engine, stockHandler *StockHandler) {
	storeGroup := router.Group("/data")

	// get stock list
	storeGroup.GET("/stock/list", func(c *gin.Context) {
		stockList := stockHandler.GetStockList()
		OK(c, stockList)
	})

	// get stock infomation
	storeGroup.GET("/stock/info/:ticker", func(c *gin.Context) {
		ticker := c.Param("ticker")
		stockInfo, err := stockHandler.GetStockInfo(ticker)
		if err != nil {
			Fail(c, http.StatusNotFound, err.Error())
			return
		}
		OK(c, stockInfo)

	})

	// get stock prices data
	storeGroup.GET("/stock/:ticker", func(c *gin.Context) {
		ticker := c.Param("ticker")
		if err := stockHandler.ValidateTicker(ticker); err != nil {
			Fail(c, http.StatusNotFound, err.Error())
			return
		}
		window := c.Query("window")
		if err := stockHandler.ValidateWindow(window); err != nil {
			Fail(c, http.StatusBadRequest, err.Error())
			return
		}
		stockPrices, err := stockHandler.GetStockPrices(c.Request.Context(), ticker, window)
		if err != nil {
			Fail(c, http.StatusInternalServerError, err.Error())
			return
		}
		OK(c, stockPrices)
	})

	// get stock features data
	storeGroup.GET("/features/:ticker", func(c *gin.Context) {
		ticker := c.Param("ticker")
		if err := stockHandler.ValidateTicker(ticker); err != nil {
			Fail(c, http.StatusNotFound, err.Error())
			return
		}
		featureGroup := c.Query("group")
		if err := stockHandler.ValidateFeatureGroup(featureGroup); err != nil {
			Fail(c, http.StatusNotFound, err.Error())
			return
		}
		window := c.Query("window")
		if err := stockHandler.ValidateWindow(window); err != nil {
			Fail(c, http.StatusBadRequest, err.Error())
			return
		}
		stockFeatures, err := stockHandler.GetStockFeatures(
			c.Request.Context(), ticker, featureGroup, window)
		if err != nil {
			Fail(c, http.StatusInternalServerError, err.Error())
			return
		}
		OK(c, stockFeatures)
	})
}
