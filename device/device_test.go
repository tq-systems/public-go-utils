package device

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseRaucCompatible(t *testing.T) {
	tests := []struct {
		name            string
		input           string
		wantErr         bool
		wantMachine     string
		wantCompatible  string
		wantVersion     string
		wantSpecVersion uint
		wantErrContains string
	}{
		{
			name:            "valid V1 format",
			input:           "em310/production/1/2.5.0",
			wantErr:         false,
			wantMachine:     "em310",
			wantCompatible:  "production",
			wantVersion:     "2.5.0",
			wantSpecVersion: 1,
		},
		{
			name:            "valid V1 with different values",
			input:           "em400/testenv/1/1.0.0-beta",
			wantErr:         false,
			wantMachine:     "em400",
			wantCompatible:  "testenv",
			wantVersion:     "1.0.0-beta",
			wantSpecVersion: 1,
		},
		{
			name:            "empty string",
			input:           "",
			wantErr:         true,
			wantErrContains: "empty",
		},
		{
			name:            "invalid format - too few parts",
			input:           "em310/production/1",
			wantErr:         true,
			wantErrContains: "invalid format",
		},
		{
			name:            "invalid format - too many parts",
			input:           "em310/production/1/2.5.0/extra",
			wantErr:         true,
			wantErrContains: "invalid V1 format",
		},
		{
			name:            "invalid format - slash in machine",
			input:           "em/310/production/1/2.5.0",
			wantErr:         true,
			wantErrContains: "invalid format",
		},
		{
			name:            "invalid format - slash in compatible",
			input:           "em310/prod/uction/1/2.5.0",
			wantErr:         true,
			wantErrContains: "invalid format",
		},
		{
			name:            "invalid format - slash in version",
			input:           "em310/production/1/2.5.0/test",
			wantErr:         true,
			wantErrContains: "invalid V1 format",
		},
		{
			name:            "unsupported spec version 2",
			input:           "em310/production/2/2.5.0",
			wantErr:         true,
			wantSpecVersion: 2,
			wantErrContains: "unsupported spec version",
		},
		{
			name:            "unsupported spec version 0",
			input:           "em310/production/0/2.5.0",
			wantErr:         true,
			wantSpecVersion: 0,
			wantErrContains: "unsupported spec version",
		},
		{
			name:            "invalid spec version - not a number",
			input:           "em310/production/abc/2.5.0",
			wantErr:         true,
			wantErrContains: "invalid format",
		},
		{
			name:            "missing final slash",
			input:           "em310/production/1",
			wantErr:         true,
			wantErrContains: "invalid format",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseRaucCompatible(tt.input)

			// Check that RawString is always set
			assert.Equal(t, tt.input, got.RawString, "RawString should equal input")

			// Check error expectation
			if tt.wantErr {
				assert.Error(t, got.ParsingError, "ParsingError field should be set when parsing fails")
				if tt.wantErrContains != "" {
					assert.Contains(t, got.ParsingError.Error(), tt.wantErrContains, "Error message should contain expected text")
				}
			} else {
				assert.NoError(t, got.ParsingError, "ParsingError field should be nil on success")
				assert.Equal(t, tt.wantMachine, got.BundleMachine, "BundleMachine should match expected value")
				assert.Equal(t, tt.wantCompatible, got.BundleCompatible, "BundleCompatible should match expected value")
				assert.Equal(t, tt.wantVersion, got.BundleVersion, "BundleVersion should match expected value")
				assert.Equal(t, tt.wantSpecVersion, got.SpecVersion, "SpecVersion should match expected value")
			}

			// Check spec version if explicitly set in test (even for errors)
			if tt.wantSpecVersion != 0 && tt.wantErr {
				assert.Equal(t, tt.wantSpecVersion, got.SpecVersion, "SpecVersion should be parsed even when overall parsing fails")
			}
		})
	}
}
