package client

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"net"

	"file_syncer/common"
)

type NetClientImpl struct {
	addr string
}

func NewNetClient(addr string) (*NetClientImpl, error) {
	impl := NetClientImpl{
		addr: addr,
	}
	err := impl.Ping()
	if err != nil {
		return nil, fmt.Errorf("fail to ping server: %v", err)
	}

	return &impl, nil
}

func (c *NetClientImpl) GetAllPathEntries() ([]common.PathEntry, error) {
	conn, err := c.dial()
	if err != nil {
		return nil, err
	}

	req := common.Pkg(common.GetAllPathEntriesCMD)
	err = req.Write(conn)
	if err != nil {
		return nil, err
	}

	resp, err := common.ReadPkg(conn)
	if err != nil {
		return nil, fmt.Errorf("fail to get response: %v", err)
	}

	var entries []common.PathEntry
	err = json.Unmarshal(resp, &entries)
	if err != nil {
		return nil, fmt.Errorf("fail to unmarshal response: %v", err)
	}

	return entries, nil
}

func (c *NetClientImpl) DownloadEntry(entry *common.PathEntry, out io.Writer) error {
	req := common.NewDownloadPkg(entry)
	conn, err := c.dial()
	if err != nil {
		return err
	}

	err = req.Write(conn)
	if err != nil {
		return fmt.Errorf("fail to send request: %v", err)
	}

	md := md5.New()
	buf := [4096]byte{}
	for {
		count, err := conn.Read(buf[:])
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("fail to fetch data from server: %v", err)
		}

		md.Write(buf[0:count])
		_, err = out.Write(buf[0:count])
		if err != nil {
			return fmt.Errorf("fail to write data: %v", err)
		}
	}

	if fmt.Sprintf("%x", md.Sum(nil)) != entry.MD5 {
		return fmt.Errorf("md5 error")
	}
	return nil
}

func (c *NetClientImpl) Ping() error {
	conn, err := c.dial()
	if err != nil {
		return err
	}
	defer conn.Close()

	pingPkg := common.Pkg([]byte("ping"))
	err = pingPkg.Write(conn)
	if err != nil {
		return fmt.Errorf("fail to send request: %v", err)
	}

	pongPkg, err := common.ReadPkg(conn)
	if err != nil {
		return fmt.Errorf("fail to fetch response: %v", err)
	}

	if "pong" != string(pongPkg) {
		return fmt.Errorf("unexpected response(%s): %v", string(pongPkg), err)
	}

	return nil
}

func (c *NetClientImpl) dial() (net.Conn, error) {
	conn, err := net.Dial("tcp", c.addr)
	if err != nil {
		return nil, fmt.Errorf("fail to connect to [%s]: %v", c.addr, err)
	}
	return conn, nil
}
