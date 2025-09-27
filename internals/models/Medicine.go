package models

// Medicine represents a medicine in the inventory
// @Description Details of a medicine including stock and dosage form
type Medicine struct {
	ID            uint      `json:"id"`
	Name          string    `json:"name"`
	DosageForm    string    `json:"dosage_form"`
	StockQuantity int       `json:"stock_quantity"`
}
