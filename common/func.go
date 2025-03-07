/**
 * @author  tongh
 * @date  2022/7/13 4:59 下午
 */
package common

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/gogf/gf/util/gconv"
	"github.com/gogf/gf/util/grand"
	"log"
	"os"
)

func Log(v ...interface{}) {
	log.Println(v...)
	//fmt.Println(v...)
}

func CheckError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}

func GetRandomString(length int) string {
	var str string = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	var sb []byte = gconv.Bytes(str)
	var sb1 []byte
	for i := 0; i < length; i++ {
		n := grand.N(0, 61)
		sb1 = append(sb1, sb[n])
	}
	str1 := gconv.String(sb1)
	return str1
}

// 生成32位MD5
func MD5(text string) string {
	md := md5.New()
	md.Write([]byte(text))
	return hex.EncodeToString(md.Sum(nil))
}
