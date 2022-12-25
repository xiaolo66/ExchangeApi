package ExchangeApi

type MessageType int

const (
	MsgOrderBook MessageType = iota
	MsgTicker
	MsgAllTicker
	MsgTrade
	MsgKLine
	MsgBalance
	MsgOrder
	MsgPositions
	MsgMarkPrice

	MsgReConnected //重新建立连接
	MsgDisConnected//网络连接已断开
	MsgClosed //连接已关闭
	MsgError//发生了某种错误
)

type Message struct {
	Type MessageType
	Data interface{}
}
type MessageChan chan Message

var (
	ReConnectedMessage  = Message{Type: MsgReConnected}
	DisConnectedMessage = Message{Type: MsgDisConnected}
	CloseMessage        = Message{Type: MsgClosed}
	ErrorMessage        = func(err error) Message { return Message{Type: MsgError, Data: err} }
)
