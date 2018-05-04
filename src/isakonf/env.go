package main

import (
	"os"
	"strings"
)

func LoadEnv(dest map[string]interface{}) error {
	for _, e := range os.Environ() {
		if i := strings.Index(e, "="); i >= 0 {
			pair := strings.Split(e, "=")
			dest[strings.ToLower(pair[0])] = pair[1]
		}
	}
	return nil
}
