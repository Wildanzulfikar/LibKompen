package request

type UsersRequestBebasPustaka struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"`
	Status   bool   `json:"status"`
}
