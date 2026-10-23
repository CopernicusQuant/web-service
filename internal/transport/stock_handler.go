package transport

import (
	"context"
	"fmt"
	"log/slog"
	"slices"
	"strings"

	"github.com/copernicus-ml/web-service/internal/repositories"
)

const approxDaysPerYear int = 252

const (
	PriceMomentum  = "pricemomentum"
	VolumeMomentum = "volumemomentum"
	Volatility     = "volatility"
	MarketActivity = "marketactivity"
	BollingerBands = "bollingerbands"
	MACD           = "macd"
	DEMA           = "dema"
	RSI            = "rsi"
	KDJ            = "kdj"
	OBV            = "obv"
	CCI            = "cci"
)

var validFeatureGroups []string = []string{
	PriceMomentum,
	VolumeMomentum,
	Volatility,
	MarketActivity,
	BollingerBands,
	MACD,
	DEMA,
	RSI,
	KDJ,
	OBV,
	CCI,
}

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

func (h *StockHandler) GetStockList() [][]string {
	inMemStockList := h.stockStore.GetStockList()
	stockList := make([][]string, 0, len(inMemStockList))
	for ticker, info := range inMemStockList {
		item := []string{ticker, info.CompanyName}
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

// get stock prices with ticker and data window
func (h *StockHandler) GetStockPrices(
	ctx context.Context,
	ticker string,
	window string) ([]StockPriceResponse, error) {
	funcLogger := h.logger.With("path", "internal/handlers/stock/GetStockPrices")
	ticker = strings.ToUpper(ticker)
	stockPrices, err := h.stockStore.ReadStock(ctx, ticker)
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

// get stock features with ticker, feature group name, and data window
func (h *StockHandler) GetStockFeatures(
	ctx context.Context,
	ticker,
	featureGroup string,
	window string) (any, error) {
	funcLogger := h.logger.With("path", "internal/handlers/stock/GetStockFeatures")
	ticker = strings.ToUpper(ticker)
	featureGroup = strings.ToLower(featureGroup)
	// get stock feature raw data
	stockFeatures, err := h.stockStore.ReadFeatures(ctx, ticker)
	if err != nil {
		funcLogger.Error("get features", "error", err)
		return nil, fmt.Errorf("get features data")
	}
	// get return data size
	dataLength, ok := dataWindow[window]
	if !ok {
		funcLogger.Error("invalid data window", "window", window)
		return nil, fmt.Errorf("invalid data window %q", window)
	}
	var response any
	startIndex := max(0, len(stockFeatures)-dataLength)
	slicedFeatures := stockFeatures[startIndex:]
	// convert raw data to the corresponding feature group
	switch featureGroup {
	case PriceMomentum:
		response = toPriceMomentumResponse(slicedFeatures)
	case VolumeMomentum:
		response = toVolumeMomentumResponse(slicedFeatures)
	case Volatility:
		response = toVolatilityResponse(slicedFeatures)
	case MarketActivity:
		response = toMarketActivityResponse(slicedFeatures)
	case BollingerBands:
		response = toBollingerBandsResponse(slicedFeatures)
	case MACD:
		response = toMACDResponse(slicedFeatures)
	case DEMA:
		response = toDEMAResponse(slicedFeatures)
	case RSI:
		response = toRSIResponse(slicedFeatures)
	case KDJ:
		response = toKDJResponse(slicedFeatures)
	case OBV:
		response = toOBVFeatureResponse(slicedFeatures)
	case CCI:
		response = toCCIResponse(slicedFeatures)
	}
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

func (h *StockHandler) ValidateFeatureGroup(featureGroup string) error {
	featureGroup = strings.ToLower(featureGroup)
	if !slices.Contains(validFeatureGroups, featureGroup) {
		return fmt.Errorf("invalid feature group: %q", featureGroup)
	}
	return nil
}
