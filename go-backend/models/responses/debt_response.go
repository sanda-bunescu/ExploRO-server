package responses

type DebtResponse struct {
	Id          uint    `json:"id"`
	UserId      string  `json:"user_id"`
	UserName    string  `json:"user_name"`
	AmountToPay float64 `json:"amount_to_pay"`
}
