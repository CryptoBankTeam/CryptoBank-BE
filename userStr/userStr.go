package userstr

type User struct {
	Id             int     `json:"id" gorm:"primaryKey;autoIncrement"`
	Username       string  `json:"username" gorm:"type:varchar(20);unique"`
	Password       string  `json:"password" gorm:"type:varchar(100)"`
	Rating         float32 `json:"rating" gorm:"type:real"`
	AdressWallet   string  `json:"adress_wallet" gorm:"varchar(42);unique"`
	CleanLoans     int     `json:"clean_loans" gorm:"type:integer;"`
	OverdueLoans   int     `json:"overdue_loans" gorm:"type:integer;"`
	OffersAccepted int     `json:"offers_accepted" gorm:"type:integer;"`
}
