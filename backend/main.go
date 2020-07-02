package main

import (
	"context"
	"crawunit/config"
	"crawunit/database"
	"crawunit/lib/validate_bridge"
	"crawunit/middlewares"
	"crawunit/model"
	"crawunit/routes"
	"crawunit/services"
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

	// 初始化日志设置
	if err := services.InitLogService(); err != nil {
		log.Error("init log error:" + err.Error())
		panic(err)
	}
	log.Info("initialized log successfully")

	if model.IsMaster() {
		// 初始化清理服务
		if err := services.InitCleanService(); err != nil {
			log.Error("init clean service error:" + err.Error())
			debug.PrintStack()
			panic(err)
		}
		log.Info("initialized clean service successfully")
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


	// 以下为主节点服务
	if model.IsMaster() {
		// 中间件
		app.Use(middlewares.CORSMiddleware())
		anonymousGroup := app.Group("/")
		{
			anonymousGroup.POST("/login", routes.Login)       // 用户登录

			// 节点
			{
				anonymousGroup.GET("/nodes", routes.GetNodeList)                            // 节点列表
				anonymousGroup.GET("/nodes/:id", routes.GetNode)                            // 节点详情
				anonymousGroup.POST("/nodes/:id", routes.PostNode)                          // 修改节点
				anonymousGroup.GET("/nodes/:id/tasks", routes.GetNodeTaskList)              // 节点任务列表
				anonymousGroup.GET("/nodes/:id/system", routes.GetSystemInfo)               // 节点任务列表
				anonymousGroup.DELETE("/nodes/:id", routes.DeleteNode)                      // 删除节点
			}
			// 爬虫
			{
				anonymousGroup.GET("/spiders", routes.GetSpiderList)                                            // 爬虫列表
				anonymousGroup.GET("/spiders/:id", routes.GetSpider)                                            // 爬虫详情
				anonymousGroup.PUT("/spiders", routes.PutSpider)                                                // 添加爬虫
				anonymousGroup.POST("/spiders", routes.UploadSpider)                                            // 上传爬虫
				anonymousGroup.POST("/spiders/:id", routes.PostSpider)                                          // 修改爬虫
				anonymousGroup.POST("/spiders/:id/publish", routes.PublishSpider)                               // 发布爬虫
				anonymousGroup.POST("/spiders/:id/upload", routes.UploadSpiderFromId)                           // 上传爬虫（ID）
				anonymousGroup.DELETE("/spiders", routes.DeleteSelectedSpider)                                  // 删除选择的爬虫
				anonymousGroup.GET("/spiders/:id/tasks", routes.GetSpiderTasks)                                 // 爬虫任务列表
				anonymousGroup.GET("/spiders/:id/file/tree", routes.GetSpiderFileTree)                          // 爬虫文件目录树读取
				anonymousGroup.GET("/spiders/:id/file", routes.GetSpiderFile)                                   // 爬虫文件读取
				anonymousGroup.GET("/spiders/:id/dir", routes.GetSpiderDir)                                     // 爬虫目录
				anonymousGroup.GET("/spiders/:id/scrapy/spiders", routes.GetSpiderScrapySpiders)                // Scrapy 爬虫名称列表
				anonymousGroup.POST("/spiders-cancel", routes.CancelSelectedSpider)                             // 停止所选爬虫任务
				anonymousGroup.POST("/spiders-run", routes.RunSelectedSpider)                                   // 运行所选爬虫
			}
			// 任务
			{
				anonymousGroup.GET("/tasks", routes.GetTaskList)                                 // 任务列表
				anonymousGroup.GET("/tasks/:id", routes.GetTask)                                 // 任务详情
				anonymousGroup.PUT("/tasks", routes.PutTask)                                     // 派发任务
				anonymousGroup.DELETE("/tasks/:id", routes.DeleteTask)                           // 删除任务
				anonymousGroup.POST("/tasks/:id/cancel", routes.CancelTask)                      // 取消任务
				anonymousGroup.GET("/tasks/:id/log", routes.GetTaskLog)                          // 任务日志
				anonymousGroup.GET("/tasks/:id/results", routes.GetTaskResults)                  // 任务结果
				anonymousGroup.POST("/tasks/:id/restart", routes.RestartTask)                    // 重新开始任务
			}
			// 用户
			{
				anonymousGroup.GET("/me", routes.GetMe)                // 获取自己账户
			}
			// 项目
			{
				anonymousGroup.GET("/projects", routes.GetProjectList)       // 列表
				anonymousGroup.GET("/projects/tags", routes.GetProjectTags)  // 项目标签
				anonymousGroup.PUT("/projects", routes.PutProject)           // 修改
				anonymousGroup.POST("/projects/:id", routes.PostProject)     // 新增
				anonymousGroup.DELETE("/projects/:id", routes.DeleteProject) // 删除
			}

			// 文件
			anonymousGroup.GET("/file", routes.GetFile) // 获取文件

		}
	}

	// 路由ping
	app.GET("/ping", routes.Ping)

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
