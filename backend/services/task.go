package services

import (
	"bufio"
	"crawunit/constants"
	"crawunit/database"
	"crawunit/entity"
	"crawunit/lib/cron"
	"crawunit/model"
	"crawunit/services/spider_handler"
	"crawunit/utils"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/apex/log"
	"github.com/globalsign/mgo/bson"
	uuid "github.com/satori/go.uuid"
	"github.com/spf13/viper"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"runtime/debug"
	"strconv"
	"strings"
	"sync"
	"time"
)

var Exec *Executor

// 任务执行锁
//Added by cloud: 2019/09/04,solve data race
var LockList sync.Map

// 任务消息
type TaskMessage struct {
	Id  string
	Cmd string
}

// 序列化任务消息
func (m *TaskMessage) ToString() (string, error) {
	data, err := json.Marshal(&m)
	if err != nil {
		return "", err
	}
	return utils.BytesToString(data), err
}

// 任务执行器
type Executor struct {
	Cron *cron.Cron
}

// 启动任务执行器
func (ex *Executor) Start() error {
	// 启动cron服务
	ex.Cron.Start()

	// 加入执行器到定时任务
	spec := "0/1 * * * * *" // 每秒执行一次
	for i := 0; i < viper.GetInt("task.workers"); i++ {
		// WorkerID
		id := i


		// ------->初始化锁
		LockList.Store(id, false)
		// <-------


		_, err := ex.Cron.AddFunc(spec, GetExecuteTaskFunc(id))
		if err != nil {
			return err
		}
	}

	return nil
}

// 派发任务
func AssignTask(task model.Task) error {
	// 生成任务信息
	msg := TaskMessage{
		Id: task.Id,
	}

	// 序列化
	msgStr, err := msg.ToString()
	if err != nil {
		return err
	}

	// 队列名称
	var queue string
	if utils.IsObjectIdNull(task.NodeId) {
		queue = "tasks:public"
	} else {
		queue = "tasks:node:" + task.NodeId.Hex()
	}

	// 任务入队
	if err := database.RedisClient.RPush(queue, msgStr); err != nil {
		return err
	}
	return nil
}

// 设置环境变量
func SetEnv(cmd *exec.Cmd, envs []model.Env, task model.Task, spider model.Spider) *exec.Cmd {
	// 默认把Node.js的全局node_modules加入环境变量
	envPath := os.Getenv("PATH")
	homePath := os.Getenv("HOME")
	nodeVersion := "v10.19.0"
	nodePath := path.Join(homePath, ".nvm/versions/node", nodeVersion, "lib/node_modules")
	if !strings.Contains(envPath, nodePath) {
		_ = os.Setenv("PATH", nodePath+":"+envPath)
	}
	_ = os.Setenv("NODE_PATH", nodePath)

	// default results collection
	col := utils.GetSpiderCol(spider.Col, spider.Name)

	// 默认环境变量
	cmd.Env = append(os.Environ(), "CRAWLAB_TASK_ID="+task.Id)
	cmd.Env = append(cmd.Env, "CRAWLAB_COLLECTION="+col)
	cmd.Env = append(cmd.Env, "CRAWLAB_MONGO_HOST="+viper.GetString("mongo.host"))
	cmd.Env = append(cmd.Env, "CRAWLAB_MONGO_PORT="+viper.GetString("mongo.port"))
	if viper.GetString("mongo.db") != "" {
		cmd.Env = append(cmd.Env, "CRAWLAB_MONGO_DB="+viper.GetString("mongo.db"))
	}
	if viper.GetString("mongo.username") != "" {
		cmd.Env = append(cmd.Env, "CRAWLAB_MONGO_USERNAME="+viper.GetString("mongo.username"))
	}
	if viper.GetString("mongo.password") != "" {
		cmd.Env = append(cmd.Env, "CRAWLAB_MONGO_PASSWORD="+viper.GetString("mongo.password"))
	}
	if viper.GetString("mongo.authSource") != "" {
		cmd.Env = append(cmd.Env, "CRAWLAB_MONGO_AUTHSOURCE="+viper.GetString("mongo.authSource"))
	}
	cmd.Env = append(cmd.Env, "PYTHONUNBUFFERED=0")
	cmd.Env = append(cmd.Env, "PYTHONIOENCODING=utf-8")
	cmd.Env = append(cmd.Env, "TZ=Asia/Shanghai")
	cmd.Env = append(cmd.Env, "CRAWLAB_DEDUP_FIELD="+spider.DedupField)
	cmd.Env = append(cmd.Env, "CRAWLAB_DEDUP_METHOD="+spider.DedupMethod)
	if spider.IsDedup {
		cmd.Env = append(cmd.Env, "CRAWLAB_IS_DEDUP=1")
	} else {
		cmd.Env = append(cmd.Env, "CRAWLAB_IS_DEDUP=0")
	}

	//任务环境变量
	for _, env := range envs {
		cmd.Env = append(cmd.Env, env.Name+"="+env.Value)
	}

	// 全局环境变量
	variables := model.GetVariableList()
	for _, variable := range variables {
		cmd.Env = append(cmd.Env, variable.Key+"="+variable.Value)
	}
	return cmd
}

