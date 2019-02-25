package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
	"strings"
	"syscall"

	"github.com/peersafe/factoring/apiserver/router"
	"github.com/peersafe/factoring/common/metadata"
	"github.com/peersafe/factoring/common/sdk"

	"github.com/DeanThompson/ginpprof"
	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"
	"github.com/op/go-logging"
	"github.com/spf13/viper"
)

var (
	configPath = flag.String("configPath", "./", "config file path")
	configName = flag.String("configName", "client_sdk", "config file name")
	isVersion  = flag.Bool("v", false, "Show version information")
)

// package-scoped variables
var logger = logging.MustGetLogger(metadata.LogModule)

// package-scoped constants
const packageName = "apiserver"

func main() {
	// parse init param
	flag.Parse()
	if *isVersion {
		printVersion()
		return
	}

	err := sdk.InitSDKs(*configPath, *configName)
	if err != nil {
		logger.Errorf("init sdk error : %s\n", err.Error())
		panic(err)
	}

	if err := sdk.SetLogLevel(viper.GetString("log.logLevel"), metadata.LogModule); err != nil {
		logger.Errorf("SetLogLevel error : %s\n", err.Error())
		return
	}
	// 设置使用系统最大CPU
	runtime.GOMAXPROCS(runtime.NumCPU())

	// 运行模式
	gin.SetMode(gin.ReleaseMode) //DebugMode ReleaseMode

	// 构造路由器
	r := router.GetRouter()

	// 调试用,可以看到堆栈状态和所有goroutine状态
	ginpprof.Wrapper(r)

	//Get the listen port for apiserver
	listenPort := viper.GetInt("apiserver.listenport")
	logger.Debug("The listen port is", listenPort)
	listenPortString := fmt.Sprintf(":%d", listenPort)

	// 运行服务
	server := endless.NewServer(listenPortString, r)
	server.BeforeBegin = func(add string) {
		pid := syscall.Getpid()
		logger.Criticalf("Actual pid is %d", pid)
		// 保存pid文件
		pidFile := "apiserver.pid"
		if checkFileIsExist(pidFile) {
			os.Remove(pidFile)
		} else {
			if err := ioutil.WriteFile(pidFile, []byte(fmt.Sprintf("%d", pid)), 0666); err != nil {
				logger.Fatalf("Api server write pid file failed! err:%v\n", err)
			}
		}
	}
	err = server.ListenAndServe()
	if err != nil {
		if strings.Contains(err.Error(), "use of closed network connection") {
			logger.Errorf("%v\n", err)
		} else {
			logger.Errorf("Api server start failed! err:%v\n", err)
			panic(err)
		}
	}
}

func checkFileIsExist(filename string) bool {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return false
	}
	return true
}

func printVersion() {
	version := metadata.GetVersionInfo()
	fmt.Println(packageName, "with version:", version)
	fmt.Println()
}
