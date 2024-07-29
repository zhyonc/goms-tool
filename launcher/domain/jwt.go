package domain

// NewTokenResponse
type LoginRequest struct {
	Username string `json:"username" binding:"required,min=5,max=20"`
	Password string `json:"password" binding:"required,min=5,max=20"`
}

// NewTokenResponse
type RenewTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// LoginRequest & RenewRequest
type NewTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
