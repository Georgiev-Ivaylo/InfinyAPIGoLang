package networks

import (
	"api/services/authentication"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func CreateToken(c echo.Context) error {
	var sendData SendToken
	err := json.NewDecoder(c.Request().Body).Decode(&sendData)
	if err != nil {
		fmt.Println("Parse Post: ", err)
	}

	body, err := json.Marshal(sendData)
	if err != nil {
		fmt.Println("Json Body: ", err)
	}
	fmt.Println(string(body))
	
	// Build the request
	req, err := http.NewRequest("POST", "https://demo.infiny.cloud/api/oauth2/access-token", bytes.NewBuffer(body))
	if err != nil {
		fmt.Println("Error is req: ", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// create a Client
	client := &http.Client{}

	// Do sends an HTTP request and
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("error in send req: ", err)
	}

	// Defer the closing of the body
	defer resp.Body.Close()

	// Fill the baseResponse with the data from the JSON
	var tokenResponse TokenResponse
	// Use json.Decode for reading streams of JSON data
	if err := json.NewDecoder(resp.Body).Decode(&tokenResponse); err != nil {
		var baseResponse OutsideResponse
		if err := json.NewDecoder(resp.Body).Decode(&baseResponse); err != nil {
			fmt.Println(err)
			return echo.NewHTTPError(http.StatusBadRequest, "Something went really wrong")
		}
		fmt.Println("StatusCode: ", baseResponse.StatusCode)
		// if baseResponse.StatusCode == 401 {}
		return echo.NewHTTPError(http.StatusUnauthorized, "Please provide valid credentials")
	}

    claims := &authentication.Bearer{
		AccessToken: tokenResponse.AccessToken,
		TokenType: tokenResponse.TokenType,
		ExpiresIn: tokenResponse.ExpiresIn,
		RefreshToken: tokenResponse.RefreshToken,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: &jwt.NumericDate{time.Now().Add(time.Minute * 3600)},
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

    accessToken, err := token.SignedString([]byte("Awesome"))
    if err != nil {
        return err
    }


	responseData := map[string]string{
        "access_token": accessToken,
    }
	return json.NewEncoder(c.Response()).Encode(responseData)
}
type SendToken struct {
    ClientId       string `json:"client_id" validate:"required,min=10,max=128"`
    ClientSecret   string `json:"client_secret" validate:"required,min=10,max=128"`
    GrantType      string `json:"grant_type" validate:"required"`
}

type TokenResponse struct {
	AccessToken      string `json:"access_token"`
	TokenType     	 string `json:"token_type"`
	ExpiresIn   	 int64  `json:"expires_in"`
	RefreshToken     string `json:"refresh_token"`
}