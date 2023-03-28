package handlers

import (
	"interncase/database"
	"interncase/models"
	"interncase/utils"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v4"
	_ "github.com/labstack/echo-jwt/v4"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

// generateToken öğrenci için token oluşturur
func generateToken(id uint) (string, error) {
	claims := models.Claims{
		ID: id,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: &jwt.NumericDate{Time: time.Now().Add(time.Hour * 24)},
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(database.SecretKey)
}

// LoginStudent öğrenci girişi yapar
func LoginStudent() echo.HandlerFunc {
	return func(c echo.Context) error {
		student := new(models.StudentLoginForm)
		if err := c.Bind(student); err != nil {
			return c.JSON(http.StatusBadRequest, utils.NewError(err.Error()))
		}

		// Öğrenci bilgilerini kontrol et
		if err := c.Validate(student); err != nil {
			return c.JSON(http.StatusBadRequest, utils.NewError(utils.ValidateError(err.(validator.ValidationErrors))))
		}

		// Database'den öğrenci bilgilerini çek
		var studentDB models.Student
		if err := database.DB.Where("email = ?", student.Email).First(&studentDB).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, utils.NewError(err.Error()))
		}

		// Şifre doğru mu?
		if err := bcrypt.CompareHashAndPassword([]byte(studentDB.Password), []byte(student.Password)); err != nil {
			return c.JSON(http.StatusBadRequest, utils.NewError("Şifre yanlış"))
		}

		// Token oluştur
		token, err := generateToken(studentDB.ID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, utils.NewError(err.Error()))
		}

		return c.JSON(http.StatusOK, map[string]string{
			"token": token,
		})

	}
}
