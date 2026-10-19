package handlers

type StockInfoDTO struct {
	TSCode       string `json:"ts_code"`
	CompanyName  string `json:"company_name"`
	Sector       string `json:"sector"`
	SubIndustry  string `json:"sub_industry"`
	Headquarters string `json:"headquarters"`
	Founded      string `json:"founded"`
}
