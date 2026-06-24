package dto

type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type TokenResponse struct {
	Token     string `json:"token"`
	ExpiresIn int64  `json:"expires_in"`
}
