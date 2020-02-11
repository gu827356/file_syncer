package common

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPkg_Write(t *testing.T) {
	pkg := Pkg([]byte("test content"))
	writeBuf := bytes.NewBuffer(nil)
	err := pkg.Write(writeBuf)
	require.NoError(t, err)

	newPkg, err := ReadPkg(writeBuf)
	require.NoError(t, err)
	require.Equal(t, pkg, newPkg)
}
