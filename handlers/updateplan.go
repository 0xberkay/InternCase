package handlers

import (
	"interncase/database"
	"interncase/models"
	"interncase/utils"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
)

// UpdatePlan planı güncelle
func UpdatePlan() echo.HandlerFunc {
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

		StudentID, ok := claims["ID"].(uint)
		if !ok {
			return c.JSON(http.StatusBadRequest, utils.NewError("Bir hata oluştu"))
		}

		// Planı veritabanından çek
		plan := new(models.Plan)
		if err := database.DB.Where("id = ?", c.Param("id")).First(&plan).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, utils.NewError(err.Error()))
		}

		// Planın öğrencisi ile JWT'deki öğrenci aynı mı kontrol et
		if plan.StudentID != StudentID {
			return c.JSON(http.StatusUnauthorized, utils.NewError("Bu planı güncelleyemezsiniz"))
		}

		// Yeni plan bilgilerini al
		newPlan := new(models.Plan)
		if err := c.Bind(newPlan); err != nil {
			return c.JSON(http.StatusBadRequest, utils.NewError(err.Error()))
		}

		// Plan bilgilerini kontrol et
		if err := c.Validate(newPlan); err != nil {
			return c.JSON(http.StatusBadRequest, utils.NewError(utils.ValidateUpdateError(err.(validator.ValidationErrors))))
		}

		// Planı güncelle
		if err := database.DB.Model(&plan).Updates(newPlan).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, utils.NewError(err.Error()))
		}

		return c.JSON(http.StatusOK, plan)
	}
}
