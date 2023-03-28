package handlers

import (
	"interncase/database"
	"interncase/models"
	"interncase/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

// GetPlans planları görüntüle
func GetPlans() echo.HandlerFunc {
	return func(c echo.Context) error {
		// Planları veritabanından çek
		plans := new([]models.Plan)
		if err := database.DB.Find(&plans).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, utils.NewError(err.Error()))
		}
		return c.JSON(http.StatusOK, plans)
	}
}
