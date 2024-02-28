// Package snowflake 雪花算法
package snowflake

import (
	"time"

	"github.com/bwmarrin/snowflake"
)

var node *snowflake.Node

// Init 初始化
func Init(startTime string, machineID int64) (err error) {
	var st time.Time
	st, err = time.Parse("2006-01-02", startTime)
	if err != nil {
		return
	}
	snowflake.Epoch = st.UnixNano() / 1_000_000
	node, err = snowflake.NewNode(machineID)
	return
}

// GenID 生成一个 id
func GenID() int64 {
	return node.Generate().Int64()
}
