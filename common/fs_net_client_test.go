package common

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIsDirOrNotExist(t *testing.T) {
	r, err := IsDirOrNotExist("/tmp/xxxxx")
	require.NoError(t, err)
	require.True(t, r)

	f, err := os.Create("test_file")
	require.NoError(t, err)
	f.Close()

	defer func() { os.Remove("test_file") }()

	err = os.Mkdir("test_dir", os.ModePerm)
	require.NoError(t, err)
	os.Remove("test_dir")

	r, err = IsDirOrNotExist("test_file")
	require.NoError(t, err)
	require.False(t, r)

	r, err = IsDirOrNotExist("test_dir")
	require.NoError(t, err)
	require.True(t, r)
}

func TestFSProxy(t *testing.T) {
	err := os.Mkdir("TestFSProxy", os.ModePerm)
	require.NoError(t, err)
	defer os.RemoveAll("TestFSProxy")

	file, err := os.Create("TestFSProxy/1.data")
	require.NoError(t, err)
	_, err = file.Write([]byte("test content"))
	require.NoError(t, err)
	err = file.Close()
	require.NoError(t, err)

	proxy, err := NewFSProxy("TestFSProxy")
	require.NoError(t, err)

	entries, err := proxy.GetAllPathEntries()
	require.NoError(t, err)
	require.Equal(t, 1, len(entries))
	require.Equal(t, "1.data", entries[0].Path)
	require.NotEqual(t, "", entries[0].MD5)

	buf := bytes.NewBuffer(nil)
	err = proxy.CopyEntry(entries[0], buf)
	require.NoError(t, err)

	require.Equal(t, "test content", string(buf.Bytes()))
}
