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

		// Planın başlangıç ve bitiş tarihi boş değilse diğer planlarla çakışıyor mu kontrol et
		if (plan.Start != time.Time{}) || (plan.End != time.Time{}) {
			if err := database.DB.Where("student_id = ?", StudentID).Where("id != ?", plan.ID).Where("start <= ?", plan.End).Where("end >= ?", plan.Start).First(&plan).Error; err == nil {
				return c.JSON(http.StatusBadRequest, utils.ReturnMess("Bu plan başka bir planla çakışıyor"))
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
