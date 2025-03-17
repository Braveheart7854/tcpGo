/**
 * @author  tongh
 * @date  2022/7/15 5:34 下午
 */
package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"tcpGo/common"
	"tcpGo/util/msg"
	"tcpGo/util/protocol"
	"time"
)

type Client interface {
	Close() error
	Auth(data []byte)
	Ping() error
	GetReader()
	SendHeartbeat(done chan struct{})
	Send(data []byte) error
}

func NewClient(host string, port int) (cl Client, err error) {
	var conn net.Conn
	server := fmt.Sprintf("%s:%d", host, port)
	conn, err = net.Dial("tcp", server)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		return
	}
	cl = &client{
		Conn:              conn,
		HeartbeatInterval: 30 * time.Second,
		ReconnectInterval: 5 * time.Second,
	}
	return
}

type client struct {
	Conn              net.Conn
	ReadChan          chan msg.Response
	WriteChan         chan msg.Request
	HeartbeatInterval time.Duration // 心跳间隔
	ReconnectInterval time.Duration // 重连间隔
}

func (c *client) Close() error {
	return c.Conn.Close()
}

// 发送心跳包
func (c *client) SendHeartbeat(done chan struct{}) {
	ticker := time.NewTicker(c.HeartbeatInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			err := c.Ping()
			if err != nil {
				fmt.Printf("Failed to send heartbeat: %v\n", err)
				return
			}
			fmt.Println("Sent heartbeat to server")
		case <-done:
			fmt.Println("Stopping heartbeat...")
			return
		}
	}
}

func (c *client) GetReader() {
	for {
		if c.Conn == nil {
			break
		}
		data, err := c.ReadData()
		if err != nil {
			common.Log(err)
			break
		}

		fmt.Println(string(data.ReturnJson.Data))
	}
}

func (c *client) ReadData() (req msg.Response, err error) {
	dataLenByte := make([]byte, 8)
	_, err = c.Conn.Read(dataLenByte)
	if err != nil {
		common.Log(c.Conn.RemoteAddr().String(), " read error: ", err)
		return
	}
	dataLen := protocol.BytesToInt64(dataLenByte)
	//fmt.Println("dataLen=", dataLen)

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

func (c *client) WriteData(msgByte []byte) (n int, err error) {
	data := make([]byte, 0)
	data = append(data, protocol.Int64ToBytes(int64(len(msgByte)))...)
	data = append(data, msgByte...)
	n, err = c.Conn.Write(data)
	return
}

func (c *client) Auth(data []byte) {
	message := msg.Request{
		Method: "connect",
		Data:   data,
	}
	msgByte, _ := json.Marshal(message)

	n1, err1 := c.WriteData(msgByte)
	if err1 != nil {
		fmt.Println(n1, "Write failed:", err1)
		return
	}
}

func (c *client) Ping() error {
	message := msg.Request{
		Method: "ping",
		Data:   []byte("心跳包 ping"),
	}
	msgByte, _ := json.Marshal(message)

	n1, err1 := c.WriteData(msgByte)
	if err1 != nil {
		fmt.Println(n1, "Write failed:", err1)
		return err1
	}
	return nil
}

func (c *client) Send(data []byte) error {
	n1, err1 := c.WriteData(data)
	if err1 != nil {
		fmt.Println(n1, "Write failed:", err1)
		return err1
	}
	return nil
}
