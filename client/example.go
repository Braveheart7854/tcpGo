/**
 * @author  tongh
 * @date  2025/3/12 15:13
 */
package main

import (
	"encoding/json"
	"fmt"
	"tcpGo/client/util"
	"tcpGo/util/msg"
	"time"
)

func main() {
	var cs = make([]util.Client, 0)
	for i := 0; i < 100; i++ {
		c, err := util.NewClient("127.0.0.1", 46001)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer c.Close()

		var heartbeatDone = make(chan struct{})
		go c.SendHeartbeat(heartbeatDone)

		go c.GetReader()

		c.Auth(nil)

		cs = append(cs, c)
	}

	for i := 0; i < 100; i++ {
		data := msg.Request{
			Method: "hello",
			Data:   []byte(fmt.Sprintf("hello %d", i)),
		}
		databyte, _ := json.Marshal(data)
		cs[i].Send(databyte)
	}

	time.Sleep(10 * time.Second)
	//close(heartbeatDone)

	select {}
}
