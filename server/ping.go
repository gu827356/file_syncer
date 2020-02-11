package server

import (
	"io"

	"file_syncer/common"
)

type pingHandle struct {
}

func (h *pingHandle) CanHandle(pkg common.Pkg) bool {
	return "ping" == string(pkg)
}

func (h *pingHandle) H(_ common.Pkg, out io.Writer) error {
	pongPkg := common.Pkg([]byte("pong"))
	return pongPkg.Write(out)
}
