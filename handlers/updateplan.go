package handlers

import (
	"interncase/database"
	"interncase/models"
	"interncase/utils"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

// UpdatePlan planı güncelle
func UpdatePlan() echo.HandlerFunc {
	return func(c echo.Context) error {
		// Planı veritabanından çek
		plan := new(models.Plan)
		if err := database.DB.Where("id = ?", c.Param("id")).First(&plan).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, utils.NewError(err.Error()))
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
