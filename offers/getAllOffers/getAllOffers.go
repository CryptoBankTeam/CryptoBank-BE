package getalloffers

import (
	"be/conf"
	offerstr "be/offers/offerStr"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetAllOffers(c *gin.Context) {

	offers := []offerstr.Offer{}
	err := conf.DB.
		Where("is_give", false).
		Preload("Creditor").
		Find(&offers).
		Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not get offers"})
		log.Println(http.StatusInternalServerError, gin.H{"error": "Could not get offers"})
		return
	}

	for i := range offers {
		offers[i].Creditor.Id = 0
		offers[i].Creditor.Password = ""

	}

	log.Println(offers)

	c.JSON(http.StatusOK, offers)
}
