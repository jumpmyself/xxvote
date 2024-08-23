package model

import (
	"fmt"
	"gorm.io/gorm"
	"sync"
	"time"
)

func GetVotes() []Vote {
	ret := make([]Vote, 0)
	if err := Conn.Table("vote").Find(&ret).Error; err != nil {
		fmt.Printf("err:%s", err.Error())
	}
	return ret
}

func GetVote(id int64) VoteWithOpt {
	var ret Vote
	if err := Conn.Table("vote").Where("id = ?", id).First(&ret).Error; err != nil {
		fmt.Printf("err:%s", err.Error())
	}

	opt := make([]VoteOpt, 0)
	if err := Conn.Table("vote_opt").Where("vote_id = ?", id).Find(&opt).Error; err != nil {
		fmt.Printf("err:%s", err.Error())
	}
	return VoteWithOpt{
		Vote: ret,
		Opt:  opt,
	}
}

// GetVoteV1 改造为原生SQL
func GetVoteV1(id int64) VoteWithOpt {
	var ret Vote
	err := Conn.Raw("select * from vote where id = ?", id).Scan(&ret).Error
	if err != nil {
		fmt.Printf("err:%s", err.Error())
	}

	opt := make([]VoteOpt, 0)
	err1 := Conn.Raw("select *from vote_opt where vote_id = ?", id).Scan(&opt).Error
	if err1 != nil {
		fmt.Printf("err:%s", err.Error())
	}
	return VoteWithOpt{
		Vote: ret,
		Opt:  opt,
	}

}

// GetVoteV2 改造为gorm预加载模式
func GetVoteV2(id int64) (*Vote, error) {
	var ret Vote
	//select * from opt where vote_id = 2;
	//select * from vote where id = 2;
	err := Conn.Preload("Opt").Table("vote").Where("id = ?", id).First(&ret).Error
	if err != nil {
		fmt.Printf("err:%s", err.Error())
		return &ret, err
	}
	return &ret, nil
}

// GetVoteV3 使用join链接
func GetVoteV3(id int64) (*VoteWithOpt, error) {
	var ret VoteWithOpt
	//select * from vote join vote_opt on vote.id = vote_opt.vote_id where vote.id = 2
	sql := "select vote.*,vote_opt.id as vid,vote_opt.name,vote_opt.count from vote join vote_opt on vote.id = vote_opt.vote_id where vote.id = ?"
	//ret1 := make(map[string]any)
	row, err := Conn.Raw(sql, id).Rows()

	if err != nil {
		return nil, err
	}
	for row.Next() {
		tmp := make(map[string]any)
		_ = Conn.ScanRows(row, &tmp)
		fmt.Printf("tmp:%+v\n", tmp)

		if v, ok := tmp["id"]; ok {
			ret.Vote.ID = v.(int64)
		}

		//将map先转为json再转为结构体，也可以写一个反射，直接实现
	}

	fmt.Printf("map:%+v\n", ret)
	return &ret, err
}

// GetVoteV4 使用第一种并发方式
func GetVoteV4(id int64) (*VoteWithOpt, error) {
	var ret Vote
	ch := make(chan struct{}, 2)
	go func() {
		err := Conn.Raw("select * from vote where id = ?", id).Scan(&ret).Error
		if err != nil {
			fmt.Printf("err:%s", err.Error())
		}
		ch <- struct{}{}
	}()

	opt := make([]VoteOpt, 0)
	go func() {
		err1 := Conn.Raw("select *from vote_opt where vote_id = ?", id).Scan(&opt).Error
		if err1 != nil {
			fmt.Printf("err:%s", err1.Error())
		}
		ch <- struct{}{}
	}()

	var i int
	for _ = range ch {
		i++
		if i >= 2 {
			break
		}
	}
	return &VoteWithOpt{
		Vote: ret,
		Opt:  opt,
	}, nil

}

