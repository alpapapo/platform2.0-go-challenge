package models

import (
	"database/sql"
	"fmt"
	c "github.com/GlobalWebIndex/platform2.0-go-challenge/common"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"os"
	"strings"
	"time"
)

var db *gorm.DB

func InitDB(projectEnv string) {
	e := c.LoadDotenv(projectEnv)
	if e != nil {
		fmt.Print(e)
	}

	username := os.Getenv("DATABASE_USER")
	password := os.Getenv("DATABASE_PASSWORD")
	dbName := os.Getenv("DATABASE_NAME")
	dbHost := os.Getenv("DATABASE_HOST")
	dbPort := os.Getenv("DATABASE_PORT")

	// create database if does not exist
	noDbUri := fmt.Sprintf("%s:%s@tcp(%s:%s)/", username, password, dbHost, dbPort)
	noDb, err := sql.Open(os.Getenv("DATABASE_TYPE"), noDbUri)
	defer noDb.Close()
	if err != nil {
		fmt.Println(err)
	} else {
		create_db := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", dbName)
		_,noDbErr := noDb.Exec(create_db)
		if noDbErr != nil {
			for i := 1; i<=10; i++ {
				if noDbErr != nil && strings.Contains(noDbErr.Error(), "connection refused") {
					fmt.Println("Connection Refused, wait 2 secs to reconnect..")
					time.Sleep(2 * time.Second)
					_,noDbErr = noDb.Exec(create_db)

				} else if noDbErr == nil {
					fmt.Println("Connected with Database!")
					break
				} else {
					fmt.Println(noDbErr)
					break
				}
			}
		}
	}

	dbUri := fmt.Sprintf("%s:%s@(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", username, password, dbHost, dbName )
	fmt.Println(dbUri)

	conn, err := gorm.Open(os.Getenv("DATABASE_TYPE"), dbUri)
	if err != nil {
		fmt.Println(err)
	}


	db = conn
	if os.Getenv("DATABASE_DEBUG") == "enabled" {
		db = db.Debug()
	}
	db.Debug().AutoMigrate(&User{}, &Asset{})
}

func GetDB() *gorm.DB {
	return db
}