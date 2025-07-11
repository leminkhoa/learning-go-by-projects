package database

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var (
	DBConn *gorm.DB
)

func Connect() {
	d, err := gorm.Open("postgres", "host=localhost port=5432 user=leminkhoa dbname=crm password=123456 sslmode=disable")
	if err != nil {
		panic(err)
	}
	DBConn = d
}

func GetDB() *gorm.DB {
	return DBConn
}
