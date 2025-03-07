/**
 * @author  tongh
 * @date  2024/11/8 14:48
 */
package server

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"tcpGo/common"
	"tcpGo/util/log"
	"tcpGo/util/msg"
	"tcpGo/util/protocol"
)

func (s *Server) GetReader(c *Connect) {

	go s.HeartBeatNew(c)

	for {

		if c.Conn == nil {
			break
		}
		data, err := s.ReadData(c)
		if err != nil {
			if err == io.EOF || strings.Contains(err.Error(), "use of closed network connection") {
				common.Log("连接未建立或已断开")
				break
			}
			if errors.Is(err, os.ErrDeadlineExceeded) {
				s.SendRetry(c)
				break
			}
			//common.Log(err)
			log.MainLogger.Error(err.Error())
			break
		}

		s.SaveConnect(c, data)

		//fmt.Println(data)
		c.ReadChan <- data
		s.ResetTimeTick(c)
	}
}

func (s *Server) HeartBeatNew(c *Connect) {
	defer func() {
		if p := recover(); p != nil {
			common.Log(p)
		}
	}()
	for range c.TimeTick.C {
		c.TimeTick.Stop()
		err := c.Conn.Close()
		if err != nil {
			common.Log(err)
		}
		delete(s.Conns, c.ConnNo)
		return
	}
}

func (s *Server) ResetTimeTick(c *Connect) {
	c.TimeTick.Reset(c.Options.HeartTimeOut)
}

func (s *Server) ReadData(c *Connect) (req msg.Request, err error) {
	dataLenByte := make([]byte, 8)
	_, err = c.Conn.Read(dataLenByte)
	if err != nil {
		common.Log(c.Conn.RemoteAddr().String(), " read error: ", err)
		return
	}
	dataLen := protocol.BytesToInt64(dataLenByte)
	fmt.Println("dataLen=", dataLen)

	data := make([]byte, dataLen)
	var readLen int64 = 0
	for readLen < dataLen {
		var temp = make([]byte, dataLen-readLen)
		var n int
		n, err = c.Conn.Read(temp)
		if err != nil {
			common.Log(c.Conn.RemoteAddr().String(), " read error: ", err)
			return
		}
		//fmt.Println(n, string(temp))
		data = append(data, temp...)
		readLen += int64(n)
	}

	cleanData := bytes.Replace(data, []byte("\x00"), []byte{}, -1)
	//fmt.Println(string(cleanData))
	err = json.Unmarshal(cleanData, &req)
	return
}

func (s *Server) SendRetry(c *Connect) {
	c.ReadChan <- msg.Request{
		ConnNo: c.ConnNo,
		Method: "retry",
		Data:   []byte("read i/o timeout"),
	}
}
