//go:build !bench

package hw10programoptimization

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetDomainStatAcceptance(t *testing.T) {
	data := `{"Id":1,"Name":"Howard Mendoza","Username":"0Oliver","Email":"aliquid_qui_ea@Browsedrive.gov","Phone":"6-866-899-36-79","Password":"InAQJvsq","Address":"Blackbird Place 25"}
{"Id":2,"Name":"Jesse Vasquez","Username":"qRichardson","Email":"mLynch@broWsecat.com","Phone":"9-373-949-64-00","Password":"SiZLeNSGn","Address":"Fulton Hill 80"}
{"Id":3,"Name":"Clarence Olson","Username":"RachelAdams","Email":"RoseSmith@Browsecat.com","Phone":"988-48-97","Password":"71kuz3gA5w","Address":"Monterey Park 39"}
{"Id":4,"Name":"Gregory Reid","Username":"tButler","Email":"5Moore@Teklist.net","Phone":"520-04-16","Password":"r639qLNu","Address":"Sunfield Park 20"}
{"Id":5,"Name":"Janice Rose","Username":"KeithHart","Email":"nulla@Linktype.com","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}`

	t.Run("find 'com'", func(t *testing.T) {
		result, err := GetDomainStat(bytes.NewBufferString(data), "com")
		require.NoError(t, err)
		require.Equal(t, DomainStat{
			"browsecat.com": 2,
			"linktype.com":  1,
		}, result)
	})

	t.Run("find 'gov'", func(t *testing.T) {
		result, err := GetDomainStat(bytes.NewBufferString(data), "gov")
		require.NoError(t, err)
		require.Equal(t, DomainStat{"browsedrive.gov": 1}, result)
	})

	t.Run("find 'unknown'", func(t *testing.T) {
		result, err := GetDomainStat(bytes.NewBufferString(data), "unknown")
		require.NoError(t, err)
		require.Equal(t, DomainStat{}, result)
	})
}

func TestGetDomainStat(t *testing.T) {
	tests := []struct {
		name    string
		data    string
		domain  string
		want    DomainStat
		wantErr bool
	}{
		{
			"user has domain email",
			`{"Email":"aliquid_qui_ea@Browsedrive.gov"}`,
			"gov",
			DomainStat{"browsedrive.gov": 1},
			false,
		},
		{
			"two users has domain email",
			`{"Email":"aliquid_qui_ea@Browsedrive.gov"}
{"Email":"RoseSmith@browseDrive.gov"}`,
			"gov",
			DomainStat{"browsedrive.gov": 2},
			false,
		},
		{
			"all emails",
			`{"Email":"aliquid_qui_ea@Browsedrive.gov"}
{"Email":"mLynch@broWsecat.com"}`,
			"",
			DomainStat{"browsedrive.gov": 1, "browsecat.com": 1},
			false,
		},
		{
			"user hasn't domain email",
			`{}`,
			"gov",
			DomainStat{},
			false,
		},
		{
			"returns parse error",
			`{"Email":"aliquid_qui_ea@Browsedrive.gov",}`,
			"gov",
			DomainStat{},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := strings.NewReader(tt.data)
			got, err := GetDomainStat(reader, tt.domain)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want, got)
			}
		})
	}
}
