/**
 * @author  tongh
 * @date  2022/7/20 11:33 上午
 */
package route

import (
	. "tcpGo/srv/api"
	"tcpGo/util/msg"
)

var PathUrl = make(map[string]HandlerFunc)

type HandlerFunc func(message msg.Request) msg.ReturnJson

func addPath(method string, function HandlerFunc) {
	PathUrl[method] = function
}

func InitPath() {
	addPath("connect", Connect)
	addPath("ping", Ping)
	addPath("retry", Retry)

	addPath("hello", Hello)

	addPath("uploadFile", UploadFile)
}
