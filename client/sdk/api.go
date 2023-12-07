package sdk

// 设置IM基础API 数据（信息、聊天session结构体）

const (
	MsgType_Text = "text"
)

type Chat struct {
	Nick      string
	UserID    string
	SessionID string
	conn      *connect
}

type Message struct {
	Type       string
	Name       string
	FromUserID string
	ToUserId   string
	Content    string
	Session    string
}

func NewChat(serverAddr, nick, userID, sessionID string) *Chat {
	return &Chat{
		Nick:      nick,
		UserID:    userID,
		SessionID: sessionID,
		conn:      newConnect(serverAddr),
	}
}

func (chat *Chat) Send(msg *Message) {
	chat.conn.send(msg)
}

func (chat *Chat) Recv() <-chan *Message {
	return chat.conn.recv()
}

func (chat *Chat) Close() {
	chat.conn.close()
}
