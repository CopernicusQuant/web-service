package transport

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/CopernicusQuant/web-service/internal/repositories"
)

const approxDaysPerYear int = 252

var dataWindow map[string]int = map[string]int{
	"60days":  60,
	"90days":  90,
	"180days": 180,
	"1year":   approxDaysPerYear,
	"3years":  3 * approxDaysPerYear,
}

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
	stockList := make([]map[string]string, 0, len(inMemStockList))
	for ticker, info := range inMemStockList {
		item := map[string]string{ticker: info.CompanyName}
		stockList = append(stockList, item)
	}
	return stockList
}

func (h *StockHandler) GetStockInfo(ticker string) (*StockInfoResponse, error) {
	funcLogger := h.logger.With("path", "internal/handlers/stock/GetStockInfo")
	ticker = strings.ToUpper(ticker)
	stockInfo, err := h.stockStore.GetStockInfo(ticker)
	if err != nil {
		funcLogger.Error("failed to get stockInfo", "error", err)
		return nil, err
	}
	return &StockInfoResponse{
		TSCode:       stockInfo.TSCode,
		CompanyName:  stockInfo.CompanyName,
		Sector:       stockInfo.Sector,
		SubIndustry:  stockInfo.SubIndustry,
		Headquarters: stockInfo.Headquarters,
		Founded:      stockInfo.Founded,
	}, nil

}

func (h *StockHandler) GetStockPrices(
	ctx context.Context,
	ticker string,
	window string) ([]StockPriceResponse, error) {
	funcLogger := h.logger.With("path", "internal/handlers/stock/GetStockPrices")
	ticker = strings.ToUpper(ticker)
	stockPrices, err := h.stockStore.ReadStock(ctx, ticker, window)
	if err != nil {
		funcLogger.Error("get stock prices data", "error", err)
		return nil, fmt.Errorf("get stock prices data: %w", err)
	}
	dataLength, ok := dataWindow[window]
	if !ok {
		funcLogger.Error("invalid data window", "window", window)
		return nil, fmt.Errorf("invalid data window %q", window)
	}
	startIndex := max(0, len(stockPrices)-dataLength)
	response := toStockPricesResponse(stockPrices[startIndex:])
	return response, nil
}

func (h *StockHandler) ValidateTicker(ticker string) error {
	ticker = strings.ToUpper(ticker)
	if !h.stockStore.ValidateTicker(ticker) {
		return fmt.Errorf("invalid ticker name: %q", ticker)
	}
	return nil
}

func (h *StockHandler) ValidateWindow(window string) error {
	if _, ok := dataWindow[window]; !ok {
		return fmt.Errorf("invalid window: %q", window)
	}
	return nil
}
