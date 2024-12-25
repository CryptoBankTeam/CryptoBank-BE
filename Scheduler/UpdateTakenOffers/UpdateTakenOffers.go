package updatetakenoffers

import (
	"be/conf"
	offertakenstr "be/offers/offerTakenStr"
	"log"
	"time"
)

func UpdateTakenOffers() {
	var offersTaken []offertakenstr.OfferTaken

	err := conf.DB.Preload("Offer").Find(&offersTaken).Error
	if err != nil {
		log.Println("Error fetching offers:", err)
	}

	for _, offerTaken := range offersTaken {
		var dateStart time.Time
		var dateEnd time.Time
		var dateNow = time.Now()

		dateStart, _ = time.Parse("2006-01-02", offerTaken.DateGive)
		dateEnd = dateStart.AddDate(0, 0, offerTaken.Offer.Loan)

		if dateEnd.Before(dateNow) || dateEnd.Equal(dateNow) {
			offerTaken.IsExpired = true
		} else {
			offerTaken.NewAmmount = offerTaken.NewAmmount * (1 + float64(offerTaken.Offer.Percent)/100)
		}
		errSave := conf.DB.Save(&offerTaken).Error
		if errSave != nil {
			log.Println("Error saving offer:", err)
		}
	}
}
