package handlers

import (
	"fmt"
	"interncase/database"
	"interncase/models"
	"interncase/utils"
	"net/http"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
)

// StudentProfile öğrenci profil
func StudentProfile() echo.HandlerFunc {
	return func(c echo.Context) error {

		// JWT'den öğrenci bilgilerini çek
		//local
		token, ok := c.Get("user").(*jwt.Token)
		fmt.Println("Tok", token)
		if !ok {
			return c.JSON(http.StatusBadRequest, utils.NewError("Bir hata oluştu"))
		}
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return c.JSON(http.StatusBadRequest, utils.NewError("Bir hata oluştu"))
		}

		fmt.Println("ID : ", claims["ID"])

		// Database'den öğrenci bilgilerini çek
		student := new(models.Student)
		if err := database.DB.Where("id = ?", claims["ID"]).First(&student).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, utils.NewError(err.Error()))
		}
		return c.JSON(http.StatusOK, student)
	}
}
