package db

import (
	"TechmasterProject/03/global"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectToDB() {
	dsn := "host=localhost user=root password=123techmaster dbname=techmaster port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		SkipDefaultTransaction: false,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("Successfully connected to database")

	global.DB = db
}
