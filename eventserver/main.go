package main

import (
	"flag"
	"fmt"

	"github.com/peersafe/factoring/common/metadata"
	"github.com/peersafe/factoring/common/sdk"
	"github.com/peersafe/factoring/eventserver/handle"
	mq "github.com/peersafe/factoring/eventserver/messagequeue"

	"github.com/op/go-logging"
	"github.com/spf13/viper"
)

// package-scoped constants
const packageName = "eventserver"

var (
	//logOutput  = os.Stderr
	configPath = flag.String("configPath", "./", "config path")
	configName = flag.String("configName", "client_sdk", "config file name")
	isVersion  = flag.Bool("v", false, "Show version information")
	logger     = logging.MustGetLogger(metadata.LogModule)
)

func main() {
	var err error

	flag.Parse()
	// print the version info and exit the program
	if *isVersion {
		printVersion()
		return
	}

	if err = sdk.InitSDKs(*configPath, *configName); err != nil {
		fmt.Printf("init sdk failed: %s\n", err.Error())
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

	if !mq.InitMQBaseInfo() {
		logger.Errorf("init mq's base info failed!")
		return
	}
	defer mq.Close()

	eventChannels := sdk.GetEventChannels()
	eventChannelInfos := sdk.GetEventChannelInfos()

	if err := handle.InitRecordInfos(eventChannels); err != nil {
		logger.Errorf("Init record info failed: %s", err.Error())
		return
	}
	if err = handle.InitBlockHeight(eventChannels); err != nil {
		logger.Errorf("Init block height failed: %s", err.Error())
		return
	}

	//listen the block event and parse the message
	for _, channel := range eventChannels {
		chaincodes := make(map[string]bool)
		for _, chaincode := range eventChannelInfos[channel].Chaincodes {
			chaincodes[chaincode] = true
		}
		go handle.ListenEvent(channel, chaincodes)
	}

	//check recover height according to This checkTime
	handle.CheckBlockSyncState(checkTime, eventChannels)
}

func printVersion() {
	version := metadata.GetVersionInfo()
	fmt.Println(packageName, " with version: ", version)
	fmt.Println()
}
