package common

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
)

const (
	GetAllPathEntriesCMD = "getAllPathEntries"
	DownloadCMD          = "download"
)

type Pkg []byte

func (g Pkg) Write(out io.Writer) error {
	lenData := [4]byte{}
	binary.BigEndian.PutUint32(lenData[:], uint32(len(g)))
	_, err := out.Write(lenData[:])
	if err != nil {
		return err
	}

	_, err = out.Write(g)
	return err
}

func ReadPkg(in io.Reader) (Pkg, error) {
	lenData := [4]byte{}
	count, err := in.Read(lenData[:])
	if err != nil {
		return nil, err
	}
	if count != 4 {
		return nil, fmt.Errorf("pkg format error, no date length")
	}

	size := binary.BigEndian.Uint32(lenData[:])

	if size > 1024*1024 {
		return nil, fmt.Errorf("pkg is overload, size = %d", size)
	}

	data := make([]byte, int(size))
	count, err = io.ReadFull(in, data)
	if err != nil {
		return nil, err
	}
	if count != len(data) {
		return nil, fmt.Errorf("pkg form error, there are not enough data[expected=%d, actual=%d]", len(data), count)
	}
	return data, nil
}

func NewDownloadPkg(entry *PathEntry) Pkg {
	d, err := json.Marshal(entry)
	if err != nil {
		panic(err)
	}
	return Pkg(DownloadCMD + ":" + string(d))
}
