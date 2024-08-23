package model

import (
	"fmt"
	"strconv"
	"testing"
	"time"
)

func TestGetVotes(t *testing.T) {
	NewMysql()
	//测试用例
	r := GetVotes()
	fmt.Printf("ret:%+v", r)
	Close()
}

func TestGetVote(t *testing.T) {
	NewMysql()
	//测试用例
	r := GetVote(2)
	fmt.Printf("ret:%+v", r)
	Close()
}

func TestDoVote(t *testing.T) {
	NewMysql()
	defer Close()
	r := DoVote(1, 1, []int64{1, 2})
	fmt.Printf("ret:%+v", r)
}

func TestAddVote(t *testing.T) {
	NewMysql()
	defer Close()
	//测试用例
	vote := Vote{
		Title:       "测试用例",
		Type:        strconv.Itoa(0),
		Status:      strconv.Itoa(0),
		Time:        0,
		UserId:      0,
		CreatedTime: time.Now(),
		UpdatedTime: time.Now(),
	}

	opt := make([]VoteOpt, 0)
	opt = append(opt, VoteOpt{
		Name:        "测试选项1",
		Count:       0,
		CreatedTime: time.Now(),
		UpdatedTime: time.Now(),
	})
	opt = append(opt, VoteOpt{
		Name:        "测试选项2",
		Count:       0,
		CreatedTime: time.Now(),
		UpdatedTime: time.Now(),
	})
	r := AddVote(vote, opt)

	fmt.Printf("ret:%+v", r)
	Close()
}

func TestGetUserV1(t *testing.T) {
	NewMysql()
	a := GetUserV1("admin")
	fmt.Printf("a:%+v", a)
}
