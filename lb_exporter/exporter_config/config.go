package exporter_config

import (
	"exporter-demo/logger"
	"fmt"
	"github.com/spf13/viper"
)

// LB config
func LbConfig() string {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		logger.Error("Error reading config file, %s", err)
	}

	return viper.GetString("exporter.port")

}

// boss config
func ConfigBoss() (host, protocol, access_key, secret_key, zone string, port int, ignore_verify bool) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		logger.Error("Error reading config file, %s", err)
		fmt.Println("Error reading config file, %s", err)
	}

	HOST := viper.GetString("boss.host")
	BossPORT := viper.GetInt("boss.boss_port")
	PROTOCOL := viper.GetString("boss.protocol")
	ACCESS_KEY := viper.GetString("boss.access_key")
	SECRET_KEY := viper.GetString("boss.secret_key")
	ZONE := viper.GetString("boss.zone")
	IGNORE_VERIFY := false

	return HOST, PROTOCOL, ACCESS_KEY, SECRET_KEY, ZONE, BossPORT, IGNORE_VERIFY

}
