package authentication

import "github.com/golang-jwt/jwt/v5"


type Bearer struct {
	AccessToken      string `json:"access_token"`
	TokenType     	 string `json:"token_type"`
	ExpiresIn   	 int64  `json:"expires_in"`
	RefreshToken     string `json:"refresh_token"`
    jwt.RegisteredClaims
}