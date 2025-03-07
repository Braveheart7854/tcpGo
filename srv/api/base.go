/**
 * @author  tongh
 * @date  2022/7/20 2:41 下午
 */
package api

import (
	"tcpGo/common"
	"tcpGo/util/msg"
)

func Connect(message msg.Request) msg.ReturnJson {
	return msg.ReturnJson{
		Code: common.OK,
		Msg:  "register success",
		Data: nil,
	}
}

func Ping(message msg.Request) msg.ReturnJson {
	return msg.ReturnJson{
		Code: common.OK,
		Msg:  "pong",
		Data: nil,
	}
}

func Retry(message msg.Request) msg.ReturnJson {
	return msg.ReturnJson{
		Code: common.OK,
		Msg:  "retry",
		Data: nil,
	}
}
