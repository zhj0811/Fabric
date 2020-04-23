package handle

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sync"

	"github.com/zhj0811/fabric/common/sdk"
	"github.com/zhj0811/fabric/define"
	"github.com/zhj0811/fabric/eventserver/messagequeue"

	"github.com/zhj0811/gohfc"
)

const (
	fileSaveName = "current.info"
)

var rwMutex sync.RWMutex
var globalFile *os.File

func SetCurrentInfo(info *define.BlockInfo) error {
	//make a json
	data, err := json.Marshal(info)
	if err != nil {
		return err
	}

	//write into file
	f, err := os.Create(fileSaveName)
	if err != nil {
		return err
	}

	rwMutex.Lock()
	defer rwMutex.Unlock()
	_, err = f.Write(data)
	if err == nil {
		err = f.Sync()
	} else if err == nil {
		err = f.Close()
	}
	return err
}

func GetBlockInfo() (define.BlockInfo, error) {
	var blockInfo define.BlockInfo

	rwMutex.RLock()
	data, err := ioutil.ReadFile(fileSaveName)
	rwMutex.RUnlock()
	if err != nil {
		return blockInfo, err
	}

	//parse from json
	err = json.Unmarshal(data, &blockInfo)
	return blockInfo, err
}

func ListenMsgChannel(fromListen chan define.BlockInfoAll) {
	for {
		select {
		case blockInfo := <-fromListen:
			handleMessage(&blockInfo)
		}
	}
}

//The funciton sends the message to mq firstly and then sets block's number and index to file.
//If send fails, the program will exit to avoid recording wrong number or index.
func handleMessage(msgInfo *define.BlockInfoAll) {
	if nil == msgInfo {
		logger.Error("msgInfo is null, nothing to be handled!")
		return
	}

	if sdk.GetMqEnable() && msgInfo.MsgInfo != nil {
		err := messagequeue.SendMessage(msgInfo.MsgInfo)
		if err != nil {
			logger.Errorf("Send message to mq failed: %s.", err.Error())
			os.Exit(1)
		}
	}
	logger.Debugf("handleMessage is %s", msgInfo.MsgInfo)
	logger.Infof("handleMessage with blockNo: %d and txIndex: %d successfully.", msgInfo.BlockNumber, msgInfo.TxIndex)

	err := SetCurrentInfo(&define.BlockInfo{BlockNumber: msgInfo.BlockNumber, TxIndex: msgInfo.TxIndex})
	if err != nil {
		logger.Criticalf("Set block info(number: %d, index: %d) failed: %s.", msgInfo.BlockNumber, msgInfo.TxIndex, err.Error())
	}
}

func CheckAndRecoverEvent(msgChan chan define.BlockInfoAll) error {
	err, curInfo := GetCurrentInfo()
	if err != nil {
		return fmt.Errorf("get block info from current.info file failed:%s\n", err.Error())
	}
	peerBlockHeight, err := gohfc.GetHandler().GetBlockHeight("")
	if err != nil {
		return err
	}

	logger.Debugf("event processed blockNum: %d txIndex: %d ,and peer node's blockNum:%d",
		curInfo.BlockNumber, curInfo.TxIndex, peerBlockHeight-1)

	for i := curInfo.BlockNumber; i < peerBlockHeight; i++ {
		block, err := gohfc.GetHandler().GetBlockByNumber(i, "")
		if err != nil {
			return err
		}
		plainBlock, err := gohfc.GetHandler().ParseCommonBlock(block)
		if err != nil {
			return err
		}

		for index, tx := range plainBlock.Transactions {
			//本地已经处理过的交易要过滤
			if i == curInfo.BlockNumber && index <= curInfo.TxIndex {
				continue
			}
			msg, err := ParaseInput(tx.ChaincodeSpec.Input.Args)
			if err != nil {
				logger.Errorf("parse %d block with %d index tx failed.", plainBlock.Header.Number, index)
				logger.Error(err)
				continue
			} else {
				blockInfo := define.BlockInfoAll{
					BlockInfo: define.BlockInfo{BlockNumber: plainBlock.Header.Number, TxIndex: index},
					MsgInfo:   msg,
				}
				msgChan <- blockInfo
			}
		}
	}
	return nil
}

func ListenEvent(msgChan chan define.BlockInfoAll) error {
	logger.Debug("enter ListenEvent!")
	ch, err := gohfc.GetHandler().ListenEventFullBlock("")
	if err != nil {
		logger.Errorf("Get ListenEventFullBlock err = %s", err.Error())
		return err
	}
	for {
		select {
		case b := <-ch:
			for index, tx := range b.Transactions {
				msg, err := ParaseInput(tx.ChaincodeSpec.Input.Args)
				if err != nil {
					logger.Error(err)
					continue
				} else {
					blockInfo := define.BlockInfoAll{
						BlockInfo: define.BlockInfo{BlockNumber: b.Header.Number, TxIndex: index},
						MsgInfo:   msg,
					}
					logger.Debugf("ListenEventFullBlock Success : BlockNumber = %d", blockInfo.BlockInfo.BlockNumber)
					msgChan <- blockInfo
				}
			}
		}
	}
}

func GetCurrentInfo() (error, define.BlockInfo) {
	var blockInfo define.BlockInfo
	var err error
	if checkFileIsExist(fileSaveName) {
		globalFile, err = os.OpenFile(fileSaveName, os.O_RDWR, 0666)
		rwMutex.RLock()
		data, err := ioutil.ReadAll(globalFile)
		rwMutex.RUnlock()
		if err != nil {
			return err, blockInfo
		}
		err = json.Unmarshal(data, &blockInfo)
		//parse from json
		if blockInfo.BlockNumber < 2 {
			blockInfo.BlockNumber = 2
		}
		return err, blockInfo
	} else {
		blockInfo.TxIndex = 0
		blockInfo.BlockNumber = 2
		if err = SetCurrentInfo(&blockInfo); err != nil {
			return err, blockInfo
		}
		return nil, blockInfo
	}
}

func checkFileIsExist(filename string) bool {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return false
	}
	return true
}
