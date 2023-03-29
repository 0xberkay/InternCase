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

// CreatePlan plan oluştur
func CreatePlan() echo.HandlerFunc {
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

		// Yeni bir plan modeli oluştur
		plan := new(models.Plan)
		if err := c.Bind(plan); err != nil {
			return c.JSON(http.StatusBadRequest, utils.NewError(err.Error()))
		}

		StudentID, ok := claims["ID"].(uint)
		if !ok {
			return c.JSON(http.StatusBadRequest, utils.NewError("Bir hata oluştu"))
		}

		plan.StudentID = StudentID
		// Plan bilgilerini kontrol et
		if err := c.Validate(plan); err != nil {
			return c.JSON(http.StatusBadRequest, utils.NewError(utils.ValidateError(err.(validator.ValidationErrors))))
		}

		//Eğer plan başka bir planın başka bir plan ile tarih aralığında çakışıyorsa ve kullanıcya aitse
		//hata mesajı döndür
		var count int64
		if err := database.DB.Model(&models.Plan{}).Where("student_id = ?", plan.StudentID).Where("((start_date <= ? AND end_date >= ?) OR (start_date <= ? AND end_date >= ?))", plan.Start, plan.Start, plan.End, plan.End).Count(&count).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, utils.NewError(err.Error()))
		}
		// Planı veritabanına kaydet
		if err := database.DB.Create(&plan).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, utils.NewError(err.Error()))
		}

		return c.JSON(http.StatusCreated, plan)
	}
}
