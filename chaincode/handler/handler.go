package handler

import (
	"encoding/json"
	"fmt"

	"github.com/peersafe/factoring/chaincode/utils"
	"github.com/peersafe/factoring/define"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

var KEEPALIVETEST = "test"

var myLogger = shim.NewLogger("hanldler")

func init() {
	myLogger.SetLevel(shim.LogDebug)
}

func SaveData(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	var err error

	//Save (request.Key,request.Value) as k,v
	request := &define.InvokeRequest{}
	if err = json.Unmarshal([]byte(args[0]), request); err != nil {
		return utils.InvokeResponse(stub, err, function, nil, false)
	}

	request.Key = stub.GetTxID()
	err = stub.PutState(request.Key, []byte(request.Value))

	if err != nil {
		myLogger.Errorf("saveData err: %s", err.Error())
		return utils.InvokeResponse(stub, err, function, nil, false)
	}

	//Save (BusinessNo,request.Key) as k,v
	message := &define.Message{}
	err = json.Unmarshal([]byte(request.Value), message)
	if err != nil {
		myLogger.Errorf("Failed to unmarshal request.Value:%s", err)
		return utils.InvokeResponse(stub, err, function, nil, false)
	}

	myLogger.Debugf("message.BusinessNo:%s", message.BusinessNo)
	myLogger.Debugf("request.Key:%s", request.Key)

	err = stub.PutState(message.BusinessNo, []byte(request.Key))
	if err != nil {
		myLogger.Errorf("Failed to Save businessNo and request.Key:%s", err)
		return utils.InvokeResponse(stub, err, function, nil, false)
	}

	return utils.InvokeResponse(stub, nil, function, []string{request.Value}, false)
}

func QueryDataByFabricTxId(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	var err error

	myLogger.Debugf("QueryDataByFabricTxId FabricTxId: %s\n", args[0])
	txdata, err := stub.GetState(args[0])
	if err != nil {
		myLogger.Errorf("Failed to GetState() as key %s: %s", args[0], err)
		utils.QueryResponse(nil, nil, define.Page{})
	}
	//myLogger.Debugf("QueryDataByFabricTxId value: %s\n", b)
	response := []string{string(txdata)}
	return utils.QueryResponse(nil, response, define.Page{})
}

func QueryDataByBusinessNo(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	var err error
	//get fabricTxid by businessNo
	fabricTxId, err := stub.GetState(args[0])
	if err != nil {
		myLogger.Errorf("Failed to GetState() as key %s: %s", args[0], err)
		utils.QueryResponse(nil, nil, define.Page{})
	}
	myLogger.Debugf("QueryDataByBusinessNo fabricTxid: %s businessNo: %s\n", fabricTxId, args[0])
	return QueryDataByFabricTxId(stub, function, []string{string(fabricTxId)})
}

func KeepaliveQuery(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	targetValue, err := stub.GetState(KEEPALIVETEST)
	if err != nil {
		err = fmt.Errorf("ERROR! KeepaliveQuery get failed, err = %s", err.Error())
		return []byte("UnReached"), err
	}

	if string(targetValue) != KEEPALIVETEST {
		err = fmt.Errorf("ERROR! KeepaliveQuery get result is %s", string(targetValue))
		return []byte("UnReached"), err
	}

	return []byte("Reached"), nil
}
