package types

type User struct {
	Fname     string `json:"fname"`
	Lname     string `json:"lname"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	Cpassword string `json:"cpassword"`
}
