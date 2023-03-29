package handlers

import (
	"interncase/database"
	"interncase/models"
	"interncase/utils"
	"net/http"
	"time"

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
			return c.JSON(http.StatusBadRequest, utils.ReturnMess("Bir hata oluştu"))
		}
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return c.JSON(http.StatusBadRequest, utils.ReturnMess("Bir hata oluştu"))
		}

		StudentID := uint(claims["ID"].(float64))

		// Planı veritabanından çek
		plan := new(models.Plan)
		if err := database.DB.Where("id = ?", c.Param("id")).First(&plan).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, utils.ReturnMess(err.Error()))
		}

		// Planın öğrencisi ile JWT'deki öğrenci aynı mı kontrol et
		if plan.StudentID != StudentID {
			return c.JSON(http.StatusUnauthorized, utils.ReturnMess("Bu planı güncelleyemezsiniz"))
		}

		if (plan.Start != time.Time{}) || (plan.End != time.Time{}) {
			var count int64
			if err := database.DB.Model(&models.Plan{}).Where("student_id = ?", plan.StudentID).Where("((start <= ? AND end >= ?) OR (start <= ? AND end >= ?))", plan.Start, plan.Start, plan.End, plan.End).Count(&count).Error; err != nil {
				return c.JSON(http.StatusInternalServerError, utils.ReturnMess(err.Error()))
			}
			if count > 0 {
				return c.JSON(http.StatusBadRequest, utils.ReturnMess("Bu tarihler arasında başka bir planınız var"))
			}
		}

		// Yeni plan bilgilerini al
		newPlan := new(models.Plan)
		if err := c.Bind(newPlan); err != nil {
			return c.JSON(http.StatusBadRequest, utils.ReturnMess(err.Error()))
		}

		// Plan bilgilerini kontrol et
		if err := c.Validate(newPlan); err != nil {
			if mess, err := utils.ValidateUpdateError(err.(validator.ValidationErrors)); err != nil {
				return c.JSON(http.StatusBadRequest, utils.ReturnMess(mess))
			}
		}

		newPlan.ID = plan.ID

		// Planı güncelle
		if err := database.DB.Model(&plan).Updates(newPlan).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, utils.ReturnMess(err.Error()))
		}

		return c.JSON(http.StatusOK, plan)
	}
}
