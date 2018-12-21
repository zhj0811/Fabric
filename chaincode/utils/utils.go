package utils

import (
	"encoding/json"

	"github.com/peersafe/tradetrain/define"

	"github.com/hyperledger/fabric/core/chaincode/shim"
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
			myLogger.Errorf("set event error : %s", errTmp.Error())
		}
	}
	return nil, err
}

func QueryResponse(err error, data interface{}, pageItem define.Page) ([]byte, error) {
	response := define.QueryResponse{
		Payload: data,
		Page:    pageItem,
	}

	payload, Marshalerr := json.Marshal(response)
	myLogger.Debug("**************************QueryResponse****************************")
	myLogger.Debug(string(payload))
	if Marshalerr != nil {
		myLogger.Debug("QueryResponse Json  encode error.")
	}

	return payload, err
}
