package sdk

type connect struct {
	serverAddr         string
	sendChan, recvChan chan *Message
}

func newConnect(serverAddr string) *connect {
	return &connect{
		serverAddr: serverAddr,
		sendChan:   make(chan *Message),
		recvChan:   make(chan *Message),
	}
}

// 向目标通道发送消息
func (c *connect) send(msg *Message) {
	c.recvChan <- msg
}

// 将想要接受的通道的消息返回
func (c *connect) recv() <-chan *Message {
	return c.recvChan
}

func (c *connect) close() {

}
