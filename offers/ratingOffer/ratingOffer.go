package ratingoffer

import (
	"be/conf"
	offerstr "be/offers/offerStr"
	"math"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Rating struct {
	Id     int `json:"id"`
	Rating int `json:"rating"`
}

func RatingOffer(c *gin.Context) {
	var rating Rating
	if err := c.ShouldBindJSON(&rating); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	var offer offerstr.Offer
	errOffer := conf.DB.Preload("Creditor").Where("id = ?", rating.Id).First(&offer).Error
	if errOffer != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Offer not found"})
		return
	}

	var tempRating int
	tempRating = int(math.Round(float64(offer.Creditor.Rating) * float64(offer.Creditor.CountRating)))
	offer.Creditor.CountRating++
	tempRating += rating.Rating
	offer.Creditor.Rating = float32(tempRating) / float32(offer.Creditor.CountRating)

	userCreditor := offer.Creditor
	err := conf.DB.Save(&userCreditor).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not rate offer"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Offer rated successfully"})
}
