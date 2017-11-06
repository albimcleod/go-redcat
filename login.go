package goredcat

//LoginRequest is the struct to log in a user into Redcat
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"psw"`
	AuthType string `json:"auth_type"`
}

//LoginResponse is the response for a login request
type LoginResponse struct {
	Token   string `json:"token"`
	Success bool   `json:"success"`
}
