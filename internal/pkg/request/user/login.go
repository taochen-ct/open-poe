package user

type Login struct {
	UserEmail string `json:"user_email"`
	Password  string `json:"password"`
}
