package main

import (
	"be/auth/auth"
	"be/auth/middleware"
	"be/auth/refresh"
	"be/auth/reg"
	"be/conf"
	getallloans "be/offers/getAllLoans"
	getmyloans "be/profile/getMyLoans"
	getprofile "be/profile/getProfile"
	"log"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {

	conf.InitDB()
	conf.MigrationTables()

	serv := gin.Default()

	serv.Use(gin.Logger())
	serv.Use(conf.Cors())

	serv.POST("/auth", auth.Auth)
	serv.POST("/refresh", refresh.RefreshToken)
	serv.POST("/reg", reg.Registration)

	protectedServ := serv.Group("/protected")
	protectedServ.Use(middleware.Middleware())

	protectedServ.GET("/getAllOffers", getallloans.GetAllLoans)
	protectedServ.GET("/getProfile", getprofile.GetProfile)
	protectedServ.GET("/getMyCredits", getmyloans.GetMyAcceptedLoans)
	protectedServ.GET("/getMyOffers", getmyloans.GetMyCreatedLoans)

	protectedServ.POST("/setWallet", reg.SetWallet)

	log.Println("Server starting at :8080")
	log.Fatal(serv.Run(":8080"))

}
