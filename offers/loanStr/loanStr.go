package loanStr

import (
	userstr "be/userStr"
	"strconv"
)

type Loan struct {
	Id         int64   `json:"id" gorm:"primaryKey;autoIncrement"`
	LenderId   string  `json:"lender_id" gorm:"type:varchar(42);not null"`
	BorrowerId string  `json:"borrower_id" gorm:"type:varchar(42);"`
	Amount     float64 `json:"amount" gorm:"type:numeric(78,0);not null"`
	Interest   int     `json:"interest" gorm:"type:integer;not null"`
	Collateral float64 `json:"collateral" gorm:"type:numeric(78,0);not null"`
	DueDate    int64   `json:"due_date" gorm:"type:integer;"`
	Status     uint8   `json:"status" gorm:"type:integer;not null"`
	Duration   int64   `json:"duration" gorm:"type:integer;not null"`

	Lender   userstr.User `json:"lender" gorm:"foreignKey:LenderId;references:Id;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Borrower userstr.User `json:"borrower" gorm:"foreignKey:BorrowerId;references:Id;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

// Маппинг статусов
var StatusMap = map[uint8]string{
	0: "Created",
	1: "Accepted",
	2: "Repaid",
	3: "Overdue",
	4: "Closed",
}

// Для сериализации в JSON строкой
func (l Loan) StatusString() string {
	if s, ok := StatusMap[l.Status]; ok {
		return s
	}
	return strconv.Itoa(int(l.Status))
}

type LoanResponse struct {
	Loan
	Status string `json:"status"`
}
