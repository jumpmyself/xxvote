package logic

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
	"xxvote/app/model"
	"xxvote/app/tools"
)

func Index(ctx *gin.Context) {
	ret := model.GetVotes()
	ctx.HTML(http.StatusOK, "index.tmpl", gin.H{"vote": ret})
}

func GetVotes(ctx *gin.Context) {
	ret := model.GetVotes()
	ctx.JSON(http.StatusOK, tools.ECode{

		Data: ret,
	})
}

// GetVoteInfo godoc
// @Summary         获取投票信息
// @Description     获取投票信息
// @Tags            vote
// @Accept          json
// @Produce         json
// @Param        id    query     int  true  "vote ID"
// @Success         200   {object}    tools.ECode
// @Router          /vote   [get]
func GetVoteInfo(ctx *gin.Context) {
	var id int64
	idStr := ctx.Query("id")
	id, _ = strconv.ParseInt(idStr, 10, 64)
	ret := model.GetVote(id)
	//log.Printf("[print]ret:%+v\n", ret)
	//log.Panicf("[panic]ret:%+v\n", ret)
	//log.Fatalf("[fatal]ret:%+v\n", ret)

	//logrus.Errorf("[error]ret:%+v", ret)
	tools.Logger.Errorf("[error]ret:%+v", ret)
	if ret.Vote.ID <= 0 {
		ctx.JSON(http.StatusNotFound, tools.ECode{})
		return
	}

	ctx.JSON(http.StatusOK, tools.ECode{

		Data: ret,
	})
}

func Dovote(ctx *gin.Context) {
	userIDStr, _ := ctx.Cookie("id")
	voteIdStr, _ := ctx.GetPostForm("vote_id")
	optStr, _ := ctx.GetPostFormArray("opt[]")

	userID, _ := strconv.ParseInt(userIDStr, 10, 64)
	voteId, _ := strconv.ParseInt(voteIdStr, 10, 64)

	//第二种：前置查询，一般使用这种方式
	old := model.GetVoteHistory(userID, voteId)
	if len(old) >= 1 {
		ctx.JSON(http.StatusOK, tools.ECode{
			Code:    10010,
			Message: "您已投过票",
		})
	}
	opt := make([]int64, 0)
	for _, v := range optStr {
		optId, _ := strconv.ParseInt(v, 10, 64)
		opt = append(opt, optId)
	}

	model.DoVote(userID, voteId, opt)
	model.Rdb.Set(ctx, fmt.Sprintf("ket_vote_%d", voteId), "", 0)
	ctx.JSON(http.StatusOK, tools.ECode{
		Message: "投票完成",
	})

}

func checkXYZ(ctx *gin.Context) bool {
	//拿到ip+ua
	ip := ctx.ClientIP()
	ua := ctx.GetHeader("user-agent")
	fmt.Printf("ip:%s\n ua:%s\n", ip, ua)

	//转为md5
	hash := md5.New()
	hash.Write([]byte(ip + ua))
	hashBytes := hash.Sum(nil)
	hashString := hex.EncodeToString(hashBytes)

	flag, _ := model.Rdb.Get(ctx, "ban-"+hashString).Bool()
	if flag {
		return false
	}

	i, _ := model.Rdb.Get(ctx, "xyz-"+hashString).Int()
	if i > 5 {
		model.Rdb.SetEx(ctx, "ban-"+hashString, true, 30*time.Second)
		return false
	}

	model.Rdb.Incr(ctx, "xyz-"+hashString)                   //自增加一
	model.Rdb.Expire(ctx, "xyz-"+hashString, 50*time.Second) //过期时间

	return true
}

func GetCaptcha(ctx *gin.Context) {
	if !checkXYZ(ctx) {
		ctx.JSON(http.StatusOK, tools.ECode{
			Code:    10005,
			Message: "您的手速实在是太快了",
		})
		return
	}
	captcha, err := tools.CaptchaGenerate()
	if err != nil {
		ctx.JSON(http.StatusOK,
			tools.ECode{
				Code:    10005,
				Message: err.Error(),
			})
		return
	}
	ctx.JSON(http.StatusOK, tools.ECode{
		Data: captcha,
	})
}
