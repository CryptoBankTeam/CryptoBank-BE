package getmycredits

import (
	"be/conf"
	offerstr "be/offers/offerStr"
	offertakenstr "be/offers/offerTakenStr"
	userstr "be/userStr"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetMyCredits(c *gin.Context) {
	id, exists := c.Get("id")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Id not found in context"})
		log.Println(http.StatusInternalServerError, gin.H{"error": "Id not found in context"})
		return
	}
	id = int(id.(float64))

	var user userstr.User
	user.Id = id.(int)

	var offersTaken []offertakenstr.OfferTaken
	errMyCredits := conf.DB.Preload("Offer.Creditor").Preload("Offer").Where("id_l = ?", user.Id).Find(&offersTaken).Error
	if errMyCredits != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		log.Println(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	var offers []offerstr.Offer
	for i := range offersTaken {
		offers = append(offers, offersTaken[i].Offer)
	}

	c.JSON(http.StatusOK, offersTaken)
}
