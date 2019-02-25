package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/zhaojianpeerfintech/fabric/apiserver/utils"
	"github.com/zhaojianpeerfintech/fabric/common/sdk"
	"github.com/zhaojianpeerfintech/fabric/define"

	"github.com/gin-gonic/gin"
)


func SaveRegistration(c *gin.Context) {
	logger.Debug("SaveRegistration......")

	var request define.Registration
	responseStatus := &define.ResponseStatus{
		StatusCode: 0,
		StatusMsg:  "SUCCESS",
	}
	var err error

	logger.Debugf("SaveRegistration header : %v", c.Request.Header)
	fabricBaseInfo, err := parseFabricBaseInfo(c.Request.Header, define.CNameAndCCNameCFBI)
	if err != nil {
		responseStatus.StatusCode = http.StatusBadRequest
		responseStatus.StatusMsg = err.Error()
		logger.Errorf("parseFabricBaseInfo failed: %s", err.Error())
		utils.Response(responseStatus, c, http.StatusBadRequest, responseStatus, nil)
		return
	}

	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		responseStatus.StatusCode = 1
		responseStatus.StatusMsg = err.Error()
		logger.Errorf("SaveRegistration read body failed: %s", err.Error())
		utils.Response(responseStatus, c, http.StatusBadRequest, responseStatus, nil)
		return
	}
	logger.Debugf("SaveRegistration body : %s", string(body))

	if err = json.Unmarshal(body, &request); err != nil {
		responseStatus.StatusCode = 1
		responseStatus.StatusMsg = err.Error()
		logger.Errorf("SaveRegistration Unmarshal failed : %s", err.Error())
		utils.Response(responseStatus, c, http.StatusBadRequest, responseStatus, nil)
		return
	}


	b, err := json.Marshal(request)
	if err != nil {
		responseStatus.StatusCode = 1
		responseStatus.StatusMsg = err.Error()
		logger.Errorf("SaveRegistration Marshal : %s", err.Error())
		utils.Response(responseStatus, c, http.StatusBadRequest, responseStatus, nil)
		return
	}

	// invoke
	var info []string
	info = append(info, define.SaveRegistration, string(b))
	txId, err := sdk.Invoke(info, fabricBaseInfo.ChannelName, fabricBaseInfo.ChaincodeName)
	if err != nil {
		responseStatus.StatusCode = 1
		responseStatus.StatusMsg = err.Error()
		logger.Errorf("SaveRegistration Invoke : %s", err.Error())
		utils.Response(responseStatus, c, http.StatusBadRequest, responseStatus, nil)
		return
	}
	logger.Infof("invoke the advanceOrderNo %v successful with txid: %s", request.AdvanceOrderNo, txId)

	utils.Response(txId, c, http.StatusOK, responseStatus, nil)
}


func QueryRegistration(c *gin.Context) {
	logger.Debug("QueryRegistration...")
	responseStatus := &define.ResponseStatus{
		StatusCode: 0,
		StatusMsg:  "SUCCESS",
	}
	var err error

	logger.Debugf("QueryRegistration header : %v", c.Request.Header)
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
	fabricTxID := values.Get("fabricTxID")
	advanceOrderNo := values.Get("advanceOrderNo")
	logger.Debug("fabricTxID:", fabricTxID)
	logger.Debug("advanceOrderNo:", advanceOrderNo)

	if fabricTxID == "" && advanceOrderNo == "" {
		errStr := fmt.Sprintf("fabricTxID and advanceOrderN are all empty")
		utils.Response(errStr, c, http.StatusBadRequest, responseStatus, nil)
		return
	}

	var responseData []byte
	

	if fabricTxID != "" {
		logger.Infof("query data with fabricTxID: %s", fabricTxID)
		responseData, err = sdk.Query([]string{define.QueryDataByFabricTxID, fabricTxID}, fabricBaseInfo.ChannelName, fabricBaseInfo.ChaincodeName)
	} else if advanceOrderNo != "" {
		logger.Infof("query data with advanceOrderNo: %s", advanceOrderNo)
		responseData, err = sdk.Query([]string{define.QueryRegistrationByNo, advanceOrderNo}, fabricBaseInfo.ChannelName, fabricBaseInfo.ChaincodeName)
	}

	if err != nil {
		logger.Errorf("Failed to query data: %s", err.Error())
		responseStatus.StatusCode = http.StatusBadRequest
		responseStatus.StatusMsg = err.Error()
		utils.Response(nil, c, http.StatusBadRequest, responseStatus, nil)
		return
	}

	queryRespose := &define.QueryResponse{}
	str := new(string)
	queryRespose.Payload = str
	err = utils.Unmarshal(responseData, queryRespose)
	if err != nil {
		logger.Errorf("Failed to Unmarshal QueryResponse: %s", err.Error())
		responseStatus.StatusCode = http.StatusBadRequest
		responseStatus.StatusMsg = err.Error()
		utils.Response(nil, c, http.StatusBadRequest, responseStatus, nil)
		return
	}
	msg := &define.Registration{}
	err = utils.Unmarshal([]byte(*str), &msg)
	if err != nil {
		logger.Errorf("Failed to Unmarshal message: %s", err.Error())
		responseStatus.StatusCode = http.StatusBadRequest
		responseStatus.StatusMsg = err.Error()
		utils.Response(nil, c, http.StatusBadRequest, responseStatus, nil)
		return
	}
	utils.Response(msg, c, http.StatusOK, responseStatus, nil)
}

