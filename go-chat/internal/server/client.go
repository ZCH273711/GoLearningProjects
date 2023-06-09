package server

import (
	"examples/go-chat/config"
	"examples/go-chat/internal/kafka"
	"examples/go-chat/pkg/common/constant"
	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
)
import "examples/go-chat/pkg/global/log"
import "examples/go-chat/pkg/protocol"

type Client struct {
	Conn *websocket.Conn
	Name string
	Send chan []byte
}

func (c *Client) Read() {
	defer func() {
		MyServer.Ungister <- c
		c.Conn.Close()
	}()

	for {
		c.Conn.PongHandler()
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			log.Logger.Error("client read message error", log.Any("client read message error", err.Error()))
			MyServer.Ungister <- c
			c.Conn.Close()
			break
		}
		msg := &protocol.Message{}
		proto.Unmarshal(message, msg)

		if msg.Type == constant.HEAR_BEAT {
			pong := &protocol.Message{
				Content: constant.PONG,
				Type:    constant.HEAR_BEAT,
			}
			pongByte, err2 := proto.Marshal(pong)
			if err2 != nil {
				log.Logger.Error("client marshal message error", log.Any("client marshal message error", err2.Error()))
			}
			c.Conn.WriteMessage(websocket.BinaryMessage, pongByte)
		} else {
			if config.GetConfig().MsgChannelType.ChannelType == constant.KAFKA {
				kafka.Send(message)
			} else {
				MyServer.Broadcast <- message
			}
		}
	}
}

func (c *Client) Write() {
	defer func() {
		c.Conn.Close()
	}()

	for message := range c.Send {
		c.Conn.WriteMessage(websocket.BinaryMessage, message)
	}
}
