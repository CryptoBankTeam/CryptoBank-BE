package offertakenstr

import (
	offerstr "be/offers/offerStr"
	userstr "be/userStr"
)

type OfferTaken struct {
	Id         int     `json:"id" gorm:"primaryKey;autoIncrement"`
	Id_l       int     `json:"id_l" gorm:"type:bigint"`
	Id_c       int     `json:"id_c" gorm:"type:bigint"`
	Id_o       int     `json:"id_o" gorm:"type:bigint"`
	DateGive   string  `json:"date_give" gorm:"type:date"`
	DateClosed *string `json:"date_closed" gorm:"type:date"`
	NewAmmount float64 `json:"new_ammount" gorm:"type:numeric"`
	IsClose    bool    `json:"is_close" gorm:"type:boolean"`
	IsExpired  bool    `json:"is_expired" gorm:"type:boolean"`
	DayPass    int     `json:"day_pass" gorm:"type:integer"`

	Offer  offerstr.Offer `json:"offer" gorm:"foreignKey:Id_o;references:Id;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Loaner userstr.User   `json:"user" gorm:"foreignKey:Id_l;references:Id;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Credit userstr.User   `json:"credit" gorm:"foreignKey:Id_c;references:Id;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}
