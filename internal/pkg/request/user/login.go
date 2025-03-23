package user

type Login struct {
	UserEmail string `json:"email"`
	Password  string `json:"password"`
}
