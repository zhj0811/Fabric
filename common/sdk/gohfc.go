package sdk

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/zhj0811/fabric/common/metadata"
	"github.com/zhj0811/fabric/define"

	logging "github.com/op/go-logging"
	"github.com/zhj0811/gohfc"
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
	level := logging.GetLevel(name).String()
	logger.Debugf("GetLogLevel level: %s, LogModule: %s\n", level, name)
	return level
}

func Invoke(in []string) (string, error) {
	Handler := gohfc.GetHandler()
	res, err := Handler.Invoke(in, "", "")
	if err != nil {
		return "", err
	}
	logger.Info("Invoke Response Status : ", res.Status)
	return res.TxID, err
}

func GetBlockHeightByEndorserPeer() (uint64, error) {
	blockHeight, err := gohfc.GetHandler().GetBlockHeight("")
	if err != nil {
		return 0, err
	}

	return blockHeight, nil
}

func GetBlockHeightByEventPeer() (uint64, error) {
	blockHeight, err := gohfc.GetHandler().GetBlockHeightByEventName("")
	if err != nil {
		return 0, err
	}

	return blockHeight, nil
}

func PeerKeepalive() error {
	var info []string
	info = append(info, define.KeepaliveQuery, "reduPara", "reduPara")
	Handler := gohfc.GetHandler()
	res, err := Handler.Query(info, "", "")
	if err != nil {
		return err
	} else {
		keepaliveResult := string(res[0].Response.Response.Payload)
		if keepaliveResult == "Reached" {
			return nil
		} else {
			return fmt.Errorf("peer cann't be reached")
		}
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

func QueryData(dateType string, txid string) ([]byte, error) {
	var info []string
	info = append(info, dateType, txid)
	return Query(info)
}

func Query(in []string) ([]byte, error) {
	Handler := gohfc.GetHandler()
	res, err := Handler.Query(in, "", "")
	if err != nil {
		return nil, err
	} else {
		if res[0].Error != nil {
			fmt.Println("error info %s", res[0].Error.Error())
			return nil, res[0].Error
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
