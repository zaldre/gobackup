package main

import (
	"os"
)

func GetEnv(key string, defvalue string) string {
	value := os.Getenv(key)
	if value != "" {
		return value
	}
	return defvalue
}
