package router

import (
	"fmt"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"net/http"
	"xxvote/app/logic"
	"xxvote/app/model"
	"xxvote/app/tools"

	_ "xxvote/docs"
)

func New() {
	r := gin.Default()
	r.LoadHTMLGlob("app/view/*")
	//所有的路径放在这里

	r.GET("/redis", func(ctx *gin.Context) {
		s := model.GetVoteCache(ctx, 3)
		fmt.Printf("redis:%+v\n", s)
	})

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	index := r.Group("")
	//index.Use(CheckUser)

	{
		//vote
		index.GET("/index", logic.Index) //静态页面

		index.POST("/vote/add", logic.AddVote)
		index.POST("/vote/update", logic.UpdateVote)
		index.POST("/vote/delete", logic.DelVote)

		index.GET("/result", logic.ResultInfo)
		index.GET("/result/info", logic.ResultVote)

	}

	//restful风格接口
	{
		//读
		index.GET("/votes", logic.GetVotes)
		index.GET("/vote", logic.GetVoteInfo)

		index.POST("/vote", logic.AddVote)
		index.PUT("/vote", logic.UpdateVote)
		index.DELETE("/vote", logic.DelVote)

		index.GET("/vote/result", logic.ResultVote)

		index.POST("/do_vote", logic.Dovote)
	}

	r.GET("/", logic.Index)

	{
		//login
		r.GET("/login", logic.GetLogin)
		r.POST("/login", logic.DoLogin)
		r.GET("/logout", logic.Logout)
		r.GET("/captcha", logic.GetCaptcha)

		r.POST("/captcha/verify", func(ctx *gin.Context) {
			var param tools.CaptchaData
			if err := ctx.ShouldBind(&param); err != nil {
				ctx.JSON(http.StatusOK, tools.ParamErr)
				return
			}

			fmt.Printf("参数为:%+v", param)
			if !tools.CaptchaVerify(param) {
				ctx.JSON(http.StatusOK,
					tools.ECode{
						Code:    10008,
						Message: "验证失败",
					})
				return
			}
			ctx.JSON(http.StatusOK, tools.OK)
		})

		//user
		r.POST("/user/create", logic.CreateUser)
	}
	if err := r.Run(":8080"); err != nil {
		panic("gin启动失败")
	}
}

func CheckUser(ctx *gin.Context) {
	var name string
	var id int64 //TODO  存在一个bu g
	values := model.GetSession(ctx)

	if v, ok := values["name"]; ok {
		name = v.(string)
	}
	if v, ok := values["id"]; ok {
		id = v.(int64)
	}

	if name == "" || id <= 0 {
		ctx.JSON(http.StatusUnauthorized, tools.NotLogin)
		ctx.Abort()

	}
	ctx.Next()
}
