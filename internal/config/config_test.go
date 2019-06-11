package config

import (
	"fmt"
	"testing"
)

func TestInitConfig(t *testing.T) {

	err := InitConfig("config")
	if err != nil {
		panic(err)
	}

	fmt.Println(ServerConfig.Debug)
	fmt.Println(NameServerConfig)
}
