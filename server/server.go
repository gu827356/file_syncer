package server

import (
	"fmt"
	"io"
	"net"

	"file_syncer/common"
)

type Handle interface {
	CanHandle(common.Pkg) bool
	H(common.Pkg, io.Writer) error
}

type SyncServer struct {
	port     int
	handlers []Handle
	fsProxy  *common.FSProxy
}

func NewSyncServer(port int, root string) (*SyncServer, error) {
	proxy, err := common.NewFSProxy(root)
	if err != nil {
		return nil, err
	}

	ser := &SyncServer{
		port:    port,
		fsProxy: proxy,
	}

	registerHandlers(ser)

	return ser, nil
}

func (s *SyncServer) Run() error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", s.port))
	if err != nil {
		return err
	}

	fmt.Printf("start listennig %d\n", s.port)

	for {
		conn, err := lis.Accept()
		if err != nil {
			fmt.Printf("fail to accept remote request: %v\n", err)
		}
		s.handle(conn)
	}
}

func (s *SyncServer) handle(conn net.Conn) {
	defer func() {
		err := recover()
		if err != nil {
			fmt.Printf("fail to handle request: %v\n", err)
		}
	}()

	pkg, err := common.ReadPkg(conn)
	if err != nil {
		fmt.Printf("failt to read pkg: %v\n", err)
		return
	}

	err = s.handlePkg(pkg, conn)
	if err != nil {
		fmt.Printf("fail to handle pkg(%s): %v\n", string(pkg), err)
	}

	err = conn.Close()
	if err != nil {
		fmt.Printf("fail to close conn: %v\n", err)
	}
}

func (s *SyncServer) handlePkg(pkg common.Pkg, out io.Writer) error {
	for _, h := range s.handlers {
		if h.CanHandle(pkg) {
			return h.H(pkg, out)
		}
	}

	return fmt.Errorf("unknown pkg")
}

func registerHandlers(s *SyncServer) {
	s.handlers = append(s.handlers, &pingHandle{})
	s.handlers = append(s.handlers, &downloadHandle{
		fsProxy: s.fsProxy,
	})
	s.handlers = append(s.handlers, &getAllPathEntriesHandle{
		fsProxy: s.fsProxy,
	})
}