func SetLogConfig(cmd *exec.Cmd, t model.Task, u model.User) error {
	// get stdout reader
	stdout, err := cmd.StdoutPipe()
	readerStdout := bufio.NewReader(stdout)
	if err != nil {
		log.Errorf("get stdout error: %s", err.Error())
		debug.PrintStack()
		return err
	}

	// get stderr reader
	stderr, err := cmd.StderrPipe()
	readerStderr := bufio.NewReader(stderr)
	if err != nil {
		log.Errorf("get stdout error: %s", err.Error())
		debug.PrintStack()
		return err
	}

	var seq int64
	var logs []model.LogItem
	isStdoutFinished := false
	isStderrFinished := false

	// periodically (1 sec) insert log items
	go func() {
		for {
			_ = model.AddLogItems(logs)
			logs = []model.LogItem{}
			if isStdoutFinished && isStderrFinished {
				break
			}
			time.Sleep(5 * time.Second)
		}
	}()

	// expire duration (in seconds)
	expireDuration := u.Setting.LogExpireDuration
	if expireDuration == 0 {
		// by default not expire
		expireDuration = constants.Infinite
	}

	// read stdout
	go func() {
		for {
			line, err := readerStdout.ReadString('\n')
			if err != nil {
				isStdoutFinished = true
				break
			}
			line = strings.Replace(line, "\n", "", -1)
			seq++
			l := model.LogItem{
				Id:       bson.NewObjectId(),
				Seq:      seq,
				Message:  line,
				TaskId:   t.Id,
				Ts:       time.Now(),
				ExpireTs: time.Now().Add(time.Duration(expireDuration) * time.Second),
			}
			logs = append(logs, l)
		}
	}()

	// read stderr
	go func() {
		for {
			line, err := readerStderr.ReadString('\n')
			if err != nil {
				isStderrFinished = true
				break
			}
			line = strings.Replace(line, "\n", "", -1)
			seq++
			l := model.LogItem{
				Id:       bson.NewObjectId(),
				Seq:      seq,
				Message:  line,
				TaskId:   t.Id,
				Ts:       time.Now(),
				ExpireTs: time.Now().Add(time.Duration(expireDuration) * time.Second),
			}
			logs = append(logs, l)
		}
	}()

	return nil
}

func FinishOrCancelTask(ch chan string, cmd *exec.Cmd, s model.Spider, t model.Task) {
	// 传入信号，此处阻塞
	signal := <-ch
	log.Infof("process received signal: %s", signal)

	if signal == constants.TaskCancel && cmd.Process != nil {
		var err error
		// 兼容windows
		if runtime.GOOS == constants.Windows {
			err = cmd.Process.Kill()
		} /*else {
			err = syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
		}*/
		// 取消进程
		if err != nil {
			log.Errorf("process kill error: %s", err.Error())
			debug.PrintStack()

			t.Error = "kill process error: " + err.Error()
			t.Status = constants.StatusError
		} else {
			t.Error = "user kill the process ..."
			t.Status = constants.StatusCancelled
		}
	} else {
		// 保存任务
		t.Status = constants.StatusFinished
	}

	t.FinishTs = time.Now()
	_ = t.Save()

	go FinishUpTask(s, t)
}

