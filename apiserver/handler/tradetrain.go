package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/zhj0811/fabric/apiserver/utils"
	"github.com/zhj0811/fabric/common/metadata"
	"github.com/zhj0811/fabric/common/sdk"
	"github.com/zhj0811/fabric/define"

	"github.com/gin-gonic/gin"
	logging "github.com/op/go-logging"
)

var logger = logging.MustGetLogger(metadata.LogModule)

func SaveData(c *gin.Context) {
	logger.Debugf("SaveData.....")
	var request define.Factor
	var err error
	responseStatus := &define.ResponseStatus{
		StatusCode: 0,
		StatusMsg:  "SUCCESS",
	}
	ResponeData := define.FactorResponse{}

	body, err := ioutil.ReadAll(c.Request.Body)
	logger.Debugf("SaveData header : %v", c.Request.Header)
	logger.Debugf("SaveData body : %s", string(body))
	if err != nil {
		responseStatus.StatusCode = 1
		responseStatus.StatusMsg = err.Error()
		logger.Errorf("SaveData read body : %s", err.Error())
		ResponeData.ResponseCode = define.Other
		ResponeData.ResponseExplain = err.Error()
		utils.Response(ResponeData, c, http.StatusNoContent, responseStatus, nil)
		return
	}

	if err = json.Unmarshal(body, &request); err != nil {
		responseStatus.StatusCode = 1
		responseStatus.StatusMsg = err.Error()
		logger.Errorf("SaveData Unmarshal : %s", err.Error())
		ResponeData.ResponseCode = define.Other
		ResponeData.ResponseExplain = err.Error()
		utils.Response(ResponeData, c, http.StatusBadRequest, responseStatus, nil)
		return
	}
	logger.Debug(request)

	b, err := utils.FormatRequestMessage(request)
	if err != nil {
		responseStatus.StatusCode = 1
		responseStatus.StatusMsg = err.Error()
		logger.Errorf("SaveData FormatRequestMessage : %s", err.Error())
		ResponeData.ResponseCode = define.ParameterError
		ResponeData.ResponseExplain = codeMessage[define.ParameterError] + err.Error()
		utils.Response(ResponeData, c, http.StatusBadRequest, responseStatus, nil)
		return
	}

	// invoke
	var info []string
	info = append(info, define.SaveData, string(b))
	txId, err := sdk.Invoke(info)
	if err != nil {
		responseStatus.StatusCode = 1
		responseStatus.StatusMsg = err.Error()
		logger.Errorf("SaveData Invoke : %s", err.Error())
		ResponeData.ResponseCode = define.Other
		ResponeData.ResponseExplain = err.Error()
		utils.Response(ResponeData, c, http.StatusBadRequest, responseStatus, nil)
		return
	}
	logger.Debugf("This SaveData tx return's txID is %s", txId)
	ResponeData.ResponseCode = define.Success
	ResponeData.ResponseExplain = codeMessage[define.Success]
	ResponeData.FactorBaseInfo = request.FactorBaseInfo
	utils.Response(ResponeData, c, http.StatusOK, responseStatus, nil)
}

