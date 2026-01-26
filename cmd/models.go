package main

type Backup struct {
	Name            string   `json:"Name"`
	Source          string   `json:"Source"`
	Destination     string   `json:"Destination"`
	Retain          int      `json:"Retain"`
	User            string   `json:"User"`
	Verbose         bool     `json:"Verbose"`
	Type            string   `json:"Type"`
	ChangeDir       bool     `json:"ChangeDir"`
	CompressionType string   `json:"CompressionType"`
	Excludes        []string `json:"Excludes"`
}
