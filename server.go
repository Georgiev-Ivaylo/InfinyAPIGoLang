package main

import (
	"api/services/authentication"
	"api/services/networks"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()
	e.Use(middleware.CORS())

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.POST("/api/oauth2/access-token", networks.CreateToken)


	r := e.Group("/api/services")
	{
		config := echojwt.Config{
			NewClaimsFunc: func(c echo.Context) jwt.Claims {
				return new(authentication.Bearer)
			},
			SigningKey: []byte("Awesome"),
		}
		
		r.Use(echojwt.WithConfig(config))

		r.GET("", networks.List)
		r.GET("/:id/service", networks.Show)
	}

	
	e.Logger.Fatal(e.Start(":1323"))
}