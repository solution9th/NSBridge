package app

import (
	"sync"

	"github.com/solution9th/NSBridge/internal/config"
	"github.com/solution9th/NSBridge/internal/oneapm"
	"github.com/solution9th/NSBridge/internal/service/cache"
	"github.com/solution9th/NSBridge/internal/service/mysql"
	"github.com/solution9th/NSBridge/internal/utils"
)

// Run run app
func Run(fileName string) (err error) {

	err = Init(fileName)
	if err != nil {
		return err
	}

	exitCh := make(chan error)
	var once sync.Once
	wg := &WaitGroupWrapper{}

	apm, err := runOneAPM()
	if err != nil {
		return err
	}

	oneapm.AppTransmission(apm)

	exitFunc := func(err error) {
		once.Do(func() {
			if err != nil {
				utils.Error("run error:", err)
			}
			exitCh <- err
		})
	}

	wg.Wrap(func() {
		exitFunc(runGRpc())
	})

	// wg.Wrap(func() {
	// 	exitFunc(runGateway())
	// })

	wg.Wrap(func() {
		exitFunc(runWeb())
	})

	return <-exitCh
}

// Init 初始化项目
func Init(fileName string) (err error) {

	err = config.InitConfig(fileName)
	if err != nil {
		return err
	}

	var (
		redisHost   = config.RedisConfig.Host
		redisPort   = config.RedisConfig.Port
		redisPasswd = config.RedisConfig.Passwd
		redisDB     = config.RedisConfig.DB

		mysqlHost   = config.MySQLConfig.Host
		mysqlPort   = config.MySQLConfig.Port
		mysqlUser   = config.MySQLConfig.User
		mysqlPasswd = config.MySQLConfig.Passwd
		mysqlDB     = config.MySQLConfig.DBName
	)

	// 初始化 redis
	err = cache.InitRedis(redisHost, redisPort, redisPasswd, redisDB)
	if err != nil {
		utils.Error("[Init] initRedis error:", err)
		return err
	}

	// 初始化 默认MySQL
	err = mysql.InitDefaultDB(mysqlDB, mysqlUser, mysqlPasswd, mysqlHost, mysqlPort)
	if err != nil {
		utils.Error("[init] initMySQL error:", err)
		return err
	}

	return nil
}

type WaitGroupWrapper struct {
	sync.WaitGroup
}

// Wrap 包装后的待执行函数都会在协成中执行
func (w *WaitGroupWrapper) Wrap(cb func()) {
	w.Add(1)
	go func() {
		cb()
		w.Done()
	}()
}
