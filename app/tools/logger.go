package tools

import (
	"github.com/sirupsen/logrus"
	"io"
	"os"
)

var Logger *logrus.Entry

func NewLogger() {
	LogStore := logrus.New()
	LogStore.SetLevel(logrus.DebugLevel)

	//同时写到多个输出
	w1 := os.Stdout //写到控制台
	w2, _ := os.OpenFile("./vote.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	LogStore.SetOutput(io.MultiWriter(w1, w2)) //io.MultiWriter 返回一个io.Writer 对象

	LogStore.SetFormatter(&logrus.JSONFormatter{})
	Logger = LogStore.WithFields(logrus.Fields{
		"name": "香香编程喵喵喵",
		"app":  "voteV2",
	})

	//LogStore.AddHook( )  出现非常严重的问题，直接邮箱或者微信警报；日志分割，塞入一些自定义的字段
}
