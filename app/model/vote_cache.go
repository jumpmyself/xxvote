package model

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

func GetVoteCache(c context.Context, id int64) VoteWithOpt {
	var ret VoteWithOpt
	key := fmt.Sprintf("key_vote_%d", id)
	voteStr, err := Rdb.Get(c, key).Result()
	if err == nil || len(voteStr) > 0 {
		//存在数据
		_ = json.Unmarshal([]byte(voteStr), &ret)
		return ret
	}

	fmt.Printf("err:%s", err.Error())
	ret = GetVote(id)
	if ret.Vote.ID > 0 {
		//更新缓存
		s, _ := json.Marshal(ret)
		err1 := Rdb.Set(c, key, s, 600*time.Second).Err()
		if err1 != nil {
			fmt.Printf("err1:%s", err.Error())
		}
	}

	return ret
}

func GetVoteHistoryV1(c context.Context, userId, voteId int64) []VoteOptUser {
	ret := make([]VoteOptUser, 0)

	//先查下redis
	K := fmt.Sprintf("vote-user-%d-%d", userId, voteId)
	str, err := Rdb.Get(c, K).Result()
	if len(str) == 0 && err == nil {
		fmt.Printf("不回溯数据库\n")
		json.Unmarshal([]byte(str), &ret)
		return ret
	}
	fmt.Printf("回溯数据库\n")
	//回溯数据库
	if err := Conn.Table("vote_opt_user").Where("vote_id = ? and user_id = ?", voteId, userId).First(&ret).Error; err != nil {
		fmt.Printf("err:%s", err.Error())
		return ret
	}

	//if len(ret) > 0 {
	retStr, _ := json.Marshal(ret)
	err = Rdb.Set(c, K, retStr, 3600*time.Second).Err()
	if err != nil {
		fmt.Printf("err1:%s\n", err.Error())
	}
	return ret
}
