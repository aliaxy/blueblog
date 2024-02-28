// Package settings 导入配置文件
// 使用 viper
package settings

import (
	"fmt"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// Conf 全局访问
var Conf = new(AppConfig)

// AppConfig 应用配置
type AppConfig struct {
	Name      string `mapstructure:"name"`
	Mode      string `mapstructure:"mode"`
	Version   string `mapstructure:"version"`
	Port      int    `mapstructure:"port"`
	StartTime string `mapstructure:"start_time"`
	MachineID int64  `mapstructure:"machine_id"`

	*LogConfig   `mapstructure:"log"`
	*MySQLConfig `mapstructure:"mysql"`
	*RedisConfig `mapstructure:"redis"`
}

// LogConfig 日志配置
type LogConfig struct {
	Level      string `mapstructure:"level"`
	Filename   string `mapstructure:"filename"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxAge     int    `mapstructure:"max_age"`
	MaxBackups int    `mapstructure:"max_backups"`
}

// MySQLConfig mysql 配置
type MySQLConfig struct {
	Host         string `mapstructure:"host"`
	Port         int    `mapstructure:"port"`
	Password     string `mapstructure:"password"`
	DB           string `mapstructure:"dbname"`
	User         string `mapstructure:"user"`
	MaxOpenConns int    `mapstructure:"max_open_conns"`
	MaxIdleConns int    `mapstructure:"max_idle_conns"`
}

// RedisConfig Redis 配置
type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
	PoolSize int    `mapstructure:"pool_size"`
}

// Init 初始化
func Init(filePath string) (err error) {
	viper.SetConfigFile(filePath)

	err = viper.ReadInConfig()
	if err != nil {
		fmt.Println("viper.ReadInConfig() failed, err: " + err.Error())
		return err
	}

	if err := viper.Unmarshal(Conf); err != nil {
		fmt.Println("viper.Unmarshal() failed, err: " + err.Error())
	}

	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		fmt.Println("Config changed")
		if err := viper.Unmarshal(Conf); err != nil {
			fmt.Println("viper.Unmarshal() failed, err: " + err.Error())
		}
	})
	return
}
