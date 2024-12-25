package userstr

type User struct {
	Id          int     `json:"id" gorm:"primaryKey;autoIncrement"`
	Username    string  `json:"username" gorm:"type:varchar(20);unique"`
	Password    string  `json:"password" gorm:"type:varchar(100)"`
	Rating      float32 `json:"rating" gorm:"type:real"`
	CountRating int     `json:"count_rating" gorm:"type:integer"`
	DevWallet   float64 `json:"dev_wallet" gorm:"type:numeric"`
}
