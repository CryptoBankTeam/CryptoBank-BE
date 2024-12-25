package getprofile

import (
	"be/conf"
	userstr "be/userStr"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetProfile(c *gin.Context) {
	id, exists := c.Get("id")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Id not found in context"})
		return
	}
	id = int(id.(float64))

	var user userstr.User
	user.Id = id.(int)

	err := conf.DB.Where("id = ?", user.Id).First(&user).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	user.Password = ""

	c.JSON(http.StatusOK, user)
}
