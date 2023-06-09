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

// CreatePlan plan oluştur
func CreatePlan() echo.HandlerFunc {
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

		// Yeni bir plan modeli oluştur
		plan := new(models.Plan)
		if err := c.Bind(plan); err != nil {
			return c.JSON(http.StatusBadRequest, utils.ReturnMess(err.Error()))
		}

		if plan.Start.After(plan.End) {
			return c.JSON(http.StatusBadRequest, utils.ReturnMess("Başlangıç tarihi bitiş tarihinden sonra olamaz"))
		}

		if plan.Start.Before(time.Now()) {
			return c.JSON(http.StatusBadRequest, utils.ReturnMess("Başlangıç tarihi bugünden önce olamaz"))
		}

		StudentID, ok := claims["ID"]
		if !ok {
			return c.JSON(http.StatusBadRequest, utils.ReturnMess("Bir hata oluştu"))
		}

		plan.StudentID = uint(StudentID.(float64))

		// Plan bilgilerini kontrol et

		if err := c.Validate(plan); err != nil {
			return c.JSON(http.StatusBadRequest, utils.ReturnMess(utils.ValidateError(err.(validator.ValidationErrors))))
		}

		//Eğer plan başka bir planın başka bir plan ile tarih aralığında çakışıyorsa ve kullanıcya aitse
		//hata mesajı döndür
		var count int64

		if err := database.DB.Model(&models.Plan{}).Where("student_id = ?", StudentID).Where("start < ?", plan.End).Where("end > ?", plan.Start).Count(&count).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, utils.ReturnMess(err.Error()))
		}

		if count > 0 {
			return c.JSON(http.StatusBadRequest, utils.ReturnMess("Bu tarih aralığında başka bir planınız var"))
		}

		// Planı veritabanına kaydet
		if err := database.DB.Create(&plan).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, utils.ReturnMess(err.Error()))
		}

		return c.JSON(http.StatusCreated, plan)
	}
}
