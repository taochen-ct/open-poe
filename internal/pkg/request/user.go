package request

type Register struct {
	Username  string `json:"user_name"`
	UserEmail string `json:"user_email"`
	Password  string `json:"password"`
}
