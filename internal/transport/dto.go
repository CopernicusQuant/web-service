package transport

import "github.com/CopernicusQuant/web-service/internal/repositories"

type StockInfoResponse struct {
	TSCode       string `json:"ts_code"`
	CompanyName  string `json:"company_name"`
	Sector       string `json:"sector"`
	SubIndustry  string `json:"sub_industry"`
	Headquarters string `json:"headquarters"`
	Founded      string `json:"founded"`
}

type StockPriceResponse struct {
	TsCode    string  `json:"ts_code"`
	TradeDate string  `json:"trade_date"`
	AdjOpen   float32 `json:"adj_open"`
	AdjClose  float32 `json:"adj_close"`
	AdjHigh   float32 `json:"adj_high"`
	AdjLow    float32 `json:"adj_low"`
	AdjVol    float32 `json:"adj_vol"`
	Turnover  float32 `json:"turnover"`
}

func toStockPricesResponse(raw []repositories.StockPrice) []StockPriceResponse {
	response := make([]StockPriceResponse, len(raw))
	for i, price := range raw {
		response[i] = StockPriceResponse{
			TsCode:    price.TsCode,
			TradeDate: price.TradeDate,
			AdjOpen:   price.AdjOpen,
			AdjClose:  price.AdjClose,
			AdjHigh:   price.AdjHigh,
			AdjLow:    price.AdjLow,
			AdjVol:    price.AdjVol,
			Turnover:  price.Turnover,
		}
	}
	return response
}
