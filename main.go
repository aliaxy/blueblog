package main

import (
	"fmt"
	"os"

	"blueblog/controller"
	"blueblog/dao/mysql"
	"blueblog/dao/redis"
	"blueblog/logger"
	"blueblog/pkg/snowflake"
	"blueblog/router"
	"blueblog/settings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("need config file.eg: blueblog config.yaml")
		return
	}
	fmt.Println(1)
	if err := settings.Init(os.Args[1]); err != nil {
		fmt.Println("init settings failed, err:", err.Error())
		return
	}
	fmt.Println(2)

	if err := logger.Init(settings.Conf.LogConfig, settings.Conf.Mode); err != nil {
		fmt.Println("init logger failed, err:", err.Error())
		return
	}
	fmt.Println(3)

	if err := mysql.Init(settings.Conf.MySQLConfig); err != nil {
		fmt.Println("init mysql failed, err:", err.Error())
		return
	}
	defer mysql.Close()

	if err := redis.Init(settings.Conf.RedisConfig); err != nil {
		fmt.Println("init redis failed, err:", err.Error())
		return
	}
	defer redis.Close()
	fmt.Println(4)

	if err := snowflake.Init(settings.Conf.StartTime, settings.Conf.MachineID); err != nil {
		fmt.Println("init snowflake failed, err:", err.Error())
		return
	}
	fmt.Println(5)

	if err := controller.InitTrans("zh"); err != nil {
		fmt.Println("init validator trans failed, err:", err.Error())
		return
	}
	fmt.Println(6)

	r := router.Setup(settings.Conf.Mode)

	err := r.Run(fmt.Sprintf(":%d", settings.Conf.Port))
	if err != nil {
		fmt.Println("run server failed, err:", err.Error())
		return
	}
}
