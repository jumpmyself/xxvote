package logic

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
	"xxvote/app/model"
	"xxvote/app/tools"
)

type ResultData struct {
	Title string
	Count int64
	Opt   []*ResultVoteOpt
}

type ResultVoteOpt struct {
	Name  string
	Count int64
}

func AddVote(ctx *gin.Context) {
	idStr := ctx.Query("title")
	optStr, _ := ctx.GetPostFormArray("opt_name[]")
	//构造结构体
	vote := model.Vote{
		Title:       idStr,
		Type:        strconv.Itoa(0),
		Status:      strconv.Itoa(0),
		CreatedTime: time.Now(),
		UpdatedTime: time.Now(), // 设置更新时间为当前时间
	}
	if vote.Title == "" {
		ctx.JSON(http.StatusBadRequest, tools.ParamErr)
		return
	}

	oldVote := model.GetVoteByName(idStr)
	if oldVote.ID > 0 {
		//ctx.JSON(http.StatusCreated,tools.OK)

		ctx.JSON(http.StatusOK, tools.ECode{
			Code:    10006,
			Message: "投票已存在",
		})
		return
	}
	opt := make([]model.VoteOpt, 0)
	for _, v := range optStr {
		opt = append(opt, model.VoteOpt{
			Name:        v,
			CreatedTime: time.Now(),
		})
	}

	if err := model.AddVote(vote, opt); err != nil {
		ctx.JSON(http.StatusOK, tools.ECode{
			Code:    10006,
			Message: err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusCreated, tools.OK)
	return
}

func UpdateVote(ctx *gin.Context) {

}

func DelVote(ctx *gin.Context) {
	var id int64
	idStr := ctx.Query("id")
	id, _ = strconv.ParseInt(idStr, 10, 64)

	vote := model.GetVote(id)
	if vote.Vote.ID <= 0 {
		ctx.JSON(http.StatusNoContent, tools.OK)
		return
	}

	if err := model.DelVote(id); err != true {
		ctx.JSON(http.StatusOK, tools.ECode{
			Code:    10006,
			Message: "删除失败",
		})
		return
	}

	ctx.JSON(http.StatusNoContent, tools.OK)
	return
}

func ResultInfo(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "result.tmpl", nil)
}

// ResultData 新定义返回结构

func ResultVote(ctx *gin.Context) {
	var id int64
	idStr := ctx.Query("id")
	id, _ = strconv.ParseInt(idStr, 10, 64)
	ret := model.GetVote(id)
	data := ResultData{
		Title: ret.Vote.Title,
	}

	for _, v := range ret.Opt {
		data.Count = data.Count + v.Count
		tmp := ResultVoteOpt{
			Name:  v.Name,
			Count: v.Count,
		}
		data.Opt = append(data.Opt, &tmp)
	}
	ctx.JSON(http.StatusOK, tools.ECode{
		Data: data,
	})
}
