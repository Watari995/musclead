package authdto

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AccessTokenResponse struct {
	AccessToken          string `json:"access_token"`
	AccessTokenExpiresAt string `json:"access_token_expires_at"`
}
