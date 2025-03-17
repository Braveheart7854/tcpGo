/**
 * @author  tongh
 * @date  2024/11/8 15:47
 */
package server

import (
	"net"
	"tcpGo/common"
	"tcpGo/util/msg"
	"time"
)

type Connect struct {
	ConnNo        string
	Uid           string
	Conn          net.Conn
	Options       OptionsConnect
	ReadChan      chan msg.Request
	WriteChan     chan msg.Response
	TimeTick      *time.Ticker
	HeartBeatChan chan struct{}
}

type OptionsConnect struct {
	ReadTimeOut  time.Duration // read timeout
	WriteTimeOut time.Duration // write timeout
	HeartTimeOut time.Duration
}

func NewConnect(conn net.Conn, options OptionsConnect) *Connect {
	//conn.SetReadDeadline(time.Now().Add(options.ReadTimeOut))
	//conn.SetWriteDeadline(time.Now().Add(options.WriteTimeOut))

	return &Connect{
		ConnNo:        common.MD5(conn.RemoteAddr().String() + "_" + common.GetRandomString(10)),
		Conn:          conn,
		Options:       options,
		ReadChan:      make(chan msg.Request),
		WriteChan:     make(chan msg.Response),
		TimeTick:      time.NewTicker(options.HeartTimeOut),
		HeartBeatChan: make(chan struct{}),
	}
}
