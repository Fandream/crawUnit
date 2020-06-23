package main

import (
	"context"
	"crawlab/config"
	"crawlab/database"
	"crawlab/lib/validate_bridge"
	"crawlab/middlewares"
	"crawlab/model"
	"crawlab/routes"
	"crawlab/services"
	"crawlab/services/challenge"
	"crawlab/services/rpc"
	"github.com/apex/log"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/spf13/viper"
	"net"
	"net/http"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"
	"time"
)

func main() {
	binding.Validator = new(validate_bridge.DefaultValidator)
	app := gin.Default()

	// 初始化配置
	if err := config.InitConfig(""); err != nil {
		log.Error("init config error:" + err.Error())
		panic(err)
	}
	log.Info("initialized config successfully")
	// 初始化Mongodb数据库
	if err := database.InitMongo(); err != nil {
		log.Error("init mongodb error:" + err.Error())
		debug.PrintStack()
		panic(err)
	}
	log.Info("initialized mongodb successfully")

	// 初始化Redis数据库
	if err := database.InitRedis(); err != nil {
		log.Error("init redis error:" + err.Error())
		debug.PrintStack()
		panic(err)
	}
	log.Info("initialized redis successfully")



	if model.IsMaster() {
		// 初始化定时任务
		if err := services.InitScheduler(); err != nil {
			log.Error("init scheduler error:" + err.Error())
			debug.PrintStack()
			panic(err)
		}
		log.Info("initialized schedule successfully")

		// 初始化用户服务
		if err := services.InitUserService(); err != nil {
			log.Error("init user service error:" + err.Error())
			debug.PrintStack()
			panic(err)
		}
		log.Info("initialized user service successfully")

		// 初始化依赖服务
		if err := services.InitDepsFetcher(); err != nil {
			log.Error("init dependency fetcher error:" + err.Error())
			debug.PrintStack()
			panic(err)
		}
		log.Info("initialized dependency fetcher successfully")

	}

	// 初始化任务执行器
	if err := services.InitTaskExecutor(); err != nil {
		log.Error("init task executor error:" + err.Error())
		debug.PrintStack()
		panic(err)
	}
	log.Info("initialized task executor successfully")

	// 初始化节点服务
	if err := services.InitNodeService(); err != nil {
		log.Error("init node service error:" + err.Error())
		panic(err)
	}
	log.Info("initialized node service successfully")

	// 初始化爬虫服务
	if err := services.InitSpiderService(); err != nil {
		log.Error("init spider service error:" + err.Error())
		debug.PrintStack()
		panic(err)
	}
	log.Info("initialized spider service successfully")

	// 初始化RPC服务
	if err := rpc.InitRpcService(); err != nil {
		log.Error("init rpc service error:" + err.Error())
		debug.PrintStack()
		panic(err)
	}
	log.Info("initialized rpc service successfully")


	// 运行服务器
	host := viper.GetString("server.host")
	port := viper.GetString("server.port")
	address := net.JoinHostPort(host, port)
	srv := &http.Server{
		Handler: app,
		Addr:    address,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			if err != http.ErrServerClosed {
				log.Error("run server error:" + err.Error())
			} else {
				log.Info("server graceful down")
			}
		}
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	ctx2, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx2); err != nil {
		log.Error("run server error:" + err.Error())
	}
}
