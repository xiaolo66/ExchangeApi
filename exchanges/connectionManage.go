package exchanges

import (
	"fmt"
	"github.com/xiaolo66/ExchangeApi/exchanges/websocket"
	"sync"
	"github.com/xiaolo66/ExchangeApi"
	set "github.com/deckarep/golang-set"
)

type ConnectFunc func(url string) (*Connection, error)
type Connection struct {
	websocket.WsConn
	MsgChannels set.Set
}

func NewConnection() *Connection {
	return &Connection{
		MsgChannels: set.NewSet(),
	}
}

func (c *Connection) Subscribe(msgChan ExchangeApi.MessageChan) {
	c.MsgChannels.Add(msgChan)
}

func (c *Connection) UnSubscribe(msgChan ExchangeApi.MessageChan) {
	c.MsgChannels.Remove(msgChan)
}

func (c *Connection) Close() {
	c.WsConn.Close()
}

func (c *Connection) Publish(msg ExchangeApi.Message, clear bool) {
	tmp := c.MsgChannels
	if clear {
		c.MsgChannels = set.NewSet()
	}
	tmp.Each(func(item interface{}) bool {
		msgChan, ok := item.(ExchangeApi.MessageChan)
		if ok && msgChan != nil {
			//must use go routine here, otherwise the "Each" method may be blocked, caused dead lock if someone call Subscribe/UnSubscribe at same time.
			go func() { msgChan <- msg }()
		}
		return false
	})
}

type ConnectionManager struct {
	sync.RWMutex
	once  sync.Once
	conns map[string]*Connection // key: ws url
}

func NewConnectionManager() *ConnectionManager {
	return &ConnectionManager{conns: make(map[string]*Connection)}
}

func (c *ConnectionManager) SetConnection(url string, connection *Connection) {
	c.Lock()
	defer c.Unlock()
	conn, ok := c.conns[url]
	if ok {
		conn.Close()
	}
	c.conns[url] = connection
}

func (c *ConnectionManager) RemoveConnection(url string) {
	c.Lock()
	defer c.Unlock()
	delete(c.conns, url)
}

func (c *ConnectionManager) Close() {
	c.Lock()
	defer c.Unlock()
	for _, conn := range c.conns {
		conn.Close()
	}
	c.conns = make(map[string]*Connection)
}

func (c *ConnectionManager) GetConnection(url string, connectFunc ConnectFunc) (*Connection, error) {
	c.Lock()
	defer c.Unlock()
	conn, ok := c.conns[url]
	if !ok {
		if connectFunc != nil {
			var err error
			conn, err = connectFunc(url)
			if err != nil {
				return nil, err
			}
			c.conns[url] = conn
			return conn, nil
		}
		return nil, fmt.Errorf("not found websocket session, url:%s", url)
	}
	return conn, nil
}

func (c *ConnectionManager) Publish(url string, message ExchangeApi.Message) {
	conn, _ := c.GetConnection(url, nil)
	if conn != nil {
		conn.Publish(message, false)
	}
}

// PublishAfterClear clear the subscribers and notify them
func (c *ConnectionManager) PublishAfterClear(url string, message ExchangeApi.Message) {
	conn, _ := c.GetConnection(url, nil)
	if conn != nil {
		conn.Publish(message, true)
	}
}

