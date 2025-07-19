package response

import "github.com/google/uuid"

type Dashboard struct {
	TotalSales         float64 `json:"totalSales"`
	TotalOrder         int64   `json:"totalOrder"`
	TotalPendingOrder  int64   `json:"totalPendingOrder"`
	TotalActiveProduct int64   `json:"totalActiveProduct"`

	StockAlert         []StockAlert         `json:"stockAlert"`
	BestSellingProduct []BestSellingProduct `json:"bestSellingProduct"`
}

type StockAlert struct {
	Id    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Stock int       `json:"stock"`
}

type BestSellingProduct struct {
	Id   uuid.UUID `json:"id"`
	Name string    `json:"name"`
	Sold int       `json:"sold"`
}
