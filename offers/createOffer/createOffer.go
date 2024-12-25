package createoffer

import (
	"be/conf"
	offerstr "be/offers/offerStr"
	userstr "be/userStr"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateOffer(c *gin.Context) {
	id, exists := c.Get("id")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Id not found in context"})
		return
	}
	id = int(id.(float64))
	var user userstr.User
	user.Id = id.(int)

	var offer offerstr.Offer
	if err := c.ShouldBindJSON(&offer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}
	offer.Id_c = user.Id
	offer.IsGive = false
	offer.IsPub = true

	errCreate := conf.DB.Create(&offer).Error
	if errCreate != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "offer not create"})
		log.Println(http.StatusInternalServerError, gin.H{"error": "offer not create"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Offer create"})

}
