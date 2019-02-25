package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os/exec"
	"strings"

	"github.com/zhaojianpeerfintech/fabric/apiserver/utils"
	"github.com/zhaojianpeerfintech/fabric/common/metadata"
	"github.com/zhaojianpeerfintech/fabric/common/sdk"
	"github.com/zhaojianpeerfintech/fabric/define"

	"github.com/gin-gonic/gin"
)

func parseFabricBaseInfo(header http.Header, checkValue define.CheckFabricBaseInfo) (*define.FabricBaseInfo, error) {
	var fabricBaseInfo define.FabricBaseInfo
/*
	logger.Warningf("len of define.ChannelName is %d", len(header[define.ChannelName]))

	switch checkValue {
	case define.CNameCFBI:
		if 0 == len(header[define.ChannelName]) {
			err := errors.New("must specify channel name")
			logger.Error(err.Error())
			return nil, err
		}
	case define.CCNameCFBI:
		if 0 == len(header[define.ChaincodeName]) {
			err := errors.New("must specify chaincode name")
			logger.Error(err.Error())
			return nil, err
		}
	case define.CNameAndCCNameCFBI:
	default:
		if 0 == len(header[define.ChannelName]) || 0 == len(header[define.ChaincodeName]) {
			err := errors.New("must specify channel and chaincode name")
			logger.Error(err.Error())
			return nil, err
		}
	}*/
	if len(header[define.ChannelName]) > 0 {
		logger.Debugf("channel name is %s", header[define.ChannelName])
		fabricBaseInfo.ChannelName = header[define.ChannelName][0]
	} else if checkValue == define.CNameCFBI || checkValue == define.CNameAndCCNameCFBI {
		err := errors.New("must specify channel name")
		logger.Error(err.Error())
		return nil, err
	}

	if len(header[define.ChaincodeName]) > 0 {
		fabricBaseInfo.ChaincodeName = (header[define.ChaincodeName])[0]
	} else if checkValue == define.CCNameCFBI || checkValue == define.CNameAndCCNameCFBI {
		err := errors.New("must specify chaincode name")
		logger.Error(err.Error())
		return nil, err
	}

	if len(header[define.ChaincodeVersion]) > 0 {
		fabricBaseInfo.ChaincodeVersion = (header[define.ChaincodeVersion])[0]
	}

	logger.Infof("the channel name is %s, chaincode name is %s and chaincode version is %s.",
		fabricBaseInfo.ChannelName, fabricBaseInfo.ChaincodeName, fabricBaseInfo.ChaincodeVersion)

	return &fabricBaseInfo, nil
}

func getReqInfo(c *gin.Context, targetInfo interface{}) (*define.ResponseStatus, int) {
	logger.Debugf("enter getRequestInfo function")
	defer logger.Debugf("exit getRequestInfo function")

	status := http.StatusOK
	resStatus := &define.ResponseStatus{
		StatusCode: 0,
		StatusMsg:  "success",
	}

	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		status = http.StatusBadRequest
		resStatus.StatusCode = define.ReadReuqestError
		resStatus.StatusMsg = err.Error()
		logger.Errorf("read request body failed: %s", err.Error())
		return resStatus, status
	}
	logger.Debugf("read request body is: %s", string(body))

	if err = json.Unmarshal(body, targetInfo); err != nil {
		status = http.StatusBadRequest
		resStatus.StatusCode = define.UnmarshalError
		resStatus.StatusMsg = err.Error()
		logger.Errorf("unmarshal request body failed: %s", err.Error())
		return resStatus, status
	}
	logger.Infof("The request is %v", targetInfo)

	return resStatus, status
}

func BlockHeight(c *gin.Context) {
	responseStatus := &define.ResponseStatus{
		StatusCode: 0,
		StatusMsg:  "SUCCESS",
	}

	logger.Debugf("SaveData header: %v", c.Request.Header)
	fabricBaseInfo, err := parseFabricBaseInfo(c.Request.Header, define.CNameCFBI)
	if err != nil {
		responseStatus.StatusCode = http.StatusBadRequest
		responseStatus.StatusMsg = err.Error()
		logger.Errorf("parseFabricBaseInfo failed: %s", err.Error())
		utils.Response(nil, c, http.StatusBadRequest, responseStatus, nil)
		return
	}

	height, err := sdk.GetBlockHeightByEndorserPeer(fabricBaseInfo.ChannelName)
	if nil != err {
		logger.Errorf("get block height failed: %s", err.Error())
		responseStatus.StatusCode = 1
		responseStatus.StatusMsg = err.Error()
		utils.Response(0, c, http.StatusBadRequest, responseStatus, nil)
		return
	}
	logger.Infof("the block height is %d", height)

	utils.Response(height, c, http.StatusOK, responseStatus, nil)
}

func KafkaNumber(c *gin.Context) {
	logger.Debug("enter KafkaNumber function.")

	responseStatus := &define.ResponseStatus{
		StatusCode: 0,
		StatusMsg:  "SUCCESS",
	}

	returnCode, err := sdk.GetKafkaNumber()
	if nil != err {
		logger.Errorf("get kafka number failed: %s", err.Error())
		responseStatus.StatusCode = returnCode
		responseStatus.StatusMsg = err.Error()
		utils.Response(nil, c, http.StatusInternalServerError, responseStatus, nil)
		return
	}
	utils.Response(nil, c, http.StatusOK, responseStatus, nil)
}

