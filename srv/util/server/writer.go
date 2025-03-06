/**
 * @author  tongh
 * @date  2024/11/8 14:49
 */
package server

import (
	"encoding/json"
	"tcpGo/srv/common"
	"tcpGo/srv/util/protocol"
)

func (s *Server) Writer(c *Connect) {
	defer func() {
		if p := recover(); p != nil {
			common.Log(p)
		}
	}()
	for {
		select {
		case data := <-c.WriteChan:
			data.ConnNo = c.ConnNo
			content, _ := json.Marshal(data)

			_, err := s.WriteData(c, content)
			if err != nil {
				common.Log(err)
				continue
			}
			//common.Log(n, data.ConnNo, data.ReturnJson)
		}
	}
}

func (s *Server) WriteData(c *Connect, msgByte []byte) (n int, err error) {
	data := make([]byte, 0)
	data = append(data, protocol.Int64ToBytes(int64(len(msgByte)))...)
	data = append(data, msgByte...)
	n, err = c.Conn.Write(data)
	return
}
