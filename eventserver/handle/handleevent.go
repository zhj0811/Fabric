package handle

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sync"

	"github.com/peersafe/factoring/common/sdk"
	"github.com/peersafe/factoring/define"
	"github.com/peersafe/factoring/eventserver/messagequeue"

	"github.com/peersafe/gohfc"

)

type RecordInfo struct {
	Channel string `json:"channel"`
	RWMutex sync.RWMutex `json:"rwMutex"`
	GlobalFile *os.File `json:"globalFile"`
	RecordBlockInfo *define.BlockInfo `json:"recordBlockInfo"`
}

const (
	constantFileName = ".current.info"
)

var recordInfos = make(map[string]*RecordInfo)
var blockHeights = make(map[string]uint64)

func (recordInfo *RecordInfo)SetCurrentInfo(info *define.BlockInfo) error {
	var err error
	var data []byte

	if data, err = json.Marshal(info); err != nil {
		logger.Errorf("marshal block info failed: %s", err.Error())
		return err
	}

	fileName := getRecordFileName(recordInfo.Channel)
	// if the file handle is empty, first try to open it
	// we do not close the file handle in order to increase efficiency
	if recordInfo.GlobalFile == nil {
		if recordInfo.GlobalFile, err = os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0666); err != nil {
			logger.Errorf("open file failed: %s", err.Error())
			return err
		}
	}

	// use the opened file handle to record the data
	recordInfo.RWMutex.Lock()
	defer recordInfo.RWMutex.Unlock()
	if err := recordInfo.GlobalFile.Truncate(0); err != nil {
		logger.Errorf("truncate to file %s failed: %s", recordInfo.Channel, err.Error())
		return err
	}
	if _, err = recordInfo.GlobalFile.WriteAt(data, 0); err != nil {
		logger.Errorf("write to file %s failed: %s", recordInfo.Channel, err.Error())
		return err
	}
	if err := recordInfo.GlobalFile.Sync(); err != nil {
		logger.Errorf("sync to file %s failed: %s", recordInfo.Channel, err.Error())
		return err
	}
	recordInfo.RecordBlockInfo = info

	return nil
}

func (recordInfo *RecordInfo)GetRecordBlockInfo() *define.BlockInfo {
	recordInfo.RWMutex.RLock()
	defer recordInfo.RWMutex.RUnlock()

	return recordInfo.RecordBlockInfo
}

//The function sends the message to mq firstly and then sets block's number and index to file.
//If send fails, the program will exit to avoid recording wrong number or index.
func handleMessage(msgInfo *define.BlockInfoAll, channelName string) error{
	if sdk.GetMqEnable() && msgInfo.MsgInfo != nil {
		err := messagequeue.SendMessage(channelName, msgInfo.MsgInfo)
		if err != nil {
			return fmt.Errorf("send message to mq failed: %s", err.Error())
		}
		logger.Debugf("Send message to mq success!  blockNo: %d txindex: %d", msgInfo.BlockNumber, msgInfo.TxIndex)
	}
	err := recordInfos[channelName].SetCurrentInfo(&define.BlockInfo{BlockNumber: msgInfo.BlockNumber, TxIndex: msgInfo.TxIndex})
	if err != nil {
		return fmt.Errorf("set block info(number: %d, index: %d) failed: %s", msgInfo.BlockNumber, msgInfo.TxIndex, err.Error())
	}
	logger.Debugf("handleMessage with blockNo: %d and txIndex: %d success.", msgInfo.BlockNumber, msgInfo.TxIndex)
	return nil
}


