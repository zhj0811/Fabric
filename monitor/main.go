package main

import (
	"encoding/binary"
	"flag"
	"fmt"
	"os"

	"github.com/zhaojianpeerfintech/fabric/common/metadata"
	"github.com/zhaojianpeerfintech/fabric/common/sdk"
	//	"github.com/zhaojianpeerfintech/fabric/define"
	//"github.com/zhaojianpeerfintech/tradetrain/eventserver/check"
	//	"github.com/zhaojianpeerfintech/tradetrain/eventserver/handle"

	//"github.com/hyperledger/fabric/common/flogging"
	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/fabric/protos/common"
	ab "github.com/hyperledger/fabric/protos/orderer"
	"github.com/op/go-logging"
	"github.com/peersafe/gohfc"
	"github.com/spf13/viper"
)

// package-scoped constants

const CHANNELBUFFER = 1000

const packageName = "monitor"

var (
	logOutput  = os.Stderr
	configPath = flag.String("configPath", "./", "config path")
	configName = flag.String("configName", "client_sdk", "config file name")
	isVersion  = flag.Bool("v", false, "Show version information")
	logger     = logging.MustGetLogger(metadata.LogModule)
	blockNum   = flag.Uint64("blockNum", 4, "get block number")
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

	checkTime := viper.GetDuration("other.check_time")
	logger.Infof("checkTime is %v.", checkTime)

	block := &common.Block{}
	block, err = gohfc.GetHandler().GetBlockByNumber(*blockNum, "")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(block)
	fmt.Println()

	valueMetadata := &common.Metadata{}
	if err = proto.Unmarshal(block.Metadata.Metadata[1], valueMetadata); err != nil {
		return
	}
	lastConfig := &common.LastConfig{}
	if err = proto.Unmarshal(valueMetadata.Value, lastConfig); err != nil {
		return
	}
	fmt.Println(lastConfig)

	if err = proto.Unmarshal(block.Metadata.Metadata[3], valueMetadata); err != nil {
		return
	}
	kafkaMetadata := &ab.KafkaMetadata{}
	if err = proto.Unmarshal(valueMetadata.Value, kafkaMetadata); err != nil {
		return
	}
	fmt.Println(kafkaMetadata)

	b := getValueFromBlockMetadata(block, 2)
	fmt.Println(b)
	fmt.Println()
	fmt.Println([]uint8(block.Metadata.Metadata[2]))

	b = getValueFromBlockMetadata(block, 3)
	fmt.Println("metadata[3]'s value:", b)
	fmt.Println()

	//	msgChan := make(chan define.BlockInfoAll, CHANNELBUFFER)
	//	go handle.ListenMsgChannel(msgChan)
	//	handle.CheckAndRecoverEvent(msgChan)
	//	go handle.ListenEvent(msgChan)

	//listen the block event and parse the message
	//go le.ListenEvent(eventAddress, chainID, handle.FilterEvent, listenToHandle)
	//check and recover the message
	//go handle.CheckAndRecoverEvent(peerClients, chainID, handle.FilterEvent, listenToHandle, currentBlockHeight)

	//	check.CheckRecover(checkTime)
}

func printVersion() {
	version := metadata.GetVersionInfo()
	fmt.Println(packageName, " with version: ", version)
	fmt.Println()
}

func getValueFromBlockMetadata(block *common.Block, index common.BlockMetadataIndex) []byte {
	valueMetadata := &common.Metadata{}
	if index == common.BlockMetadataIndex_LAST_CONFIG {
		if err := proto.Unmarshal(block.Metadata.Metadata[index], valueMetadata); err != nil {
			return nil
		}

		lastConfig := &common.LastConfig{}
		if err := proto.Unmarshal(valueMetadata.Value, lastConfig); err != nil {
			return nil
		}
		b := make([]byte, 8)
		binary.LittleEndian.PutUint64(b, uint64(lastConfig.Index))
		return b
	} else if index == common.BlockMetadataIndex_ORDERER {
		if err := proto.Unmarshal(block.Metadata.Metadata[index], valueMetadata); err != nil {
			return nil
		}

		kafkaMetadata := &ab.KafkaMetadata{}
		if err := proto.Unmarshal(valueMetadata.Value, kafkaMetadata); err != nil {
			return nil
		}
		b := make([]byte, 8)
		binary.LittleEndian.PutUint64(b, uint64(kafkaMetadata.LastOffsetPersisted))
		return b
	} else if index == common.BlockMetadataIndex_TRANSACTIONS_FILTER {
		return block.Metadata.Metadata[index]
	}
	return valueMetadata.Value
}