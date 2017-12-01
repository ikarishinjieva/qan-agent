package util

import (
	"fmt"

	pc "github.com/percona/pmm/proto/config"
)

func GetMySQLConfig(config pc.QAN) ([]string, []string, error) {
	switch config.CollectFrom {
	case "slowlog":
		return makeSlowLogConfig()
	case "perfschema":
		return makePerfSchemaConfig()
	default:
		return nil, nil, fmt.Errorf("invalid CollectFrom: '%s'; expected 'slowlog' or 'perfschema'", config.CollectFrom)
	}
}

func makeSlowLogConfig() ([]string, []string, error) {
	on := []string{
	}
	off := []string{
	}

	on = append(on,
		"SET time_zone='+0:00'",
	)
	return on, off, nil
}

func makePerfSchemaConfig() ([]string, []string, error) {
	return []string{"SET time_zone='+0:00'"}, []string{}, nil
}
