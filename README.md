# tcpGo

- 让网络编程变得及其简单，像http编程一样进行tcp编程，开发者只需关注业务逻辑实现。
- 框架底层已为你处理好tcp粘包、tcp分包、心跳、无用连接自动销毁等基础功能。

## 快速开始
将代码库拉取到本地
```bash
git clone git@github.com:Braveheart7854/tcpGo.git
```

运行服务端代码
```bash
go run ./main.go
```

出现以下提示代表运行成功
```bash
GOROOT=/usr/local/go #gosetup
GOPATH=/Users/XXX/web/go/module #gosetup
/usr/local/go/bin/go build -o /Users/XXX/Library/Caches/JetBrains/GoLand2023.3/tmp/GoLand/___1tcpGo /Users/XXX/web/go/tcpGo/main.go #gosetup
/Users/XXX/Library/Caches/JetBrains/GoLand2023.3/tmp/GoLand/___1tcpGo
2025/03/07 21:16:47 Waiting for clients
2025/03/07 21:16:57 当前连接数为:0
```

运行客户端代码
```bash
go run ./client/example.go
```

## 代码讲解

### 客户端
客户端代码在client目录下

util/client.go

注：Method:"hello" 与服务端的路由 "hello" 匹配
```
func (c *client) Hello(data []byte) error {
	message := msg.Request{
		Method: "hello",
		Data:   data,
	}
	return c.Send(message)
}
```

### 服务端
服务端开发者只需要关注srv目录的业务代码

在route/path.go里定义路由

```
addPath("hello", Hello)
```

在api/目录里写业务代码

```
func Hello(message msg.Request) msg.ReturnJson {

	log.Println(string(message.Data)) // print hello world

	return msg.ReturnJson{
		Code: common.OK,
		Msg:  "success",
		Data: message.Data,
	}
}
```
message.Data就是客户端传过来的数据




