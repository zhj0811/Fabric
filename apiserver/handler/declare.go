package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/peersafe/tradetrain/apiserver/utils"
	"github.com/peersafe/tradetrain/common/sdk"
	"github.com/peersafe/tradetrain/define"

	"github.com/gin-gonic/gin"
)

func SaveCustomsForm(c *gin.Context) {
	logger.Debugf("SaveData.....")
	var request define.CustomsDeclarationInfo
	var err error
	responseStatus := &define.ResponseStatus{
		StatusCode: 0,
		StatusMsg:  "SUCCESS",
	}
	ResponeData := define.QueryFormDataResponse{}

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

	b, err := utils.FormatRequestFormMessage(request)
	if err != nil {
		responseStatus.StatusCode = 1
		responseStatus.StatusMsg = err.Error()
		logger.Errorf("SaveData FormatRequestMessage : %s", err.Error())
		ResponeData.ResponseCode = define.Other
		ResponeData.ResponseExplain = err.Error()
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
	ResponeData.EntryID = request.EntryID
	utils.Response(ResponeData, c, http.StatusOK, responseStatus, nil)
}

// QueryFormDataByKey 报关单表体、表头信息查询
func QueryFormDataByKey(c *gin.Context) {
	var request define.QueryFormData
	var err error
	responseStatus := &define.ResponseStatus{
		StatusCode: 0,
		StatusMsg:  "SUCCESS",
	}
	ResponeData := define.QueryFormDataResponse{}

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
	if request.BusinessType == "" || request.DataType == "" || request.Key == "" || request.Reader == "" || request.WriteRoleType == "" || request.EntryID == "" {
		ResponeData.ResponseCode = define.ParameterError
		ResponeData.ResponseExplain = codeMessage[define.ParameterError]
		responseStatus.StatusCode = 1
		responseStatus.StatusMsg = codeMessage[define.ParameterError]
		logger.Errorf("QueryData: %s", codeMessage[define.ParameterError])
		utils.Response(ResponeData, c, http.StatusBadRequest, responseStatus, nil)
		return
	}

	logger.Debug(request)
	var info []string
	info = append(info, define.QueryDataByKey, request.Key, request.BusinessType, request.DataType, request.WriteRoleType, request.Reader, request.EntryID)
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
	var msg define.CustomsDeclarationMessage
	err = utils.Unmarshal([]byte(queryRespose.Payload.(string)), &msg)
	if err != nil {
		errStr := fmt.Sprintf("Failed to Unmarshal message:%s", err.Error())
		ResponeData.ResponseCode = define.Other
		ResponeData.ResponseExplain = errStr
		utils.Response(ResponeData, c, http.StatusBadRequest, responseStatus, nil)
		return
	}

	var deckaretionList define.CustomsDeclarationInfo
	if err := utils.FormatResponseFormMessage(sdk.GetUserId(), &deckaretionList, &msg); err != nil {
		errStr := fmt.Sprintf("Failed to FormatResponseAccessMessage:%s", err.Error())
		ResponeData.ResponseCode = define.Other
		ResponeData.ResponseExplain = errStr
		utils.Response(ResponeData, c, http.StatusBadRequest, responseStatus, nil)
		return
	}
	ResponeData.ResponseCode = "900"
	ResponeData.ResponseExplain = "Success"
	ResponeData.FactorBaseInfo = deckaretionList.FactorBaseInfo
	ResponeData.EntryID = deckaretionList.EntryID
	ResponeData.BusinessData = deckaretionList.BusinessData
	utils.Response(ResponeData, c, http.StatusOK, responseStatus, nil)
}
