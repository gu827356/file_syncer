package common

import (
	"crypto/md5"
	"fmt"
	"io"
	"os"
)

func FileMD5(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", fmt.Errorf("fail to calculate MD5 of %s: %v", path, err)
	}
	defer file.Close()

	md := md5.New()
	buf := [4096]byte{}
	for {
		count, err := file.Read(buf[:])
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", fmt.Errorf("fail to calculate MD5 of %s: %v", path, err)
		}

		_, err = md.Write(buf[0:count])
		if err != nil {
			return "", fmt.Errorf("fail to calculate MD5 of %s: %v", path, err)
		}
	}

	return fmt.Sprintf("%x", md.Sum(nil)), nil
}
