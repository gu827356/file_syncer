package common

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type FSProxy struct {
	root string
}

func IsDirOrNotExist(path string) (bool, error) {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return true, nil
		}
		return false, fmt.Errorf("fail to get path info[%s]: %v", path, err)
	}
	return info.IsDir(), nil
}

func NewFSProxy(root string) (*FSProxy, error) {
	root = strings.ReplaceAll(root, `\`, "/")
	for ; strings.HasSuffix(root, "/"); {
		root = root[0 : len(root)-1]
	}
	info, err := os.Stat(root)
	if err != nil {
		return nil, err
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("dir[%s] is no a directory", root)
	}
	return &FSProxy{root: root + "/"}, nil
}

func (p *FSProxy) GetAllPathEntries() ([]PathEntry, error) {
	result := make([]PathEntry, 0)

	lock := sync.Mutex{}
	latch := sync.WaitGroup{}

	err := filepath.Walk(p.root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		latch.Add(1)
		lastModifiedTime := info.ModTime().Unix()
		go func() {
			md5, err := FileMD5(path)
			if err != nil {
				fmt.Printf("fail to calculate files[%s]'s md5: %v\n", path, err)
			}

			path = strings.ReplaceAll(path, `\`, "/")
			if strings.HasPrefix(path, p.root) {
				path = path[len(p.root):]
			}

			lock.Lock()
			result = append(result, PathEntry{
				Path:             path,
				MD5:              md5,
				LastModifiedTime: lastModifiedTime,
			})
			lock.Unlock()
			latch.Done()
		}()

		return nil
	})

	latch.Wait()

	return result, err
}

func (p *FSProxy) CopyEntry(entry PathEntry, dist io.Writer) error {
	filePath := p.root + "/" + entry.Path
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("fail to open file[%s]: %v", filePath, err)
	}
	defer file.Close()

	_, err = io.Copy(dist, file)
	return err
}
