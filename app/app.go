package app

import (
	"xxvote/app/model"
	"xxvote/app/router"
	"xxvote/app/schedule"
	"xxvote/app/tools"
)

// Start 这只是一个启动器方法
func Start() {
	model.NewMysql()
	model.NewRdb()
	defer func() {
		model.Close()
	}()

	schedule.Start()

	tools.NewLogger()
	router.New()
}
