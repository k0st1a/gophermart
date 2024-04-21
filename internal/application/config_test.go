package application

import (
	"os"
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

func TestFromFlags(t *testing.T) {
	resetArgsFun := func() {
		func(args []string) {
			os.Args = args
		}(os.Args)
	}
	defer resetArgsFun()

	tests := []struct {
		name string
		args []string
		cfg  Config
	}{
		{
			name: "Check config from flags",
			args: []string{
				"cmd",
				"-a", "RUN_ADDRESS_VALUE_FROM_FLAG",
				"-d", "DATABASE_URI_VALUE_FROM_FLAG",
				"-r", "ACCRUAL_SYSTEM_ADDRESS_VALUE_FROM_FLAG",
			},
			cfg: Config{
				RunAddress:           "RUN_ADDRESS_VALUE_FROM_FLAG",
				DatabaseURI:          "DATABASE_URI_VALUE_FROM_FLAG",
				AccrualSystemAddress: "ACCRUAL_SYSTEM_ADDRESS_VALUE_FROM_FLAG",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			os.Args = test.args
			cfg := newConfig()
			err := parseFlags(cfg)
			assert.NoError(t, err)
			assert.Equal(t, test.cfg, *cfg)
			resetArgsFun()
		})
	}
}
