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
