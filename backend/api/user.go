package api

type RegisterRequest struct {
	Name     string `json: "name" binding: "required"`
	Password string `json: "password" binding: "required"`
	Phone    string `json: "phone" binding: "required"`
	Email    string `json: "email" binding: "required,email"`
}

type LoginRequest struct {
	Name     string `json: "name"`
	Password string `json: "password"`
}

type LoginResponse struct {
	Token string
}

// json:"field, omitempty" omitempty的作用：当字段值为零值（如 ""、0、nil、false 等）时，该字段不会出现在 JSON 中。
// omitempty 只适用于 JSON/YAML 序列化，form表单数据需要手动检查空值
type UpdateUserRequest struct {
    Name     *string `json:"name,omitempty"`
    Password *string `json:"password,omitempty"`
    Email    *string `json:"email,omitempty" valid:"email"`
    Phone    *string `json:"phone,omitempty" valid:"matches(^1[3-9]{1}\\d{9}$)"`
}
