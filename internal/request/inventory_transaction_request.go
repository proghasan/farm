package request

type CreateInventoryTransactionRequest struct {
	InventoryItemID uint    `json:"inventory_item_id" validate:"required"`
	TransactionType string  `json:"transaction_type" validate:"required,oneof=Purchase Sale Consumption Adjustment Return Damage"`
	Quantity        float64 `json:"quantity" validate:"required"`
	TransactionDate string  `json:"transaction_date" validate:"required"`
	Remarks         *string `json:"remarks"`
}

type UpdateInventoryTransactionRequest struct {
	TransactionType *string  `json:"transaction_type"`
	Quantity        *float64 `json:"quantity"`
	TransactionDate *string  `json:"transaction_date"`
	Remarks         *string  `json:"remarks"`
}
