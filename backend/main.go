package main

import (
	"context"
	"crawlab/config"
	"crawlab/database"
	_ "crawlab/docs"
	validate2 "crawlab/lib/validate"
	"crawlab/middlewares"
	"crawlab/model"
	"crawlab/routes"
	"crawlab/services"
	"crawlab/services/challenge"
	"crawlab/services/rpc"
	"github.com/apex/log"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/olivere/elastic/v7"
	"github.com/spf13/viper"
	"github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"net"
	"net/http"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"
	"time"
)

var swagHandler gin.HandlerFunc

func init() {
	swagHandler = ginSwagger.WrapHandler(swaggerFiles.Handler)
}
func main() {

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


}
