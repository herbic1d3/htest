package model

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"htest/config"
)

var (
	DBConn *gorm.DB
)

func GormInit() error {
	conf := &config.DBConfig{}
	err := conf.Read()
	DBConn, err = gorm.Open("postgres",
		fmt.Sprintf("host=localhost user=%s dbname=%s sslmode=disable password=%s",
			conf.DBUser, conf.DBName, conf.DBPass))

	if err != nil {
		return err
	} else {
		// drop table for easy tests
		DBConn.DropTable(&User{})
		DBConn.AutoMigrate(&User{})
		user := User{
			Login:      "admin",
			Pass:       "1000000",
			WorkNumber: 1,
		}

		if err = DBConn.First(&user).Error; err != nil {
			DBConn.Create(&user)
		}
	}

	return nil
}

func GormClose() error {
	if DBConn != nil {
		return DBConn.Close()
	}
	return nil
}
