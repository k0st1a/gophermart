package application

import (
	"flag"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigFromEnv(t *testing.T) {
	origCommandLineFun := func() {
		func(cl *flag.FlagSet) {
			flag.CommandLine = cl
		}(flag.CommandLine)
	}
	defer origCommandLineFun()

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
			flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

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

func TestConfigFromFlags(t *testing.T) {
	origCommandLineFun := func() {
		func(cl *flag.FlagSet) {
			flag.CommandLine = cl
		}(flag.CommandLine)
	}
	defer origCommandLineFun()

	origArgsFun := func() {
		func(args []string) {
			os.Args = args
		}(os.Args)
	}
	defer origArgsFun()

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
			flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

			os.Args = test.args

			cfg := newConfig()
			err := parseFlags(cfg)
			assert.NoError(t, err)
			assert.Equal(t, test.cfg, *cfg)
			origArgsFun()
			origCommandLineFun()
		})
	}
}

func TestConfig(t *testing.T) {
	origCommandLineFun := func() {
		func(cl *flag.FlagSet) {
			flag.CommandLine = cl
		}(flag.CommandLine)
	}
	defer origCommandLineFun()

	origArgsFun := func() {
		func(args []string) {
			os.Args = args
		}(os.Args)
	}
	defer origArgsFun()

	tests := []struct {
		name string
		env  map[string]string
		args []string
		cfg  Config
	}{
		{
			name: "Check config from env",
			env: map[string]string{
				"RUN_ADDRESS":            "RUN_ADDRESS_VALUE_FROM_ENV",
				"DATABASE_URI":           "DATABASE_URI_VALUE_FROM_ENV",
				"ACCRUAL_SYSTEM_ADDRESS": "ACCRUAL_SYSTEM_ADDRESS_VALUE_FROM_ENV",
			},
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
			flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

			for k, v := range test.env {
				t.Setenv(k, v)
			}
			os.Args = test.args

			cfg, err := collectConfig()
			assert.NoError(t, err)
			assert.Equal(t, test.cfg, *cfg)
			origArgsFun()
			origCommandLineFun()
		})
	}
}
