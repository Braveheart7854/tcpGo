/**
 * @author  tongh
 * @date  2022/7/20 1:29 下午
 */
package api

import (
	"log"
	"tcpGo/common"
	"tcpGo/util/msg"
)

func Hello(message msg.Request) msg.ReturnJson {

	log.Println(string(message.Data)) // print hello world

	return msg.ReturnJson{
		Code: common.OK,
		Msg:  "success",
		Data: message.Data,
	}
}
