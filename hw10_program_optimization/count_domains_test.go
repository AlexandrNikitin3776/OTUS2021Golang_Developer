package hw10programoptimization

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func BenchmarkCountDomains(b *testing.B) {
	testUsers := users{
		{1, "Howard Mendoza", "0Oliver", "aliquid_qui_ea@Browsedrive.com", "6-866-899-36-79", "InAQJvsq", "Blackbird Place 25"},
	}
	domain := "com"
	for i := 0; i < b.N; i++ {
		_, err := countDomains(testUsers, domain)
		require.NoError(b, err)
	}
}
