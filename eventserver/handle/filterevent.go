package handle

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os/exec"
	"strings"
	"time"

	"github.com/peersafe/factoring/apiserver/utils"
	"github.com/peersafe/factoring/common/crypto"
	"github.com/peersafe/factoring/common/metadata"
	"github.com/peersafe/factoring/common/sdk"
	"github.com/peersafe/factoring/define"

	"github.com/op/go-logging"
	"github.com/streadway/amqp"
)

const (
	prefixUserAlias = "user_alias:"
	NoCheck         = "NOCHECK"
	broadcastTag    = "ALL"
	firstStageNum   = 10
	secondStageNum  = 50
	firstStageTime  = 1
	secondStageTime = 5
	MaxFail = 3
)

var (
	logger      = logging.MustGetLogger(metadata.LogModule)
	userAlias   []string
	noCheckSend = false
)

func SetUserAlias(alias string) {
	if "" == alias {
		logger.Warning("alias from file is empty!")
		return
	} else {
		alias = strings.Replace(alias, " ", "", -1)
		userAlias = strings.Split(alias, ",")
		logger.Infof("get %d alias from file, and the alias is %s", len(userAlias), userAlias)
	}

	if len(userAlias) != 0 {
		for _, v := range userAlias {
			if NoCheck == v {
				noCheckSend = true
				logger.Info("get NOCHECK tag, set NOCheckSend to be ture.")
			}
		}
	}

	return
}

func SetUserAliastoFile(configPath, configFile, alias string) bool {
	var command, echoCommand string

	if "" == alias {
		logger.Error("alias is nil, cat set it to the file!")
		return false
	}

	// The command is just for ubuntu and centos, other operating system may use other command to do it.
	// TODO: adapt to other operating systems.
	// echoCommand just to adapt to docker environment
	if strings.HasSuffix(configPath, "/") {
		command = fmt.Sprintf("$(sed 's/^    alias.*/    alias: %s/g' %s%s.yaml)", alias, configPath, configFile)
		echoCommand = fmt.Sprintf("echo %q > %s%s.yaml", command, configPath, configFile)
	} else {
		command = fmt.Sprintf("$(sed 's/^    alias.*/    alias: %s/g' %s/%s.yaml)", alias, configPath, configFile)
		echoCommand = fmt.Sprintf("echo %q > %s/%s.yaml", command, configPath, configFile)
	}

	logger.Info("the echoCommand is:", echoCommand)

	cmd := exec.Command("/bin/bash", "-c", echoCommand)
	err := cmd.Run()
	if nil != err {
		logger.Error("set alias:", alias, "to the file failed:", err.Error())
		return false
	}

	logger.Info("set alias:", alias, "to the file successful")
	return true
}

func GetUserAlias(mqAddrs []string, mqQueue, configPath, configFile string) {
	var tryNum = 0
	var addrsNum = len(mqAddrs)

	if 0 == len(mqAddrs) || "" == mqQueue || "" == configPath || "" == configFile {
		logger.Error("mqAddr or mqQueue or configPath or configFile is nil!")
		return
	}
	logger.Infof("mqAddr is: %s and mqQueue is: %s", mqAddrs, mqQueue)
	logger.Infof("configPaht is: %s, and configFile is: %s", configPath, configFile)
	logger.Infof("The systems has %d addresses.", addrsNum)

	for {
		tryNum++
		logger.Critical("Get userAlias try times:", tryNum)
		mqAddr := mqAddrs[(tryNum-1)%addrsNum]

		if firstStageNum <= tryNum && tryNum < secondStageNum {
			time.Sleep(time.Second * firstStageTime)
		} else if tryNum >= secondStageNum {
			time.Sleep(time.Second * secondStageTime)
		}

		conn, err := amqp.Dial(mqAddr)
		if err != nil {
			logger.Error("Failed to connect to RabbitMQ:", err.Error())
			continue
		}

		channel, err := conn.Channel()
		if err != nil {
			logger.Error("Failed to open a channel:", err.Error())
			conn.Close()
			continue
		}
		err = channel.Qos(1, 0, false)
		if err != nil {
			logger.Error("Failed to set the channle's qos:", err.Error())
			conn.Close()
			continue
		}
		queue, err := channel.QueueDeclare(
			mqQueue, // name
			true,    // durable
			false,   // delete when unused
			false,   // exclusive
			false,   // no-wait
			nil,     // arguments
		)
		if err != nil {
			logger.Error("Failed to declare a queue:", err.Error())
			conn.Close()
			continue
		}

		msgs, err := channel.Consume(
			queue.Name, // queue
			"",         // consumer
			true,       // auto-ack
			false,      // exclusive
			false,      // no-local
			false,      // no-wait
			nil,        // args
		)
		if err != nil {
			logger.Error("Failed to consume a queue:", err.Error())
			conn.Close()
			continue
		}

		for msg := range msgs {
			logger.Infof("Received a message: %s", msg.Body)
			aliasFromRabbitmq := string(msg.Body)
			if strings.HasPrefix(aliasFromRabbitmq, prefixUserAlias) {
				aliasFromRabbitmq = aliasFromRabbitmq[len(prefixUserAlias):len(aliasFromRabbitmq)]
			} else {
				logger.Error("The message form rabbitmq doesn't meet the pre-specified format:", aliasFromRabbitmq)
				continue
			}
			aliasFromRabbitmq = strings.Replace(aliasFromRabbitmq, " ", "", -1)
			SetUserAlias(aliasFromRabbitmq)
			logger.Info("Get the alias:", aliasFromRabbitmq, "from rabbitmq")
			if setOk := SetUserAliastoFile(configPath, configFile, aliasFromRabbitmq); !setOk {
				logger.Error("Get the alias:", aliasFromRabbitmq, "from rabbitmq but set it to the config file failed.")
				continue
			}
			logger.Info("Get the alias", aliasFromRabbitmq, "from rabbitmq and set it to the config file.")
		}

		conn.Close()
	}

	logger.Error("----------GetUserAlias exit----------")
}

