package request

type CreateAccountTransactionRequest struct {
	AccountHeadID   uint    `json:"account_head_id" validate:"required"`
	TransactionDate string  `json:"transaction_date" validate:"required"`
	Amount          float64 `json:"amount" validate:"required"`
	PaymentMethod   string  `json:"payment_method" validate:"omitempty,oneof=Cash Bank 'Mobile Banking' Other"`
	ReferenceNo     *string `json:"reference_no"`
	Description     *string `json:"description"`
}

type UpdateAccountTransactionRequest struct {
	TransactionDate *string  `json:"transaction_date"`
	Amount          *float64 `json:"amount"`
	PaymentMethod   *string  `json:"payment_method"`
	ReferenceNo     *string  `json:"reference_no"`
	Description     *string  `json:"description"`
}