func StartTaskProcess(cmd *exec.Cmd, t model.Task) error {
	if err := cmd.Start(); err != nil {
		log.Errorf("start spider error:{}", err.Error())
		debug.PrintStack()

		t.Error = "start task error: " + err.Error()
		t.Status = constants.StatusError
		t.FinishTs = time.Now()
		_ = t.Save()
		return err
	}
	return nil
}

func WaitTaskProcess(cmd *exec.Cmd, t model.Task, s model.Spider) error {
	if err := cmd.Wait(); err != nil {
		log.Errorf("wait process finish error: %s", err.Error())
		debug.PrintStack()

		if exitError, ok := err.(*exec.ExitError); ok {
			exitCode := exitError.ExitCode()
			log.Errorf("exit error, exit code: %d", exitCode)

			// 非kill 的错误类型
			if exitCode != -1 {
				// 非手动kill保存为错误状态
				t.Error = err.Error()
				t.FinishTs = time.Now()
				t.Status = constants.StatusError
				_ = t.Save()

				FinishUpTask(s, t)
			}
		}

		return err
	}

	return nil
}

// 执行shell命令
func ExecuteShellCmd(cmdStr string, cwd string, t model.Task, s model.Spider, u model.User) (err error) {
	log.Infof("cwd: %s", cwd)
	log.Infof("cmd: %s", cmdStr)

	// 生成执行命令
	var cmd *exec.Cmd
	if runtime.GOOS == constants.Windows {
		cmd = exec.Command("cmd", "/C", cmdStr)
	} else {
		cmd = exec.Command("sh", "-c", cmdStr)
	}

	// 工作目录
	cmd.Dir = cwd

	// 日志配置
	if err := SetLogConfig(cmd, t, u); err != nil {
		return err
	}

	// 环境变量配置
	envs := s.Envs
	if s.Type == constants.Configurable {
		// 数据库配置
		envs = append(envs, model.Env{Name: "CRAWLAB_MONGO_HOST", Value: viper.GetString("mongo.host")})
		envs = append(envs, model.Env{Name: "CRAWLAB_MONGO_PORT", Value: viper.GetString("mongo.port")})
		envs = append(envs, model.Env{Name: "CRAWLAB_MONGO_DB", Value: viper.GetString("mongo.db")})
		envs = append(envs, model.Env{Name: "CRAWLAB_MONGO_USERNAME", Value: viper.GetString("mongo.username")})
		envs = append(envs, model.Env{Name: "CRAWLAB_MONGO_PASSWORD", Value: viper.GetString("mongo.password")})
		envs = append(envs, model.Env{Name: "CRAWLAB_MONGO_AUTHSOURCE", Value: viper.GetString("mongo.authSource")})

		// 设置配置
		for envName, envValue := range s.Config.Settings {
			envs = append(envs, model.Env{Name: "CRAWLAB_SETTING_" + envName, Value: envValue})
		}
	}
	cmd = SetEnv(cmd, envs, t, s)

	// 起一个goroutine来监控进程
	ch := utils.TaskExecChanMap.ChanBlocked(t.Id)

	go FinishOrCancelTask(ch, cmd, s, t)

/*	// kill的时候，可以kill所有的子进程
	if runtime.GOOS != constants.Windows {
		cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	}*/

	// 启动进程
	if err := StartTaskProcess(cmd, t); err != nil {
		return err
	}

	// 同步等待进程完成
	if err := WaitTaskProcess(cmd, t, s); err != nil {
		return err
	}
	ch <- constants.TaskFinish
	return nil
}

