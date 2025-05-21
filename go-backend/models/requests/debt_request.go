package requests

type DebtRequest struct {
	UserId      string  `json:"user_id"`
	AmountToPay float64 `json:"amount_to_pay"`
}