// GetVoteV5 使用sync.waitgroup并发方式
func GetVoteV5(id int64) (*VoteWithOpt, error) {
	var ret Vote

	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		err := Conn.Raw("select * from vote where id = ?", id).Scan(&ret).Error
		if err != nil {
			fmt.Printf("err:%s", err.Error())
		}
	}()

	wg.Add(1)
	opt := make([]VoteOpt, 0)
	go func() {
		defer wg.Done()
		err1 := Conn.Raw("select *from vote_opt where vote_id = ?", id).Scan(&opt).Error
		if err1 != nil {
			fmt.Printf("err:%s", err1.Error())
		}
	}()

	wg.Wait()
	return &VoteWithOpt{
		Vote: ret,
		Opt:  opt,
	}, nil
}

// DoVote gorm里面最常用的事物方法
func DoVote(userId, voteId int64, optIDs []int64) bool {
	tx := Conn.Begin()
	var ret Vote
	if err := tx.Table("vote").Where("id = ?", voteId).First(&ret).Error; err != nil {
		fmt.Printf("err:%s", err.Error())
		tx.Rollback()
	}
	for _, value := range optIDs {
		if err := tx.Table("vote_opt").Where("id = ?", value).Update("count", gorm.Expr("count + ?", 1)).Error; err != nil {
			fmt.Printf("err:%s", err.Error())
			tx.Rollback()
		}
		//第一种方式，在事务中查询，增加了事务的逻辑，成本非常高
		var oldVoteUser VoteOptUser
		if err := tx.Table("vote_opt_user").Where("vote_id = ? and user_id = ?", voteId, userId).First(&oldVoteUser).Error; err != nil {
			fmt.Printf("err:%s", err.Error())
			tx.Rollback()
		}

		if oldVoteUser.ID > 0 {
			fmt.Printf("用户已投票")
			tx.Rollback()
		}
		user := VoteOptUser{
			VoteId:      voteId,
			UserId:      userId,
			VoteOptId:   value,
			CreatedTime: time.Now(),
		}

		if err := tx.Create(&user).Error; err != nil {
			fmt.Printf("err:%s", err.Error())
			tx.Rollback()
		}
	}
	tx.Commit()
	return true
}

func DoVoteV1(userId, voteId int64, optIDs []int64) bool {
	Conn.Exec("begin").
		Exec("select * from vote where id = ?", voteId).
		Exec("commit")
	return false
}

// DoVoteV2 匿名函数，用的最多，使用起来最方便
func DoVoteV2(userId, voteId int64, optIDs []int64) bool {
	if err := Conn.Transaction(func(tx *gorm.DB) error {
		var ret Vote
		if err := tx.Table("vote").Where("id=?", voteId).First(&ret).Error; err != nil {
			fmt.Printf("err:#(err.Error()}")
			return err
		}

		for _, value := range optIDs {
			if err := tx.Table("vote_opt").Where("id = ?", value).Update("count", gorm.Expr("count + ?", 1)).Error; err != nil {
				fmt.Printf("err:%s", err.Error())
				return err
			}
			user := VoteOptUser{
				VoteId:      voteId,
				UserId:      userId,
				VoteOptId:   value,
				CreatedTime: time.Now(),
			}

			if err := tx.Create(&user).Error; err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		fmt.Printf("err:%s", err.Error())
		return false
	}
	return true
}

// DoVoteV3 原生sql语句
func DoVoteV3(userId, voteId int64, optIDs []int64) bool {
	tx := Conn.Begin()
	var ret Vote

	if err := tx.Raw("select * from vote where id = ? ", voteId).Scan(&ret).Error; err != nil {
		fmt.Printf("err:%s", err.Error())
		tx.Rollback()
	}

	var oldVoteUser VoteOptUser

	if err := tx.Raw("select * from vote_opt_user where vote_id = ? and user_id = ?", voteId, userId).Scan(&oldVoteUser).Error; err != nil {
		fmt.Printf("err:%s", err.Error())
		tx.Rollback()
	}

	if oldVoteUser.ID > 0 {
		fmt.Printf("用户已投票")
		tx.Rollback()
	}

	for _, value := range optIDs {
		if err := tx.Exec("update vote_opt set count = count+1 where id = ? limit 1", value).Error; err != nil {
			fmt.Printf("err:%s", err.Error())
			tx.Rollback()
		}
		user := VoteOptUser{
			VoteId:      voteId,
			UserId:      userId,
			VoteOptId:   value,
			CreatedTime: time.Now(),
		}
		if err := tx.Create(&user).Error; err != nil {
			fmt.Printf("err:%s", err.Error())
			tx.Rollback()
		}
	}
	tx.Commit()
	return true

}
func AddVote(vote Vote, opt []VoteOpt) error {
	err := Conn.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&vote).Error; err != nil {
			return err
		}
		for _, VoteOpt := range opt {
			VoteOpt.VoteId = vote.ID
			if err := tx.Create(&VoteOpt).Error; err != nil {
				return err
			}
		}
		return nil
	})
	return err
}