// 生成日志目录
func MakeLogDir(t model.Task) (fileDir string, err error) {
	// 日志目录
	fileDir = filepath.Join(viper.GetString("log.path"), t.SpiderId.Hex())

	// 如果日志目录不存在，生成该目录
	if !utils.Exists(fileDir) {
		if err := os.MkdirAll(fileDir, 0777); err != nil {
			log.Errorf("execute task, make log dir error: %s", err.Error())
			debug.PrintStack()
			return "", err
		}
	}

	return fileDir, nil
}

// 获取日志文件路径
func GetLogFilePaths(fileDir string, t model.Task) (filePath string) {
	// 时间戳
	ts := time.Now()
	tsStr := ts.Format("20060102150405")

	// stdout日志文件
	filePath = filepath.Join(fileDir, t.Id+"_"+tsStr+".log")

	return filePath
}

// 生成执行任务方法
func GetExecuteTaskFunc(id int) func() {
	return func() {
		ExecuteTask(id)
	}
}

func GetWorkerPrefix(id int) string {
	return "[Worker " + strconv.Itoa(id) + "] "
}

// 统计任务结果数
func SaveTaskResultCount(id string) func() {
	return func() {
		if err := model.UpdateTaskResultCount(id); err != nil {
			log.Errorf(err.Error())
			debug.PrintStack()
			return
		}
	}
}

// Scan Error Logs
func ScanErrorLogs(t model.Task) func() {
	return func() {
		u, err := model.GetUser(t.UserId)
		if err != nil {
			return
		}
		if err := model.UpdateTaskErrorLogs(t.Id, u.Setting.ErrorRegexPattern); err != nil {
			return
		}
		if err := model.UpdateErrorLogCount(t.Id); err != nil {
			return
		}
	}
}

