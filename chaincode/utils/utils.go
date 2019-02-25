package utils

import (
	"encoding/json"
	"github.com/peersafe/factoring/define"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"time"
)

var myLogger = shim.NewLogger("utils")

func init() {
	myLogger.SetLevel(shim.LogDebug)
}

func InvokeResponse(stub shim.ChaincodeStubInterface, err error, function string, data interface{}, eventFlag bool) ([]byte, error) {
	if eventFlag {
		myLogger.Debugf("Set event  %s\n.", function)
		factordata, _ := json.Marshal(data)
		if errTmp := stub.SetEvent(function, factordata); errTmp != nil {
			myLogger.Errorf("the transaction is %s,set event error : %s", stub.GetTxID(),errTmp.Error())
		}
		myLogger.Infof("The invoke response txId is %s",stub.GetTxID())
		myLogger.Debugf("The invoke response time is %d",time.Now().UnixNano()/1000000)
	}
	return nil, err
}

func QueryResponse(err error, data interface{}, pageItem define.Page) ([]byte, error) {
	response := define.QueryResponse{
		Payload: data,
		Page:    pageItem,
	}

	payload, err := json.Marshal(response)
	myLogger.Debug("**************************QueryResponse****************************")
	myLogger.Debug(string(payload))
	if err != nil {
		myLogger.Debug("QueryResponse Json  encode error.")
	}

	return payload, err
}
