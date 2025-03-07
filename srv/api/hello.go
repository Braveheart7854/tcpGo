/**
 * @author  tongh
 * @date  2022/7/20 1:29 下午
 */
package api

import (
	"encoding/json"
	"fmt"
	"log"
	"tcpGo/common"
	"tcpGo/util/msg"
	"time"
)

func Hello(message msg.Request) msg.ReturnJson {
	//time.Sleep(4*time.Second)
	//log.Println("get ", time.Now(), string(message.Data))

	var data struct {
		FileName string
		FileSize int
		FileNo   int
	}
	json.Unmarshal(message.Data, &data)
	log.Println(fmt.Sprintf("%s get data: fileName: %s, fileSize: %d, fileNo: %d",
		time.Now().Format("2006-01-02 15:04:05"), data.FileName, data.FileSize, data.FileNo))
	return msg.ReturnJson{
		Code: common.OK,
		Msg:  "success",
		Data: message.Data,
	}
}