// 执行任务
func ExecuteTask(id int) {
	//------>判断该执行器是否被占用:若被占用,写入日志并返回;若没被占用才可以继续操作
	if flag, ok := LockList.Load(id); ok {
		if flag.(bool) {
			log.Debugf(GetWorkerPrefix(id) + "running tasks...")
			return
		}
	}


	//------>上锁
	LockList.Store(id, true)

	//------>程序return的时候,解锁,初始化锁
	defer func() {
		LockList.Delete(id)
		LockList.Store(id, false)
	}()
	//<------

	// 开始计时
	tic := time.Now()

	// 获取当前节点
	node, err := model.GetCurrentNode()
	if err != nil {
		log.Errorf("execute task get current node error: %s", err.Error())
		debug.PrintStack()
		return
	}

	// 节点队列
	queueCur := "tasks:node:" + node.Id.Hex()

	// 节点队列任务
	var msg string
	if msg, err = database.RedisClient.LPop(queueCur); err != nil {
		// 节点队列没有任务，获取公共队列任务
		queuePub := "tasks:public"
		if msg, err = database.RedisClient.LPop(queuePub); err != nil {
		}
	}

	// 如果没有获取到任务，返回
	if msg == "" {
		return
	}

	// 反序列化
	tMsg := TaskMessage{}
	if err := json.Unmarshal([]byte(msg), &tMsg); err != nil {
		log.Errorf("json string to struct error: %s", err.Error())
		return
	}

	// 获取任务
	t, err := model.GetTask(tMsg.Id)
	if err != nil {
		log.Errorf("execute task, get task error: %s", err.Error())
		return
	}

	// 获取爬虫
	spider, err := t.GetSpider()
	if err != nil {
		log.Errorf("execute task, get spider error: %s", err.Error())
		return
	}

	// 创建日志目录
	var fileDir string
	if fileDir, err = MakeLogDir(t); err != nil {
		return
	}
	// 获取日志文件路径
	t.LogPath = GetLogFilePaths(fileDir, t)

	// 工作目录
	cwd := filepath.Join(
		viper.GetString("spider.path"),
		spider.Name,
	)

	// 执行命令
	var cmd string
	if spider.Type == constants.Configurable {
		// 可配置爬虫命令
		cmd = "scrapy crawl config_spider"
	} else {
		// 自定义爬虫命令
		cmd = spider.Cmd
	}

	// 加入参数
	if t.Param != "" {
		cmd += " " + t.Param
	}

	// 获得触发任务用户
	user, err := model.GetUser(t.UserId)
	if err != nil {
		log.Errorf(GetWorkerPrefix(id) + err.Error())
		return
	}

	// 任务赋值
	t.NodeId = node.Id                                   // 任务节点信息
	t.StartTs = time.Now()                               // 任务开始时间
	t.Status = constants.StatusRunning                   // 任务状态
	t.WaitDuration = t.StartTs.Sub(t.CreateTs).Seconds() // 等待时长



	// 文件检查
	if err := SpiderFileCheck(t, spider); err != nil {
		log.Errorf("spider file check error: %s", err.Error())
		return
	}

	// 开始执行任务
	log.Infof(GetWorkerPrefix(id) + "start task (id:" + t.Id + ")")

	// 储存任务
	_ = t.Save()

	// 起一个cron执行器来统计任务结果数
	cronExec := cron.New(cron.WithSeconds())
	_, err = cronExec.AddFunc("*/5 * * * * *", SaveTaskResultCount(t.Id))
	if err != nil {
		log.Errorf(GetWorkerPrefix(id) + err.Error())
		debug.PrintStack()
		return
	}
	cronExec.Start()
	defer cronExec.Stop()

	// 起一个cron来更新错误日志
	cronExecErrLog := cron.New(cron.WithSeconds())
	_, err = cronExecErrLog.AddFunc("*/30 * * * * *", ScanErrorLogs(t))
	if err != nil {
		log.Errorf(GetWorkerPrefix(id) + err.Error())
		debug.PrintStack()
		return
	}
	cronExecErrLog.Start()
	defer cronExecErrLog.Stop()

	// 执行Shell命令
	if err := ExecuteShellCmd(cmd, cwd, t, spider, user); err != nil {
		log.Errorf(GetWorkerPrefix(id) + err.Error())

		// 如果发生错误，则发送通知
		t, _ = model.GetTask(t.Id)



		return
	}

	// 完成进程
	t, err = model.GetTask(t.Id)
	if err != nil {
		log.Errorf(GetWorkerPrefix(id) + err.Error())
		return
	}

	// 统计数据
	t.Status = constants.StatusFinished                     // 任务状态: 已完成
	t.FinishTs = time.Now()                                 // 结束时间
	t.RuntimeDuration = t.FinishTs.Sub(t.StartTs).Seconds() // 运行时长
	t.TotalDuration = t.FinishTs.Sub(t.CreateTs).Seconds()  // 总时长





	// 保存任务
	if err := t.Save(); err != nil {
		log.Errorf(GetWorkerPrefix(id) + err.Error())
		return
	}

	// 完成任务收尾工作
	go FinishUpTask(spider, t)

	// 结束计时
	toc := time.Now()

	// 统计时长
	duration := toc.Sub(tic).Seconds()
	durationStr := strconv.FormatFloat(duration, 'f', 6, 64)
	log.Infof(GetWorkerPrefix(id) + "task (id:" + t.Id + ")" + " finished. elapsed:" + durationStr + " sec")
}

func FinishUpTask(s model.Spider, t model.Task) {
	// 更新任务结果数
	go func() {
		if err := model.UpdateTaskResultCount(t.Id); err != nil {
			return
		}
	}()

	// 更新任务错误日志
	go func() {
		ScanErrorLogs(t)()
	}()
}

