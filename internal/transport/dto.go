package transport

import "github.com/copernicus-ml/web-service/internal/repositories"

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

type PriceMomentumResponse struct {
	TradeDate  string  `json:"trade_date"`
	Ma5        float32 `json:"ma_5"`
	Ma20       float32 `json:"ma_20"`
	Ma60       float32 `json:"ma_60"`
	Ma5Bias    float32 `json:"ma_5_bias"`
	Ma20Bias   float32 `json:"ma_20_bias"`
	Ma60Bias   float32 `json:"ma_60_bias"`
	Return5d   float32 `json:"return_5d"`
	Return20d  float32 `json:"return_20d"`
	Return60d  float32 `json:"return_60d"`
	UpRatio5d  float32 `json:"up_ratio_5d"`
	UpRatio20d float32 `json:"up_ratio_20d"`
}

func toPriceMomentumResponse(raw []repositories.StockFeature) []PriceMomentumResponse {
	response := make([]PriceMomentumResponse, len(raw))
	for i, features := range raw {
		response[i] = PriceMomentumResponse{
			TradeDate:  features.TradeDate,
			Ma5:        features.Ma5,
			Ma20:       features.Ma20,
			Ma60:       features.Ma60,
			Ma5Bias:    features.Ma5Bias,
			Ma20Bias:   features.Ma20Bias,
			Ma60Bias:   features.Ma60Bias,
			Return5d:   features.Return5d,
			Return20d:  features.Return20d,
			Return60d:  features.Return60d,
			UpRatio5d:  features.UpRatio5d,
			UpRatio20d: features.UpRatio20d,
		}
	}
	return response
}

type VolumeMomentumResponse struct {
	TradeDate      string  `json:"trade_date"`
	VolMa5Bias     float32 `json:"vol_ma_5_bias"`
	VolMa10Bias    float32 `json:"vol_ma_10_bias"`
	VolMa20Bias    float32 `json:"vol_ma_20_bias"`
	VolMa60Bias    float32 `json:"vol_ma_60_bias"`
	VolChange5d    float32 `json:"vol_change_5d"`
	VolChange10d   float32 `json:"vol_change_10d"`
	VolChange20d   float32 `json:"vol_change_20d"`
	VolChange60d   float32 `json:"vol_change_60d"`
	VolMa510Ratio  float32 `json:"vol_ma_5_10_ratio"`
	VolMa1020Ratio float32 `json:"vol_ma_10_20_ratio"`
	VolMa2060Ratio float32 `json:"vol_ma_20_60_ratio"`
}

func toVolumeMomentumResponse(raw []repositories.StockFeature) []VolumeMomentumResponse {
	response := make([]VolumeMomentumResponse, len(raw))
	for i, features := range raw {
		response[i] = VolumeMomentumResponse{
			TradeDate:      features.TradeDate,
			VolMa5Bias:     features.VolMa5Bias,
			VolMa10Bias:    features.VolMa10Bias,
			VolMa20Bias:    features.VolMa20Bias,
			VolMa60Bias:    features.VolMa60Bias,
			VolChange5d:    features.VolChange5d,
			VolChange10d:   features.VolChange10d,
			VolChange20d:   features.VolChange20d,
			VolChange60d:   features.VolChange60d,
			VolMa510Ratio:  features.VolMa510Ratio,
			VolMa1020Ratio: features.VolMa1020Ratio,
			VolMa2060Ratio: features.VolMa2060Ratio,
		}
	}
	return response
}

type VolatilityResponse struct {
	TradeDate      string  `json:"trade_date"`
	VolatilityLog  float32 `json:"volatility_log"`
	VolatilityRank float32 `json:"volatility_rank"`
}

func toVolatilityResponse(raw []repositories.StockFeature) []VolatilityResponse {
	response := make([]VolatilityResponse, len(raw))
	for i, features := range raw {
		response[i] = VolatilityResponse{
			TradeDate:      features.TradeDate,
			VolatilityLog:  features.VolatilityLog,
			VolatilityRank: features.VolatilityRank,
		}
	}
	return response
}

type OBVFeatureResponse struct {
	TradeDate           string  `json:"trade_date"`
	ObvStrength5        float32 `json:"obv_strength_5"`
	ObvStrength20       float32 `json:"obv_strength_20"`
	ObvStrengthAccel520 float32 `json:"obv_strength_accel_5_20"`
}

