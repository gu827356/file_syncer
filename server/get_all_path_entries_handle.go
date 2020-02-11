package server

import (
	"encoding/json"
	"io"

	"file_syncer/common"
)

type getAllPathEntriesHandle struct {
	fsProxy *common.FSProxy
}

func (h *getAllPathEntriesHandle) CanHandle(pkg common.Pkg) bool {
	return common.GetAllPathEntriesCMD == string(pkg)
}

func (h *getAllPathEntriesHandle) H(_ common.Pkg, out io.Writer) error {
	entries, err := h.fsProxy.GetAllPathEntries()
	if err != nil {
		return err
	}
	data, err := json.Marshal(entries)
	if err != nil {
		return err
	}
	pkg := common.Pkg(data)
	return pkg.Write(out)
}
