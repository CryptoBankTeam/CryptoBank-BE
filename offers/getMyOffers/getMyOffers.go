package getmyoffers

import (
	"be/conf"
	offerstr "be/offers/offerStr"
	userstr "be/userStr"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetMyOffers(c *gin.Context) {
	id, exists := c.Get("id")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Id not found in context"})
		return
	}
	id = int(id.(float64))

	var user userstr.User
	user.Id = id.(int)

	log.Println(user.Id)

	var offers []offerstr.Offer
	errOffer := conf.DB.Where("id_c =?", user.Id).Find(&offers).Error
	if errOffer != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "offers not found"})
		log.Println(http.StatusNotFound, gin.H{"error": "offers not found"}, errOffer)
		return
	}

	c.JSON(http.StatusOK, offers)
}
