package controller

type ResetPasswordRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type NewPasswordRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,passwrd"`
	Code     string `json:"code" validate:"required"`
}
