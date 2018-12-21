package handle

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/peersafe/tradetrain/apiserver/utils"
	"github.com/peersafe/tradetrain/common/metadata"
	"github.com/peersafe/tradetrain/define"

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

func ParaseInput(input []string) (interface{}, error) {
	if len(input) < 2 {
		return nil, fmt.Errorf("get cc input error")
	}
	var request define.InvokeRequest
	var b []byte
	err := json.Unmarshal([]byte(input[1]), &request)
	if err != nil {
		return nil, err
	}

	message := define.Message{}
	err = json.Unmarshal([]byte(request.Value), &message)
	if err != nil {
		return nil, err
	}
	if message.BusinessData != "" {
		var payload define.Factor
		err = utils.FormatResponseMessage(&payload, &message)
		if err != nil {
			return nil, err
		}
		eventResponse := define.Event{Payload: payload}
		b, err = json.Marshal(eventResponse)
		if err != nil {
			return nil, err
		}
		logger.Debugf("the msg is %s\n", b)
		return b, nil
	}

	accessMmessage := define.AccessMessage{}
	err = json.Unmarshal([]byte(request.Value), &accessMmessage)
	if err != nil {
		return nil, err
	}
	if len(accessMmessage.ReaderList) != 0 {
		var payload define.Access
		err = utils.FormatResponseAccessMessage(&payload, &accessMmessage)
		if err != nil {
			return nil, err
		}
		eventResponse := define.Event{Payload: payload}
		b, err = json.Marshal(eventResponse)
		if err != nil {
			return nil, err
		}
		logger.Debugf("the msg is %s\n", b)
		return b, nil
	}

	userInfoMessage := define.UserInfoMessage{}
	err = json.Unmarshal([]byte(request.Value), &userInfoMessage)
	if err != nil {
		return nil, err
	}
	if userInfoMessage.UserID != "" {
		payload := userInfoMessage.UserInfo
		eventResponse := define.Event{Payload: payload}
		b, err = json.Marshal(eventResponse)
		if err != nil {
			return nil, err
		}
		logger.Debugf("the msg is %s\n", b)
		return b, nil
	}

	return b, nil
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

		if containsStr(userAlias, v) {
			logger.Debugf("The request's receiver meet the node's id: %s.", v)
			return true, nil
		}
	}

	logger.Debug("Receiver doesn't contain the node's all user id, will not send to mq")
	return false, nil
}

func containsStr(strList []string, str string) bool {
	for _, v := range strList {
		if v == str {
			return true
		}
	}
	return false
}
