package client

import (
	"fmt"
	"io"
	"os"
	"path"
	"strings"
	"sync"
	"sync/atomic"

	"file_syncer/common"

	"github.com/google/uuid"
)

const downloadWorkerSize = 5

type NetClient interface {
	GetAllPathEntries() ([]common.PathEntry, error)
	DownloadEntry(*common.PathEntry, io.Writer) error
}

type SyncClient struct {
	netClient NetClient
	root      string
}

func NewSyncClient(addr, root string) (*SyncClient, error) {
	netClient, err := NewNetClient(addr)
	if err != nil {
		return nil, err
	}

	for strings.HasSuffix(root, "/") {
		root = root[:len(root)-1]
	}

	return &SyncClient{
		netClient: netClient,
		root:      root,
	}, nil
}

func (c *SyncClient) Sync() error {
	err := c.checkRootDir()
	if err != nil {
		return err
	}

	entries, err := c.netClient.GetAllPathEntries()
	if err != nil {
		return fmt.Errorf("fail to get entries: %v", err)
	}

	entries, err = c.diff(entries)
	if err != nil {
		return fmt.Errorf("fail to diff entries: %v", err)
	}

	err = c.insureAllPath(entries)
	if err != nil {
		return fmt.Errorf("fail to create directory: %v", err)
	}

	workingTokens := make(chan int, downloadWorkerSize)
	for i := 0; i < downloadWorkerSize; i++ {
		workingTokens <- 1
	}

	var completed int32
	reportProgress := func() {
		fmt.Printf("progress: %d / %d\n", completed, len(entries))
	}
	reportProgress()

	latch := sync.WaitGroup{}
	for _, entry := range entries {
		latch.Add(1)

		go func(entry common.PathEntry) {
			<-workingTokens

			fmt.Printf("start download %s\n", entry.String())
			err := c.syncFile(&entry)
			if err != nil {
				fmt.Println(err.Error())
			}
			fmt.Printf("download %s successfully\n", entry.String())
			atomic.AddInt32(&completed, 1)
			latch.Done()
			reportProgress()

			workingTokens <- 1
		}(entry)
	}

	latch.Wait()
	return nil
}

func (c *SyncClient) diff(entries []common.PathEntry) ([]common.PathEntry, error) {
	result := make([]common.PathEntry, 0)
	for _, entry := range entries {
		needSync, err := c.isNeedSync(&entry)
		if err != nil {
			return nil, err
		}

		if needSync {
			result = append(result, entry)
		}
	}
	return result, nil
}

func (c *SyncClient) insureAllPath(entries []common.PathEntry) error {
	dirPaths := make([]string, 0)
	addPath := func(p string) {
		dir, _ := path.Split(p)
		for _, dp := range dirPaths {
			if dp == dir {
				return
			}
		}
		dirPaths = append(dirPaths, dir)
	}
	for _, entry := range entries {
		addPath(c.pathOf(&entry))
	}

	for _, p := range dirPaths {
		err := os.MkdirAll(p, os.ModePerm)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *SyncClient) isNeedSync(serverEntry *common.PathEntry) (bool, error) {
	path := c.root + "/" + serverEntry.Path
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return true, nil
		}
		return false, fmt.Errorf("fail to get file info[%s]: %v", path, err)
	}

	if info.IsDir() {
		return false, fmt.Errorf("path[%s] is a file", path)
	}

	md5, err := common.FileMD5(path)
	if err != nil {
		return false, err
	}

	if md5 == serverEntry.MD5 {
		return false, nil
	}

	return serverEntry.LastModifiedTime > info.ModTime().Unix(), nil
}

func (c *SyncClient) syncFile(entry *common.PathEntry) error {
	id, err := uuid.NewRandom()
	if err != nil {
		return err
	}
	tmpFilePath := os.TempDir() + "/" + id.String()
	file, err := os.Create(tmpFilePath)
	if err != nil {
		return fmt.Errorf("fail to create temporary file[%s]: %v", tmpFilePath, err)
	}
	defer func() {
		file.Close()
		os.Remove(tmpFilePath)
	}()

	err = c.netClient.DownloadEntry(entry, file)
	if err != nil {
		return fmt.Errorf("fail to download %s", entry)
	}

	tarFilePath := c.pathOf(entry)
	err = os.Rename(tmpFilePath, tarFilePath)
	if err != nil {
		return fmt.Errorf("fail to move %s to %s: %v", tmpFilePath, tarFilePath, err)
	}
	return nil
}

func (c *SyncClient) pathOf(entry *common.PathEntry) string {
	return c.root + "/" + entry.Path
}

func (c *SyncClient) checkRootDir() error {
	isDir, err := common.IsDirOrNotExist(c.root)
	if err != nil {
		return err
	}
	if !isDir {
		return fmt.Errorf("dir[%s] is not a directory", c.root)
	}
	return nil
}
