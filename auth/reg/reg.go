package reg

import (
	"be/conf"
	userstr "be/userStr"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/argon2"
)

func HashPassword(user *userstr.User) error {
	// Генерация соли
	salt := make([]byte, 16)
	_, err := rand.Read(salt)
	if err != nil {
		return err
	}

	// Хеширование пароля с использованием Argon2
	hash := argon2.IDKey([]byte(user.Password), salt, 1, 64*1024, 4, 32)
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)
	user.Password = fmt.Sprintf("%s.%s", b64Salt, b64Hash)

	return nil
}

func Registration(c *gin.Context) {
	var user userstr.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неудалось декодировать JSON"})
		return
	}

	if user.Username == "" || user.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Пустые поля"})
		return
	}

	errUserCheck := conf.DB.Where("email = ?", user.Username).First(&user).Error
	if errUserCheck == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Пользователь с таким email уже существует"})
		return
	}

	err := HashPassword(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось захешировать пароль"})
		return
	}

	errUser := conf.DB.Create(&user).Error
	if errUser != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось создать пользователя"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Пользователь успешно создан"})
}
