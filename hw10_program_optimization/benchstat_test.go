//go:build bench

package hw10programoptimization

import (
	"archive/zip"
	"github.com/stretchr/testify/require"
	"testing"
)

func BenchmarkGetDomainStat(b *testing.B) {
	r, err := zip.OpenReader("testdata/users.dat.zip")
	require.NoError(b, err)
	defer func() {
		err = r.Close()
		require.NoError(b, err)
	}()

	require.Equal(b, 1, len(r.File))

	for i := 0; i < b.N; i++ {
		data, err := r.File[0].Open()
		require.NoError(b, err)

		stat, err := GetDomainStat(data, "biz")
		require.NoError(b, err)
		require.Equal(b, expectedBizStat, stat)
	}
}