func SpiderFileCheck(t model.Task, spider model.Spider) error {
	// 判断爬虫文件是否存在
	gfFile := model.GetGridFs(spider.FileId)
	if gfFile == nil {
		t.Error = "cannot find spider files, please re-upload"
		t.Status = constants.StatusError
		t.FinishTs = time.Now()                                 // 结束时间
		t.RuntimeDuration = t.FinishTs.Sub(t.StartTs).Seconds() // 运行时长
		t.TotalDuration = t.FinishTs.Sub(t.CreateTs).Seconds()  // 总时长
		_ = t.Save()
		return errors.New(t.Error)
	}

	// 判断md5值是否一致
	path := filepath.Join(viper.GetString("spider.path"), spider.Name)
	md5File := filepath.Join(path, spider_handler.Md5File)
	md5 := utils.GetSpiderMd5Str(md5File)
	if gfFile.Md5 != md5 {
		spiderSync := spider_handler.SpiderSync{Spider: spider}
		spiderSync.RemoveDownCreate(gfFile.Md5)
	}
	return nil
}

func GetTaskLog(id string, keyword string, page int, pageSize int) (logItems []model.LogItem, logTotal int, err error) {
	task, err := model.GetTask(id)
	if err != nil {
		return
	}

	logItems, logTotal, err = task.GetLogItems(keyword, page, pageSize)
	if err != nil {
		return logItems, logTotal, err
	}

	return logItems, logTotal, nil
}

func GetTaskErrorLog(id string, n int) (errLogItems []model.ErrorLogItem, err error) {
	if n == 0 {
		n = 1000
	}

	task, err := model.GetTask(id)
	if err != nil {
		return
	}
	errLogItems, err = task.GetErrorLogItems(n)
	if err != nil {
		return
	}
	return errLogItems, nil
}

func CancelTask(id string) (err error) {
	// 获取任务
	task, err := model.GetTask(id)
	if err != nil {
		log.Errorf("task not found, task id : %s, error: %s", id, err.Error())
		debug.PrintStack()
		return err
	}

	// 如果任务状态不为pending或running，返回错误
	if task.Status != constants.StatusPending && task.Status != constants.StatusRunning {
		return errors.New("task is not cancellable")
	}

	// 获取当前节点（默认当前节点为主节点）
	node, err := model.GetCurrentNode()
	if err != nil {
		log.Errorf("get current node error: %s", err.Error())
		debug.PrintStack()
		return err
	}

	log.Infof("current node id is: %s", node.Id.Hex())
	log.Infof("task node id is: %s", task.NodeId.Hex())

	if node.Id == task.NodeId {
		// 任务节点为主节点

		// 获取任务执行频道
		ch := utils.TaskExecChanMap.ChanBlocked(id)
		if ch != nil {
			// 发出取消进程信号
			ch <- constants.TaskCancel
		} else {
			if err := model.UpdateTaskToAbnormal(node.Id); err != nil {
				log.Errorf("update task to abnormal : {}", err.Error())
				debug.PrintStack()
				return err
			}
		}
	} else {
		// 任务节点为工作节点

		// 序列化消息
		msg := entity.NodeMessage{
			Type:   constants.MsgTypeCancelTask,
			TaskId: id,
		}
		msgBytes, err := json.Marshal(&msg)
		if err != nil {
			return err
		}

		// 发布消息
		if _, err := database.RedisClient.Publish("nodes:"+task.NodeId.Hex(), utils.BytesToString(msgBytes)); err != nil {
			return err
		}
	}

	return nil
}

func RestartTask(id string, uid bson.ObjectId) (err error) {
	// 获取任务
	oldTask, err := model.GetTask(id)
	if err != nil {
		log.Errorf("task not found, task id : %s, error: %s", id, err.Error())
		debug.PrintStack()
		return err
	}

	newTask := model.Task{
		SpiderId:   oldTask.SpiderId,
		NodeId:     oldTask.NodeId,
		Param:      oldTask.Param,
		UserId:     uid,
		RunType:    oldTask.RunType,
		ScheduleId: bson.ObjectIdHex(constants.ObjectIdNull),
	}

	// 加入任务队列
	_, err = AddTask(newTask)
	if err != nil {
		log.Errorf(err.Error())
		debug.PrintStack()
		return err
	}

	return nil
}

