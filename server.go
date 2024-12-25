package main

import (
	updatetakenoffers "be/Scheduler/UpdateTakenOffers"
	"be/auth/auth"
	"be/auth/middleware"
	"be/auth/refresh"
	"be/auth/reg"
	"be/conf"
	closeoffer "be/offers/closeOffer"
	createoffer "be/offers/createOffer"
	getalloffers "be/offers/getAllOffers"
	getmyoffers "be/offers/getMyOffers"
	puboffer "be/offers/pubOffer"
	ratingoffer "be/offers/ratingOffer"
	takeoffer "be/offers/takeOffer"
	unpuboffer "be/offers/unpubOffer"
	getmycredits "be/profile/getMyCredits"
	getprofile "be/profile/getProfile"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-co-op/gocron"
	_ "github.com/lib/pq"
)

func main() {

	conf.InitDB()
	conf.MigrationTables()

	scheduler := gocron.NewScheduler(time.UTC)
	scheduler.Every(1).Day().Do(updatetakenoffers.UpdateTakenOffers)
	scheduler.StartAsync()

	serv := gin.Default()

	serv.Use(gin.Logger())
	serv.Use(conf.Cors())

	serv.POST("/auth", auth.Auth)
	serv.POST("/refresh", refresh.RefreshToken)
	serv.POST("/reg", reg.Registration)

	protectedServ := serv.Group("/protected")
	protectedServ.Use(middleware.Middleware())

	protectedServ.GET("/getAllOffers", getalloffers.GetAllOffers)
	protectedServ.GET("/getProfile", getprofile.GetProfile)
	protectedServ.GET("/getMyCredits", getmycredits.GetMyCredits)
	protectedServ.GET("/getMyOffers", getmyoffers.GetMyOffers)

	protectedServ.POST("/takeOffer", takeoffer.TakeOffer)
	protectedServ.POST("/rateOffer", ratingoffer.RatingOffer)
	protectedServ.POST("/closeOffer", closeoffer.CloseOffer)
	protectedServ.POST("/createOffer", createoffer.CreateOffer)

	protectedServ.POST("/pubOffer", puboffer.PubOffer)
	protectedServ.POST("/unpubOffer", unpuboffer.UnPubOffer)

	log.Println("Server starting at :8080")
	log.Fatal(serv.Run(":8080"))

}
