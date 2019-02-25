package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/zhaojianpeerfintech/fabric/apiserver/utils"
	"github.com/zhaojianpeerfintech/fabric/common/metadata"
	"github.com/zhaojianpeerfintech/fabric/common/sdk"
	"github.com/zhaojianpeerfintech/fabric/define"

	"github.com/gin-gonic/gin"
	"github.com/op/go-logging"
)

var logger = logging.MustGetLogger(metadata.LogModule)

// SaveData 保存保理信息
func SaveData(c *gin.Context) {
	logger.Debug("SaveData......")

	var request []define.Factor
	var txIds []string
	var err error
	responseStatus := &define.ResponseStatus{
		StatusCode: 0,
		StatusMsg:  "SUCCESS",
	}

	logger.Debugf("SaveData header : %v", c.Request.Header)
	fabricBaseInfo, err := parseFabricBaseInfo(c.Request.Header, define.CNameAndCCNameCFBI)
	if err != nil {
		responseStatus.StatusCode = http.StatusBadRequest
		responseStatus.StatusMsg = err.Error()
		logger.Errorf("parseFabricBaseInfo failed: %s", err.Error())
		utils.Response(nil, c, http.StatusBadRequest, responseStatus, nil)
		return
	}

	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		responseStatus.StatusCode = 1
		responseStatus.StatusMsg = err.Error()
		logger.Errorf("SaveData read body failed: %s", err.Error())
		utils.Response(nil, c, http.StatusNoContent, responseStatus, nil)
		return
	}
	logger.Debugf("SaveData body : %s", string(body))
	if err = json.Unmarshal(body, &request); err != nil {
		responseStatus.StatusCode = 1
		responseStatus.StatusMsg = err.Error()
		logger.Errorf("SaveData Unmarshal : %s", err.Error())
		utils.Response(nil, c, http.StatusBadRequest, responseStatus, nil)
		return
	}
	logger.Debug(request)

	for _, factor := range request {
		logger.Infof("receive factor with businessNO: %s", factor.BusinessNo)
		b, err := utils.FormatRequestMessage(factor)
		if err != nil {
			responseStatus.StatusCode = 1
			responseStatus.StatusMsg = err.Error()
			logger.Errorf("SaveData FormatRequestMessage : %s", err.Error())
			utils.Response(nil, c, http.StatusBadRequest, responseStatus, nil)
			return
		}

		// invoke
		var info []string
		info = append(info, define.SaveData, string(b))
		logger.Infof("invoke in channel %s and chaincode %s", fabricBaseInfo.ChannelName, fabricBaseInfo.ChaincodeName)
		txId, err := sdk.Invoke(info, fabricBaseInfo.ChannelName, fabricBaseInfo.ChaincodeName)
		if err != nil {
			responseStatus.StatusCode = 1
			responseStatus.StatusMsg = err.Error()
			logger.Errorf("SaveData Invoke failed: %s and txid is %s", err.Error(), txId)
			utils.Response(nil, c, http.StatusBadRequest, responseStatus, nil)
			return
		}
		logger.Infof("invoke the businessNO %s successful with txid: %s", factor.BusinessNo, txId)
		txIds = append(txIds, txId)
	}

	utils.Response(txIds, c, http.StatusOK, responseStatus, nil)
}

