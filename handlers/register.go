package handlers

import (
	"interncase/database"
	"interncase/models"
	"interncase/utils"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"

	"golang.org/x/crypto/bcrypt"
)

// RegisterStudent öğrenci kaydı oluşturur
func RegisterStudent() echo.HandlerFunc {
	return func(c echo.Context) error {
		// Yeni bir öğrenci modeli oluştur
		student := new(models.Student)
		if err := c.Bind(student); err != nil {
			return c.JSON(http.StatusBadRequest, utils.NewError(err.Error()))
		}

		// Öğrenci bilgilerini kontrol et
		if err := c.Validate(student); err != nil {
			return c.JSON(http.StatusBadRequest, utils.NewError(utils.ValidateError(err.(validator.ValidationErrors))))
		}

		// Öğrencinin email adresi daha önce kullanılmış mı?
		var count int64
		if err := database.DB.Model(&models.Student{}).Where("email = ?", student.Email).Count(&count).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, utils.NewError(err.Error()))
		}

		if count > 0 {
			return c.JSON(http.StatusBadRequest, utils.NewError("Bu email adresi daha önce kullanılmış"))
		}

		// Şifreyi hashle
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(student.Password), bcrypt.DefaultCost)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, utils.NewError(err.Error()))
		}
		student.Password = string(hashedPassword)

		// Öğrenciyi veritabanına kaydet
		if err := database.DB.Create(&student).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, utils.NewError(err.Error()))
		}

		return c.JSON(http.StatusCreated, student)
	}
}
