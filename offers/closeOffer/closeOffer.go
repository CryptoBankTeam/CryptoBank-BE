package closeoffer

import (
	"be/conf"
	offerstr "be/offers/offerStr"
	offertakenstr "be/offers/offerTakenStr"
	userstr "be/userStr"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func CloseOffer(c *gin.Context) {
	id, exists := c.Get("id")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Id not found in context"})
		log.Println(http.StatusInternalServerError, gin.H{"error": "Id not found in context"})
		return
	}
	id = int(id.(float64))

	var user userstr.User
	user.Id = id.(int)

	errUser := conf.DB.Where("id=?", user.Id).First(&user).Error
	if errUser != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		log.Println(http.StatusNotFound, gin.H{"error": "User not found"})
	}

	log.Println("user:", &user)

	var offer offerstr.Offer
	if err := c.ShouldBindJSON(&offer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		log.Println(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	var offerTaken offertakenstr.OfferTaken

	errOfT := conf.DB.Where("id_o = ?", offer.Id).First(&offerTaken).Error
	if errOfT != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "OfferTaken not found"})
		log.Println(http.StatusNotFound, gin.H{"error": "OfferTaken not found"})
	}

	user.DevWallet -= offerTaken.NewAmmount
	offerTaken.IsClose = true
	dateClosed := time.Now().Format("2006-01-02")
	offerTaken.DateClosed = &dateClosed

	errSaveUser := conf.DB.Save(&user).Error
	if errSaveUser != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Cannot update user"})
		log.Println(http.StatusInternalServerError, gin.H{"error": "Cannot update user"})
	}

	errSaveOffer := conf.DB.Save(&offerTaken).Error
	if errSaveOffer != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Cannot update offer"})
		log.Println(http.StatusInternalServerError, gin.H{"error": "Cannot update offer"})
	}

	c.JSON(http.StatusOK, gin.H{"message": "Сделка закрыта успешно"})
}