func QueryByInfo(c *gin.Context) {
	logger.Debug("QueryByInfo...")

	var err error
	var responseData []byte
	var infos []string
	responseStatus := &define.ResponseStatus{
		StatusCode: 0,
		StatusMsg:  "SUCCESS",
	}

	logger.Debugf("SaveData header : %v", c.Request.Header)
	fabricBaseInfo, err := parseFabricBaseInfo(c.Request.Header, define.CNameAndCCNameCFBI)
	if err != nil {
		responseStatus.StatusCode = http.StatusBadRequest
		responseStatus.StatusMsg = err.Error()
		logger.Errorf("parseFabricBaseInfo failed: %s", err.Error())
		utils.Response(nil, c, http.StatusBadRequest, responseStatus, nil)
		return
	}

	// query
	values, _ := url.ParseQuery(c.Request.URL.RawQuery)
	businessNo := values.Get("businessno")
	fabricTxId := values.Get("fabrictxid")
	logger.Debug("fabricTxId:", fabricTxId)
	logger.Debug("businessNo:", businessNo)

	if businessNo == "" && fabricTxId == "" {
		errStr := fmt.Sprintf("businessno and fabrictxid are all empty")
		utils.Response(errStr, c, http.StatusBadRequest, responseStatus, nil)
		return
	}

	if fabricTxId != "" {
		logger.Infof("query data with txid: %s", fabricTxId)
		infos = append(infos, define.QueryDataByFabricTxId, fabricTxId)
	} else if businessNo != "" {
		logger.Infof("query data with businessNO: %s", businessNo)
		infos = append(infos, define.QueryDataByBusinessNo, businessNo)
	}
	if responseData, err = sdk.Query(infos, fabricBaseInfo.ChannelName, fabricBaseInfo.ChaincodeName); err != nil {
		logger.Errorf("Failed to query data: %s", err.Error())
		errStr := fmt.Sprintf("Failed to query data: %s", err.Error())
		utils.Response(errStr, c, http.StatusBadRequest, responseStatus, nil)
		return
	}

	queryRespose := &define.QueryResponse{}
	strlist := new([]string)
	queryRespose.Payload = strlist
	err = utils.Unmarshal(responseData, queryRespose)
	if err != nil {
		logger.Errorf("Failed to Unmarshal QueryResponse: %s", err.Error())
		errStr := fmt.Sprintf("Failed to Unmarshal QueryResponse:%s", err.Error())
		utils.Response(errStr, c, http.StatusBadRequest, responseStatus, nil)
		return
	}
	var messageList []define.Message
	for _, v := range *strlist {
		var msg define.Message
		err = utils.Unmarshal([]byte(v), &msg)
		if err != nil {
			logger.Errorf("Failed to Unmarshal message: %s", err.Error())
			errStr := fmt.Sprintf("Failed to Unmarshal message:%s", err.Error())
			utils.Response(errStr, c, http.StatusBadRequest, responseStatus, nil)
			return
		}
		messageList = append(messageList, msg)
	}

	var factorList []define.Factor
	if err := utils.FormatResponseMessage(sdk.GetUserId(), &factorList, &messageList); err != nil {
		logger.Errorf("Failed to FormatResponseMessage: %s", err.Error())
		errStr := fmt.Sprintf("Failed to FormatResponseMessage:%s", err.Error())
		utils.Response(errStr, c, http.StatusBadRequest, responseStatus, nil)
		return
	}
	for _, factor := range factorList {
		logger.Infof("find txid(%s) with businessNO(%s)", factor.FabricTxId, factor.BusinessNo)
	}

	utils.Response(factorList, c, http.StatusOK, responseStatus, nil)
}

