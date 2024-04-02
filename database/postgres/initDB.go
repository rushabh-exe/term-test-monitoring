package postgres

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func PostgresInitializer() {
	var err error
	dsn := "host=localhost user=postgres password=mysecretpassword dbname=postgres port=5431"
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("Postgres not connected")
	} else {
		fmt.Println("Prostgres Connected")
	}
}
