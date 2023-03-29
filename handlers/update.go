package handlers

import (
	"fmt"
	"interncase/database"
	"interncase/models"
	"interncase/utils"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

// UpdateStudent öğrenci bilgilerini güncelle
func UpdateStudent() echo.HandlerFunc {
	return func(c echo.Context) error {
		// JWT'den öğrenci bilgilerini çek
		token, ok := c.Get("user").(*jwt.Token)
		fmt.Println("Tok", token)
		if !ok {
			return c.JSON(http.StatusBadRequest, utils.ReturnMess("Bir hata oluştu"))
		}
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return c.JSON(http.StatusBadRequest, utils.ReturnMess("Bir hata oluştu"))
		}

		// Database'den öğrenci bilgilerini çek
		student := new(models.Student)
		if err := database.DB.Where("id = ?", claims["ID"]).First(&student).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, utils.ReturnMess(err.Error()))
		}

		// Öğrenci bilgilerini güncelle
		if err := c.Bind(student); err != nil {
			return c.JSON(http.StatusBadRequest, utils.ReturnMess(err.Error()))
		}

		// Öğrenci bilgilerini kontrol et ve reqiuired ise
		if err := c.Validate(student); err != nil {
			if mess, err := utils.ValidateUpdateError(err.(validator.ValidationErrors)); err != nil {
				return c.JSON(http.StatusBadRequest, utils.ReturnMess(mess))
			}
		}

		// Öğrencinin email adresi daha önce kullanılmış mı?
		var count int64
		if err := database.DB.Model(&models.Student{}).Where("email = ?", student.Email).Count(&count).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, utils.ReturnMess(err.Error()))
		}

		if count > 0 {
			return c.JSON(http.StatusBadRequest, utils.ReturnMess("Bu email adresi daha önce kullanılmış"))
		}

		// Şifreyi hashle
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(student.Password), bcrypt.DefaultCost)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, utils.ReturnMess(err.Error()))
		}
		student.Password = string(hashedPassword)

		// Öğrenciyi veritabanına kaydet
		if err := database.DB.Create(&student).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, utils.ReturnMess(err.Error()))
		}

		return c.JSON(http.StatusCreated, student)
	}
}
