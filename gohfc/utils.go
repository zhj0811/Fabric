package gohfc

import (
	"fmt"
	"math/rand"
	"time"
)

func getChainCodeObj(args []string, transientMap map[string][]byte, channelName, chaincodeName string) (*ChainCode, error) {
	if len(channelName) == 0 {
		channelName = handler.client.Channel.ChannelId
	}
	if len(chaincodeName) == 0 {
		chaincodeName = handler.client.Channel.ChaincodeName
	}
	mspId := handler.client.Channel.LocalMspId
	if channelName == "" || chaincodeName == "" || mspId == "" {
		return nil, fmt.Errorf("channelName or chaincodeName or mspId is empty")
	}

	chaincode := ChainCode{
		ChannelId:    channelName,
		Type:         ChaincodeSpec_GOLANG,
		Name:         chaincodeName,
		Args:         args,
		TransientMap: transientMap,
	}

	return &chaincode, nil
}

//解析背书策略
func parsePolicy() error {
	policyOrgs := handler.client.Channel.Orgs
	policyRule := handler.client.Channel.Rule
	if len(policyOrgs) == 0 || policyRule == "" {
		for _, v := range handler.client.Peers {
			peerNames = append(peerNames, v.Name)
		}
	}
	for ordname := range handler.client.Orderers {
		orderNames = append(orderNames, ordname)
	}
	for _, v := range handler.client.EventPeers {
		eventName = v.Name
		break
	}
	if len(policyOrgs) > 0 {
		for _, v := range handler.client.Peers {
			if containsStr(policyOrgs, v.OrgName) {
				orgPeerMap[v.OrgName] = append(orgPeerMap[v.OrgName], v.Name)
				if policyRule == "or" {
					orRulePeerNames = append(orRulePeerNames, v.Name)
				}
			}
		}
	}

	return nil
}

func getSendOrderName() string {
	return orderNames[generateRangeNum(0, len(orderNames))]
}

func getSendPeerName() []string {
	if len(orRulePeerNames) > 0 {
		return []string{orRulePeerNames[generateRangeNum(0, len(orRulePeerNames))]}
	}
	if len(peerNames) > 0 {
		return peerNames
	}
	var sendNameList []string
	policyRule := handler.client.Channel.Rule
	if policyRule == "and" {
		for _, peerNames := range orgPeerMap {
			sendNameList = append(sendNameList, peerNames[generateRangeNum(0, len(peerNames))])
			continue
		}
	}

	return sendNameList
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
