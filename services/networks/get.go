package networks

import (
	"api/services/authentication"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func List(c echo.Context) error {	
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*authentication.Bearer)
	accessToken := claims.AccessToken

	start := time.Now().UnixNano() / int64(time.Millisecond)
	// Build the request
	req, err := http.NewRequest("GET", "https://demo.infiny.cloud/api/services", nil)
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
	end := time.Now().UnixNano() / int64(time.Millisecond)
	diff := end - start
	log.Printf("Duration(ms): %d", diff) // Slow external API!!!

	// Fill the baseResponse with the data from the JSON
	var baseResponse Services
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

	return json.NewEncoder(c.Response()).Encode(baseResponse.Service)
}

type Services struct {
	Service []struct {
		ID             uint64 `json:"id"`
		Name		   string `json:"name"`
		Expired		   bool   `json:"expired"`
		Paused		   bool   `json:"paused"`
		CreatedAt	   *string `json:"created"`
		CanceledAt     *string `json:"cancellation_date"`
		// CanceledAt     sql.NullString `json:"cancellation_date"`
		Type		   string `json:"type"`
	} `json:"services"`
}

type OutsideResponse struct {
	Message        string `json:"message"`
	StatusCode     uint16 `json:"status_code"`
}