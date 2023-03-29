package routes

import (
	"interncase/database"
	"interncase/handlers"

	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

// LoadRoutes routes dosyası
func LoadRoutes(e *echo.Echo) {
	// Öğrenci kaydı oluştur

	auth := e.Group("/auth")

	auth.POST("/register", handlers.RegisterStudent())
	// Öğrenci girişi yap
	auth.POST("/login", handlers.LoginStudent())

	student := e.Group("/student")

	// jwt middleware
	student.Use(echojwt.JWT(database.SecretKey))

	// Öğrenci bilgilerini görüntüle
	student.GET("/profile", handlers.StudentProfile())
	// Öğrenci bilgilerini güncelle
	student.PUT("/profile", handlers.UpdateStudent())

	// Plan oluştur
	student.POST("/createplan", handlers.CreatePlan())

	// Planları görüntüle
	student.GET("/plans", handlers.GetPlans())

	// Planı görüntüle
	student.GET("/plan/:id", handlers.GetPlan())

	// Planı güncelle
	student.PUT("/plan/:id", handlers.UpdatePlan())

	// Planı sil
	student.DELETE("/plan/:id", handlers.DeletePlan())

	// Statuleri görüntüle
	student.GET("/statuses", handlers.GetStatuDatas())

}
