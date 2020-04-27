package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/zhj0811/fabric/apiserver/utils"
	"github.com/zhj0811/fabric/common/sdk"
	"github.com/zhj0811/fabric/define"
)

func SaveData(c *gin.Context) {
	logger.Debugf("Enter SaveData...")
	var request define.Factory
	resData := &define.ResData{
		ResCode: 0,
		ResMsg:  define.Success,
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		resData.ResCode = 1
		resData.ResMsg = errors.WithMessage(err, define.ErrRequestBody).Error()
		logger.Error(resData.ResMsg)
		res, _ := resData.MarshalJSON()
		utils.Response(res, c, http.StatusBadRequest)
	}
	logger.Debugf("request is %v", request)

	param, err := request.MarshalJSON()
	if err != nil {
		resData.ResCode = 1
		resData.ResMsg = errors.WithMessage(err, define.ErrMarshalRequest).Error()
		logger.Error(resData.ResMsg)
		res, _ := resData.MarshalJSON()
		utils.Response(res, c, http.StatusBadRequest)
		return
	}

	// invoke
	var info []string
	info = append(info, define.SaveData, string(param))
	txID, err := sdk.Invoke(info)
	if err != nil {
		resData.ResCode = 1
		resData.ResMsg = errors.WithMessage(err, define.ErrInvoke).Error()
		logger.Error(resData.ResMsg)
		res, _ := resData.MarshalJSON()
		utils.Response(res, c, http.StatusBadRequest)
		return
	}
	logger.Infof("This %s func success with txID %s", define.SaveData, txID)
	resData.ResMsg = txID
	res, _ := resData.MarshalJSON()
	utils.Response(res, c, http.StatusOK)
	return
}

func QueryData(c *gin.Context) {
	resData := &define.ResData{
		ResCode: 0,
		ResMsg:  "SUCCESS",
	}

	key := c.Query("key")
	var info []string
	info = append(info, define.QueryData, key)
	result, err := sdk.Query(info)
	if err != nil {
		resData.ResCode = 1
		resData.ResMsg = errors.WithMessage(err, define.ErrQuery).Error()
		logger.Error(resData.ResMsg)
		res, _ := resData.MarshalJSON()
		utils.Response(res, c, http.StatusBadRequest)
		return
	}
	logger.Infof("This %s func success with result %s", define.QueryData, string(result))
	resData.ResMsg = string(result)
	res, _ := resData.MarshalJSON()
	utils.Response(res, c, http.StatusOK)
}

func KeepaliveQuery(c *gin.Context) {
	resData := &define.ResData{
		ResCode: 0,
		ResMsg:  "SUCCESS",
	}

	//if err := sdk.PeerKeepalive(); err != nil {
	//	resData.ResCode = 1
	//	resData.ResMsg = errors.WithMessage(err, "PeerFailed").Error()
	//	logger.Error(resData.ResMsg)
	//	res, _ := resData.MarshalJSON()
	//	utils.Response(res, c, http.StatusBadRequest)
	//	return
	//}

	if err := sdk.OrderKeepalive(); err != nil {
		resData.ResCode = 1
		resData.ResMsg = errors.WithMessage(err, "OrderFailed").Error()
		logger.Error(resData.ResMsg)
		res, _ := resData.MarshalJSON()
		utils.Response(res, c, http.StatusBadRequest)
		return
	}
	res, _ := resData.MarshalJSON()
	utils.Response(res, c, http.StatusOK)
	return
}

func BlockHeight(c *gin.Context) {
	resData := &define.ResData{
		ResCode: 0,
		ResMsg:  "SUCCESS",
	}

	height, err := sdk.GetBlockHeightByEndorserPeer()
	if nil != err {
		resData.ResCode = 1
		resData.ResMsg = errors.WithMessage(err, "ErrGetBlockHeight").Error()
		logger.Error(resData.ResMsg)
		res, _ := resData.MarshalJSON()
		utils.Response(res, c, http.StatusBadRequest)
		return
	}
	resData.ResMsg = strconv.FormatUint(height, 10)
	res, _ := resData.MarshalJSON()
	utils.Response(res, c, http.StatusOK)
	return
}

func KafkaNumber(c *gin.Context) {
	logger.Debug("enter KafkaNumber function.")

	resData := &define.ResData{
		ResCode: 0,
		ResMsg:  "SUCCESS",
	}

	returnCode, err := sdk.GetKafkaNumber()
	if nil != err {
		resData.ResCode = 1
		resData.ResMsg = errors.WithMessage(err, "ErrGetKafkaNumber").Error()
		logger.Error(resData.ResMsg)
		res, _ := resData.MarshalJSON()
		utils.Response(res, c, http.StatusBadRequest)
		return
	}
	resData.ResMsg = strconv.Itoa(returnCode)
	res, _ := resData.MarshalJSON()
	utils.Response(res, c, http.StatusOK)
	return
}
