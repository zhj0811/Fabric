package handler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/zhj0811/fabric/apiserver/utils"
	"github.com/zhj0811/fabric/common/sdk"
	"github.com/zhj0811/fabric/define"

	"github.com/gin-gonic/gin"
)

var codeMessage = map[string]string{
	"900": "Success",
	"901": "parameter error",
	"902": "Permission not found",
	"903": "The value of key is null",
	"904": "do not have permission",
	"905": "Other",
}

func SaveACL(c *gin.Context) {
	logger.Debugf("SaveACL.....")
	var request define.Access
	var err error
	responseStatus := &define.ResponseStatus{
		StatusCode: 0,
		StatusMsg:  "SUCCESS",
	}

	ResponeData := define.AccessResponse{}
	body, err := ioutil.ReadAll(c.Request.Body)
	logger.Debugf("SaveACL header : %v", c.Request.Header)
	logger.Debugf("SaveACL body : %s", string(body))
	if err != nil {
		responseStatus.StatusCode = 1
		responseStatus.StatusMsg = err.Error()
		logger.Errorf("SaveACL read body : %s", err.Error())
		ResponeData.ResponseCode = define.Other
		ResponeData.ResponseExplain = err.Error()
		utils.Response(ResponeData, c, http.StatusNoContent, responseStatus, nil)
		return
	}

	if err = json.Unmarshal(body, &request); err != nil {
		responseStatus.StatusCode = 1
		responseStatus.StatusMsg = err.Error()
		logger.Errorf("define.Access Unmarshal : %s", err.Error())
		ResponeData.ResponseCode = define.Other
		ResponeData.ResponseExplain = err.Error()
		utils.Response(ResponeData, c, http.StatusBadRequest, responseStatus, nil)
		return
	}
	logger.Debug(request)

	b, err := utils.FormatRequestAccessMessage(request)
	if err != nil {
		responseStatus.StatusCode = 1
		responseStatus.StatusMsg = err.Error()
		logger.Errorf("SaveACL FormatRequestMessage : %s", err.Error())
		ResponeData.ResponseCode = define.ParameterError
		ResponeData.ResponseExplain = codeMessage[define.ParameterError] + err.Error()
		utils.Response(ResponeData, c, http.StatusBadRequest, responseStatus, nil)
		return
	}

	// invoke
	var info []string
	info = append(info, define.SaveACL, string(b))
	txId, err := sdk.Invoke(info)

	if err != nil {
		responseStatus.StatusCode = 1
		responseStatus.StatusMsg = err.Error()
		logger.Errorf("SaveACL Invoke : %s", err.Error())
		ResponeData.ResponseCode = define.Other
		ResponeData.ResponseExplain = err.Error()
		utils.Response(ResponeData, c, http.StatusBadRequest, responseStatus, nil)
		return
	}
	logger.Debugf("This SaveACL tx return's txID is %s", txId)
	ResponeData.ResponseCode = define.Success
	ResponeData.ResponseExplain = codeMessage[define.Success]
	ResponeData.ACLBaseInfo = request.ACLBaseInfo
	utils.Response(ResponeData, c, http.StatusOK, responseStatus, nil)
}

func QueryListById(c *gin.Context) {
	var request define.QueryACL
	var err error
	responseStatus := &define.ResponseStatus{
		StatusCode: 0,
		StatusMsg:  "SUCCESS",
	}
	body, err := ioutil.ReadAll(c.Request.Body)
	logger.Debugf("QueryList header : %v", c.Request.Header)
	logger.Debugf("QueryList body : %s", string(body))
	ResponeData := define.QueryACLResponse{}
	if err != nil {
		responseStatus.StatusCode = 1
		responseStatus.StatusMsg = err.Error()
		logger.Errorf("QueryList read body : %s", err.Error())
		ResponeData.ResponseCode = define.Other
		ResponeData.ResponseExplain = err.Error()
		utils.Response(ResponeData, c, http.StatusNoContent, responseStatus, nil)
		return
	}

	if err = json.Unmarshal(body, &request); err != nil {
		responseStatus.StatusCode = 1
		responseStatus.StatusMsg = err.Error()
		logger.Errorf("QueryListById define.QueryACL Unmarshal : %s", err.Error())
		ResponeData.ResponseCode = define.Other
		ResponeData.ResponseExplain = err.Error()
		utils.Response(ResponeData, c, http.StatusBadRequest, responseStatus, nil)
		return
	}
	logger.Debug(request)
	if err = utils.VerifyQueryAclRequestFormat(&request); err != nil {
		responseStatus.StatusCode = 1
		responseStatus.StatusMsg = err.Error()
		logger.Errorf("verifyQueryAclRequestFormat: %s", err.Error())
		ResponeData.ResponseCode = define.ParameterError
		ResponeData.ResponseExplain = codeMessage[define.ParameterError] + err.Error()
		utils.Response(ResponeData, c, http.StatusBadRequest, responseStatus, nil)
		return
	}
	var info []string
	info = append(info, define.QueryListById, request.Key, request.BusinessType, request.DataType, request.Writer)
	var responseData []byte
	responseData, err = sdk.Query(info)
	if err != nil {
		logger.Errorf("Failed to query access control list: %s", err.Error())
		ResponeData.ResponseCode = define.Other
		strinfo := "Failed to query access control list: "
		ResponeData.ResponseExplain = strinfo + err.Error()
		valueNil := strings.Contains(err.Error(), define.ValueOfKeyNil)
		if valueNil {
			ResponeData.ResponseCode = define.ValueOfKeyNil
			ResponeData.ResponseExplain = strinfo + codeMessage[define.ValueOfKeyNil]
		}
		utils.Response(ResponeData, c, http.StatusBadRequest, responseStatus, nil)
		return
	}

	queryRespose := &define.QueryResponse{}
	strlist := new(string)
	queryRespose.Payload = strlist
	err = utils.Unmarshal(responseData, queryRespose)
	if err != nil {
		logger.Errorf("Failed to Unmarshal QueryResponse:%s", err.Error())
		ResponeData.ResponseCode = define.Other
		ResponeData.ResponseExplain = err.Error()
		utils.Response(ResponeData, c, http.StatusBadRequest, responseStatus, nil)
		return
	}

	var msg define.AccessMessage
	err = utils.Unmarshal([]byte(*strlist), &msg)
	if err != nil {
		logger.Errorf("Failed to Unmarshal access message:%s", err.Error())
		ResponeData.ResponseCode = define.Other
		ResponeData.ResponseExplain = err.Error()
		utils.Response(ResponeData, c, http.StatusBadRequest, responseStatus, nil)
		return
	}

	var accessList define.Access
	if err := utils.FormatResponseAccessMessage(&accessList, &msg); err != nil {
		logger.Errorf("Failed to FormatResponseAccessMessage:%s", err.Error())
		ResponeData.ResponseCode = define.Other
		ResponeData.ResponseExplain = err.Error()
		utils.Response(ResponeData, c, http.StatusBadRequest, responseStatus, nil)
		return
	}
	ResponeData.ResponseCode = define.Success
	ResponeData.ResponseExplain = codeMessage[define.Success]
	ResponeData.Access = accessList
	utils.Response(ResponeData, c, http.StatusOK, responseStatus, nil)
}
