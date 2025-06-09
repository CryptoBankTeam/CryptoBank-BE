package getMyLoans

import (
	"be/conf"
	loanstr "be/offers/loanStr"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Займы, созданные пользователем
func GetMyCreatedLoans(c *gin.Context) {
	id, exists := c.Get("id")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Id not found in context"})
		log.Println(http.StatusInternalServerError, gin.H{"error": "Id not found in context"})
		return
	}
	userId := int(id.(float64))

	var loans []loanstr.Loan
	err := conf.DB.
		Where("lender_id = ?", userId).
		Preload("Borrower").
		Find(&loans).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Loans not found"})
		log.Println(http.StatusNotFound, gin.H{"error": "Loans not found"}, err)
		return
	}

	var resp []loanstr.LoanResponse
	for _, loan := range loans {
		loan.Borrower.Password = ""
		lonResp := loanstr.LoanResponse{
			Loan:   loan,
			Status: loan.StatusString(),
		}
		resp = append(resp, lonResp)
		log.Print("lonResp", lonResp)
	}

	log.Print("resp", resp)
	c.JSON(http.StatusOK, resp)
}

// Займы, принятые пользователем
func GetMyAcceptedLoans(c *gin.Context) {
	id, exists := c.Get("id")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Id not found in context"})
		log.Println(http.StatusInternalServerError, gin.H{"error": "Id not found in context"})
		return
	}
	userId := int(id.(float64))

	var loans []loanstr.Loan
	err := conf.DB.
		Where("borrower_id = ?", userId).
		Preload("Lender").
		Find(&loans).Error
	if err != nil {
		log.Println(http.StatusInternalServerError, gin.H{"error": "DB error"}, err)
		c.JSON(http.StatusOK, []loanstr.LoanResponse{})
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

	c.JSON(http.StatusOK, resp)
}
