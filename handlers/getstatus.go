package handlers

import (
	"interncase/database"
	"interncase/models"
	"interncase/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

// GetStatus plan durumunu seçeneklerini görüntüle
func GetStatuDatas() echo.HandlerFunc {
	return func(c echo.Context) error {

		// Plan durumlarını veritabanından çek
		statuData := new([]models.StatuData)
		if err := database.DB.Find(&statuData).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, utils.ReturnMess(err.Error()))
		}
		return c.JSON(http.StatusOK, statuData)

	}
}
