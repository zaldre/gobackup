package main

import "strings"

// shellQuote properly quotes a string for use in shell commands
func shellQuote(s string) string {
	// Replace single quotes with '\'' (close quote, escaped quote, open quote)
	quoted := strings.ReplaceAll(s, "'", "'\\''")
	return "'" + quoted + "'"
}
