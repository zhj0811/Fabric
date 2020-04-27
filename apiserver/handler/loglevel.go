package handler

import (
	"net/http"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/gin-gonic/gin"
	"github.com/zhj0811/fabric/apiserver/utils"
	"github.com/zhj0811/fabric/define"
	"github.com/zhj0811/fabric/pkg/logging"
)

var logger *zap.SugaredLogger
var atomicLevel = zap.NewAtomicLevelAt(zapcore.InfoLevel)

func init() {
	config := logging.NewDefaultConfig()
	config.Level = atomicLevel
	delogger, _ := config.Build()
	logger = delogger.Named("handler").Sugar()
}

func SetLogLevel(c *gin.Context) {
	logger.Debug("Enter SetLogLevel function")
	defer logger.Debug("Exit SetLogLevel function")
	resData := &define.ResData{
		ResCode: 0,
		ResMsg:  "SUCCESS",
	}
	level := c.Query("level")
	logger.Infof("Module handler set to new loglevel: %s", level)
	atomicLevel.SetLevel(logging.NameToLevel(level))
	logger.Debug("Set log level successfully.")
	res, _ := resData.MarshalJSON()
	utils.Response(res, c, http.StatusOK)
	return
}

func GetLogLevel(c *gin.Context) {
	resData := &define.ResData{
		ResCode: 0,
		ResMsg:  "SUCCESS",
	}
	level := atomicLevel.String()
	logger.Infof("Get log level %s successful.", level)
	resData.ResMsg = level
	res, _ := resData.MarshalJSON()
	utils.Response(res, c, http.StatusOK)
	return
}
