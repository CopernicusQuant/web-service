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