func Version(c *gin.Context) {
	status := http.StatusOK
	responseStatus := &define.ResponseStatus{
		StatusCode: 200,
		StatusMsg:  "SUCCESS",
	}

	version := metadata.GetVersionInfo()
	logger.Infof("the version is %s", version)

	utils.Response(version, c, status, responseStatus, nil)
}

func KeepaliveQuery(c *gin.Context) {
	status := http.StatusOK
	responseStatus := &define.ResponseStatus{
		StatusCode: 200,
		StatusMsg:  "SUCCESS",
	}

	logger.Debugf("KeepaliveQuery header : %v", c.Request.Header)
	fabricBaseInfo, err := parseFabricBaseInfo(c.Request.Header, define.CNameAndCCNameCFBI)
	if err != nil {
		responseStatus.StatusCode = http.StatusBadRequest
		responseStatus.StatusMsg = err.Error()
		logger.Errorf("parseFabricBaseInfo failed: %s", err.Error())
		utils.Response(nil, c, http.StatusBadRequest, responseStatus, nil)
		return
	}

	if err := sdk.PeerKeepalive(fabricBaseInfo.ChannelName, fabricBaseInfo.ChaincodeName); err != nil {
		responseStatus.StatusCode = define.PeerFailed
		responseStatus.StatusMsg = "Peer FAILED"
		status = define.PeerFailed
		logger.Error("peer can not be reached.")
		utils.Response(nil, c, status, responseStatus, nil)
		return
	}

	if err := sdk.OrderKeepalive(); err != nil {
		responseStatus.StatusCode = define.OrdererFailed
		responseStatus.StatusMsg = "Orderer FAILED"
		status = define.OrdererFailed
		logger.Error("orderer can not be reached.")
		utils.Response(nil, c, status, responseStatus, nil)
		return
	}
	logger.Debug("the peer and orderer is normal as usual")

	utils.Response(nil, c, status, responseStatus, nil)
}

func Keepalive(c *gin.Context) {
	status := http.StatusOK
	responseStatus := &define.ResponseStatus{
		StatusCode: 200,
		StatusMsg:  "SUCCESS",
	}

	logger.Debugf("KeepaliveQuery header : %v", c.Request.Header)
	fabricBaseInfo, err := parseFabricBaseInfo(c.Request.Header, define.CNameAndCCNameCFBI)
	if err != nil {
		responseStatus.StatusCode = http.StatusBadRequest
		responseStatus.StatusMsg = err.Error()
		logger.Errorf("parseFabricBaseInfo failed: %s", err.Error())
		utils.Response(nil, c, http.StatusBadRequest, responseStatus, nil)
		return
	}

	if err := sdk.PeerKeepalive(fabricBaseInfo.ChannelName, fabricBaseInfo.ChaincodeName); err != nil {
		responseStatus.StatusCode = define.PeerFailed
		responseStatus.StatusMsg = "Peer FAILED"
		status = define.PeerFailed
		logger.Error("peer can not be reached.")
		utils.Response(responseStatus.StatusMsg, c, status, responseStatus, nil)
		return
	}

	if err := sdk.OrderKeepalive(); err != nil {
		responseStatus.StatusCode = define.OrdererFailed
		responseStatus.StatusMsg = "Orderer FAILED"
		status = define.OrdererFailed
		logger.Error("orderer can not be reached.")
		utils.Response(responseStatus.StatusMsg, c, status, responseStatus, nil)
		return
	}
	logger.Debug("the peer and orderer is normal as usual")

	utils.Response(responseStatus.StatusMsg, c, status, responseStatus, nil)
}

func SetLogLevel(c *gin.Context) {
	var request define.LogLevel
	var command, echoCommand string
	status := http.StatusOK
	responseStatus := &define.ResponseStatus{
		StatusCode: 0,
		StatusMsg:  "SUCCESS",
	}

	logger.Debug("enter SetLogLevel function")
	defer logger.Debug("exit SetLogLevel function")

	if responseStatus, status = getReqInfo(c, &request); http.StatusOK != status {
		logger.Errorf("get request info failed: %s.", responseStatus.StatusMsg)
		utils.Response(nil, c, status, responseStatus, nil)
		return
	}
	logger.Debugf("set module %s's level to %s.", metadata.LogModule, request.Level)
	err := sdk.SetLogLevel(request.Level, metadata.LogModule)
	if nil != err {
		errMsg := fmt.Sprintf("set module %s's level to %s failed.", metadata.LogModule, request.Level)
		logger.Error(errMsg)
		status = http.StatusInternalServerError
		responseStatus.StatusCode = define.LogModuleSetError
		responseStatus.StatusMsg = errMsg
		utils.Response("Set Log Level Failed.", c, http.StatusBadRequest, responseStatus, nil)
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
	utils.Response("Set Log Level Successfully.", c, http.StatusOK, responseStatus, nil)
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