func QueryDataByKey(c *gin.Context) {
	var request define.QueryData
	var err error
	responseStatus := &define.ResponseStatus{
		StatusCode: 0,
		StatusMsg:  "SUCCESS",
	}
	ResponeData := define.QueryDataResponse{}

	body, err := ioutil.ReadAll(c.Request.Body)
	logger.Debugf("QueryData header : %v", c.Request.Header)
	logger.Debugf("QueryData body : %s", string(body))
	if err != nil {
		responseStatus.StatusCode = 1
		responseStatus.StatusMsg = err.Error()
		logger.Errorf("QueryData read body : %s", err.Error())
		ResponeData.ResponseCode = define.Other
		ResponeData.ResponseExplain = err.Error()
		utils.Response(ResponeData, c, http.StatusNoContent, responseStatus, nil)
		return
	}

	err = json.Unmarshal(body, &request)
	if err != nil {
		responseStatus.StatusCode = 1
		responseStatus.StatusMsg = err.Error()
		logger.Errorf("QueryData Unmarshal : %s", err.Error())
		ResponeData.ResponseCode = define.Other
		ResponeData.ResponseExplain = err.Error()
		utils.Response(ResponeData, c, http.StatusBadRequest, responseStatus, nil)
		return
	}
	if err = utils.VerifyQueryDataRequestFormat(&request); err != nil {
		ResponeData.ResponseCode = define.ParameterError
		ResponeData.ResponseExplain = codeMessage[define.ParameterError] + err.Error()
		responseStatus.StatusCode = 1
		responseStatus.StatusMsg = codeMessage[define.ParameterError]
		logger.Errorf("QueryData: %s", codeMessage[define.ParameterError])
		utils.Response(ResponeData, c, http.StatusBadRequest, responseStatus, nil)
		return
	}

	logger.Debug(request)
	var info []string
	info = append(info, define.QueryDataByKey, request.Key, request.BusinessType, request.DataType, request.WriteRoleType, request.Reader)
	var responseData []byte
	queryRespose := &define.QueryResponse{}

	responseData, err = sdk.Query(info)

	if err != nil {
		//errStr := fmt.Sprintf("Failed to query access control list: %s", err.Error())
		errstr := err.Error()
		fmt.Println("errInfo is: ", errstr)
		permissionNo := strings.Contains(errstr, define.PermissionNotFound)
		valueNil := strings.Contains(errstr, define.ValueOfKeyNil)
		NoPermission := strings.Contains(errstr, define.NoPermission)
		if permissionNo {
			ResponeData.ResponseCode = define.PermissionNotFound
			ResponeData.ResponseExplain = codeMessage[define.PermissionNotFound]
		} else if valueNil {
			ResponeData.ResponseCode = define.ValueOfKeyNil
			ResponeData.ResponseExplain = codeMessage[define.ValueOfKeyNil]
		} else if NoPermission {
			ResponeData.ResponseCode = define.NoPermission
			ResponeData.ResponseExplain = codeMessage[define.NoPermission]
		}

		utils.Response(ResponeData, c, http.StatusBadRequest, responseStatus, nil)
		return
	}

	err = utils.Unmarshal(responseData, queryRespose)
	if err != nil {
		errStr := fmt.Sprintf("Failed to Unmarshal QueryResponse:%s", err.Error())
		ResponeData.ResponseCode = "1"
		ResponeData.ResponseExplain = errStr
		utils.Response(ResponeData, c, http.StatusBadRequest, responseStatus, nil)
		return
	}
	var msg define.Message
	err = utils.Unmarshal([]byte(queryRespose.Payload.(string)), &msg)
	if err != nil {
		errStr := fmt.Sprintf("Failed to Unmarshal message:%s", err.Error())
		ResponeData.ResponseCode = define.Other
		ResponeData.ResponseExplain = errStr
		utils.Response(ResponeData, c, http.StatusBadRequest, responseStatus, nil)
		return
	}

	var factorList define.Factor
	if err := utils.FormatResponseMessage(&factorList, &msg); err != nil {
		errStr := fmt.Sprintf("Failed to FormatResponseAccessMessage:%s", err.Error())
		ResponeData.ResponseCode = define.Other
		ResponeData.ResponseExplain = errStr
		utils.Response(ResponeData, c, http.StatusBadRequest, responseStatus, nil)
		return
	}
	ResponeData.ResponseCode = define.Success
	ResponeData.ResponseExplain = codeMessage[define.Success]
	ResponeData.FactorBaseInfo = factorList.FactorBaseInfo
	ResponeData.BusinessData = factorList.BusinessData
	utils.Response(ResponeData, c, http.StatusOK, responseStatus, nil)
}

func KeepaliveQuery(c *gin.Context) {
	status := http.StatusOK
	responseStatus := &define.ResponseStatus{
		StatusCode: 200,
		StatusMsg:  "SUCCESS",
	}

	if err := sdk.PeerKeepalive(); err != nil {
		responseStatus.StatusCode = define.PeerFailed
		responseStatus.StatusMsg = "Peer FAILED"
		status = define.PeerFailed
		logger.Error("peer cann't be reached.")
		utils.Response(nil, c, status, responseStatus, nil)
		return
	}

	if err := sdk.OrderKeepalive(); err != nil {
		responseStatus.StatusCode = define.OrdererFailed
		responseStatus.StatusMsg = "Order FAILED"
		status = define.OrdererFailed
		logger.Error("order cann't be reached.")
	}

	utils.Response(nil, c, status, responseStatus, nil)
}

func BlockHeight(c *gin.Context) {
	responseStatus := &define.ResponseStatus{
		StatusCode: 0,
		StatusMsg:  "SUCCESS",
	}

	height, err := sdk.GetBlockHeightByEndorserPeer()
	if nil != err {
		responseStatus.StatusCode = 1
		responseStatus.StatusMsg = err.Error()
		utils.Response(0, c, http.StatusBadRequest, responseStatus, nil)
		return
	}

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
		responseStatus.StatusCode = returnCode
		responseStatus.StatusMsg = err.Error()
		utils.Response(nil, c, http.StatusInternalServerError, responseStatus, nil)
		return
	}
	utils.Response(nil, c, http.StatusOK, responseStatus, nil)
}
