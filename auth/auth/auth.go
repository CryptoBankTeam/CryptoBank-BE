package auth

import (
	"be/conf"
	secretconf "be/secretConf"
	userStr "be/userStr"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/argon2"
)

// Вспомогательная функция для хеширования пароля с использованием соли
func hashPassword(password string, salt []byte) string {
	hash := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32)
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)
	return fmt.Sprintf("%s.%s", b64Salt, b64Hash)
}

func checkPassword(passwordClient string, passwordDB string) bool {
	parts := strings.Split(passwordDB, ".")
	if len(parts) != 2 {
		return false
	}
	salt, err := base64.RawStdEncoding.DecodeString(parts[0])
	if err != nil {
		return false
	}
	expectedHash := hashPassword(passwordClient, salt)
	return passwordDB == expectedHash
}

func Auth(c *gin.Context) {

	var userClient userStr.User
	if err := c.ShouldBindJSON(&userClient); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		log.Println(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	var userDB userStr.User
	errUser := conf.DB.Where("username=?", userClient.Username).First(&userDB).Error
	if errUser != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		log.Println(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if !checkPassword(userClient.Password, userDB.Password) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid password"})
		log.Println(http.StatusBadRequest, gin.H{"error": "Invalid password"})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  userDB.Id,
		"exp": time.Now().Add(30 * time.Second).Unix(), // Токен действует 30 секунд
	})
	tokenString, err := token.SignedString(secretconf.JWT_KEY)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate access token"})
		return
	}

	// Генерация refresh токена
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  userDB.Id,
		"exp": time.Now().Add(time.Hour * 24 * 7).Unix(), // Токен действует 7 дней
	})
	refreshTokenString, err := refreshToken.SignedString(secretconf.JWT_KEY)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate refresh token"})
		return
	}

	// Отправка токенов на фронт
	c.JSON(http.StatusOK, gin.H{
		"token":        tokenString,
		"refreshToken": refreshTokenString,
	})

}
