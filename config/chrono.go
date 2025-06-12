package config

import "time"

func GetNext5AM() time.Time {
	now := time.Now().UTC()
	today5AM := time.Date(now.Year(), now.Month(), now.Day(), 5, 0, 0, 0, time.UTC)
	if now.After(today5AM) {
		today5AM = today5AM.Add(24 * time.Hour)
	}

	return today5AM
}
