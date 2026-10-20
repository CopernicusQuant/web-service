package repositories

type StockInfo struct {
	TSCode       string
	CompanyName  string
	Sector       string
	SubIndustry  string
	Headquarters string
	DateAdded    string // Preserve YYYYMMDD exactly; parse later if needed.
	CIK          string
	Founded      string // string, the value can be like 2019 (1860)
}

type StockPrice struct {
	TsCode    string  `parquet:"ts_code"`
	TradeDate string  `parquet:"trade_date"`
	AdjOpen   float32 `parquet:"adj_open"`
	AdjClose  float32 `parquet:"adj_close"`
	AdjHigh   float32 `parquet:"adj_high"`
	AdjLow    float32 `parquet:"adj_low"`
	AdjVol    float32 `parquet:"adj_vol"`
	Turnover  float32 `parquet:"turnover"`
}

type StockFeature struct {
	TsCode    string `parquet:"ts_code"`
	TradeDate string `parquet:"trade_date"`
	// Price momentum features
	Ma5        float32 `parquet:"ma_5"`
	Ma20       float32 `parquet:"ma_20"`
	Ma60       float32 `parquet:"ma_60"`
	Ma5Bias    float32 `parquet:"ma_5_bias"`
	Ma20Bias   float32 `parquet:"ma_20_bias"`
	Ma60Bias   float32 `parquet:"ma_60_bias"`
	Return5d   float32 `parquet:"return_5d"`
	Return20d  float32 `parquet:"return_20d"`
	Return60d  float32 `parquet:"return_60d"`
	UpRatio5d  float32 `parquet:"up_ratio_5d"`
	UpRatio20d float32 `parquet:"up_ratio_20d"`

	// Volume momentum features
	VolMa5Bias     float32 `parquet:"vol_ma_5_bias"`
	VolMa10Bias    float32 `parquet:"vol_ma_10_bias"`
	VolMa20Bias    float32 `parquet:"vol_ma_20_bias"`
	VolMa60Bias    float32 `parquet:"vol_ma_60_bias"`
	VolChange5d    float32 `parquet:"vol_change_5d"`
	VolChange10d   float32 `parquet:"vol_change_10d"`
	VolChange20d   float32 `parquet:"vol_change_20d"`
	VolChange60d   float32 `parquet:"vol_change_60d"`
	VolMa510Ratio  float32 `parquet:"vol_ma_5_10_ratio"`
	VolMa1020Ratio float32 `parquet:"vol_ma_10_20_ratio"`
	VolMa2060Ratio float32 `parquet:"vol_ma_20_60_ratio"`

	// OBV features
	ObvStrength5        float32 `parquet:"obv_strength_5"`
	ObvStrength20       float32 `parquet:"obv_strength_20"`
	ObvStrengthAccel520 float32 `parquet:"obv_strength_accel_5_20"`

	// Market activity features
	Amplitude                  float32 `parquet:"amplitude"`
	AmplitudeMa5               float32 `parquet:"amplitude_ma_5"`
	AmplitudeMa20              float32 `parquet:"amplitude_ma_20"`
	ActivityScore60            float32 `parquet:"activity_score_60"`
	ThinTradeAmplitudeScore60  float32 `parquet:"thin_trade_amplitude_score_60"`
	TurnoverWithoutMoveScore60 float32 `parquet:"turnover_without_move_score_60"`
	AmplitudeQuantile60        float32 `parquet:"amplitude_quantile_60"`
	TurnoverQuantile60         float32 `parquet:"turnover_quantile_60"`

	// MACD features
	MacdEmaSlow float32 `parquet:"macd_ema_slow"`
	MacdEmaFast float32 `parquet:"macd_ema_fast"`
	MacdDiff    float32 `parquet:"macd_diff"`
	MacdDea     float32 `parquet:"macd_dea"`
	MacdHist    float32 `parquet:"macd_hist"`
	MacdGold    int32   `parquet:"macd_gold"`
	MacdDead    int32   `parquet:"macd_dead"`

	// Bollinger Bands features
	BbUpper    float32 `parquet:"bb_upper"`
	BbLower    float32 `parquet:"bb_lower"`
	BbMid      float32 `parquet:"bb_mid"`
	BbWidth    float32 `parquet:"bb_width"`
	BbPosition float32 `parquet:"bb_position"`

	// DEMA features
	Dema510Fast    float32 `parquet:"dema_5_10_fast"`
	Dema1020Fast   float32 `parquet:"dema_10_20_fast"`
	Dema2060Fast   float32 `parquet:"dema_20_60_fast"`
	Dema510Slow    float32 `parquet:"dema_5_10_slow"`
	Dema1020Slow   float32 `parquet:"dema_10_20_slow"`
	Dema2060Slow   float32 `parquet:"dema_20_60_slow"`
	Dema510Spread  float32 `parquet:"dema_5_10_spread"`
	Dema1020Spread float32 `parquet:"dema_10_20_spread"`
	Dema2060Spread float32 `parquet:"dema_20_60_spread"`
	Dema510Gold    int32   `parquet:"dema_5_10_gold"`
	Dema1020Gold   int32   `parquet:"dema_10_20_gold"`
	Dema2060Gold   int32   `parquet:"dema_20_60_gold"`
	Dema510Dead    int32   `parquet:"dema_5_10_dead"`
	Dema1020Dead   int32   `parquet:"dema_10_20_dead"`
	Dema2060Dead   int32   `parquet:"dema_20_60_dead"`

	// RSI features
	Rsi6  float32 `parquet:"rsi_6"`
	Rsi14 float32 `parquet:"rsi_14"`

	// KDJ features
	KdjLow9  float32 `parquet:"kdj_low_9"`
	KdjHigh9 float32 `parquet:"kdj_high_9"`
	KdjK     float32 `parquet:"kdj_k"`
	KdjD     float32 `parquet:"kdj_d"`
	KdjJ     float32 `parquet:"kdj_j"`
}
