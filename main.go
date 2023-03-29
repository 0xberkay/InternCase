package main

import (
	"interncase/database"
	"interncase/envs"
	"interncase/routes"
	"interncase/utils"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func init() {
	envs.LoadEnv()
	database.Connect()
}

func main() {
	// Echo instance
	e := echo.New()
	e.Validator = &utils.CustomValidator{
		Validator: validator.New(),
	}

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	routes.LoadRoutes(e)

	// Start server
	e.Logger.Fatal(e.Start(":3001"))
}
