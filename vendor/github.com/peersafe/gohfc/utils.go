package gohfc

import (
	"fmt"
	"github.com/op/go-logging"
	"math/rand"
	"os"
	"time"
)


func getChainCodeObj(args []string, channelName, chaincodeName string) (*ChainCode, error) {
	mspId := handler.client.Channel.LocalMspId
	if channelName == "" || chaincodeName == "" || mspId == "" {
		return nil, fmt.Errorf("channelName or chaincodeName or mspId is empty")
	}

	chaincode := ChainCode{
		ChannelId: channelName,
		Type:      ChaincodeSpec_GOLANG,
		Name:      chaincodeName,
		Args:      args,
	}

	return &chaincode, nil
}

//设置log级别
func SetLogLevel(level, name string) error {
	format := logging.MustStringFormatter("%{shortfile} %{time:2006-01-02 15:04:05.000} [%{module}] %{level:.4s} : %{message}")
	backend := logging.NewLogBackend(os.Stderr, "", 0)
	backendFormatter := logging.NewBackendFormatter(backend, format)
	logLevel, err := logging.LogLevel(level)
	if err != nil {
		return err
	}
	logging.SetBackend(backendFormatter).SetLevel(logLevel, name)
	logger.Debugf("SetLogLevel level: %s, levelName: %s\n", level, name)
	return nil
}

func GetLogLevel(name string) string {
	level := logging.GetLevel(name).String()
	logger.Debugf("GetLogLevel level: %s, LogModule: %s\n", level, name)
	return level
}

//解析背书策略
func parsePolicy() error {
	// orderer and event peer do not have relationship with endorsement policy
	// init orderer name
	for ordname := range handler.client.Orderers {
		orderNames = append(orderNames, ordname)
	}
	// init event peer name
	for _, v := range handler.client.EventPeers {
		eventName = v.Name
		break
	}
    //init org's all peer name
	for _, v := range handler.client.Peers {
		if (!containsStr(orgPeerMap[v.OrgName], v.Name)){
			orgPeerMap[v.OrgName] = append(orgPeerMap[v.OrgName], v.Name)
		}
	}
	
	for channelName, channelInfo := range handler.client.Channel.ChannelInfos {
		for chaincodeName, ccInfo := range  channelInfo.CCInfos {
			policyOrgs := ccInfo.Orgs
			//policyRule := ccInfo.Rule
			// not specify org in endorsement policy
			// TODO think wrong configuration
			if len(policyOrgs) == 0 {
				for _, v := range handler.client.Peers {
					peerNames[channelName+chaincodeName] = append(peerNames[channelName+chaincodeName], v.Name)
				}
			} else {
				for _, v := range handler.client.Peers {
					if containsStr(policyOrgs, v.OrgName) {
						rulePeerNames[channelName+chaincodeName] = append(rulePeerNames[channelName+chaincodeName], v.Name)
						logger.Debugf("add peer %s to rulePeerNames in channel %s and chaincode %s", v.Name, channelName, chaincodeName)
						if !containsStr(channelPeerNames[channelName], v.Name) {
							channelPeerNames[channelName] = append(channelPeerNames[channelName], v.Name)
							logger.Debugf("add peer %s to channelPeerNames in channel %s", v.Name, channelName)
						}
					}
				}
			}
		}
	}
	return nil
}

func getSendOrderName() string {
	return orderNames[generateRangeNum(0, len(orderNames))]
}

// we can support "and" and "or" policies only currently
func getSendPeerName(channelName, chaincodeName string) []string {
	var sendNameList []string

	// not specify chaincode name, choose one peer randomly in this channel's peer
	if "" == chaincodeName {
		return []string{channelPeerNames[channelName][generateRangeNum(0, len(channelPeerNames[channelName]))]}
	}

	targetName := channelName + chaincodeName
	numRulePeerNames := len(rulePeerNames[targetName])
	policyRule := handler.client.Channel.ChannelInfos[channelName].CCInfos[chaincodeName].Rule

	//not specify org, send to all peers
	if len(peerNames[targetName]) != 0 {
		return peerNames[targetName]
	}

	if policyRule == "and" {
		// and policy, choose one peer randomly in every org
		for _, orgName := range handler.client.Channel.ChannelInfos[channelName].CCInfos[chaincodeName].Orgs {
			logger.Debugf("orgName is %s and orgPeerMap is %s", orgName, orgPeerMap[orgName])
			sendNameList = append(sendNameList, orgPeerMap[orgName][generateRangeNum(0, len(orgPeerMap[orgName]))])
		}
		return sendNameList
	} else {
		// or policy, choose one peer randomly
		return []string{rulePeerNames[targetName][generateRangeNum(0, numRulePeerNames)]}
	}
}

func generateRangeNum(min, max int) int {
	rand.Seed(time.Now().Unix())
	randNum := rand.Intn(max-min) + min
	return randNum
}

func containsStr(strList []string, str string) bool {
	for _, v := range strList {
		if v == str {
			return true
		}
	}
	return false
}
