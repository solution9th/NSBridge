package config

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"

	"github.com/spf13/viper"
)

var (
	IsDebug = false
)

var (
	MySQLConfig      ConfigMySQL
	RedisConfig      ConfigRedis
	ServerConfig     ConfigServer
	FoneConfig       ConfigFOne
	GRpcConfig       ConfigGRpc
	SamlConfig       ConfigSaml
	NameServerConfig ConfigNameServer
)

// InitConfig init config
func InitConfig(fileName string) error {

	viper.SetConfigName(fileName)
	viper.AddConfigPath(".")
	viper.AddConfigPath("/etc/ns_bridge/")
	viper.AddConfigPath("../")
	viper.AddConfigPath("../../")

	err := viper.ReadInConfig()
	if err != nil {
		log.Println("read config error:", err)
		return err
	}

	err = get(viper.GetStringMap("mysql"), &MySQLConfig)
	if err != nil {
		return err
	}

	err = get(viper.GetStringMap("redis"), &RedisConfig)
	if err != nil {
		return err
	}

	err = get(viper.GetStringMap("server"), &ServerConfig)
	if err != nil {
		return err
	}

	err = get(viper.GetStringMap("grpc"), &GRpcConfig)
	if err != nil {
		return err
	}

	err = get(viper.GetStringMap("fone"), &FoneConfig)
	if err != nil {
		return err
	}

	err = get(viper.GetStringMap("saml"), &SamlConfig)
	if err != nil {
		return err
	}

	NameServerConfig = viper.GetStringSlice("nameserver")

	if ServerConfig.Debug {
		IsDebug = true
	}

	return nil
}

func get(config map[string]interface{}, ptr interface{}) error {

	if reflect.TypeOf(ptr).Kind() != reflect.Ptr {
		return fmt.Errorf("params need ptr")
	}

	tmpTmp, err := json.Marshal(config)
	if err != nil {
		return err
	}

	return json.Unmarshal(tmpTmp, ptr)
}
