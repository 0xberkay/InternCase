package handlers

import (
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
		token, ok := c.Get("user").(*jwt.Token)
		if !ok {
			return c.JSON(http.StatusBadRequest, utils.ReturnMess("Bir hata oluştu"))
		}
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return c.JSON(http.StatusBadRequest, utils.ReturnMess("Bir hata oluştu"))
		}

		// Database'den öğrenci bilgilerini çek
		student := new(models.Student)
		if err := database.DB.Where("id = ?", claims["ID"]).First(&student).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, utils.ReturnMess(err.Error()))
		}

		student.Password = ""
		return c.JSON(http.StatusOK, student)
	}
}
