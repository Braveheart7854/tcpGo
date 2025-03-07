/**
 * @author  tongh
 * @date  2022/7/13 4:49 下午
 */
package server

import (
	"context"
	"errors"
	"net"
	"sync"
	. "tcpGo/common"
	"tcpGo/srv/route"
	. "tcpGo/util/msg"
	"time"
)

type Server struct {
	Conns       map[string]*Connect
	MaxConn     int
	M           *sync.Mutex
	ConnTimeOut time.Duration
}

func NewServer() *Server {
	return &Server{
		Conns:       make(map[string]*Connect),
		MaxConn:     10000,
		M:           new(sync.Mutex),
		ConnTimeOut: 10 * time.Second,
	}
}

func (s *Server) HandleConnection(conn net.Conn) {
	c := NewConnect(conn, OptionsConnect{
		//ReadTimeOut:  75 * time.Second,
		//WriteTimeOut: 75 * time.Second,
		HeartTimeOut: 60 * time.Second,
	})

	//go s.Writer(c.WriteChan)
	go s.Writer(c)
	go s.Handler(c)
	//go s.Reader(conn, c.ReadChan)
	go s.GetReader(c)
}

func (s *Server) Connections(c *Connect) error {
	s.M.Lock()
	defer s.M.Unlock()
	if len(s.Conns) >= s.MaxConn {
		return errors.New("Too many connections!")
	}
	s.Conns[c.ConnNo] = c
	return nil
}

func (s *Server) SaveConnect(c *Connect, msg Request) {
	if msg.Method == "connect" {
		//保存连接
		err := s.Connections(c)
		if err != nil {
			Log(c.ConnNo, err.Error())
			c.Conn.Write([]byte(err.Error()))
			c.Conn.Close()
		} else {
			Log(c.ConnNo, "connect success")
		}
	}
}

/*
	func (s *Server) Reader(conn net.Conn, requestChan chan Request) {
		defer func() {
			if p := recover(); p != nil {
				Log(p)
			}
		}()
		tick := time.NewTicker(s.ConnTimeOut)
		defer tick.Stop()
		for {

			dataLenByte := make([]byte, 4)
			_, err := conn.Read(dataLenByte)
			if err != nil {
				Log(conn.RemoteAddr().String(), " connection error: ", err)
				return
			}
			dataLen := protocol.BytesToInt(dataLenByte)
			//fmt.Println("dataLen=",dataLen)

			data := make([]byte, dataLen)
			_, err = conn.Read(data)
			//fmt.Println(string(data))

			var msg Request
			json.Unmarshal(data, &msg)
			uid := msg.Uid
			if msg.Method == "connect" {
				//保存连接
				err = s.Connections(uid, &Connect{
					ConnNo:    "",
					Uid:       "",
					conn:      conn,
					options:   OptionsConnect{},
					ReadChan:  nil,
					WriteChan: nil,
					TimeTick:  nil,
				})
				if err != nil {
					Log(uid, err.Error())
					conn.Write([]byte(err.Error()))
					conn.Close()
					continue
				} else {
					Log(uid, "connect success")
					go s.HeartBeat(uid, tick)
				}
			}

			tick.Reset(s.ConnTimeOut)
			requestChan <- msg
		}
	}

	func (s *Server) Writer(responseChan chan Response) {
		defer func() {
			if p := recover(); p != nil {
				Log(p)
			}
		}()
		for {
			select {
			case data := <-responseChan:
				connect, ok := s.Conns[data.Uid]
				if !ok {
					Log("连接不存在")
					return
				}
				content, _ := json.Marshal(data.ReturnJson)

				n, err := connect.conn.Write(content)
				if err != nil {
					Log(err)
					continue
				}
				Log(n, data.Uid, data.ReturnJson)
			}
		}
	}

	func (s *Server) HeartBeat(uid string, t *time.Ticker) {
		defer func() {
			if p := recover(); p != nil {
				Log(p)
			}
		}()
		for {
			select {
			case <-t.C:
				//fmt.Println("uid=",uid)
				t.Stop()
				connect, ok := s.Conns[uid]
				if !ok {
					Log("连接不存在")
				}
				connect.conn.Close()
				delete(s.Conns, uid)
				return
			}
		}
	}
*/
func (s *Server) Handler(c *Connect) {
	defer func() {
		if p := recover(); p != nil {
			Log(p)
		}
	}()
	for {
		select {
		case params := <-c.ReadChan:

			go s.Execute(params, c.WriteChan)

		}
	}
}

func (s *Server) Execute(params Request, responseChan chan Response) {
	defer func() {
		if p := recover(); p != nil {
			Log(p)
		}
	}()
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var finish = make(chan Response, 1)
	var panicChan = make(chan interface{}, 1)
	go func() {
		defer func() {
			if p := recover(); p != nil {
				panicChan <- p
			}
		}()

		var returnJson ReturnJson

		//对应的函数
		handler, ok := route.PathUrl[params.Method]
		if !ok {
			returnJson = ReturnJson{Code: MethodErr, Msg: "method error"}
		} else {
			returnJson = handler(params)
		}

		finish <- Response{
			ConnNo:     params.ConnNo,
			ReturnJson: returnJson,
		}
	}()

	select {
	case <-ctx.Done():
		responseChan <- Response{
			ConnNo:     params.ConnNo,
			ReturnJson: ReturnJson{Code: TimeOut, Msg: "服务忙，请稍候"},
		}
	case data := <-finish:
		responseChan <- data
	case p := <-panicChan:
		Log(p)
		responseChan <- Response{
			ConnNo:     params.ConnNo,
			ReturnJson: ReturnJson{Code: ServerErr, Msg: "服务端异常"},
		}
	}
}
