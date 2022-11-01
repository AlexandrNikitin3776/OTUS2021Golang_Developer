package hw10programoptimization

import (
	"github.com/stretchr/testify/require"
	"io"
	"strings"
	"testing"
)

func BenchmarkGetUsers(b *testing.B) {
	reader := strings.NewReader("{\"ID\":1,\"Name\":\"Howard Mendoza\",\"Username\":\"0Oliver\",\"Email\":\"aliquid_qui_ea@Browsedrive.gov\",\"Phone\":\"6-866-899-36-79\",\"Password\":\"InAQJvsq\",\"Address\":\"Blackbird Place 25\"}")
	for i := 0; i < b.N; i++ {
		_, err := reader.Seek(0, io.SeekStart)
		require.NoError(b, err)

		_, err = getUsers(reader)
		require.NoError(b, err)
	}
}
