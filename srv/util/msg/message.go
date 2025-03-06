/**
 * @author  tongh
 * @date  2022/7/19 5:08 下午
 */
package msg

type Request struct {
	ConnNo string `json:"conn_no"`
	Method string `json:"method"`
	Data   []byte `json:"data"`
}

type Response struct {
	ConnNo     string     `json:"conn_no"`
	ReturnJson ReturnJson `json:"return_json"`
}

type ReturnJson struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data []byte `json:"data"`
}
