package hw09structvalidator

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestParseIntRule(t *testing.T) {
	tests := []struct {
		name    string
		rule    string
		wantErr bool
	}{
		{"min ok", "min:25", false},
		{"max ok", "max:36", false},
		{"in ok", "in:15", false},
		{"without : ", "min,25", true},
		{"invalid rule", "invalid:rule", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseIntRule(tt.rule)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.NotNil(t, got)
			}
		})
	}
}

func TestIntMinRule(t *testing.T) {
	tests := []struct {
		name           string
		controlValue   string
		checkValueOk   int64
		checkValueFail int64
		wantErr        bool
	}{
		{"ok", "2", 2, 1, false},
		{"invalid arg", "short", 0, 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ir := intRule{}
			got, err := ir.getMinRule(tt.controlValue)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.NotNil(t, got)
				require.NoError(t, got(tt.checkValueOk))
				require.ErrorIs(t, got(tt.checkValueFail), InvalidIntMin)
			}

		})
	}
}
func TestIntMaxRule(t *testing.T) {
	tests := []struct {
		name           string
		controlValue   string
		checkValueOk   int64
		checkValueFail int64
		wantErr        bool
	}{
		{"ok", "2", 2, 3, false},
		{"invalid arg", "short", 0, 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ir := intRule{}
			got, err := ir.getMaxRule(tt.controlValue)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.NotNil(t, got)
				require.NoError(t, got(tt.checkValueOk))
				require.ErrorIs(t, got(tt.checkValueFail), InvalidIntMax)
			}

		})
	}
}

func TestIntInRule(t *testing.T) {
	tests := []struct {
		name           string
		controlValue   string
		checkValueOk   int64
		checkValueFail int64
		wantErr        bool
	}{
		{"ok", "10,12", 10, 11, false},
		{"empty arg", "", 0, 0, true},
		{"invalid arg", ",10", 0, 0, true},
		{"invalid arg", ",a", 0, 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ir := intRule{}
			got, err := ir.getInRule(tt.controlValue)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.NotNil(t, got)
				require.NoError(t, got(tt.checkValueOk))
				require.ErrorIs(t, got(tt.checkValueFail), InvalidIntIn)
			}

		})
	}
}
