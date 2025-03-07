// 通讯协议处理，主要处理封包和解包的过程
package protocol

import (
	"bytes"
	"encoding/binary"
)

const (
	ConstHeader         = "xxxx"
	ConstHeaderLength   = 15
	ConstSaveDataLength = 4
)

// 封包
func Packet(message []byte) []byte {
	return append(append([]byte(ConstHeader), IntToBytes(len(message))...), message...)
}

// 解包
func Unpack(buffer []byte, readerChannel chan []byte) []byte {
	length := len(buffer)

	var i int
	for i = 0; i < length; i = i + 1 {
		if length < i+ConstHeaderLength+ConstSaveDataLength {
			break
		}
		if string(buffer[i:i+ConstHeaderLength]) == ConstHeader {
			messageLength := BytesToInt(buffer[i+ConstHeaderLength : i+ConstHeaderLength+ConstSaveDataLength])
			if length < i+ConstHeaderLength+ConstSaveDataLength+messageLength {
				break
			}
			data := buffer[i+ConstHeaderLength+ConstSaveDataLength : i+ConstHeaderLength+ConstSaveDataLength+messageLength]
			readerChannel <- data

			i += ConstHeaderLength + ConstSaveDataLength + messageLength - 1
		}
	}

	if i == length {
		return make([]byte, 0)
	}
	return buffer[i:]
}

// 整形转换成字节
func IntToBytes(n int) []byte {
	x := int32(n)

	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.BigEndian, x)
	return Rsort(bytesBuffer.Bytes())
}

// 字节转换成整形
func BytesToInt(b []byte) int {
	b = Rsort(b)

	bytesBuffer := bytes.NewBuffer(b)

	var x int32
	binary.Read(bytesBuffer, binary.BigEndian, &x)

	return int(x)
}

// 字节转换成整形
func BytesToInt16(b []byte) int {
	b = Rsort(b)

	bytesBuffer := bytes.NewBuffer(b)

	var x int16
	binary.Read(bytesBuffer, binary.BigEndian, &x)

	return int(x)
}

// 字节转换成整形
func BytesToInt8(b []byte) int {
	b = Rsort(b)

	bytesBuffer := bytes.NewBuffer(b)

	var x int8
	binary.Read(bytesBuffer, binary.BigEndian, &x)

	return int(x)
}

func Int64ToBytes(i int64) []byte {
	var buf = make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(i))
	return buf
}

func BytesToInt64(b []byte) int64 {
	return int64(binary.BigEndian.Uint64(b))
}

func Rsort(b []byte) (list []byte) {
	length := len(b)
	list = make([]byte, length)
	for i := 0; i < length; i++ {
		list[i] = b[length-i-1]
	}
	return
}
