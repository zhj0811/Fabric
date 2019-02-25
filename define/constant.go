package define

const (
	SaveData = "SaveData"
	KeepaliveQuery = "KeepaliveQuery"
	QueryDataByFabricTxId = "QueryDataByFabricTxId"
	QueryDataByBusinessNo = "QueryDataByBusinessNo"
	DSL_QUERY      = "DslQuery"
	CRYPTO_PATH    = "./crypto/"

	ChannelName = "Channelname"
	ChaincodeName = "Chaincodename"
	ChaincodeVersion = "Chaincodeversion"

	PeerFailed            = 601
	OrdererFailed         = 602
	KafkaNormal           = 603
	KafkaConfigFailed     = 604
	KafkaConnectionFailed = 605
	KafkaBrokerAbnormal   = 606
	ReadReuqestError      = 607
	UnmarshalError        = 608
	LogModuleInvalid      = 609
	LogModuleSetError     = 610
)

type CheckFabricBaseInfo int
const (
	NOCFBI CheckFabricBaseInfo = iota
	CNameCFBI
	CCNameCFBI
	CCVersionCFBI
	CNameAndCCNameCFBI
)
