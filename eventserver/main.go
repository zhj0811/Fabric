package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/zhj0811/fabric/common/metadata"
	"github.com/zhj0811/fabric/common/sdk"
	"github.com/zhj0811/fabric/define"
	"github.com/zhj0811/fabric/eventserver/check"
	"github.com/zhj0811/fabric/eventserver/handle"
	mq "github.com/zhj0811/fabric/eventserver/messagequeue"

	//"github.com/hyperledger/fabric/common/flogging"
	"github.com/op/go-logging"
	"github.com/spf13/viper"
)

// package-scoped constants

const CHANNELBUFFER = 1000

const packageName = "eventserver"

var (
	logOutput  = os.Stderr
	configPath = flag.String("configPath", "./", "config path")
	configName = flag.String("configName", "client_sdk", "config file name")
	isVersion  = flag.Bool("v", false, "Show version information")
	logger     = logging.MustGetLogger(metadata.LogModule)
)

func main() {
	flag.Parse()
	if *isVersion {
		printVersion()
		return
	}

	err := sdk.InitSDKs(*configPath, *configName)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	/*
		//setup system-wide logging backend based on settings from core.yaml
		flogging.InitBackend(flogging.SetFormat(viper.GetString("logging.format")), logOutput)
		logging.SetLevel(logging.DEBUG, "client_sdk")
	*/
	checkTime := viper.GetDuration("other.check_time")
	logger.Infof("checkTime is %v.", checkTime)
	userAlias := viper.GetString("user.id")
	handle.SetUserAlias(userAlias)

	mqEnable := viper.GetBool("mq.mqEnable")
	logger.Debugf("mq enable is %v.", mqEnable)
	if mqEnable {
		mqAddresses := viper.GetStringSlice("mq.mqAddress")
		queueName := viper.GetString("mq.queueName")
		systemQueueName := viper.GetString("mq.systemQueueName")

		if len(mqAddresses) == 0 {
			logger.Panic("The mq_address is empty!")
		}

		if !mq.InitMQ(queueName, mqAddresses...) {
			logger.Panic("init message queue failed!")
		}
		defer mq.Close()

		go handle.GetUserAlias(mqAddresses, systemQueueName, *configPath, *configName)
	}

	msgChan := make(chan define.BlockInfoAll, CHANNELBUFFER)
	go handle.ListenMsgChannel(msgChan)
	handle.CheckAndRecoverEvent(msgChan)
	go handle.ListenEvent(msgChan)

	//listen the block event and parse the message
	//go le.ListenEvent(eventAddress, chainID, handle.FilterEvent, listenToHandle)
	//check and recover the message
	//go handle.CheckAndRecoverEvent(peerClients, chainID, handle.FilterEvent, listenToHandle, currentBlockHeight)

	check.CheckRecover(checkTime)
}

func printVersion() {
	version := metadata.GetVersionInfo()
	fmt.Println(packageName, " with version: ", version)
	fmt.Println()
}
