package api

type RegisterRequest struct {
	Name     string `json: "name" binding: "required"`
	Password string `json: "password" binding: "required"`
	Phone    string `json: "phone" binding: "required"`
	Email    string `json: "email" binding: "required, email"`
}

type LoginRequest struct {
	Name     string `json: "name"`
	Password string `json: "password"`
}

type LoginResponse struct {
	Token string
}
