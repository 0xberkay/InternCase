package handlers

import (
	"interncase/database"
	"interncase/models"
	"interncase/utils"
	"net/http"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
)

// DeletePlan planı sil
func DeletePlan() echo.HandlerFunc {
	return func(c echo.Context) error {
		// JWT'den öğrenci bilgilerini çek
		token, ok := c.Get("user").(*jwt.Token)
		if !ok {
			return c.JSON(http.StatusBadRequest, utils.NewError("Bir hata oluştu"))
		}
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return c.JSON(http.StatusBadRequest, utils.NewError("Bir hata oluştu"))
		}

		// Planı veritabanından çek
		plan := new(models.Plan)
		if err := database.DB.Where("id = ?", c.Param("id")).First(&plan).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, utils.NewError(err.Error()))
		}

		StudentID := uint(claims["ID"].(float64))

		// Planı öğrenciye ait mi kontrol et
		if plan.StudentID != StudentID {
			return c.JSON(http.StatusUnauthorized, utils.NewError("Bu planı silemezsiniz"))
		}

		// Planı veritabanından sil
		if err := database.DB.Delete(&plan).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, utils.NewError(err.Error()))
		}

		return c.JSON(http.StatusOK, plan)
	}
}