func toOBVFeatureResponse(raw []repositories.StockFeature) []OBVFeatureResponse {
	response := make([]OBVFeatureResponse, len(raw))
	for i, features := range raw {
		response[i] = OBVFeatureResponse{
			TradeDate:           features.TradeDate,
			ObvStrength5:        features.ObvStrength5,
			ObvStrength20:       features.ObvStrength20,
			ObvStrengthAccel520: features.ObvStrengthAccel520,
		}
	}
	return response
}

type MarketActivityResponse struct {
	TradeDate                  string  `json:"trade_date"`
	Amplitude                  float32 `json:"amplitude"`
	AmplitudeMa5               float32 `json:"amplitude_ma_5"`
	AmplitudeMa20              float32 `json:"amplitude_ma_20"`
	ActivityScore60            float32 `json:"activity_score_60"`
	ThinTradeAmplitudeScore60  float32 `json:"thin_trade_amplitude_score_60"`
	TurnoverWithoutMoveScore60 float32 `json:"turnover_without_move_score_60"`
	AmplitudeQuantile60        float32 `json:"amplitude_quantile_60"`
	TurnoverQuantile60         float32 `json:"turnover_quantile_60"`
}

func toMarketActivityResponse(raw []repositories.StockFeature) []MarketActivityResponse {
	response := make([]MarketActivityResponse, len(raw))
	for i, features := range raw {
		response[i] = MarketActivityResponse{
			TradeDate:                  features.TradeDate,
			Amplitude:                  features.Amplitude,
			AmplitudeMa5:               features.AmplitudeMa5,
			AmplitudeMa20:              features.AmplitudeMa20,
			ActivityScore60:            features.ActivityScore60,
			ThinTradeAmplitudeScore60:  features.ThinTradeAmplitudeScore60,
			TurnoverWithoutMoveScore60: features.TurnoverWithoutMoveScore60,
			AmplitudeQuantile60:        features.AmplitudeQuantile60,
			TurnoverQuantile60:         features.TurnoverQuantile60,
		}
	}
	return response
}

type MACDResponse struct {
	TradeDate   string  `json:"trade_date"`
	MacdEmaSlow float32 `json:"macd_ema_slow"`
	MacdEmaFast float32 `json:"macd_ema_fast"`
	MacdDiff    float32 `json:"macd_diff"`
	MacdDea     float32 `json:"macd_dea"`
	MacdHist    float32 `json:"macd_hist"`
	MacdGold    int32   `json:"macd_gold"`
	MacdDead    int32   `json:"macd_dead"`
}

func toMACDResponse(raw []repositories.StockFeature) []MACDResponse {
	response := make([]MACDResponse, len(raw))
	for i, features := range raw {
		response[i] = MACDResponse{
			TradeDate:   features.TradeDate,
			MacdEmaSlow: features.MacdEmaSlow,
			MacdEmaFast: features.MacdEmaFast,
			MacdDiff:    features.MacdDiff,
			MacdDea:     features.MacdDea,
			MacdHist:    features.MacdHist,
			MacdGold:    features.MacdGold,
			MacdDead:    features.MacdDead,
		}
	}
	return response
}

type BollingerBandsResponse struct {
	TradeDate  string  `json:"trade_date"`
	BbUpper    float32 `json:"bb_upper"`
	BbLower    float32 `json:"bb_lower"`
	BbMid      float32 `json:"bb_mid"`
	BbWidth    float32 `json:"bb_width"`
	BbPosition float32 `json:"bb_position"`
}

func toBollingerBandsResponse(raw []repositories.StockFeature) []BollingerBandsResponse {
	response := make([]BollingerBandsResponse, len(raw))
	for i, features := range raw {
		response[i] = BollingerBandsResponse{
			TradeDate:  features.TradeDate,
			BbUpper:    features.BbUpper,
			BbLower:    features.BbLower,
			BbMid:      features.BbMid,
			BbWidth:    features.BbWidth,
			BbPosition: features.BbPosition,
		}
	}
	return response
}