func ListenEvent(channelName string, chaincodeName map[string]bool) error {
	recordInfo := recordInfos[channelName].RecordBlockInfo
	logger.Infof("for channel %s, the last records height is %d and txIndex is %d, and the chaincodes are %s",
		channelName, recordInfo.BlockNumber, recordInfo.TxIndex, chaincodeName)

	ch, err := gohfc.GetHandler().ListenEventFullBlock(channelName, int(recordInfo.BlockNumber))
	if err != nil {
		logger.Errorf("listen event for full block failed:%s", err.Error())
		return err
	}
	for {
		select {
		case b := <-ch:
			for index, tx := range b.Transactions {
				msg := interface{}(nil)
				logger.Debugf("ListenEventFullBlock BlockNumber= %d, TxIndex=%d", b.Header.Number, index)
				if recordInfo.BlockNumber == b.Header.Number && recordInfo.TxIndex > index {
					logger.Debugf("the record txindex=%d > current txindex=%d, wait for next", recordInfo.TxIndex, index)
					continue
				}
				if _, ok := chaincodeName[tx.ChaincodeSpec.ChaincodeId.Name]; ok {
					msg, err = ParseTransactions(tx.ChaincodeSpec.Input.Args)
					if err != nil {
						panic(fmt.Errorf("Parse Transactions failed with chaincode name: %s and error is: %s",
							tx.ChaincodeSpec.ChaincodeId.Name, err.Error()))
					}
				} else {
					logger.Debugf("the transaction chaincode:%s is not in the scope of monitoring", tx.ChaincodeSpec.ChaincodeId.Name)
				}
				blockInfo := define.BlockInfoAll{
					BlockInfo: define.BlockInfo{BlockNumber: b.Header.Number, TxIndex: index},
					MsgInfo:   msg,
				}
				if err := handleMessage(&blockInfo, channelName); err != nil {
					panic(err)
				}
				logger.Debugf("handleMessage is pushed successful %s\n", msg)
			}
		}
	}
}

func InitRecordInfos(channels []string) (error) {
	var err error

	channelNum := len(channels)
	for i:= 0; i < channelNum; i++ {
		var tempRecordInfo RecordInfo
		var tempBlockInfo define.BlockInfo

		logger.Debugf("try to init channel(%s)'s record info", channels[i])
		fileName := getRecordFileName(channels[i])
		// if a record file already exists, read the file to init block info for the channel
		// otherwise, the channel has not been listened and set the block info to the default value
		if checkFileIsExist(fileName) {
			if tempRecordInfo.GlobalFile, err = os.OpenFile(fileName, os.O_RDWR, 0666); err != nil {
				logger.Errorf("open file(%s) failed: %s", fileName, err.Error())
				return err
			}
			data, err := ioutil.ReadAll(tempRecordInfo.GlobalFile)
			if err != nil {
				logger.Errorf("read file(%s) failed: %s", fileName, err.Error())
				return err
			}
			if err = json.Unmarshal(data, &tempBlockInfo); err != nil {
				logger.Errorf("unmarshal to block info failed: %s", err.Error())
				return nil
			}
			logger.Infof("channel(%s) has record height %d and index %d",
				channels[i], tempBlockInfo.BlockNumber, tempBlockInfo.TxIndex)
		}

		if tempBlockInfo.BlockNumber < 2 {
			tempBlockInfo.BlockNumber = 2
			tempBlockInfo.TxIndex = 0
			logger.Infof("set the block number to 2 for channel(%s) forcibly", channels[i])
		}

		tempRecordInfo.RecordBlockInfo = &tempBlockInfo
		tempRecordInfo.Channel = channels[i]
		recordInfos[channels[i]] = &tempRecordInfo
	}

	logger.Infof("Init channel(%s %d)'s record info successfully", channels, channelNum)
	return nil
}

func InitBlockHeight(channels []string)  error {
	channelNum := len(channels)
	logger.Debugf("try to get the block height for channel %s", channels)

	for i := 0; i < channelNum; i++ {
		if blockHeight, err := sdk.GetBlockHeightByEventPeer(channels[i]); err != nil {
			logger.Errorf("get block height for channel %s failed: %s", channels[i], err.Error())
			return err
		} else {
			blockHeights[channels[i]] = blockHeight
			logger.Infof("get block height for channel %s is %d", channels[i], blockHeight)
		}
	}

	logger.Infof("get all the block height for channel(%s-%d) successfully", channels, channelNum)
	return nil
}

func checkFileIsExist(filename string) bool {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return false
	}
	return true
}

func getRecordFileName(channel string) string {
	return channel + constantFileName
}
