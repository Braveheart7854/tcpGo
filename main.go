/**
 * @author  tongh
 * @date  2022/7/13 4:21 下午
 */
package main

import (
	"fmt"
	"net"
	"tcpGo/common"
	"tcpGo/srv/route"
	"tcpGo/util/log"
	"tcpGo/util/server"
	"time"
)

func main() {
	log.Init()
	netListen, err := net.Listen("tcp", ":46001")
	common.CheckError(err)

	defer netListen.Close()

	common.Log("Waiting for clients")

	route.InitPath()

	srv := server.NewServer()

	// 调试用
	go func() {
		for range time.Tick(10 * time.Second) {
			common.Log(fmt.Sprintf("当前连接数为:%d，当前连接列表为：%v", len(srv.Conns), srv.Conns))
		}
	}()

	for {
		conn, err := netListen.Accept()
		if err != nil {
			continue
		}

		tcpConn, ok := conn.(*net.TCPConn)
		if !ok {
			common.Log("not a TCP connection")
			continue
		}

		if err := tcpConn.SetNoDelay(true); err != nil {
			common.Log("failed to set TCP_NODELAY")
			continue
		}

		common.Log(time.Now(), tcpConn.RemoteAddr().String(), " tcp connect success")
		go srv.HandleConnection(tcpConn)
	}
}
