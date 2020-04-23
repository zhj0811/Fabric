package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/zhj0811/fabric/apiserver/utils"
	"github.com/zhj0811/fabric/common/sdk"
	"github.com/zhj0811/fabric/define"

	"github.com/gin-gonic/gin"
)

func SaveUserInfo(c *gin.Context) {
	logger.Debugf("SaveUserInfo.....")
	var request define.UserInfo
	var err error
	responseStatus := &define.ResponseStatus{
		StatusCode: 0,
		StatusMsg:  "SUCCESS",
	}
	responseData := define.SaveUserInfoRespone{}

	body, err := ioutil.ReadAll(c.Request.Body)
	logger.Debugf("SaveUserInfo header : %v", c.Request.Header)
	logger.Debugf("SaveUserInfo body : %s", string(body))
	if err != nil {
		responseStatus.StatusCode = 1
		responseStatus.StatusMsg = err.Error()
		responseData.ResponseCode = define.Other
		responseData.ResponseExplain = err.Error()
		logger.Errorf("SaveUserInfo read body : %s", err.Error())
		utils.Response(responseData, c, http.StatusNoContent, responseStatus, nil)
		return
	}

	if err = json.Unmarshal(body, &request); err != nil {
		responseStatus.StatusCode = 1
		responseStatus.StatusMsg = err.Error()
		responseData.ResponseCode = define.Other
		responseData.ResponseExplain = err.Error()
		logger.Errorf("SaveUserInfo define.UserInfo Unmarshal : %s", err.Error())
		utils.Response(responseData, c, http.StatusBadRequest, responseStatus, nil)
		return
	}
	logger.Debug(request)
	b, err := utils.FormatUserInfoRequestMessage(request)
	if err != nil {
		responseStatus.StatusCode = 1
		responseStatus.StatusMsg = err.Error()
		responseData.ResponseCode = define.ParameterError
		responseData.ResponseExplain = codeMessage[define.ParameterError] + err.Error()
		logger.Errorf("SaveUserInfo FormatRequestMessage : %s", err.Error())
		utils.Response(responseData, c, http.StatusBadRequest, responseStatus, nil)
		return
	}

	// invoke
	var info []string
	info = append(info, define.SaveUserInfo, string(b))
	txId, err := sdk.Invoke(info)
	if err != nil {
		responseStatus.StatusCode = 1
		responseStatus.StatusMsg = err.Error()
		logger.Errorf("SaveData Invoke : %s", err.Error())
		responseData.ResponseCode = define.Other
		responseData.ResponseExplain = err.Error()
		utils.Response(responseData, c, http.StatusBadRequest, responseStatus, nil)
		return
	}
	logger.Debugf("This tx return's txID is %s", txId)
	responseData.ResponseCode = define.Success
	responseData.ResponseExplain = codeMessage[define.Success]
	responseData.UserBaseInfo = request.UserBaseInfo
	utils.Response(responseData, c, http.StatusOK, responseStatus, nil)
}

func QueryUserInfo(c *gin.Context) {
	logger.Debugf("QueryUserInfo.....")
	var request define.QueryUserInfo
	var err error
	responseStatus := &define.ResponseStatus{
		StatusCode: 0,
		StatusMsg:  "SUCCESS",
	}
	ResponeData := define.QueryUserInfoRespone{}
	body, err := ioutil.ReadAll(c.Request.Body)
	logger.Debugf("QueryUserInfo header : %v", c.Request.Header)
	logger.Debugf("QueryUserInfo body : %s", string(body))
	if err != nil {
		responseStatus.StatusCode = 1
		responseStatus.StatusMsg = err.Error()
		ResponeData.ResponseCode = define.Other
		ResponeData.ResponseExplain = err.Error()
		logger.Errorf("QueryUserInfo read body : %s", err.Error())
		utils.Response(ResponeData, c, http.StatusNoContent, responseStatus, nil)
		return
	}

	if err = json.Unmarshal(body, &request); err != nil {
		responseStatus.StatusCode = 1
		responseStatus.StatusMsg = err.Error()
		logger.Errorf("QueryUserInfo Unmarshal : %s", err.Error())
		ResponeData.ResponseCode = define.Other
		ResponeData.ResponseExplain = err.Error()
		utils.Response(ResponeData, c, http.StatusBadRequest, responseStatus, nil)
		return
	}
	logger.Debug(request)

	if err = utils.VerifyQueryUserInfoRequestFormat(&request); err != nil {
		responseStatus.StatusCode = 1
		responseStatus.StatusMsg = err.Error()
		logger.Errorf("VerifyQueryUserInfoRequestFormat: %s", err.Error())
		ResponeData.ResponseCode = define.ParameterError
		ResponeData.ResponseExplain = codeMessage[define.ParameterError] + err.Error()
		utils.Response(ResponeData, c, http.StatusBadRequest, responseStatus, nil)
		return
	}
	b, _ := json.Marshal(request)
	var responseData []byte
	responseData, err = sdk.QueryData(define.QueryUserdata, string(b))
	if err != nil {
		errStr := fmt.Sprintf("Failed to query user info : %s", err.Error())
		ResponeData.ResponseCode = define.Other
		ResponeData.ResponseExplain = errStr
		valueNil := strings.Contains(err.Error(), define.ValueOfKeyNil)
		strinfo := "Failed to query user info :"
		if valueNil {
			ResponeData.ResponseCode = define.ValueOfKeyNil
			ResponeData.ResponseExplain = strinfo + codeMessage[define.ValueOfKeyNil]
		}
		utils.Response(ResponeData, c, http.StatusBadRequest, responseStatus, nil)
		return
	}

	queryRespose := &define.QueryResponse{}
	strInfo := new(string)
	queryRespose.Payload = strInfo
	err = utils.Unmarshal(responseData, queryRespose)
	if err != nil {
		errStr := fmt.Sprintf("Failed to Unmarshal QueryUserInfoResponse:%s", err.Error())
		ResponeData.ResponseCode = define.Other
		ResponeData.ResponseExplain = errStr
		utils.Response(ResponeData, c, http.StatusBadRequest, responseStatus, nil)
		return
	}

	var msg define.UserInfo
	err = utils.Unmarshal([]byte(*strInfo), &msg)
	if err != nil {
		errStr := fmt.Sprintf("Failed to Unmarshal access message:%s", err.Error())
		ResponeData.ResponseCode = define.Other
		ResponeData.ResponseExplain = errStr
		utils.Response(ResponeData, c, http.StatusBadRequest, responseStatus, nil)
		return
	}
	ResponeData.ResponseCode = define.Success
	ResponeData.ResponseExplain = codeMessage[define.Success]
	ResponeData.UserInfo = msg
	utils.Response(ResponeData, c, http.StatusOK, responseStatus, nil)
}
