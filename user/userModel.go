package user

type CreateUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Status   string `json:"status"`
}
