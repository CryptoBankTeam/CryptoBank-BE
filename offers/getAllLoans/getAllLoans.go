package getAllLoans

import (
	"be/conf"
	loanstr "be/offers/loanStr"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetAllLoans(c *gin.Context) {
	loans := []loanstr.Loan{}
	err := conf.DB.
		Where("status = ?", 0). // 0 = Created, только открытые займы
		Preload("Lender").
		Find(&loans).
		Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not get loans"})
		log.Println(http.StatusInternalServerError, gin.H{"error": "Could not get loans"})
		return
	}

	var resp []loanstr.LoanResponse
	for _, loan := range loans {
		loan.Lender.Password = ""
		resp = append(resp, loanstr.LoanResponse{
			Loan:   loan,
			Status: loan.StatusString(),
		})
	}

	log.Println(resp)

	c.JSON(http.StatusOK, resp)
}
