package responses

type GroupUserResponse struct {
	UserID    string `json:"user_id"`
	UserName  string `json:"user_name"`
	UserEmail string `json:"user_email"`
}
