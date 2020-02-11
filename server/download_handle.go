package server

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"file_syncer/common"
)

type downloadHandle struct {
	fsProxy *common.FSProxy
}

func (h *downloadHandle) CanHandle(pkg common.Pkg) bool {
	return strings.HasPrefix(string(pkg), common.DownloadCMD)
}

func (h *downloadHandle) H(pkg common.Pkg, out io.Writer) error {
	jsonData := pkg[len(common.DownloadCMD)+1:]
	entry := common.PathEntry{}
	err := json.Unmarshal(jsonData, &entry)
	if err != nil {
		return fmt.Errorf("fail to unmarshal path entry[%s]: %v", string(jsonData), err)
	}

	return h.fsProxy.CopyEntry(entry, out)
}
