package takeoffer

import (
	"be/conf"
	offerstr "be/offers/offerStr"
	offertakenstr "be/offers/offerTakenStr"
	userstr "be/userStr"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func TakeOffer(c *gin.Context) {

	id, exists := c.Get("id")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Id not found in context"})
		return
	}
	id = int(id.(float64))

	var user userstr.User
	user.Id = id.(int)

	errUser := conf.DB.Where("id = ?", user.Id).First(&user).Error
	if errUser != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	var offer offerstr.Offer
	if err := c.ShouldBindJSON(&offer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	errFind := conf.DB.Where("id = ? AND is_give = ?", offer.Id, false).First(&offer).Error
	if errFind != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Offer not found or already taken"})
		return
	}

	offer.IsGive = true
	user.DevWallet += offer.Ammount

	err := conf.DB.Save(&offer).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not take offer"})
		return
	}

	err = conf.DB.Save(&user).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not take offer"})
		return
	}

	var offerTaken offertakenstr.OfferTaken

	offerTaken.Id_l = user.Id
	offerTaken.Id_c = offer.Id_c
	offerTaken.Id_o = offer.Id
	offerTaken.NewAmmount = user.DevWallet * (1 + float64(offer.Percent)/100)
	offerTaken.DateGive = time.Now().Format("2006-01-02")

	err = conf.DB.Save(&offerTaken).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not take offer"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Offer taken"})
}
