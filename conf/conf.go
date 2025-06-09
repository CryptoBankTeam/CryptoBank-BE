package conf

import (
	"be/offers/loanStr"
	userstr "be/userStr"
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, PATCH")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}

func InitDB() {
	var err error
	dsn := "host=localhost user=postgres password=5121508 dbname=cr_bank_db port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	DB = db
	if err != nil {
		log.Fatal(err)
	}
}

func MigrationTables() {
	if err := DB.AutoMigrate(&userstr.User{}); err != nil {
		log.Fatalf("Failed to migrate user table: %v", err)
	}

	if err := DB.AutoMigrate(&loanStr.Loan{}); err != nil {
		log.Fatalf("Failed to migrate Loan table: %v", err)
	}
}