/*
// DslQuery 按条件查询信息
func DslQuery(c *gin.Context) {
	logger.Debug("DslQuery.....")
	var request define.QueryRequest
	var err error
	responseStatus := &define.ResponseStatus{
		StatusCode: 0,
		StatusMsg:  "SUCCESS",
	}
	body, err := ioutil.ReadAll(c.Request.Body)
	logger.Debugf("DslQuery header : %v", c.Request.Header)
	logger.Debugf("DslQuery body : %s", string(body))
	if err != nil {
		responseStatus.StatusCode = 1
		responseStatus.StatusMsg = err.Error()
		logger.Errorf("DslQuery read body : %s", err.Error())
		utils.Response(nil, c, http.StatusNoContent, responseStatus, nil)
		return
	}

	// query
	status := http.StatusOK
	requestPage := c.Request.Header.Get("page")
	json.Unmarshal([]byte(requestPage), &request.SplitPage)
	request.DslSyntax = string(body)
	response, retStatus, page, err := sdk.Handler.DSL("", define.DSL_QUERY, nil, request)
	if err != nil {
		responseStatus.StatusCode = 1
		responseStatus.StatusMsg = err.Error()
		logger.Errorf("DslQuery Query : %s", err.Error())
		status = http.StatusBadRequest
	} else {
		responseStatus = retStatus
	}
	utils.Response(response, c, status, responseStatus, page)
}

// BlockQuery 获取指定业务编号相关的区块信息
func BlockQuery(c *gin.Context) {
	logger.Debug("BlockQuery.....")
	var request define.QueryRequest
	var err error
	responseStatus := &define.ResponseStatus{
		StatusCode: 0,
		StatusMsg:  "SUCCESS",
	}
	// query
	status := http.StatusOK
	businessNo := c.Param("id")
	requestPage := c.Request.Header.Get("page")
	json.Unmarshal([]byte(requestPage), &request.SplitPage)
	request.DslSyntax = fmt.Sprintf("{\"selector\":{\"businessNo\":\"%s\"}}", businessNo)
	request.BlockFlag = true
	b, _ := json.Marshal(request)
	responseData, retStatus, page, err := sdk.Handler.Query("", define.DSL_QUERY, nil, string(b), true)
	if err != nil {
		responseStatus.StatusCode = 1
		responseStatus.StatusMsg = err.Error()
		logger.Errorf("DslQuery Query : %s", err.Error())
		status = http.StatusBadRequest
	} else {
		responseStatus = retStatus
	}

	txIdList, ok := responseData.Payload.([]string)
	if !ok {
		responseStatus.StatusMsg = err.Error()
		logger.Errorf("DslQuery Query : %s", err.Error())
		status = http.StatusBadRequest
	}

	if len(txIdList) > 0 {
		response, retStatus, _, err := sdk.Handler.BlockQuery(txIdList)
		if err != nil {
			responseStatus.StatusCode = 1
			responseStatus.StatusMsg = err.Error()
			logger.Errorf("DslQuery Query : %s", err.Error())
			status = http.StatusBadRequest
		} else {
			responseStatus = retStatus
		}
		utils.Response(response, c, status, responseStatus, page)
	}
	utils.Response(define.QueryContents{}, c, status, responseStatus, page)
}

// BlockQuery 获取指定业务编号相关的区块信息
func BlockQueryEx(c *gin.Context) {
	logger.Debug("BlockQuery.....")
	var request define.QueryRequest
	var err error
	responseStatus := &define.ResponseStatus{
		StatusCode: 0,
		StatusMsg:  "SUCCESS",
	}
	// query
	status := http.StatusOK
	businessNo := c.Param("id")
	requestPage := c.Request.Header.Get("page")
	json.Unmarshal([]byte(requestPage), &request.SplitPage)
	request.DslSyntax = fmt.Sprintf("{\"selector\":{\"businessNo\":\"%s\"}}", businessNo)
	request.BlockFlag = true
	b, _ := json.Marshal(request)
	responseData, retStatus, page, err := sdk.Handler.Query("", define.DSL_QUERY, nil, string(b), true)
	if err != nil {
		responseStatus.StatusCode = 1
		responseStatus.StatusMsg = err.Error()
		logger.Errorf("DslQuery Query : %s", err.Error())
		status = http.StatusBadRequest
	} else {
		responseStatus = retStatus
	}

	txIdList, ok := responseData.Payload.([]string)
	if !ok {
		responseStatus.StatusMsg = err.Error()
		logger.Errorf("DslQuery Query : %s", err.Error())
		status = http.StatusBadRequest
	}

	if len(txIdList) > 0 {
		response, retStatus, _, err := sdk.Handler.BlockQueryEx(txIdList)
		if err != nil {
			responseStatus.StatusCode = 1
			responseStatus.StatusMsg = err.Error()
			logger.Errorf("DslQuery Query : %s", err.Error())
			status = http.StatusBadRequest
		} else {
			responseStatus = retStatus
		}
		utils.Response(response, c, status, responseStatus, page)
	}
	utils.Response(define.QueryContents{}, c, status, responseStatus, page)
}
*/
