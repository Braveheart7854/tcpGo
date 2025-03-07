/**
 * @author  tongh
 * @date  2024/11/21 15:27
 */
package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"sync"
	"tcpGo/common"
	log2 "tcpGo/util/log"
	"tcpGo/util/msg"
	"time"
)

var UserTempFiles = make(map[string]map[int]TempFile)
var UploadFileMutex = new(sync.Mutex)

type TempFile struct {
	FileName string
	FileSize int
	FileNo   int
	Content  []byte
	Secret   string
}

func UploadFile(message msg.Request) msg.ReturnJson {
	UploadFileMutex.Lock()
	defer UploadFileMutex.Unlock()

	var data struct {
		FileName string
		FileSize int
		FileNo   int
		Content  []byte
		Account  string
		End      bool
		Secret   string
	}
	err := json.Unmarshal(message.Data, &data)
	if err != nil {
		//common.Log(err)
		log2.MainLogger.Error(err.Error())
		return msg.ReturnJson{
			Code: common.FAILED,
			Msg:  "参数错误",
			Data: nil,
		}
	}

	resp := struct {
		FileName string
		FileSize int
		FileNo   int
	}{
		FileName: data.FileName,
		FileSize: data.FileSize,
		FileNo:   data.FileNo,
	}
	b, _ := json.Marshal(resp)

	if data.FileSize > 0 {

		//common.Log(data.FileNo, common.MD5(string(data.Content)), data.Secret)

		if common.MD5(string(data.Content)) != data.Secret {
			common.Log("文件校验失败", data)
			return msg.ReturnJson{
				Code: common.FAILED,
				Msg:  "文件校验失败",
				Data: b,
			}
		}
		if UserTempFiles[data.Account] == nil {
			UserTempFiles[data.Account] = make(map[int]TempFile)
		}
		UserTempFiles[data.Account][data.FileNo] = TempFile{
			FileName: data.FileName,
			FileSize: data.FileSize,
			FileNo:   data.FileNo,
			Content:  data.Content,
			Secret:   data.Secret,
		}
	}

	if data.End {
		var content []byte
		fileMap, ok := UserTempFiles[data.Account]
		if !ok {
			return msg.ReturnJson{
				Code: common.FAILED,
				Msg:  "合并文件失败",
				Data: b,
			}
		}
		for i := 0; i < len(fileMap); i++ {
			content = append(content, fileMap[i].Content...)
		}

		//fmt.Println(common.MD5(string(content)))

		fileName := fmt.Sprintf("./temp/%s", data.FileName)
		err = saveFile(content, fileName)
		if err != nil {
			//common.Log(err)
			log2.MainLogger.Error(err.Error())
			return msg.ReturnJson{
				Code: common.FAILED,
				Msg:  "保存文件失败",
				Data: b,
			}
		}
		delete(UserTempFiles, data.Account)
	}

	//log.Println("UserTempFiles ", len(UserTempFiles[data.Account]))

	log.Println(fmt.Sprintf("%s get data: fileName: %s, fileSize: %d, fileNo: %d, end: %v",
		time.Now().Format("2006-01-02 15:04:05"), data.FileName, data.FileSize, data.FileNo, data.End))
	return msg.ReturnJson{
		Code: common.OK,
		Msg:  "success",
		Data: b,
	}
}

func saveFile(data []byte, filename string) error {
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		file, err := os.Create(filename)
		if err != nil {
			return err
		}
		defer file.Close()
		_, err = file.Write(data)
		if err != nil {
			return err
		}
		err = file.Sync()
		if err != nil {
			return err
		}
		return nil
	}
	return errors.New("文件已存在")
}
