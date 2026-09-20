package dto

type CreateWalletRequest struct {
	UserID string `json:"user_id" binding:"required"`
}

type WalletResponse struct {
	ID        string  `json:"id"`
	UserID    string  `json:"user_id"`
	Balance   float64 `json:"balance"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
}

type TransactionRequest struct {
	Amount float64 `json:"amount" binding:"required"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
