package util

//nolint:gosec
import (
	"crypto/md5"
	"fmt"
	"io"
	"os"
)

func Exists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		return os.IsExist(err)
	}

	return true
}

func CreateMutiDir(filePath string) error {
	if !Exists(filePath) {
		err := os.MkdirAll(filePath, os.ModePerm)
		if err != nil {
			return err
		}

		return err
	}

	return nil
}

//nolint:gosec
func MD5(filePath string) (string, error) {
	hashString := ""
	// 创建MD5哈希对象
	hash := md5.New()

	file, err := os.Open(filePath)
	if err == nil {
		// 复制文件内容到哈希对象
		_, err = io.Copy(hash, file)
		if err == nil {
			// 计算MD5哈希值
			hashBytes := hash.Sum(nil)
			// 将哈希值转换为16进制字符串
			hashString = fmt.Sprintf("%x", hashBytes)
		}
	}

	defer file.Close()

	return hashString, err
}
