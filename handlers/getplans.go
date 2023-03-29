package handlers

import (
	"interncase/database"
	"interncase/models"
	"interncase/utils"
	"net/http"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
)

// GetPlans planları görüntüle
func GetPlans() echo.HandlerFunc {
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

		StudentID := uint(claims["ID"].(float64))

		plans := new([]models.Plan)

		// Planları veritabanından çek (tarih aralığı varsa)
		if c.QueryParam("start") != "" && c.QueryParam("end") != "" {
			if err := database.DB.Where("student_id = ? AND start >= ? AND end <= ?", StudentID, c.QueryParam("start"), c.QueryParam("end")).Find(&plans).Error; err != nil {
				return c.JSON(http.StatusInternalServerError, utils.ReturnMess(err.Error()))
			}
		} else {
			if err := database.DB.Where("student_id = ?", StudentID).Find(&plans).Error; err != nil {
				return c.JSON(http.StatusInternalServerError, utils.ReturnMess(err.Error()))
			}
		}

		return c.JSON(http.StatusOK, plans)
	}
}
