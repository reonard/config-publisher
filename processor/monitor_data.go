package processor

const (
	STATUS_PENDING = iota
	STATUS_PUSHED
	STATUS_ACK
)

const (
	ActionUpdateConfig = iota + 1
	ActionRequestData
	ActionReset
	ActionSelfCheck
	ActionRequestConfig
)

type MonitorData struct {
	DeviceId          int64                    `json:"ID"`
	DeviceStatus      int8                     `json:"status"`
	TimeStamp         string                   `json:"once"`
	Project           int64                    `json:"Project"`
	Data              []map[string]interface{} `json:"ProbeData"`
	AlarmHandleStatus int8                     `json:"AlarmHandleStatus"`
	AlarmHandle       string                   `json:"AlarmHandle"`
	AlarmDescription  []string                 `json:"AlarmDesc"`
}

type ActionMsg struct {

	DeviceId    int64  `json:"ID"`
	MsgId       int64  `json:"MsgId"`
	ActionType  int8   `json:"ActionType"`
	Secret      string  `json:"Secret"`
	ProbeConfig []map[string]interface{} `json:"ProbeConfig"`
}

type DeviceInfo struct {
	DeviceId   int64
	DeviceName string
	Project    int64
}


type DeviceActionMsg struct {
	DeviceId           int64  `json:"ID"`
	IssueStatusId      int64  `json:"MsgID"`
	ActionType         int8   `json:"ActionType"`
	Secret             string  `json:"Secret"`
	ProbeConfig        map[string]interface{} `json:"ProbeConfig"`
}