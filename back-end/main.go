package main

import (
	"com/mittacy/gomeet/config"
	"com/mittacy/gomeet/controller"
	"com/mittacy/gomeet/database"
	"com/mittacy/gomeet/logger"
	"com/mittacy/gomeet/router"
	"flag"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

func main() {
	// 1. 初始化配置文件
	if err := config.InitConfig(); err != nil {
		panic(err)
	}

	// 2. 初始化日志系统
	if err := logger.InitLogger(); err != nil {
		panic(err)
	}

	// 3. 初始化数据库
	if err := database.InitMysql(); err != nil {
		logger.Record("连接数据库错误", err)
		panic(err)
	}
	defer database.CloseMysql()

	if err := database.InitRedis(); err != nil {
		logger.Record("连接Redis错误", err)
		panic(err)
	}
	defer database.CloseRedis()

	// 4. 路由
	r := router.InitRouter()

	// 5. 启动定时任务
	controller.InitCron()

	// 6. 服务器配置，启动服务
	// 命令行指定端口，如果没有则使用配置文件的端口
	serverConfig, err := config.Cfg.GetSection("server")
	if err != nil {
		logger.Record("获取服务配置信息错误", err)
		panic(err)
	}
	confPort, err := serverConfig.Key("port").Int()
	if err != nil {
		panic(err)
	}
	port := flag.Int("port", confPort, "监听端口，int类型")
	flag.Parse()
	fmt.Println("监听端口: " + strconv.Itoa(*port))

	readTimeout, err := serverConfig.Key("read_timeout").Int()
	if err != nil {
		logger.Record("获取服务配置信息错误", err)
		panic(err)
	}
	writeTimeout, err := serverConfig.Key("write_timeout").Int()
	if err != nil {
		logger.Record("获取服务配置信息错误", err)
		panic(err)
	}
	s := &http.Server{
		Addr:           ":" + strconv.Itoa(*port),
		Handler:        r,
		ReadTimeout:    time.Duration(readTimeout) * time.Second,
		WriteTimeout:   time.Duration(writeTimeout) * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	s.ListenAndServe()
}
