package config

import (
	"fmt"
	"github.com/spf13/viper"
)

type TomlConfig struct {
	AppName        string
	MySQL          MySQLConfig
	Log            LogConfig
	StaticPath     PathConfig
	MsgChannelType MsgChannelType
}

type MySQLConfig struct {
	Host        string
	Name        string
	Password    string
	Port        int
	TablePrefix string
	User        string
}

type LogConfig struct {
	Path  string
	Level string
}

type PathConfig struct {
	FilePath string
}

// 消息队列类型及其消息队列相关信息
// gochannel为单机使用go默认的channel进行消息传递
// kafka是使用kafka作为消息队列，可以分布式扩展消息聊天程序
type MsgChannelType struct {
	ChannelType string
	KafkaHosts  string
	KafkaTopic  string
}

var c TomlConfig

func init() {
	//设置文件名
	viper.SetConfigName("config")
	//设置文件类型
	viper.SetConfigType("toml")
	//设置文件路径，可以多个viper会根据设置顺序一次查找
	viper.AddConfigPath(".")
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %s", err))
	}
	viper.Unmarshal(&c)
}

func GetConfig() TomlConfig {
	return c
}