func decryptData(ms define.Message, recv string) (string, error) {
	cryptoKey, ok := ms.PeersafeData.Keys[recv]
	if !ok {
		return "", fmt.Errorf("The message is for %s, but it does not has the user's cryptoKey", recv)
	}
	path := define.CRYPTO_PATH + recv + "/" + "enrollment.key"
	privateKey, err := ioutil.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("ReadFile %s key error : %s", recv, err.Error())
	}

	key, err := crypto.EciesDecrypt(cryptoKey, privateKey)
	if err != nil {
		return "", fmt.Errorf("EciesDecrypt %s random key error : %s", recv, err.Error())
	}

	tempData, err := base64.StdEncoding.DecodeString(ms.TxData)
	if err != nil {
		return "", fmt.Errorf("base64 decode txdata failed.")
	}

	data, err := crypto.AesDecrypt(tempData, key)
	if err != nil {
		return "", fmt.Errorf("AesDecrypt %s data error : %s", recv, err.Error())
	}

	return string(data), nil
}

func isSendInfo(receiver []string) (bool, error) {
	logger.Debugf("The receiver is %s.", receiver)
	if noCheckSend {
		return true, nil
	}

	for _, v := range receiver {
		if broadcastTag == v {
			logger.Debug("Receiver contains broadcast tag ALL")
			return true, nil
		}

		if ContainsStr(userAlias, v) {
			logger.Debugf("The request's receiver meet the node's id: %s.", v)
			return true, nil
		}
	}

	logger.Debug("Receiver doesn't contain the node's all user id, will not send to mq")
	return false, nil
}

func ContainsStr(strList []string, str string) bool {
	for _, v := range strList {
		if v == str {
			return true
		}
	}
	return false
}

func ParseTransactions(input []string) (interface{}, error) {
	if len(input) < 2 {
		return nil, fmt.Errorf("get cc input error")
	}

	var request define.InvokeRequest
	err := json.Unmarshal([]byte(input[1]), &request)
	if err != nil {
		return nil, err
	}
	var payload []define.Factor
	message := define.Message{}
	err = json.Unmarshal([]byte(request.Value), &message)
	if err != nil {
		return nil, err
	}
	//currentTime := time.Now()
	//timeDiff := currentTime.UnixNano() / 1000000 - int64(message.CreateTime)
	//fmt.Println(timeDiff)

	err = utils.FormatResponseMessage(sdk.GetUserId(), &payload, &[]define.Message{message})
	if err != nil {
		return nil, err
	}

	eventResponse := define.Event{Payload: payload}
	b, err := json.Marshal(eventResponse)
	if err != nil {
		return nil, err
	}
	logger.Debugf("the msg is %s\n", b)
	tempReceiverList := append(message.Receiver, message.Sender)
	if sdk.GetMqEnable() {
		isSendMq, err := isSendInfo(tempReceiverList)
		if err != nil {
			return nil, err
		} else if !isSendMq {
			return nil, nil
		}
	}
	logger.Infof("the businessNO %s successful with txid: %s is pushed",message.BusinessNo,message.FabricTxId)
	return b, nil
}

func CheckBlockSyncState(checkTime time.Duration, channels []string) {
	var curChainBlockHeight uint64
	var curChainBlockNum uint64
	var err error

	failedTimes := make(map[string]uint8)
	preRecordInfo := make(map[string]uint64)
	channelNum := len(channels)

	logger.Infof("start to check block number on chain and record file interval %v", checkTime)

	ticker := time.NewTicker(checkTime)
	for {
		select {
		case <-ticker.C:
			for i := 0; i < channelNum; i++ {
				recordBlockNum := recordInfos[channels[i]].RecordBlockInfo.BlockNumber
				// if the previous block height in record file is not equal to the current value,
				// which shows that listen gorouting is working
				if recordBlockNum != preRecordInfo[channels[i]] {
					logger.Debugf("the previous record height is %d and current record height is %d, normal",
						preRecordInfo[channels[i]], recordBlockNum)
					preRecordInfo[channels[i]] = recordBlockNum
					continue
				} else {
					logger.Debugf("record height is %d, which stay the same.", recordBlockNum )
				}

				if curChainBlockHeight, err = sdk.GetBlockHeightByEventPeer(channels[i]); err != nil {
					logger.Errorf("get block height for channel %s failed: %s", channels[i], err.Error())
					failedTimes[channels[i]]++
					goto Check
				}
				curChainBlockNum = curChainBlockHeight - 1
				logger.Debugf("block height in peer for channel %s is %d", channels[i], curChainBlockNum)

				if recordBlockNum != curChainBlockNum {
					logger.Warningf("in channel %s, curChainBlockNum=%d, recordBlockNum=%d, which are not the same",
						channels[i], curChainBlockNum, recordBlockNum)
					failedTimes[channels[i]]++
				} else {
					failedTimes[channels[i]] = 0
					continue
				}

				Check:
				logger.Warningf("CheckBlockSyncState failed for %d times.", failedTimes)
				if failedTimes[channels[i]] >= 3 {
					panic(fmt.Errorf("the failed time of checking channel %s has reached the max times(3)", channels[i]))
				}
			}
		}
	}
}
