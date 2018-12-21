package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os/exec"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/peersafe/tradetrain/apiserver/utils"
	"github.com/peersafe/tradetrain/common/metadata"
	"github.com/peersafe/tradetrain/common/sdk"
	"github.com/peersafe/tradetrain/define"
)

func SetLogLevel(c *gin.Context) {
	var request define.LogLevel
	var command, echoCommand string
	responseStatus := &define.ResponseStatus{
		StatusCode: 0,
		StatusMsg:  "SUCCESS",
	}

	logger.Debug("enter SetLogLevel function")
	defer logger.Debug("exit SetLogLevel function")

	body, err := ioutil.ReadAll(c.Request.Body)
	logger.Debugf("SetLogLevel header : %v", c.Request.Header)
	logger.Debugf("SetLogLevel body : %s", string(body))
	if err != nil {
		responseStatus.StatusCode = 1
		responseStatus.StatusMsg = err.Error()
		logger.Errorf("SetLogLevel read body : %s", err.Error())
		utils.Response(responseStatus, c, http.StatusNoContent, responseStatus, nil)
		return
	}
	if err = json.Unmarshal(body, &request); err != nil {
		responseStatus.StatusCode = 1
		responseStatus.StatusMsg = err.Error()
		logger.Errorf("SetLogLevel define.LogLevel Unmarshal : %s", err.Error())
		utils.Response(responseStatus, c, http.StatusBadRequest, responseStatus, nil)
	}

	logger.Infof("set module %s's level to %s.", metadata.LogModule, request.Level)
	err = sdk.SetLogLevel(request.Level, metadata.LogModule)
	if nil != err {
		errMsg := fmt.Sprintf("set module %s's level to %s failed.", metadata.LogModule, request.Level)
		logger.Error(errMsg)
		responseStatus.StatusCode = define.LogModuleSetError
		responseStatus.StatusMsg = errMsg
		utils.Response(responseStatus, c, http.StatusBadRequest, responseStatus, nil)
		return
	}
	logger.Debug("Set log level successfully.")
	configPath := "."
	configFile := "client_sdk"
	if strings.HasSuffix(configPath, "/") {
		command = fmt.Sprintf("$(sed 's/^    logLevel.*/    logLevel: %s/g' %s%s.yaml)", request.Level, configPath, configFile)
		echoCommand = fmt.Sprintf("echo %q > %s%s.yaml", command, configPath, configFile)
	} else {
		command = fmt.Sprintf("$(sed 's/^    logLevel.*/    logLevel: %s/g' %s/%s.yaml)", request.Level, configPath, configFile)
		echoCommand = fmt.Sprintf("echo %q > %s/%s.yaml", command, configPath, configFile)
	}

	logger.Info("the echoCommand is:", echoCommand)

	cmd := exec.Command("/bin/bash", "-c", echoCommand)
	err = cmd.Run()
	if nil != err {
		logger.Error("set logLevel:", request.Level, "to the file failed:", err.Error())
		return
	}
	utils.Response(responseStatus, c, http.StatusOK, responseStatus, nil)
	return
}

func GetLogLevel(c *gin.Context) {
	status := http.StatusOK
	responseStatus := &define.ResponseStatus{
		StatusCode: 0,
		StatusMsg:  "SUCCESS",
	}

	logger.Debug("enter GetLogLevel function")
	defer logger.Debug("exit GetLogLevel function")

	logger.Debugf("try to get %s's log level.", metadata.LogModule)
	level := sdk.GetLogLevel(metadata.LogModule)
	if level == "" {
		errMsg := fmt.Sprintf("get module %s's level is failed.", metadata.LogModule)
		logger.Error(errMsg)
		status = http.StatusInternalServerError
		responseStatus.StatusCode = define.LogModuleInvalid
		responseStatus.StatusMsg = errMsg
		utils.Response(nil, c, status, responseStatus, nil)
	}
	logger.Infof("module %s's log level is %s", metadata.LogModule, level)

	logger.Debug("Get log level successfully.")
	utils.Response(level, c, status, responseStatus, nil)

	return
}
