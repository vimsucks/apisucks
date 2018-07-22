package models

import (
	"log"
	. "github.com/vimsucks/apisucks/config"
	"github.com/go-redis/redis"
	"github.com/jmoiron/sqlx"
	_ "github.com/go-sql-driver/mysql"
)

var Redis *redis.Client
var DB *sqlx.DB

func init() {
	var err error
	DB, err = sqlx.Open("mysql", Conf.DatabaseUrl)
	if err != nil {
		log.Fatal(err)
	}
	err = DB.Ping()

	Redis = redis.NewClient(&redis.Options{
		Addr:     Conf.RedisAddr,
		Password: Conf.RedisPassword,
		DB:       Conf.RedisDB,
	})
	_, err = Redis.Ping().Result()
	if err != nil {
		log.Fatal(err)
	}
}
