package networks

import (
	"api/services/authentication"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func Show(c echo.Context) error {	
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*authentication.Bearer)
	accessToken := claims.AccessToken

	// Build the request
	req, err := http.NewRequest("GET", fmt.Sprintf("https://demo.infiny.cloud/api/services/%s/service", c.Param("id")), nil)
	if err != nil {
		fmt.Println("Error is req: ", err)
	}

    var bearer = "Bearer " + accessToken
	req.Header.Set("Authorization", bearer)

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
	var baseResponse map[string]interface{}
	// Use json.Decode for reading streams of JSON data
	if err := json.NewDecoder(resp.Body).Decode(&baseResponse); err != nil {
		var baseResponse OutsideResponse
		if err := json.NewDecoder(resp.Body).Decode(&baseResponse); err != nil {
			fmt.Println(err)
			return echo.NewHTTPError(http.StatusBadRequest, "Something went really wrong")
		}
		fmt.Println("StatusCode: ", baseResponse.StatusCode)
		// if baseResponse.StatusCode == 401 {}
		return echo.NewHTTPError(http.StatusUnauthorized, "Please provide valid credentials")
	}

	return json.NewEncoder(c.Response()).Encode(baseResponse)
}