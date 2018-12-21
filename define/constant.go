package define

const (
	SaveData       = "SaveData"
	SaveACL        = "SaveACL"
	KeepaliveQuery = "KeepaliveQuery"
	QueryDataByKey = "QueryDataByKey"
	QueryListById  = "QueryListById"
	SaveUserInfo   = "SaveUserInfo"
	QueryUserdata  = "QueryUserdata"
	DSL_QUERY      = "DslQuery"
	CRYPTO_PATH    = "./crypto/"
	SaveTable      = "SaveTable"

	PeerFailed            = 601
	OrdererFailed         = 602
	KafkaNormal           = 603
	KafkaConfigFailed     = 604
	KafkaConnectionFailed = 605
	KafkaBrokerAbnormal   = 606
	LogModuleInvalid      = 607
	LogModuleSetError     = 608

	Success            = "900"
	ParameterError     = "901"
	PermissionNotFound = "902"
	ValueOfKeyNil      = "903"
	NoPermission       = "904"
	Other              = "905"
)
