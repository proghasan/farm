package request

type CreateInventoryItemRequest struct {
	CategoryID    uint    `json:"category_id" validate:"required"`
	Name          string  `json:"name" validate:"required,min=1,max=200"`
	SKU           *string `json:"sku"`
	Unit          string  `json:"unit" validate:"required,min=1,max=50"`
	PurchasePrice float64 `json:"purchase_price"`
	SellingPrice  float64 `json:"selling_price"`
}

type UpdateInventoryItemRequest struct {
	Name          *string  `json:"name" validate:"omitempty,min=1,max=200"`
	SKU           *string  `json:"sku"`
	Unit          *string  `json:"unit" validate:"omitempty,min=1,max=50"`
	PurchasePrice *float64 `json:"purchase_price"`
	SellingPrice  *float64 `json:"selling_price"`
}
