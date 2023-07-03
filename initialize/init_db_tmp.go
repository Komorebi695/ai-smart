package initialize

import (
	"fmt"
	"strings"
)

func Load(env string, service string) string {
	env = strings.ToLower(env)
	passEnv := map[string]bool{
		"dev":  true,
		"prod": true,
		"test": true,
	}
	var conf string
	if _, ok := passEnv[env]; ok {
		conf = fmt.Sprintf("conf/%s.yaml", env)
	} else {
		conf = fmt.Sprintf("conf/%s.yaml", env)
	}
	return conf
}
