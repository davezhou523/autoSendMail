package helper

import (
	"archive/zip"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"io"
	"os"
)

type ZipStruct struct {
	Name string
	Body string
	Flag int //1、附件 2:字符
}

func ZipCompress(files []ZipStruct, zipFile string) (string, error) {
	archive, err := os.Create(zipFile)
	if err != nil {
		msg := "创建压缩文件文件失败"
		logx.Errorf("%s:%s", msg, err)
		return "", fmt.Errorf("%s:%s", msg, err)
	}
	defer archive.Close()
	zipW := zip.NewWriter(archive)
	for _, file := range files {
		err := dealCompress(zipW, file)
		if err != nil {
			logx.Errorf("%s", err)
			return "", fmt.Errorf("%s", err)
		}
	}

	err = zipW.Close()
	if err != nil {
		msg := "压缩失败"
		logx.Errorf("%s:%s", msg, err)
		return "", fmt.Errorf("%s:%s", msg, err)
	}
	return zipFile, nil
}

func dealCompress(zipW *zip.Writer, file ZipStruct) error {
	zipFileW, err := zipW.Create(file.Name)
	if err != nil {
		msg := "压缩文件打开文件失败"
		logx.Errorf("%s:%s", msg, err)
		return fmt.Errorf("%s:%s", msg, err)
	}
	//Flag 1、附件 2:字符
	if file.Flag == 1 {
		attachmentFile, err := os.Open(file.Body)
		defer attachmentFile.Close()
		if err != nil {
			msg := "打开附件文件失败"
			logx.Errorf("%s:%s", msg, err)
			return fmt.Errorf("%s:%s", msg, err)
		}
		//附件文件复制到压缩文件
		_, err = io.Copy(zipFileW, attachmentFile)
		if err != nil {
			msg := "复制文件失败"
			logx.Errorf("%s:%s", msg, err)
			return fmt.Errorf("%s:%s", msg, err)
		}
	} else if file.Flag == 2 {
		_, err = zipFileW.Write([]byte(file.Body))
		if err != nil {
			msg := "压缩文件写入文件失败"
			logx.Errorf("%s:%s", msg, err)
			return fmt.Errorf("%s:%s", msg, err)
		}
	}
	return nil
}
