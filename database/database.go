package database

import (
	"fmt"
	"interncase/envs"
	"interncase/models"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	DB        *gorm.DB
	SecretKey = []byte("secret")
)

// Bağlantı fonksiyonu
func Connect() {

	//connect
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", envs.DBuser, envs.DBpass, envs.DBhost, envs.DBport, envs.DBname)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalln(err)
	}

	//ping
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalln(err)
	}
	err = sqlDB.Ping()
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("Database connected")

	//migrate
	if err := db.AutoMigrate(&models.Plan{}, &models.Student{}, &models.StatuData{}); err != nil {
		log.Fatalln(err)
	}

	DB = db

}
