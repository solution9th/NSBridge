package main

import (
	"time"

	"github.com/solution9th/NSBridge/app"
	"github.com/solution9th/NSBridge/internal/utils"
)

func main() {

	logPath, ok := utils.FindFile("./dns.log", "/var/log/ns_bridge/dns.log")
	if !ok {
		panic("not found log")
	}

	// 初始化日志
	utils.NewLogFile(logPath, time.RFC3339)

	// without extension
	err := app.Run("config")
	if err != nil {
		panic(err)
	}
}
