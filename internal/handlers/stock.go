package handlers

import (
	"log/slog"
	"strings"

	"github.com/CopernicusQuant/web-service/internal/repositories"
)

type StockHandler struct {
	stockStore *repositories.StockStore
	logger     *slog.Logger
}

func NewStockStoreHandler(
	stockStore *repositories.StockStore,
	logger *slog.Logger,
) *StockHandler {
	return &StockHandler{stockStore: stockStore, logger: logger}
}

func (h *StockHandler) GetStockList() []map[string]string {
	inMemStockList := h.stockStore.GetStockList()
	stockList := make([]map[string]string, 0, len(*inMemStockList))
	for ticker, info := range *inMemStockList {
		item := map[string]string{ticker: info.CompanyName}
		stockList = append(stockList, item)
	}
	return stockList
}

func (h *StockHandler) GetStockInfo(ticker string) (*StockInfoDTO, error) {
	funcLogger := h.logger.With("path", "internal/handlers/stock/GetStockInfo")
	ticker = strings.ToUpper(ticker)
	stockInfo, err := h.stockStore.GetStockInfo(ticker)
	if err != nil {
		funcLogger.Error("failed to get stockInfo", "error", err)
		return nil, err
	}
	return &StockInfoDTO{
		TSCode:       stockInfo.TSCode,
		CompanyName:  stockInfo.CompanyName,
		Sector:       stockInfo.Sector,
		SubIndustry:  stockInfo.SubIndustry,
		Headquarters: stockInfo.Headquarters,
		Founded:      stockInfo.Founded,
	}, nil

}
