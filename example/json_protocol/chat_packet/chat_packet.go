package chat_packet

const (
	CHAT_PACKET_VERSION uint32 = 100
	PID_REQ_LOGIN       uint32 = 1001
	PID_RES_LOGIN       uint32 = 1002
	PID_REQ_USER_LIST   uint32 = 1003
	PID_RES_USER_LIST   uint32 = 1004
	PID_REQ_SEND_CHAT   uint32 = 1005
	PID_NOTIFY_CHAT     uint32 = 1006
)

const (
	RESULT_OK   = 0
	RESULT_FAIL = 1
)

type UserInfo struct {
	Nickname string `json:"nickname"`
}

// Request Login.
type ReqLogin struct {
	PacketId uint32 `json:"packet_id"`
	Nickname string `json:"nickname"`
}

func NewReqLogin() ReqLogin {
	return ReqLogin{PacketId: PID_REQ_LOGIN}
}

// Response Login.
type ResLogin struct {
	PacketId uint32 `json:"packet_id"`
	Result   uint16 `json:"result"`
}

func NewResLogin() ResLogin {
	return ResLogin{PacketId: PID_RES_LOGIN}
}

// Request user list
type ReqUserList struct {
	PacketId uint32 `json:"packet_id"`
}

func NewReqUserList() ReqUserList {
	return ReqUserList{PacketId: PID_REQ_USER_LIST}
}

// Response user list
type ResUserList struct {
	PacketId      uint32     `json:"packet_id"`
	NumberOfUsers int        `json:"number_of_users"`
	Users         []UserInfo `json:"users"`
}

func NewResUserList() ResUserList {
	return ResUserList{PacketId: PID_RES_USER_LIST}
}

// Send chat message.
// Response isn't required.
type ReqSendChat struct {
	PacketId uint32 `json:"packet_id"`
	Chat     string `json:"chat"`
}

func NewReqSendChat() ReqSendChat {
	return ReqSendChat{PacketId: PID_REQ_SEND_CHAT}
}

// Notify chat message to all users
type NotifyChat struct {
	PacketId       uint32 `json:"packet_id"`
	SenderNickname string `json:"sender_nickname"`
	Chat           string `json:"chat"`
}

func NewNotifyChat() NotifyChat {
	return NotifyChat{PacketId: PID_NOTIFY_CHAT}
}
