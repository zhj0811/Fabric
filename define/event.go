package define

type Event struct {
	Payload interface{} `json:"payload"`
}

type Header struct {
	ContentDef     ContentDef     `json:"contentDef"`
	Ack            Ack            `json:"ack"`
	ResponseStatus ResponseStatus `json:"responseStatus"`
}

type ContentDef struct {
	ContentType string `json:"contentType"`
	TrackId     string `json:"trackId"`
	Language    string `json:"language"`
}

type Ack struct {
	Level    string `json:"level"`
	Callback string `json:"callback"`
}

type Contents struct {
	Schema  string      `json:"$schema"`
	Payload interface{} `json:"payload"`
	Command Command     `json:"command, omitempty"`
}

type Command struct {
	Uri    string `json:"uri, omitempty"`
	Action string `json:"action, omitempty"`
	Desc   string `json:"desc, omitempty"`
}

type ResponseData struct {
	TrackId        string         `json:"trackId"`
	ResponseStatus ResponseStatus `json:"responseStatus"`
	Page           Page           `json:"page"`
	Payload        interface{}    `json:"payload"`
}

type BlockInfo struct {
	BlockNumber uint64 `json:"block_number"`
	TxIndex     int    `json:"tx_index"`
}

type BlockInfoAll struct {
	BlockInfo
	MsgInfo interface{} `json:"msgInfo"`
}
