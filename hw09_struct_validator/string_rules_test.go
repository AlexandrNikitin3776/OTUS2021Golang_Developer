package hw09structvalidator

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestParseStringRule(t *testing.T) {
	tests := []struct {
		name    string
		rule    string
		wantErr bool
	}{
		{"len ok", "len:32", false},
		{"regexp ok", "regexp:\\d", false},
		{"in ok", "in:a", false},
		{"without : ", "len32", true},
		{"invalid rule", "invalid:rule", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseStringRule(tt.rule)

			if tt.wantErr {
				require.ErrorIs(t, err, InvalidStringRule)
			} else {
				require.NoError(t, err)
				require.NotNil(t, got)
			}
		})
	}
}

func TestStringLenRule(t *testing.T) {
	tests := []struct {
		name           string
		controlValue   string
		checkValueOk   string
		checkValueFail string
		wantErr        bool
	}{
		{"ok", "2", "aa", "aaa", false},
		{"invalid arg", "short", "", "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sr := stringRule{}
			got, err := sr.getLenRule(tt.controlValue)

			if tt.wantErr {
				require.ErrorIs(t, err, InvalidStringRule)
			} else {
				require.NoError(t, err)
				require.NotNil(t, got)
				require.NoError(t, got(tt.checkValueOk))
				require.ErrorIs(t, got(tt.checkValueFail), InvalidStringLen)
			}

		})
	}
}

func TestStringRegexpRule(t *testing.T) {
	tests := []struct {
		name           string
		controlValue   string
		checkValueOk   string
		checkValueFail string
		wantErr        bool
	}{
		{"ok", "\\d+", "1234", "aaa", false},
		{"invalid arg", "?d", "", "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sr := stringRule{}
			got, err := sr.getRegexpRule(tt.controlValue)

			if tt.wantErr {
				require.ErrorIs(t, err, InvalidStringRule)
			} else {
				require.NoError(t, err)
				require.NotNil(t, got)
				require.NoError(t, got(tt.checkValueOk))
				require.ErrorIs(t, got(tt.checkValueFail), InvalidStringRegexp)
			}

		})
	}
}

func TestStringInRule(t *testing.T) {
	tests := []struct {
		name           string
		controlValue   string
		checkValueOk   string
		checkValueFail string
		wantErr        bool
	}{
		{"ok", "a,b", "a", "c", false},
		{"empty arg", "", "", "", true},
		{"invalid arg", ",a", "", "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sr := stringRule{}
			got, err := sr.getInRule(tt.controlValue)

			if tt.wantErr {
				require.ErrorIs(t, err, InvalidStringRule)
			} else {
				require.NoError(t, err)
				require.NotNil(t, got)
				require.NoError(t, got(tt.checkValueOk))
				require.ErrorIs(t, got(tt.checkValueFail), InvalidStringIn)
			}

		})
	}
}
