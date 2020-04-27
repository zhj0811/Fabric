package define

const (
	SaveData       = "SaveData"
	QueryData      = "QueryData"
	KeepaliveQuery = "KeepaliveQuery"
	CRYPTO_PATH    = "./crypto/"

	ErrRequestBody    = "request body error"
	ErrMarshalRequest = "marshaljson request error"
	ErrInvoke         = "invoke func error"
	ErrQuery          = "query func error"

	PeerFailed            = 601
	OrdererFailed         = 602
	KafkaNormal           = 603
	KafkaConfigFailed     = 604
	KafkaConnectionFailed = 605
	KafkaBrokerAbnormal   = 606
	LogModuleInvalid      = 607
	LogModuleSetError     = 608

	Success = "Success"

	ParameterError     = "901"
	PermissionNotFound = "902"
	ValueOfKeyNil      = "903"
	NoPermission       = "904"
	Other              = "905"
)