func UpdateVote(vote Vote, opt []VoteOpt) error {
	err := Conn.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(&vote).Error; err != nil {
			return err
		}
		for _, VoteOpt := range opt {
			if err := tx.Save(&VoteOpt).Error; err != nil {
				return err
			}
		}
		return nil
	})
	return err
}

func DelVote(id int64) bool {
	if err := Conn.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&Vote{}, id).Error; err != nil {
			fmt.Printf("err:%s", err.Error())
			return err
		}

		if err := tx.Where("vote_id = ?", id).Delete(&VoteOpt{}).Error; err != nil {
			fmt.Printf("err:%s", err.Error())
			return err
		}

		if err := tx.Where("vote_id = ?", id).Delete(&VoteOpt{}).Error; err != nil {
			fmt.Printf("err:%s", err.Error())
			return err
		}
		return nil
	}); err != nil {
		fmt.Printf("err:%s", err.Error())
		return false
	}
	return true
}

// DelVoteV1 原生sql优化
func DelVoteV1(id int64) bool {
	if err := Conn.Transaction(func(tx *gorm.DB) error {

		if err := tx.Exec("delete from vote where id = ? limit 1", id).Error; err != nil {
			fmt.Printf("err:%s", err.Error())
			return err
		}

		if err := tx.Exec("delete from vote_opt where vote_id = ?", id).Error; err != nil {
			fmt.Printf("err:%s", err.Error())
			return err
		}

		if err := tx.Exec("delete from vote_opt_user where vote_id = ?", id).Error; err != nil {
			fmt.Printf("err:%s", err.Error())
			return err
		}
		return nil
	}); err != nil {
		fmt.Printf("err:%s", err.Error())
		return false
	}
	return true
}

func GetVoteHistory(userId, voteId int64) []VoteOptUser {
	ret := make([]VoteOptUser, 0)
	if err := Conn.Table("vote_opt_user").Where("vote_id = ? and user_id = ?", voteId, userId).First(&ret).Error; err != nil {
		fmt.Printf("err:%s", err.Error())
	}
	return ret
}

func EndVote() {
	votes := make([]Vote, 0)
	if err := Conn.Table("vote").Where("status = ?", 1).Find(&votes).Error; err != nil {
		return
	}

	now := time.Now().Unix()
	for _, vote := range votes {
		if vote.Time+vote.CreatedTime.Unix() <= now {
			Conn.Table("vote").Where("id = ?", vote.ID).Update("status", 0)

		}
	}
	return
}

// EndVoteV1 原生sql优化
func EndVoteV1() {
	votes := make([]Vote, 0)

	if err := Conn.Raw("select * from vote where status = ?", 1).Scan(&votes).Error; err != nil {
		return
	}

	now := time.Now().Unix()
	for _, vote := range votes {
		if vote.Time+vote.CreatedTime.Unix() <= now {
			Conn.Exec("update vote set status = 0 where id = ? limit 1", vote.ID)
		}
	}
	return
}

func GetVoteByName(name string) *Vote {
	var ret Vote
	if err := Conn.Table("vote").Where("title = ?", name).First(&ret).Error; err != nil {
		fmt.Printf("err:%s", err.Error())
	}
	return &ret
}
