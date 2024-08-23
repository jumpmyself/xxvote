package model

import "time"

// User undefined
type User struct {
	ID          int64     `json:"id" gorm:"id"`
	Uid         int64     `gorm:"column:uid;default:NULL"`
	Name        string    `json:"name" gorm:"name"`
	Password    string    `json:"password" gorm:"password"`
	CreatedTime time.Time `json:"created_time" gorm:"created_time"`
	UpdatedTime time.Time `json:"updated_time" gorm:"updated_time"`
	Uuid        string    `gorm:"column:uuid;default:NULL"`
}

// TableName 表名称
func (v *User) TableName() string {
	return "user"
}

// Vote undefined
type Vote struct {
	ID          int64     `json:"id" gorm:"id"  form:"ID"`
	Title       string    `json:"title" gorm:"title"  form:"Title"`
	Type        string    `json:"type" gorm:"type"  form:"Type"`
	Status      string    `json:"status" gorm:"status"  form:"Statue"`
	Time        int64     `json:"time" gorm:"time"  form:"Time"`
	UserId      int64     `json:"user_id" gorm:"user_id"  form:"UserId"`
	CreatedTime time.Time `json:"created_time" gorm:"created_time"  form:"CreatedTime"`
	UpdatedTime time.Time `json:"updated_time" gorm:"updated_time"  form:"UpdatedTime"`
}

// TableName 表名称
func (*Vote) TableName() string {
	return "vote"
}

// VoteOpt undefined
type VoteOpt struct {
	ID          int64     `json:"id" gorm:"id"`
	Name        string    `json:"name" gorm:"name"`
	VoteId      int64     `json:"vote_id" gorm:"vote_id"`
	Count       int64     `json:"count" gorm:"count"`
	CreatedTime time.Time `json:"created_time" gorm:"created_time"`
	UpdatedTime time.Time `json:"updated_time" gorm:"updated_time"`
}

// TableName 表名称
func (*VoteOpt) TableName() string {
	return "vote_opt"
}

// VoteOptUser undefined
type VoteOptUser struct {
	ID          int64     `json:"id" gorm:"id"`
	UserId      int64     `json:"user_id" gorm:"user_id"`
	VoteId      int64     `json:"vote_id" gorm:"vote_id"`
	VoteOptId   int64     `json:"vote_opt_id" gorm:"vote_opt_id"`
	CreatedTime time.Time `json:"created_time" gorm:"created_time"`
}

// TableName 表名称
func (*VoteOptUser) TableName() string {
	return "vote_opt_user"
}

type VoteWithOpt struct {
	Vote Vote
	Opt  []VoteOpt
}

type VoteWithOptV1 struct {
	Vote
	Opt []VoteOpt
}
