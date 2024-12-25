package offerstr

import userstr "be/userStr"

type Offer struct {
	Id      int     `json:"id" gorm:"primaryKey;autoIncrement"`
	Ammount float64 `json:"ammount" gorm:"type:numeric"`
	Percent float32 `json:"percent" gorm:"type:real"`
	Loan    int     `json:"loan" gorm:"type:integer"`
	Id_c    int     `json:"id_c" gorm:"type:integer"`
	IsGive  bool    `json:"is_give" gorm:"type:boolean"`
	IsPub   bool    `json:"is_pub" gorm:"type:boolean;default:true"`

	Creditor userstr.User `json:"creditor" gorm:"foreignKey:Id_c;references:Id;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}
