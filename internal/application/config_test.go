package application

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFromEnv(t *testing.T) {
	tests := []struct {
		name string
		env  map[string]string
		cfg  Config
	}{
		{
			name: "Check config from env",
			env: map[string]string{
				"RUN_ADDRESS":            "RUN_ADDRESS_VALUE_FROM_ENV",
				"DATABASE_URI":           "DATABASE_URI_VALUE_FROM_ENV",
				"ACCRUAL_SYSTEM_ADDRESS": "ACCRUAL_SYSTEM_ADDRESS_VALUE_FROM_ENV",
			},
			cfg: Config{
				RunAddress:           "RUN_ADDRESS_VALUE_FROM_ENV",
				DatabaseURI:          "DATABASE_URI_VALUE_FROM_ENV",
				AccrualSystemAddress: "ACCRUAL_SYSTEM_ADDRESS_VALUE_FROM_ENV",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			for k, v := range test.env {
				t.Setenv(k, v)
			}
			cfg := newConfig()
			err := parseEnv(cfg)
			assert.NoError(t, err)
			assert.Equal(t, test.cfg, *cfg)
		})
	}
}
