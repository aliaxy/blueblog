package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"main/dao/mysql"
	"main/dao/redis"
	"main/logger"
	"main/pkg/snowflake"
	"main/router"
	"main/settings"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func main() {
	if err := settings.Init(); err != nil {
		fmt.Println("init settings failed, err: " + err.Error())
	}

	if err := logger.Init(settings.Conf.LogCofig); err != nil {
		fmt.Println("init logger failed, err: " + err.Error())
	}
	defer zap.L().Sync()
	zap.L().Debug("logger init success...")

	if err := mysql.Init(settings.Conf.MySQLConfig); err != nil {
		fmt.Println("init mysql failed, err: " + err.Error())
	}
	defer mysql.Close()

	if err := redis.Init(settings.Conf.RedisConfig); err != nil {
		fmt.Println("init redis failed, err: " + err.Error())
	}
	defer redis.Close()

	if err := snowflake.Init(settings.Conf.StartTime, settings.Conf.MachineID); err != nil {
		fmt.Println("init snowflake failed, err: " + err.Error())
	}

	r := router.Setup()

	r.Run()

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", viper.GetInt("app.port")),
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)                      // 创建一个接收信号的通道
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM) // 此处不会阻塞
	<-quit                                               // 阻塞在此，当接收到上述两种信号时才会往下执行
	log.Println("Shutdown Server ...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown: ", err)
	}

	log.Println("Server exiting")
}
