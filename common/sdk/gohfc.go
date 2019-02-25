package sdk

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/peersafe/factoring/common/metadata"
	"github.com/peersafe/factoring/define"

	"github.com/op/go-logging"
	"github.com/peersafe/gohfc"
	"github.com/spf13/viper"
)

var (
	logger = logging.MustGetLogger(metadata.LogModule)
)

func InitSDKs(path, name string) error {
	viper.SetEnvPrefix("core")
	viper.AutomaticEnv()
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)
	configFilePath := filepath.Join(path, name+".yaml")
	viper.SetConfigFile(configFilePath)

	err := viper.ReadInConfig()
	if err != nil {
		return fmt.Errorf("Fatal error when initializing %s config : %s\n", "SDK", err)
	}
	err = gohfc.InitSDK(configFilePath)
	if err != nil {
		return err
	}

	logger.Debugf("config file is %s", configFilePath)
	return nil
}

func SetLogLevel(level, name string) error {
	return gohfc.SetLogLevel(level, name)
}

func GetLogLevel(name string) string {
	return gohfc.GetLogLevel(name)
}

func GetEventChannels() []string {
	return gohfc.GetEventChannels()
}

func GetEventChannelInfos() map[string]gohfc.EventChannelInfo{
	return gohfc.GetEventChannelInfos()
}

func Invoke(in []string,  channelName, chaincodeName string) (string, error) {
	Handler := gohfc.GetHandler()
	res, err := Handler.Invoke(in, channelName, chaincodeName)
	if err != nil {
		logger.Errorf("Invoke Response failed with error: %s", err.Error())
		return "", err
	}
	logger.Info("Invoke Response Status: ", res.Status)
	return res.TxID, err
}

func GetBlockHeightByEndorserPeer(channelName string) (uint64, error) {
	blockHeight, err := gohfc.GetHandler().GetBlockHeight(channelName)
	if err != nil {
		return 0, err
	}

	return blockHeight, nil
}

func GetBlockHeightByEventPeer(channelName string) (uint64, error) {
	if blockHeight, err := gohfc.GetHandler().GetBlockHeightByEventName(channelName); err != nil {
		logger.Errorf("get block height in channel %s failed: %s", channelName, err.Error())
		return 0, err
	} else {
		logger.Debugf("the block height in channel %s is %d", channelName, blockHeight)
		return blockHeight, nil
	}
}

func PeerKeepalive(channelName, chaincodeName string) error {
	var info []string
	info = append(info, define.KeepaliveQuery, "reduPara")
	Handler := gohfc.GetHandler()
	res, err := Handler.Query(info, channelName, chaincodeName)
	if err != nil || len(res) == 0 {
		return fmt.Errorf("peer cann't found the value by key")
	}
	if res[0].Error != nil {
		return res[0].Error
	  } 
	if res[0].Response.Response.GetStatus() != 200 {
		return fmt.Errorf("peer status is wrong!")
	}
	keepaliveResult := string(res[0].Response.Response.Payload)
	if keepaliveResult == "Reached" {
		return nil
	} else {
		return fmt.Errorf("peer cann't be reached")
	}
}

func OrderKeepalive() error {
	isconnect, err := gohfc.GetHandler().GetOrdererConnect()
	if err != nil {
		return err
	}
	if !isconnect {
		return fmt.Errorf("orderer cann't be reached")
	}
	return nil
}

func Query(in []string, channelName, chaincodeName string) ([]byte, error) {
	Handler := gohfc.GetHandler()
	res, err := Handler.Query(in, channelName, chaincodeName)
	if err != nil {
		return nil, err
	} else {
		if res[0].Error != nil {
			return nil, err
		}
	}
	return res[0].Response.Response.Payload, nil
}

func GetUserId() string {
	return viper.GetString("user.id")
}

func GetMqEnable() bool {
	return viper.GetBool("mq.mqEnable")
}

func GetApiserverAuthorization() map[string]string {
	return viper.GetStringMapString("apiserver.authorization")
}
