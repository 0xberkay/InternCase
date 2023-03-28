package handlers

import (
	"interncase/database"
	"interncase/models"
	"interncase/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

// GetPlan planı görüntüle
func GetPlan() echo.HandlerFunc {
	return func(c echo.Context) error {
		// Planı veritabanından çek
		plan := new(models.Plan)
		if err := database.DB.Where("id = ?", c.Param("id")).First(&plan).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, utils.NewError(err.Error()))
		}
		return c.JSON(http.StatusOK, plan)
	}
}