package models

import (
	"fmt"
	"os"
	"log"
	"errors"
	"time"
	"github.com/gomodule/redigo/redis"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/getsentry/sentry-go"
	beeline "github.com/honeycombio/beeline-go"
	"github.com/sirupsen/logrus"
	"github.com/joho/godotenv"

)

var db *gorm.DB

var Cache redis.Conn

func init() {

	initSentry()

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
		sentry.WithScope(func(scope *sentry.Scope){
        	scope.SetLevel(sentry.LevelFatal)
        	sentry.CaptureException(errors.New("Database connection failed, thrown in init function"))
        })
		fmt.Print(err)
	}

	db = conn
	db.Debug().AutoMigrate(&Account{}, &AccountTypes{}, &ArkPCServer{}, &MinecraftBRServer{}, &Arma2Server{}, &AssettoCCServer{}, &AvailServer{})

	initCache()

}


func initSentry() {
	/*
	   =================
	    Beehive initialization:
	   =================
	*/
	sentryErr := sentry.Init(sentry.ClientOptions{
		Dsn: "https://337fbcd4a3ac46f0b30ca550d188bfe1@o551311.ingest.sentry.io/5679680",
	})
	if sentryErr != nil {
		log.Fatalf("sentry.Init: %s", sentryErr)
	}

	defer sentry.Flush(2 * time.Second)

	/*
	   =================
	    Beehive initialization:
	   =================
	*/

	wk := "8da63485f3ab2c6983c14a2c46de1456" //os.Getenv("HONEYCOMB_WRITEKEY")
	if wk == "" {
		logrus.Error("got empty writekey from the environment. Please set HONEYCOMB_WRITEKEY")
	}
	beeline.Init(beeline.Config{
		WriteKey: wk,
		Dataset:  "nultat-api",
	})
	// addCommonLibhoneyFields()

}

func GetDB() *gorm.DB {
	return db
}

func initCache() {
	conn, err := redis.DialURL("redis://redis.nultat.net:6379")
	if err != nil {
		sentry.WithScope(func(scope *sentry.Scope){
        	scope.SetLevel(sentry.LevelFatal)
        	sentry.CaptureException(errors.New("Redis connection failed, thrown in initCache function"))
        })
		fmt.Print(err)
	}
	// conn.Do("AUTH", "S@tUrn")
	Cache = conn
}