type DEMAResponse struct {
	TradeDate      string  `json:"trade_date"`
	Dema510Fast    float32 `json:"dema_5_10_fast"`
	Dema1020Fast   float32 `json:"dema_10_20_fast"`
	Dema2060Fast   float32 `json:"dema_20_60_fast"`
	Dema510Slow    float32 `json:"dema_5_10_slow"`
	Dema1020Slow   float32 `json:"dema_10_20_slow"`
	Dema2060Slow   float32 `json:"dema_20_60_slow"`
	Dema510Spread  float32 `json:"dema_5_10_spread"`
	Dema1020Spread float32 `json:"dema_10_20_spread"`
	Dema2060Spread float32 `json:"dema_20_60_spread"`
	Dema510Gold    int32   `json:"dema_5_10_gold"`
	Dema1020Gold   int32   `json:"dema_10_20_gold"`
	Dema2060Gold   int32   `json:"dema_20_60_gold"`
	Dema510Dead    int32   `json:"dema_5_10_dead"`
	Dema1020Dead   int32   `json:"dema_10_20_dead"`
	Dema2060Dead   int32   `json:"dema_20_60_dead"`
}

func toDEMAResponse(raw []repositories.StockFeature) []DEMAResponse {
	response := make([]DEMAResponse, len(raw))
	for i, features := range raw {
		response[i] = DEMAResponse{
			TradeDate:      features.TradeDate,
			Dema510Fast:    features.Dema510Fast,
			Dema1020Fast:   features.Dema1020Fast,
			Dema2060Fast:   features.Dema2060Fast,
			Dema510Slow:    features.Dema510Slow,
			Dema1020Slow:   features.Dema1020Slow,
			Dema2060Slow:   features.Dema2060Slow,
			Dema510Spread:  features.Dema510Spread,
			Dema1020Spread: features.Dema1020Spread,
			Dema2060Spread: features.Dema2060Spread,
			Dema510Gold:    features.Dema510Gold,
			Dema1020Gold:   features.Dema1020Gold,
			Dema2060Gold:   features.Dema2060Gold,
			Dema510Dead:    features.Dema510Dead,
			Dema1020Dead:   features.Dema1020Dead,
			Dema2060Dead:   features.Dema2060Dead,
		}
	}
	return response
}

type RSIResponse struct {
	TradeDate string  `json:"trade_date"`
	Rsi6      float32 `json:"rsi_6"`
	Rsi14     float32 `json:"rsi_14"`
}

func toRSIResponse(raw []repositories.StockFeature) []RSIResponse {
	response := make([]RSIResponse, len(raw))
	for i, features := range raw {
		response[i] = RSIResponse{
			TradeDate: features.TradeDate,
			Rsi6:      features.Rsi6,
			Rsi14:     features.Rsi14,
		}
	}
	return response
}

type KDJResponse struct {
	TradeDate string  `json:"trade_date"`
	KdjLow9   float32 `json:"kdj_low_9"`
	KdjHigh9  float32 `json:"kdj_high_9"`
	KdjK      float32 `json:"kdj_k"`
	KdjD      float32 `json:"kdj_d"`
	KdjJ      float32 `json:"kdj_j"`
}

func toKDJResponse(raw []repositories.StockFeature) []KDJResponse {
	response := make([]KDJResponse, len(raw))
	for i, features := range raw {
		response[i] = KDJResponse{
			TradeDate: features.TradeDate,
			KdjLow9:   features.KdjLow9,
			KdjHigh9:  features.KdjHigh9,
			KdjK:      features.KdjK,
			KdjD:      features.KdjD,
			KdjJ:      features.KdjJ,
		}
	}
	return response
}

type CCIResponse struct {
	TradeDate string  `json:"trade_date"`
	Tp        float32 `json:"tp"`
	TpSma     float32 `json:"tp_sma"`
	Cci       float32 `json:"cci"`
	CciMa5    float32 `json:"cci_ma_5"`
}

func toCCIResponse(raw []repositories.StockFeature) []CCIResponse {
	response := make([]CCIResponse, len(raw))
	for i, features := range raw {
		response[i] = CCIResponse{
			TradeDate: features.TradeDate,
			Tp:        features.Tp,
			TpSma:     features.TpSma,
			Cci:       features.Cci,
			CciMa5:    features.CciMa5,
		}
	}
	return response
}
