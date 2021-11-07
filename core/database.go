package core

import (
	"github.com/jinzhu/gorm"
)

func ConnectDatabase(config *Config) (*gorm.DB, error) {
	return gorm.Open("mysql", config.MYSQL_CONNECTION)
}

func CloseDatabase(connection *gorm.DB) {
	sqldb := connection.DB()
	sqldb.Close()
}
