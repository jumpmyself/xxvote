package tools

import (
	"fmt"
	"github.com/bwmarrin/snowflake"
	"github.com/google/uuid"
)

var snowNode *snowflake.Node

func GetUUID() string {
	id := uuid.New() //默认v4版本 基于随机数,有的基于时间和硬件地址，比如mark
	fmt.Printf("uuid:%s,version:%s", id.String(), id.Version().String())
	return id.String()
}

func GetUid() int64 {
	if snowNode == nil {
		snowNode, _ = snowflake.NewNode(1)
	}

	return snowNode.Generate().Int64()
}
