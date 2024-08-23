package model

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rbcervilla/redisstore/v9"
)

// session存到本机上，然后将session——name 通过cookie传给前端
var store *redisstore.RedisStore
var sessionName = "session-name"

func GetSession(ctx *gin.Context) map[interface{}]interface{} {
	session, _ := store.Get(ctx.Request, sessionName)
	fmt.Printf("session:%+v\n", session.Values)
	return session.Values
}

func SetSession(ctx *gin.Context, name string, id int64) error {
	session, _ := store.Get(ctx.Request, sessionName)
	session.Values["name"] = name
	session.Values["id"] = id
	return session.Save(ctx.Request, ctx.Writer)
}

func FlushSession(ctx *gin.Context) error {
	session, _ := store.Get(ctx.Request, sessionName)
	fmt.Printf("session:%+v\n", session.Values)
	session.Values["name"] = ""
	session.Values["id"] = int64(0)
	return session.Save(ctx.Request, ctx.Writer)
}
