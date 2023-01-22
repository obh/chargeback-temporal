package main

import (
	"chargebackapp/delivery"
	"chargebackapp/models"
	"chargebackapp/temporal"
	"fmt"
	"log"

	"github.com/labstack/echo/v4"

	"go.temporal.io/sdk/client"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	db.AutoMigrate(&models.ChargebackRequest{})

	if err != nil {
		panic("failed to connect database")
	}

	c, err := temporal.NewClient(client.Options{})
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	defer c.Close()

	e := echo.New()
	delivery.AddChargebackHandler(e, db, c)
	fmt.Println("running server...")
	e.Logger.Fatal(e.Start(":1323"))
}
