package client

import (
	"io"
	"os"
	"testing"

	"file_syncer/common"

	"github.com/stretchr/testify/require"
)

type fsNetClient struct {
	fsProxy *common.FSProxy
}

func (c *fsNetClient) GetAllPathEntries() ([]common.PathEntry, error) {
	return c.fsProxy.GetAllPathEntries()
}

func (c *fsNetClient) DownloadEntry(entry *common.PathEntry, out io.Writer) error {
	return c.fsProxy.CopyEntry(*entry, out)
}

func TestSyncClient(t *testing.T) {
	proxy, err := common.NewFSProxy("../")
	require.NoError(t, err)

	err = os.RemoveAll("/tmp/TestSyncClient")
	require.NoError(t, err)
	defer os.RemoveAll("/tmp/TestSyncClient")

	syncClient := SyncClient{
		netClient: &fsNetClient{fsProxy: proxy,},
		root:      "/tmp/TestSyncClient",
	}

	err = syncClient.Sync()
	require.NoError(t, err)
}
