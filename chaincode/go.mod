module github.com/zhj0811/fabric/chaincode

go 1.14

require (
	github.com/Shopify/sarama v1.26.1 // indirect
	github.com/gin-gonic/gin v1.6.2 // indirect
	github.com/hyperledger/fabric-chaincode-go v0.0.0-20200330074746-2584993c3b5e
	github.com/hyperledger/fabric-protos-go v0.0.0-20200422100619-316dc6798e96
	github.com/op/go-logging v0.0.0-20160315200505-970db520ece7 // indirect
	github.com/pkg/errors v0.8.1
	github.com/spf13/viper v1.6.3 // indirect
	github.com/zhj0811/fabric v0.0.0-00010101000000-000000000000
	github.com/zhj0811/fabric/apiserver v0.0.0-20200423084516-4a809feac5d4 // indirect
)

replace github.com/zhj0811/fabric => ../../fabric

replace github.com/zhj0811/fabric/apiserver => ../apiserver
