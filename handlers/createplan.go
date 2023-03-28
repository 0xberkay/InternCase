package handlers

import (
	"interncase/database"
	"interncase/models"
	"interncase/utils"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

// CreatePlan plan oluştur
func CreatePlan() echo.HandlerFunc {
	return func(c echo.Context) error {
		// Yeni bir plan modeli oluştur
		plan := new(models.Plan)
		if err := c.Bind(plan); err != nil {
			return c.JSON(http.StatusBadRequest, utils.NewError(err.Error()))
		}

		// Plan bilgilerini kontrol et
		if err := c.Validate(plan); err != nil {
			return c.JSON(http.StatusBadRequest, utils.NewError(utils.ValidateError(err.(validator.ValidationErrors))))
		}

		// Planı veritabanına kaydet
		if err := database.DB.Create(&plan).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, utils.NewError(err.Error()))
		}

		return c.JSON(http.StatusCreated, plan)
	}
}
