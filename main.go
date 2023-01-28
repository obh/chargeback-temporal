package main

import (
	"fmt"
	"log"

	"github.com/obh/chargebackapp/delivery"
	"github.com/obh/chargebackapp/models"
	"github.com/obh/chargebackapp/temporal"
	"github.com/obh/chargebackapp/utils"

	"github.com/labstack/echo/v4"

	"go.temporal.io/sdk/client"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	if res := db.Exec("PRAGMA foreign_keys = ON", nil); res.Error != nil {
		panic("failed to initiate foreign keys")
	}
	db.AutoMigrate(&models.Payment{})
	db.AutoMigrate(&models.Customer{})
	db.AutoMigrate(&models.Chargeback{})

	utils.InsertPayment(db)

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
