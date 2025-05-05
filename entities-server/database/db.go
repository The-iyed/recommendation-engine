package database

import (
	"entities-server/config"
	"entities-server/modules/product"
	"fmt"
	"log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() *gorm.DB {

    var err error
    cfg := config.LoadConfig()
    dsn := fmt.Sprintf(
        "host=%s user=%s password=%s dbname=%s port=%s",
        cfg.DB_HOST, cfg.DB_USER,
        cfg.DB_PASSWORD , cfg.DB_NAME,
        cfg.DB_PORT,
    )

    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatal("Failed to connect to database:", err)

    }
    
    log.Printf("Connectd to database successfully")
	err = db.AutoMigrate(&product.Product{})
	if err != nil {
		log.Fatal("failed to migrate database:", err)
	}

    return db
}

func GetDB() *gorm.DB {
    return DB
}

