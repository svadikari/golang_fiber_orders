package database

import (
	"os"

	orderModels "github.com/svadikari/golang_fiber_orders/src/orders/models"
	productModels "github.com/svadikari/golang_fiber_orders/src/products/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DbInstance struct {
	Db *gorm.DB
}

var Database DbInstance

func ConnectDB() error {
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		dsn = "host=localhost user=svadikari password=yourpassword dbname=orders-db port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	}
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}

	// Migrate the schema
	db.AutoMigrate(&productModels.Product{}, &orderModels.Order{}, &orderModels.OrderItem{})

	Database = DbInstance{
		Db: db,
	}
	return nil
}
