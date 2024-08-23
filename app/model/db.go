package model

import (
	"fmt"
	"github.com/rbcervilla/redisstore/v9"
	"github.com/redis/go-redis/v9"
	"golang.org/x/net/context"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Conn 所有的数据库操作放在这里
var Conn *gorm.DB
var Rdb *redis.Client

func NewMysql() {
	my := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local", "root", "431630", "127.0.0.1:3306", "xxvote")
	conn, err := gorm.Open(mysql.Open(my), &gorm.Config{})
	if err != nil {
		fmt.Printf("err:%s\n", err)
		panic(err)
	}
	Conn = conn
}

func NewRdb() {
	rdb := redis.NewClient(&redis.Options{
		//Network: "",
		Addr: "127.0.0.1:6379",
		//ClientName:            "",
		//Dialer:                nil,
		//OnConnect:             nil,
		//Protocol:              0,
		//Username:              "",
		Password: "", //no password  set
		//CredentialsProvider:   nil,
		DB: 0, //use  default  DB
		//MaxRetries:            0,
		//MinRetryBackoff:       0,
		//MaxRetryBackoff:       0,
		//DialTimeout:           0,
		//ReadTimeout:           0,
		//WriteTimeout:          0,
		//ContextTimeoutEnabled: false,
		//PoolFIFO:              false,
		//PoolSize:              0,
		//PoolTimeout:           0,
		//MinIdleConns:          0,
		//MaxIdleConns:          0,
		//MaxActiveConns:        0,
		//ConnMaxIdleTime:       0,
		//ConnMaxLifetime:       0,
		//TLSConfig:             nil,
		//Limiter:               nil,
		//DisableIndentity:      false,
		//IdentitySuffix:        "",
	})
	Rdb = rdb

	//初始化session
	store, _ = redisstore.NewRedisStore(context.TODO(), Rdb)

	return
}
func Close() {
	db, _ := Conn.DB()
	_ = db.Close()
	_ = Rdb.Close()
}