func AddTask(t model.Task) (string, error) {
	// 生成任务ID
	id := uuid.NewV4()
	t.Id = id.String()

	// 设置任务状态
	t.Status = constants.StatusPending

	// 如果没有传入node_id，则置为null
	if t.NodeId.Hex() == "" {
		t.NodeId = bson.ObjectIdHex(constants.ObjectIdNull)
	}

	// 将任务存入数据库
	if err := model.AddTask(t); err != nil {
		log.Errorf(err.Error())
		debug.PrintStack()
		return t.Id, err
	}

	// 加入任务队列
	if err := AssignTask(t); err != nil {
		log.Errorf(err.Error())
		debug.PrintStack()
		return t.Id, err
	}

	return t.Id, nil
}

func GetTaskEmailMarkdownContent(t model.Task, s model.Spider) string {
	n, _ := model.GetNode(t.NodeId)
	errMsg := ""
	statusMsg := fmt.Sprintf(`<span style="color:green">%s</span>`, t.Status)
	if t.Status == constants.StatusError {
		errMsg = " with errors"
		statusMsg = fmt.Sprintf(`<span style="color:red">%s</span>`, t.Status)
	}
	return fmt.Sprintf(`
Your task has finished%s. Please find the task info below.

 | 
--: | :--
**Task ID:** | %s
**Task Status:** | %s
**Task Param:** | %s
**Spider ID:** | %s
**Spider Name:** | %s
**Node:** | %s
**Create Time:** | %s
**Start Time:** | %s
**Finish Time:** | %s
**Wait Duration:** | %.0f sec
**Runtime Duration:** | %.0f sec
**Total Duration:** | %.0f sec
**Number of Results:** | %d
**Error:** | <span style="color:red">%s</span>

`,
		errMsg,
		t.Id,
		statusMsg,
		t.Param,
		s.Id.Hex(),
		s.Name,
		n.Name,
		utils.GetLocalTimeString(t.CreateTs),
		utils.GetLocalTimeString(t.StartTs),
		utils.GetLocalTimeString(t.FinishTs),
		t.WaitDuration,
		t.RuntimeDuration,
		t.TotalDuration,
		t.ResultCount,
		t.Error,
	)
}

func GetTaskMarkdownContent(t model.Task, s model.Spider) string {
	n, _ := model.GetNode(t.NodeId)
	errMsg := ""
	errLog := "-"
	statusMsg := fmt.Sprintf(`<font color="#00FF00">%s</font>`, t.Status)
	if t.Status == constants.StatusError {
		errMsg = `（有错误）`
		errLog = fmt.Sprintf(`<font color="#FF0000">%s</font>`, t.Error)
		statusMsg = fmt.Sprintf(`<font color="#FF0000">%s</font>`, t.Status)
	}
	return fmt.Sprintf(`
您的任务已完成%s，请查看任务信息如下。

> **任务ID:** %s  
> **任务状态:** %s  
> **任务参数:** %s  
> **爬虫ID:** %s  
> **爬虫名称:** %s  
> **节点:** %s  
> **创建时间:** %s  
> **开始时间:** %s  
> **完成时间:** %s  
> **等待时间:** %.0f秒   
> **运行时间:** %.0f秒  
> **总时间:** %.0f秒  
> **结果数:** %d  
> **错误:** %s  

`,
		errMsg,
		t.Id,
		statusMsg,
		t.Param,
		s.Id.Hex(),
		s.Name,
		n.Name,
		utils.GetLocalTimeString(t.CreateTs),
		utils.GetLocalTimeString(t.StartTs),
		utils.GetLocalTimeString(t.FinishTs),
		t.WaitDuration,
		t.RuntimeDuration,
		t.TotalDuration,
		t.ResultCount,
		errLog,
	)
}



func InitTaskExecutor() error {
	// 构造任务执行器
	c := cron.New(cron.WithSeconds())
	Exec = &Executor{
		Cron: c,
	}

	// 如果不允许主节点运行任务，则跳过
	if model.IsMaster() && viper.GetString("setting.runOnMaster") == "N" {
		return nil
	}

	// 运行定时任务
	if err := Exec.Start(); err != nil {
		return err
	}
	return nil
}
