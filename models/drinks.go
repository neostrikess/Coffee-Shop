package models

// Drink représente une boisson du menu
type Drink struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	Category  string  `json:"category"` // coffee, tea, cold
	BasePrice float64 `json:"base_price"`
}
