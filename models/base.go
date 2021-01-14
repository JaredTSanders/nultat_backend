package models

import (
	"fmt"
	"os"

	"github.com/gomodule/redigo/redis"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/joho/godotenv"
)

var db *gorm.DB

var Cache redis.Conn

func init() {

	e := godotenv.Load(".env")
	if e != nil {
		fmt.Print(e)
	}

	username := os.Getenv("db_user")
	password := os.Getenv("db_pass")
	dbName := os.Getenv("db_name")
	dbHost := os.Getenv("db_host")
	dbPort := os.Getenv("db_port")
	dbUri := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=require password=%s", dbHost, dbPort, username, dbName, password)
	fmt.Println(dbUri)

	conn, err := gorm.Open("postgres", dbUri)
	if err != nil {
		fmt.Print(err)
	}

	db = conn
	db.Debug().AutoMigrate(&Account{}, &AccountTypes{}, &ArkPCServer{}, &MinecraftBRServer{}, &Arma2Server{}, &AssettoCCServer{}, &AvailServer{})

	initCache()

}

func GetDB() *gorm.DB {
	return db
}

func initCache() {
	conn, err := redis.DialURL("redis://10.10.30.30:6379")
	if err != nil {
		fmt.Print(err)
	}
	// conn.Do("AUTH", "S@tUrn")
	Cache = conn
}
